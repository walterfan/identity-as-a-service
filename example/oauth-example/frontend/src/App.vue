<script setup>
import { ref } from 'vue'
import axios from 'axios'

const user = ref(null)
const error = ref(null)

const loginWithGitHub = () => {
  window.location.href = 'http://localhost:8000/auth/github'
}

const checkAuth = async () => {
  try {
    // 检查 URL 中是否有 token (回调后)
    const urlParams = new URLSearchParams(window.location.search)
    const accessToken = urlParams.get('access_token')
    
    if (accessToken) {
      // 存储 token
      localStorage.setItem('github_access_token', accessToken)
      // 清除 URL 中的 token
      window.history.replaceState({}, document.title, window.location.pathname)
    }
    
    const token = localStorage.getItem('github_access_token')
    if (token) {
      const response = await axios.get('http://localhost:8000/api/user', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      user.value = response.data
    }
  } catch (err) {
    error.value = err.response?.data?.detail || err.message
    localStorage.removeItem('github_access_token')
  }
}

const logout = () => {
  localStorage.removeItem('github_access_token')
  user.value = null
}

// 初始化检查认证状态
checkAuth()
</script>

<template>
  <div>
    <h1>GitHub OAuth Demo</h1>
    
    <div v-if="error" class="error">
      Error: {{ error }}
    </div>
    
    <div v-if="user">
      <h2>Welcome, {{ user.login }}!</h2>
      <img :src="user.avatar_url" width="100" style="border-radius: 50%;">
      <p>Name: {{ user.name || 'Not provided' }}</p>
      <p>GitHub: <a :href="user.html_url" target="_blank">{{ user.html_url }}</a></p>
      <button @click="logout">Logout</button>
    </div>
    
    <button v-else @click="loginWithGitHub">Login with GitHub</button>
  </div>
</template>

<style>
.error {
  color: red;
}
</style>