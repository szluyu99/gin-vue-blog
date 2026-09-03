import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import RolePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getRoles: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    getMenuOption: vi.fn(),
    getResourceOption: vi.fn(),
    saveOrUpdateRole: vi.fn().mockResolvedValue({ code: 0 }),
    deleteRole: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

const menus = [{ key: 1, label: '系统管理', children: [{ key: 2, label: '菜单管理' }] }]
const resources = [{ key: 10, label: '文章模块', children: [{ key: 11, label: '新增文章' }] }]

function mountPage() {
  return mount(RolePage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

describe('角色管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getMenuOption.mockReset().mockResolvedValue({ code: 0, data: menus })
    api.getResourceOption.mockReset().mockResolvedValue({ code: 0, data: resources })
    api.saveOrUpdateRole.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() }
  })

  it('新建角色会预取菜单与资源选项', async () => {
    // 回归 A11: 原来两棵树只在 modalAction === 'edit' 时渲染, option 预取还被注释着,
    // 用户必须先建角色再编辑一次权限
    const wrapper = mountPage()
    await wrapper.vm.handleAddRole()

    expect(api.getMenuOption).toHaveBeenCalled()
    expect(api.getResourceOption).toHaveBeenCalled()
    expect(wrapper.vm.modalAction).toBe('add')
    expect(wrapper.vm.menuOption).toEqual(menus)
    expect(wrapper.vm.resourceOption).toEqual(resources)
  })

  it('新建表单里 menu_ids / resource_ids 是数组', async () => {
    const wrapper = mountPage()
    await wrapper.vm.handleAddRole()

    expect(wrapper.vm.modalForm.menu_ids).toEqual([])
    expect(wrapper.vm.modalForm.resource_ids).toEqual([])
  })

  it('新建时勾选的权限会一起提交', async () => {
    const wrapper = mountPage()
    await wrapper.vm.handleAddRole()

    wrapper.vm.modalForm.name = '编辑'
    wrapper.vm.modalForm.label = 'editor'
    wrapper.vm.modalForm.menu_ids = [1, 2]
    wrapper.vm.modalForm.resource_ids = [11]

    await wrapper.vm.handleSave()

    expect(api.saveOrUpdateRole).toHaveBeenCalledWith(
      expect.objectContaining({
        name: '编辑',
        label: 'editor',
        menu_ids: [1, 2],
        resource_ids: [11],
      }),
    )
  })

  it('选项已拉取过就不重复请求', async () => {
    const wrapper = mountPage()
    await wrapper.vm.handleAddRole()
    await wrapper.vm.handleAddRole()

    expect(api.getMenuOption).toHaveBeenCalledTimes(1)
    expect(api.getResourceOption).toHaveBeenCalledTimes(1)
  })
})
