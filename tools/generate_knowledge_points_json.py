import json
from bs4 import BeautifulSoup
import re

def parse_html_to_json(html_content):
    soup = BeautifulSoup(html_content, 'html.parser')
    root_ul = soup.find('ul', class_='tree-top-ul')
    
    def parse_node(element, level=1):
        result = []
        for li in element.find_all('li', recursive=False):
            node = {
                "id": int(li.get('tree-id')),
                "name": li.find('a', class_='tree-anchor').text.strip(),
                "is_leaf": 'tree-leaf' in li.get('class', []),
                "level": level
            }
            if not node['is_leaf']:
                child_ul = li.find('ul', class_='tree-ul')
                if child_ul:
                    node['children'] = parse_node(child_ul, level + 1)
            result.append(node)
        return result

    return parse_node(root_ul)

def main():
    # 读取HTML文件
    with open('knowledge_points.html', 'r', encoding='utf-8') as f:
        html_content = f.read()

    # 解析HTML并生成JSON结构
    knowledge_points = parse_html_to_json(html_content)

    # 将JSON结构写入文件
    with open('../backend/sql/knowledge_points.json', 'w', encoding='utf-8') as f:
        json.dump(knowledge_points, f, ensure_ascii=False, indent=2)

    print("JSON文件已生成: ../backend/sql/knowledge_points.json")

if __name__ == "__main__":
    main()