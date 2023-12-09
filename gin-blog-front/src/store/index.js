import { createPinia } from 'pinia'

// https://github.com/prazdevs/pinia-plugin-persistedstate
// pinia 数据持久化，解决刷新数据丢失的问题
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

export const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

export * from './app'
export * from './user'
