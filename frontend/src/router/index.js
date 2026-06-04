import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../pages/LoginView.vue'
import RegisterView from '../pages/RegisterView.vue'
import ProductMaster from '../pages/ProductMaster.vue'

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
    path: '/products',
    name: 'Products',
    component: ProductMaster,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
