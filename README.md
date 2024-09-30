# 初中道德与法制知识库系统

## 项目简介

本项目是一个综合性的 OCR 和知识处理系统，旨在帮助初中学生更好地理解道德与法制科目。它结合了先进的 OCR 技术和知识库管理，为学习和教学提供强大支持。

### 项目特点

- 多模型 OCR 支持：集成 GOT-OCR 和 PaddleOCR，自动选择最佳模型
- 智能图像处理：支持多种图像格式，动态预处理优化识别效果
- 灵活的文本输出：支持 LaTeX、纯文本、Markdown 和 HTML 格式
- 用户友好的 API 接口：基于 FastAPI，易于集成和使用
- 详细的日志记录：便于调试和性能优化
- RAG 功能：支持从 Word 文档提取文本，为知识库构建奠定基础
- 环境灵活配置：使用 .env 文件进行系统配置

### 技术亮点

- 使用先进的 AI 模型：如 GOTQwenForCausalLM
- 动态图像分割和缩放：优化处理各种复杂图像
- 智能空格处理：提高文本可读性
- 多模型协同：结合多个 OCR 模型提高识别准确率

## 配置

本项目使用 `.env` 文件进行配置。在开始之前，请复制 `.env.example` 文件并重命名为 `.env`，然后根据您的环境设置相应的值。

## 安装和运行

1. 克隆仓库：
   ```
   git clone [仓库URL]
   cd [项目目录]
   ```

2. 安装依赖：
   ```
   pip install -r requirements.txt
   ```

3. 配置环境变量：
   复制 `.env.example` 为 `.env` 并填写必要的配置信息。

4. 运行 OCR 服务：
   ```
   python tools/ocr/api.py
   ```

5. 运行 RAG 处理（如需要）：
   ```
   python tools/rag/main.py
   ```

## API 使用

OCR API 端点：`POST /ocr`

参数：
- `file`: 上传的图像文件
- `multi_page`: 是否处理多页（默认 False）
- `render`: 是否渲染结果（默认 False）
- `output_format`: 输出格式，可选 "latex"、"plain"、"markdown"、"html"（默认 "latex"）

## 贡献

欢迎对本项目进行贡献！如果您有任何改进意见或发现了 bug，请提交 issue 或 pull request。

## 许可证

[添加许可证信息]

## 致谢

特别感谢所有为本项目做出贡献的开发者和研究人员。本项目使用了多个开源库和模型，包括但不限于 FastAPI、PaddleOCR、PyTorch 等。
