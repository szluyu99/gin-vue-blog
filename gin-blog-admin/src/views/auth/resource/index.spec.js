import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import ResourcePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getResources: vi.fn().mockResolvedValue({ code: 0, data: [] }),
    saveOrUpdateResource: vi.fn().mockResolvedValue({ code: 0 }),
    deleteResource: vi.fn().mockResolvedValue({ code: 0 }),
    updateResourceAnonymous: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(ResourcePage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

// 模块(父资源)没有 url / 请求方式, 列表里用 children 区分
const module_ = { id: 1, name: '文章模块', parent_id: 0, url: '', request_method: '', is_anonymous: false, children: [] }
const item = { id: 2, name: '文章列表', parent_id: 1, url: '/article/list', request_method: 'GET', is_anonymous: false }

describe('接口管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getResources.mockReset().mockResolvedValue({ code: 0, data: [] })
    api.saveOrUpdateResource.mockReset().mockResolvedValue({ code: 0 })
    api.updateResourceAnonymous.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('新增模块保存成功后关闭弹窗并刷新列表', async () => {
    const wrapper = mountPage()
    wrapper.vm.handleAddModule()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.moduleModalVisible).toBe(true)

    wrapper.vm.modalForm.name = '新模块'
    api.getResources.mockClear()
    await wrapper.vm.handleModuleSave()

    expect(api.saveOrUpdateResource).toHaveBeenCalledWith(expect.objectContaining({ name: '新模块' }))
    expect(window.$message.success).toHaveBeenCalledWith('新增成功')
    expect(wrapper.vm.moduleModalVisible).toBe(false)
    await vi.waitFor(() => expect(api.getResources).toHaveBeenCalled())
  })

  // 回归: 以前 handleSave() 不 await 就把弹窗关掉, 保存失败也会关
  it('模块保存失败时弹窗不关, 也不提示成功', async () => {
    api.saveOrUpdateResource.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    wrapper.vm.handleAddModule()
    await wrapper.vm.$nextTick()

    await wrapper.vm.handleModuleSave()

    expect(wrapper.vm.moduleModalVisible).toBe(true)
    expect(window.$message.success).not.toHaveBeenCalled()
    expect(wrapper.vm.modalLoading).toBe(false)
  })

  it('编辑模块带上原有数据, 标题走编辑分支', async () => {
    const wrapper = mountPage()
    wrapper.vm.handleEditModule(module_)
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.modalAction).toBe('edit')
    expect(wrapper.vm.modalForm.id).toBe(1)

    await wrapper.vm.handleModuleSave()
    expect(window.$message.success).toHaveBeenCalledWith('编辑成功')
  })

  it('匿名开关成功后保留状态, 失败回滚', async () => {
    const wrapper = mountPage()
    const row = { ...item }

    await wrapper.vm.handleUpdateAnonymous(row)
    expect(api.updateResourceAnonymous).toHaveBeenCalledWith(expect.objectContaining({ id: 2, is_anonymous: true }))
    expect(row.is_anonymous).toBe(true)
    expect(row.publishing).toBe(false)

    api.updateResourceAnonymous.mockRejectedValue(new Error('boom'))
    await wrapper.vm.handleUpdateAnonymous(row)
    expect(row.is_anonymous).toBe(true)
  })

  it('模块行不展示路径/请求方式/匿名开关', () => {
    const wrapper = mountPage()
    const url = wrapper.vm.columns.find(e => e.key === 'url')
    const method = wrapper.vm.columns.find(e => e.key === 'request_method')
    const anonymous = wrapper.vm.columns.find(e => e.key === 'is_anonymous')

    expect(url.render(module_)).toBe('-')
    expect(method.render(module_)).toBe('-')
    expect(anonymous.render(module_)).toBe('-')
    // 接口行正常渲染
    expect(url.render(item).children).toBe('/article/list')
    expect(method.render(item).children.default()).toBe('GET')
  })
})
