import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useUserStore } from '@/store'
import UserPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getUser: vi.fn(),
    updateUser: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

const push = vi.fn()
vi.mock('vue-router', () => ({ useRouter: () => ({ push }) }))

const userData = {
  id: 7,
  nickname: '阵雨',
  avatar: 'public/uploaded/a.jpg',
  website: 'https://a.com',
  intro: '我的简介',
  email: 'a@b.com',
  article_like_set: [],
  comment_like_set: [],
}

function mountPage() {
  return mount(UserPage, {
    global: { stubs: { BannerPage: { template: '<div><slot /></div>' } } },
  })
}

describe('个人中心', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    push.mockClear()
    window.$message = { success: vi.fn(), error: vi.fn(), warning: vi.fn() }
  })

  it('getUserInfo 返回后表单会同步真实资料', async () => {
    // 回归 FE1: form 是 setup 阶段拷的快照, getUserInfo 是 onMounted 里异步取的,
    // 不重新同步的话页面显示默认头像和空表单
    api.getUser.mockResolvedValue({ code: 0, data: userData })
    const store = useUserStore()
    store.setToken('t')

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.form.nickname).toBe('阵雨'))

    expect(wrapper.vm.form.avatar).toBe('public/uploaded/a.jpg')
    expect(wrapper.vm.form.intro).toBe('我的简介')
    expect(wrapper.vm.form.website).toBe('https://a.com')
    expect(wrapper.vm.form.email).toBe('a@b.com')
    expect(push).not.toHaveBeenCalled()
  })

  it('提交时发送的是原始相对路径, 不是展示地址', async () => {
    api.getUser.mockResolvedValue({ code: 0, data: userData })
    const store = useUserStore()
    store.setToken('t')

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.form.nickname).toBe('阵雨'))
    await wrapper.vm.updateUserInfo()

    // FE2: 存进库的必须是 public/uploaded/a.jpg, 不能是带域名或带前导 / 的展示地址
    expect(api.updateUser).toHaveBeenCalledWith(
      expect.objectContaining({ avatar: 'public/uploaded/a.jpg' }),
    )
  })

  it('未登录时跳回首页', async () => {
    api.getUser.mockResolvedValue({ code: 0, data: { ...userData, id: '' } })
    mountPage()
    await vi.waitFor(() => expect(push).toHaveBeenCalledWith('/'))
  })
})
