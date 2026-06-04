<template>
  <div class="page">
    <section class="card">
      <h2>Dashboard</h2>
      <p class="subtitle">Welcome back 👋</p>

      <!-- Loading -->
      <p v-if="!user && !error" class="loading">Loading profile...</p>

      <!-- Profile -->
      <div v-if="user" class="profile">
        <div class="avatar">
          {{ user.name?.charAt(0)?.toUpperCase() }}
        </div>

        <div class="info">
          <p><span>Name</span> {{ user.name }}</p>
          <p><span>Email</span> {{ user.email }}</p>
        </div>

        <button @click="logout" class="logout">
          Logout
        </button>
      </div>

      <!-- Error -->
      <p v-if="error" class="error">{{ error }}</p>
    </section>
  </div>
</template>



<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const router = useRouter()
const user = ref(null)
const error = ref('')

const loadProfile = async () => {
  try {
    const response = await api.get('/user')
    user.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load profile'
    logout()
  }
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/')
}

onMounted(() => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/')
    return
  }
  loadProfile()
})
</script>

<style scoped>
/* Page layout */
.page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #fef2f2, #f1f5f9);
  padding: 20px;
}

/* Card */
.card {
  width: 100%;
  max-width: 460px;
  background: white;
  padding: 28px;
  border-radius: 16px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.1);
}

/* Titles */
h2 {
  font-size: 24px;
  margin-bottom: 6px;
  color: #111827;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 20px;
}

/* Loading */
.loading {
  color: #6b7280;
  font-size: 14px;
}

/* Profile layout */
.profile {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}

/* Avatar */
.avatar {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: #ef4444;
  color: white;
  font-size: 26px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Info block */
.info {
  width: 100%;
  background: #f9fafb;
  padding: 14px;
  border-radius: 12px;
}

.info p {
  margin: 8px 0;
  font-size: 14px;
  color: #111827;
}

.info span {
  display: block;
  font-size: 12px;
  color: #6b7280;
}

/* Logout button */
.logout {
  width: 100%;
  padding: 12px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s;
}

.logout:hover {
  background: #dc2626;
}

/* Error */
.error {
  margin-top: 14px;
  text-align: center;
  color: #b91c1c;
  font-size: 14px;
}
</style>