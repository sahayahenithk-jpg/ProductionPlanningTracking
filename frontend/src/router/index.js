import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../pages/LoginView.vue'
import RegisterView from '../pages/RegisterView.vue'
import ProductMaster from '../pages/ProductMaster.vue'
import DashboardView from '../pages/DashboardView.vue'
import ProductionPlanning from '../pages/ProductionPlanning.vue'
import ProductionEntry from '../pages/ProductionEntry.vue'
import VarianceReport from '../pages/VarianceReport.vue'

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
    name: 'VarianceReport',
    component: VarianceReport,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
