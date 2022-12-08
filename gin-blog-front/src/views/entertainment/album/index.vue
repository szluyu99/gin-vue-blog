<script setup lang="ts">
const router = useRouter()
</script>

<template>
  <BannerCard label="album">
    <n-result
      status="403"
      title="403 禁止访问"
      description="这里暂时还不能给你看"
    >
      <template #footer>
        <n-button @click="router.push('/')">
          放轻松
        </n-button>
      </template>
    </n-result>
  </BannerCard>
</template>
