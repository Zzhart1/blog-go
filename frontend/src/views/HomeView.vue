<script setup>
// HomeView.vue —— 首页：文章列表（搜索 + 分类筛选 + 分页）

import { ref, computed, onMounted } from 'vue'
import { getArticleList, getCategoryList, formatTime } from '../api'

const articles = ref([])   // 文章列表
const categories = ref([]) // 分类（用于下拉筛选）
const total = ref(0)       // 总条数
const page = ref(1)        // 当前页
const size = 10            // 每页条数
const keyword = ref('')    // 搜索关键词
const categoryId = ref('') // 选中的分类 id（'' = 全部）

// 计算总页数（分页按钮用）
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size)))

// 加载文章列表
async function loadArticles() {
  try {
    const res = await getArticleList({
      page: page.value,
      size,
      category_id: categoryId.value,
      keyword: keyword.value,
    })
    articles.value = res.data
    total.value = res.total
  } catch (e) {
    console.error(e)
  }
}

// 加载分类（下拉选项）
async function loadCategories() {
  try {
    const res = await getCategoryList()
    categories.value = res.data
  } catch (e) {
    console.error(e)
  }
}

// 点搜索或换分类：回到第一页重新查
function doSearch() {
  page.value = 1
  loadArticles()
}

// 跳转第 n 页
function goPage(n) {
  if (n < 1 || n > totalPages.value || n === page.value) return
  page.value = n
  loadArticles()
}

onMounted(() => {
  loadArticles()
  loadCategories()
})
</script>

<template>
  <div>
    <!-- 搜索 + 筛选区 -->
    <div class="toolbar">
      <select v-model="categoryId" class="select" style="max-width: 180px" @change="doSearch">
        <option value="">全部分类</option>
        <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
      <input
        v-model="keyword"
        class="input"
        style="max-width: 260px"
        placeholder="搜索标题 / 内容..."
        @keyup.enter="doSearch"
      />
      <button class="btn btn-primary" @click="doSearch">搜索</button>
    </div>

    <!-- 文章列表 -->
    <div v-if="articles.length">
      <router-link v-for="a in articles" :key="a.ID" :to="`/article/${a.ID}`" class="card">
        <div class="card-title">{{ a.title }}</div>
        <div class="card-meta">
          <span>作者：{{ a.author ? a.author.username : '未知' }}</span>
          <span>分类：{{ a.category ? a.category.name : '未分类' }}</span>
          <span>{{ formatTime(a.CreatedAt) }}</span>
        </div>
      </router-link>

      <!-- 分页 -->
      <div class="pager">
        <button class="btn" :disabled="page === 1" @click="goPage(page - 1)">上一页</button>
        <button
          v-for="p in totalPages"
          :key="p"
          class="btn"
          :class="{ active: p === page }"
          @click="goPage(p)"
        >
          {{ p }}
        </button>
        <button class="btn" :disabled="page === totalPages" @click="goPage(page + 1)">下一页</button>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty">还没有文章，登录后点击右上角「写文章」开始创作</div>
  </div>
</template>

<style scoped>
/* 本页专用样式：搜索区横排 */
.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}
</style>
