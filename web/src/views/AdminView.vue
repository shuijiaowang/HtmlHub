<script setup>
import { computed, ref } from 'vue'
import UserManagement from '@/components/Admin/UserManagement.vue'
import HtmlManagement from '@/components/Admin/HtmlManagement.vue'

const activeMenu = ref('users')
const currentTitle = computed(() => (activeMenu.value === 'users' ? '用户管理' : 'HTML管理'))
</script>

<template>
  <el-container class="admin-layout">
    <el-aside width="220px" class="admin-aside">
      <div class="admin-logo">HtmlHub Admin</div>
      <el-menu :default-active="activeMenu" class="admin-menu" @select="activeMenu = $event">
        <el-menu-item index="users">用户管理</el-menu-item>
        <el-menu-item index="html">HTML管理</el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="admin-header">
        <h1>{{ currentTitle }}</h1>
      </el-header>
      <el-main class="admin-main">
        <UserManagement v-if="activeMenu === 'users'" />
        <HtmlManagement v-else />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: #f5f7fa;
}

.admin-aside {
  background: #fff;
  border-right: 1px solid #e4e7ed;
}

.admin-logo {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  font-size: 18px;
  font-weight: 700;
  border-bottom: 1px solid #e4e7ed;
}

.admin-menu {
  border-right: none;
}

.admin-header {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}

.admin-header h1 {
  margin: 0;
  font-size: 22px;
}

.admin-main {
  padding: 24px;
}
</style>