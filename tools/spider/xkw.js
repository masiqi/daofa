// ==UserScript==
// @name         组卷网题目抓取
// @namespace    http://tampermonkey.net/
// @version      0.4
// @description  抓取组卷网题目,包括题目类型和相关知识点
// @match        https://www.zujuan.com/*
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    function cleanHtml(html) {
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = html;
        const allElements = tempDiv.getElementsByTagName("*");
        for (let i = 0; i < allElements.length; i++) {
            allElements[i].removeAttribute('style');
            allElements[i].removeAttribute('class');
        }
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

            // 获取题目类型
            const questionType = element.querySelector('.msg-box .left-msg .addi-info:first-child .info-cnt')?.textContent.trim() || '';

            // 获取相关知识点
            const knowledgePoints = Array.from(element.querySelectorAll('.knowledge-list .knowledge-item')).map(item => item.textContent.trim());

            const answerImageSrc = answerElement ? answerElement.src : '';
            const parseImageSrc = parseElement ? parseElement.src : '';

            questions.push({
                id: questionId,
                question: cleanHtml(questionText),
                answerImage: answerImageSrc,
                parseImage: parseImageSrc,
                type: questionType,
                knowledgePoints: knowledgePoints
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
