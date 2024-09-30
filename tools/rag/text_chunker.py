def chunk_text(text, chunk_size=4000, overlap=200):
    chunks = []
    start = 0
    text_length = len(text)

    while start < text_length:
        end = start + chunk_size
        if end > text_length:
            end = text_length
        
        # 找到最近的句子结束点
        while end < text_length and text[end] not in ['.', '!', '?', '\n']:
            end += 1
        
        # 如果找不到句子结束点，就直接在chunk_size处截断
        if end == text_length or end - start > chunk_size + 200:
            end = start + chunk_size

        chunk = text[start:end]
        chunks.append(chunk)

        # 移动起始点，考虑重叠
        start = end - overlap

    return chunks