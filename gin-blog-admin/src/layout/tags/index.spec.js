import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { useTagStore } from '@/store'
import AppTags from './index.vue'

// vi.mock 的工厂会被提升到文件顶部, 普通 const 在那时还没初始化
const { push } = vi.hoisted(() => ({ push: vi.fn() }))
vi.mock('@/router', async importOriginal => ({
  ...await importOriginal(),
  router: { push },
}))
vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ path: '/b', fullPath: '/b', name: 'B', meta: { title: 'B' } }),
  useRouter: () => ({ push }),
}))

const tags = [
  { name: 'A', path: '/a', title: 'A' },
  { name: 'B', path: '/b', title: 'B' },
  { name: 'C', path: '/c', title: 'C' },
]

describe('标签栏', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    push.mockClear()
    window.$loadingBar = { start: vi.fn(), finish: vi.fn(), error: vi.fn() }
  })

  it('tabRefs 的顺序与 tags 一致', async () => {
    // A21: v-for 上的模板 ref 数组顺序 Vue 不保证, 而 watch activeIndex 里是按下标取的。
    // 一旦顺序错位, 激活标签的滚动定位就会指到别的标签上
    const tagStore = useTagStore()
    tagStore.setTags(tags)
    await tagStore.setActiveTag('/b')

    // 模板里除了 useRoute() 还直接用了 $route(来自 router 插件的全局属性),
    // 没装插件时它是 undefined, 要单独 mock
    const wrapper = mount(AppTags, {
      global: { stubs: { ContextMenu: true }, mocks: { $route: { name: 'B' } } },
    })
    await wrapper.vm.$nextTick()

    const titles = wrapper.vm.tabRefs.map(c => c.$el.textContent.trim())
    expect(titles).toEqual(['A', 'B', 'C'])
    expect(tagStore.activeIndex).toBe(1)
  })

  it('关闭当前标签会跳到左边那个', async () => {
    const tagStore = useTagStore()
    tagStore.setTags(tags)
    await tagStore.setActiveTag('/b')

    tagStore.removeTag('/b')

    expect(push).toHaveBeenCalledWith('/a')
    expect(tagStore.tags.map(t => t.path)).toEqual(['/a', '/c'])
  })

  it('关闭第一个标签会跳到第二个', async () => {
    const tagStore = useTagStore()
    tagStore.setTags(tags)
    await tagStore.setActiveTag('/a')

    tagStore.removeTag('/a')

    expect(push).toHaveBeenCalledWith('/b')
  })
})
