<script setup>
// App.vue —— 整个站点的外壳：顶部导航栏 + 内容区 + 页脚
// 登录状态存 localStorage，这里统一管理"登录/退出"

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const isLogin = computed(() => !!username.value)

function logout() {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  username.value = ''
  router.push('/')
}
</script>

<template>
  <div class="layout">
    <header class="navbar">
      <div class="nav-inner">
        <router-link to="/" class="logo">张子皓的博客</router-link>
        <nav class="nav-links">
          <router-link to="/">文章</router-link>
          <router-link to="/edit">写文章</router-link>
          <router-link to="/categories">分类管理</router-link>
        </nav>
        <div class="nav-user">
          <template v-if="isLogin">
            <span class="username">{{ username }}</span>
            <button class="btn-link" @click="logout">退出</button>
          </template>
          <template v-else>
            <router-link to="/login" class="btn btn-primary">登录</router-link>
          </template>
        </div>
      </div>
    </header>

    <main class="container">
      <!-- 每个路由对应的页面会渲染在这里 -->
      <router-view />
    </main>

    <footer class="footer">© 2026 张子皓的博客 · Go + Vue 练手项目</footer>
  </div>
</template>
