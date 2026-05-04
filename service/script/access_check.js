// HtmlHub 访问校验脚本模板（用于未公开页面的登录态校验）。
(async function () {
  const TOKEN_KEY = "htmlhub_sync_token";
  const HOME_URL = __HOME_URL__;
  const currentUrl = new URL(window.location.href);
  const queryToken = currentUrl.searchParams.get("token");
  if (queryToken) {
    localStorage.setItem(TOKEN_KEY, queryToken);
    currentUrl.searchParams.delete("token");
    window.history.replaceState(null, "", currentUrl.toString());
  }

  const token = localStorage.getItem(TOKEN_KEY) || "";
  if (!token) {
    document.body.innerHTML = '<div style="padding:20px">该页面未公开，如果是开发者请从控制台跳转打开</div>';
    setTimeout(() => {
      location.href = HOME_URL;
    }, 2000);
    return;
  }

  try {
    const res = await fetch(currentUrl.toString(), {
      headers: { Authorization: "Bearer " + token },
    });
    if (!res.ok) {
      localStorage.removeItem(TOKEN_KEY);
      document.body.innerHTML = '<div style="padding:20px">您的登录已过期或无权访问，请从控制台重新打开</div>';
      setTimeout(() => {
        location.href = HOME_URL;
      }, 1500);
      return;
    }

    const html = await res.text();
    document.open();
    document.write(html);
    document.close();
  } catch (e) {
    document.body.textContent = "加载失败";
    return;
  }
})();
