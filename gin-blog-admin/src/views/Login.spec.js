import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { getLocal, setLocal } from '@/utils'
import Login from './Login.vue'

vi.mock('@/api', () => ({
  default: {
    login: vi.fn(),
    getUserInfo: vi.fn().mockResolvedValue({ code: 0, data: { id: 1, nickname: 'admin' } }),
  },
}))

vi.mock('@/router', async importOriginal => ({
  ...await importOriginal(),
  addDynamicRoutes: vi.fn().mockResolvedValue(undefined),
}))

const push = vi.fn()
vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ query: {} }),
  useRouter: () => ({ push }),
}))

function mountPage() {
  return mount(Login, {
    global: { stubs: { AppPage: { template: '<div><slot /></div>' } } },
  })
}

describe('登录页', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    sessionStorage.clear()
    push.mockClear()
    window.$message = { success: vi.fn(), error: vi.fn(), warning: vi.fn(), info: vi.fn() }
    api.login.mockReset()
    api.login.mockResolvedValue({ code: 0, data: { token: 'tk' } })
  })

  it('不勾选记住我时不会保存账号密码', async () => {
    // 回归 A8: isRemember 是 useStorage 返回的 Ref, 少写 .value 会恒为真,
    // 取消勾选也照样把账号密码写进 localStorage
    setLocal('loginInfo', { username: 'old', password: 'old' })
    const wrapper = mountPage()
    wrapper.vm.isRemember = false

    await wrapper.vm.handleLogin()
    await vi.waitFor(() => expect(push).toHaveBeenCalled())

    expect(getLocal('loginInfo')).toBeNull()
  })

  it('勾选记住我时保存账号密码', async () => {
    const wrapper = mountPage()
    wrapper.vm.isRemember = true
    wrapper.vm.loginForm.username = 'admin'
    wrapper.vm.loginForm.password = '123456'

    await wrapper.vm.handleLogin()
    await vi.waitFor(() => expect(push).toHaveBeenCalled())

    expect(getLocal('loginInfo')).toEqual({ username: 'admin', password: '123456' })
  })

  it('登录失败时 loading 会复位, 且不产生未捕获的 rejection', async () => {
    // doLogin 里只有 try/finally, 调用处必须 .catch 兜住
    api.login.mockRejectedValue({ code: 1004, message: '用户名或密码错误' })
    const wrapper = mountPage()

    await expect(wrapper.vm.handleLogin()).resolves.toBeUndefined()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))
    expect(window.$message.success).not.toHaveBeenCalled()
  })

  it('用户名或密码为空时直接提示', async () => {
    const wrapper = mountPage()
    wrapper.vm.loginForm.username = ''

    await wrapper.vm.handleLogin()

    expect(window.$message.warning).toHaveBeenCalledWith('请输入用户名和密码')
    expect(api.login).not.toHaveBeenCalled()
  })
})
