import { createApp } from 'vue'
import UToast from './UToast.vue'

// Create a mount node for the toast component
const mountNode = document.createElement('div')
document.body.appendChild(mountNode)

let toastInstance = createInstance({
  position: 'top',
  align: 'center',
  closeable: true,
  timeout: 3500,
})

let toastMounted = toastInstance.mount(mountNode)

// Cache the toast position and align
// When the props change, unmount the current toast instance
const _toastCache = {
  position: 'top',
  align: 'center',
  closeable: true,
}

export function useToast(options = {
  position: 'top',
  align: 'center',
  closeable: true,
}) {
  if (!compareObjects(_toastCache, options)) {
    _toastCache.position = options.position
    _toastCache.align = options.align

    toastInstance.unmount()
    toastInstance = createInstance(options)
    toastMounted = toastInstance.mount(mountNode)
  }

  return {
    success: msg => toastMounted.success(msg),
    info: msg => toastMounted.info(msg),
    warning: msg => toastMounted.warning(msg),
    error: msg => toastMounted.error(msg),
  }
}

function createInstance(options) {
  return createApp(UToast, options)
}

function compareObjects(o1, o2) {
  return ['position', 'align', 'closeable'].every(key => o1[key] === o2[key])
}

export const toast = useToast()
