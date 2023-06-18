<script setup>
import { NResult } from 'naive-ui'
import CommonPage from '@/components/page/CommonPage.vue'
</script>

<template>
  <CommonPage>
    <NResult
      class="mt-150"
      status="404"
      title="该页面还在开发中..."
      description="稍安毋躁"
      size="huge"
    />
  </CommonPage>
</template>
