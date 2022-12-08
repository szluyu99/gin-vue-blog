import { createPinia } from 'pinia'
import type { App } from 'vue'

// TODO: pinia 持久化
export function setupStore(app: App) {
  app.use(createPinia())
}

export * from './modules'
