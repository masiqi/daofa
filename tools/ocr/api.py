from fastapi import FastAPI, File, UploadFile, Form
from pydantic import BaseModel
import torch
from PIL import Image, ImageStat, ImageEnhance
import io
import numpy as np
from crop import eval_model, dynamic_preprocess, convert_output_format
from argparse import Namespace
from transformers import AutoTokenizer
from model.GOT_ocr_2_0 import GOTQwenForCausalLM
from GOT.utils.conversation import conv_templates
from GOT.utils.utils import disable_torch_init
import logging
import os
from paddleocr import PaddleOCR
import paddle
import uvicorn

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()

# 初始化模型
disable_torch_init()
model_name = "stepfun-ai/GOT-OCR2_0"
tokenizer = AutoTokenizer.from_pretrained(model_name, trust_remote_code=True)
model = GOTQwenForCausalLM.from_pretrained(model_name, low_cpu_mem_usage=True, device_map='cuda', use_safetensors=True, pad_token_id=151643).eval()
model.to(device='cuda', dtype=torch.bfloat16)

# 确保 PaddlePaddle 使用 GPU
paddle.set_device('gpu')

# 初始化 PaddleOCR
paddle_ocr = PaddleOCR(use_angle_cls=True, lang='ch', use_gpu=True)

def convert_p_to_rgb(image):
    if image.mode == 'P':
        # 获取调色板
        palette = image.getpalette()
        if palette is None:
            print("调色板为空，无法转换。")
            return image

        # 将调色板转换为RGB颜色列表
        rgb_palette = []
        for i in range(0, len(palette), 3):
            rgb_palette.append((palette[i], palette[i+1], palette[i+2]))

        # 手动修改调色板中的颜色
        modified_rgb_palette = []
        for i, color in enumerate(rgb_palette):
            if i < len(palette)/6:
                modified_rgb_palette.append((255, 255, 255))  # 将索引0设置为白色
            else:
                modified_rgb_palette.append((0, 0, 0))  # 其他索引设置为黑色

        # 创建一个新的RGB图像
        rgb_image = Image.new("RGB", image.size)

        # 将调色板应用到图像上
        pixels = image.load()
        rgb_pixels = rgb_image.load()
        for y in range(image.size[1]):
            for x in range(image.size[0]):
                index = pixels[x, y]
                rgb_pixels[x, y] = modified_rgb_palette[index]

        # 保存转换后的图片
        return rgb_image
    else:
        return image

def convert_rgba_to_rgb(image):
    if image.mode == 'RGBA':
        # 创建一个新的RGB图像
        rgb_image = Image.new("RGB", image.size, (255, 255, 255))  # 使用白色背景
        
        # 将RGBA图像粘贴到RGB图像上
        rgb_image.paste(image, mask=image.split()[3])  # 使用alpha通道作为掩码
        
        return rgb_image
    else:
        return image

def process_image(image):
    # 保存原始图像
    logger.info(f"原始图像大小: {image.size}, 模式: {image.mode}")

    if image.mode in ['P', 'RGBA']:
        # 如果是 'P' 模式，转换为 'L' 模式
        if image.mode == 'P':
            image = convert_p_to_rgb(image)
        
        if image.mode == "RGBA":
            image = convert_rgba_to_rgb(image)
    else:
        logger.info(f"图像模式为 {image.mode}，不进行处理")

    logger.info(f"处理后的图像大小: {image.size}, 模式: {image.mode}")
    
    return image

def paddle_ocr_process(image):
    # 使用 PaddleOCR 进行识别
    result = paddle_ocr.ocr(np.array(image), cls=True)
    
    # 打印结果结构以便调试
    logger.info(f"PaddleOCR 原始结果: {result}")
    
    # 修改结果处理逻辑
    text_lines = []
    if isinstance(result, list) and len(result) > 0:
        for item in result[0]:
            if isinstance(item, list) and len(item) == 2:
                text, confidence = item[1]
                text_lines.append(text)
    
    # 直接拼接所有文本行，不添加换行符
    text = ''.join(text_lines)
    logger.info(f"处理后的文本: {text}")
    return text

@app.post("/ocr")
async def ocr(file: UploadFile = File(...), multi_page: bool = Form(False), render: bool = Form(False), output_format: str = Form("latex")):
    try:
        contents = await file.read()
        image = Image.open(io.BytesIO(contents))
        logger.info(f"上传的图像大小: {image.size}, 模式: {image.mode}, 格式: {image.format}")
        
        # 创建一个目录来保存处理后的图像
        os.makedirs("processed_images", exist_ok=True)
        save_path = f"processed_images/processed_{file.filename}"
        
        # 处理图像格式
        image = process_image(image)
        logger.info(f"处理后的图像大小: {image.size}, 模式: {image.mode}")
        
        # 处理图像
        if multi_page:
            sub_images = [image]  # 在API中,我们假设只处理单个上传的图像
        else:
            sub_images = dynamic_preprocess(image)
        logger.info(f"子图像数量: {len(sub_images)}")
        
        args = Namespace(
            model_name=model_name,
            multi_page=multi_page,
            render=render,
            conv_mode="mpt",
            output_format=output_format
        )
        
        # 准备对话模板
        conv = conv_templates[args.conv_mode].copy()
        
        # 调用eval_model函数，传递处理好的参数
        logger.info("开始 GOT_OCR 处理")
        result = eval_model(args, tokenizer, model, sub_images, conv)
        
        logger.info(f"GOT_OCR 原始结果: {result}")
        logger.info(f"GOT_OCR result length: {len(result)}")
        
        # 转换输出格式
        converted_result = convert_output_format(result, output_format)
        ocr_method = "GOT_OCR"
        
        # 修改判断条件，检查是否只包含空白字符或LaTeX表格标记
        if len(converted_result.strip()) == 0 or converted_result.strip() == r"\begin{tabular}{l|l|l|l|l|}\hline\end{tabular}":
            logger.info("GOT_OCR 未检测到有效文本，使用PaddleOCR")
            converted_result = paddle_ocr_process(image)
            logger.info(f"PaddleOCR 处理后结果: {converted_result}")
            logger.info(f"PaddleOCR result length: {len(converted_result)}")
            ocr_method = "PADDLE_OCR"
        
        logger.info(f"最终转换结果: {converted_result}")
        
        return {
            "result": converted_result,
            "ocr_method": ocr_method
        }
    except Exception as e:
        logger.error(f"处理过程中发生错误: {str(e)}", exc_info=True)
        raise

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)