<template>
  <div class="app-shell">
    <header>
      <h1>Project Planning Tracking</h1>
      <nav v-if="loggedIn" class="nav-links">
        <router-link to="/dashboard">Dashboard</router-link>
        <router-link to="/products">Products</router-link>
        <router-link to="/plans">Production Plans</router-link>
      </nav>
    </header>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const loggedIn = ref(!!localStorage.getItem('token'))

watch(route, () => {
  loggedIn.value = !!localStorage.getItem('token')
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
  text-align: center;
  margin-bottom: 24px;
}

h1 {
  margin: 0;
  font-size: 28px;
}

.nav-links {
  margin-top: 16px;
  display: flex;
  justify-content: center;
  gap: 18px;
  flex-wrap: wrap;
}

.nav-links a {
  color: #2563eb;
  text-decoration: none;
  font-weight: 600;
}

.nav-links a.router-link-active {
  color: #111827;
}
</style>
