// custom style
import './styles/index.css'
import './styles/common.css'
import './styles/animate.css'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

// vue
import { createApp } from 'vue'

// hljs
import hljs from 'highlight.js/lib/core'
import 'highlight.js/styles/a11y-dark.css'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import bash from 'highlight.js/lib/languages/bash'

import { router } from './router'
import { pinia } from './store'
import App from './App.vue'

// 全局 toast
import { useToast } from './components/ui/useToast'
import { useNotify } from './components/ui/useNotify'
window.$message = useNotify({ position: 'top', align: 'center', timeout: 3000, closeable: true })
window.$notify = useToast({ position: 'top', align: 'right', timeout: 3000, closeable: true })

hljs.registerLanguage('go', go)
hljs.registerLanguage('java', java)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)

const app = createApp(App)
app.use(router).use(pinia)
app.mount('#app')
