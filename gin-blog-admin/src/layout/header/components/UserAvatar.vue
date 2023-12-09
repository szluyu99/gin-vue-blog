<script setup>
import { useRouter } from 'vue-router'
import { NDropdown } from 'naive-ui'

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
  <NDropdown :options="options" @select="handleSelect">
    <div class="flex cursor-pointer items-center">
      <img :src="userStore.avatar" class="mr-10 h-35 w-35 rounded-full">
      <span>{{ userStore.nickname }}</span>
    </div>
  </NDropdown>
</template>
