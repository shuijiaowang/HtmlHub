<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <div class="head">
          <div>
            <h3 class="title">上传 HTML</h3>
            <p class="tips">可填写子域名（可选），支持上传 html 文件或直接粘贴 html 文本。</p>
<!--            <p class="tips">初始默认为私密，仅可自己访问；到“管理”可切换为公开（公开需等待管理员审核）。</p>-->
          </div>
          <div class="badge" aria-label="上传提示">建议：优先上传完整 HTML 文档</div>
        </div>

        <form class="upload-form" @submit.prevent="submitUpload">
          <div class="form-grid">
            <div class="left">
              <label class="field">
                <span class="label">子域名（可选）</span>
                <span class="hint">如 `todolist`，会生成：todolist.{{ htmlPublicHost }}</span>
                <input v-model="uploadForm.subdomain" placeholder="例如：todolist" inputmode="url" autocomplete="off">
              </label>
              <label class="field">
                <span class="label">简介（可选）</span>
                <input v-model="uploadForm.description" maxlength="500" placeholder="写一点介绍（方便后续管理）">
              </label>
              <label class="field field-file">
                <span class="label">选择本地 HTML 文件（可选）</span>
                <span class="file-row">
                  <input class="file" type="file" accept=".html,text/html" @change="onSelectFile">
                </span>
              </label>
              <div class="actions">
                <button class="primary-btn" type="submit">上传</button>
                <span class="mini">上传后可在“管理”里编辑简介、更新 HTML、切换公开/私密。</span>
              </div>
            </div>

            <label class="field right">
              <span class="label">HTML 文本</span>
              <textarea
                v-model="uploadForm.htmlContent"
                rows="14"
                placeholder="粘贴完整 HTML（含 <html>、<head>、<body> 及闭合标签）"
                required
              />
            </label>
          </div>
        </form>
      </section>
    </main>
  </div>
</template>

<script setup>
import { onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { uploadHtml } from '@/api/html'
import { useUserStore } from '@/stores/user'

/** 子域名 HTML 访问用 host，本地与线上由 Vite 环境变量区分 */
const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'

const uploadForm = reactive({
  subdomain: '',
  description: '',
  fileSize: 0,
  htmlContent: ''
})

const userStore = useUserStore()

onMounted(() => {
  if (userStore.token) return
  window.setTimeout(() => {
    if (!userStore.token) userStore.openAuthDialog('login')
  }, 3000)
})

const onSelectFile = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  uploadForm.fileSize = file.size
  uploadForm.htmlContent = await file.text()
}

const generateAutoFileName = () => {
  const now = new Date()
  const pad = (num) => String(num).padStart(2, '0')
  const datePart = `${now.getFullYear()}${pad(now.getMonth() + 1)}${pad(now.getDate())}`
  const timePart = `${pad(now.getHours())}${pad(now.getMinutes())}${pad(now.getSeconds())}`
  const randomPart = Math.random().toString(36).slice(2, 6)
  return `htmlhub-${datePart}-${timePart}-${randomPart}.html`
}

const submitUpload = async () => {
  if (!userStore.token) {
    ElMessage.warning('请先登录后再上传')
    userStore.openAuthDialog('login')
    return
  }

  const htmlContent = uploadForm.htmlContent.trim()
  if (!htmlContent) {
    ElMessage.warning('请填写HTML文本或上传HTML文件')
    return
  }
  const autoFileName = generateAutoFileName()
  const fileSize = uploadForm.fileSize > 0 ? uploadForm.fileSize : new Blob([htmlContent]).size

  await uploadHtml({
    subdomain: uploadForm.subdomain.trim(),
    fileName: autoFileName,
    description: uploadForm.description.trim(),
    fileSize,
    htmlContent
  })

  ElMessage.success('上传成功')
  uploadForm.subdomain = ''
  uploadForm.description = ''
  uploadForm.fileSize = 0
  uploadForm.htmlContent = ''
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: var(--color-background);
}

.content {
  padding: 32px 24px;
  max-width: 1040px;
  margin: 0 auto;
}

.card {
  background: var(--hh-surface-solid);
  border: 1px solid var(--hh-border);
  border-radius: var(--hh-radius-md);
  padding: 18px 18px 16px;
  margin-bottom: 18px;
  box-shadow: var(--hh-shadow-md);
}

.head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
  margin-bottom: 14px;
}

.title {
  margin: 0 0 6px;
  font-size: 18px;
  letter-spacing: -0.02em;
}

.tips {
  margin: 0 0 8px;
  color: var(--hh-text-2);
  font-size: 14px;
}

.badge {
  flex: 0 0 auto;
  padding: 8px 10px;
  border-radius: 999px;
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.18);
  background: rgb(var(--hh-brand-rgb) / 0.08);
  color: color-mix(in srgb, var(--hh-brand) 82%, var(--hh-text) 18%);
  font-size: 12px;
  line-height: 1.2;
  white-space: nowrap;
}

.upload-form {
  display: grid;
  gap: 10px;
}

.form-grid {
  display: grid;
  grid-template-columns: 0.95fr 1.05fr;
  gap: 14px;
}

.left {
  display: grid;
  gap: 12px;
  align-content: start;
}

.field {
  display: grid;
  gap: 6px;
  color: var(--hh-text);
}

.label {
  font-weight: 650;
  font-size: 13px;
}

.hint {
  font-size: 12px;
  color: var(--hh-text-3);
}

.upload-form input,
.upload-form textarea,
.upload-form .file {
  width: 100%;
  border: 1px solid var(--hh-border-2);
  border-radius: var(--hh-radius-sm);
  padding: 10px 12px;
  background: color-mix(in srgb, var(--hh-surface-solid) 94%, #000 0%);
}

.field-file .file-row {
  display: block;
}

.upload-form input.file[type='file'] {
  padding: 8px 10px;
  cursor: pointer;
  font-size: 14px;
}

.upload-form input.file::file-selector-button {
  margin-right: 12px;
  padding: 8px 14px;
  border-radius: var(--hh-radius-sm);
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.22);
  background: rgb(var(--hh-brand-rgb) / 0.10);
  color: color-mix(in srgb, var(--hh-brand) 88%, var(--hh-text) 12%);
  font-weight: 600;
  font-size: 14px;
  font-family: inherit;
  cursor: pointer;
}

.upload-form input.file::-webkit-file-upload-button {
  margin-right: 12px;
  padding: 8px 14px;
  border-radius: var(--hh-radius-sm);
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.22);
  background: rgb(var(--hh-brand-rgb) / 0.10);
  color: color-mix(in srgb, var(--hh-brand) 88%, var(--hh-text) 12%);
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
}

.upload-form textarea {
  resize: vertical;
  min-height: 320px;
}

.actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.mini {
  font-size: 12px;
  color: var(--hh-text-3);
}

.primary-btn {
  width: fit-content;
  padding: 10px 16px;
  background: linear-gradient(135deg, var(--hh-brand), color-mix(in srgb, var(--hh-brand) 60%, var(--hh-brand-2) 40%));
  color: #fff;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  box-shadow: 0 10px 22px rgb(var(--hh-brand-rgb) / 0.18);
}

.primary-btn:hover {
  filter: brightness(0.98);
}

@media (max-width: 640px) {
  .content {
    padding: 18px 12px;
  }

  .head {
    flex-direction: column;
    align-items: stretch;
  }

  .badge {
    width: fit-content;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .upload-form textarea {
    min-height: 240px;
  }
}
</style>
