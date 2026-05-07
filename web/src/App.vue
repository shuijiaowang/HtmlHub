<script setup>
import { computed, watch } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import AuthDialog from '@/components/AuthDialog.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isHomeRoute = computed(() => route.path.startsWith('/home'))
const isLoggedIn = computed(() => !!userStore.token)
const displayNickname = computed(() => userStore.userInfo.nickname || '已登录用户')

const logout = async () => {
  await userStore.logout()
}

const openLogin = () => userStore.openAuthDialog('login')
const openRegister = () => userStore.openAuthDialog('register')

watch(
  () => route.query?.auth,
  async (val) => {
    if (val === 'login' || val === 'register') {
      userStore.openAuthDialog(val)
      const nextQuery = { ...route.query }
      delete nextQuery.auth
      await router.replace({ query: nextQuery })
    }
  },
  { immediate: true }
)
</script>

<template>
  <div id="app">
    <header v-if="isHomeRoute" class="top-nav">
      <div class="left-links">
        <router-link to="/home">主页</router-link>
        <router-link to="/home/upload">上传 HTML</router-link>
        <router-link to="/home/manage">个人 HTML 管理</router-link>
        <router-link to="/home/showcase">展示页</router-link>
      </div>
      <div class="right-actions">
        <template v-if="isLoggedIn">
          <span class="nickname">{{ displayNickname }}</span>
          <button class="text-btn" @click="logout">退出</button>
        </template>
        <template v-else>
          <span class="not-login">未登录</span>
          <button class="text-btn" @click="openRegister">注册</button>
          <button class="text-btn" @click="openLogin">登录</button>
        </template>
      </div>
    </header>
    <RouterView />
    <AuthDialog />
  </div>
</template>

<style scoped>
#app {
  min-height: 100vh;
}

.top-nav {
  position: sticky;
  top: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 64px;
  padding: 10px 24px;
  border-bottom: 1px solid var(--hh-border);
  background: color-mix(in srgb, var(--hh-surface) 88%, transparent);
  backdrop-filter: blur(10px);
}

.left-links,
.right-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.left-links {
  flex: 1;
  flex-wrap: wrap;
  min-width: 0;
}

.left-links a {
  padding: 8px 10px;
  border-radius: 10px;
  color: var(--hh-text-2);
  border: 1px solid transparent;
}

.left-links a:hover {
  text-decoration: none;
  color: var(--hh-text);
  background: rgba(16, 185, 129, 0.10);
}

.nickname {
  color: var(--hh-text);
  font-weight: 600;
}

.text-btn {
  border: none;
  background: rgba(16, 185, 129, 0.10);
  color: color-mix(in srgb, var(--hh-brand) 88%, #000 0%);
  cursor: pointer;
  font-size: 14px;
  padding: 8px 10px;
  border-radius: 10px;
}

.text-btn:hover {
  background: rgba(16, 185, 129, 0.16);
}

.not-login {
  color: var(--hh-text-3);
  font-size: 14px;
}

a.router-link-active {
  color: var(--hh-text);
  font-weight: 600;
  background: rgba(16, 185, 129, 0.16);
  border-color: rgba(16, 185, 129, 0.22);
}

@media (max-width: 640px) {
  .top-nav {
    padding: 10px 12px;
    gap: 10px;
    align-items: flex-start;
  }

  .right-actions {
    flex-shrink: 0;
    gap: 8px;
  }

  .text-btn {
    padding: 8px 10px;
  }
}
</style>