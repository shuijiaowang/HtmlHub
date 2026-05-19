<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <div class="head-row">
          <div>
            <h1>作品展示</h1>
            <p class="subtitle">浏览所有用户公开的 HTML 页面</p>
          </div>
          <div class="sort-bar">
            <button
              type="button"
              class="sort-btn"
              :class="{ active: sortBy === 'latest' }"
              @click="changeSort('latest')"
            >
              最新
            </button>
            <button
              type="button"
              class="sort-btn"
              :class="{ active: sortBy === 'likes' }"
              @click="changeSort('likes')"
            >
              最多点赞
            </button>
          </div>
        </div>

        <p v-if="loading" class="hint">加载中...</p>
        <p v-else-if="records.length === 0" class="hint">暂无公开页面</p>

        <div v-else class="record-list">
          <article v-for="item in records" :key="item.id" class="record-item">
            <div class="record-top">
              <div class="author">
                <span class="author-label">上传人</span>
                <strong>{{ item.nickname || '匿名用户' }}</strong>
              </div>
              <button
                type="button"
                class="like-btn"
                :class="{ liked: item.liked }"
                :disabled="likeLoadingId === item.id"
                @click="toggleLike(item)"
              >
                <span aria-hidden="true">{{ item.liked ? '♥' : '♡' }}</span>
                {{ item.likeCount || 0 }}
              </button>
            </div>

            <p class="record-desc">
              {{ item.description || '暂无简介' }}
            </p>

            <p class="record-route">
              访问链接：{{ item.subdomain }}.{{ htmlPublicHost }}
              <button type="button" class="open-link" @click="openPage(item)">点击打开</button>
            </p>

            <div class="record-meta">
              <span>访问 {{ item.visitCount || 0 }} 次</span>
              <span>{{ formatDate(item.createdAt) }}</span>
            </div>
          </article>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getPublicHtmlList, likeHtmlRecord, unlikeHtmlRecord } from '@/api/html'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'
const userStore = useUserStore()

const records = ref([])
const sortBy = ref('latest')
const loading = ref(false)
const likeLoadingId = ref(null)

const loadRecords = async () => {
  loading.value = true
  try {
    const res = await getPublicHtmlList({ sort: sortBy.value, page: 1, pageSize: 100 })
    records.value = Array.isArray(res.data?.list) ? res.data.list : []
  } finally {
    loading.value = false
  }
}

const changeSort = async (sort) => {
  if (sortBy.value === sort) return
  sortBy.value = sort
  await loadRecords()
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (Number.isNaN(date.getTime())) return dateStr
  return date.toLocaleString()
}

const openPage = (item) => {
  if (!item?.subdomain) return
  const url = `${window.location.protocol}//${item.subdomain}.${htmlPublicHost}`
  window.open(url, '_blank', 'noopener,noreferrer')
}

const toggleLike = async (item) => {
  if (!userStore.token) {
    ElMessage.warning('请先登录后再点赞')
    userStore.openAuthDialog('login')
    return
  }
  likeLoadingId.value = item.id
  try {
    if (item.liked) {
      await unlikeHtmlRecord(item.id)
      item.liked = false
      item.likeCount = Math.max(0, (item.likeCount || 0) - 1)
    } else {
      await likeHtmlRecord(item.id)
      item.liked = true
      item.likeCount = (item.likeCount || 0) + 1
    }
  } finally {
    likeLoadingId.value = null
  }
}

onMounted(loadRecords)
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
  padding: 18px;
  box-shadow: var(--hh-shadow-md);
}

.head-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 20px;
}

h1 {
  margin: 0;
  font-size: 22px;
}

.subtitle {
  margin: 6px 0 0;
  color: var(--hh-text-3);
  font-size: 13px;
}

.sort-bar {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.sort-btn {
  border: 1px solid var(--hh-border);
  background: transparent;
  border-radius: 999px;
  padding: 6px 14px;
  cursor: pointer;
  font-size: 13px;
  font-family: inherit;
  color: var(--hh-text-2);
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.sort-btn:hover {
  border-color: rgb(var(--hh-brand-rgb) / 0.28);
}

.sort-btn.active {
  background: rgb(var(--hh-brand-rgb) / 0.12);
  border-color: rgb(var(--hh-brand-rgb) / 0.28);
  color: var(--hh-brand);
  font-weight: 600;
}

.record-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.record-item {
  border: 1px solid var(--hh-border);
  border-radius: var(--hh-radius-md);
  padding: 14px;
  background: color-mix(in srgb, var(--hh-surface-solid) 94%, #000 0%);
  box-shadow: var(--hh-shadow-sm);
  transition: transform 0.1s ease, box-shadow 0.15s ease, border-color 0.15s ease;
}

.record-item:hover {
  transform: translateY(-1px);
  border-color: rgb(var(--hh-brand-rgb) / 0.22);
  box-shadow: 0 14px 30px rgba(15, 23, 42, 0.1);
}

.record-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.author-label {
  display: block;
  font-size: 12px;
  color: var(--hh-text-3);
}

.record-desc {
  margin: 10px 0 0;
  color: var(--hh-text-2);
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
}

.record-route {
  margin: 8px 0 0;
  font-size: 13px;
  color: var(--hh-text-2);
  word-break: break-all;
}

.open-link {
  margin-left: 8px;
  border: none;
  background: none;
  padding: 0;
  color: color-mix(in srgb, var(--hh-brand) 86%, #000 0%);
  cursor: pointer;
  text-decoration: underline;
  font: inherit;
}

.open-link:hover {
  color: var(--hh-brand);
}

.like-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid var(--hh-border);
  border-radius: 999px;
  padding: 4px 10px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
  font-family: inherit;
  flex-shrink: 0;
  transition: background-color 0.15s ease, border-color 0.15s ease, color 0.15s ease;
}

.like-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.like-btn.liked {
  color: #e11d48;
  border-color: #fecdd3;
  background: #fff1f2;
}

.record-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 10px;
  font-size: 12px;
  color: var(--hh-text-3);
}

.hint {
  color: var(--hh-text-3);
  padding: 20px 0;
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

  .head-row {
    flex-direction: column;
  }

  .sort-bar {
    width: 100%;
  }

  .sort-btn {
    flex: 1;
    text-align: center;
  }
}
</style>
