import {createRouter, createWebHistory} from 'vue-router'
import LoginView from "@/views/LoginView.vue";
import RegisterView from "@/views/RegisterView.vue";
import { useUserStore } from '@/stores/user'

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
            component: LoginView
        },
        {
            path: '/register',
            name: 'register',
            component: RegisterView
        },
        {
            path: '/home',
            name: 'home',
            component: () => import('../views/HomeView.vue'),
            meta: {requiresAuth: true}  // 需要登录才能访问
        },
        {
            path: '/home/upload',
            name: 'home-upload',
            component: () => import('../views/UploadView.vue'),
            meta: {requiresAuth: true}
        },
        {
            path: '/home/manage',
            name: 'home-manage',
            component: () => import('../views/ManageView.vue'),
            meta: {requiresAuth: true}
        },
        {
            path: '/home/showcase',
            name: 'home-showcase',
            component: () => import('../views/ShowcaseView.vue'),
            meta: {requiresAuth: true}
        },
        {
            path: '/admin',
            name: 'admin',
            component: () => import('../views/AdminView.vue'),
            meta: {requiresAuth: true, requiresAdmin: true}
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
export default router
