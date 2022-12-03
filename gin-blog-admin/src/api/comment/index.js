import { request } from '@/utils'

export function useCommentApi() {
  return {
    getComments: (params = {}) => request.get('/comment/list', { params }),
    deleteComments: (data = []) => request.delete(`/comment`, { data }),
    updateCommentReview: (data) => request.put('/comment/review', data), // 修改评论审核
  }
}
