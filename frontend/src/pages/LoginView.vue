<template>
  <div class="page">
    <div class="card">
      <header class="card-header">
        <h2>Welcome back</h2>
        <p class="subtitle">Login to continue to your dashboard</p>
      </header>

      <form @submit.prevent="submitLogin" class="login-form">
        <div class="form-group">
          <label for="username">Username</label>
          <input 
            id="username"
            v-model="username" 
            type="text" 
            placeholder="Enter your username" 
            required 
            autocomplete="username"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input 
            id="password"
            v-model="password" 
            type="password" 
            placeholder="••••••••" 
            required 
            autocomplete="current-password"
          />
        </div>

        <button type="submit" :disabled="loading" class="submit-btn">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? "Logging in..." : "Sign in" }}
        </button>

        <p v-if="error" class="error-msg" role="alert">
          {{ error }}
        </p>

        <p class="footer-note">
          Don't have an account?
          <router-link to="/register">Register</router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const router = useRouter()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const submitLogin = async () => {
  error.value = ''
  loading.value = true

  try {
    const response = await api.post('/login', {
      username: username.value,
      password: password.value,
    })

    localStorage.setItem('token', response.data.token)
    router.push('/products')
  } catch (err) {
    error.value = err.response?.data?.error || 'Invalid username or password'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f5f5f5;
  font-family: Arial, sans-serif;
}

.card {
  width: 100%;
  max-width: 400px;
  background: #ffffff;
  padding: 24px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.card-header h2 {
  text-align: center;
  margin-bottom: 5px;
  font-size: 22px;
}

.subtitle {
  text-align: center;
  font-size: 13px;
  color: #666;
  margin-bottom: 20px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

label {
  font-size: 13px;
  color: #333;
}

input {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  font-size: 14px;
}

input:focus {
  outline: none;
  border-color: #4a90e2;
}

.submit-btn {
  padding: 10px;
  background: #4a90e2;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.submit-btn:disabled {
  background: #a0c4f0;
  cursor: not-allowed;
}

.error-msg {
  color: red;
  font-size: 13px;
  text-align: center;
}

.footer-note {
  text-align: center;
  font-size: 13px;
  margin-top: 10px;
}

.footer-note a {
  color: #4a90e2;
  text-decoration: none;
}
</style>