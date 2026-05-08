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
      <div class="nav-inner">
        <div class="brand">
          <router-link class="brand-link" to="/home">
            <span class="brand-mark" aria-hidden="true" />
            <span class="brand-text">HtmlHub</span>
          </router-link>
        </div>
        <nav class="left-links" aria-label="主导航">
          <router-link to="/home">首页</router-link>
          <router-link to="/home/upload">上传</router-link>
          <router-link to="/home/manage">管理</router-link>
          <router-link to="/home/showcase">展示</router-link>
        </nav>
        <div class="right-actions">
          <template v-if="isLoggedIn">
            <span class="nickname" :title="displayNickname">{{ displayNickname }}</span>
            <button type="button" class="ghost-btn nav-logout" @click="logout">退出</button>
          </template>
          <template v-else>
            <button type="button" class="ghost-btn nav-register" @click="openRegister">注册</button>
            <button type="button" class="primary-btn nav-login" @click="openLogin">登录</button>
          </template>
        </div>
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
  min-height: 64px;
  border-bottom: 1px solid var(--hh-border);
  background: color-mix(in srgb, var(--hh-surface) 88%, transparent);
  backdrop-filter: blur(10px);
}

.nav-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  min-height: 64px;
  padding: 10px 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.brand {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
}

.brand-link {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 12px;
  color: var(--hh-text);
}

.brand-link:hover {
  text-decoration: none;
  background: rgb(var(--hh-brand-rgb) / 0.10);
}

.brand-mark {
  width: 14px;
  height: 14px;
  border-radius: 6px;
  background: radial-gradient(
    circle at 30% 30%,
    rgb(var(--hh-brand-2-rgb) / 0.90) 0%,
    var(--hh-brand) 60%,
    rgb(var(--hh-brand-rgb) / 0.90) 140%
  );
  box-shadow: 0 8px 18px rgb(var(--hh-brand-rgb) / 0.22);
}

.brand-text {
  font-weight: 800;
  letter-spacing: 0.2px;
}

.left-links,
.right-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.left-links {
  flex: 1;
  min-width: 0;
  overflow-x: auto;
  padding: 2px;
  scrollbar-width: none;
}

.left-links::-webkit-scrollbar {
  display: none;
}

.left-links a {
  white-space: nowrap;
  padding: 8px 12px;
  border-radius: 10px;
  color: var(--hh-text-2);
  border: 1px solid transparent;
}

.left-links a:hover {
  text-decoration: none;
  color: var(--hh-text);
  background: rgb(var(--hh-brand-rgb) / 0.10);
}

.nickname {
  color: var(--hh-text);
  font-weight: 600;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ghost-btn,
.primary-btn {
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.22);
  cursor: pointer;
  font-size: 14px;
  padding: 8px 12px;
  border-radius: 12px;
  transition: transform 0.08s ease, background-color 0.15s ease, border-color 0.15s ease;
}

.ghost-btn {
  background: rgb(var(--hh-brand-rgb) / 0.10);
  color: color-mix(in srgb, var(--hh-brand) 88%, #000 0%);
}

.primary-btn {
  background: color-mix(in srgb, var(--hh-brand) 86%, #000 0%);
  color: white;
  border-color: rgb(var(--hh-brand-rgb) / 0.30);
}

.ghost-btn:hover {
  background: rgb(var(--hh-brand-rgb) / 0.16);
}

.primary-btn:hover {
  background: color-mix(in srgb, var(--hh-brand) 92%, #000 0%);
}

.ghost-btn:active,
.primary-btn:active {
  transform: translateY(1px);
}

a.router-link-active {
  color: var(--hh-text);
  font-weight: 600;
  background: rgb(var(--hh-brand-rgb) / 0.16);
  border-color: rgb(var(--hh-brand-rgb) / 0.22);
}

@media (max-width: 640px) {
  .top-nav {
    min-height: 48px;
  }

  .nav-inner {
    padding: 6px 8px;
    gap: 6px;
    min-height: 48px;
  }

  .brand {
    display: none;
  }

  .left-links {
    flex: 1 1 0;
    min-width: 0;
    overflow-x: visible;
    justify-content: space-between;
    gap: 2px;
    padding: 0;
  }

  .left-links a {
    flex: 1 1 0;
    min-width: 0;
    padding: 6px 2px;
    border-radius: 8px;
    font-size: 12px;
    text-align: center;
    line-height: 1.25;
  }

  .right-actions {
    flex: 0 0 auto;
    flex-shrink: 0;
    gap: 6px;
  }

  .nav-register {
    display: none;
  }

  .nickname {
    display: none;
  }

  .nav-login,
  .nav-logout {
    padding: 6px 8px;
    font-size: 12px;
    border-radius: 10px;
  }

  .ghost-btn:not(.nav-register),
  .primary-btn {
    padding: 6px 8px;
  }
}
</style>