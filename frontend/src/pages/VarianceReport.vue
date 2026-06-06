<template>
  <div class="page">
    <section class="card">
      <div class="header-row">
        <div>
          <h2>Production Report</h2>
          <p class="subtitle">Filter production performance and compare planned vs actual output.</p>
        </div>
        <button class="logout" @click="logout">Logout</button>
      </div>

      <form @submit.prevent="loadReport" class="filter-form">
        <div class="grid">
          <select v-model.number="filters.productId">
            <option value="">All Products</option>
            <option v-for="product in products" :key="product.productId" :value="product.productId">
              {{ product.productName }}
            </option>
          </select>
          <input type="date" v-model="filters.startDate" />
          <input type="date" v-model="filters.endDate" />
          <input type="month" v-model="filters.month" />
        </div>

        <div class="buttons">
          <button type="submit" class="primary">Apply Filters</button>
          <button type="button" class="secondary" @click="resetFilters">Reset</button>
        </div>
      </form>

      <div v-if="error" class="error">{{ error }}</div>

      <div class="table">
        <div class="row header">
          <span>Plan #</span>
          <span>Product</span>
          <span>Date</span>
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
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const router = useRouter()
const products = ref([])
const rows = ref([])
const error = ref('')
const filters = ref({ productId: null, startDate: '', endDate: '', month: '' })

const loadProducts = async () => {
  try {
    const response = await api.get('/products')
    products.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load products'
  }
}

const loadReport = async () => {
  error.value = ''
  try {
    const params = {}
    if (filters.value.productId) params.productId = filters.value.productId
    if (filters.value.startDate) params.startDate = filters.value.startDate
    if (filters.value.endDate) params.endDate = filters.value.endDate
    if (filters.value.month) params.month = filters.value.month

    const response = await api.get('/reports/variance', { params })
    rows.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load report'
  }
}

const resetFilters = async () => {
  filters.value = { productId: null, startDate: '', endDate: '', month: '' }
  await loadReport()
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
  await loadProducts()
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
  margin-bottom: 20px;
}

.header-row h2 {
  font-size: 24px;
  margin: 0;
  color: #222;
}

.subtitle {
  margin: 4px 0 0;
  color: #666;
}

.filter-form {
  background: #f9fbff;
  border: 1px solid #e6ecf5;
  border-radius: 14px;
  padding: 18px;
  margin-bottom: 22px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(160px, 1fr));
  gap: 14px;
}

.grid select,
.grid input {
  width: 100%;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 10px;
  font-size: 14px;
}

.buttons {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
}

.primary,
.secondary {
  border: none;
  border-radius: 10px;
  padding: 12px 18px;
  cursor: pointer;
}

.primary {
  background: #4a90e2;
  color: white;
}

.secondary {
  background: #e5e7eb;
  color: #111827;
}

.table {
  width: 100%;
  border-radius: 14px;
  overflow: hidden;
  border: 1px solid #e6ecf5;
}

.row {
  display: grid;
  grid-template-columns: 1.4fr 1.4fr 1fr 0.9fr 0.9fr 0.9fr 1fr;
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
  margin-top: 16px;
  color: #dc2626;
}

.logout {
  background: #ff4d4f;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 10px;
  cursor: pointer;
}
</style>
