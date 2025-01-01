import { createRouter, createWebHistory } from 'vue-router'
import { loadUserFromStorage, loadWorkspaceFromStorage } from '@/utils/storage'
import { canManageWorkspace } from '@/lib/permission'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: () => import('@/views/Landing.vue'),
    meta: {
      title: 'The Autonomous Merge Layer for Agent-Scale Delivery',
      description:
        'The autonomous merge layer for agent-scale delivery. Merge agent branches safely, resolve conflicts automatically, and keep shipping continuous.',
    },
  },
  {
    path: '/cases/customer-support',
    name: 'CaseCustomerSupport',
    component: () => import('@/views/cases/CustomerSupport.vue'),
  },
  {
    path: '/cases/sales',
    name: 'CaseSales',
    component: () => import('@/views/cases/Sales.vue'),
  },
  {
    path: '/cases/healthcare',
    name: 'CaseHealthcare',
    component: () => import('@/views/cases/Healthcare.vue'),
  },
  {
    path: '/cases/education',
    name: 'CaseEducation',
    component: () => import('@/views/cases/Education.vue'),
  },
  {
    path: '/cases/devtools',
    name: 'CaseDevTools',
    component: () => import('@/views/cases/DevTools.vue'),
  },
  {
    path: '/cases/e-commerce',
    name: 'CaseECommerce',
    component: () => import('@/views/cases/ECommerce.vue'),
  },
  {
    path: '/status',
    name: 'Status',
    component: () => import('@/views/Status.vue'),
    meta: { title: 'Status', description: 'Ziee platform status' }
  },
  {
    path: '/docs',
    name: 'Docs',
    component: () => import('@/views/Docs.vue'),
    meta: { title: 'Docs', description: 'Ziee documentation' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresGuest: true, title: 'Login', description: 'Login to platform' }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/ForgotPassword.vue'),
    meta: { requiresGuest: true, title: 'Forgot password', description: 'Request password reset' }
  },
  {
    path: '/reset-password/:token',
    name: 'ResetPassword',
    component: () => import('@/views/ResetPassword.vue'),
    meta: { requiresGuest: true, title: 'Reset password', description: 'Set new password' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresGuest: true, title: 'Register', description: 'Create an account' }
  },
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('@/views/Setup.vue'),
    meta: { requiresGuest: true, title: 'Setup', description: 'Setup platform' }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, title: 'Dashboard', description: 'Platform dashboard' }
  },
  {
    path: '/select-workspace',
    name: 'SwitchWorkspaces',
    component: () => import('@/views/SelectWorkspace.vue'),
    meta: { requiresAuth: true, title: 'Switch Workspaces', description: 'Choose or create a workspace' }
  },
  {
    path: '/invite/:token',
    name: 'Invite',
    component: () => import('@/views/Invite.vue'),
    meta: { requiresAuth: true, title: 'Workspace Invite', description: 'Accept or reject a workspace invite' }
  },
  {
    path: '/workspaces',
    name: 'Workspaces',
    component: () => import('@/views/Workspaces.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, title: 'Workspaces', description: 'Workspaces' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { requiresAuth: true, title: 'Profile', description: 'Your account profile' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/Settings.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, requiresWorkspaceAdmin: true, title: 'Settings', description: 'Settings' }
  },
  {
    path: '/members',
    name: 'Members',
    component: () => import('@/views/Members.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, requiresWorkspaceAdmin: true, title: 'Members', description: 'Members' }
  },
  {
    path: '/billing',
    name: 'Billing',
    component: () => import('@/views/Billing.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, requiresWorkspaceAdmin: true, title: 'Billing', description: 'Billing' }
  },
  {
    path: '/integrations',
    name: 'Integrations',
    component: () => import('@/views/Integrations.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, requiresWorkspaceAdmin: true, title: 'Integrations', description: 'Integrations' }
  },
  {
    path: '/audits',
    name: 'Audits',
    component: () => import('@/views/Audits.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, requiresWorkspaceAdmin: true, title: 'Audits', description: 'Workspace audit log' }
  },
  {
    path: '/knowledge',
    name: 'Knowledge',
    component: () => import('@/views/Knowledge.vue'),
    meta: { requiresAuth: true, requiresWorkspace: true, title: 'Knowledge', description: 'Knowledge' }
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: 'Page not found' }
  },
  {
    path: '/500',
    name: 'ServerError',
    component: () => import('@/views/ServerError.vue'),
    meta: { title: 'Server error' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFoundCatchAll',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: 'Page not found' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return { top: 0 }
  },
})

function isSafeRedirect(redirect) {
  return typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//')
}

// Navigation guard - reads auth state directly from storage
router.beforeEach((to, from, next) => {
  const currentUser = loadUserFromStorage()
  const currentWorkspace = loadWorkspaceFromStorage()
  const isAuthenticated = !!currentUser

  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else if (to.meta.requiresGuest && isAuthenticated) {
    const redirect = to.query.redirect
    next(isSafeRedirect(redirect) ? redirect : '/select-workspace')
  } else if (to.meta.requiresWorkspace && !currentWorkspace) {
    next('/select-workspace')
  } else if (to.meta.requiresWorkspaceAdmin && !canManageWorkspace(currentUser, currentWorkspace)) {
    next('/404')
  } else {
    next()
  }
})

export default router
