<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <h3>上传 HTML</h3>
        <p class="tips">
          使用提示：可填写子域名（可选），支持上传 html 文件或直接粘贴 html 文本。
        </p>
        <p class="tips">
          初始默认为私密，仅可自己访问。可到管理切换为公开。公开需等待管理员审核才能被真正公开！
        </p>

        <form class="upload-form" @submit.prevent="submitUpload">
          <label>
            子域名（可选） 如todolist,会生成：todolist.htmlhub.lyyxy.top
            <input v-model="uploadForm.subdomain" placeholder="例如：todolist">
          </label>
          <label>
            简介（可选）
            <input v-model="uploadForm.description" maxlength="500" placeholder="写一点介绍">
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
          <button class="primary-btn" type="submit">上传</button>
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
</style>
