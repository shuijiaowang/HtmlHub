// HtmlHub 发布模式：访客强制同步发布者数据；发布者若本地有 token 则先拉取带同步助手的页面。
(function () {
  const SUBDOMAIN = __SUBDOMAIN__;
  const TOKEN_KEY = "htmlhub_sync_token";
  const PANEL_MINIMIZED_KEY = "htmlhub_sync_panel_minimized";
  const AUTO_SYNC_ENABLED_KEY = "htmlhub_sync_auto_enabled";
  const APPLY_FLAG = "htmlhub_pub_applied_" + SUBDOMAIN;
  const PUBLISH_LOAD_API =
    "/api/html/data/publish?subdomain=" + encodeURIComponent(SUBDOMAIN);

  try {
    var ingestUrl = new URL(window.location.href);
    var queryToken = ingestUrl.searchParams.get("token");
    if (queryToken) {
      localStorage.setItem(TOKEN_KEY, queryToken);
      ingestUrl.searchParams.delete("token");
      window.history.replaceState(null, "", ingestUrl.toString());
    }
  } catch (e) {}

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
    Object.keys(data || {}).forEach(function (key) {
      localStorage.setItem(key, data[key]);
    });
  }

  function clearDataLocal() {
    const token = localStorage.getItem(TOKEN_KEY) || "";
    const minimized = localStorage.getItem(PANEL_MINIMIZED_KEY);
    const autoSync = localStorage.getItem(AUTO_SYNC_ENABLED_KEY);
    localStorage.clear();
    if (token) localStorage.setItem(TOKEN_KEY, token);
    if (minimized) localStorage.setItem(PANEL_MINIMIZED_KEY, minimized);
    if (autoSync) localStorage.setItem(AUTO_SYNC_ENABLED_KEY, autoSync);
  }

  async function applyPublisherData() {
    try {
      const res = await fetch(PUBLISH_LOAD_API);
      const body = await res.json();
      if (!body || body.code !== 0) return;

      let parsed = {};
      try {
        parsed = JSON.parse((body.data && body.data.dataJson) || "{}");
      } catch (e) {
        parsed = {};
      }

      const before = JSON.stringify(readAllLocal());
      clearDataLocal();
      writeAllLocal(parsed);
      const after = JSON.stringify(readAllLocal());

      if (before === after) return;
      if (sessionStorage.getItem(APPLY_FLAG) === "1") {
        sessionStorage.removeItem(APPLY_FLAG);
        return;
      }
      sessionStorage.setItem(APPLY_FLAG, "1");
      location.reload();
    } catch (e) {
      console.warn("[HtmlHub publish sync]", e);
    }
  }

  async function tryOwnerReload() {
    const token = localStorage.getItem(TOKEN_KEY) || "";
    if (!token) return false;
    try {
      const res = await fetch(window.location.href, {
        headers: { Authorization: "Bearer " + token },
      });
      if (!res.ok) return false;
      const html = await res.text();
      if (html.indexOf("htmlhub-sync-panel") === -1) return false;
      document.open();
      document.write(html);
      document.close();
      return true;
    } catch (e) {
      return false;
    }
  }

  (async function () {
    const isOwnerPage = await tryOwnerReload();
    if (isOwnerPage) return;
    await applyPublisherData();
  })();
})();
