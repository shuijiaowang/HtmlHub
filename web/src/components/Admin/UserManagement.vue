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
      <el-table-column label="HTML限制" width="120">
        <template #default="{ row }">{{ formatSize(row.maxHtmlContentBytes) }}</template>
      </el-table-column>
      <el-table-column label="同步数据限制" width="130">
        <template #default="{ row }">{{ formatSize(row.maxHtmlDataBytes) }}</template>
      </el-table-column>
      <el-table-column label="未删除页面" width="130">
        <template #default="{ row }">{{ row.activeUploadCount || 0 }} / {{ row.maxActiveHtmlRecords }}</template>
      </el-table-column>
      <el-table-column label="累计页面" width="120">
        <template #default="{ row }">{{ row.totalUploadCount || 0 }} / {{ row.maxTotalHtmlRecords }}</template>
      </el-table-column>
      <el-table-column label="HTML占用" width="120">
        <template #default="{ row }">{{ formatSize(row.activeHtmlBytes) }}</template>
      </el-table-column>
      <el-table-column label="同步数据占用" width="130">
        <template #default="{ row }">{{ formatSize(row.htmlDataBytes) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">修改</el-button>
        </template>
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

    <el-dialog v-model="editVisible" title="修改用户" width="520px">
      <el-form :model="form" label-width="150px">
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" maxlength="20" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="留空表示不修改" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
            <el-option label="超级管理员" value="super_admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="单HTML限制(KB)">
          <el-input-number v-model="form.maxHtmlContentKb" :min="1" :step="1" />
        </el-form-item>
        <el-form-item label="单HTML同步数据(KB)">
          <el-input-number v-model="form.maxHtmlDataKb" :min="1" :step="1" />
        </el-form-item>
        <el-form-item label="未删除页面上限">
          <el-input-number v-model="form.maxActiveHtmlRecords" :min="1" :step="1" />
        </el-form-item>
        <el-form-item label="累计上传上限">
          <el-input-number v-model="form.maxTotalHtmlRecords" :min="1" :step="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submitEdit">保存</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminUsers, updateAdminUser } from '@/api/admin'

const KB = 1024

const users = ref([])
const pager = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const editVisible = ref(false)
const saving = ref(false)
const form = reactive({
  id: 0,
  nickname: '',
  email: '',
  password: '',
  role: 'user',
  maxHtmlContentKb: 500,
  maxHtmlDataKb: 100,
  maxActiveHtmlRecords: 5,
  maxTotalHtmlRecords: 10
})

const loadUsers = async () => {
  const res = await getAdminUsers({
    page: pager.page,
    pageSize: pager.pageSize
  })
  users.value = res.data?.list || []
  pager.total = Number(res.data?.total) || 0
}

const openEdit = (row) => {
  form.id = row.id
  form.nickname = row.nickname
  form.email = row.email
  form.password = ''
  form.role = row.role || 'user'
  form.maxHtmlContentKb = Math.round((Number(row.maxHtmlContentBytes) || 0) / KB) || 1
  form.maxHtmlDataKb = Math.round((Number(row.maxHtmlDataBytes) || 0) / KB) || 1
  form.maxActiveHtmlRecords = Number(row.maxActiveHtmlRecords) || 1
  form.maxTotalHtmlRecords = Number(row.maxTotalHtmlRecords) || 1
  editVisible.value = true
}

const submitEdit = async () => {
  saving.value = true
  try {
    await updateAdminUser(form.id, {
      nickname: form.nickname.trim(),
      email: form.email.trim(),
      password: form.password,
      role: form.role,
      maxHtmlContentBytes: Math.round(form.maxHtmlContentKb * KB),
      maxHtmlDataBytes: Math.round(form.maxHtmlDataKb * KB),
      maxActiveHtmlRecords: form.maxActiveHtmlRecords,
      maxTotalHtmlRecords: form.maxTotalHtmlRecords
    })
    ElMessage.success('已更新')
    editVisible.value = false
    await loadUsers()
  } finally {
    saving.value = false
  }
}

const formatRole = (role) => {
  const roleMap = {
    admin: '管理员',
    user: '普通用户',
    super_admin: '超级管理员'
  }
  return roleMap[role] || '普通用户'
}

const formatSize = (size) => {
  const value = Number(size) || 0
  if (value <= 0) return '0 B'
  if (value < 1024) return `${value} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
  return `${(value / (1024 * 1024)).toFixed(2)} MB`
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
