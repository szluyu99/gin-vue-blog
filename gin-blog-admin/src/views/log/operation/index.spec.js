import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import OperationLogPage from './index.vue'

const { copy } = vi.hoisted(() => ({ copy: vi.fn() }))
vi.mock('@vueuse/core', async importOriginal => ({
  ...await importOriginal(),
  useClipboard: () => ({ copy }),
}))

vi.mock('@/api', () => ({
  default: {
    getOperationLogs: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    deleteOperationLogs: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(OperationLogPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const log = {
  id: 7,
  opt_module: '文章模块',
  opt_type: '新增',
  request_method: 'POST',
  request_param: '{"title":"x"}',
  response_data: '{"code":0}',
  nickname: 'admin',
  ip_address: '10.0.0.1',
  ip_source: '内网IP',
  created_at: '2026-09-01T00:00:00Z',
}

describe('操作日志', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    copy.mockClear()
    api.getOperationLogs.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.deleteOperationLogs.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('复制的是格式化后的 JSON', () => {
    const wrapper = mountPage()

    wrapper.vm.copyFormatCode('{"a":1}')

    expect(copy).toHaveBeenCalledWith('{\n  "a": 1\n}')
    expect(window.$message.success).toHaveBeenCalled()
  })

  // 回归 A4: 以前是 JSON.stringify(JSON.parse(x)), 空串或非 JSON 会抛
  it('请求参数为空串或非 JSON 时复制不抛异常', () => {
    const wrapper = mountPage()

    expect(() => wrapper.vm.copyFormatCode('')).not.toThrow()
    expect(copy).toHaveBeenLastCalledWith('')

    expect(() => wrapper.vm.copyFormatCode('<html>网关拦截</html>')).not.toThrow()
    expect(copy).toHaveBeenLastCalledWith('<html>网关拦截</html>')
  })

  it('请求方式按方法映射标签颜色', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'request_method')

    expect(column.render(log).props.type).toBe('success')
    expect(column.render({ ...log, request_method: 'GET' }).props.type).toBe('info')
    expect(column.render({ ...log, request_method: 'DELETE' }).props.type).toBe('error')
    // 未知方法退回 info, 不抛异常
    expect(column.render({ ...log, request_method: 'PATCH' }).props.type).toBe('info')
  })

  it('查看带出整行数据, 删除调删除接口', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')

    actions.render(log)[0].props.onClick()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.modalForm).toMatchObject({ id: 7, request_param: '{"title":"x"}' })

    actions.render(log)[1].props.onPositiveClick()
    await vi.waitFor(() => expect(api.deleteOperationLogs).toHaveBeenCalled())
    expect(api.deleteOperationLogs).toHaveBeenCalledWith(JSON.stringify([7]))
  })
})
