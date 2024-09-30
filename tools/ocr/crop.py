import argparse
from transformers import AutoTokenizer, AutoModelForCausalLM
import torch
import os
from GOT.utils.conversation import conv_templates, SeparatorStyle
from GOT.utils.utils import disable_torch_init
from transformers import CLIPVisionModel, CLIPImageProcessor, StoppingCriteria
from GOT.model import *
from GOT.utils.utils import KeywordsStoppingCriteria

from PIL import Image

import os
import requests
from PIL import Image
from io import BytesIO
from GOT.model.plug.blip_process import BlipImageEvalProcessor
from transformers import TextStreamer
from natsort import natsorted
import glob
import re
from bs4 import BeautifulSoup
import markdown
import logging

DEFAULT_IMAGE_TOKEN = "<image>"
DEFAULT_IMAGE_PATCH_TOKEN = '<imgpad>'
DEFAULT_IM_START_TOKEN = '<img>'
DEFAULT_IM_END_TOKEN = '</img>'

logger = logging.getLogger(__name__)

def load_image(image_file):
    if image_file.startswith('http') or image_file.startswith('https'):
        response = requests.get(image_file)
        image = Image.open(BytesIO(response.content)).convert('RGB')
    else:
        image = Image.open(image_file).convert('RGB')
    return image

def find_closest_aspect_ratio(aspect_ratio, target_ratios, width, height, image_size):
    best_ratio_diff = float('inf')
    best_ratio = (1, 1)
    area = width * height
    for ratio in target_ratios:
        target_aspect_ratio = ratio[0] / ratio[1]
        ratio_diff = abs(aspect_ratio - target_aspect_ratio)
        if ratio_diff < best_ratio_diff:
            best_ratio_diff = ratio_diff
            best_ratio = ratio
        elif ratio_diff == best_ratio_diff:
            if area > 0.5 * image_size * image_size * ratio[0] * ratio[1]:
                best_ratio = ratio
    # print(f'width: {width}, height: {height}, best_ratio: {best_ratio}')
    return best_ratio


