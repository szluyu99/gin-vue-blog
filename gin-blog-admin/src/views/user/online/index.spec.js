import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import OnlineUserPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getOnlineUsers: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    kickOutUser: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(OnlineUserPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

describe('在线用户', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  // 曾经的 bug: 写的是 row.info.avatar / row.info.nickname, 没有可选链。
  // 在线数据是后端从 Redis 反序列化的 UserAuth, 没有 Preload("UserInfo"),
  // info 为空时列的 render 抛 TypeError, NDataTable 整张表都渲染不出来
  it('info 为空时列渲染不抛错', () => {
    const wrapper = mountPage()
    const columns = wrapper.vm.columns

    for (const row of [{ id: 1 }, { id: 1, info: null }, { id: 1, info: {} }]) {
      for (const col of columns) {
        if (typeof col.render === 'function') {
          expect(() => col.render(row)).not.toThrow()
        }
      }
    }
  })

  it('有 info 时展示昵称', () => {
    const wrapper = mountPage()
    const nickname = wrapper.vm.columns.find(e => e.key === 'nickname')

    expect(nickname.render({ id: 1, info: { nickname: '游客' } }).children).toBe('游客')
    // 缺昵称时退回「未知」而不是空白
    expect(nickname.render({ id: 1, info: {} }).children).toBe('未知')
  })
})
