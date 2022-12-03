import { resolveToken } from '../utils'

const users = {
  admin: {
    id: 1,
    username: '大脸怪(admin)',
    avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
    email: 'Ronnie@123.com',
    role: ['admin'],
  },
  editor: {
    id: 2,
    username: '大脸怪(editor)',
    avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
    email: 'Ronnie@123.com',
    role: ['editor'],
  },
  tester: {
    id: 3,
    username: '测试管理员(test admin)',
    avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
    role: ['admin'],
  },
  guest: {
    id: 4,
    username: '访客(guest)',
    avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
    role: ['admin'],
  },
}
// 根据 token 获取用户信息
export default [
  {
    url: '/api/user',
    method: 'get',
    response: ({ headers }) => {
      const token = resolveToken(headers?.authorization)
      return {
        code: 0,
        data: {
          ...(users[token] || users.guest),
        },
      }
    },
  },
]
