<script setup>
import { useRouter } from 'vue-router'
import { NButton, NResult } from 'naive-ui'

import BannerPage from '@/components/page/BannerPage.vue'

const route = useRouter()
</script>

<template>
  <BannerPage label="404" card>
    <NResult status="404" title="404 资源不存在" description="生活总归带点荒谬">
      <template #footer>
        <NButton @click="route.push('/')">
          回到首页
        </NButton>
      </template>
    </NResult>
  </BannerPage>
</template>
