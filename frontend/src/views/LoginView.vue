<script setup>
// LoginView.vue —— 登录 / 注册（一个页面两个模式切换）

import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login, register } from '../api'

const router = useRouter()

const mode = ref('login')      // 'login' 登录 | 'register' 注册
const username = ref('')
const password = ref('')
const msg = ref('')
const msgType = ref('error')   // 'error' | 'ok'

// 切换登录/注册
function switchMode(m) {
  mode.value = m
  msg.value = ''
}

// 提交
async function submit() {
  msg.value = ''
  // 简单前端校验（后端也会再校验一遍）
  if (username.value.length < 3 || password.value.length < 6) {
    msg.value = '用户名至少3位，密码至少6位'
    return
  }
  try {
    if (mode.value === 'login') {
      const res = await login({ username: username.value, password: password.value })
      // 登录成功：把 token 和用户名存到 localStorage，全局都能读到
      localStorage.setItem('token', res.token)
      localStorage.setItem('username', username.value)
      router.push('/')
    } else {
      await register({ username: username.value, password: password.value })
      msg.value = '注册成功，请登录'
      msgType.value = 'ok'
      mode.value = 'login' // 自动切到登录
    }
  } catch (e) {
    msg.value = e.response?.data?.msg || '操作失败'
    msgType.value = 'error'
  }
}
</script>

<template>
  <div style="max-width: 380px; margin: 60px auto">
    <h1 style="font-size: 24px; text-align: center; margin-bottom: 24px">
      {{ mode === 'login' ? '登录' : '注册' }}
    </h1>

    <!-- 模式切换 -->
    <div style="display: flex; gap: 12px; justify-content: center; margin-bottom: 20px">
      <button
        class="btn"
        :class="{ 'btn-primary': mode === 'login' }"
        @click="switchMode('login')"
      >登录</button>
      <button
        class="btn"
        :class="{ 'btn-primary': mode === 'register' }"
        @click="switchMode('register')"
      >注册</button>
    </div>

    <div v-if="msg" class="msg" :class="msgType === 'error' ? 'msg-error' : 'msg-ok'">{{ msg }}</div>

    <div class="form-item">
      <label>用户名</label>
      <input v-model="username" class="input" placeholder="3-20 位" />
    </div>
    <div class="form-item">
      <label>密码</label>
      <input v-model="password" type="password" class="input" placeholder="6-72 位" @keyup.enter="submit" />
    </div>

    <button class="btn btn-primary" style="width: 100%" @click="submit">
      {{ mode === 'login' ? '登 录' : '注 册' }}
    </button>
  </div>
</template>
