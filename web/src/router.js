// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('./views/Home.vue'),
    redirect: '/chat',
    children: [
      {
        path: '/chat',
        name: 'Chat',
        meta: { title: 'AI-对话' },
        component: () => import('./views/ChatPage.vue'),
      },
      {
        path: '/chat/:id',
        name: 'ChatDetail',
        meta: { title: 'AI-对话' },
        component: () => import('./views/ChatPage.vue'),
      },
      {
        path: '/workflow',
        name: 'Workflow',
        meta: { title: '工作流空间' },
        component: () => import('./views/WorkflowPage.vue'),
      },
    ]
  },
  {
    path: '/share/:chat_id',
    meta: { title: '对话分享' },
    component: () => import('./views/ChatShare.vue'),
  },
  {
    path: '/admin/login',
    meta: { title: '用户登录' },
    component: () => import('./views/admin/Login.vue'),
  },
  {
    path: '/admin',
    redirect: '/admin/dashboard',
    component: () => import('./views/admin/Home.vue'),
    meta: { title: '控制台' },
    children: [
      {
        path: '/admin/dashboard',
        meta: { title: '仪表盘' },
        component: () => import('./views/admin/DashBoard.vue'),
      },
      {
        path: '/admin/users',
        meta: { title: '用户列表' },
        component: () => import('./views/admin/UserList.vue'),
      },
      {
        path: '/admin/manager',
        meta: { title: '管理员列表' },
        component: () => import('./views/admin/ManagerList.vue'),
      },
      {
        path: '/admin/apps/list',
        meta: { title: '应用列表' },
        component: () => import('./views/admin/AppList.vue'),
      },
      {
        path: '/admin/apps/category',
        meta: { title: '应用分类' },
        component: () => import('./views/admin/AppCategory.vue'),
      },
      {
        path: '/admin/redeem',
        meta: { title: '兑换码' },
        component: () => import('./views/admin/Redeem.vue'),
      },
      {
        path: '/admin/score/log',
        meta: { title: '积分日志' },
        component: () => import('./views/admin/ScoreLog.vue'),
      },
      {
        path: '/admin/users/loginLog',
        meta: { title: '登录日志' },
        component: () => import('./views/admin/LoginLog.vue'),
      },
      {
        path: '/admin/settings/basic',
        meta: { title: '基础设置' },
        component: () => import('./views/admin/settings/Basic.vue'),
      },
      {
        path: '/admin/settings/notice',
        meta: { title: '公告管理' },
        component: () => import('./views/admin/settings/Notice.vue'),
      },
      {
        path: '/admin/settings/geek',
        meta: { title: '增值服务配置' },
        component: () => import('./views/admin/settings/GeekService.vue'),
      },
      {
        path: '/admin/settings/coze',
        meta: { title: 'Coze API 设置' },
        component: () => import('./views/admin/settings/Coze.vue'),
      },
      {
        path: '/admin/settings/payment',
        meta: { title: '支付设置' },
        component: () => import('./views/admin/settings/Payment.vue'),
      },
      {
        path: '/admin/settings/sms',
        meta: { title: '短信服务设置' },
        component: () => import('./views/admin/settings/Sms.vue'),
      },
      {
        path: '/admin/settings/oss',
        meta: { title: '文件存储设置' },
        component: () => import('./views/admin/settings/FileStore.vue'),
      },
      {
        path: '/admin/settings/smtp',
        meta: { title: '邮件服务设置' },
        component: () => import('./views/admin/settings/Smtp.vue'),
      },
      {
        path: '/admin/creators',
        meta: { title: '创作者管理' },
        component: () => import('./views/admin/creator/CreatorList.vue'),
      },
      {
        path: '/admin/creator/apps',
        meta: { title: '创作者应用管理' },
        component: () => import('./views/admin/creator/CreatorAppList.vue'),
      },
      {
        path: '/admin/creator/withdraws',
        meta: { title: '创作者提现管理' },
        component: () => import('./views/admin/creator/WithdrawList.vue'),
      },
      {
        path: '/admin/products',
        meta: { title: '充值产品列表' },
        component: () => import('./views/admin/Product.vue'),
      },
      {
        path: '/admin/orders',
        name: 'OrderList',
        component: () => import('./views/admin/OrderList.vue'),
        meta: { title: '订单管理' },
      },
      {
        path: '/admin/workflows',
        meta: { title: '工作流管理' },
        component: () => import('./views/admin/Workflows.vue'),
      },
    ],
  },
  {
    path: '/creator/console',
    meta: { title: '创作者控制台' },
    component: () => import('./views/creator/Console.vue'),
  },
  {
    path: '/creator/:username',
    meta: { title: '创作者首页' },
    component: () => import('./views/creator/CreatorHome.vue'),
  },
  {
    path: '/test',
    meta: { title: '测试' },
    component: () => import('./views/Test.vue'),
  },
  {
    path: '/500',
    meta: { title: '系统开小差了' },
    component: () => import('./views/errors/500.vue'),
  },
  {
    path: '/503',
    meta: { title: '服务暂时不可用' },
    component: () => import('./views/errors/503.vue'),
  },
  {
    name: 'NotFound',
    path: '/:all(.*)',
    meta: { title: '页面没有找到' },
    component: () => import('./views/errors/404.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes: routes,
})

let prevRoute = null
// dynamic change the title when router change
router.beforeEach((to, from, next) => {
  document.title = to.meta.title
  prevRoute = from
  next()
})

export { prevRoute, router }
