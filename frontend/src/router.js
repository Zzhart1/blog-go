// router.js —— 页面路由表
// 每个路由：path 网址、name 名字、component 对应的页面组件

import { createRouter, createWebHashHistory } from 'vue-router'

// 页面组件（懒加载：访问到了才加载，首页更快）
const HomeView = () => import('./views/HomeView.vue')
const ArticleView = () => import('./views/ArticleView.vue')
const LoginView = () => import('./views/LoginView.vue')
const EditView = () => import('./views/EditView.vue')
const CategoryView = () => import('./views/CategoryView.vue')

const routes = [
  { path: '/', name: 'home', component: HomeView },          // 首页：文章列表
  { path: '/article/:id', name: 'article', component: ArticleView }, // 文章详情 + 评论
  { path: '/login', name: 'login', component: LoginView },   // 登录 / 注册
  { path: '/edit', name: 'edit', component: EditView },      // 写文章
  { path: '/edit/:id', name: 'editArticle', component: EditView }, // 编辑文章
  { path: '/categories', name: 'categories', component: CategoryView }, // 分类管理
]

// 使用 hash 路由（网址带 #），好处：静态文件部署不需要服务器额外配置
const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
