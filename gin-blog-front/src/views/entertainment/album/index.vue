<script setup>
import { useRouter } from 'vue-router'
import { NButton, NResult } from 'naive-ui'
import BannerPage from '@/components/page/BannerPage.vue'
const router = useRouter()
</script>

<template>
  <BannerPage label="album" card>
    <NResult
      status="403"
      title="403 禁止访问"
      description="这里暂时还不能给你看"
    >
      <template #footer>
        <NButton @click="router.push('/')">
          放轻松
        </NButton>
      </template>
    </NResult>
  </BannerPage>
</template>
