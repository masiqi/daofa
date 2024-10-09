chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: "sendImageToOCR",
    title: "发送图片到OCR处理队列",
    contexts: ["image"]
  });
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (info.menuItemId === "sendImageToOCR") {
    chrome.tabs.sendMessage(tab.id, {
      action: "sendImageToOCR",
      imageUrl: info.srcUrl
    });
  }
});

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "getApiUrl") {
    chrome.storage.sync.get("apiUrl", (data) => {
      sendResponse({ apiUrl: data.apiUrl || "" });
    });
    return true;
  } else if (request.action === "proxyRequest") {
    fetch(request.url, {
      method: request.method,
      headers: request.headers,
      body: request.body
    })
    .then(response => response.json())
    .then(data => sendResponse({ success: true, data: data }))
    .catch(error => sendResponse({ success: false, error: error.message }));
    return true;
  }
});