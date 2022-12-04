// 浏览器持久化类封装 (LocalStorage 和 SessionStorage API 完全一样)

class Storage {
  /**
   * @param {*} prefixKey 前缀
   * @param {*} storage LocalStorage 或 SessionStorage
   */
  constructor({ prefixKey, storage }) {
    this.storage = storage
    this.prefixKey = prefixKey
  }

  getKey(key) {
    return `${this.prefixKey}${key}`.toUpperCase()
  }

  /**
   * @param {*} key 键
   * @param {*} value 值
   * @param {*} expire  过期时间
   */
  set(key, value, expire) {
    const stringData = JSON.stringify({
      value, // 数据
      time: Date.now(), // 存储时间
      expire: expire ? new Date().getTime() + expire * 1000 : null, // 过期时间
    })
    this.storage.setItem(this.getKey(key), stringData)
  }

  get(key) {
    const { value } = this.getItem(key, {})
    return value
  }

  getItem(key, def = null) {
    const val = this.storage.getItem(this.getKey(key))
    if (!val)
      return def
    try {
      const { value, time, expire } = JSON.parse(val)
      if (!expire || expire > new Date().getTime())
        return { value, time }

      this.remove(key)
      return def
    }
    catch (error) {
      this.remove(key)
      return def
    }
  }

  remove(key) {
    this.storage.removeItem(this.getKey(key))
  }

  clear() {
    this.storage.clear()
  }
}

export function createStorage({ prefixKey = '', storage = sessionStorage }) {
  return new Storage({ prefixKey, storage })
}
