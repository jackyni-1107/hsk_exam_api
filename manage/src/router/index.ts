import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/index.vue'),
        },
        {
          path: 'system/user',
          name: 'User',
          component: () => import('@/views/system/user/index.vue'),
        },
        {
          path: 'system/member',
          name: 'Member',
          component: () => import('@/views/system/member/index.vue'),
          meta: { title: '会员管理' },
        },
        {
          path: 'system/client-user',
          redirect: '/system/member',
        },
        {
          path: 'system/role',
          name: 'Role',
          component: () => import('@/views/system/role/index.vue'),
        },
        {
          path: 'system/menu',
          name: 'Menu',
          component: () => import('@/views/system/menu/index.vue'),
        },
        {
          path: 'audit-log/request',
          name: 'AuditLog',
          component: () => import('@/views/audit-log/request/index.vue'),
        },
        {
          path: 'audit-log/login-log',
          name: 'LoginLog',
          component: () => import('@/views/audit-log/login-log/index.vue'),
        },
        {
          path: 'audit-log/exception-log',
          name: 'ExceptionLog',
          component: () => import('@/views/audit-log/exception-log/index.vue'),
        },
        {
          path: 'audit-log/security-event-log',
          name: 'SecurityEventLog',
          component: () => import('@/views/audit-log/security-event-log/index.vue'),
        },
        {
          path: 'system/notification',
          name: 'Notification',
          component: () => import('@/views/system/notification/index.vue'),
        },
        {
          path: 'system/file',
          name: 'File',
          component: () => import('@/views/system/file/index.vue'),
        },
        {
          path: 'system/config',
          name: 'Config',
          component: () => import('@/views/system/config/index.vue'),
        },
        {
          path: 'system/dict',
          name: 'Dict',
          component: () => import('@/views/system/dict/index.vue'),
        },
        {
          path: 'task',
          name: 'Task',
          component: () => import('@/views/task/index.vue'),
        },
        {
          path: 'task/log',
          name: 'TaskLog',
          component: () => import('@/views/task/log.vue'),
        },
        {
          path: 'exam/paper',
          name: 'ExamPaper',
          meta: { title: '试卷管理' },
          component: () => import('@/views/exam/paper/index.vue'),
        },
        {
          path: 'exam/batch',
          name: 'ExamBatch',
          meta: { title: '考试批次' },
          component: () => import('@/views/exam/batch/index.vue'),
        },
        {
          path: 'exam/result',
          name: 'ExamResult',
          meta: { title: '考试结果' },
          component: () => import('@/views/exam/result/index.vue'),
        },
        {
          path: 'exam/hls-debug',
          name: 'ExamHlsDebug',
          meta: { title: 'HLS 联调（临时）' },
          component: () => import('@/views/exam/hls-debug/index.vue'),
        },
        // 兼容旧菜单路径，避免 router-view 无匹配导致空白
        {
          path: 'system/exam/paper',
          redirect: '/exam/paper',
        },
        {
          path: 'system/exam/batch',
          redirect: '/exam/batch',
        },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  if (to.meta.public) {
    if (userStore.token) {
      next({ path: '/' })
    } else {
      next()
    }
  } else {
    if (!userStore.token) {
      next({ path: '/login', query: { redirect: to.fullPath } })
    } else {
      next()
    }
  }
})

export default router
