import { defineStore } from 'pinia'
import { useStorage } from '@vueuse/core'
import { login } from '@/api/user'
import { ElLoading, ElMessage } from 'element-plus'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
    const token = useStorage('token', '')
    const userInfo = useStorage('userInfo', {
        nickname: '',
        email: '',
        uuid: '',
        role: ''
    })

    // 登录/注册弹窗（全局）
    const authDialogVisible = useStorage('authDialogVisible', false)
    const authDialogMode = useStorage('authDialogMode', 'login') // 'login' | 'register'

    const setUserInfo = (info) => {
        userInfo.value = {
            nickname: info.nickname || userInfo.value.nickname,
            email: info.email || userInfo.value.email,
            uuid: info.uuid || userInfo.value.uuid,
            role: info.role || userInfo.value.role
        }
    }

    // 设置token
    const setToken = (newToken) => {
        token.value = newToken
    }

    const clearUserState = () => {
        token.value = ''
        userInfo.value = {
            nickname: '',
            email: '',
            uuid: '',
            role: ''
        }
    }

    const openAuthDialog = (mode = 'login') => {
        authDialogMode.value = mode
        authDialogVisible.value = true
    }

    const closeAuthDialog = () => {
        authDialogVisible.value = false
    }

    const loginIn = async (loginForm) => {
        const loading = ElLoading.service({
            fullscreen: true,
            text: '登录中...'
        })

        try {
            if (!loginForm.email || !loginForm.password) {
                ElMessage.warning('请输入邮箱和密码')
                return false
            }

            const res = await login(loginForm)
            if (res.code !== 0) {
                return false
            }

            setToken(res.data.token)
            setUserInfo({
                nickname: res.data.nickname,
                email: res.data.email,
                uuid: res.data.uuid,
                role: res.data.role
            })

            const targetRoute = 'home'
            if (router.hasRoute(targetRoute)) {
                await router.replace({ name: targetRoute })
            } else {
                ElMessage.error('首页路由不存在，请联系管理员')
                return false
            }

            ElMessage.success('登录成功')
            closeAuthDialog()
            return true
        } catch (error) {
            console.error('登录失败:', error)
            return false
        } finally {
            loading.close()
        }
    }

    // 退出登录
    const logout = async () => {
        clearUserState()
        if (router.hasRoute('home')) {
            await router.replace({ name: 'home' })
        } else {
            await router.replace({ path: '/home' })
        }
        ElMessage.success('已退出登录')
    }

    return {
        token,
        userInfo,
        authDialogVisible,
        authDialogMode,
        setUserInfo,
        setToken,
        clearUserState,
        openAuthDialog,
        closeAuthDialog,
        loginIn,
        logout
    }
})