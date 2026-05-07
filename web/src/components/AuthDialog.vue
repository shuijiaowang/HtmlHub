<template>
  <el-dialog
    v-model="visible"
    :title="mode === 'login' ? '登录' : '注册'"
    width="420px"
    destroy-on-close
    :close-on-click-modal="false"
  >
    <el-form :model="form" label-position="top" @submit.prevent>
      <el-form-item v-if="mode === 'register'" label="昵称">
        <el-input v-model="form.nickname" placeholder="请输入昵称" autocomplete="nickname" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="form.email" placeholder="请输入邮箱" autocomplete="email" />
      </el-form-item>
      <el-form-item label="密码">
        <el-input
          v-model="form.password"
          type="password"
          show-password
          placeholder="请输入密码"
          autocomplete="current-password"
          @keyup.enter="onSubmit"
        />
      </el-form-item>

      <div class="actions">
        <el-button @click="visible = false">关闭</el-button>
        <el-button type="primary" :loading="submitting" @click="onSubmit">
          {{ mode === 'login' ? '登录' : '注册' }}
        </el-button>
      </div>

      <div class="switch-row">
        <template v-if="mode === 'login'">
          还没有账号？
          <el-link type="primary" :underline="false" @click="switchMode('register')">去注册</el-link>
        </template>
        <template v-else>
          已有账号？
          <el-link type="primary" :underline="false" @click="switchMode('login')">去登录</el-link>
        </template>
      </div>
    </el-form>
  </el-dialog>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { register } from '@/api/user'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const visible = computed({
  get: () => userStore.authDialogVisible,
  set: (val) => {
    if (!val) userStore.closeAuthDialog()
  }
})

const mode = computed(() => userStore.authDialogMode)

const submitting = ref(false)
const form = reactive({
  nickname: '',
  email: '',
  password: ''
})

const resetForm = () => {
  form.nickname = ''
  form.email = ''
  form.password = ''
}

watch(
  () => userStore.authDialogVisible,
  (v) => {
    if (v) resetForm()
  }
)

const switchMode = (nextMode) => {
  userStore.openAuthDialog(nextMode)
}

const onSubmit = async () => {
  const email = form.email.trim()
  const password = form.password

  if (!email || !password) {
    ElMessage.warning('请输入邮箱和密码')
    return
  }

  submitting.value = true
  try {
    if (mode.value === 'login') {
      const ok = await userStore.loginIn({ email, password })
      if (!ok) ElMessage.error('登录失败，请检查邮箱和密码')
      return
    }

    const nickname = form.nickname.trim()
    if (!nickname) {
      ElMessage.warning('请输入昵称')
      return
    }
    await register({ nickname, email, password })
    ElMessage.success('注册成功，请登录')
    userStore.openAuthDialog('login')
  } catch (e) {
    console.error(e)
    ElMessage.error(mode.value === 'login' ? '登录失败' : '注册失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 4px;
}

.switch-row {
  margin-top: 12px;
  font-size: 13px;
  color: #666;
  text-align: right;
}
</style>
