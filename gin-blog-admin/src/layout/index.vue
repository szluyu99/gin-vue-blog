<script setup>
import { storeToRefs } from 'pinia'
import { NLayout, NLayoutSider } from 'naive-ui'

import { header, tags } from '../../setting/theme.json'
import AppHeader from './components/header/index.vue'
import AppMain from './components/AppMain.vue'
import Sidebar from './components/sidebar/index.vue'
import AppTags from './components/tags/index.vue'

import { useThemeStore } from '@/store'

const { collapsed } = storeToRefs(useThemeStore())
</script>

<template>
  <NLayout has-sider class="h-full w-full">
    <!-- 左侧边栏 -->
    <NLayoutSider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      :native-scrollbar="false"
      :collapsed="collapsed"
    >
      <Sidebar />
    </NLayoutSider>
    <!-- 右半部分 -->
    <article class="flex flex-1 flex-col overflow-hidden">
      <!-- 头部 -->
      <header
        :style="`height: ${header.height}px`"
        class="flex items-center border-b-1 border-gray-200 border-b-solid px-15"
      >
        <AppHeader />
      </header>
      <!-- 标签栏 -->
      <section v-if="tags.visible" class="border-b border-gray-200 border-b-solid">
        <AppTags :style="{ height: `${tags.height}px` }" />
      </section>
      <!-- 主体内容 -->
      <section class="flex-1 overflow-hidden">
        <AppMain />
      </section>
    </article>
  </NLayout>
</template>
