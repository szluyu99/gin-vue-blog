<script setup>
import { NIcon } from 'naive-ui'

function handleLinkClick() {
  window.open('https://github.com/szluyu99/gin-vue-blog')
}
</script>

<template>
  <NIcon class="cursor-pointer" size="18" @click="handleLinkClick">
    <span class="i-mdi:github" />
  </NIcon>
</template>
