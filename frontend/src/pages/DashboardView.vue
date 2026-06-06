<template>
  <div class="page">
    <section class="card">
      <div class="header-row">
        <div>
          <h2>Dashboard</h2>
          <p class="subtitle">Summary cards for products, plans, and today’s production.</p>
        </div>
        <button class="logout" @click="logout">Logout</button>
      </div>

      <div v-if="error" class="error">{{ error }}</div>

      <div v-if="!error" class="cards-grid">
        <div class="summary-card">
          <div class="card-label">Total Products</div>
          <div class="card-value">{{ overview.totalProducts }}</div>
        </div>
        <div class="summary-card">
          <div class="card-label">Total Plans</div>
          <div class="card-value">{{ overview.totalPlans }}</div>
        </div>
        <div class="summary-card">
          <div class="card-label">Today's Production</div>
          <div class="card-value">{{ today.productionTotal }}</div>
        </div>
        <div class="summary-card">
          <div class="card-label">Achievement %</div>
          <div class="card-value">{{ today.achievementPct.toFixed(2) }}%</div>
        </div>
      </div>

      <div class="section">
        <h3>Your profile</h3>
        <div class="profile-card">
          <p><strong>Name:</strong> {{ user?.name || '-' }}</p>
          <p><strong>Email:</strong> {{ user?.email || '-' }}</p>
        </div>
      </div>

      <div class="actions-row">
        <router-link to="/products" class="action-button">Manage Products</router-link>
        <router-link to="/plans" class="action-button">Create Plans</router-link>
        <router-link to="/variance" class="action-button">View Production Report</router-link>
      </div>
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
const overview = ref({ totalProducts: 0, totalPlans: 0 })
const today = ref({ productionTotal: 0, achievementPct: 0 })

const formatToday = () => new Date().toISOString().split('T')[0]

const loadProfile = async () => {
  try {
    const response = await api.get('/user')
    user.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load profile'
  }
}

const loadOverview = async () => {
  try {
    const response = await api.get('/reports/summary')
    overview.value = {
      totalProducts: response.data.totalProducts,
      totalPlans: response.data.totalPlans,
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load overview'
  }
}

const loadTodayStats = async () => {
  try {
    const response = await api.get('/reports/summary', {
      params: {
        startDate: formatToday(),
        endDate: formatToday(),
      },
    })
    today.value = {
      productionTotal: response.data.productionTotal,
      achievementPct: response.data.achievementPct,
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load today summary'
  }
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/')
}

onMounted(async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/')
    return
  }
  await loadProfile()
  await loadOverview()
  await loadTodayStats()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 30px;
  display: flex;
  justify-content: center;
  background: linear-gradient(135deg, #eef2f3, #d9e4f5);
  font-family: Arial, sans-serif;
}

.card {
  width: 100%;
  max-width: 1100px;
  background: white;
  border-radius: 16px;
  padding: 28px;
  box-shadow: 0 12px 35px rgba(0, 0, 0, 0.12);
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.header-row h2 {
  font-size: 26px;
  margin: 0;
  color: #111827;
}

.subtitle {
  color: #6b7280;
  margin: 4px 0 0;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.summary-card {
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 20px;
}

.card-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
}

.card-value {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
}

.section {
  margin-top: 20px;
}

.profile-card {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 18px;
}

.profile-card p {
  margin: 8px 0;
  color: #111827;
}

.actions-row {
  margin-top: 28px;
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
}

.action-button {
  background: #2563eb;
  color: white;
  padding: 12px 18px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
}

.logout {
  background: #ef4444;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 10px;
  cursor: pointer;
}

.error {
  color: #b91c1c;
  margin-bottom: 16px;
}
</style>
