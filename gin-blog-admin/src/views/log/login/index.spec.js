import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import LoginLogPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getLoginLogs: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    deleteLoginLogs: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(LoginLogPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const success = { id: 1, user_id: 1, username: 'admin', nickname: '管理员', ip_address: '127.0.0.1', ip_source: '内网IP', status: 1, message: '', created_at: '2026-09-01T10:00:00Z' }
const failed = { id: 2, user_id: 0, username: 'root', nickname: '', ip_address: '1.2.3.4', ip_source: '美国', status: 2, message: '用户名或密码错误', created_at: '2026-09-01T03:00:00Z' }

describe('登录日志', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getLoginLogs.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.deleteLoginLogs.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('挂载后拉取列表', async () => {
    mountPage()
    await vi.waitFor(() => expect(api.getLoginLogs).toHaveBeenCalled())
  })

  it('状态列区分成功与失败', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'status')

    expect(column.render(success).props.type).toBe('success')
    expect(column.render(success).children.default()).toBe('成功')
    expect(column.render(failed).props.type).toBe('error')
    expect(column.render(failed).children.default()).toBe('失败')
  })

  // 登录失败时拿不到昵称(用户可能都不存在), 失败原因在成功时也是空的
  it('昵称与失败原因为空时显示占位符', () => {
    const wrapper = mountPage()
    const nickname = wrapper.vm.columns.find(e => e.key === 'nickname')
    const message = wrapper.vm.columns.find(e => e.key === 'message')

    expect(nickname.render(success).children).toBe('管理员')
    expect(nickname.render(failed).children).toBe('-')
    expect(message.render(failed).children).toBe('用户名或密码错误')
    expect(message.render(success).children).toBe('-')
  })

  it('登录时间精确到秒', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'created_at')

    expect(column.render(success).children.default()).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/)
  })

  it('行内删除不弹二次确认, 直接调接口', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')

    actions.render(failed).props.onPositiveClick()
    await vi.waitFor(() => expect(api.deleteLoginLogs).toHaveBeenCalled())

    expect(api.deleteLoginLogs).toHaveBeenCalledWith(JSON.stringify([2]))
  })

  it('批量删除为空时只提示不发请求', async () => {
    const wrapper = mountPage()

    await wrapper.vm.handleDelete([])

    expect(api.deleteLoginLogs).not.toHaveBeenCalled()
    expect(window.$message.info).toHaveBeenCalled()
  })
})
