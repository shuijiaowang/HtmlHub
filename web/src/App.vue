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
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 64px;
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

.not-login {
  color: #888;
  font-size: 14px;
}

a.router-link-active {
  color: hsla(160, 100%, 30%, 1);
  font-weight: 600;
}
</style>