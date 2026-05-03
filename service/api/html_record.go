package api

import (
	"fmt"
	"htmlhub/config"
	"htmlhub/dao"
	"htmlhub/util"
	"htmlhub/util/response"
	"os"
	"path/filepath"
	"runtime"
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

func (h *HTMLRecordApi) UpdateDescription(c *gin.Context) {
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
		Description string `json:"description" binding:"max=500"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("简介参数错误", c)
		return
	}

	record, err := htmlRecordService.UpdateDescriptionByUserID(uint(userInfo.ID), uint(id), req.Description)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(record, c)
}

func (h *HTMLRecordApi) UpdateHTMLContent(c *gin.Context) {
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
		HTMLContent string `json:"htmlContent" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请提供 HTML 内容", c)
		return
	}

	record, err := htmlRecordService.UpdateHTMLContentByUserID(uint(userInfo.ID), uint(id), req.HTMLContent)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(record, c)
}

func (h *HTMLRecordApi) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	records, total, err := htmlRecordService.AdminList(dao.AdminHTMLRecordQuery{
		Nickname:       strings.TrimSpace(c.Query("nickname")),
		Email:          strings.TrimSpace(c.Query("email")),
		Subdomain:      strings.TrimSpace(c.Query("subdomain")),
		Visibility:     strings.TrimSpace(c.Query("visibility")),
		ApprovalStatus: strings.TrimSpace(c.Query("approvalStatus")),
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"list":     records,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, c)
}

func (h *HTMLRecordApi) AdminDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("记录ID无效", c)
		return
	}
	record, err := htmlRecordService.AdminGet(uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(record, c)
}

func (h *HTMLRecordApi) AdminUpdateApprovalStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("记录ID无效", c)
		return
	}
	var req struct {
		ApprovalStatus string `json:"approvalStatus" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("审核状态参数错误", c)
		return
	}
	record, err := htmlRecordService.AdminUpdateApprovalStatus(uint(id), req.ApprovalStatus)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(record, c)
}

func (h *HTMLRecordApi) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("记录ID无效", c)
		return
	}
	if err := htmlRecordService.AdminDelete(uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"deleted": true}, c)
}

func (h *HTMLRecordApi) AdminUpdateSubdomain(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("记录ID无效", c)
		return
	}
	var req struct {
		Subdomain string `json:"subdomain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("子域名参数错误", c)
		return
	}
	record, err := htmlRecordService.AdminUpdateSubdomain(uint(id), req.Subdomain)
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
		_ = htmlRecordService.IncrementVisitCount(record.ID)
		return
	}

	userID := requestUserID(c)
	if htmlRecordService.CanOwnerAccess(record, userID) {
		renderHTML(c, record.HTMLContent, subdomain)
		_ = htmlRecordService.IncrementVisitCount(record.ID)
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
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return htmlContent
	}
	scriptPath := filepath.Join(filepath.Dir(currentFile), "..", "script", "sync_helper.js")
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return htmlContent
	}

	scriptBody := strings.ReplaceAll(string(content), "__SUBDOMAIN__", strconv.Quote(subdomain))
	scriptBody = strings.ReplaceAll(scriptBody, "__REGISTER_URL__", strconv.Quote(registerURL))
	script := "<script>\n" + scriptBody + "\n</script>"

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
