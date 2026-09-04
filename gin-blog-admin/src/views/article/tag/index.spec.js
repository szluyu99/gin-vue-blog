import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import TagPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getTags: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    saveOrUpdateTag: vi.fn().mockResolvedValue({ code: 0 }),
    deleteTag: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(TagPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const tag = { id: 4, name: 'Go', article_count: 3, created_at: '2026-09-01T00:00:00Z', updated_at: '2026-09-02T00:00:00Z' }

describe('标签管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getTags.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.saveOrUpdateTag.mockReset().mockResolvedValue({ code: 0 })
    api.deleteTag.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
    window.$dialog = { confirm: ({ confirm }) => confirm() }
  })

  it('新增走空表单, 保存调新增接口', async () => {
    const wrapper = mountPage()
    wrapper.vm.handleAdd()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.modalForm).toEqual({})

    wrapper.vm.modalForm.name = 'Vue'
    await wrapper.vm.handleSave()

    expect(api.saveOrUpdateTag).toHaveBeenCalledWith({ name: 'Vue' })
    expect(window.$message.success).toHaveBeenCalledWith('新增成功')
  })

  it('编辑带上原行数据, 删除文案说的是标签', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')

    actions.render(tag)[0].props.onClick()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.modalForm).toMatchObject({ id: 4, name: 'Go' })

    expect(actions.render(tag)[1].children.default().children).toBe('确定删除该标签吗?')
  })

  it('行内删除直接调接口', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')

    actions.render(tag)[1].props.onPositiveClick()
    await vi.waitFor(() => expect(api.deleteTag).toHaveBeenCalled())

    expect(api.deleteTag).toHaveBeenCalledWith(JSON.stringify([4]))
  })

  it('批量删除为空时只提示不发请求', async () => {
    const wrapper = mountPage()

    await wrapper.vm.handleDelete([])

    expect(api.deleteTag).not.toHaveBeenCalled()
    expect(window.$message.info).toHaveBeenCalled()
  })

  it('保存失败不提示成功', async () => {
    api.saveOrUpdateTag.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    wrapper.vm.handleAdd()
    await wrapper.vm.$nextTick()

    await wrapper.vm.handleSave()

    expect(window.$message.success).not.toHaveBeenCalled()
    expect(wrapper.vm.modalLoading).toBe(false)
  })
})
