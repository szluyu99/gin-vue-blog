import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import TalkingCarousel from './TalkingCarousel.vue'

vi.mock('@/api', () => ({
  default: {
    getTalks: vi.fn(),
  },
}))

// 一言接口会真的发请求, 测试里固定住, 别让网络进来
vi.mock('@/utils', () => ({
  getOneSentence: vi.fn().mockResolvedValue('一言文案'),
  getRandomSentence: vi.fn(() => '内置文案'),
}))

describe('首页说说轮播', () => {
  beforeEach(() => {
    api.getTalks.mockReset()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('有说说时轮播真实内容, 而不是一言', async () => {
    api.getTalks.mockResolvedValue({
      code: 0,
      data: { page_data: [{ id: 1, content: '第一条说说' }, { id: 2, content: '第二条说说' }], total: 2 },
    })
    vi.useFakeTimers()

    const wrapper = mount(TalkingCarousel, { global: { stubs: { RouterLink: true } } })
    await vi.advanceTimersByTimeAsync(0)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('第一条说说')

    // 5 秒后轮到第二条
    await vi.advanceTimersByTimeAsync(5000)
    expect(wrapper.text()).toContain('第二条说说')
  })

  it('没有说说时退回一言兜底', async () => {
    api.getTalks.mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    vi.useFakeTimers()

    const wrapper = mount(TalkingCarousel, { global: { stubs: { RouterLink: true } } })
    await vi.advanceTimersByTimeAsync(0)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('一言文案')
  })
})
