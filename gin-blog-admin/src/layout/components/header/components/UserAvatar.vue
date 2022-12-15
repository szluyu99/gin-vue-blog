<script setup>
import { renderIcon } from '@/utils'
import { useUserStore } from '@/store'

const userStore = useUserStore()
const router = useRouter()

const options = [
  {
    label: '个人中心',
    key: 'profile',
    icon: renderIcon('mdi:account', { size: '14px' }),
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon('mdi:exit-to-app', { size: '14px' }),
  },
]

function handleSelect(key) {
  if (key === 'logout') {
    $dialog.confirm({
      title: '提示',
      type: 'info',
      content: '确认退出？',
      confirm() {
        userStore.logout()
      },
    })
  }
  else if (key === 'profile') {
    router.push('/profile')
  }
}
</script>

<template>
  <n-dropdown :options="options" @select="handleSelect">
    <div flex items-center cursor-pointer>
      <img :src="userStore.avatar" mr-10 w-35 h-35 rounded-full>
      <span>{{ userStore.nickname }}</span>
    </div>
  </n-dropdown>
</template>
