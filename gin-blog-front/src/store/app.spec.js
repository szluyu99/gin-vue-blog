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
    expect(store.blogConfig.website_avatar).toBe('http://test-server/public/uploaded/avatar.png')
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
    expect(store.pageList[0].cover).toBe('http://test-server/public/uploaded/home.png')
    // 已经是完整 URL 的不动
    expect(store.pageList[1].cover).toBe('https://cdn.com/archive.png')
  })
})
