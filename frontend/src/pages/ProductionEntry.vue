<template>
  <div class="page">
    <section class="card">
      <div class="header-row">
        <div>
          <h2>Production Entry</h2>
          <p class="subtitle">Record daily production performance and review recent entries.</p>
        </div>
        <button class="logout" @click="logout">Logout</button>
      </div>

      <div class="tabs">
        <button :class="{ active: view === 'create' }" @click="view = 'create'">Create Entry</button>
        <button :class="{ active: view === 'list' }" @click="view = 'list'">View Entries</button>
      </div>

      <div v-if="view === 'create'" class="panel">
        <form @submit.prevent="saveEntry" class="form">
          <div class="grid">
            <select v-model.number="form.planId" required>
              <option value="" disabled>Select Plan</option>
              <option v-for="plan in plans" :key="plan.planId" :value="plan.planId">
                {{ plan.planNumber }} - {{ plan.product?.productName || 'Unknown' }}
              </option>
            </select>
            <input type="date" v-model="form.productionDate" required />
            <input type="number" min="1" v-model.number="form.producedQuantity" placeholder="Produced Quantity" required />
            <input v-model="form.shift" placeholder="Shift" />
            <textarea v-model="form.remarks" placeholder="Remarks" class="full"></textarea>
          </div>
          <div class="buttons">
            <button type="submit" class="primary">Save Entry</button>
            <button type="button" class="secondary" @click="resetForm">Reset</button>
          </div>
        </form>
      </div>

      <div v-if="view === 'list'" class="panel">
        <div class="table">
          <div class="row header">
            <span>Plan #</span>
            <span>Product</span>
            <span>Date</span>
            <span>Qty</span>
            <span>Shift</span>
            <span>Remarks</span>
          </div>

          <div v-for="entry in entries" :key="entry.entryId" class="row">
            <span>{{ entry.plan?.planNumber || 'N/A' }}</span>
            <span>{{ entry.plan?.product?.productName || 'N/A' }}</span>
            <span>{{ formatDate(entry.productionDate) }}</span>
            <span>{{ entry.producedQuantity }}</span>
            <span>{{ entry.shift || '-' }}</span>
            <span>{{ entry.remarks || '-' }}</span>
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
const view = ref('create')
const plans = ref([])
const entries = ref([])
const error = ref('')
const form = ref({
  planId: null,
  productionDate: '',
  producedQuantity: 0,
  shift: '',
  remarks: '',
})

const loadPlans = async () => {
  try {
    const response = await api.get('/plans')
    plans.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load plans'
  }
}

const loadEntries = async () => {
  try {
    const response = await api.get('/production')
    entries.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load production entries'
  }
}

const saveEntry = async () => {
  error.value = ''
  try {
    await api.post('/production', {
      planId: form.value.planId,
      productionDate: form.value.productionDate,
      producedQuantity: form.value.producedQuantity,
      shift: form.value.shift,
      remarks: form.value.remarks,
    })
    resetForm()
    view.value = 'list'
    await loadEntries()
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to save production entry'
  }
}

const resetForm = () => {
  form.value = {
    planId: null,
    productionDate: '',
    producedQuantity: 0,
    shift: '',
    remarks: '',
  }
  error.value = ''
}

const formatDate = (value) => {
  if (!value) return ''
  const date = new Date(value)
  return date.toISOString().split('T')[0]
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
  await loadPlans()
  await loadEntries()
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

.tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.tabs button {
  padding: 12px 18px;
  border: 1px solid #d1d5db;
  border-radius: 10px;
  background: white;
  cursor: pointer;
}

.tabs button.active {
  background: #4a90e2;
  color: white;
  border-color: #4a90e2;
}

.form {
  background: #f9fbff;
  padding: 18px;
  border-radius: 12px;
  border: 1px solid #e6ecf5;
  margin-bottom: 20px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14px;
}

.grid input,
.grid select,
.grid textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
}

.grid textarea {
  resize: vertical;
  min-height: 100px;
}

.full {
  grid-column: span 2;
}

.buttons {
  display: flex;
  gap: 10px;
  margin-top: 16px;
}

.primary {
  background: #4a90e2;
  color: white;
  border: none;
  padding: 12px 18px;
  border-radius: 10px;
  cursor: pointer;
}

.secondary {
  background: #eaeaea;
  border: none;
  padding: 12px 18px;
  border-radius: 10px;
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
  grid-template-columns: 1.2fr 1.4fr 1fr 0.8fr 0.8fr 1.5fr;
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
