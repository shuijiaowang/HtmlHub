<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <div class="records-title-row">
          <div class="title-wrap">
            <h3 class="title">我的足迹</h3>
            <p class="subtitle">这里汇总你云同步过数据的页面，以及你点赞收藏过的页面。</p>
          </div>
          <button class="text-btn refresh" :disabled="loading" @click="reloadActive">刷新</button>
        </div>

        <p v-if="!isLoggedIn" class="empty">请先登录后查看我的足迹。</p>

        <el-tabs v-else v-model="activeTab" class="foot-tabs">
          <!-- 同步数据 -->
          <el-tab-pane name="data">
            <template #label>
              <span>同步数据<span class="tab-count">{{ dataList.length }}</span></span>
            </template>

            <div class="data-summary">
              <div class="summary-item">
                <span class="summary-label">已同步页面</span>
                <span class="summary-value">{{ dataSummary.count }}</span>
              </div>
              <div class="summary-item">
                <span class="summary-label">同步数据占用</span>
                <span class="summary-value">{{ formatBytes(dataSummary.totalBytes) }}</span>
              </div>
              <div class="summary-item">
                <span class="summary-label">单页限制</span>
                <span class="summary-value">{{ formatBytes(dataSummary.maxDataBytes) }}</span>
              </div>
              <button
                class="text-btn danger-btn clear-all"
                type="button"
                :disabled="dataList.length === 0"
                @click="clearAll"
              >
                清空全部
              </button>
            </div>

            <p v-if="dataList.length === 0" class="empty">暂无同步数据。在你自己的页面上登录并云同步后，会显示在这里。</p>
            <div v-else class="record-list">
              <article v-for="item in dataList" :key="item.id" class="record-item">
                <div class="record-head">
                  <strong>{{ item.subdomain }}</strong>
                  <span>{{ formatBytes(item.dataBytes) }}</span>
                </div>
                <p class="record-route">
                  访问链接：{{ item.subdomain }}.{{ htmlPublicHost }}
                  <button
                    type="button"
                    class="share-link"
                    title="在新标签页打开"
                    @click="openPage(item, $event)"
                    @auxclick="openPage(item, $event)"
                  >
                    点击打开
                  </button>
                </p>
                <p v-if="item.description" class="record-desc">简介：{{ item.description }}</p>
                <el-progress
                  :percentage="percent(item.dataBytes, dataSummary.maxDataBytes)"
                  :stroke-width="6"
                  :color="barColor(item.dataBytes, dataSummary.maxDataBytes)"
                  :show-text="false"
                />
                <div class="record-meta">
                  <span>已用额度：{{ percent(item.dataBytes, dataSummary.maxDataBytes) }}%</span>
                  <span>更新时间：{{ formatDate(item.updatedAt) }}</span>
                </div>
                <div class="record-actions">
                  <button class="text-btn" type="button" :disabled="item._exporting" @click="exportData(item)">
                    导出 JSON
                  </button>
                  <button class="text-btn danger-btn" type="button" @click="removeData(item)">删除数据</button>
                </div>
              </article>
            </div>
          </el-tab-pane>

          <!-- 我的点赞 -->
          <el-tab-pane name="likes">
            <template #label>
              <span>我的点赞<span class="tab-count">{{ likeList.length }}</span></span>
            </template>

            <p v-if="likeList.length === 0" class="empty">还没有点赞过任何页面。去展示页逛逛吧。</p>
            <div v-else class="record-list">
              <article v-for="item in likeList" :key="item.id" class="record-item">
                <div class="record-head">
                  <strong>{{ item.fileName || item.subdomain }}</strong>
                  <span v-if="item.nickname">by {{ item.nickname }}</span>
                </div>
                <p class="record-route">
                  访问链接：{{ item.subdomain }}.{{ htmlPublicHost }}
                  <button
                    type="button"
                    class="share-link"
                    title="在新标签页打开"
                    @click="openPage(item, $event)"
                    @auxclick="openPage(item, $event)"
                  >
                    点击打开
                  </button>
                </p>
                <p v-if="item.description" class="record-desc">简介：{{ item.description }}</p>
                <div class="record-meta">
                  <span class="meta-metric" title="访问次数">
                    <el-icon class="metric-icon" aria-hidden="true"><View /></el-icon>
                    {{ item.visitCount || 0 }}
                  </span>
                  <span class="meta-metric" title="点赞次数">
                    <svg class="metric-icon heart-icon" viewBox="0 0 24 24" aria-hidden="true">
                      <path
                        fill="currentColor"
                        d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                      />
                    </svg>
                    {{ item.likeCount || 0 }}
                  </span>
                  <span>点赞时间：{{ formatDate(item.likedAt) }}</span>
                </div>
                <div class="record-actions">
                  <button class="text-btn danger-btn" type="button" @click="removeLike(item)">取消点赞</button>
                </div>
              </article>
            </div>
          </el-tab-pane>
        </el-tabs>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import {
  getMySyncDataList,
  exportMySyncData,
  deleteMySyncData,
  clearMySyncData,
  getMyLikedList,
  unlikeHtmlRecord
} from '@/api/html'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View } from '@element-plus/icons-vue'

const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'

const userStore = useUserStore()
const isLoggedIn = computed(() => !!userStore.token)

const activeTab = ref('data')
const loading = ref(false)

const dataList = ref([])
const dataSummary = ref({ count: 0, totalBytes: 0, maxDataBytes: 0 })
const likeList = ref([])
const dataLoaded = ref(false)
const likesLoaded = ref(false)

