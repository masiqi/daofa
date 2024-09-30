import json
from dotenv import load_dotenv
import os
from openai import OpenAI

# 加载.env文件中的环境变量
load_dotenv()

# 从环境变量中获取配置
api_key = os.getenv('OLLAMA_API_KEY', '')
base_url = os.getenv('OLLAMA_BASE_URL')
model_name = os.getenv('OLLAMA_MODEL_NAME')

# 创建OpenAI客户端
client = OpenAI(api_key=api_key, base_url=base_url)

def process_chunk(chunk):
    prompt = f"""
    请分析以下教材内容，识别并提取知识点和习题。对于每个识别出的项目，请严格按照以下JSON格式输出：

    [
        {{
            "type": "knowledge_point",
            "title": "知识点标题",
            "content": "知识点内容"
        }},
        {{
            "type": "exercise",
            "title": "习题标题",
            "content": "习题内容"
        }}
    ]

    请注意以下几点：
    1. 忽略非知识点和非习题内容，如"关注微信公众号 初高教辅站"获取更多初高中教辅资料"。
    2. 大知识点（父知识点）和子知识点都应该被识别出来。
    3. 对于父知识点，content 字段应为空字符串。
    4. 对于子知识点，请完整保留原文内容，不要进行摘要或总结。
    5. 如果遇到描述相同但名称不同的知识点，请分别列出。
    6. 请确保输出的JSON格式完全正确，可以被直接解析。

    文本内容：
    {chunk}

    JSON输出：
    """
    
    try:
        response = client.chat.completions.create(
            model=model_name,
            messages=[
                {"role": "system", "content": "你是一个专业的教育内容分析助手，擅长从教材中提取结构化的知识点信息。请严格按照指定的JSON格式输出结果。"},
                {"role": "user", "content": prompt}
            ],
            stream=False
        )
        
        output = response.choices[0].message.content.strip()
        
        # 尝试解析JSON输出
        try:
            parsed_output = json.loads(output)
            # 验证输出格式
            if not isinstance(parsed_output, list):
                raise ValueError("输出应该是一个列表")
            for item in parsed_output:
                if not all(key in item for key in ["type", "title", "content"]):
                    raise ValueError("每个项目都应该包含 type, title, 和 content 字段")
            return parsed_output
        except json.JSONDecodeError:
            print(f"无法解析LLM输出为JSON: {output}")
            return None
        except ValueError as e:
            print(f"输出格式不正确: {e}")
            return None
        
    except Exception as e:
        print(f"处理chunk时出错: {e}")
        return None