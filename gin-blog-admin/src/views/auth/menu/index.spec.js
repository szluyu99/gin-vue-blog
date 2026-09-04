import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import MenuPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getMenus: vi.fn().mockResolvedValue({ code: 0, data: [] }),
    saveOrUpdateMenu: vi.fn().mockResolvedValue({ code: 0 }),
    deleteMenu: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(MenuPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const catalogue = { id: 1, name: '文章管理', parent_id: 0, is_catalogue: true, component: 'Layout', redirect: '/article/list', keep_alive: false, is_hidden: false }
const topMenu = { id: 2, name: '首页', parent_id: 0, is_catalogue: false, component: '/home', redirect: '/home', keep_alive: true, is_hidden: false }

describe('菜单管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getMenus.mockReset().mockResolvedValue({ code: 0, data: [] })
    api.saveOrUpdateMenu.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('顶部新增: 父 id 为 0, 默认目录, 组件路径为 Layout', async () => {
    const wrapper = mountPage()
    wrapper.vm.handleClickAdd()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.modalForm.parent_id).toBe(0)
    expect(wrapper.vm.modalForm.is_catalogue).toBe(true)
    expect(wrapper.vm.modalForm.component).toBe('Layout')
  })

  it('行内新增子菜单: 父 id 为该行, 不是目录, 组件路径清空', async () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'actions')
    // 第一个按钮是「新增」(只在非目录的一级菜单上显示)
    column.render(topMenu)[0].props.onClick()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.modalForm.parent_id).toBe(2)
    expect(wrapper.vm.modalForm.is_catalogue).toBe(false)
    expect(wrapper.vm.modalForm.component).toBe('')
  })

  // initForm 是模块级对象, 两种新增都在改它, 顺序不同不能互相污染
  it('两种新增交替点击互不污染', async () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'actions')

    column.render(topMenu)[0].props.onClick()
    await wrapper.vm.$nextTick()
    wrapper.vm.handleClickAdd()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.modalForm.parent_id).toBe(0)
    expect(wrapper.vm.modalForm.component).toBe('Layout')

    column.render(topMenu)[0].props.onClick()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.modalForm.parent_id).toBe(2)
    expect(wrapper.vm.modalForm.component).toBe('')
  })

  it('保活开关提交整行并在失败时回滚', async () => {
    const wrapper = mountPage()
    const row = { ...topMenu, keep_alive: false }

    await wrapper.vm.handleUpdateKeepAlive(row)
    expect(api.saveOrUpdateMenu).toHaveBeenCalledWith(expect.objectContaining({ id: 2, keep_alive: true }))
    expect(row.keep_alive).toBe(true)
    expect(row.publishing).toBe(false)

    api.saveOrUpdateMenu.mockRejectedValue(new Error('boom'))
    await wrapper.vm.handleUpdateKeepAlive(row)
    expect(row.keep_alive).toBe(true) // 回滚
  })

  it('隐藏开关失败时回滚', async () => {
    api.saveOrUpdateMenu.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    const row = { ...topMenu, is_hidden: false }

    await wrapper.vm.handleUpdateHidden(row)

    expect(row.is_hidden).toBe(false)
    expect(window.$message.success).not.toHaveBeenCalled()
  })

  it('目录行不展示跳转路径与组件路径', () => {
    const wrapper = mountPage()
    const redirect = wrapper.vm.columns.find(e => e.key === 'redirect')
    const component = wrapper.vm.columns.find(e => e.key === 'component')

    expect(redirect.render(catalogue).children).toBe('-')
    expect(component.render(catalogue).children).toBe('-')
    expect(redirect.render(topMenu).children).toBe('/home')
    expect(component.render(topMenu).children).toBe('/home')
  })
})
