<script setup>
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()

const isHomeRoute = computed(() => route.path.startsWith('/home'))
const isLoggedIn = computed(() => !!userStore.token)
const displayNickname = computed(() => userStore.userInfo.nickname || '已登录用户')

const logout = async () => {
  await userStore.logout()
}
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
          <router-link to="/register">注册</router-link>
          <router-link to="/login">登录</router-link>
        </template>
      </div>
    </header>
    <RouterView />
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

a.router-link-active {
  color: hsla(160, 100%, 30%, 1);
  font-weight: 600;
}
</style>