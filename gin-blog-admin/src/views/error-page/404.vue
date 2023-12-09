<script setup>
import { NButton, NResult } from 'naive-ui'
import AppPage from '@/components/common/AppPage.vue'
</script>

<template>
  <AppPage>
    <NResult m-auto status="404" description="抱歉，您访问的页面不存在。">
      <template #icon>
        <img src="/image/404.webp" width="500">
      </template>
      <template #footer>
        <NButton @click="$router.replace('/')">
          返回首页
        </NButton>
      </template>
    </NResult>
  </AppPage>
</template>
