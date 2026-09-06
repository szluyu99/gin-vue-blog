import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useAppStore, useUserStore } from '@/store'
import LoginModal from './LoginModal.vue'

vi.mock('@/api', () => ({
  default: {
    login: vi.fn(),
    getUser: vi.fn().mockResolvedValue({
      code: 0,
      data: { id: 1, nickname: 'me', article_like_set: [], comment_like_set: [] },
    }),
    logout: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

// UModal 用 Teleport 挂到 body, 所以从 document 上取输入框
function inputs() {
  return [...document.querySelectorAll('input')]
}

async function openModal(wrapper, appStore) {
  appStore.setLoginFlag(true)
  await wrapper.vm.$nextTick()
  const [username, password] = inputs()
  return { username: username.value, password: password.value }
}

// v-model 靠 input 事件, 直接改 DOM 值再派发
async function fill(username, password) {
  const [u, p] = inputs()
  u.value = username
  u.dispatchEvent(new Event('input'))
  p.value = password
  p.dispatchEvent(new Event('input'))
}

describe('登录框的用户名回填', () => {
  let appStore
  let userStore

  beforeEach(() => {
    document.body.innerHTML = ''
    setActivePinia(createPinia())
    appStore = useAppStore()
    userStore = useUserStore()
    api.login.mockReset().mockResolvedValue({ code: 0, data: { token: 'tk' } })
    window.$message = { warning: vi.fn() }
    window.$notify = { success: vi.fn() }
  })

  // 回归: 原来写死了作者的示例账号 test@qq.com / 11111
  it('没登录成功过时是空的, 不带任何示例账号', async () => {
    const wrapper = mount(LoginModal, { attachTo: document.body })
    expect(await openModal(wrapper, appStore)).toEqual({ username: '', password: '' })
  })

  it('登录成功后记住用户名, 下次打开自动填上但密码留空', async () => {
    const wrapper = mount(LoginModal, { attachTo: document.body })
    appStore.setLoginFlag(true)
    await wrapper.vm.$nextTick()

    await fill('me@example.com', 'pwd')
    await wrapper.vm.$nextTick()

    // 模态框里的按钮顺序: 关闭 / 登录 / 立即注册 / 忘记密码
    document.querySelectorAll('button')[1].click()
    await vi.waitFor(() => expect(api.login).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(api.login).toHaveBeenCalledWith({ username: 'me@example.com', password: 'pwd' })
    expect(userStore.lastUsername).toBe('me@example.com')

    // 登录成功后模态框自己关掉, 再打开时才会重新回填
    await vi.waitFor(() => expect(appStore.loginFlag).toBe(false))
    expect(await openModal(wrapper, appStore)).toEqual({
      username: 'me@example.com',
      password: '',
    })
  })

  it('退出登录不忘记用户名', () => {
    userStore.setLastUsername('me@example.com')
    userStore.resetLoginState()
    expect(userStore.token).toBeNull()
    expect(userStore.lastUsername).toBe('me@example.com')
  })

  it('刚注册完优先填注册用的邮箱', async () => {
    userStore.setLastUsername('old@example.com')
    appStore.setPrefillUsername('new@example.com')
    const wrapper = mount(LoginModal, { attachTo: document.body })
    expect((await openModal(wrapper, appStore)).username).toBe('new@example.com')
  })
})
