<template>
  <div class="page">
    <section class="card">
      <div class="header-row">
        <h2>Product Master</h2>
        <button class="logout" @click="logout">Logout</button>
      </div>
      <p class="subtitle">Manage your product catalog</p>

      <!-- FORM -->
      <form @submit.prevent="saveProduct" class="form">
        <div class="grid">
          <input v-model="form.productCode" placeholder="Product Code" required />
          <input v-model="form.productName" placeholder="Product Name" required />
          <input v-model="form.unit" placeholder="Unit (kg, pcs...)" />
          <input v-model="form.status" placeholder="Status" />
          <input v-model="form.description" placeholder="Description" class="full" />
        </div>

        <div class="buttons">
          <button type="submit" class="primary">
            {{ editId ? 'Update Product' : 'Add Product' }}
          </button>
          <button type="button" class="secondary" @click="resetForm">
            Cancel
          </button>
        </div>
      </form>

      <!-- ERROR -->
      <p v-if="error" class="error">{{ error }}</p>

      <!-- TABLE -->
      <div class="table">
        <div class="row header">
          <span>Code</span>
          <span>Name</span>
          <span>Unit</span>
          <span>Status</span>
          <span>Actions</span>
        </div>

        <div v-for="p in products" :key="p.productId" class="row">
          <span class="code">{{ p.productCode }}</span>
          <span>{{ p.productName }}</span>
          <span>{{ p.unit || '-' }}</span>
          <span>
            <span class="badge">{{ p.status || 'Active' }}</span>
          </span>

          <div class="actions">
            <button class="edit" @click="startEdit(p)">Edit</button>
            <button class="delete" @click="remove(p.productId)">Delete</button>
          </div>
        </div>
      </div>
    </section>
  </div>

  
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const products = ref([])
const form = ref({
  productCode: '',
  productName: '',
  unit: '',
  description: '',
  status: ''
})
const editId = ref(null)
const error = ref('')
const router = useRouter()

const load = async () => {
  try {
    const res = await api.get('/products')
    products.value = res.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to load products'
  }
}

const saveProduct = async () => {
  error.value = ''
  try {
    if (editId.value) {
      await api.put(`/products/${editId.value}`, form.value)
    } else {
      await api.post('/products', form.value)
    }
    resetForm()
    load()
  } catch (err) {
    error.value = err.response?.data?.error || 'Save failed'
  }
}

const startEdit = (p) => {
  editId.value = p.productId
  form.value = {
    productCode: p.productCode,
    productName: p.productName,
    unit: p.unit,
    description: p.description,
    status: p.status
  }
}

const remove = async (id) => {
  try {
    await api.delete(`/products/${id}`)
    load()
  } catch (err) {
    error.value = err.response?.data?.error || 'Delete failed'
  }
}

const resetForm = () => {
  editId.value = null
  form.value = {
    productCode: '',
    productName: '',
    unit: '',
    description: '',
    status: ''
  }
}

onMounted(load)

const logout = () => {
  localStorage.removeItem('token')
  router.push('/')
}
</script>

<style scoped>
/* PAGE */
.page {
  min-height: 100vh;
  padding: 30px;
  display: flex;
  justify-content: center;
  background: linear-gradient(135deg, #eef2f3, #d9e4f5);
  font-family: Arial, sans-serif;
}

/* MAIN CARD */
.card {
  width: 100%;
  max-width: 1100px;
  background: #fff;
  border-radius: 16px;
  padding: 28px;
  box-shadow: 0 12px 35px rgba(0, 0, 0, 0.12);
}

/* HEADER ROW */
.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.header-row h2 {
  font-size: 24px;
  color: #222;
}

.subtitle {
  font-size: 14px;
  color: #666;
  margin-bottom: 20px;
}

/* LOGOUT */
.logout {
  background: #ff4d4f;
  color: white;
  border: none;
  padding: 8px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: 0.3s;
}

.logout:hover {
  background: #d9363e;
}

/* FORM */
.form {
  background: #f9fbff;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 20px;
  border: 1px solid #e6ecf5;
}

/* GRID INPUTS */
.grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.grid input {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  transition: 0.2s;
}

.grid input:focus {
  border-color: #4a90e2;
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.15);
  outline: none;
}

/* FULL WIDTH FIELD */
.full {
  grid-column: span 2;
}

/* BUTTONS */
.buttons {
  display: flex;
  gap: 10px;
  margin-top: 12px;
}

.primary {
  background: #4a90e2;
  color: white;
  border: none;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: 0.3s;
}

.primary:hover {
  background: #357bd8;
}

.secondary {
  background: #eaeaea;
  border: none;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
}

.secondary:hover {
  background: #d6d6d6;
}

/* ERROR */
.error {
  background: #ffe5e5;
  color: #d93025;
  padding: 10px;
  border-radius: 8px;
  margin-bottom: 15px;
  font-size: 13px;
}

/* TABLE */
.table {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #eee;
}

/* ROW */
.row {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr 1fr 1fr;
  padding: 12px;
  align-items: center;
  border-bottom: 1px solid #eee;
  background: white;
  transition: 0.2s;
}

.row:hover {
  background: #f7faff;
}

/* HEADER ROW */
.row.header {
  background: #f1f5ff;
  font-weight: bold;
  color: #333;
}

/* CODE */
.code {
  font-weight: 600;
  color: #333;
}

/* BADGE */
.badge {
  background: #e6f0ff;
  color: #2b6cb0;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
}

/* ACTION BUTTONS */
.actions {
  display: flex;
  gap: 8px;
}

.edit {
  background: #4a90e2;
  color: white;
  border: none;
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
}

.edit:hover {
  background: #357bd8;
}

.delete {
  background: #ff4d4f;
  color: white;
  border: none;
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
}

.delete:hover {
  background: #d9363e;
}

/* RESPONSIVE */
@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }

  .full {
    grid-column: span 1;
  }

  .row {
    grid-template-columns: 1fr;
    gap: 6px;
  }

  .actions {
    justify-content: flex-start;
  }
}
</style>