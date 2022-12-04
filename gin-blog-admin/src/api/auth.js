import { request } from '@/utils'

export default {
  login: data => request.post('/login', data, { noNeedToken: true }),
  logout: () => request.get('/logout'),
  // 菜单
  getUserMenus: () => request.get('/menu/user/list'), // 获取当前用户的菜单
  getMenus: (params = {}) => request.get('/menu/list', { params }),
  saveOrUpdateMenu: data => request.post('/menu', data),
  deleteMenu: id => request.delete(`/menu/${id}`),
  getMenuOption: () => request.get('/menu/option'),
  // 资源
  getResources: (params = {}) => request.get('/resource/list', { params }),
  saveOrUpdateResource: data => request.post('/resource', data),
  deleteResource: id => request.delete(`/resource/${id}`),
  updateResourceAnonymous: data => request.put('/resource/anonymous', data),
  getResourceOption: () => request.get('/resource/option'),
  // 角色
  getRoles: (params = {}) => request.get('/role/list', { params }),
  saveOrUpdateRole: data => request.post('/role', data),
  deleteRole: (data = []) => request.delete('/role', { data }),
  getRoleOption: () => request.get('/role/option'),
}
