import { defineStore } from 'pinia'
import { shallowRef } from 'vue'
import { asyncRoutes, basicRoutes, vueModules } from '@/router/routes'
import Layout from '@/layout/index.vue'
import api from '@/api'

export const usePermissionStore = defineStore('permission', {
  state() {
    return {
      accessRoutes: [],
    }
  },
  getters: {
    routes() {
      // 在 基础路由 上组合 需要权限的路由
      return basicRoutes.concat(this.accessRoutes)
    },
    menus() {
      // 过滤掉 isHidden 的路由
      return this.routes.filter(route => route.name && !route.isHidden)
    },
  },
  actions: {
    // ! 后端生成路由
    async generateRoutesFromBack() {
      const res = await api.getUserMenus()
      this.accessRoutes = res.data.map(e => ({
        name: e.name,
        path: e.component !== 'Layout' ? '/' : e.path, // 处理目录是一级菜单的情况
        component: shallowRef(Layout), // 不使用 shallowRef 控制台会有 warning
        isHidden: e.is_hidden,
        redirect: e.redirect,
        meta: {
          title: e.name,
          icon: e.icon,
          order: e.order_num,
          keepAlive: e.keep_alive,
        },
        children: e.children.map(ee => ({
          name: ee.name,
          path: ee.path, // 父路径 + 当前菜单路径
          component: vueModules[`/src/views${ee.component}/index.vue`], // ! 读取动态加载的路由模块
          isHidden: ee.is_hidden,
          meta: {
            title: ee.name,
            icon: ee.icon,
            order: ee.order_num,
            keepAlive: ee.keep_alive,
          },
        })),
      }))
      return this.accessRoutes
    },
    // ! 前端控制路由权限: 根据角色过滤路由
    generateRoutes(role = []) {
      return filterAsyncRoutes(asyncRoutes, role)
    },
    resetPermission() {
      this.$reset()
    },
  },
})

// 前端路由相关函数
// 判断用户角色是否有权限访问路由
function hasPermission(route, role) {
  // 路由不需要权限直接返回 true
  if (!route.meta?.requireAuth)
    return true
  // 路由需要的角色
  const routeRole = route.meta?.role ?? []
  // 登录用户没有角色 或者 路由没有设置角色判断, 为没有权限
  if (!role.length || !routeRole.length)
    return false
  // 路由指定的角色包含任登录用户角色则判定有权限
  return role.some(item => routeRole.includes(item))
}

// 从路由中过滤出有权限访问的路由
function filterAsyncRoutes(routes = [], role) {
  const res = []
  routes.forEach((route) => {
    // FIXME: 需要先判断 route 不为 null, 什么原因会加入 null? 初步估计是 route.js 中的配置...
    if (route && hasPermission(route, role)) {
      const curRoute = { ...route, children: [] }
      if (route.children?.length)
        curRoute.children = filterAsyncRoutes(route.children, role)
      else
        Reflect.deleteProperty(curRoute, 'children')
      res.push(curRoute)
    }
  })
  return res
}
