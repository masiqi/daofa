import os
import json
import time
import traceback
import undetected_chromedriver as uc
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from bs4 import BeautifulSoup
from dotenv import load_dotenv, set_key
from selenium.common.exceptions import TimeoutException

load_dotenv()

url = os.getenv('SCRAPE_URL')
cookie = os.getenv('SCRAPE_COOKIE')
last_page = int(os.getenv('LAST_PAGE', '1'))

def scrape_page(driver):
    try:
        WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.CSS_SELECTOR, 'div.tk-quest-item')))
        
        # 点击复选框并等待数据加载
        checkbox = WebDriverWait(driver, 10).until(
            EC.element_to_be_clickable((By.CSS_SELECTOR, 'em.show-answer__checkbox'))
        )
        
        # 检查复选框是否已经被选中
        if 'checked' not in checkbox.get_attribute('class'):
            print("正在点击复选框...")
            driver.execute_script("arguments[0].click();", checkbox)
            print("已点击复选框")
            
            # 等待一段时间
            time.sleep(5)
        else:
            print("复选框已经被选中")
        
        # 使用JavaScript获取页面内容，确保我们获取的是最新的DOM
        page_source = driver.execute_script("return document.body.innerHTML;")
        soup = BeautifulSoup(page_source, 'html.parser')
        
        questions = []
        question_elements = soup.select('div.tk-quest-item')
        print(f"找到 {len(question_elements)} 个问题元素")
        
        for index, element in enumerate(question_elements):
            question_id = element.get('id') or f'question_{index + 1}'
            question_text = element.select_one('.exam-item__cnt')
            question_text = question_text.decode_contents() if question_text else ''
            
            answer_element = element.select_one('.exam-item__opt .item.answer img')
            parse_element = element.select_one('.exam-item__opt .item.parse img')
            
            question_type = element.select_one('.msg-box .left-msg .addi-info:first-child .info-cnt')
            question_type = question_type.text.strip() if question_type else ''
            
            knowledge_points = [item.text.strip() for item in element.select('.knowledge-list .knowledge-item')]
            
            answer_image = answer_element['src'] if answer_element else ''
            parse_image = parse_element['src'] if parse_element else ''
            
            print(f"问题 {question_id}:")
            print(f"  答案图片URL: {answer_image}")
            print(f"  解析图片URL: {parse_image}")
            
            if not answer_image and not parse_image:
                print(f"警告：问题 {question_id} 没有找到答案或解析图片")
            
            questions.append({
                'id': question_id,
                'question': question_text,
                'answerImage': answer_image,
                'parseImage': parse_image,
                'type': question_type,
                'knowledgePoints': knowledge_points
            })
        
        return questions
    except Exception as e:
        print(f"错误：在 scrape_page 函数中发生异常: {str(e)}")
        print("详细错误信息:")
        traceback.print_exc()
        return []

def save_questions(questions):
    with open('temp_questions.jsonl', 'a', encoding='utf-8') as f:
        for question in questions:
            f.write(json.dumps(question, ensure_ascii=False) + '\n')

def generate_final_json():
    questions = []
    with open('temp_questions.jsonl', 'r', encoding='utf-8') as f:
        for line in f:
            questions.append(json.loads(line))
    
    with open('questions.json', 'w', encoding='utf-8') as f:
        json.dump(questions, f, ensure_ascii=False, indent=2)
    
    # 删除临时文件
    os.remove('temp_questions.jsonl')

def main():
    global url, cookie, last_page
    
    print(f"URL: {url}")
    print(f"Cookie: {cookie}")
    print(f"Last Page: {last_page}")

    chrome_options = uc.ChromeOptions()
    chrome_options.add_argument("--headless")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_options.add_argument("--disable-gpu")
    chrome_options.add_argument("--disable-extensions")
    chrome_options.add_argument("--disable-setuid-sandbox")
    chrome_options.add_argument("--remote-debugging-port=9222")
    chrome_options.add_argument("--disable-web-security")
    chrome_options.add_argument("--allow-running-insecure-content")
    chrome_options.add_argument("--window-size=1920,1080")

    try:
        driver = uc.Chrome(options=chrome_options)
    except Exception as e:
        print(f"创建 Chrome WebDriver 时出错: {str(e)}")
        print("详细错误信息:")
        traceback.print_exc()
        return

    try:
        driver.get(url)

        if cookie:
            for cookie_item in cookie.split(';'):
                cookie_item = cookie_item.strip()
                if '=' in cookie_item:
                    name, value = cookie_item.split('=', 1)
                    driver.add_cookie({'name': name, 'value': value})
                else:
                    print(f"警告: 无效的cookie项: {cookie_item}")

            # 打印所有cookie以确认
            print("当前的cookies:")
            for cookie in driver.get_cookies():
                print(f"{cookie['name']}: {cookie['value']}")

        driver.refresh()

        page = last_page

        # 如果是从第一页开始，删除可能存在的旧的临时文件
        if page == 1 and os.path.exists('temp_questions.jsonl'):
            os.remove('temp_questions.jsonl')

        while True:
            print(f"正在抓取第 {page} 页...")
            current_url = driver.current_url
            print(f"当前页面URL: {current_url}")
            
            # 更新 .env 文件
            set_key('.env', 'LAST_PAGE', str(page))
            set_key('.env', 'SCRAPE_URL', current_url)
            
            questions = scrape_page(driver)
            if questions:
                save_questions(questions)
            else:
                print(f"警告：第 {page} 页没有抓取到问题")

            try:
                next_button = WebDriverWait(driver, 10).until(
                    EC.presence_of_element_located((By.CSS_SELECTOR, 'a.pager-item.next-page'))
                )
                if 'disabled' in next_button.get_attribute('class'):
                    break
                next_button.click()
                time.sleep(2)
                page += 1
            except:
                print("没有找到下一页按钮或已到达最后一页")
                break

        print(f"抓取完成,共抓取到第 {page} 页")
        
        # 生成最终的 JSON 文件
        generate_final_json()
        print("已生成最终的 JSON 文件")

    except Exception as e:
        print(f"在主循环中发生错误: {str(e)}")
        print("详细错误信息:")
        traceback.print_exc()

    finally:
        driver.quit()

if __name__ == "__main__":
    main()