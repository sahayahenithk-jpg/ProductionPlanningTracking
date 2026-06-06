<template>
  <div class="page">
    <section class="card">
      <div class="header-row">
        <div>
          <h2>Production Planning</h2>
          <p class="subtitle">Create and manage production plans for your products.</p>
        </div>
        <button class="logout" @click="logout">Logout</button>
      </div>

      <div class="tabs">
        <button v-if="canEditPlans" :class="{ active: view === 'create' }" @click="view = 'create'">Create Plan</button>
        <button :class="{ active: view === 'list' }" @click="view = 'list'">View / Edit Plans</button>
      </div>

      <div v-if="view === 'create' && canEditPlans" class="panel">
        <form @submit.prevent="savePlan" class="form">
          <div class="grid">
            <input v-model="form.planNumber" placeholder="Plan Number" required />
            <select v-model.number="form.productId" required>
              <option value="" disabled>Select Product</option>
              <option v-for="product in products" :key="product.productId" :value="product.productId">
                {{ product.productName }}
              </option>
            </select>
            <input type="date" v-model="form.planDate" required />
            <input type="number" min="1" v-model.number="form.plannedQuantity" placeholder="Planned Quantity" required />
            <select v-model.number="form.assignedTo">
              <option value="" disabled>Select operator to assign</option>
              <option v-for="user in operators" :key="user.id" :value="user.id">
                {{ user.name }} ({{ user.email }})
              </option>
            </select>
            <input v-model="form.status" placeholder="Status" required />
            <textarea v-model="form.remarks" placeholder="Remarks" class="full"></textarea>
          </div>

          <div class="buttons">
            <button type="submit" class="primary">
              {{ editId ? 'Update Plan' : 'Create Plan' }}
            </button>
            <button type="button" class="secondary" @click="resetForm">Reset</button>
          </div>
        </form>
      </div>
      <div v-else-if="view === 'create'" class="panel">
        <p class="info">Only Production Planners and Admins can create or edit plans.</p>
      </div>

      <div v-if="view === 'list'" class="panel">
        <div class="table">
          <div class="row header">
            <span>Plan #</span>
            <span>Product</span>
            <span>Date</span>
            <span>Quantity</span>
            <span>Assigned</span>
            <span>Status</span>
            <span v-if="canEditPlans">Actions</span>
          </div>

          <div v-for="plan in plans" :key="plan.planId" class="row">
            <span>{{ plan.planNumber }}</span>
            <span>{{ plan.product?.productName || 'Unknown' }}</span>
            <span>{{ formatDate(plan.planDate) }}</span>
            <span>{{ plan.plannedQuantity }}</span>
            <span>{{ plan.assignedUser?.name || '-' }}</span>
            <span><span class="badge">{{ plan.status || 'Pending' }}</span></span>
            <div class="actions" v-if="canEditPlans">
              <button class="edit" @click="startEdit(plan)">Edit</button>
            </div>
          </div>
        </div>
      </div>

      <p v-if="error" class="error">{{ error }}</p>
    </section>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'
import { getUserRole } from '../services/auth'

const router = useRouter()
const view = ref('create')
const products = ref([])
const plans = ref([])
const users = ref([])
const editId = ref(null)
const error = ref('')
const role = ref(getUserRole())
const form = ref({
  planNumber: '',
  productId: null,
  planDate: '',
  plannedQuantity: 0,
  assignedTo: null,
  remarks: '',
  status: '',
})

const canEditPlans = computed(() => role.value === 'planner' || role.value === 'admin')
const operators = computed(() => users.value.filter((user) => user.role === 'operator'))

const loadProducts = async () => {
  try {
    const response = await api.get('/products')
    products.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load products'
  }
}

const loadPlans = async () => {
  try {
    const response = await api.get('/plans')
    plans.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load plans'
  }
}

const loadUsers = async () => {
  if (!canEditPlans.value) {
    return
  }
  try {
    const response = await api.get('/users')
    users.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load users'
  }
}

const resetForm = () => {
  editId.value = null
  form.value = {
    planNumber: '',
    productId: null,
    planDate: '',
    plannedQuantity: 0,
    assignedTo: null,
    remarks: '',
    status: '',
  }
  view.value = 'create'
  error.value = ''
}

const savePlan = async () => {
  error.value = ''
  const payload = {
    planNumber: form.value.planNumber,
    productId: form.value.productId,
    planDate: form.value.planDate,
    plannedQuantity: form.value.plannedQuantity,
    assignedTo: form.value.assignedTo,
    remarks: form.value.remarks,
    status: form.value.status,
  }

  try {
    if (editId.value) {
      await api.put(`/plans/${editId.value}`, payload)
    } else {
      await api.post('/plans', payload)
    }
    resetForm()
    await loadPlans()
    view.value = 'list'
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to save plan'
  }
}

const startEdit = (plan) => {
  editId.value = plan.planId
  form.value = {
    planNumber: plan.planNumber,
    productId: plan.productId,
    planDate: formatDate(plan.planDate),
    plannedQuantity: plan.plannedQuantity,
    assignedTo: plan.assignedTo || null,
    remarks: plan.remarks,
    status: plan.status,
  }
  view.value = 'create'
}

const formatDate = (value) => {
  if (!value) return ''
  const date = new Date(value)
  return date.toISOString().split('T')[0]
}

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('userRole')
  router.push('/')
}

onMounted(async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/')
    return
  }
  await loadProducts()
  await loadPlans()
  await loadUsers()
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
  padding: 0;
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
  padding: 14px 16px;
  border: 1px solid #cbd5e1;
  border-radius: 14px;
  font-size: 15px;
  background: #ffffff;
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
  background: #1d4ed8;
  color: white;
  border: none;
  padding: 12px 18px;
  border-radius: 10px;
  cursor: pointer;
}

 .secondary {
  background: #e2e8f0;
  border: none;
  padding: 12px 18px;
  border-radius: 10px;
  cursor: pointer;
}

.table {
  border-radius: 12px;
  border: 1px solid #e6ecf5;
  overflow: hidden;
}

.row {
  display: grid;
  grid-template-columns: 1.3fr 1.7fr 1fr 1fr 1fr 1fr;
  gap: 10px;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid #eef2f5;
}

.row.header {
  background: #f3f7fb;
  font-weight: 600;
}

.row:last-child {
  border-bottom: none;
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px 10px;
  background: #eef4ff;
  border-radius: 999px;
  color: #1e40af;
  font-size: 13px;
}

.actions {
  display: flex;
  gap: 8px;
}

.edit {
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
}

.error {
  margin-top: 16px;
  color: #b91c1c;
}

@media (max-width: 900px) {
  .row {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .row.header {
    display: none;
  }

  .row {
    padding: 12px;
  }
}
</style>
