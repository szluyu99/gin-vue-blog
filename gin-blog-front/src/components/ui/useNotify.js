import { createApp } from 'vue'
import UToast from './UToast.vue'

// Create a mount node for the toast component
const mountNode = document.createElement('div')
document.body.appendChild(mountNode)

let notifyInstance = createInstance({
  position: 'top',
  align: 'right',
  closeable: true,
  timeout: 3500,
})

let notifyMounted = notifyInstance.mount(mountNode)

// Cache the toast position and align
// When the props change, unmount the current toast instance
const _notifyCache = {
  position: 'top',
  align: 'right',
  closeable: true,
}

export function useNotify(options = {
  position: 'top',
  align: 'right',
  closeable: true,
}) {
  if (!compareObjects(_notifyCache, options)) {
    _notifyCache.position = options.position
    _notifyCache.align = options.align

    notifyInstance.unmount()
    notifyInstance = createInstance(options)
    notifyMounted = notifyInstance.mount(mountNode)
  }

  return {
    success: msg => notifyMounted.success(msg),
    info: msg => notifyMounted.info(msg),
    warning: msg => notifyMounted.warning(msg),
    error: msg => notifyMounted.error(msg),
  }
}

function createInstance(options) {
  return createApp(UToast, options)
}
function compareObjects(o1, o2) {
  return ['position', 'align', 'closeable'].every(key => o1[key] === o2[key])
}

export const notify = useNotify()
