import { request } from '@/utils'

export default {
  getOperationLogs: (params = {}) => request.get('/operation/log/list', { params }),
  deleteOperationLogs: (data = []) => request.delete(`/operation/log`, { data }),
}
