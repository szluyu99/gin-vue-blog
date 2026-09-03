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

  // 回归 A15: 开关原来写死 checkedValue: 1 / uncheckedValue: 0, 而后端 is_disable 是布尔值,
  // 已禁用的角色也一直显示成关闭
  it('禁用开关按布尔值渲染', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'is_disable')
    const vnode = column.render({ id: 1, is_disable: true })

    expect(vnode.props.value).toBe(true)
    expect(vnode.props.checkedValue).toBeUndefined()
    expect(vnode.props.uncheckedValue).toBeUndefined()
  })

  it('切换禁用会带上原有的资源与菜单权限', async () => {
    const wrapper = mountPage()
    const row = { id: 3, name: '编辑', label: 'editor', is_disable: false, resource_ids: [11], menu_ids: [1, 2] }

    await wrapper.vm.handleUpdateDisable(row)

    // 后端 UpdateRole 会整体替换关联, 不带上这两个字段等于清空该角色的权限
    expect(api.saveOrUpdateRole).toHaveBeenCalledWith({
      id: 3,
      name: '编辑',
      label: 'editor',
      is_disable: true,
      resource_ids: [11],
      menu_ids: [1, 2],
    })
    expect(row.is_disable).toBe(true)
  })

  it('切换禁用失败时不改变本地状态', async () => {
    api.saveOrUpdateRole.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    const row = { id: 3, name: '编辑', label: 'editor', is_disable: false }

    await wrapper.vm.handleUpdateDisable(row)

    expect(row.is_disable).toBe(false)
    expect(row.publishing).toBe(false)
  })
})
