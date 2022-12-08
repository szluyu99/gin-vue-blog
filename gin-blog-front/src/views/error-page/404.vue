<script setup lang="ts">
const route = useRouter()
</script>

<template>
  <BannerCard label="404">
    <n-result status="404" title="404 资源不存在" description="生活总归带点荒谬">
      <template #footer>
        <n-button @click="route.push('/')">
          回到首页
        </n-button>
      </template>
    </n-result>
  </BannerCard>
</template>
