<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <div class="records-title-row">
          <h3>我的上传记录</h3>
          <button class="text-btn" @click="loadRecords">刷新</button>
        </div>
        <p class="total-visits">总访问次数：{{ totalVisitCount }}</p>
        <p v-if="records.length === 0" class="empty">暂无记录</p>
        <div v-else class="record-list">
          <article v-for="item in records" :key="item.id" class="record-item">
            <div class="record-head">
              <strong>{{ item.fileName }}</strong>
              <span>{{ formatSize(item.fileSize) }}</span>
            </div>
            <p class="record-route">
              访问链接：{{ item.subdomain }}.{{ htmlPublicHost }}
              <a class="share-link" :href="buildShareUrl(item)" target="_blank">点击打开</a>
            </p>
            <p class="record-desc">{{ "简介："+item.description || '无简介' }}</p>
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
      <el-input v-model="htmlDraft" type="textarea" :rows="16" placeholder="粘贴完整 HTML" />
      <template #footer>
        <el-button @click="htmlDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="htmlSaving" @click="saveHtmlContent">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import {
  deleteHtmlRecord,
  getMyHtmlList,
  updateHtmlContent,
  updateHtmlDescription,
  updateHtmlVisibility
} from '@/api/html'
import { useUserStore } from '@/stores/user'

/** 子域名 HTML 访问用 host，本地与线上由 Vite 环境变量区分 */
const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'

const records = ref([])
const userStore = useUserStore()

const descDialogVisible = ref(false)
const htmlDialogVisible = ref(false)
const editingRecord = ref(null)
const descDraft = ref('')
const htmlDraft = ref('')
const descSaving = ref(false)
const htmlSaving = ref(false)

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
  if (!subdomain) return '#'
  const url = new URL(`${window.location.protocol}//${subdomain}.${htmlPublicHost}`)
  if ((item.visibility !== 'public' || item.approvalStatus !== 'approved') && userStore.token) {
    url.searchParams.set('token', userStore.token)
  }
  return url.toString()
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
  await loadRecords()
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: var(--color-background);
}

.content {
  padding: 32px 24px;
  max-width: 980px;
  margin: 0 auto;
}

.card {
  background: #fff;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  padding: 18px;
  margin-bottom: 18px;
}

.text-btn {
  border: none;
  background: transparent;
  color: hsla(160, 100%, 37%, 1);
  cursor: pointer;
  font-size: 14px;
}

.records-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.total-visits {
  margin: 4px 0 14px;
  color: #555;
  font-size: 14px;
}

.record-list {
  display: grid;
  gap: 10px;
}

.record-item {
  border: 1px solid #edf0f2;
  border-radius: 8px;
  padding: 12px;
  background: #fafbfc;
}

.record-head {
  display: flex;
  justify-content: space-between;
}

.record-desc {
  margin: 8px 0;
  color: #555;
}

.record-route {
  margin: 8px 0 0;
  color: #666;
  font-size: 14px;
}

.share-link {
  margin-left: 10px;
}

.record-meta {
  display: flex;
  gap: 16px;
  color: #888;
  font-size: 13px;
}

.record-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.danger-btn {
  color: #d93025;
}

.empty {
  color: #888;
}

.dialog-hint {
  margin: 0 0 12px;
  color: #666;
  font-size: 13px;
  line-height: 1.5;
}
</style>
