import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { h } from 'vue'

import { useTagStore } from '@/store'
import AppTags from './index.vue'

// vi.mock 的工厂会被提升到文件顶部, 普通 const 在那时还没初始化
const { push, handleScroll } = vi.hoisted(() => ({ push: vi.fn(), handleScroll: vi.fn() }))
vi.mock('@/router', async importOriginal => ({
  ...await importOriginal(),
  router: { push },
}))
vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ path: '/b', fullPath: '/b', name: 'B', meta: { title: 'B' } }),
  useRouter: () => ({ push }),
}))

// 替换 ScrollX, 只关心它被要求滚动到什么位置
const ScrollXStub = {
  name: 'ScrollX',
  setup(_props, { slots, expose }) {
    expose({ handleScroll })
    return () => h('div', slots.default?.())
  },
}

const tags = [
  { name: 'A', path: '/a', title: 'A' },
  { name: 'B', path: '/b', title: 'B' },
  { name: 'C', path: '/c', title: 'C' },
]

// jsdom 里所有元素的 offsetLeft / offsetWidth 都是 0, 手动给每个标签造不同的几何值,
// 这样才能从 handleScroll 的入参反推出用的是哪一个标签
function stubGeometry(wrapper) {
  wrapper.findAll('.n-tag').forEach((tag, i) => {
    Object.defineProperty(tag.element, 'offsetLeft', { value: (i + 1) * 100 })
    Object.defineProperty(tag.element, 'offsetWidth', { value: (i + 1) * 10 })
  })
}

async function mountTags() {
  const wrapper = mount(AppTags, {
    global: { stubs: { ScrollX: ScrollXStub, ContextMenu: true } },
  })
  await wrapper.vm.$nextTick()
  stubGeometry(wrapper)
  return wrapper
}

describe('标签栏', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    push.mockClear()
    handleScroll.mockClear()
    window.$loadingBar = { start: vi.fn(), finish: vi.fn(), error: vi.fn() }
  })

  it('滚动定位取的是激活标签自己的元素', async () => {
    // A21: 原来是 tabRefs.value[activeIndex], 而 v-for 上的模板 ref 数组
    // Vue 不保证顺序与 tags 一致, 顺序一旦错位就会滚到别的标签上。
    // 现在按 path 存元素, 与渲染顺序无关
    const tagStore = useTagStore()
    tagStore.setTags(tags)
    const wrapper = await mountTags()

    handleScroll.mockClear()
    await tagStore.setActiveTag('/c')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    // 第三个标签: offsetLeft 300, offsetWidth 30
    expect(handleScroll).toHaveBeenLastCalledWith(330, 30)

    await tagStore.setActiveTag('/a')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(handleScroll).toHaveBeenLastCalledWith(110, 10)
  })

  it('标签被移除后不再持有它的元素', async () => {
    const tagStore = useTagStore()
    tagStore.setTags(tags)
    const wrapper = await mountTags()

    tagStore.removeTag('/c')
    await wrapper.vm.$nextTick()
    handleScroll.mockClear()

    // 激活一个已经不存在的标签: 找不到元素就不滚动, 也不能抛异常
    await tagStore.setActiveTag('/c')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(handleScroll).not.toHaveBeenCalled()
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
