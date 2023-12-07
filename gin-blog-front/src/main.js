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
import json from 'highlight.js/lib/languages/json'
import javascript from 'highlight.js/lib/languages/javascript'
import bash from 'highlight.js/lib/languages/bash'

import { router } from './router'
import { pinia } from './store'
import App from './App.vue'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)

const app = createApp(App)
app.use(router)
app.use(pinia)
app.mount('#app')
