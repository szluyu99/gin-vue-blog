import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import LinkList from './LinkList.vue'

const links = [
  { id: 1, name: 'Gin', avatar: 'public/uploaded/a.png', address: 'https://gin-gonic.com', intro: 'Go 的 Web 框架' },
  { id: 2, name: 'Vue', avatar: '', address: 'https://cn.vuejs.org', intro: '渐进式框架' },
]

function mountList(linkList = links) {
  return mount(LinkList, { props: { linkList } })
}

describe('友链列表', () => {
  it('渲染每个友链的名称、简介与外链', () => {
    const wrapper = mountList()

    expect(wrapper.findAll('.link-wrapper')).toHaveLength(2)
    expect(wrapper.text()).toContain('Gin')
    expect(wrapper.text()).toContain('Go 的 Web 框架')

    const a = wrapper.find('a')
    expect(a.attributes('href')).toBe('https://gin-gonic.com')
    expect(a.attributes('target')).toBe('_blank')
  })

  // 头像原来只写了 w-[65px], 图床挂掉或网络不通时高度是 0(实测 65x0),
  // 整行文字会跟着塌上去
  it('头像同时指定宽高, 不会因为加载失败塌成 0 高', () => {
    const img = mountList().find('img')

    expect(img.classes()).toContain('w-[65px]')
    expect(img.classes()).toContain('h-[65px]')
    expect(img.classes()).toContain('object-cover')
    // 有 alt 才能在加载失败时看出这是谁
    expect(img.attributes('alt')).toBe('Gin')
  })

  it('没有友链时展示空状态', () => {
    const wrapper = mountList([])

    expect(wrapper.findAll('.link-wrapper')).toHaveLength(0)
    expect(wrapper.text()).toContain('暂无友情链接')
    expect(wrapper.find('img[alt="暂无友情链接"]').exists()).toBe(true)
  })

  it('默认 props 为空数组, 不传也不会崩', () => {
    const wrapper = mount(LinkList)
    expect(wrapper.text()).toContain('暂无友情链接')
  })
})
