<script setup>
// EditView.vue —— 写文章 / 编辑文章（一个页面两个用途）
// 路由带 /edit/:id 是编辑，否则是新建

import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticle, createArticle, updateArticle, getCategoryList } from '../api'

const route = useRoute()
const router = useRouter()

const articleId = computed(() => route.params.id) // 有值 = 编辑模式
const categories = ref([])
const title = ref('')
const content = ref('')
const categoryId = ref('')
const msg = ref('')

// 未登录：直接跳登录页
if (!localStorage.getItem('token')) {
  router.replace('/login')
}

// 加载分类下拉 + （编辑模式）加载文章原内容
async function load() {
  try {
    const res = await getCategoryList()
    categories.value = res.data
  } catch (e) {
    console.error(e)
  }

  if (articleId.value) {
    try {
      const res = await getArticle(articleId.value)
      title.value = res.data.title
      content.value = res.data.content
      categoryId.value = String(res.data.category_id)
    } catch (e) {
      msg.value = '文章不存在'
    }
  }
}

// 提交：新建走 createArticle，编辑走 updateArticle
async function submit() {
  if (!title.value.trim() || !content.value.trim()) {
    msg.value = '标题和内容不能为空'
    return
  }
  if (!categoryId.value) {
    msg.value = '请选择分类'
    return
  }
  const data = {
    title: title.value,
    content: content.value,
    category_id: Number(categoryId.value),
  }
  try {
    if (articleId.value) {
      await updateArticle(articleId.value, data)
      router.push(`/article/${articleId.value}`) // 编辑完回详情页
    } else {
      const res = await createArticle(data)
      router.push(`/article/${res.data.ID}`) // 新建完进详情页
    }
  } catch (e) {
    msg.value = e.response?.data?.msg || '保存失败'
  }
}

onMounted(load)
</script>

<template>
  <div style="max-width: 720px; margin: 0 auto">
    <h1 style="font-size: 24px; margin-bottom: 20px">
      {{ articleId ? '编辑文章' : '写文章' }}
    </h1>

    <div v-if="msg" class="msg msg-error">{{ msg }}</div>

    <div class="form-item">
      <label>标题</label>
      <input v-model="title" class="input" maxlength="100" placeholder="文章标题" />
    </div>

    <div class="form-item">
      <label>分类</label>
      <select v-model="categoryId" class="select">
        <option value="">请选择分类</option>
        <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
    </div>

    <div class="form-item">
      <label>正文（纯文本，支持换行）</label>
      <textarea v-model="content" class="textarea" placeholder="开始写作..."></textarea>
    </div>

    <button class="btn btn-primary" @click="submit">保存</button>
  </div>
</template>
