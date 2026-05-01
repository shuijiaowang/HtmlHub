package api

import (
	"fmt"
	"htmlhub/config"
	"htmlhub/util"
	"htmlhub/util/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type HTMLRecordApi struct{}

func (h *HTMLRecordApi) Upload(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}

	var req struct {
		Subdomain   string `json:"subdomain"`
		FileName    string `json:"fileName" binding:"required,max=255"`
		Description string `json:"description" binding:"max=500"`
		FileSize    int64  `json:"fileSize"`
		HTMLContent string `json:"htmlContent" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误，请检查文件名和HTML内容", c)
		return
	}

	subdomain, err := htmlRecordService.Upload(uint(userInfo.ID), req.Subdomain, req.FileName, req.Description, req.HTMLContent, req.FileSize)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"subdomain": subdomain,
	}, c)
}

func (h *HTMLRecordApi) MyList(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}

	records, err := htmlRecordService.ListByUserID(uint(userInfo.ID))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(records, c)
}

func (h *HTMLRecordApi) Delete(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("记录ID无效", c)
		return
	}

	if err := htmlRecordService.DeleteByUserID(uint(userInfo.ID), uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"deleted": true,
	}, c)
}

func (h *HTMLRecordApi) UpdateVisibility(c *gin.Context) {
	userInfo := util.GetUserInfo(c)
	if userInfo == nil || userInfo.ID <= 0 {
		response.FailWithMessage("未获取到用户信息", c)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("记录ID无效", c)
		return
	}

	var req struct {
		Visibility string `json:"visibility" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("可见性参数错误", c)
		return
	}

	record, err := htmlRecordService.UpdateVisibilityByUserID(uint(userInfo.ID), uint(id), req.Visibility)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(record, c)
}

func (h *HTMLRecordApi) PublicHTML(c *gin.Context) {
	subdomain := extractSubdomain(c.Request.Host)
	record, err := htmlRecordService.GetBySubdomain(subdomain)
	if err != nil {
		c.String(404, "Not Found")
		return
	}

	if htmlRecordService.CanPublicAccess(record) {
		renderHTML(c, record.HTMLContent, subdomain)
		return
	}

	userID := requestUserID(c)
	if htmlRecordService.CanOwnerAccess(record, userID) {
		renderHTML(c, record.HTMLContent, subdomain)
		return
	}

	if strings.TrimSpace(c.GetHeader("Authorization")) == "" {
		c.Data(200, "text/html; charset=utf-8", []byte(accessCheckHTML(portalHomeURL())))
		return
	}

	c.String(404, "Not Found")
}

func requestUserID(c *gin.Context) uint {
	userInfo := util.GetUserInfo(c)
	if userInfo != nil && userInfo.ID > 0 {
		return uint(userInfo.ID)
	}

	token := strings.TrimSpace(c.GetHeader("Authorization"))
	if token == "" {
		return 0
	}
	parts := strings.Fields(token)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		token = parts[1]
	}

	claims, err := util.ParseToken(token)
	if err != nil || claims == nil || claims.ID <= 0 {
		return 0
	}
	return uint(claims.ID)
}

func renderHTML(c *gin.Context, recordHTML, subdomain string) {
	injected := injectSyncScript(recordHTML, subdomain, portalRegisterURL())
	c.Data(200, "text/html; charset=utf-8", []byte(injected))
}

func extractSubdomain(host string) string {
	withoutPort := strings.ToLower(strings.Split(host, ":")[0])
	suffix := strings.TrimSpace(config.AppConfig.App.HtmlPublicHost)
	if suffix != "" {
		suffix = strings.ToLower(suffix)
		if strings.HasSuffix(withoutPort, "."+suffix) {
			prefix := strings.TrimSuffix(withoutPort, "."+suffix)
			if prefix != "" && !strings.Contains(prefix, ".") {
				return prefix
			}
		}
	}
	parts := strings.Split(withoutPort, ".")
	if len(parts) >= 2 && parts[len(parts)-1] == "localhost" {
		return parts[0]
	}
	if len(parts) >= 3 {
		return parts[0]
	}
	return ""
}

