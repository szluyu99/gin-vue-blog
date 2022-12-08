import { lStorage } from '@/utils'

// Token 存储在 LocalStorage 中, 封装其操作
const TOKEN_CODE = 'access_token'
const DURATION = 6 * 60 * 60

export function getToken() {
  return lStorage.get(TOKEN_CODE)
}

export function setToken(token) {
  lStorage.set(TOKEN_CODE, token, DURATION)
}

export function removeToken() {
  lStorage.remove(TOKEN_CODE)
}

// TODO:
// export async function refreshAccessToken() {
//   const tokenItem = lStorage.getItem(TOKEN_CODE)
//   if (!tokenItem) return
//   const { time } = tokenItem
//   // token 生成或者刷新后 30 分钟内不执行刷新
//   if (new Date().getTime() - time <= 1000 * 60 * 30) return
//   try {
//     const res = await api.refreshToken()
//     setToken(res.data.token)
//   } catch (err) {
//     console.error(err)
//   }
// }
