<template>
  <section class="admin-card">
    <div class="card-header">
      <h2>同步数据</h2>
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
    </el-form>

    <el-table :data="records" border>
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="subdomain" label="子域名" min-width="140" />
      <el-table-column label="HTML链接" min-width="240">
        <template #default="{ row }">
          <el-link type="primary" :href="buildPublicUrl(row.subdomain)" target="_blank" rel="noopener noreferrer">
            {{ buildPublicUrl(row.subdomain) }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column label="数据体积" width="120">
        <template #default="{ row }">{{ formatSize(row.dataBytes) }}</template>
      </el-table-column>
      <el-table-column label="更新时间" width="180">
        <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
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
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { getAdminHtmlDataList } from '@/api/admin'

const htmlPublicHost = import.meta.env.VITE_HTML_PUBLIC_HOST || 'localhost:7789'

const records = ref([])
const filters = reactive({
  nickname: '',
  email: '',
  subdomain: ''
})
const pager = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadRecords = async () => {
  const res = await getAdminHtmlDataList({
    ...filters,
    page: pager.page,
    pageSize: pager.pageSize
  })
  records.value = res.data?.list || []
  pager.total = Number(res.data?.total) || 0
}

const formatSize = (size) => {
  const value = Number(size) || 0
  if (value <= 0) return '0 B'
  if (value < 1024) return `${value} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
  return `${(value / (1024 * 1024)).toFixed(2)} MB`
}

const formatDate = (value) => {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

const buildPublicUrl = (subdomain) => {
  const host = String(htmlPublicHost).replace(/^https?:\/\//, '').replace(/\/$/, '')
  const protocol = host.includes('localhost') || host.includes('127.0.0.1') ? 'http:' : window.location.protocol
  return `${protocol}//${subdomain}.${host}`
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
.pagination-row {
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
</style>
