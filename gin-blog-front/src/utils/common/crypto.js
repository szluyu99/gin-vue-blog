import CryptoJS from 'crypto-js'

const CryptoSecret = '__SecretKey__'

/**
 * 加密数据
 * @param {any} data - 数据
 */
export function encrypto(data) {
  const newData = JSON.stringify(data)
  return CryptoJS.AES.encrypt(newData, CryptoSecret).toString()
}

/**
 * 解密数据
 * @param {String} cipherText - 密文
 */
export function decrypto(cipherText) {
  const bytes = CryptoJS.AES.decrypt(cipherText, CryptoSecret)
  const originalText = bytes.toString(CryptoJS.enc.Utf8)
  if (originalText)
    return JSON.parse(originalText)
  return null
}
