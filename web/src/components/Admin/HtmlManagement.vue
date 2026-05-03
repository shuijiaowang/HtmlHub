<template>
  <section class="admin-card">
    <div class="card-header">
      <h2>HTML管理</h2>
      <el-button type="primary" @click="loadRecords">搜索</el-button>
    </div>

    <el-form :model="filters" class="filter-form" label-width="86px">
      <el-form-item label="昵称">
        <el-input v-model="filters.nickname" clearable placeholder="按昵称搜索" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="filters.email" clearable placeholder="按邮箱搜索" />
      </el-form-item>
      <el-form-item label="子域名">
        <el-input v-model="filters.subdomain" clearable placeholder="按子域名搜索" />
      </el-form-item>
      <el-form-item label="公开私密">
        <el-select v-model="filters.visibility" clearable placeholder="全部">
          <el-option label="公开" value="public" />
          <el-option label="私密" value="private" />
        </el-select>
      </el-form-item>
      <el-form-item label="审核状态">
        <el-select v-model="filters.approvalStatus" clearable placeholder="全部">
          <el-option label="未审核" value="pending" />
          <el-option label="通过" value="approved" />
          <el-option label="拒绝" value="rejected" />
        </el-select>
      </el-form-item>
    </el-form>

    <el-table :data="records" border>
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="subdomain" label="子域名" min-width="140" />
      <el-table-column prop="visibility" label="公开私密" width="110">
        <template #default="{ row }">{{ formatVisibility(row.visibility) }}</template>
      </el-table-column>
      <el-table-column prop="approvalStatus" label="审核状态" width="110">
        <template #default="{ row }">{{ formatApprovalStatus(row.approvalStatus) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="400" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openDetail(row)">查看</el-button>
          <el-dropdown @command="(status) => updateApproval(row, status)">
            <el-button link type="primary">审核</el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="pending">未审核</el-dropdown-item>
                <el-dropdown-item command="approved">通过</el-dropdown-item>
                <el-dropdown-item command="rejected">拒绝</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button link type="primary" @click="editSubdomain(row)">改子域名</el-button>
          <el-button link type="danger" @click="removeRecord(row)">删除</el-button>
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
        @size-change="loadRecords"
        @current-change="loadRecords"
      />
    </div>

    <el-dialog v-model="detailVisible" title="HTML详情" width="70%">
      <div v-if="detail" class="detail-meta">
        <span>文件：{{ detail.fileName }}</span>
        <span>子域名：{{ detail.subdomain }}</span>
        <span>用户：{{ detail.nickname }} / {{ detail.email }}</span>
      </div>
      <pre class="html-preview">{{ detail?.htmlContent || '' }}</pre>
    </el-dialog>
  </section>
</template>

<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import {
  deleteAdminHtml,
  getAdminHtmlDetail,
  getAdminHtmlList,
  updateAdminHtmlApproval,
  updateAdminHtmlSubdomain
} from '@/api/admin'

const records = ref([])
const detail = ref(null)
const detailVisible = ref(false)
const filters = reactive({
  nickname: '',
  email: '',
  subdomain: '',
  visibility: '',
  approvalStatus: ''
})
const pager = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadRecords = async () => {
  const res = await getAdminHtmlList({
    ...filters,
    page: pager.page,
    pageSize: pager.pageSize
  })
  records.value = res.data?.list || []
  pager.total = Number(res.data?.total) || 0
}

const openDetail = async (row) => {
  const res = await getAdminHtmlDetail(row.id)
  detail.value = res.data
  detailVisible.value = true
}

const updateApproval = async (row, approvalStatus) => {
  const res = await updateAdminHtmlApproval(row.id, approvalStatus)
  row.approvalStatus = res.data?.approvalStatus || approvalStatus
  ElMessage.success('审核状态已更新')
}

const editSubdomain = async (row) => {
  try {
    const value = await ElMessageBox.prompt('请输入新的子域名前缀', '修改子域名', {
      inputValue: row.subdomain,
      inputPattern: /^[a-z0-9][a-z0-9-]{2,63}$/,
      inputErrorMessage: '仅支持小写字母、数字和中划线，长度3-64'
    })
    const subdomain = value.value
    const res = await updateAdminHtmlSubdomain(row.id, subdomain)
    row.subdomain = res.data?.subdomain || subdomain
    ElMessage.success('子域名已更新')
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') throw error
  }
}

const removeRecord = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除「${row.fileName}」吗？`, '删除确认', {
      type: 'warning'
    })
    await deleteAdminHtml(row.id)
    await loadRecords()
    ElMessage.success('已删除')
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') throw error
  }
}

const formatVisibility = (visibility) => {
  return visibility === 'public' ? '公开' : '私密'
}

const formatApprovalStatus = (status) => {
  const statusMap = {
    pending: '未审核',
    approved: '通过',
    rejected: '拒绝'
  }
  return statusMap[status] || '未审核'
}

onMounted(loadRecords)
</script>

<style scoped>
.admin-card {
  background: #fff;
  border-radius: 10px;
  padding: 20px;
}

.card-header,
.pagination-row,
.detail-meta {
  display: flex;
  align-items: center;
}

.card-header {
  justify-content: space-between;
  margin-bottom: 16px;
}

.card-header h2 {
  margin: 0;
  font-size: 20px;
}

.filter-form {
  display: grid;
  grid-template-columns: repeat(3, minmax(220px, 1fr));
  gap: 0 12px;
}

.pagination-row {
  justify-content: flex-end;
  margin-top: 16px;
}

.detail-meta {
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 12px;
  color: #606266;
}

.html-preview {
  max-height: 520px;
  overflow: auto;
  padding: 12px;
  border-radius: 8px;
  background: #f6f8fa;
  white-space: pre-wrap;
}
</style>
