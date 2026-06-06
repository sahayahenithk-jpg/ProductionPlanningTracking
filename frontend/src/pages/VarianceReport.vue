<template>
  <div class="page">
    <section class="card">
      <div class="header-row">
        <div>
          <h2>Variance Report</h2>
          <p class="subtitle">View daily production performance against planned quantity.</p>
        </div>
        <button class="logout" @click="logout">Logout</button>
      </div>

      <div class="panel">
        <div class="table">
          <div class="row header">
            <span>Plan #</span>
            <span>Product</span>
            <span>Prod Date</span>
            <span>Produced</span>
            <span>Planned</span>
            <span>Difference</span>
            <span>Achievement %</span>
          </div>

          <div v-for="row in rows" :key="row.entryId" class="row">
            <span>{{ row.planNumber }}</span>
            <span>{{ row.productName }}</span>
            <span>{{ row.productionDate }}</span>
            <span>{{ row.producedQuantity }}</span>
            <span>{{ row.plannedQuantity }}</span>
            <span>{{ row.difference }}</span>
            <span>{{ row.achievementPct.toFixed(2) }}%</span>
          </div>
        </div>
      </div>

      <p v-if="error" class="error">{{ error }}</p>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const router = useRouter()
const rows = ref([])
const error = ref('')

const loadReport = async () => {
  try {
    const response = await api.get('/reports/variance')
    rows.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load variance report'
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
  await loadReport()
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
  background: #fff;
  border-radius: 16px;
  padding: 28px;
  box-shadow: 0 12px 35px rgba(0, 0, 0, 0.12);
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.header-row h2 {
  font-size: 24px;
  color: #222;
  margin: 0;
}

.subtitle {
  font-size: 14px;
  color: #666;
  margin-top: 6px;
}

.logout {
  background: #ff4d4f;
  color: white;
  border: none;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
}

.table {
  width: 100%;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e6ecf5;
}

.row {
  display: grid;
  grid-template-columns: 1.2fr 1.4fr 1fr 0.8fr 0.8fr 0.9fr 1fr;
  gap: 12px;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
}

.row.header {
  background: #f1f5f9;
  font-weight: 700;
}

.error {
  color: #dc2626;
  margin-top: 16px;
}
</style>
