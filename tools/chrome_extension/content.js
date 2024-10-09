chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "sendImageToOCR") {
    sendImageToOCR(request.imageUrl);
  }
});

function sendImageToOCR(imageUrl) {
  chrome.runtime.sendMessage({ action: "getApiUrl" }, (response) => {
    const apiUrl = response.apiUrl;
    if (!apiUrl) {
      showNotification("错误", "未设置API URL，请在插件选项中设置");
      return;
    }

    fetch(apiUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Referer": document.location.href
      },
      body: JSON.stringify({
        imageUrl: imageUrl,
        cookie: document.cookie
      })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        showNotification("成功", "图片已成功发送到OCR处理队列");
      } else {
        showNotification("错误", data.message || "发送失败");
      }
    })
    .catch(error => {
      showNotification("错误", "发送请求时出错：" + error.message);
    });
  });
}

function showNotification(title, message) {
  const notification = document.createElement("div");
  notification.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    background-color: #f0f0f0;
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 10px;
    z-index: 9999;
  `;
  notification.innerHTML = `<strong>${title}</strong>: ${message}`;
  document.body.appendChild(notification);

  setTimeout(() => {
    document.body.removeChild(notification);
  }, 3000);
}