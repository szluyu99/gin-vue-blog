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
    global: {
      stubs: {
        BannerPage: { template: '<div><slot /></div>' },
        RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' },
      },
    },
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

  // 这页原来只有一张头像加一段正文, 空得慌; 名片区复用的都是已有数据
  it('名片区展示作者信息、社交链接与统计', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.html).toContain('<h1'))

    const appStore = useAppStore()
    appStore.blog_config = {
      website_author: '阵雨',
      website_intro: '往事随风而去',
      website_avatar: '/images/common/header.jpeg',
      github: 'https://github.com/szluyu99',
      gitee: 'https://gitee.com/szluyu99',
      qq: '123456789',
    }
    appStore.blogInfo = { article_count: 13, category_count: 3, tag_count: 6, view_count: 42 }
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('阵雨')
    expect(wrapper.text()).toContain('往事随风而去')

    // 三项统计各自链到对应页面, 访问量没有对应页面所以不做链接
    const hrefs = wrapper.findAll('a').map(a => a.attributes('href'))
    expect(hrefs).toContain('/archives')
    expect(hrefs).toContain('/categories')
    expect(hrefs).toContain('/tags')
    expect(hrefs).toContain('https://github.com/szluyu99')
    expect(wrapper.text()).toContain('13')
    expect(wrapper.text()).toContain('42')
  })

  // 外链一律要带 noopener, 否则新窗口能通过 window.opener 操作原页面
  it('社交外链带 noopener', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.html).toContain('<h1'))

    const external = wrapper.findAll('a').filter(a => a.attributes('target') === '_blank')
    expect(external.length).toBe(3)
    for (const a of external) {
      expect(a.attributes('rel')).toContain('noopener')
    }
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
