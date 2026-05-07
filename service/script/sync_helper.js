// HtmlHub 同步助手脚本：注入页面后提供本地数据同步能力。
(function () {
  const SUBDOMAIN = __SUBDOMAIN__;
  const TOKEN_KEY = "htmlhub_sync_token";
  const PANEL_MINIMIZED_KEY = "htmlhub_sync_panel_minimized";
  const AUTO_SYNC_ENABLED_KEY = "htmlhub_sync_auto_enabled";
  const REGISTER_URL = __REGISTER_URL__;
  const HTML_PUBLIC_HOST = __HTML_PUBLIC_HOST__;
  const LOGIN_API = "/api/user/login";
  const SAVE_API = "/api/html/data/save";
  const LOAD_API = "/api/html/data/load?subdomain=" + encodeURIComponent(SUBDOMAIN);
  const HOME_MANAGE_PATH = "/home/manage";

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
      if (!key) continue;
      if (key === TOKEN_KEY) continue;
      if (key === PANEL_MINIMIZED_KEY) continue;
      if (key === AUTO_SYNC_ENABLED_KEY) continue;
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

  async function doUpload(options) {
    const opts = options || {};
    const token = getToken();
    if (!token) {
      if (!opts.silent) alert("请先登录");
      return;
    }
    const res = await request(SAVE_API, {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: "Bearer " + token },
      body: JSON.stringify({ subdomain: SUBDOMAIN, dataJson: JSON.stringify(readAllLocal()) }),
    });
    if (!opts.silent) {
      alert(res.code === 0 ? "上传到云端成功" : (res.msg || "上传到云端失败"));
    }
    return res;
  }

  async function doLoad() {
    const token = getToken();
    if (!token) {
      alert("请先登录");
      return;
    }
    const res = await request(LOAD_API, { method: "GET", headers: { Authorization: "Bearer " + token } });
    if (res.code !== 0) {
      alert(res.msg || "同步到本地失败");
      return;
    }
    try {
      const parsed = JSON.parse(res.data.dataJson || "{}");
      clearAllLocal();
      writeAllLocal(parsed);
      alert("同步到本地成功，页面将刷新");
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

  // 支持最小化，避免遮挡页面内容；状态持久化见 PANEL_MINIMIZED_KEY。
  let minimized = localStorage.getItem(PANEL_MINIMIZED_KEY) === "1";
  let autoSyncEnabled = localStorage.getItem(AUTO_SYNC_ENABLED_KEY) === "1";
  let lastLocalSnapshot = "";
  let autoSyncTimer = null;
  const host = document.createElement("div");
  host.style.cssText = minimized
    ? "position:fixed;right:0;top:50%;transform:translateY(-50%);z-index:2147483647;"
    : "position:fixed;right:16px;bottom:16px;z-index:2147483647;";
  document.body.appendChild(host);
  const shadow = host.attachShadow({ mode: "open" });

  shadow.innerHTML =
    "" +
    "<style>" +
    ".htmlhub-sync-panel{width:320px;background:#fff;border:1px solid #dcdfe6;border-radius:10px;box-shadow:0 10px 28px rgba(0,0,0,.15);font-size:14px;color:#333;padding:12px;font-family:Arial,sans-serif;box-sizing:border-box;}" +
    ".htmlhub-sync-panel--minimized{width:26px;height:52px;padding:0;border-radius:26px 0 0 26px;border-right:none;overflow:hidden;display:flex;flex-direction:column;align-items:center;justify-content:center;}" +
    ".htmlhub-sync-panel--minimized .htmlhub-sync-header{margin:0;flex:1;display:flex;flex-direction:column;align-items:center;justify-content:center;width:100%;min-height:0;}" +
    ".htmlhub-sync-panel--minimized .htmlhub-sync-title{display:none;}" +
    ".htmlhub-sync-panel--minimized .htmlhub-sync-mini{border:none;background:transparent;padding:0;width:100%;flex:1;font-size:20px;line-height:1;display:flex;align-items:center;justify-content:center;-webkit-text-fill-color:#222;}" +
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
    '<button class="htmlhub-sync-btn" id="h-home" title="跳转到主页">主页</button>'+
    '<button class="htmlhub-sync-btn" id="h-login" title="登录账号，用于云端数据同步">登录</button>'+
    '<button class="htmlhub-sync-btn" id="h-register" title="跳转到注册页面创建账号">注册</button>'+
    '<button class="htmlhub-sync-btn" id="h-logout" title="退出当前登录的账号">退出</button>'+
    '<button class="htmlhub-sync-btn" id="h-export" title="导出本地所有数据到JSON文件">导出</button>'+
    '<button class="htmlhub-sync-btn" id="h-import" title="从本地文件导入数据，会覆盖现有全部数据">导入</button>'+
    '<button class="htmlhub-sync-btn" id="h-clear" title="清空当前页面的所有本地存储数据">清空</button>'+
    '<button class="htmlhub-sync-btn" id="h-upload" title="将本地数据同步上传到云端存储">上传到云端</button>'+
    '<button class="htmlhub-sync-btn" id="h-load" title="从云端加载数据，会覆盖本地全部数据">同步到本地</button>'+
    '<button class="htmlhub-sync-btn" id="h-auto" title="开启后每30秒检测本地存储变化并自动上传云端">自动同步：关</button>'+
    "</div></div></div>";

  function getEl(id) {
    return shadow.getElementById(id);
  }
  function updateStatus() {
    const status = getEl("htmlhub-sync-status");
    if (!status) return;
    status.textContent = getToken() ? "状态：已登录" : "状态：未登录";
  }

  function updateAutoSyncBtn() {
    const btn = getEl("h-auto");
    if (!btn) return;
    btn.textContent = autoSyncEnabled ? "自动同步：开" : "自动同步：关";
  }

  function takeLocalSnapshot() {
    try {
      return JSON.stringify(readAllLocal());
    } catch (e) {
      return "";
    }
  }

  async function tickAutoSync() {
    if (!autoSyncEnabled) return;
    if (!getToken()) return;
    const current = takeLocalSnapshot();
    if (!current) return;
    if (lastLocalSnapshot === "") {
      lastLocalSnapshot = current;
      return;
    }
    if (current === lastLocalSnapshot) return;
    lastLocalSnapshot = current;
    try {
      const res = await doUpload({ silent: true });
      if (!res || res.code !== 0) {
        console.warn("[HtmlHub sync] auto upload failed:", res && res.msg ? res.msg : res);
      }
    } catch (e) {
      console.warn("[HtmlHub sync] auto upload exception:", e);
    }
  }

  function startAutoSyncTimer() {
    stopAutoSyncTimer();
    lastLocalSnapshot = takeLocalSnapshot();
    autoSyncTimer = setInterval(tickAutoSync, 30 * 1000);
  }

  function stopAutoSyncTimer() {
    if (autoSyncTimer) {
      clearInterval(autoSyncTimer);
      autoSyncTimer = null;
    }
  }

  function updatePanelVisibility() {
    const content = getEl("htmlhub-sync-content");
    const toggle = getEl("h-toggle");
    const panel = shadow.querySelector(".htmlhub-sync-panel");
    if (!content || !toggle || !panel) return;
    content.classList.toggle("htmlhub-sync-hidden", minimized);
    toggle.textContent = minimized ? "\u2039" : "最小化";
    if (minimized) {
      toggle.setAttribute("title", "展开");
      toggle.setAttribute("aria-label", "展开同步助手");
    } else {
      toggle.removeAttribute("title");
      toggle.removeAttribute("aria-label");
    }
    panel.classList.toggle("htmlhub-sync-panel--minimized", minimized);
    if (minimized) {
      host.style.cssText =
        "position:fixed;right:0;top:50%;transform:translateY(-50%);z-index:2147483647;";
    } else {
      host.style.cssText =
        "position:fixed;right:16px;bottom:16px;top:auto;transform:none;z-index:2147483647;";
    }
  }

  getEl("h-toggle").onclick = function () {
    minimized = !minimized;
    localStorage.setItem(PANEL_MINIMIZED_KEY, minimized ? "1" : "0");
    updatePanelVisibility();
  };
  getEl("h-home").onclick = function () {
    if (HTML_PUBLIC_HOST) {
      const scheme = location.protocol === "https:" ? "https://" : "http://";
      location.href = scheme + HTML_PUBLIC_HOST + HOME_MANAGE_PATH;
      return;
    }
    location.href = HOME_MANAGE_PATH;
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
  getEl("h-upload").onclick = function () {
    return doUpload({ silent: false });
  };
  getEl("h-load").onclick = doLoad;
  getEl("h-auto").onclick = function () {
    autoSyncEnabled = !autoSyncEnabled;
    localStorage.setItem(AUTO_SYNC_ENABLED_KEY, autoSyncEnabled ? "1" : "0");
    updateAutoSyncBtn();
    if (autoSyncEnabled) startAutoSyncTimer();
    else stopAutoSyncTimer();
  };
  updatePanelVisibility();
  updateStatus();
  updateAutoSyncBtn();
  if (autoSyncEnabled) startAutoSyncTimer();
})();
