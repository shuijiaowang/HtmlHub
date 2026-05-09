<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <div class="records-title-row">
          <div class="title-wrap">
            <h3 class="title">我的上传记录</h3>
            <p class="subtitle">管理你的 HTML：编辑简介、更新内容、切换可见性、查看访问统计。</p>
          </div>
          <button class="text-btn refresh" @click="loadRecords">刷新</button>
        </div>
        <div class="stats-row">
          <div class="stat">
            <div class="stat-label">总访问次数</div>
            <div class="stat-value">{{ totalVisitCount }}</div>
          </div>
          <div class="stat">
            <div class="stat-label">记录数</div>
            <div class="stat-value">{{ records.length }}</div>
          </div>
        </div>
        <p v-if="records.length === 0" class="empty">暂无记录</p>
        <div v-else class="record-list">
          <article v-for="item in records" :key="item.id" class="record-item">
            <div class="record-head">
              <strong>{{ item.fileName }}</strong>
              <span>{{ formatSize(item.fileSize) }}</span>
            </div>
            <p class="record-route">
              访问链接：{{ item.subdomain }}.{{ htmlPublicHost }}
              <button
                type="button"
                class="share-link"
                title="在新标签页打开"
                @click="openShare(item, $event)"
                @auxclick="openShare(item, $event)"
              >
                点击打开
              </button>
            </p>
            <p class="record-desc">{{ item.description ? `简介：${item.description}` : '简介：无简介' }}</p>
            <div class="record-meta">
              <span>审核：{{ formatApprovalStatus(item.approvalStatus) }}</span>
              <span>可见性：{{ formatVisibility(item.visibility) }}</span>
              <span>访问次数：{{ item.visitCount || 0 }}</span>
              <span>创建时间：{{ formatDate(item.createdAt) }}</span>
            </div>
            <div class="record-actions">
              <button class="text-btn" type="button" @click="openDescDialog(item)">编辑简介</button>
              <button class="text-btn" type="button" @click="openHtmlDialog(item)">更新 HTML</button>
              <button class="text-btn" @click="toggleVisibility(item)">
                切换为{{ item.visibility === 'public' ? '私密' : '公开' }}
              </button>
              <button class="text-btn danger-btn" @click="removeRecord(item)">删除</button>
            </div>
          </article>
        </div>
      </section>
    </main>

    <el-dialog v-model="descDialogVisible" title="编辑简介" width="520px" destroy-on-close @closed="editingRecord = null">
      <el-input v-model="descDraft" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="简介（选填）" />
      <template #footer>
        <el-button @click="descDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="descSaving" @click="saveDescription">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="htmlDialogVisible" title="更新 HTML" width="720px" destroy-on-close @closed="editingRecord = null">
      <p class="dialog-hint">保存后仅替换 HTML 正文，审核状态将重置为「未审核」。须为完整文档（含 &lt;html&gt;、&lt;head&gt;、&lt;body&gt; 及闭合标签）。</p>
      <div class="html-dialog-editor">
        <div class="html-toolbar" role="toolbar" aria-label="HTML 文本操作">
          <button type="button" class="toolbar-btn" @click="copyHtmlDraft">复制</button>
          <button type="button" class="toolbar-btn" @click="pasteHtmlDraft">粘贴</button>
          <button type="button" class="toolbar-btn" @click="clearHtmlDraft">清空</button>
        </div>
        <el-input
          ref="htmlInputRef"
          v-model="htmlDraft"
          type="textarea"
          :rows="16"
          placeholder="粘贴完整 HTML"
          @paste="onHtmlDraftNativePaste"
        />
      </div>
      <template #footer>
        <el-button @click="htmlDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="htmlSaving" @click="saveHtmlContent">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, unref } from 'vue'
import {
  deleteHtmlRecord,
  getMyHtmlList,
  updateHtmlContent,
  updateHtmlDescription,
  updateHtmlVisibility
} from '@/api/html'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

/** 子域名 HTML 访问用 host，本地与线上由 Vite 环境变量区分 */
const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'

const records = ref([])
const userStore = useUserStore()
const isLoggedIn = computed(() => !!userStore.token)

const descDialogVisible = ref(false)
const htmlDialogVisible = ref(false)
const editingRecord = ref(null)
const descDraft = ref('')
const htmlDraft = ref('')
const htmlInputRef = ref(null)
const descSaving = ref(false)
const htmlSaving = ref(false)

const resolveHtmlDraftTextarea = () => {
  const inst = htmlInputRef.value
  if (!inst) return null
  const exposed = inst.textarea
  if (exposed != null) return unref(exposed)
  return inst.$el?.querySelector?.('textarea') ?? null
}

const focusTextareaStartAndReveal = (ta) => {
  if (!ta || !document.contains(ta)) return
  ta.scrollTop = 0
  ta.setSelectionRange(0, 0)
  try {
    ta.focus({ preventScroll: true })
  } catch {
    ta.focus()
  }
  ta.scrollIntoView({ block: 'nearest', inline: 'nearest' })
}

