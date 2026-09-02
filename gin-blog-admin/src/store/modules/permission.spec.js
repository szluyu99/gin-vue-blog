import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { basicRoutes } from '@/router/routes'
import { usePermissionStore } from '@/store'

const getUserMenus = vi.fn()

vi.mock('@/api', () => ({ default: { getUserMenus: (...args) => getUserMenus(...args) } }))
vi.mock('@/router', () => ({
  router: { push: vi.fn(), replace: vi.fn(), currentRoute: { query: {} } },
  resetRouter: vi.fn(),
}))
// Layout 组件只是被塞进路由的 component 字段, 测试里不需要真实渲染
vi.mock('@/layout/index.vue', () => ({ default: { name: 'Layout' } }))

describe('permission store', () => {
  let store

  beforeEach(() => {
    setActivePinia(createPinia())
    store = usePermissionStore()
    getUserMenus.mockReset()
  })

  it('routes 是基础路由 + 可访问路由', () => {
    expect(store.routes).toEqual(basicRoutes)

    store.accessRoutes = [{ name: 'Home', path: '/home' }]

    expect(store.routes).toHaveLength(basicRoutes.length + 1)
  })

  it('menus 过滤掉隐藏路由和没有 name 的路由', () => {
    store.accessRoutes = [
      { name: 'Home', path: '/home' },
      { name: 'Hidden', path: '/hidden', isHidden: true },
      { path: '/no-name' },
    ]

    // basicRoutes 里的 Login / 404 都是 isHidden
    expect(store.menus.map(e => e.name)).toEqual(['Home'])
  })

  it('generateRoutesFront 按角色过滤路由', () => {
    const routes = store.generateRoutesFront(['admin'])

    expect(routes.length).toBeGreaterThan(0)
    expect(store.accessRoutes).toBe(routes)
  })

  it('generateRoutesFront 没有角色时只剩不需要鉴权的路由', () => {
    const routes = store.generateRoutesFront([])

    // views 下的路由都带 requireAuth, 没有角色时一条都拿不到
    expect(routes.every(route => !route.meta?.requireAuth)).toBe(true)
  })

  it('generateRoutesBack 把后端菜单转成前端路由', async () => {
    getUserMenus.mockResolvedValue({
      code: 0,
      data: [
        {
          name: '文章管理',
          path: '/article',
          component: 'Layout',
          icon: 'i-mdi:file',
          order_num: 2,
          is_catalogue: false,
          is_hidden: false,
          children: [
            { name: '文章列表', path: 'list', component: '/article/list', icon: 'i-mdi:list', order_num: 1, keep_alive: true },
          ],
        },
        {
          name: '首页',
          path: '/home',
          component: 'Layout',
          order_num: 1,
          is_catalogue: true,
          is_hidden: false,
          keep_alive: false,
        },
      ],
    })

    const routes = await store.generateRoutesBack()

    expect(routes).toHaveLength(2)
    // 非目录: 自己作为带 Layout 的父路由, 子菜单挂在 children 上
    expect(routes[0].name).toBe('文章管理Parent')
    expect(routes[0].meta).toMatchObject({ title: '文章管理', order: 2, keepAlive: undefined })
    expect(routes[0].children.map(e => e.name)).toEqual(['文章列表'])
    // 目录: 父路由名加 -catalogue 后缀, 避免和子路由同名
    expect(routes[1].name).toBe('首页-catalogue')
    expect(routes[1].isCatalogue).toBe(true)
    expect(routes[1].children[0].name).toBe('首页')
  })

  it('resetPermission 清空可访问路由', () => {
    store.accessRoutes = [{ name: 'Home', path: '/home' }]

    store.resetPermission()

    expect(store.accessRoutes).toEqual([])
  })
})