func portalRegisterURL() string {
	origin := strings.TrimSpace(config.AppConfig.App.PortalOrigin)
	origin = strings.TrimSuffix(origin, "/")
	if origin == "" {
		origin = "http://localhost:5173"
	}
	return origin + "/register"
}

func portalHomeURL() string {
	origin := strings.TrimSpace(config.AppConfig.App.PortalOrigin)
	origin = strings.TrimSuffix(origin, "/")
	if origin == "" {
		origin = "http://localhost:5173"
	}
	return origin + "/home"
}

func injectSyncScript(htmlContent, subdomain, registerURL string) string {
	script := fmt.Sprintf(`
<script>
(function () {
  const SUBDOMAIN = %q;
  const TOKEN_KEY = 'htmlhub_sync_token';
  const REGISTER_URL = %q;
  const LOGIN_API = '/api/user/login';
  const SAVE_API = '/api/html/data/save';
  const LOAD_API = '/api/html/data/load?subdomain=' + encodeURIComponent(SUBDOMAIN);

  function getToken() { return localStorage.getItem(TOKEN_KEY) || ''; }
  function request(url, options) { return fetch(url, options).then((res) => res.json()); }

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
    const email = prompt('请输入登录邮箱');
    if (!email) return;
    const password = prompt('请输入登录密码');
    if (!password) return;
    const res = await request(LOGIN_API, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({ email, password })
    });
    if (res.code !== 0) { alert(res.msg || '登录失败'); return; }
    localStorage.setItem(TOKEN_KEY, res.data.token);
    alert('登录成功');
    updateStatus();
  }

  function doRegisterJump() { window.open(REGISTER_URL, '_blank'); }

  function doExport() {
    const blob = new Blob([JSON.stringify(readAllLocal(), null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = SUBDOMAIN + '-local-export.json';
    a.click();
    URL.revokeObjectURL(url);
  }

  function doImport() {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json,application/json';
    input.onchange = function () {
      const file = input.files && input.files[0];
      if (!file) return;
      file.text().then((text) => {
        try {
          const parsed = JSON.parse(text);
          clearAllLocal();
          writeAllLocal(parsed || {});
          alert('导入成功，页面将刷新');
          location.reload();
        } catch (e) { alert('导入失败：JSON格式错误'); }
      });
    };
    input.click();
  }

  async function doUpload() {
    const token = getToken();
    if (!token) { alert('请先登录'); return; }
    const res = await request(SAVE_API, {
      method: 'POST',
      headers: {'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token},
      body: JSON.stringify({ subdomain: SUBDOMAIN, dataJson: JSON.stringify(readAllLocal()) })
    });
    alert(res.code === 0 ? '上传成功' : (res.msg || '上传失败'));
  }

  async function doLoad() {
    const token = getToken();
    if (!token) { alert('请先登录'); return; }
    const res = await request(LOAD_API, { method: 'GET', headers: { 'Authorization': 'Bearer ' + token } });
    if (res.code !== 0) { alert(res.msg || '加载失败'); return; }
    try {
      const parsed = JSON.parse(res.data.dataJson || '{}');
      clearAllLocal();
      writeAllLocal(parsed);
      alert('加载成功，页面将刷新');
      location.reload();
    } catch (e) { alert('服务端数据格式异常'); }
  }

  function doLogout() { localStorage.removeItem(TOKEN_KEY); alert('已退出'); updateStatus(); }

  const host = document.createElement('div');
  host.style.cssText = 'position:fixed;right:16px;bottom:16px;z-index:2147483647;';
  document.body.appendChild(host);
  const shadow = host.attachShadow({ mode: 'open' });

  shadow.innerHTML = ''
    + '<style>'
    + '.htmlhub-sync-panel{width:320px;background:#fff;border:1px solid #dcdfe6;border-radius:10px;box-shadow:0 10px 28px rgba(0,0,0,.15);font-size:14px;color:#333;padding:12px;font-family:Arial,sans-serif;}'
    + '.htmlhub-sync-title{font-weight:600;margin-bottom:6px;color:#222;font-size:15px;}'
    + '.htmlhub-sync-status{margin-bottom:8px;color:#666;font-size:13px;}'
    + '.htmlhub-sync-actions{display:flex;flex-wrap:wrap;gap:6px;}'
    + '.htmlhub-sync-btn{border:1px solid #dcdfe6;background:#fff;padding:4px 8px;border-radius:6px;cursor:pointer;color:#222;font-size:13px;line-height:1.2;text-shadow:none;-webkit-text-fill-color:#222;}'
    + '.htmlhub-sync-btn:hover{background:#f5f7fa;}'
    + '</style>'
    + '<div class="htmlhub-sync-panel">'
    + '<div class="htmlhub-sync-title">HtmlHub 同步助手</div>'
    + '<div id="htmlhub-sync-status" class="htmlhub-sync-status">状态：未登录</div>'
    + '<div class="htmlhub-sync-actions">'
    + '<button class="htmlhub-sync-btn" id="h-login">登录</button><button class="htmlhub-sync-btn" id="h-register">注册</button><button class="htmlhub-sync-btn" id="h-logout">退出</button>'
    + '<button class="htmlhub-sync-btn" id="h-export">导出</button><button class="htmlhub-sync-btn" id="h-import">导入</button><button class="htmlhub-sync-btn" id="h-clear">清空</button>'
    + '<button class="htmlhub-sync-btn" id="h-upload">上传</button><button class="htmlhub-sync-btn" id="h-load">加载</button>'
    + '</div></div>';

  function getEl(id) { return shadow.getElementById(id); }
  function updateStatus() {
    const status = getEl('htmlhub-sync-status');
    if (!status) return;
    status.textContent = getToken() ? '状态：已登录' : '状态：未登录';
  }

  getEl('h-login').onclick = doLogin;
  getEl('h-register').onclick = doRegisterJump;
  getEl('h-logout').onclick = doLogout;
  getEl('h-export').onclick = doExport;
  getEl('h-import').onclick = doImport;
  getEl('h-clear').onclick = function () {
    if (confirm('确定清空当前页面本地数据？')) { clearAllLocal(); alert('已清空'); }
  };
  getEl('h-upload').onclick = doUpload;
  getEl('h-load').onclick = doLoad;
  updateStatus();
})();
</script>
`, subdomain, registerURL)

	lower := strings.ToLower(htmlContent)
	idx := strings.LastIndex(lower, "</body>")
	if idx >= 0 {
		return htmlContent[:idx] + script + htmlContent[idx:]
	}
	return htmlContent + script
}

