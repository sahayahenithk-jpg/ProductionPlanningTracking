<template>
  <div class="page">
    <section class="card">
      <div class="header-row">
        <div>
          <h2>User Management</h2>
          <p class="subtitle">Create, update, and manage system users and roles.</p>
        </div>
      </div>

      <form @submit.prevent="saveUser" class="form">
        <div class="grid">
          <input v-model="form.name" placeholder="Full Name" required />
          <input v-model="form.email" type="email" placeholder="Email" required />
          <select v-model="form.role" required>
            <option value="">Select Role</option>
            <option value="admin">Admin</option>
            <option value="planner">Production Planner</option>
            <option value="operator">Operator</option>
          </select>
          <input v-model="form.password" type="password" placeholder="Password" :required="!editId" />
        </div>
        <div class="buttons">
          <button type="submit" class="primary">
            {{ editId ? 'Update User' : 'Create User' }}
          </button>
          <button type="button" class="secondary" @click="resetForm">Cancel</button>
        </div>
      </form>

      <div class="table">
        <div class="row header">
          <span>Name</span>
          <span>Email</span>
          <span>Role</span>
          <span>Actions</span>
        </div>

        <div v-for="user in users" :key="user.id" class="row">
          <span>{{ user.name }}</span>
          <span>{{ user.email }}</span>
          <span>{{ user.role }}</span>
          <div class="actions">
            <button class="edit" @click="startEdit(user)">Edit</button>
            <button class="delete" @click="removeUser(user.id)">Delete</button>
          </div>
        </div>
      </div>

      <p v-if="error" class="error">{{ error }}</p>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'

const users = ref([])
const editId = ref(null)
const error = ref('')
const form = ref({
  name: '',
  email: '',
  password: '',
  role: '',
})

const loadUsers = async () => {
  error.value = ''
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
    name: '',
    email: '',
    password: '',
    role: '',
  }
  error.value = ''
}

const saveUser = async () => {
  error.value = ''
  try {
    if (editId.value) {
      const payload = {
        name: form.value.name,
        email: form.value.email,
        role: form.value.role,
      }
      if (form.value.password) {
        payload.password = form.value.password
      }
      await api.put(`/users/${editId.value}`, payload)
    } else {
      await api.post('/users', {
        name: form.value.name,
        email: form.value.email,
        password: form.value.password,
        role: form.value.role,
      })
    }
    resetForm()
    await loadUsers()
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to save user'
  }
}

const startEdit = (user) => {
  editId.value = user.id
  form.value = {
    name: user.name,
    email: user.email,
    password: '',
    role: user.role,
  }
}

const removeUser = async (id) => {
  error.value = ''
  try {
    await api.delete(`/users/${id}`)
    await loadUsers()
  } catch (err) {
    error.value = err.response?.data?.error || 'Unable to delete user'
  }
}

onMounted(loadUsers)
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
  margin-bottom: 18px;
}

.header-row h2 {
  font-size: 24px;
  margin: 0;
}

.subtitle {
  color: #666;
  margin-top: 6px;
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
.grid select {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
}

.buttons {
  display: flex;
  gap: 10px;
  margin-top: 16px;
}

.primary,
.secondary,
.edit,
.delete {
  border: none;
  border-radius: 10px;
  padding: 10px 14px;
  cursor: pointer;
}

.primary {
  background: #4a90e2;
  color: white;
}

.secondary {
  background: #f3f4f6;
  color: #334155;
}

.table {
  border-top: 1px solid #e2e8f0;
}

.row {
  display: grid;
  grid-template-columns: 2fr 2fr 1fr 1fr;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid #e2e8f0;
  align-items: center;
}

.row.header {
  font-weight: 700;
  text-transform: uppercase;
  color: #475569;
}

.actions {
  display: flex;
  gap: 8px;
}

.edit {
  background: #fbbf24;
  color: white;
}

.delete {
  background: #ef4444;
  color: white;
}

.error {
  color: #b91c1c;
  margin-top: 14px;
}
</style>
