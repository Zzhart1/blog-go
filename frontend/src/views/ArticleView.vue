<script setup>
// ArticleView.vue —— 文章详情 + 评论区

import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  getArticle,
  getComments,
  createComment,
  deleteComment,
  formatTime,
} from '../api'

const route = useRoute()
const router = useRouter()

const article = ref(null)    // 文章详情
const comments = ref([])     // 评论列表
const myComment = ref('')    // 输入框内容
const errorMsg = ref('')     // 错误提示

const articleId = computed(() => route.params.id)
const myName = ref(localStorage.getItem('username') || '')
const isLogin = computed(() => !!myName.value)

// 加载文章 + 评论
async function load() {
  try {
    const res = await getArticle(articleId.value)
    article.value = res.data
  } catch (e) {
    errorMsg.value = '文章不存在或已被删除'
    return
  }
  loadComments()
}

async function loadComments() {
  try {
    const res = await getComments(articleId.value)
    comments.value = res.data
  } catch (e) {
    console.error(e)
  }
}

// 发表评论
async function submitComment() {
  if (!myComment.value.trim()) return
  try {
    await createComment({ article_id: Number(articleId.value), content: myComment.value })
    myComment.value = ''
    loadComments() // 重新拉取评论列表
  } catch (e) {
    errorMsg.value = e.response?.data?.msg || '评论失败'
  }
}

// 删除自己的评论
async function removeComment(id) {
  if (!confirm('确定删除这条评论吗？')) return
  try {
    await deleteComment(id)
    loadComments()
  } catch (e) {
    alert(e.response?.data?.msg || '删除失败')
  }
}

onMounted(load)
</script>

<template>
  <div v-if="article">
    <h1 class="article-title">{{ article.title }}</h1>
    <div class="card-meta">
      <span>作者：{{ article.author ? article.author.username : '未知' }}</span>
      <span>分类：{{ article.category ? article.category.name : '未分类' }}</span>
      <span>{{ formatTime(article.CreatedAt) }}</span>
    </div>

    <!-- 正文：纯文本显示，保留换行 -->
    <div class="article-content">{{ article.content }}</div>

    <!-- 评论区 -->
    <div class="comment-box">
      <h2 style="font-size: 20px; margin-bottom: 12px">评论（{{ comments.length }}）</h2>

      <!-- 发表评论：登录了才能写 -->
      <div v-if="isLogin" style="margin-bottom: 16px">
        <textarea v-model="myComment" class="textarea" style="min-height: 80px" placeholder="写下你的评论..."></textarea>
        <button class="btn btn-primary" style="margin-top: 8px" @click="submitComment">发表评论</button>
      </div>
      <div v-else style="margin-bottom: 16px; color: #8a919f; font-size: 14px">
        <router-link to="/login">登录</router-link> 后可评论
      </div>

      <div v-if="errorMsg" class="msg msg-error">{{ errorMsg }}</div>

      <!-- 评论列表 -->
      <div v-if="comments.length">
        <div v-for="c in comments" :key="c.ID" class="comment-item">
          <div class="comment-head">
            <span class="who">{{ c.user ? c.user.username : '匿名' }}</span>
            <span>{{ formatTime(c.CreatedAt) }}</span>
            <!-- 只能删自己的评论 -->
            <button v-if="c.user && c.user.username === myName" class="btn-link" @click="removeComment(c.ID)">删除</button>
          </div>
          <div style="margin-top: 4px">{{ c.content }}</div>
        </div>
      </div>
      <div v-else class="empty" style="padding: 20px 0">还没有评论，来抢沙发~</div>
    </div>
  </div>

  <div v-else-if="errorMsg" class="msg msg-error">{{ errorMsg }}</div>
</template>