func accessCheckHTML(homeURL string) string {
	return fmt.Sprintf(`
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>页面未公开</title>
</head>
<body>
  <script>
  (async function () {
    const TOKEN_KEY = 'htmlhub_sync_token';
    const HOME_URL = %q;
    const currentUrl = new URL(window.location.href);
    const queryToken = currentUrl.searchParams.get('token');
    if (queryToken) {
      localStorage.setItem(TOKEN_KEY, queryToken);
      currentUrl.searchParams.delete('token');
      window.history.replaceState(null, '', currentUrl.toString());
    }

    const token = localStorage.getItem(TOKEN_KEY) || '';
    if (!token) {
      document.body.innerHTML = '<div style="padding:20px">该页面未公开，请从控制台打开</div>';
      setTimeout(() => { location.href = HOME_URL; }, 2000);
      return;
    }

    try {
      const res = await fetch(currentUrl.toString(), {
        headers: { 'Authorization': 'Bearer ' + token }
      });
      if (!res.ok) {
        localStorage.removeItem(TOKEN_KEY);
        document.body.innerHTML = '<div style="padding:20px">登录已过期或无权访问，请重新打开</div>';
        setTimeout(() => { location.href = HOME_URL; }, 1500);
        return;
      }

      const html = await res.text();
      document.open();
      document.write(html);
      document.close();
    } catch (e) {
      document.body.textContent = '加载失败';
      return;
    }
  })();
  </script>
</body>
</html>`, homeURL)
}
