import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useAppStore } from '@/store'
import AboutPage from './index.vue'

vi.mock('@/api', () => ({
  default: { about: vi.fn() },
}))

vi.mock('@/utils/mathjax', () => ({
  typesetMath: vi.fn().mockResolvedValue(undefined),
  hasMath: () => false,
}))

function mountPage() {
  return mount(AboutPage, {
    global: { stubs: { BannerPage: { template: '<div><slot /></div>' } } },
  })
}

describe('关于我', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.about.mockReset().mockResolvedValue({ code: 0, data: '# 关于我\n\n正文' })
  })

  it('markdown 内容被解析成 HTML', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.html).toContain('<h1'))
  })

  it('接口失败或内容为空都不抛异常', async () => {
    api.about.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.about).toHaveBeenCalled())
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.html).toBe('')

    api.about.mockResolvedValue({ code: 0, data: null })
    const wrapper2 = mountPage()
    await vi.waitFor(() => expect(api.about).toHaveBeenCalledTimes(2))
    await wrapper2.vm.$nextTick()
    expect(wrapper2.vm.html).toBe('')
  })

  // 后台还没填「关于我」时, 页面上只剩一张头像, 看着像坏了
  it('内容为空时给出提示', async () => {
    api.about.mockResolvedValue({ code: 0, data: '' })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.about).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('博主还没有填写关于信息')
  })

  it('有内容时不显示空状态提示', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.html).toContain('<h1'))
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).not.toContain('博主还没有填写关于信息')
  })

  // 回归: 以前是 const { blogConfig } = useAppStore(), blogConfig 是 getter 且
  // store 里 blog_config 整体重新赋值, 解构后拿到的是旧对象, 接口回来后头像不更新
  it('站点配置后到也能更新头像', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.html).toContain('<h1'))

    const appStore = useAppStore()
    appStore.blog_config = { website_avatar: '/images/common/header.jpeg' }
    await wrapper.vm.$nextTick()

    expect(wrapper.find('img').attributes('src')).toBe('/images/common/header.jpeg')
  })
})