const afterHtmlDraftUpdate = () => {
  nextTick(() => {
    requestAnimationFrame(() => {
      focusTextareaStartAndReveal(resolveHtmlDraftTextarea())
    })
  })
}

const copyHtmlDraft = async () => {
  const text = htmlDraft.value
  if (!text) {
    ElMessage.warning('暂无内容可复制')
    return
  }
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败：无法写入剪贴板，请检查权限或浏览器限制')
  }
}

const pasteHtmlDraft = async () => {
  let text
  try {
    text = await navigator.clipboard.readText()
  } catch {
    ElMessage.error('粘贴失败：无法读取剪贴板，请检查权限或在内置浏览器中重试')
    return
  }
  htmlDraft.value = text
  afterHtmlDraftUpdate()
}

const clearHtmlDraft = () => {
  htmlDraft.value = ''
  afterHtmlDraftUpdate()
}

const onHtmlDraftNativePaste = (e) => {
  const el = e.target?.tagName === 'TEXTAREA' ? e.target : resolveHtmlDraftTextarea()
  const cd = e.clipboardData
  if (!cd || !el) return
  const text = cd.getData('text/plain')
  if (text === '') return
  e.preventDefault()
  const start = el.selectionStart
  const end = el.selectionEnd
  const prev = htmlDraft.value
  htmlDraft.value = prev.slice(0, start) + text + prev.slice(end)
  afterHtmlDraftUpdate()
}

const openDescDialog = (item) => {
  editingRecord.value = item
  descDraft.value = item.description || ''
  descDialogVisible.value = true
}

const openHtmlDialog = (item) => {
  editingRecord.value = item
  htmlDraft.value = item.htmlContent || ''
  htmlDialogVisible.value = true
}

const saveDescription = async () => {
  const item = editingRecord.value
  if (!item?.id) return
  descSaving.value = true
  try {
    const res = await updateHtmlDescription(item.id, descDraft.value)
    const next = res.data?.description
    if (next !== undefined) item.description = next
    else item.description = descDraft.value
    descDialogVisible.value = false
  } finally {
    descSaving.value = false
  }
}

const saveHtmlContent = async () => {
  const item = editingRecord.value
  if (!item?.id) return
  htmlSaving.value = true
  try {
    const res = await updateHtmlContent(item.id, htmlDraft.value)
    if (res.data?.htmlContent !== undefined) item.htmlContent = res.data.htmlContent
    if (res.data?.approvalStatus !== undefined) item.approvalStatus = res.data.approvalStatus
    htmlDialogVisible.value = false
  } finally {
    htmlSaving.value = false
  }
}

const totalVisitCount = computed(() => {
  return records.value.reduce((total, item) => total + (Number(item.visitCount) || 0), 0)
})

const loadRecords = async () => {
  if (!userStore.token) {
    ElMessage.warning('请先登录后查看个人记录')
    userStore.openAuthDialog('login')
    return
  }
  const res = await getMyHtmlList()
  records.value = Array.isArray(res.data) ? res.data : []
}

const formatSize = (size) => {
  if (!size || size <= 0) return '0 B'
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / (1024 * 1024)).toFixed(2)} MB`
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (Number.isNaN(date.getTime())) return dateStr
  return date.toLocaleString()
}

const formatVisibility = (visibility) => {
  return visibility === 'public' ? '公开' : '私密'
}

const toggleVisibility = async (item) => {
  const nextVisibility = item.visibility === 'public' ? 'private' : 'public'
  const res = await updateHtmlVisibility(item.id, nextVisibility)
  item.visibility = res.data?.visibility || nextVisibility
}

const removeRecord = async (item) => {
  if (!window.confirm(`确定删除「${item.fileName}」吗？`)) return
  await deleteHtmlRecord(item.id)
  await loadRecords()
}

const buildShareUrl = (item) => {
  const subdomain = item?.subdomain
  if (!subdomain) return ''
  const url = new URL(`${window.location.protocol}//${subdomain}.${htmlPublicHost}`)
  // 已登录时始终附带 token，便于子域注入的同步脚本将 JWT 写入 localStorage（与 access_check 一致）
  if (userStore.token) {
    url.searchParams.set('token', userStore.token)
  }
  return url.toString()
}

/** 不把带 token 的 URL 写在 href 上，避免右键「复制链接」泄露 JWT；仅在真实点击时打开 */
const openShare = (item, event) => {
  if (event) {
    if (event.type === 'auxclick' && event.button !== 1) return
    if (event.type === 'click' && event.button !== 0) return
    event.preventDefault()
  }
  const url = buildShareUrl(item)
  if (!url) return
  window.open(url, '_blank', 'noopener,noreferrer')
}

