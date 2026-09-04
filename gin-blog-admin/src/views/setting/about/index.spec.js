import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import AboutPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getAbout: vi.fn(),
    updateAbout: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

// md-editor-v3 在 jsdom 下初始化较重, 用桩替掉, 只关心 v-model 的值
vi.mock('md-editor-v3', () => ({
  MdEditor: { name: 'MdEditor', props: ['modelValue'], template: '<div class="md-editor-stub" />' },
}))
vi.mock('md-editor-v3/lib/style.css', () => ({}))

function mountPage() {
  return mount(AboutPage, {
    global: { stubs: { CommonPage: { template: '<div><slot /></div>' } } },
  })
}

describe('关于我', () => {
  beforeEach(() => {
    api.getAbout.mockReset().mockResolvedValue({ code: 0, data: '# 关于我' })
    api.updateAbout.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('挂载后填充内容', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.aboutContent).toBe('# 关于我'))
  })

  it('拉取失败不抛出, 内容保持空串', async () => {
    api.getAbout.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getAbout).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.aboutContent).toBe('')
  })

  it('后端返回 null 时内容是空串而不是 null', async () => {
    api.getAbout.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getAbout).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.aboutContent).toBe('')
  })

  it('保存提交当前内容并复位按钮 loading', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.aboutContent).toBe('# 关于我'))

    wrapper.vm.aboutContent = '# 改过了'
    await wrapper.vm.handleSave()

    expect(api.updateAbout).toHaveBeenCalledWith({ content: '# 改过了' })
    expect(window.$message.success).toHaveBeenCalledWith('更新成功')
    expect(wrapper.vm.btnLoading).toBe(false)
  })

  it('保存失败不提示成功, loading 也要复位', async () => {
    api.updateAbout.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.aboutContent).toBe('# 关于我'))

    await wrapper.vm.handleSave()

    expect(window.$message.success).not.toHaveBeenCalled()
    expect(wrapper.vm.btnLoading).toBe(false)
  })
})
