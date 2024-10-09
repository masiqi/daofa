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

    const requestData = {
      imageUrl: imageUrl,
      cookie: document.cookie
    };

    chrome.runtime.sendMessage({
      action: "proxyRequest",
      url: apiUrl,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Referer": document.location.href
      },
      body: JSON.stringify(requestData)
    }, (response) => {
      if (response.success) {
        showNotification("成功", "图片已成功发送到OCR处理队列");
      } else {
        showNotification("错误", response.error || "发送失败");
      }
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

// 添加以下代码来启用右键菜单
document.addEventListener('contextmenu', function(e) {
    e.stopPropagation();
}, true);

// 移除可能禁用右键菜单的事件监听器
document.oncontextmenu = null;
document.body.oncontextmenu = null;

// 遍历所有元素,移除可能的contextmenu事件监听器
function removeContextMenuListeners(element) {
    element.oncontextmenu = null;
    for (let child of element.children) {
        removeContextMenuListeners(child);
    }
}
removeContextMenuListeners(document.body);