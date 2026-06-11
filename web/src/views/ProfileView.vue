<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <div class="records-title-row">
          <div class="title-wrap">
            <h3 class="title">个人中心</h3>
            <p class="subtitle">查看你的账号信息、资源限制与当前使用情况。</p>
          </div>
          <button class="text-btn refresh" :disabled="loading" @click="loadProfile">刷新</button>
        </div>

        <p v-if="!isLoggedIn" class="empty">请先登录后查看个人中心。</p>

        <template v-else>
          <!-- 账号信息 -->
          <div class="profile-head">
            <div class="avatar" aria-hidden="true">{{ avatarText }}</div>
            <div class="profile-main">
              <div class="nickname-row">
                <template v-if="!editingNickname">
                  <span class="profile-nickname">{{ profile.nickname || '-' }}</span>
                  <button class="text-btn" type="button" @click="startEditNickname">修改昵称</button>
                </template>
                <template v-else>
                  <el-input
                    v-model="nicknameDraft"
                    class="nickname-input"
                    size="small"
                    maxlength="20"
                    show-word-limit
                    placeholder="2-20 个字符"
                    @keyup.enter="saveNickname"
                  />
                  <button class="text-btn" type="button" :disabled="nicknameSaving" @click="saveNickname">保存</button>
                  <button class="text-btn" type="button" :disabled="nicknameSaving" @click="cancelEditNickname">取消</button>
                </template>
                <span class="role-badge" :class="`role-${profile.role}`">{{ roleLabel }}</span>
              </div>
              <div class="profile-meta">
                <span>邮箱：{{ profile.email || '-' }}</span>
                <span>UUID：{{ profile.uuid || '-' }}</span>
                <span>注册时间：{{ profile.createdAt || '-' }}</span>
              </div>
            </div>
          </div>

          <!-- 使用情况 -->
          <h4 class="section-title">页面数量</h4>
          <div class="usage-grid">
            <div class="usage-card">
              <div class="usage-label">未删除页面</div>
              <div class="usage-value">
                {{ profile.activeUploadCount }}<span class="usage-limit"> / {{ profile.maxActiveHtmlRecords }}</span>
              </div>
              <el-progress
                :percentage="percent(profile.activeUploadCount, profile.maxActiveHtmlRecords)"
                :stroke-width="8"
                :color="barColor(profile.activeUploadCount, profile.maxActiveHtmlRecords)"
                :show-text="false"
              />
              <div class="usage-hint">当前生效的页面数量上限</div>
            </div>
            <div class="usage-card">
              <div class="usage-label">累计页面</div>
              <div class="usage-value">
                {{ profile.totalUploadCount }}<span class="usage-limit"> / {{ profile.maxTotalHtmlRecords }}</span>
              </div>
              <el-progress
                :percentage="percent(profile.totalUploadCount, profile.maxTotalHtmlRecords)"
                :stroke-width="8"
                :color="barColor(profile.totalUploadCount, profile.maxTotalHtmlRecords)"
                :show-text="false"
              />
              <div class="usage-hint">含已删除页面，删除后不会释放累计额度</div>
            </div>
          </div>

          <h4 class="section-title">单项限制</h4>
          <div class="usage-grid">
            <div class="usage-card">
              <div class="usage-label">HTML 单页限制</div>
              <div class="usage-value">{{ formatBytes(profile.maxHtmlContentBytes) }}</div>
              <div class="usage-hint">单个页面 HTML 内容的最大体积</div>
            </div>
            <div class="usage-card">
              <div class="usage-label">同步数据单页限制</div>
              <div class="usage-value">{{ formatBytes(profile.maxHtmlDataBytes) }}</div>
              <div class="usage-hint">单个页面单次云同步数据的最大体积</div>
            </div>
          </div>

          <h4 class="section-title">空间占用</h4>
          <div class="usage-grid">
            <div class="usage-card">
              <div class="usage-label">HTML 占用</div>
              <div class="usage-value">{{ formatBytes(profile.activeHtmlBytes) }}</div>
              <div class="usage-hint">未删除页面 HTML 内容的总体积</div>
            </div>
            <div class="usage-card">
              <div class="usage-label">同步数据占用</div>
              <div class="usage-value">{{ formatBytes(profile.htmlDataBytes) }}</div>
              <div class="usage-hint">全部云同步数据的总体积</div>
            </div>
          </div>
        </template>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { getUserProfile, updateUserProfile } from '@/api/user'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const isLoggedIn = computed(() => !!userStore.token)

const loading = ref(false)
const profile = ref({
  nickname: '',
  email: '',
  uuid: '',
  role: '',
  createdAt: '',
  maxHtmlContentBytes: 0,
  maxHtmlDataBytes: 0,
  maxActiveHtmlRecords: 0,
  maxTotalHtmlRecords: 0,
  activeUploadCount: 0,
  totalUploadCount: 0,
  activeHtmlBytes: 0,
  htmlDataBytes: 0
})

