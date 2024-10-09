document.addEventListener('DOMContentLoaded', () => {
  const apiUrlInput = document.getElementById('apiUrl');
  const saveBtn = document.getElementById('saveBtn');
  const status = document.getElementById('status');

  // 加载保存的API URL
  chrome.storage.sync.get('apiUrl', (data) => {
    apiUrlInput.value = data.apiUrl || '';
  });

  // 保存API URL
  saveBtn.addEventListener('click', () => {
    const apiUrl = apiUrlInput.value;
    chrome.storage.sync.set({ apiUrl }, () => {
      status.textContent = '选项已保存';
      status.style.opacity = 1;
      setTimeout(() => {
        status.style.opacity = 0;
      }, 2000);
    });
  });
});