<template>
  <n-icon mr-20 size="18" cursor-pointer @click="handleLinkClick">
    <icon-mdi:github />
  </n-icon>
</template>

<script setup>
function handleLinkClick() {
  window.open('https://github.com/szluyu99/gin-vue-blog')
}
</script>
