import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useUserStore } from '@/store'
import ProfilePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getUserInfo: vi.fn(),
    updateCurrent: vi.fn().mockResolvedValue({ code: 0 }),
    updateCurrentPassword: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

const remoteUser = {
  id: 1,
  nickname: '管理员',
  avatar: 'public/uploaded/a.png',
  intro: '简介',
  website: 'https://a.com',
}

function mountPage() {
  return mount(ProfilePage, {
    global: { stubs: { CommonPage: { template: '<div><slot /></div>' } } },
  })
}

// NTabs 只渲染激活的那个面板, 「修改密码」表单在未激活时挂不上 ref,
// 这里塞一个"校验通过"的表单 ref, 直接验证 updatePassword 自己的逻辑
function passValidation(wrapper) {
  wrapper.vm.passwordFormRef = { validate: cb => cb(null) }
}

describe('个人中心', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getUserInfo.mockReset().mockResolvedValue({ code: 0, data: { ...remoteUser } })
    api.updateCurrent.mockReset().mockResolvedValue({ code: 0 })
    api.updateCurrentPassword.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('挂载后用接口数据同步表单', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.infoForm.nickname).toBe('管理员'))

    expect(api.getUserInfo).toHaveBeenCalled()
    expect(wrapper.vm.infoForm.website).toBe('https://a.com')
  })

  // 回归 A12: 表单里的 avatar 必须是库里的原始值, 不能是跑过 convertImgUrl 的展示地址,
  // 否则提交时会把展示用地址(空头像时是占位图)写回库
  it('表单里的头像是原始相对路径, 不是展示用地址', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.infoForm.nickname).toBe('管理员'))

    const store = useUserStore()
    expect(store.avatar).toBe('/public/uploaded/a.png') // getter 转换过
    expect(wrapper.vm.infoForm.avatar).toBe('public/uploaded/a.png') // 表单里是原始值
  })

  it('保存资料提交表单并刷新用户信息', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.infoForm.nickname).toBe('管理员'))
    api.getUserInfo.mockClear()

    wrapper.vm.infoForm.nickname = '新昵称'
    await wrapper.vm.updateProfile()
    await vi.waitFor(() => expect(api.updateCurrent).toHaveBeenCalled())

    expect(api.updateCurrent).toHaveBeenCalledWith(
      expect.objectContaining({ nickname: '新昵称', avatar: 'public/uploaded/a.png' }),
    )
    expect(window.$message.success).toHaveBeenCalledWith('更新成功!')
    expect(api.getUserInfo).toHaveBeenCalled()
  })

  it('保存资料失败不提示成功', async () => {
    api.updateCurrent.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.infoForm.nickname).toBe('管理员'))

    await wrapper.vm.updateProfile()
    await vi.waitFor(() => expect(api.updateCurrent).toHaveBeenCalled())

    expect(window.$message.success).not.toHaveBeenCalled()
  })

  it('改密码成功后清空表单', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.infoForm.nickname).toBe('管理员'))

    passValidation(wrapper)
    wrapper.vm.passwordForm = { old_password: '123456', new_password: 'abcdef', confirm_password: 'abcdef' }
    await wrapper.vm.updatePassword()
    await vi.waitFor(() => expect(api.updateCurrentPassword).toHaveBeenCalled())

    expect(api.updateCurrentPassword).toHaveBeenCalledWith(
      expect.objectContaining({ old_password: '123456', new_password: 'abcdef' }),
    )
    expect(wrapper.vm.passwordForm.old_password).toBe('')
  })

  it('改密码失败不提示成功也不清空表单', async () => {
    api.updateCurrentPassword.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.infoForm.nickname).toBe('管理员'))

    passValidation(wrapper)
    wrapper.vm.passwordForm = { old_password: '123456', new_password: 'abcdef', confirm_password: 'abcdef' }
    await wrapper.vm.updatePassword()
    await vi.waitFor(() => expect(api.updateCurrentPassword).toHaveBeenCalled())

    expect(window.$message.success).not.toHaveBeenCalled()
    expect(wrapper.vm.passwordForm.old_password).toBe('123456')
  })
})
