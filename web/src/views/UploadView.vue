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
              <div class="field field-file">
                <span class="label">本地 HTML 文件（可选）</span>
                <span class="hint">支持拖拽 .html 文件到下方区域，或点击选择；内容会自动填入右侧编辑器</span>
                <div
                  class="drop-zone"
                  :class="{ 'drop-zone--active': isDragOver, 'drop-zone--filled': selectedFileName }"
                  role="button"
                  tabindex="0"
                  aria-label="拖拽或点击上传 HTML 文件"
                  @click="openFilePicker"
                  @keydown.enter.prevent="openFilePicker"
                  @keydown.space.prevent="openFilePicker"
                  @dragenter.prevent="onDragEnter"
                  @dragover.prevent="onDragOver"
                  @dragleave.prevent="onDragLeave"
                  @drop.prevent="onDrop"
                >
                  <input
                    ref="fileInputRef"
                    class="file-input-hidden"
                    type="file"
                    accept=".html,text/html"
                    @change="onSelectFile"
                  >
                  <span class="drop-zone-icon" aria-hidden="true">📄</span>
                  <span class="drop-zone-title">
                    {{ selectedFileName ? selectedFileName : '拖拽 HTML 文件到此处' }}
                  </span>
                  <span class="drop-zone-desc">
                    {{ selectedFileName ? '已载入，可重新拖拽或点击更换文件' : '或点击此区域选择文件 · 仅支持 .html' }}
                  </span>
                </div>
              </div>
              <div class="actions">
                <button class="primary-btn" type="submit">上传</button>
                <span class="mini">上传后可在“管理”里编辑简介、更新 HTML、切换公开/私密。</span>
              </div>
            </div>

            <label class="field right">
              <span class="label">HTML 文本</span>
              <div class="html-toolbar" role="toolbar" aria-label="HTML 文本操作">
                <button type="button" class="toolbar-btn" @click="copyHtmlContent">复制</button>
                <button type="button" class="toolbar-btn" @click="pasteHtmlContent">粘贴</button>
                <button type="button" class="toolbar-btn" @click="clearHtmlContent">清空</button>
              </div>
              <textarea
                ref="htmlTextareaRef"
                v-model="htmlContent"
                rows="14"
                placeholder="粘贴完整 HTML（含 <html>、<head>、<body> 及闭合标签）"
                required
                @paste="onHtmlNativePaste"
              />
            </label>
          </div>
        </form>
      </section>
    </main>
  </div>
</template>

