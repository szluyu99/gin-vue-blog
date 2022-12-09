import type { App } from 'vue'
// https://github.com/prazdevs/pinia-plugin-persistedstate
// pinia 数据持久化，解决刷新数据丢失的问题
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

export function setupStore(app: App) {
  const pinia = createPinia()
  pinia.use(piniaPluginPersistedstate)
  app.use(pinia)
}

export * from './modules'
