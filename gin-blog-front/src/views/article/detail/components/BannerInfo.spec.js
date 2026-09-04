import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import BannerInfo from './BannerInfo.vue'

const article = {
  id: 1,
  title: '第一篇',
  content: '<h1>标题</h1><p>正文正文正文</p>',
  created_at: '2026-09-01T00:00:00Z',
  updated_at: '2026-09-02T00:00:00Z',
  like_count: 3,
  view_count: 10,
  comment_count: 1,
  category: { id: 1, name: '项目' },
  tags: [{ id: 1, name: 'Go' }],
}

function mountInfo(props) {
  return mount(BannerInfo, {
    props: { article: { ...article, ...props } },
    global: { stubs: { RouterLink: { template: '<a><slot /></a>' } } },
  })
}

describe('文章详情横幅', () => {
  it('统计字数时去掉 HTML 标签', () => {
    const wrapper = mountInfo()

    // 「标题正文正文正文」共 8 个字
    expect(wrapper.vm.wordNum).toBe(8)
    expect(wrapper.vm.readTime).toBe('0分钟')
  })

  it('按 400 字每分钟估算阅读时间', () => {
    const wrapper = mountInfo({ content: '字'.repeat(1200) })

    expect(wrapper.vm.wordNum).toBe(1200)
    expect(wrapper.vm.readTime).toBe('3分钟')
  })

  // 回归: 以前直接对 article.content 调 replace, 正文为空(草稿/导入的文章)会抛
  it('正文为空或缺失时不抛异常', () => {
    expect(() => mountInfo({ content: null })).not.toThrow()
    expect(() => mountInfo({ content: undefined })).not.toThrow()
    expect(mountInfo({ content: '' }).vm.wordNum).toBe(0)
  })
})
