import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import UserListPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getUsers: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    getRoleOption: vi.fn().mockResolvedValue({ code: 0, data: [] }),
    updateUser: vi.fn().mockResolvedValue({ code: 0 }),
    updateUserDisable: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(UserListPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const user = {
  id: 3,
  username: 'guest',
  is_disable: false,
  info: { nickname: '游客', avatar: 'public/uploaded/a.png' },
  roles: [{ id: 2, name: 'guest' }],
}

describe('用户列表', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getUsers.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.getRoleOption.mockReset().mockResolvedValue({ code: 0, data: [] })
    api.updateUserDisable.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('挂载时拉取角色选项', async () => {
    mountPage()
    await vi.waitFor(() => expect(api.getRoleOption).toHaveBeenCalled())
  })

  it('拉取角色选项失败不抛出', async () => {
    api.getRoleOption.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getRoleOption).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.roleOptions).toEqual([])
  })

  it('编辑把 roles 映射成 role_ids, 昵称从 info 取', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')
    const row = { ...user }

    actions.render(row)[0].props.onClick()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.modalForm.role_ids).toEqual([2])
    expect(wrapper.vm.modalForm.nickname).toBe('游客')
  })

  // 没有任何角色的用户 roles 可能是 null, 以前直接 .map 会抛在 onClick 里
  it('没有角色的用户点编辑不会抛异常', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')
    const row = { ...user, roles: null }

    expect(() => actions.render(row)[0].props.onClick()).not.toThrow()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.modalForm.role_ids).toEqual([])
  })

  it('禁用开关成功后刷新列表, 失败回滚', async () => {
    const wrapper = mountPage()
    const row = { ...user }

    await wrapper.vm.handleUpdateDisable(row)
    expect(api.updateUserDisable).toHaveBeenCalledWith(3, true)
    expect(row.is_disable).toBe(true)
    expect(row.publishing).toBe(false)

    api.updateUserDisable.mockRejectedValue(new Error('boom'))
    await wrapper.vm.handleUpdateDisable(row)
    expect(row.is_disable).toBe(true)
  })
})
