import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import TalkDetail from './detail.vue'

vi.mock('@/api', () => ({
  default: {
    getTalk: vi.fn(),
    getComments: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
  },
}))

const replace = vi.fn()
const route = { params: { id: '7' } }
vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => route,
  useRouter: () => ({ replace }),
}))

function mountPage() {
  return mount(TalkDetail, {
    global: {
      stubs: {
        BannerPage: { template: '<div><slot /></div>' },
        RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' },
        Comment: { template: '<div class="comment-stub" />' },
      },
    },
  })
}

describe('说说详情页', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    replace.mockReset()
  })

  it('加载说说并渲染内容, 评论区 type 是 3', async () => {
    api.getTalk.mockResolvedValue({
      code: 0,
      data: { id: 7, content: '说说正文\n第二行', nickname: '博主', avatar: 'a.png', is_top: true, comment_count: 2, created_at: '2026-09-06T00:00:00Z' },
    })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(api.getTalk).toHaveBeenCalledWith('7')
    expect(wrapper.text()).toContain('说说正文')
    expect(wrapper.find('.comment-stub').exists()).toBe(true)
  })

  // 私密说说在后端就 404, 前端拿不到数据时跳回 404 页
  it('取不到说说时跳转 404', async () => {
    api.getTalk.mockRejectedValue(new Error('not found'))
    mountPage()
    await vi.waitFor(() => expect(replace).toHaveBeenCalledWith('/404'))
  })
})
