<script setup>
// CategoryView.vue —— 分类管理（新增 / 改名 / 删除）

import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '../api'

const router = useRouter()

const categories = ref([])
const newName = ref('')      // 新增输入框
const editingId = ref(null)  // 正在改名的分类 id
const editName = ref('')     // 改名输入框内容
const msg = ref('')

// 未登录跳登录页
if (!localStorage.getItem('token')) {
  router.replace('/login')
}

async function load() {
  try {
    const res = await getCategoryList()
    categories.value = res.data
  } catch (e) {
    console.error(e)
  }
}

// 新增分类
async function addCategory() {
  if (!newName.value.trim()) return
  try {
    await createCategory({ name: newName.value.trim() })
    newName.value = ''
    load()
  } catch (e) {
    msg.value = e.response?.data?.msg || '创建失败'
  }
}

// 进入编辑状态（把当前名字填进输入框）
function startEdit(c) {
  editingId.value = c.ID
  editName.value = c.name
}

// 保存改名
async function saveEdit(id) {
  if (!editName.value.trim()) return
  try {
    await updateCategory(id, { name: editName.value.trim() })
    editingId.value = null
    load()
  } catch (e) {
    msg.value = e.response?.data?.msg || '修改失败'
  }
}

// 删除分类（后端会拒绝删除还有文章的分类）
async function removeCategory(c) {
  if (!confirm(`确定删除分类「${c.name}」吗？`)) return
  try {
    await deleteCategory(c.ID)
    load()
  } catch (e) {
    alert(e.response?.data?.msg || '删除失败')
  }
}

onMounted(load)
</script>

<template>
  <div style="max-width: 640px; margin: 0 auto">
    <h1 style="font-size: 24px; margin-bottom: 20px">分类管理</h1>

    <div v-if="msg" class="msg msg-error">{{ msg }}</div>

    <!-- 新增分类 -->
    <div style="display: flex; gap: 10px; margin-bottom: 20px">
      <input v-model="newName" class="input" maxlength="50" placeholder="新分类名称" @keyup.enter="addCategory" />
      <button class="btn btn-primary" @click="addCategory">新增</button>
    </div>

    <!-- 分类列表 -->
    <table class="table" v-if="categories.length">
      <thead>
        <tr>
          <th>ID</th>
          <th>名称</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="c in categories" :key="c.ID">
          <td>{{ c.ID }}</td>
          <td>
            <!-- 普通状态显示名字；编辑状态显示输入框 -->
            <template v-if="editingId === c.ID">
              <input v-model="editName" class="input" style="max-width: 200px" @keyup.enter="saveEdit(c.ID)" />
            </template>
            <template v-else>{{ c.name }}</template>
          </td>
          <td>
            <template v-if="editingId === c.ID">
              <button class="btn" style="margin-right: 6px" @click="saveEdit(c.ID)">保存</button>
              <button class="btn" @click="editingId = null">取消</button>
            </template>
            <template v-else>
              <button class="btn" style="margin-right: 6px" @click="startEdit(c)">改名</button>
              <button class="btn btn-danger" @click="removeCategory(c)">删除</button>
            </template>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-else class="empty">还没有分类，先新增一个吧</div>
  </div>
</template>
