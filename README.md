# blog

基于 Gin + GORM + MySQL 的个人博客后端 API。

## 技术栈

- Go 1.26
- Gin（Web 框架）
- GORM（ORM，MySQL）
- JWT（golang-jwt/v5，登录鉴权）
- bcrypt（密码哈希）

## 快速开始

1. 复制配置并填写：`copy .env.example .env`
2. 确保本地 MySQL 已建库：`CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
3. 启动：

```bash
go mod tidy
go run .
```

服务默认监听 `8080` 端口（可用环境变量 `PORT` 修改），首次启动会自动建表。

## API 一览

### 公开接口

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/register | 注册（用户名3-20位，密码6-72位） |
| POST | /api/login | 登录，返回 token（每 IP 限频） |
| GET | /api/article/list | 文章列表，参数：page、size、category_id、keyword |
| GET | /api/article/:id | 文章详情 |
| GET | /api/category/list | 分类列表 |

### 需登录（请求头 `Authorization: Bearer <token>`）

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/article | 创建文章 |
| PUT | /api/article/:id | 修改文章（仅作者） |
| DELETE | /api/article/:id | 删除文章（仅作者） |
| POST | /api/category | 创建分类 |
| PUT | /api/category/:id | 修改分类 |
| DELETE | /api/category/:id | 删除分类（分类下有文章时拒绝） |

## 环境变量

| 变量 | 说明 | 默认值 |
|---|---|---|
| DB_USER / DB_PASS / DB_HOST / DB_PORT / DB_NAME | 数据库连接 | root / 空 / 127.0.0.1 / 3306 / blog_db |
| PORT | 服务端口 | 8080 |
| JWT_SECRET | JWT 签名密钥（至少32位随机串） | 随机生成（重启失效） |
| ALLOWED_ORIGIN | 允许跨域的来源，逗号分隔，`*` 放行所有 | 空（不跨域） |

## 目录结构

```
blog/
├── main.go             入口：DB 初始化 + 启动服务
├── controller/         HTTP 层：参数绑定、状态码
├── service/            业务层：校验、权限、数据组装
├── model/              GORM 模型与数据库初始化
├── middleware/         JWT 鉴权、登录限流、CORS
└── router/             路由注册
```
