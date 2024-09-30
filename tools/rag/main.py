import os
from dotenv import load_dotenv
from docx_to_text import convert_docx_to_text
# from llm_processor import process_chunk

# 加载.env文件中的环境变量
load_dotenv()

def save_to_txt(text, output_file):
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(text)

def main():
    # 从环境变量中获取文档路径
    docx_path = os.getenv('DOCX_PATH')
    
    if not docx_path:
        print("错误：未设置DOCX_PATH环境变量。请在.env文件中设置DOCX_PATH。")
        return

    # 转换Word文档为文本
    print("正在转换Word文档...")
    full_text = convert_docx_to_text(docx_path)
    
    # 将全文保存到txt文件
    output_file = 'full_text.txt'
    print(f"正在保存全文到 {output_file}...")
    save_to_txt(full_text, output_file)
    
    print(f"处理完成! 全文已保存到 {output_file}")

    """
    # 将文本分块并处理
    print("正在分割和处理文本...")
    all_data = []
    chunk_size = 4000
    overlap = 200
    start = 0
    text_length = len(full_text)

    while start < text_length:
        end = start + chunk_size
        if end > text_length:
            end = text_length
        
        chunk = full_text[start:end]
        print(f"处理第 {len(all_data)+1} 个块")
        processed_data = process_chunk(chunk)
        if processed_data:
            all_data.extend(processed_data)
        else:
            print(f"警告：第 {len(all_data)+1} 个块处理失败，跳过该块")

        start = end - overlap

    # 保存数据到JSON文件
    output_file = 'processed_data.json'
    print(f"正在保存数据到 {output_file}...")
    save_to_json(all_data, output_file)
    
    print(f"处理完成! 数据已保存到 {output_file}")
    print(f"总共处理了 {len(all_data)} 个有效块")
    """

if __name__ == "__main__":
    main()