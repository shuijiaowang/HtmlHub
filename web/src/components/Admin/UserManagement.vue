<template>
  <section class="admin-card">
    <div class="card-header">
      <h2>用户管理</h2>
      <el-button type="primary" @click="loadUsers">刷新</el-button>
    </div>

    <el-table :data="users" border>
      <el-table-column prop="nickname" label="昵称" min-width="160" />
      <el-table-column prop="email" label="邮箱" min-width="220" />
      <el-table-column prop="role" label="角色" width="140">
        <template #default="{ row }">{{ formatRole(row.role) }}</template>
      </el-table-column>
    </el-table>

    <div class="pagination-row">
      <el-pagination
        v-model:current-page="pager.page"
        v-model:page-size="pager.pageSize"
        layout="total, sizes, prev, pager, next"
        :page-sizes="[10, 20, 50]"
        :total="pager.total"
        @size-change="loadUsers"
        @current-change="loadUsers"
      />
    </div>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { getAdminUsers } from '@/api/admin'

const users = ref([])
const pager = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadUsers = async () => {
  const res = await getAdminUsers({
    page: pager.page,
    pageSize: pager.pageSize
  })
  users.value = res.data?.list || []
  pager.total = Number(res.data?.total) || 0
}

const formatRole = (role) => {
  const roleMap = {
    admin: '管理员',
    user: '普通用户',
    super_admin: '超级管理员'
  }
  return roleMap[role] || '普通用户'
}

onMounted(loadUsers)
</script>

<style scoped>
.admin-card {
  background: #fff;
  border-radius: 10px;
  padding: 20px;
}

.card-header,
.pagination-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header {
  margin-bottom: 16px;
}

.card-header h2 {
  margin: 0;
  font-size: 20px;
}

.pagination-row {
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
