import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import CategoryPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getCategorys: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    saveOrUpdateCategory: vi.fn().mockResolvedValue({ code: 0 }),
    deleteCategory: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(CategoryPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const category = { id: 3, name: '项目', article_count: 2, created_at: '2026-09-01T00:00:00Z', updated_at: '2026-09-02T00:00:00Z' }

describe('分类管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getCategorys.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.saveOrUpdateCategory.mockReset().mockResolvedValue({ code: 0 })
    api.deleteCategory.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
    window.$dialog = { confirm: ({ confirm }) => confirm() }
  })

  it('新增走空表单, 保存调新增接口', async () => {
    const wrapper = mountPage()
    wrapper.vm.handleAdd()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.modalForm).toEqual({})

    wrapper.vm.modalForm.name = '新分类'
    await wrapper.vm.handleSave()

    expect(api.saveOrUpdateCategory).toHaveBeenCalledWith({ name: '新分类' })
    expect(window.$message.success).toHaveBeenCalledWith('新增成功')
  })

  it('编辑带上原行数据, 保存调更新接口', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')

    actions.render(category)[0].props.onClick()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.modalForm).toMatchObject({ id: 3, name: '项目' })

    await wrapper.vm.handleSave()
    expect(api.saveOrUpdateCategory).toHaveBeenCalledWith(expect.objectContaining({ id: 3 }))
    expect(window.$message.success).toHaveBeenCalledWith('编辑成功')
  })

  it('行内删除不弹二次确认, 直接调接口', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')

    actions.render(category)[1].props.onPositiveClick()
    await vi.waitFor(() => expect(api.deleteCategory).toHaveBeenCalled())

    expect(api.deleteCategory).toHaveBeenCalledWith(JSON.stringify([3]))
    expect(actions.render(category)[1].children.default().children).toBe('确定删除该分类吗?')
  })

  it('批量删除为空时只提示不发请求', async () => {
    const wrapper = mountPage()

    await wrapper.vm.handleDelete([])

    expect(api.deleteCategory).not.toHaveBeenCalled()
    expect(window.$message.info).toHaveBeenCalled()
  })

  it('保存失败不提示成功', async () => {
    api.saveOrUpdateCategory.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    wrapper.vm.handleAdd()
    await wrapper.vm.$nextTick()

    await wrapper.vm.handleSave()

    expect(window.$message.success).not.toHaveBeenCalled()
    expect(wrapper.vm.modalLoading).toBe(false)
  })
})
