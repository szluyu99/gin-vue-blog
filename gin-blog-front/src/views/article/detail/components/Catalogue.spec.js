import { mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'

import Catalogue from './Catalogue.vue'

// 用真实 DOM 造一段正文, 组件靠 querySelectorAll 解析标题
function makePreview(html) {
  const el = document.createElement('div')
  el.innerHTML = html
  document.body.appendChild(el)
  return el
}

function mountCatalogue(previewRef) {
  return mount(Catalogue, { props: { previewRef } })
}

describe('文章目录', () => {
  afterEach(() => {
    document.body.innerHTML = ''
  })

  it('按标题层级生成缩进', () => {
    const preview = makePreview('<h1>一级</h1><h2>二级</h2><h3>三级</h3><h2>另一个二级</h2>')
    const wrapper = mountCatalogue(preview)

    expect(wrapper.vm.anchors.map(a => a.name)).toEqual(['一级', '二级', '三级', '另一个二级'])
    expect(wrapper.vm.anchors.map(a => a.indent)).toEqual([0, 1, 2, 1])
  })

  it('只有二三级标题时缩进从 0 开始', () => {
    const preview = makePreview('<h2>二级</h2><h3>三级</h3>')
    const wrapper = mountCatalogue(preview)

    expect(wrapper.vm.anchors.map(a => a.indent)).toEqual([0, 1])
  })

  it('重名标题也能生成不同的 id', () => {
    const preview = makePreview('<h2>安装</h2><h2>安装</h2>')
    const wrapper = mountCatalogue(preview)

    const ids = wrapper.vm.anchors.map(a => a.id)
    expect(new Set(ids).size).toBe(2)
    // id 写回到了真实的 DOM 节点上, 点击时靠它定位
    expect(document.getElementById(ids[0])).not.toBeNull()
    expect(document.getElementById(ids[1])).not.toBeNull()
  })

  // 以前 hTags 用过滤后的列表算层级, 但循环用的是未过滤的标题,
  // markdown 里写了 ## 却没写文字时目录里会多出一行空白
  it('空标题不进目录', () => {
    const preview = makePreview('<h2>正常标题</h2><h2>   </h2><h3></h3>')
    const wrapper = mountCatalogue(preview)

    expect(wrapper.vm.anchors).toHaveLength(1)
    expect(wrapper.vm.anchors[0].name).toBe('正常标题')
  })

  it('正文里没有标题时目录为空', () => {
    const preview = makePreview('<p>只有正文</p>')
    const wrapper = mountCatalogue(preview)

    expect(wrapper.vm.anchors).toEqual([])
    expect(wrapper.text()).toContain('目录')
  })

  // 正文还没渲染出来时父组件可能传进 null
  it('previewRef 为 null 时不抛异常', () => {
    expect(() => mountCatalogue(null)).not.toThrow()
  })

  it('点击目录项滚动到对应标题', () => {
    const preview = makePreview('<h2>安装</h2>')
    const wrapper = mountCatalogue(preview)
    const scrollTo = vi.fn()
    window.scrollTo = scrollTo

    wrapper.vm.handleClickAnchor(wrapper.vm.anchors[0].id)
    expect(scrollTo).toHaveBeenCalledWith(expect.objectContaining({ behavior: 'smooth' }))

    // 找不到元素时直接返回, 不抛异常
    scrollTo.mockClear()
    expect(() => wrapper.vm.handleClickAnchor('不存在的-id')).not.toThrow()
    expect(scrollTo).not.toHaveBeenCalled()
  })
})
