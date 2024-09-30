// ==UserScript==
// @name         组卷网题目抓取
// @namespace    http://tampermonkey.net/
// @version      0.3
// @description  从组卷网抓取题目并下载,去除无用HTML标签
// @match        https://zujuan.xkw.com/*
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    function cleanHtml(html) {
        // 创建一个临时的div元素
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = html;

        // 移除所有字体相关的标签,包括span
        const tagsToRemove = ['font', 'b', 'strong', 'i', 'em', 'u', 'strike', 's', 'span'];
        tagsToRemove.forEach(tag => {
            const elements = tempDiv.getElementsByTagName(tag);
            for (let i = elements.length - 1; i >= 0; i--) {
                const parent = elements[i].parentNode;
                while (elements[i].firstChild) {
                    parent.insertBefore(elements[i].firstChild, elements[i]);
                }
                parent.removeChild(elements[i]);
            }
        });

        // 移除所有样式属性
        const allElements = tempDiv.getElementsByTagName('*');
        for (let i = 0; i < allElements.length; i++) {
            allElements[i].removeAttribute('style');
            allElements[i].removeAttribute('class');
        }

        // 返回清理后的HTML
        return tempDiv.innerHTML;
    }

    function extractQuestions() {
        const questions = [];
        const questionElements = document.querySelectorAll('body > main > article > section.test-list > div.tk-quest-item');

        questionElements.forEach((element, index) => {
            const questionId = element.id || `question_${index + 1}`;
            const questionText = element.querySelector('.exam-item__cnt')?.innerHTML || '';
            const answerElement = element.querySelector('.exam-item__opt .item.answer img');
            const parseElement = element.querySelector('.exam-item__opt .item.parse img');

            const answerImageSrc = answerElement ? answerElement.src : '';
            const parseImageSrc = parseElement ? parseElement.src : '';

            questions.push({
                id: questionId,
                question: cleanHtml(questionText),
                answerImage: answerImageSrc,
                parseImage: parseImageSrc
            });
        });

        return questions;
    }

    function downloadQuestions(questions) {
        const timestamp = Date.now();
        const filename = `questions_组卷网题目抓取_${timestamp}.json`;
        const jsonStr = JSON.stringify(questions, null, 2);
        const blob = new Blob([jsonStr], { type: 'application/json' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = filename;
        link.click();
        URL.revokeObjectURL(link.href);
        alert(`文件 "${filename}" 已下载。请检查您的下载文件夹。`);

        // 下载图片
        questions.forEach((q, index) => {
            if (q.answerImage) {
                downloadImage(q.answerImage, `answer_${timestamp}_${index + 1}.jpg`);
            }
            if (q.parseImage) {
                downloadImage(q.parseImage, `parse_${timestamp}_${index + 1}.jpg`);
            }
        });
    }

    function downloadImage(url, filename) {
        const link = document.createElement('a');
        link.href = url;
        link.download = filename;
        link.click();
    }

    function addDownloadButton() {
        const button = document.createElement('button');
        button.textContent = '下载题目';
        button.style.cssText = `
            position: fixed;
            top: 10px;
            right: 10px;
            z-index: 9999;
            padding: 10px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        `;
        button.addEventListener('click', () => {
            const questions = extractQuestions();
            downloadQuestions(questions);
        });
        document.body.appendChild(button);
    }

    // 添加下载按钮
    addDownloadButton();
})();