const formatBytes = (size) => {
  const n = Number(size) || 0
  if (n <= 0) return '0 B'
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / (1024 * 1024)).toFixed(2)} MB`
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (Number.isNaN(date.getTime())) return dateStr
  return date.toLocaleString()
}

const percent = (used, limit) => {
  const u = Number(used) || 0
  const l = Number(limit) || 0
  if (l <= 0) return 0
  return Math.min(100, Math.round((u / l) * 100))
}

const barColor = (used, limit) => {
  const p = percent(used, limit)
  if (p >= 100) return '#d93025'
  if (p >= 80) return '#e6a23c'
  return '#0ea5e9'
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getMySyncDataList()
    const data = res.data || {}
    dataList.value = Array.isArray(data.list) ? data.list : []
    dataSummary.value = {
      count: data.count || dataList.value.length,
      totalBytes: data.totalBytes || 0,
      maxDataBytes: data.maxDataBytes || 0
    }
    dataLoaded.value = true
  } finally {
    loading.value = false
  }
}

const loadLikes = async () => {
  loading.value = true
  try {
    const res = await getMyLikedList()
    likeList.value = Array.isArray(res.data?.list) ? res.data.list : []
    likesLoaded.value = true
  } finally {
    loading.value = false
  }
}

const reloadActive = () => {
  if (!userStore.token) {
    ElMessage.warning('请先登录')
    userStore.openAuthDialog('login')
    return
  }
  if (activeTab.value === 'data') return loadData()
  return loadLikes()
}

const buildShareUrl = (item) => {
  const subdomain = item?.subdomain
  if (!subdomain) return ''
  const url = new URL(`${window.location.protocol}//${subdomain}.${htmlPublicHost}`)
  if (userStore.token) {
    url.searchParams.set('token', userStore.token)
  }
  return url.toString()
}

const openPage = (item, event) => {
  if (event) {
    if (event.type === 'auxclick' && event.button !== 1) return
    if (event.type === 'click' && event.button !== 0) return
    event.preventDefault()
  }
  const url = buildShareUrl(item)
  if (!url) {
    ElMessage.warning('该页面暂不可访问')
    return
  }
  window.open(url, '_blank', 'noopener,noreferrer')
}

const downloadJson = (filename, content) => {
  let text = content
  try {
    text = JSON.stringify(JSON.parse(content), null, 2)
  } catch {
    // 非标准 JSON 时按原文导出
  }
  const blob = new Blob([text], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const exportData = async (item) => {
  item._exporting = true
  try {
    const res = await exportMySyncData(item.id)
    const json = res.data?.dataJson
    if (json == null) {
      ElMessage.warning('暂无可导出的数据')
      return
    }
    const base = item.subdomain || `data-${item.id}`
    downloadJson(`${base}.json`, json)
    ElMessage.success('已开始下载')
  } finally {
    item._exporting = false
  }
}

const removeData = async (item) => {
  try {
    await ElMessageBox.confirm(
      `确定删除「${item.subdomain}」的云同步数据吗？该操作不可恢复。`,
      '删除同步数据',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }
  await deleteMySyncData(item.id)
  ElMessage.success('已删除')
  await loadData()
}

const clearAll = async () => {
  try {
    await ElMessageBox.confirm(
      '确定清空你的全部云同步数据吗？该操作不可恢复。',
      '清空全部',
      { type: 'warning', confirmButtonText: '清空', cancelButtonText: '取消' }
    )
  } catch {
    return
  }
  await clearMySyncData()
  ElMessage.success('已清空')
  await loadData()
}

const removeLike = async (item) => {
  try {
    await ElMessageBox.confirm(
      `确定取消点赞「${item.fileName || item.subdomain}」吗？`,
      '取消点赞',
      { type: 'warning', confirmButtonText: '取消点赞', cancelButtonText: '再想想' }
    )
  } catch {
    return
  }
  await unlikeHtmlRecord(item.id)
  likeList.value = likeList.value.filter((row) => row.id !== item.id)
  ElMessage.success('已取消点赞')
}

onMounted(() => {
  if (userStore.token) {
    loadData()
    return
  }
  window.setTimeout(() => {
    if (!userStore.token) userStore.openAuthDialog('login')
  }, 1500)
})

watch(activeTab, (tab) => {
  if (!userStore.token) return
  if (tab === 'likes' && !likesLoaded.value) loadLikes()
  if (tab === 'data' && !dataLoaded.value) loadData()
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

.text-btn:hover:not(:disabled) {
  background: rgb(var(--hh-brand-rgb) / 0.10);
  border-color: rgb(var(--hh-brand-rgb) / 0.18);
}

.text-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.danger-btn {
  color: #d93025;
}

.refresh {
  white-space: nowrap;
}

.tab-count {
  display: inline-block;
  margin-left: 6px;
  padding: 0 7px;
  font-size: 12px;
  line-height: 18px;
  border-radius: 999px;
  color: var(--hh-text-3);
  background: rgb(var(--hh-brand-rgb) / 0.10);
}

.data-summary {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px 22px;
  padding: 6px 0 14px;
  border-bottom: 1px solid var(--hh-border);
  margin-bottom: 14px;
}

.summary-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.summary-label {
  font-size: 12px;
  color: var(--hh-text-3);
}

.summary-value {
  font-weight: 800;
  font-size: 18px;
  letter-spacing: -0.02em;
}

.clear-all {
  margin-left: auto;
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

.record-route {
  margin: 8px 0 10px;
  color: var(--hh-text-2);
  font-size: 13px;
  word-break: break-word;
}

.record-desc {
  margin: 0 0 10px;
  color: var(--hh-text-2);
  font-size: 13px;
  line-height: 1.6;
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

.meta-metric {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.metric-icon {
  width: 1em;
  height: 1em;
  font-size: 16px;
  flex-shrink: 0;
}

.heart-icon {
  color: #e85d6a;
}

.record-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.empty {
  color: var(--hh-text-3);
  padding: 14px 0 6px;
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
}
</style>
