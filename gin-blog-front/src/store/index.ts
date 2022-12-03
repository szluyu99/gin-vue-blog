import { createPinia } from 'pinia'
import type { App } from 'vue'

export function setupStore(app: App) {
  app.use(createPinia())
}

export * from './modules'
