const CryptoSecret = '__SecretKey__'

/**
 * 存储序列化后的数据到 LocalStorage
 * @param {string} key
 * @param {any} value 对象需要序列化
 * @param {number} expire
 */
export function setLocal(key, value, expire = 60 * 60 * 24 * 7) {
  const data = JSON.stringify({
    value,
    time: Date.now(),
    expire: expire ? new Date().getTime() + expire * 1000 : null,
  })
  window.localStorage.setItem(key, encrypto(data)) // 加密存储
}

/**
 * 从 LocalStorage 中获取数据, 解密后反序列化, 根据是否过期来返回
 * @param {string} key
 */
export function getLocal(key) {
  const encryptedVal = window.localStorage.getItem(key)
  if (encryptedVal) {
    const val = decrypto(encryptedVal) // 解密
    const { value, expire } = JSON.parse(val)
    // 未过期则返回
    if (!expire || expire > new Date().getTime()) {
      return value
    }
  }
  // 过期则移除
  removeLocal(key)
  return null
}

export function removeLocal(key) {
  window.localStorage.removeItem(key)
}

export function clearLocal() {
  window.localStorage.clear()
}

/**
 * 加密数据: Base64 加密
 * @param {any} data - 数据
 */
function encrypto(data) {
  const newData = JSON.stringify(data)
  const encryptedData = btoa(CryptoSecret + newData)
  return encryptedData
}

/**
 * 解密数据: Base64 解密
 * @param {string} cipherText - 密文
 */
function decrypto(cipherText) {
  const decryptedData = atob(cipherText)
  const originalText = decryptedData.replace(CryptoSecret, '')
  try {
    const parsedData = JSON.parse(originalText)
    return parsedData
  }
  catch (error) {
    return null
  }
}
