import { h } from 'vue'
import { Icon } from '@iconify/vue'
import { NIcon } from 'naive-ui'
import SvgIcon from '@/components/icon/SvgIcon.vue'

/**
 * @param {String} icon
 * @param {{size: Number, color: String}} props
 * @returns
 */
export function renderIcon(icon, props = { size: 16 }) {
  return () => h(NIcon, props, { default: () => h(Icon, { icon }) })
}

/**
 * @param {String} icon
 * @param {{size: Number, color: String}} props
 * @returns
 */
export function renderCustomIcon(icon, props = { size: 11 }) {
  return () => h(NIcon, props, { default: () => h(SvgIcon, { icon }) })
}
