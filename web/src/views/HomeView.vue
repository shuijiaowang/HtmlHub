
<template>
  <div class="home-page">
    <header class="top-nav">
      <div class="left-links">
        <router-link to="/home">首页</router-link>
      </div>
      <div class="right-actions">
        <template v-if="isLoggedIn">
          <span class="nickname">{{ displayNickname }}</span>
          <button class="text-btn" @click="logout">退出</button>
        </template>
        <template v-else>
          <router-link to="/register">注册</router-link>
          <router-link to="/login">登录</router-link>
        </template>
      </div>
    </header>

    <main class="content">
      <h1>HtmlHub 核心 Demo</h1>
      <p class="subtitle">上传 HTML 文件文本并查看你自己的上传记录。</p>

      <section class="card">
        <div class="upload-title-row">
          <h3>上传 HTML</h3>
          <button class="text-btn" type="button" @click="uploadPanelOpen = !uploadPanelOpen">
            {{ uploadPanelOpen ? '收起' : '展开' }}
          </button>
        </div>
        <p v-if="uploadPanelOpen" class="tips">填写页面前缀后，系统会自动追加 2-3 位随机后缀（示例：todo-a9 / todo-x7k）；不填写则系统自动生成完整随机子域名。</p>
        <form v-if="uploadPanelOpen" class="upload-form" @submit.prevent="submitUpload">
          <label>
            页面前缀（可选）
            <input v-model="uploadForm.subdomain" placeholder="例如：todo">
          </label>
          <label>
            文件名
            <input v-model="uploadForm.fileName" placeholder="例如：landing-page.html" required>
          </label>
          <label>
            简介
            <input v-model="uploadForm.description" maxlength="500" placeholder="可选：写一点说明">
          </label>
          <label>
            选择本地 HTML 文件（可选）
            <input type="file" accept=".html,text/html" @change="onSelectFile">
          </label>
          <label>
            HTML 文本
            <textarea
              v-model="uploadForm.htmlContent"
              rows="10"
              placeholder="<html>...</html>"
              required
            />
          </label>
          <button class="primary-btn" type="submit">上传并写入数据库</button>
        </form>
      </section>

      <section class="card">
        <div class="records-title-row">
          <h3>我的上传记录</h3>
          <button class="text-btn" @click="loadRecords">刷新</button>
        </div>
        <p v-if="records.length === 0" class="empty">暂无记录</p>
        <div v-else class="record-list">
          <article v-for="item in records" :key="item.id" class="record-item">
            <div class="record-head">
              <strong>{{ item.fileName }}</strong>
              <span>{{ formatSize(item.fileSize) }}</span>
            </div>
            <p class="record-route">
              访问域名：{{ item.subdomain }}.{{ htmlPublicHost }}
              <a class="share-link" :href="buildShareUrl(item.subdomain)" target="_blank">打开分享链接</a>
            </p>
            <p class="record-desc">{{ item.description || '无简介' }}</p>
            <div class="record-meta">
              <span>审核：{{ item.isApproved ? '已通过' : '未审核' }}</span>
              <span>创建时间：{{ formatDate(item.createdAt) }}</span>
            </div>
          </article>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { uploadHtml, getMyHtmlList } from '@/api/html'

const userStore = useUserStore()

/** 子域名 HTML 访问用 host，本地与线上由 Vite 环境变量区分 */
const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'

const isLoggedIn = computed(() => !!userStore.token)
const displayNickname = computed(() => userStore.userInfo.nickname || '已登录用户')
const records = ref([])
const uploadPanelOpen = ref(false)
const uploadForm = reactive({
  subdomain: '',
  fileName: '',
  description: '',
  fileSize: 0,
  htmlContent: ''
})

const logout = async () => {
  await userStore.logout()
}

const onSelectFile = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  if (!uploadForm.fileName) {
    uploadForm.fileName = file.name
  }
  uploadForm.fileSize = file.size
  uploadForm.htmlContent = await file.text()
}

const submitUpload = async () => {
  if (!uploadForm.fileName.trim() || !uploadForm.htmlContent.trim()) {
    ElMessage.warning('请填写文件名和HTML文本')
    return
  }

  await uploadHtml({
    subdomain: uploadForm.subdomain.trim(),
    fileName: uploadForm.fileName.trim(),
    description: uploadForm.description.trim(),
    fileSize: uploadForm.fileSize,
    htmlContent: uploadForm.htmlContent
  })

  ElMessage.success('上传成功')
  uploadForm.subdomain = ''
  uploadForm.fileName = ''
  uploadForm.description = ''
  uploadForm.fileSize = 0
  uploadForm.htmlContent = ''
  await loadRecords()
}

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

const buildShareUrl = (subdomain) => {
  if (!subdomain) return '#'
  return `${window.location.protocol}//${subdomain}.${htmlPublicHost}`
}

onMounted(async () => {
  if (isLoggedIn.value) {
    await loadRecords()
  }
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: var(--color-background);
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
  border-bottom: 1px solid var(--color-border);
}

.left-links,
.right-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}

.nickname {
  color: var(--color-heading);
  font-weight: 600;
}

.text-btn {
  border: none;
  background: transparent;
  color: hsla(160, 100%, 37%, 1);
  cursor: pointer;
  font-size: 14px;
}

.content {
  padding: 32px 24px;
  max-width: 980px;
  margin: 0 auto;
}

.subtitle {
  color: #666;
  margin-bottom: 20px;
}

.card {
  background: #fff;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  padding: 18px;
  margin-bottom: 18px;
}

.upload-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tips {
  margin: 0 0 10px;
  color: #666;
  font-size: 14px;
}

.upload-form {
  display: grid;
  gap: 12px;
}

.upload-form label {
  display: grid;
  gap: 6px;
  color: #333;
}

.upload-form input,
.upload-form textarea {
  width: 100%;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  padding: 10px 12px;
}

.upload-form textarea {
  resize: vertical;
}

.primary-btn {
  width: fit-content;
  padding: 10px 16px;
  background-color: #42b883;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.primary-btn:hover {
  background-color: #35996d;
}

.records-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.empty {
  color: #888;
}
</style>