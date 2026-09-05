import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import ArticleCard from './ArticleCard.vue'

const article = {
  id: 1,
  title: '第一篇',
  content: '摘要',
  img: 'public/uploaded/a.png',
  created_at: '2026-09-01T00:00:00Z',
  category: { id: 1, name: '后端' },
  tags: [{ id: 1, name: 'Go' }],
}

function mountCard(props) {
  return mount(ArticleCard, {
    props: { idx: 0, article, ...props },
    global: { stubs: { RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' } } },
  })
}

describe('首页文章卡片', () => {
  // 曾经写的是 article.is_top === 1, 而后端 Article.IsTop 是 bool(JSON 里 true/false),
  // 所以真实接口下置顶标记永远不显示 —— mock 数据当时是 0/1, 恰好掩盖了这个问题
  it('置顶文章显示置顶标记', () => {
    expect(mountCard({ article: { ...article, is_top: true } }).text()).toContain('置顶')
  })

  it('非置顶文章不显示置顶标记', () => {
    expect(mountCard({ article: { ...article, is_top: false } }).text()).not.toContain('置顶')
    // 字段缺失也不能显示
    expect(mountCard().text()).not.toContain('置顶')
  })

  it('渲染标题、分类与标签, 链接指向详情页', () => {
    const wrapper = mountCard()

    expect(wrapper.text()).toContain('第一篇')
    expect(wrapper.text()).toContain('后端')
    expect(wrapper.text()).toContain('Go')
    expect(wrapper.find('a').attributes('href')).toBe('/article/1')
  })

  // 奇偶下标决定封面图在左还是在右
  it('按下标交替封面图位置', () => {
    expect(mountCard({ idx: 0 }).html()).toContain('md:order-0')
    expect(mountCard({ idx: 1 }).html()).toContain('md:order-1')
  })
})
