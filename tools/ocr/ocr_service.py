import cv2
import numpy as np
import torch
from GOT.models.got_ocr import GOT_OCR
from GOT.utils.utils import AttnLabelConverter

class OCRService:
    def __init__(self, dict_path, model_path):
        self.device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        
        # 初始化标签转换器
        with open(dict_path, 'r', encoding='utf-8') as file:
            vocabulary = [char.strip() for char in file]
        self.converter = AttnLabelConverter(vocabulary)
        
        # 初始化模型
        self.model = GOT_OCR(
            img_size=(32, 100),
            num_class=len(self.converter.character),
            hidden_size=512,
            num_block=4,
            num_head=8,
            encoder_mode='cnn',
            decoder_mode='transformer',
            loss_mode='ctc'
        )
        self.model = self.model.to(self.device)
        self.model.load_state_dict(torch.load(model_path, map_location=self.device))
        self.model.eval()

    def process_image(self, image):
        # 将图片转换为OpenCV格式
        nparr = np.frombuffer(image, np.uint8)
        img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        
        # 图像预处理
        img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
        img = cv2.resize(img, (100, 32))
        img = img.astype(np.float32) / 255.
        img = torch.from_numpy(img).permute(2, 0, 1).unsqueeze(0)
        img = img.to(self.device)

        # 执行OCR
        with torch.no_grad():
            pred = self.model(img)
        
        # 解码预测结果
        _, pred_index = pred.max(2)
        pred_str = self.converter.decode(pred_index)

        return pred_str[0]