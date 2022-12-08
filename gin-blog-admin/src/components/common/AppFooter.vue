<script setup>
import { useAppStore } from '@/store'
const { blogConfig } = storeToRefs(useAppStore())
</script>

<template>
  <footer text-14 f-c-c flex-col color="#6a6a6a">
    <p>
      Copyright©2022
      <a
        href="https://github.com/szluyu99"
        target="__blank"
        hover="decoration-underline color-primary"
      >
        szluyu99
      </a>
    </p>
    <p>
      <a
        href="http://beian.miit.gov.cn/"
        target="__blank"
        hover="decoration-underline color-primary"
      >
        {{ blogConfig.website_record ?? '备案信息未知' }}
      </a>
    </p>
  </footer>
</template>
