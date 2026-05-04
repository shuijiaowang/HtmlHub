import {createRouter, createWebHistory} from 'vue-router'
import LoginView from "@/views/LoginView.vue";
import RegisterView from "@/views/RegisterView.vue";
import { useUserStore } from '@/stores/user'

const DEFAULT_TITLE = 'HtmlHub - 轻量在线 HTML 托管平台'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            redirect: '/login'  // 默认重定向到登录页
        },
        {
            path: '/login',
            name: 'login',
            component: LoginView,
            meta: { title: '登录' }
        },
        {
            path: '/register',
            name: 'register',
            component: RegisterView,
            meta: { title: '注册' }
        },
        {
            path: '/home',
            name: 'home',
            component: () => import('../views/HomeView.vue'),
            meta: { requiresAuth: true, title: '主页' }
        },
        {
            path: '/home/upload',
            name: 'home-upload',
            component: () => import('../views/UploadView.vue'),
            meta: { requiresAuth: true, title: '上传 HTML' }
        },
        {
            path: '/home/manage',
            name: 'home-manage',
            component: () => import('../views/ManageView.vue'),
            meta: { requiresAuth: true, title: '个人 HTML 管理' }
        },
        {
            path: '/home/showcase',
            name: 'home-showcase',
            component: () => import('../views/ShowcaseView.vue'),
            meta: { requiresAuth: true, title: '展示页' }
        },
        {
            path: '/admin',
            name: 'admin',
            component: () => import('../views/AdminView.vue'),
            meta: { requiresAuth: true, requiresAdmin: true, title: '管理后台' }
        }

    ]
})
// 路由守卫：未登录用户只能访问登录页
router.beforeEach((to, from, next) => {
    const isAuthenticated = !!localStorage.getItem('token')
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login')
    } else if (to.meta.requiresAdmin) {
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