<script setup>
import { nextTick, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { uploadHtml } from '@/api/html'
import { useUserStore } from '@/stores/user'

/** 子域名 HTML 访问用 host，本地与线上由 Vite 环境变量区分 */
const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'

const uploadForm = reactive({
  subdomain: '',
  description: '',
  fileSize: 0
})

/** 大段 HTML 单独 ref，减少与表单其余字段的响应式耦合 */
const htmlContent = ref('')
const htmlTextareaRef = ref(null)
const fileInputRef = ref(null)
const isDragOver = ref(false)
const selectedFileName = ref('')

const userStore = useUserStore()

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

const afterHtmlValueUpdate = () => {
  nextTick(() => {
    requestAnimationFrame(() => {
      focusTextareaStartAndReveal(htmlTextareaRef.value)
    })
  })
}

const copyHtmlContent = async () => {
  const text = htmlContent.value
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

const pasteHtmlContent = async () => {
  let text
  try {
    text = await navigator.clipboard.readText()
  } catch {
    ElMessage.error('粘贴失败：无法读取剪贴板，请检查权限或在内置浏览器中重试')
    return
  }
  htmlContent.value = text
  afterHtmlValueUpdate()
}

const clearHtmlContent = () => {
  htmlContent.value = ''
  afterHtmlValueUpdate()
}

const onHtmlNativePaste = (e) => {
  const cd = e.clipboardData
  if (!cd) return
  const text = cd.getData('text/plain')
  if (text === '') return
  e.preventDefault()
  const el = e.target
  const start = el.selectionStart
  const end = el.selectionEnd
  const prev = htmlContent.value
  htmlContent.value = prev.slice(0, start) + text + prev.slice(end)
  afterHtmlValueUpdate()
}

onMounted(() => {
  if (userStore.token) return
  window.setTimeout(() => {
    if (!userStore.token) userStore.openAuthDialog('login')
  }, 3000)
})

const isHtmlFile = (file) => {
  if (!file) return false
  const name = file.name.toLowerCase()
  return name.endsWith('.html') || file.type === 'text/html'
}

const loadHtmlFile = async (file) => {
  if (!isHtmlFile(file)) {
    ElMessage.warning('仅支持 .html 文件，请重新选择')
    return
  }
  uploadForm.fileSize = file.size
  selectedFileName.value = file.name
  htmlContent.value = await file.text()
  afterHtmlValueUpdate()
  ElMessage.success(`已载入：${file.name}`)
}

const openFilePicker = () => {
  fileInputRef.value?.click()
}

const onSelectFile = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  await loadHtmlFile(file)
  event.target.value = ''
}

const onDragEnter = () => {
  isDragOver.value = true
}

const onDragOver = () => {
  isDragOver.value = true
}

const onDragLeave = (event) => {
  if (event.currentTarget?.contains(event.relatedTarget)) return
  isDragOver.value = false
}

const onDrop = async (event) => {
  isDragOver.value = false
  const file = event.dataTransfer?.files?.[0]
  if (!file) return
  await loadHtmlFile(file)
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

  const body = htmlContent.value.trim()
  if (!body) {
    ElMessage.warning('请填写HTML文本或上传HTML文件')
    return
  }
  const autoFileName = generateAutoFileName()
  const fileSize = uploadForm.fileSize > 0 ? uploadForm.fileSize : new Blob([body]).size

  await uploadHtml({
    subdomain: uploadForm.subdomain.trim(),
    fileName: autoFileName,
    description: uploadForm.description.trim(),
    fileSize,
    htmlContent: body
  })

  ElMessage.success('上传成功')
  uploadForm.subdomain = ''
  uploadForm.description = ''
  uploadForm.fileSize = 0
  htmlContent.value = ''
  selectedFileName.value = ''
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
.upload-form textarea {
  width: 100%;
  border: 1px solid var(--hh-border-2);
  border-radius: var(--hh-radius-sm);
  padding: 10px 12px;
  background: color-mix(in srgb, var(--hh-surface-solid) 94%, #000 0%);
}

.file-input-hidden {
  display: none;
}

.drop-zone {
  display: grid;
  justify-items: center;
  gap: 4px;
  padding: 22px 16px;
  border: 1.5px dashed var(--hh-border-2);
  border-radius: var(--hh-radius-sm);
  background: color-mix(in srgb, var(--hh-surface-solid) 96%, rgb(var(--hh-brand-rgb) / 0.04) 4%);
  cursor: pointer;
  text-align: center;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}

.drop-zone:hover,
.drop-zone:focus-visible {
  border-color: rgb(var(--hh-brand-rgb) / 0.35);
  background: rgb(var(--hh-brand-rgb) / 0.06);
  outline: none;
}

.drop-zone--active {
  border-color: var(--hh-brand);
  background: rgb(var(--hh-brand-rgb) / 0.10);
  box-shadow: 0 0 0 3px rgb(var(--hh-brand-rgb) / 0.12);
}

.drop-zone--filled {
  border-style: solid;
  border-color: rgb(var(--hh-brand-rgb) / 0.28);
}

.drop-zone-icon {
  font-size: 24px;
  line-height: 1;
}

.drop-zone-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--hh-text);
  word-break: break-all;
}

.drop-zone-desc {
  font-size: 12px;
  color: var(--hh-text-3);
  line-height: 1.4;
}

.html-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 6px;
}

.toolbar-btn {
  padding: 6px 12px;
  border-radius: var(--hh-radius-sm);
  border: 1px solid var(--hh-border-2);
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
