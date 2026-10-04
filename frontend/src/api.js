// api.js —— 所有后端接口的封装
// 统一从这里发请求，页面里只调函数，不直接写 axios

import axios from 'axios'

// 创建 axios 实例
// baseURL 为空 = 走同源 /api（开发时由 vite 代理到后端 8080，上线时由 Nginx 反代）
const http = axios.create({
  baseURL: '',
  timeout: 10000,
})

// 请求拦截器：每次请求自动带上登录 token
http.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：统一处理错误（401 = 未登录/登录过期）
http.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      // 跳到登录页（避免死循环，登录页本身不触发）
      if (!location.hash.includes('/login')) {
        location.hash = '#/login'
      }
    }
    return Promise.reject(err)
  }
)

// ---------- 用户 ----------
export const register = (data) => http.post('/api/register', data)
export const login = (data) => http.post('/api/login', data)

// ---------- 文章 ----------
export const getArticleList = (params) =>
  http.get('/api/article/list', { params }) // params: {page,size,category_id,keyword}
export const getArticle = (id) => http.get(`/api/article/${id}`)
export const createArticle = (data) => http.post('/api/article', data)
export const updateArticle = (id, data) => http.put(`/api/article/${id}`, data)
export const deleteArticle = (id) => http.delete(`/api/article/${id}`)

// ---------- 分类 ----------
export const getCategoryList = () => http.get('/api/category/list')
export const createCategory = (data) => http.post('/api/category', data)
export const updateCategory = (id, data) => http.put(`/api/category/${id}`, data)
export const deleteCategory = (id) => http.delete(`/api/category/${id}`)

// ---------- 评论 ----------
export const getComments = (articleId) => http.get(`/api/article/${articleId}/comments`)
export const createComment = (data) => http.post('/api/comment', data)
export const deleteComment = (id) => http.delete(`/api/comment/${id}`)

// ---------- 小工具 ----------
// 把后端返回的时间（如 2026-10-02T15:30:00+08:00）格式化成 2026-10-02 15:30
export function formatTime(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}
