import {createRouter, createWebHistory} from 'vue-router'
import { useUserStore } from '@/stores/user'

const DEFAULT_TITLE = 'HtmlHub - 轻量在线 HTML 托管平台'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            redirect: '/home'  // 默认重定向到主页
        },
        {
            path: '/login',
            name: 'login',
            redirect: '/home?auth=login',
            meta: { title: '登录' }
        },
        {
            path: '/register',
            name: 'register',
            redirect: '/home?auth=register',
            meta: { title: '注册' }
        },
        {
            path: '/home',
            name: 'home',
            component: () => import('../views/HomeView.vue'),
            meta: { title: '主页' }
        },
        {
            path: '/home/upload',
            name: 'home-upload',
            component: () => import('../views/UploadView.vue'),
            meta: { title: '上传 HTML' }
        },
        {
            path: '/home/manage',
            name: 'home-manage',
            component: () => import('../views/ManageView.vue'),
            meta: { title: '个人 HTML 管理' }
        },
        {
            path: '/home/showcase',
            name: 'home-showcase',
            component: () => import('../views/ShowcaseView.vue'),
            meta: { title: '展示页' }
        },
        {
            path: '/admin',
            name: 'admin',
            component: () => import('../views/AdminView.vue'),
            meta: { requiresAuth: true, requiresAdmin: true, title: '管理后台' }
        }

    ]
})
// 路由守卫：仅管理员页面需要强制校验
router.beforeEach((to, from, next) => {
    if (to.meta.requiresAdmin) {
        const userStore = useUserStore()
        const role = userStore.userInfo?.role
        if (role === 'admin' || role === 'super_admin') {
            next()
        } else {
            next('/home')
        }
    } else {
        next()
    }
})

router.afterEach((to) => {
    const piece = to.meta.title
    document.title = piece ? `${piece} | HtmlHub` : DEFAULT_TITLE
})

export default router
