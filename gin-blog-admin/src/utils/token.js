import { getLocal, removeLocal, setLocal } from './local'

const TOKEN_KEY = 'blog_admin_token'

/** Token 过期时间: 6 小时 */
const TOKEN_EXPIRE = 60 * 60 * 24 * 7

export function getToken() {
  return getLocal(TOKEN_KEY)
}

export function setToken(token) {
  setLocal(TOKEN_KEY, token, TOKEN_EXPIRE)
}

export function removeToken() {
  removeLocal(TOKEN_KEY)
}
