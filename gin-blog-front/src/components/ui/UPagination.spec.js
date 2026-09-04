import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import UPagination from './UPagination.vue'

function mountPager(props) {
  return mount(UPagination, { props: { page: 1, pageCount: 3, ...props } })
}

function labels(wrapper) {
  return wrapper.findAll('button').map(b => b.text()).filter(Boolean)
}

describe('前台分页器', () => {
  it('总页数不多时把页码全部列出来', () => {
    const wrapper = mountPager({ pageCount: 3 })
    expect(labels(wrapper)).toEqual(['1', '2', '3'])
    expect(wrapper.find('[aria-current="page"]').text()).toBe('1')
  })

  it('首末页时上一页/下一页禁用', () => {
    const first = mountPager({ page: 1, pageCount: 3 })
    expect(first.find('[aria-label="上一页"]').attributes('disabled')).toBeDefined()
    expect(first.find('[aria-label="下一页"]').attributes('disabled')).toBeUndefined()

    const last = mountPager({ page: 3, pageCount: 3 })
    expect(last.find('[aria-label="下一页"]').attributes('disabled')).toBeDefined()
  })

  it('点击页码与翻页按钮抛出 update:page', async () => {
    const wrapper = mountPager({ page: 2, pageCount: 5 })

    await wrapper.find('[aria-label="下一页"]').trigger('click')
    await wrapper.find('[aria-label="上一页"]').trigger('click')
    await wrapper.findAll('button')[3].trigger('click') // 页码 3

    expect(wrapper.emitted('update:page')).toEqual([[3], [1], [3]])
  })

  it('点当前页或越界不抛事件', async () => {
    const wrapper = mountPager({ page: 1, pageCount: 3 })

    await wrapper.find('[aria-current="page"]').trigger('click')
    await wrapper.find('[aria-label="上一页"]').trigger('click')

    expect(wrapper.emitted('update:page')).toBeUndefined()
  })

  it('页码很多时只显示当前页附近的一段, 首末页单独补上', () => {
    const wrapper = mountPager({ page: 10, pageCount: 20, window: 5 })

    // 1 ... 8 9 10 11 12 ... 20
    expect(labels(wrapper)).toEqual(['1', '8', '9', '10', '11', '12', '20'])
    expect(wrapper.findAll('span').filter(s => s.text() === '...')).toHaveLength(2)
  })

  it('靠边时窗口不越界, 长度保持不变', () => {
    expect(labels(mountPager({ page: 1, pageCount: 20, window: 5 })))
      .toEqual(['1', '2', '3', '4', '5', '20'])
    expect(labels(mountPager({ page: 20, pageCount: 20, window: 5 })))
      .toEqual(['1', '16', '17', '18', '19', '20'])
  })
})
