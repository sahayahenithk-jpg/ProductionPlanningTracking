import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../pages/LoginView.vue'
import RegisterView from '../pages/RegisterView.vue'
import ProductMaster from '../pages/ProductMaster.vue'
import DashboardView from '../pages/DashboardView.vue'
import ProductionPlanning from '../pages/ProductionPlanning.vue'
import ProductionEntry from '../pages/ProductionEntry.vue'
import VarianceReport from '../pages/VarianceReport.vue'
import UserManagement from '../pages/UserManagement.vue'
import { getUserRole, getToken } from '../services/auth'

const routes = [
  {
    path: '/',
    name: 'Login',
    component: LoginView,
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView,
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
  },
  {
    path: '/products',
    name: 'Products',
    component: ProductMaster,
  },
  {
    path: '/plans',
    name: 'ProductionPlanning',
    component: ProductionPlanning,
  },
  {
    path: '/production',
    name: 'ProductionEntry',
    component: ProductionEntry,
  },
  {
    path: '/variance',
    name: 'ProductionReport',
    component: VarianceReport,
  },
  {
    path: '/users',
    name: 'UserManagement',
    component: UserManagement,
  },
]

const roleRoutes = {
  Products: ['admin'],
  UserManagement: ['admin'],
  ProductionPlanning: ['planner', 'admin'],
  ProductionEntry: ['operator', 'admin'],
  ProductionReport: ['planner', 'admin'],
  Dashboard: ['admin', 'planner', 'operator'],
}

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const token = getToken()
  const role = getUserRole()

  if (!token && to.name !== 'Login' && to.name !== 'Register') {
    return next({ name: 'Login' })
  }

  if (token && (to.name === 'Login' || to.name === 'Register')) {
    return next({ name: 'Dashboard' })
  }

  const allowedRoles = roleRoutes[to.name]
  if (allowedRoles && !allowedRoles.includes(role)) {
    return next({ name: 'Dashboard' })
  }

  return next()
})

export default router
