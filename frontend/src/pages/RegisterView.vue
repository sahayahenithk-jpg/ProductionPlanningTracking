<template>
  <div class="page">
    <section class="card">
      <h2>Create Account</h2>
      <p class="subtitle">Join us and start managing your workflow</p>

      <form @submit.prevent="submitRegister">
        <label>
          Name
          <input v-model="name" type="text" placeholder="John Doe" required />
        </label>

        <label>
          Email
          <input v-model="email" type="email" placeholder="you@example.com" required />
        </label>

        <label>
          Password
          <input v-model="password" type="password" placeholder="••••••••" required minlength="6" />
        </label>

        <button type="submit" :disabled="loading">
          {{ loading ? "Creating..." : "Create Account" }}
        </button>

        <p class="note">
          Already registered?
          <router-link to="/">Login</router-link>
        </p>

        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
      </form>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '../services/api'

const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)

const submitRegister = async () => {
  error.value = ''
  success.value = ''
  loading.value = true

  try {
    await api.post('/auth/register', {
      name: name.value,
      email: email.value,
      password: password.value,
    })

    success.value = 'Registration successful. Please login.'

    name.value = ''
    email.value = ''
    password.value = ''
  } catch (err) {
    error.value = err.response?.data?.error || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>

<style>
/* PAGE BACKGROUND */
.page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #eef2f3, #d9e4f5);
  padding: 20px;
  font-family: Arial, sans-serif;
}

/* CARD */
.card {
  width: 100%;
  max-width: 420px;
  background: #ffffff;
  padding: 32px;
  border-radius: 14px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  animation: fadeIn 0.4s ease-in-out;
}

/* HEADER */
.card h2 {
  margin-bottom: 6px;
  font-size: 24px;
  color: #222;
  text-align: center;
}

.subtitle {
  text-align: center;
  font-size: 14px;
  color: #666;
  margin-bottom: 20px;
}

/* FORM */
form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* LABEL + INPUT */
label {
  font-size: 13px;
  color: #444;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

input {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  outline: none;
  transition: 0.2s;
  font-size: 14px;
}

input:focus {
  border-color: #4a90e2;
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.15);
}

/* BUTTON */
button {
  margin-top: 10px;
  padding: 12px;
  background: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  cursor: pointer;
  transition: 0.3s ease;
}

button:hover {
  background: #357bd8;
}

button:disabled {
  background: #9bbce6;
  cursor: not-allowed;
}

/* TEXT LINKS */
.note {
  text-align: center;
  font-size: 13px;
  color: #555;
  margin-top: 10px;
}

.note a {
  color: #4a90e2;
  text-decoration: none;
}

.note a:hover {
  text-decoration: underline;
}

/* ERROR & SUCCESS */
.error {
  background: #ffe5e5;
  color: #d93025;
  padding: 10px;
  border-radius: 8px;
  font-size: 13px;
  text-align: center;
}

.success {
  background: #e6f9ed;
  color: #1e8e3e;
  padding: 10px;
  border-radius: 8px;
  font-size: 13px;
  text-align: center;
}

/* ANIMATION */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

</style>