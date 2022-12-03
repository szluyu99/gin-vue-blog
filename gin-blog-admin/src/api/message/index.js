import { request } from '@/utils'

export function useMessageApi() {
  return {
    getMessages: (params = {}) => request.get('/message/list', { params }),
    deleteMessages: (data = []) => request.delete(`/message`, { data }),
    updateMessageReview: (data) => request.put('/message/review', data), // 修改评论审核
  }
}
