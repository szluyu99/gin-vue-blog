import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import WebsitePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
  },
}))

const remoteConfig = {
  website_name: '线上的站名',
  website_author: '线上的作者',
  website_intro: '线上的简介',
  qq: '10000',
  is_comment_review: 'true',
}

function mountPage() {
  return mount(WebsitePage, {
    global: { stubs: { CommonPage: { template: '<div><slot /></div>' } } },
  })
}

describe('网站管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getConfig.mockReset().mockResolvedValue({ code: 0, data: { ...remoteConfig } })
    api.updateConfig.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
    window.$loadingBar = { start: vi.fn(), finish: vi.fn(), error: vi.fn() }
  })

  it('挂载后用后端配置填充表单', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.form.website_name).toBe('线上的站名'))

    expect(api.getConfig).toHaveBeenCalled()
    expect(wrapper.vm.form.qq).toBe('10000')
  })

  // 配置表为空时后端返回 {}, 直接赋值会把表单默认值清成空白
  it('后端返回空配置时保留表单默认值', async () => {
    api.getConfig.mockResolvedValue({ code: 0, data: {} })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getConfig).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.form.website_name).toBe('阵、雨的个人博客')
  })

  it('拉取配置失败不产生未捕获的 rejection, 表单保持默认值', async () => {
    api.getConfig.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getConfig).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.form.website_name).toBe('阵、雨的个人博客')
  })

  it('保存提交当前表单并提示成功', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.form.website_name).toBe('线上的站名'))

    wrapper.vm.form.website_name = '改过的站名'
    // handleSave 现在按 tab 收 ref: 三个表单原来都叫 formRef, 互相覆盖
    wrapper.vm.handleSave(wrapper.vm.basicFormRef)
    await vi.waitFor(() => expect(api.updateConfig).toHaveBeenCalled())

    expect(api.updateConfig).toHaveBeenCalledWith(
      expect.objectContaining({ website_name: '改过的站名', qq: '10000' }),
    )
    expect(window.$message.success).toHaveBeenCalledWith('博客信息更新成功')
    expect(window.$loadingBar.finish).toHaveBeenCalled()
  })

  it('保存失败时走 loadingBar.error, 不提示成功', async () => {
    api.updateConfig.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.form.website_name).toBe('线上的站名'))

    wrapper.vm.handleSave(wrapper.vm.basicFormRef)
    await vi.waitFor(() => expect(window.$loadingBar.error).toHaveBeenCalled())

    expect(window.$message.success).not.toHaveBeenCalled()
  })
})