const editingNickname = ref(false)
const nicknameDraft = ref('')
const nicknameSaving = ref(false)

const roleLabel = computed(() => {
  const map = { user: '普通用户', admin: '管理员', super_admin: '超级管理员' }
  return map[profile.value.role] || '普通用户'
})

const avatarText = computed(() => {
  const name = (profile.value.nickname || '').trim()
  return name ? name.slice(0, 1).toUpperCase() : 'U'
})

const formatBytes = (size) => {
  const n = Number(size) || 0
  if (n <= 0) return '0 B'
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / (1024 * 1024)).toFixed(2)} MB`
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

const loadProfile = async () => {
  if (!userStore.token) {
    ElMessage.warning('请先登录后查看个人中心')
    userStore.openAuthDialog('login')
    return
  }
  loading.value = true
  try {
    const res = await getUserProfile()
    if (res.data) profile.value = res.data
  } finally {
    loading.value = false
  }
}

const startEditNickname = () => {
  nicknameDraft.value = profile.value.nickname || ''
  editingNickname.value = true
}

const cancelEditNickname = () => {
  editingNickname.value = false
}

const saveNickname = async () => {
  const next = nicknameDraft.value.trim()
  if (next.length < 2 || next.length > 20) {
    ElMessage.warning('昵称长度需为 2-20 位')
    return
  }
  if (next === profile.value.nickname) {
    editingNickname.value = false
    return
  }
  nicknameSaving.value = true
  try {
    const res = await updateUserProfile({ nickname: next })
    const saved = res.data?.nickname || next
    profile.value.nickname = saved
    userStore.setUserInfo({ nickname: saved })
    editingNickname.value = false
    ElMessage.success('昵称已更新')
  } finally {
    nicknameSaving.value = false
  }
}

onMounted(() => {
  if (userStore.token) {
    loadProfile()
    return
  }
  window.setTimeout(() => {
    if (!userStore.token) userStore.openAuthDialog('login')
  }, 1500)
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
  padding: 18px 18px 16px;
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

.refresh {
  white-space: nowrap;
}

.empty {
  color: var(--hh-text-3);
  padding: 14px 0 6px;
}

.profile-head {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 0 6px;
  border-bottom: 1px solid var(--hh-border);
  margin-bottom: 10px;
}

.avatar {
  flex: 0 0 auto;
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  font-size: 24px;
  font-weight: 800;
  color: #fff;
  background: linear-gradient(135deg, var(--hh-brand) 0%, color-mix(in srgb, var(--hh-brand) 60%, #000) 120%);
  box-shadow: 0 8px 18px rgb(var(--hh-brand-rgb) / 0.22);
}

.profile-main {
  min-width: 0;
  flex: 1;
}

.nickname-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.profile-nickname {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: -0.02em;
  word-break: break-word;
}

.nickname-input {
  max-width: 220px;
}

.role-badge {
  font-size: 12px;
  padding: 2px 10px;
  border-radius: 999px;
  border: 1px solid var(--hh-border);
  color: var(--hh-text-2);
  background: color-mix(in srgb, var(--hh-surface-solid) 88%, rgb(var(--hh-brand-rgb) / 0.10) 12%);
}

.role-admin,
.role-super_admin {
  color: #b06800;
  border-color: rgba(230, 162, 60, 0.4);
  background: rgba(230, 162, 60, 0.12);
}

.profile-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 18px;
  margin-top: 8px;
  color: var(--hh-text-3);
  font-size: 13px;
  word-break: break-all;
}

.section-title {
  margin: 18px 0 10px;
  font-size: 14px;
  color: var(--hh-text-2);
}

.usage-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.usage-card {
  border: 1px solid var(--hh-border);
  border-radius: var(--hh-radius-md);
  padding: 14px;
  background: color-mix(in srgb, var(--hh-surface-solid) 92%, rgb(var(--hh-brand-rgb) / 0.06) 8%);
}

.usage-label {
  font-size: 12px;
  color: var(--hh-text-3);
}

.usage-value {
  margin: 4px 0 8px;
  font-weight: 800;
  font-size: 22px;
  letter-spacing: -0.02em;
}

.usage-limit {
  font-size: 14px;
  font-weight: 600;
  color: var(--hh-text-3);
}

.usage-hint {
  margin-top: 8px;
  font-size: 12px;
  color: var(--hh-text-3);
  line-height: 1.5;
}

@media (max-width: 640px) {
  .content {
    padding: 18px 12px;
  }

  .usage-grid {
    grid-template-columns: 1fr;
  }
}
</style>