const formatApprovalStatus = (status) => {
  const statusMap = {
    pending: '未审核',
    approved: '通过',
    rejected: '拒绝'
  }
  return statusMap[status] || '未审核'
}

onMounted(async () => {
  if (userStore.token) {
    await loadRecords()
    return
  }

  window.setTimeout(() => {
    if (!userStore.token) userStore.openAuthDialog('login')
  }, 3000)
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: var(--color-background);
}

.content {
  padding: 32px 24px;
  max-width: 1100px;
  margin: 0 auto;
}

.card {
  background: var(--hh-surface-solid);
  border: 1px solid var(--hh-border);
  border-radius: var(--hh-radius-md);
  padding: 18px 18px 14px;
  margin-bottom: 18px;
  box-shadow: var(--hh-shadow-md);
}

.text-btn {
  border: 1px solid transparent;
  background: transparent;
  color: color-mix(in srgb, var(--hh-brand) 86%, #000 0%);
  cursor: pointer;
  font-size: 14px;
  padding: 6px 7px;
  border-radius: 8px;
  transition: background-color 0.15s ease, border-color 0.15s ease, transform 0.08s ease;
}

.text-btn:hover {
  background: rgb(var(--hh-brand-rgb) / 0.10);
  border-color: rgb(var(--hh-brand-rgb) / 0.18);
  text-decoration: none;
}

.text-btn:active {
  transform: translateY(1px);
}

.refresh {
  white-space: nowrap;
}

.records-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.title-wrap {
  min-width: 0;
}

.title {
  margin: 0;
  font-size: 18px;
  letter-spacing: -0.02em;
}

.subtitle {
  margin: 6px 0 0;
  color: var(--hh-text-3);
  font-size: 13px;
  line-height: 1.5;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin: 14px 0 16px;
}

.stat {
  border: 1px solid var(--hh-border);
  border-radius: var(--hh-radius-md);
  padding: 12px 12px;
  background: color-mix(in srgb, var(--hh-surface-solid) 92%, rgb(var(--hh-brand-rgb) / 0.06) 8%);
}

.stat-label {
  font-size: 12px;
  color: var(--hh-text-3);
}

.stat-value {
  margin-top: 2px;
  font-weight: 800;
  font-size: 20px;
  letter-spacing: -0.02em;
}

.record-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.record-item {
  border: 1px solid var(--hh-border);
  border-radius: var(--hh-radius-md);
  padding: 14px 14px 12px;
  background: color-mix(in srgb, var(--hh-surface-solid) 94%, #000 0%);
  box-shadow: var(--hh-shadow-sm);
  transition: transform 0.10s ease, box-shadow 0.15s ease, border-color 0.15s ease;
}

.record-item:hover {
  transform: translateY(-1px);
  border-color: rgb(var(--hh-brand-rgb) / 0.22);
  box-shadow: 0 14px 30px rgba(15, 23, 42, 0.10);
}

.record-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: baseline;
}

.record-head strong {
  font-size: 15px;
  line-height: 1.3;
  word-break: break-word;
}

.record-head span {
  color: var(--hh-text-3);
  font-size: 12px;
  flex: 0 0 auto;
}

.record-desc {
  margin: 10px 0 0;
  color: var(--hh-text-2);
  font-size: 13px;
  line-height: 1.6;
  word-break: break-word;
}

.record-route {
  margin: 8px 0 0;
  color: var(--hh-text-2);
  font-size: 13px;
  word-break: break-word;
}

.share-link {
  margin-left: 10px;
  font-weight: 650;
  font: inherit;
  color: color-mix(in srgb, var(--hh-brand) 86%, #000 0%);
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.share-link:hover {
  color: var(--hh-brand);
}

.record-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 14px;
  margin-top: 10px;
  color: var(--hh-text-3);
  font-size: 13px;
}

.record-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.danger-btn {
  color: #d93025;
}

.empty {
  color: var(--hh-text-3);
  padding: 14px 0 6px;
}

.dialog-hint {
  margin: 0 0 12px;
  color: var(--hh-text-2);
  font-size: 13px;
  line-height: 1.5;
}

.html-dialog-editor {
  display: grid;
  gap: 6px;
}

.html-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.toolbar-btn {
  padding: 6px 12px;
  border-radius: var(--hh-radius-sm);
  border: 1px solid var(--hh-border);
  background: color-mix(in srgb, var(--hh-surface-solid) 92%, rgb(var(--hh-brand-rgb) / 0.06) 8%);
  color: var(--hh-text);
  font-size: 13px;
  font-family: inherit;
  cursor: pointer;
}

.toolbar-btn:hover {
  border-color: rgb(var(--hh-brand-rgb) / 0.22);
  background: rgb(var(--hh-brand-rgb) / 0.08);
}

@media (max-width: 900px) {
  .record-list {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .content {
    padding: 18px 12px;
  }

  .stats-row {
    grid-template-columns: 1fr;
  }
}
</style>
