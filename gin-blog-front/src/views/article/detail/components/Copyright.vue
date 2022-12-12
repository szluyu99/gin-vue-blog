<script setup lang="ts">
import { useAppStore } from '@/store'
const { blogConfig } = storeToRefs(useAppStore())
const articleHref = window.location.href
</script>

<template>
  <div
    px-17 py-12 leading-30 text-14
    border-1 border-dashed border-lightblue
    lg:text-15
  >
    <!-- TODO: 点击复制 -->
    <div relative>
      <div
        absolute right-0 top-0 w-16 h-16
        rounded-full bg-white border-blue border-4px
      />
      <span color="#49b1f5" font-bold> 文章作者： </span>
      <router-link to="/" color="#99a9bf" class="!underline">
        {{ blogConfig.website_author }}
      </router-link>
    </div>
    <div>
      <span color="#49b1f5" font-bold> 文章链接： </span>
      <a :href="articleHref" color="#99a9bf" class="!underline">
        {{ articleHref }}
      </a>
    </div>
    <div>
      <span color="#49b1f5" font-bold> 版权声明： </span>
      <span>
        本博客所有文章除特别声明外，均采用
        <a
          href="https://creativecommons.org/licenses/by-nc-sa/4.0/" target="_blank"
          color="#99a9bf" i_underline
        >
          CC BY-NC-SA 4.0
        </a>
        许可协议。转载请注明文章出处。
      </span>
    </div>
  </div>
</template>
