<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'

const { blogConfig } = storeToRefs(useAppStore())
</script>

<template>
  <footer class="flex flex-col items-center justify-center text-14 text-gray-700">
    <p>
      Copyright©2023
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
