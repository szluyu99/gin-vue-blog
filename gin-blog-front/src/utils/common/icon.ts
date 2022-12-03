import { h } from 'vue'
import { Icon } from '@iconify/vue'
import { NIcon } from 'naive-ui'
import SvgIcon from '@/components/icon/SvgIcon.vue'

interface Props {
  size?: number
  color?: string
}

export function renderIcon(icon: string, props: Props = { size: 16 }) {
  return () => h(NIcon, props, { default: () => h(Icon, { icon }) })
}

export function renderCustomIcon(icon: string, props: Props = { size: 12 }) {
  return () => h(NIcon, props, { default: () => h(SvgIcon, { icon }) })
}
