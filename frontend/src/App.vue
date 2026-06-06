<template>
  <div class="app-shell">
    <header>
      <h1>Project Planning Tracking</h1>
      <nav v-if="loggedIn" class="nav-links">
        <router-link v-for="item in navItems" :key="item.path" :to="item.path">{{ item.label }}</router-link>
      </nav>
    </header>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getUserRole, getToken, refreshUserProfile } from './services/auth'

const route = useRoute()
const token = ref(getToken())
const role = ref(getUserRole())

const loggedIn = computed(() => !!token.value)
const navItems = computed(() => {
  const items = [
    { path: '/dashboard', label: 'Dashboard' },
  ]

  if (role.value === 'admin') {
    items.push({ path: '/products', label: 'Products' })
    items.push({ path: '/plans', label: 'Production Plans' })
    items.push({ path: '/production', label: 'Production Entry' })
    items.push({ path: '/variance', label: 'Reports' })
    items.push({ path: '/users', label: 'Users' })
  } else if (role.value === 'planner') {
    items.push({ path: '/plans', label: 'Production Plans' })
    items.push({ path: '/variance', label: 'Reports' })
  } else if (role.value === 'operator') {
    items.push({ path: '/production', label: 'Production Entry' })
  }

  return items
})

watch(route, () => {
  token.value = getToken()
  role.value = getUserRole()
})

onMounted(async () => {
  if (token.value && !role.value) {
    try {
      await refreshUserProfile()
      role.value = getUserRole()
    } catch (err) {
      console.warn('Unable to refresh user profile', err)
    }
  }
})
</script>

<style>
body {
  margin: 0;
  font-family: system-ui, sans-serif;
  background: #f7fafc;
}

.app-shell {
  max-width: 1100px;
  margin: 0 auto;
  padding: 32px 16px;
}

header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  padding: 18px 20px;
  border-radius: 20px;
  background: linear-gradient(135deg, #ffffff, #f3f7fc);
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.06);
  margin-bottom: 28px;
}

h1 {
  margin: 0;
  font-size: 24px;
  letter-spacing: -0.02em;
  color: #111827;
}

.nav-links {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}

.nav-links a {
  color: #475569;
  text-decoration: none;
  font-weight: 600;
  padding: 10px 16px;
  border-radius: 14px;
  transition: all 0.2s ease;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  background: #eff6ff;
  color: #1d4ed8;
}
</style>
