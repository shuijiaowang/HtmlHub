// HtmlHub 同步助手脚本：注入页面后提供本地数据同步能力。
(function () {
  const SUBDOMAIN = __SUBDOMAIN__;
  const TOKEN_KEY = "htmlhub_sync_token";
  const REGISTER_URL = __REGISTER_URL__;
  const LOGIN_API = "/api/user/login";
  const SAVE_API = "/api/html/data/save";
  const LOAD_API = "/api/html/data/load?subdomain=" + encodeURIComponent(SUBDOMAIN);

  function getToken() {
    return localStorage.getItem(TOKEN_KEY) || "";
  }
  function request(url, options) {
    return fetch(url, options).then((res) => res.json());
  }

  function readAllLocal() {
    const data = {};
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (!key || key === TOKEN_KEY) continue;
      data[key] = localStorage.getItem(key);
    }
    return data;
  }

  function writeAllLocal(data) {
    Object.keys(data).forEach((key) => localStorage.setItem(key, data[key]));
  }

  function clearAllLocal() {
    const token = getToken();
    localStorage.clear();
    if (token) localStorage.setItem(TOKEN_KEY, token);
  }

  async function doLogin() {
    const email = prompt("请输入登录邮箱");
    if (!email) return;
    const password = prompt("请输入登录密码");
    if (!password) return;
    const res = await request(LOGIN_API, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (res.code !== 0) {
      alert(res.msg || "登录失败");
      return;
    }
    localStorage.setItem(TOKEN_KEY, res.data.token);
    alert("登录成功");
    updateStatus();
  }

  function doRegisterJump() {
    window.open(REGISTER_URL, "_blank");
  }

  function doExport() {
    const blob = new Blob([JSON.stringify(readAllLocal(), null, 2)], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = SUBDOMAIN + "-local-export.json";
    a.click();
    URL.revokeObjectURL(url);
  }

  function doImport() {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = ".json,application/json";
    input.onchange = function () {
      const file = input.files && input.files[0];
      if (!file) return;
      file.text().then((text) => {
        try {
          const parsed = JSON.parse(text);
          clearAllLocal();
          writeAllLocal(parsed || {});
          alert("导入成功，页面将刷新");
          location.reload();
        } catch (e) {
          alert("导入失败：JSON格式错误");
        }
      });
    };
    input.click();
  }

  async function doUpload() {
    const token = getToken();
    if (!token) {
      alert("请先登录");
      return;
    }
    const res = await request(SAVE_API, {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: "Bearer " + token },
      body: JSON.stringify({ subdomain: SUBDOMAIN, dataJson: JSON.stringify(readAllLocal()) }),
    });
    alert(res.code === 0 ? "上传成功" : (res.msg || "上传失败"));
  }

  async function doLoad() {
    const token = getToken();
    if (!token) {
      alert("请先登录");
      return;
    }
    const res = await request(LOAD_API, { method: "GET", headers: { Authorization: "Bearer " + token } });
    if (res.code !== 0) {
      alert(res.msg || "加载失败");
      return;
    }
    try {
      const parsed = JSON.parse(res.data.dataJson || "{}");
      clearAllLocal();
      writeAllLocal(parsed);
      alert("加载成功，页面将刷新");
      location.reload();
    } catch (e) {
      alert("服务端数据格式异常");
    }
  }

  function doLogout() {
    localStorage.removeItem(TOKEN_KEY);
    alert("已退出");
    updateStatus();
  }

  const host = document.createElement("div");
  host.style.cssText = "position:fixed;right:16px;bottom:16px;z-index:2147483647;";
  document.body.appendChild(host);
  const shadow = host.attachShadow({ mode: "open" });

  // 支持最小化，避免遮挡页面内容。
  let minimized = false;

  shadow.innerHTML =
    "" +
    "<style>" +
    ".htmlhub-sync-panel{width:320px;background:#fff;border:1px solid #dcdfe6;border-radius:10px;box-shadow:0 10px 28px rgba(0,0,0,.15);font-size:14px;color:#333;padding:12px;font-family:Arial,sans-serif;}" +
    ".htmlhub-sync-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:6px;}" +
    ".htmlhub-sync-title{font-weight:600;color:#222;font-size:15px;}" +
    ".htmlhub-sync-status{margin-bottom:8px;color:#666;font-size:13px;}" +
    ".htmlhub-sync-actions{display:flex;flex-wrap:wrap;gap:6px;}" +
    ".htmlhub-sync-btn{border:1px solid #dcdfe6;background:#fff;padding:4px 8px;border-radius:6px;cursor:pointer;color:#222;font-size:13px;line-height:1.2;text-shadow:none;-webkit-text-fill-color:#222;}" +
    ".htmlhub-sync-btn:hover{background:#f5f7fa;}" +
    ".htmlhub-sync-mini{padding:2px 6px;font-size:12px;}" +
    ".htmlhub-sync-hidden{display:none;}" +
    "</style>" +
    '<div class="htmlhub-sync-panel">' +
    '<div class="htmlhub-sync-header"><div class="htmlhub-sync-title">HtmlHub 同步助手</div><button class="htmlhub-sync-btn htmlhub-sync-mini" id="h-toggle">最小化</button></div>' +
    '<div id="htmlhub-sync-content">' +
    '<div id="htmlhub-sync-status" class="htmlhub-sync-status">状态：未登录</div>' +
    '<div class="htmlhub-sync-actions">' +
    '<button class="htmlhub-sync-btn" id="h-login">登录</button><button class="htmlhub-sync-btn" id="h-register">注册</button><button class="htmlhub-sync-btn" id="h-logout">退出</button>' +
    '<button class="htmlhub-sync-btn" id="h-export">导出</button><button class="htmlhub-sync-btn" id="h-import">导入</button><button class="htmlhub-sync-btn" id="h-clear">清空</button>' +
    '<button class="htmlhub-sync-btn" id="h-upload">上传</button><button class="htmlhub-sync-btn" id="h-load">加载</button>' +
    "</div></div></div>";

  function getEl(id) {
    return shadow.getElementById(id);
  }
  function updateStatus() {
    const status = getEl("htmlhub-sync-status");
    if (!status) return;
    status.textContent = getToken() ? "状态：已登录" : "状态：未登录";
  }

  function updatePanelVisibility() {
    const content = getEl("htmlhub-sync-content");
    const toggle = getEl("h-toggle");
    if (!content || !toggle) return;
    content.classList.toggle("htmlhub-sync-hidden", minimized);
    toggle.textContent = minimized ? "展开" : "最小化";
  }

  getEl("h-toggle").onclick = function () {
    minimized = !minimized;
    updatePanelVisibility();
  };
  getEl("h-login").onclick = doLogin;
  getEl("h-register").onclick = doRegisterJump;
  getEl("h-logout").onclick = doLogout;
  getEl("h-export").onclick = doExport;
  getEl("h-import").onclick = doImport;
  getEl("h-clear").onclick = function () {
    if (confirm("确定清空当前页面本地数据？")) {
      clearAllLocal();
      alert("已清空");
    }
  };
  getEl("h-upload").onclick = doUpload;
  getEl("h-load").onclick = doLoad;
  updatePanelVisibility();
  updateStatus();
})();