def dynamic_preprocess(image, min_num=1, max_num=6, image_size=1024, use_thumbnail=True):
    orig_width, orig_height = image.size
    aspect_ratio = orig_width / orig_height

    # calculate the existing image aspect ratio
    target_ratios = set(
        (i, j) for n in range(min_num, max_num + 1) for i in range(1, n + 1) for j in range(1, n + 1) if
        i * j <= max_num and i * j >= min_num)
    # print(target_ratios)
    target_ratios = sorted(target_ratios, key=lambda x: x[0] * x[1])

    # find the closest aspect ratio to the target
    target_aspect_ratio = find_closest_aspect_ratio(
        aspect_ratio, target_ratios, orig_width, orig_height, image_size)

    # print(target_aspect_ratio)
    # calculate the target width and height
    target_width = image_size * target_aspect_ratio[0]
    target_height = image_size * target_aspect_ratio[1]
    blocks = target_aspect_ratio[0] * target_aspect_ratio[1]

    # resize the image
    resized_img = image.resize((target_width, target_height))
    processed_images = []
    for i in range(blocks):
        box = (
            (i % (target_width // image_size)) * image_size,
            (i // (target_width // image_size)) * image_size,
            ((i % (target_width // image_size)) + 1) * image_size,
            ((i // (target_width // image_size)) + 1) * image_size
        )
        # split the image
        split_img = resized_img.crop(box)
        processed_images.append(split_img)
    assert len(processed_images) == blocks
    if use_thumbnail and len(processed_images) != 1:
        thumbnail_img = image.resize((image_size, image_size))
        processed_images.append(thumbnail_img)
    return processed_images

def eval_model(args, tokenizer, model, sub_images, conv):
    image_processor_high = BlipImageEvalProcessor(image_size=1024)
    use_im_start_end = True
    image_token_len = 256

    image_list = []
    for p in sub_images:
        image_1 = p.copy()
        # 确保图像是 RGB 模式
        if image_1.mode != 'RGB':
            image_1 = image_1.convert('RGB')
        image_tensor_1 = image_processor_high(image_1)
        image_list.append(image_tensor_1)

    image_list = torch.stack(image_list)
    ll = len(sub_images)

    logger.info(f'====new images batch size======: {image_list.shape}')

    if args.multi_page:
        qs = 'OCR with format across multi pages: '
    else:
        qs = 'OCR with format upon the patch reference: '

    if use_im_start_end:
        qs = DEFAULT_IM_START_TOKEN + DEFAULT_IMAGE_PATCH_TOKEN*image_token_len*ll + DEFAULT_IM_END_TOKEN + '\n' + qs 
    else:
        qs = DEFAULT_IMAGE_TOKEN + '\n' + qs

    conv.append_message(conv.roles[0], qs)
    conv.append_message(conv.roles[1], None)
    prompt = conv.get_prompt()

    inputs = tokenizer([prompt])
    input_ids = torch.as_tensor(inputs.input_ids).cuda()

    stop_str = conv.sep if conv.sep_style != SeparatorStyle.TWO else conv.sep2
    keywords = [stop_str]
    stopping_criteria = KeywordsStoppingCriteria(keywords, tokenizer, input_ids)
    streamer = TextStreamer(tokenizer, skip_prompt=True, skip_special_tokens=True)

    with torch.autocast("cuda", dtype=torch.bfloat16):
        output_ids = model.generate(
            input_ids,
            images=[(image_list.half().cuda(), image_list.half().cuda())],
            do_sample=False,
            num_beams=1,
            streamer=streamer,
            max_new_tokens=4096,
            stopping_criteria=[stopping_criteria]
        )
    
    outputs = tokenizer.decode(output_ids[0, input_ids.shape[1]:]).strip()

    if outputs.endswith(stop_str):
        outputs = outputs[:-len(stop_str)]
    outputs = outputs.strip()

    if args.render:
        # 渲染逻辑保持不变
        print('==============rendering===============')
        # ... (渲染代码)

    return outputs  # 返回OCR结果

def convert_output_format(latex_output, output_format):
    if output_format == "latex":
        return latex_output
    elif output_format == "plain":
        return latex_to_plain(latex_output)
    elif output_format == "markdown":
        return latex_to_markdown(latex_output)
    elif output_format == "html":
        return latex_to_html(latex_output)
    else:
        raise ValueError(f"Unsupported output format: {output_format}")

def latex_to_plain(latex):
    # 移除所有LaTeX命令和环境
    plain = re.sub(r'\\[a-zA-Z]+(\[.*?\])?({.*?})?', '', latex)
    plain = re.sub(r'\\begin{.*?}|\\end{.*?}', '', plain)
    # 移除 {|c|c|c|} 和 {l} 等格式
    plain = re.sub(r'\{[|lcr]+\}', '', plain)
    # 移除剩余的花括号、反斜杠和竖线
    plain = re.sub(r'[{}\\|]', '', plain)
    # 将 & 替换为空格
    plain = plain.replace('&', ' ')
    # 将多个空白字符替换为单个空格
    plain = re.sub(r'\s+', ' ', plain)
    # 移除行首和行尾的空白字符
    plain = re.sub(r'^\s+|\s+$', '', plain, flags=re.MULTILINE)
    
    # 智能处理空格
    def smart_space(match):
        left, space, right = match.groups()
        # 定义需要保留空格的情况
        keep_space = (
            left.isdigit() or right.isdigit() or  # 数字和其他字符之间保留空格
            (left in '，。、：；！？') or (right in '，。、：；！？') or  # 标点符号前后保留空格
            (left.isalnum() and right.isalnum()) or  # 字母数字混合词之间保留空格
            (left in ['第', '共']) or (right in ['名', '个'])  # 特定词语前后保留空格
        )
        return left + (' ' if keep_space else '') + right

    # 应用智能空格处理
    plain = re.sub(r'(\S)(\s)(\S)', smart_space, plain)
    
    return plain.strip()

def latex_to_markdown(latex):
    # 将LaTeX表格转换为Markdown表格
    markdown = re.sub(r'\\begin{tabular}.*?\\end{tabular}', lambda m: latex_table_to_markdown(m.group(0)), latex, flags=re.DOTALL)
    # 移除其他LaTeX命令
    markdown = re.sub(r'\\[a-zA-Z]+(\[.*?\])?({.*?})?', '', markdown)
    markdown = re.sub(r'[{}\\]', '', markdown)
    return markdown.strip()

def latex_table_to_markdown(latex_table):
    rows = re.findall(r'\\hline(.*?)(?=\\hline|\Z)', latex_table, re.DOTALL)
    markdown_rows = []
    for row in rows:
        cells = re.split(r'\s*&\s*', row.strip())
        markdown_row = '| ' + ' | '.join(cells) + ' |'
        markdown_rows.append(markdown_row)
    
    if len(markdown_rows) > 1:
        header_separator = '| ' + ' | '.join(['---'] * len(cells)) + ' |'
        markdown_rows.insert(1, header_separator)
    
    return '\n'.join(markdown_rows)

def latex_to_html(latex):
    # 首先转换为Markdown
    markdown_text = latex_to_markdown(latex)
    # 然后将Markdown转换为HTML
    html = markdown.markdown(markdown_text, extensions=['tables'])
    return html

from PIL import Image, ImageStat
import numpy as np

def process_image(image, save_path):
    # 保存原始图像
    original_save_path = save_path.replace("processed_", "original_")
    image.save(original_save_path)
    logger.info(f"保存原始图像到 {original_save_path}")

    img_array = np.array(image)
    logger.info(f"原始图像形状: {img_array.shape}")
    logger.info(f"图像模式: {image.mode}")

    if image.mode == 'RGB':
        logger.info("处理 RGB 图像")
        # 检查是否需要从 BGR 转换为 RGB
        if img_array[0, 0, 0] == 0 and img_array[0, 0, 2] > 0:
            logger.info("将 BGR 转换为 RGB")
            img_array = img_array[:, :, ::-1]
            image = Image.fromarray(img_array)
    elif image.mode == 'RGBA':
        logger.info("将 RGBA 转换为 RGB")
        image = image.convert('RGB')
    elif image.mode in ['L', 'P']:
        logger.info(f"处理 {image.mode} 模式图像")
        
        if image.mode == 'P':
            # 检查调色板
            palette = image.getpalette()
            if palette is None or all(v == 0 for v in palette):
                logger.info("调色板为空，手动创建调色板")
                palette = [i // 3 for i in range(768)]
                image.putpalette(palette)
            
            # 转换为灰度图像
            image = image.convert('L')
        
        # 计算图像的平均亮度
        stat = ImageStat.Stat(image)
        mean_brightness = stat.mean[0]
        
        # 判断是否需要反转图像
        if mean_brightness < 128:
            logger.info("检测到深色背景，进行反转")
            image = Image.eval(image, lambda x: 255 - x)
        else:
            logger.info("检测到浅色背景，保持原样")
    else:
        logger.warning(f"意外的图像模式: {image.mode}")
    
    logger.info(f"处理后的图像形状: {np.array(image).shape}")
    logger.info(f"处理后的图像模式: {image.mode}")
    
    # 保存处理后的图像
    image.save(save_path)
    logger.info(f"保存处理后的图像到 {save_path}")
    
    return image