import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useTagStore } from '@/store'

// tag store 里的关闭标签逻辑会跳转路由, 这里只断言跳到了哪
vi.mock('@/router', () => ({
  router: { push: vi.fn(), replace: vi.fn(), currentRoute: { query: {} } },
  resetRouter: vi.fn(),
}))
vi.mock('@/api', () => ({ default: {} }))

function tag(path) {
  return { name: path.slice(1), path, title: path.slice(1) }
}

describe('tag store', () => {
  let store

  beforeEach(async () => {
    setActivePinia(createPinia())
    store = useTagStore()
    const { router } = await import('@/router')
    router.push.mockClear()
    window.$loadingBar = { start: vi.fn(), finish: vi.fn(), error: vi.fn() }
  })

  it('addTag 追加新标签并激活', async () => {
    store.addTag(tag('/home'))
    store.addTag(tag('/article'))
    await vi.waitFor(() => expect(store.activeTag).toBe('/article'))

    expect(store.tags.map(e => e.path)).toEqual(['/home', '/article'])
    expect(store.activeIndex).toBe(1)
  })

  it('addTag 重复 path 时替换而不是追加', () => {
    store.addTag(tag('/home'))
    store.addTag({ ...tag('/home'), title: '新标题' })

    expect(store.tags).toHaveLength(1)
    expect(store.tags[0].title).toBe('新标题')
  })

  it('removeTag 关闭非当前标签时不跳转', async () => {
    const { router } = await import('@/router')
    store.setTags([tag('/home'), tag('/article')])
    store.activeTag = '/article'

    store.removeTag('/home')

    expect(store.tags.map(e => e.path)).toEqual(['/article'])
    expect(router.push).not.toHaveBeenCalled()
  })

  it('removeTag 关闭当前标签时跳到左边的标签', async () => {
    const { router } = await import('@/router')
    store.setTags([tag('/home'), tag('/article'), tag('/user')])
    store.activeTag = '/article'

    store.removeTag('/article')

    expect(router.push).toHaveBeenCalledWith('/home')
    expect(store.tags.map(e => e.path)).toEqual(['/home', '/user'])
  })

  it('removeTag 关闭第一个标签时跳到第二个标签', async () => {
    const { router } = await import('@/router')
    store.setTags([tag('/home'), tag('/article')])
    store.activeTag = '/home'

    store.removeTag('/home')

    expect(router.push).toHaveBeenCalledWith('/article')
  })

  it('removeOther 只保留指定标签, 关掉当前标签时会跳转', async () => {
    const { router } = await import('@/router')
    store.setTags([tag('/home'), tag('/article'), tag('/user')])
    store.activeTag = '/user'

    store.removeOther('/home')

    expect(store.tags.map(e => e.path)).toEqual(['/home'])
    expect(router.push).toHaveBeenCalledWith('/home')
  })

  it('removeOther 默认关闭当前标签之外的所有标签', async () => {
    const { router } = await import('@/router')
    store.setTags([tag('/home'), tag('/article')])
    store.activeTag = '/article'

    store.removeOther()

    expect(store.tags.map(e => e.path)).toEqual(['/article'])
    expect(router.push).not.toHaveBeenCalled()
  })

  it('removeLeft 关掉左侧标签, 当前标签被关掉则跳到最后一个', async () => {
    const { router } = await import('@/router')
    store.setTags([tag('/home'), tag('/article'), tag('/user')])
    store.activeTag = '/home'

    store.removeLeft('/article')

    expect(store.tags.map(e => e.path)).toEqual(['/article', '/user'])
    expect(router.push).toHaveBeenCalledWith('/user')
  })

  it('removeRight 关掉右侧标签, 当前标签保留时不跳转', async () => {
    const { router } = await import('@/router')
    store.setTags([tag('/home'), tag('/article'), tag('/user')])
    store.activeTag = '/home'

    store.removeRight('/article')

    expect(store.tags.map(e => e.path)).toEqual(['/home', '/article'])
    expect(router.push).not.toHaveBeenCalled()
  })

  it('resetTags 清空标签', () => {
    store.setTags([tag('/home')])
    store.activeTag = '/home'

    store.resetTags()

    expect(store.tags).toEqual([])
    expect(store.activeTag).toBe('')
  })

  it('updateAliveKey 按 route name 刷新 keepAlive 的 key', () => {
    store.updateAliveKey('Article')
    const first = store.aliveKeys.Article
    expect(first).toBeTypeOf('number')

    store.updateAliveKey('Article')
    expect(store.aliveKeys.Article).toBeGreaterThanOrEqual(first)
  })

  it('reloadTag 通过 reloading 触发一次白屏, 结束后恢复', async () => {
    // layout/index.vue 里是 v-if="tagStore.reloading"
    expect(store.reloading).toBe(true)

    const done = store.reloadTag()
    expect(store.reloading).toBe(false)

    await done

    expect(store.reloading).toBe(true)
    expect(window.$loadingBar.start).toHaveBeenCalled()
  })
})
