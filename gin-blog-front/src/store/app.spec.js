import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useAppStore } from '@/store/app'

vi.mock('@/api', () => ({
  default: {
    getHomeData: vi.fn(),
    getPageList: vi.fn(),
  },
}))

const api = (await import('@/api')).default

describe('useAppStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
    document.documentElement.classList.remove('dark')
  })

  it('弹框开关只改对应字段', () => {
    const store = useAppStore()

    store.setLoginFlag(true)
    store.setSearchFlag(true)
    expect(store.loginFlag).toBe(true)
    expect(store.searchFlag).toBe(true)
    expect(store.registerFlag).toBe(false)

    store.setLoginFlag(false)
    expect(store.loginFlag).toBe(false)
  })

  it('getBlogInfo 写入统计数据, 并补全站点头像地址', async () => {
    api.getHomeData.mockResolvedValue({
      code: 0,
      data: {
        article_count: 3,
        category_count: 2,
        tag_count: 5,
        view_count: 100,
        blog_config: { website_name: '博客', website_avatar: 'public/uploaded/avatar.png' },
      },
    })

    const store = useAppStore()
    await store.getBlogInfo()

    expect(store.articleCount).toBe(3)
    expect(store.tagCount).toBe(5)
    expect(store.viewCount).toBe(100)
    expect(store.blogConfig.website_avatar).toBe('/public/uploaded/avatar.png')
  })

  it('getBlogInfo 遇到业务错误码时 reject, 不写脏数据', async () => {
    api.getHomeData.mockResolvedValue({ code: 9999, message: '服务异常' })

    const store = useAppStore()
    await expect(store.getBlogInfo()).rejects.toMatchObject({ code: 9999 })
    expect(store.articleCount).toBe(0)
  })

  it('getPageList 给每个页面封面拼上后端地址', async () => {
    api.getPageList.mockResolvedValue({
      code: 0,
      data: [
        { name: '首页', label: 'home', cover: 'public/uploaded/home.png' },
        { name: '归档', label: 'archive', cover: 'https://cdn.com/archive.png' },
      ],
    })

    const store = useAppStore()
    await store.getPageList()

    expect(store.pageList).toHaveLength(2)
    expect(store.pageList[0].cover).toBe('/public/uploaded/home.png')
    // 已经是完整 URL 的不动
    expect(store.pageList[1].cover).toBe('https://cdn.com/archive.png')
  })

  it('toggleTheme 同时改 store、localStorage 和 html 的 class', () => {
    const store = useAppStore()

    store.toggleTheme()
    expect(store.theme).toBe('dark')
    expect(store.isDark).toBe(true)
    expect(localStorage.getItem('blog-theme')).toBe('dark')
    expect(document.documentElement.classList.contains('dark')).toBe(true)

    store.toggleTheme()
    expect(store.theme).toBe('light')
    expect(localStorage.getItem('blog-theme')).toBe('light')
    expect(document.documentElement.classList.contains('dark')).toBe(false)
  })

  it('initTheme 优先用存过的偏好', () => {
    localStorage.setItem('blog-theme', 'dark')

    const store = useAppStore()
    store.initTheme()

    expect(store.theme).toBe('dark')
    expect(document.documentElement.classList.contains('dark')).toBe(true)
  })

  it('initTheme 没有偏好时跟随系统', () => {
    // jsdom 没有实现 matchMedia, 直接塞一个假的
    const matchMedia = vi.fn().mockReturnValue({ matches: true })
    vi.stubGlobal('matchMedia', matchMedia)

    const store = useAppStore()
    store.initTheme()

    expect(matchMedia).toHaveBeenCalledWith('(prefers-color-scheme: dark)')
    expect(store.theme).toBe('dark')
    vi.unstubAllGlobals()
  })

  it('setTheme 遇到非法值回落到浅色', () => {
    const store = useAppStore()

    store.setTheme('dark')
    store.setTheme('rainbow')

    expect(store.theme).toBe('light')
    expect(document.documentElement.classList.contains('dark')).toBe(false)
  })
})
