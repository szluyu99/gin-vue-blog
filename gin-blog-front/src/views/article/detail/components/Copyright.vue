<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/store'
const { blogConfig } = storeToRefs(useAppStore())
const articleHref = window.location.href
</script>

<template>
  <div class="border-1 border-lightblue border-dashed px-4 py-3 leading-7">
    <!-- TODO: 点击复制 -->
    <div class="relative">
      <div class="absolute right-0 top-0 h-5 w-5 border-4px border-blue rounded-full bg-white" />
      <span class="font-bold color-#49b1f5"> 文章作者： </span>
      <RouterLink to="/" class="color-#99a9bf !underline">
        {{ blogConfig.website_author }}
      </RouterLink>
    </div>
    <div>
      <span class="font-bold color-#49b1f5"> 文章链接： </span>
      <a :href="articleHref" class="color-#99a9bf !underline">
        {{ articleHref }}
      </a>
    </div>
    <div>
      <span class="font-bold color-#49b1f5"> 版权声明： </span>
      <span>
        本博客所有文章除特别声明外，均采用
        <a
          href="https://creativecommons.org/licenses/by-nc-sa/4.0/" target="_blank"
          class="color-#99a9bf !underline"
        >
          CC BY-NC-SA 4.0
        </a>
        许可协议。转载请注明文章出处。
      </span>
    </div>
  </div>
</template>
