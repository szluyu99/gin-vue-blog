<template>
  <n-layout has-sider wh-full>
    <!-- 左侧边栏 -->
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      :native-scrollbar="false"
      :collapsed="appStore.collapsed"
    >
      <Sidebar />
    </n-layout-sider>
    <!-- 右半部分 -->
    <article flex-1 flex-col overflow-hidden>
      <!-- 头部 -->
      <header
        bg-white
        px-15
        border-b
        bc-eee
        flex
        items-center
        :style="`height: ${header.height}px`"
      >
        <AppHeader />
      </header>
      <!-- 标签栏 -->
      <section v-if="tags.visible" border-b bc-eee>
        <AppTags :style="{ height: `${tags.height}px` }" />
      </section>
      <!-- 主体内容 -->
      <section flex-1 overflow-hidden>
        <AppMain />
      </section>
    </article>
  </n-layout>
</template>

<script setup>
import AppHeader from './components/header/index.vue'
import AppMain from './components/AppMain.vue'
import Sidebar from './components/sidebar/index.vue'
import AppTags from './components/tags/index.vue'
import { useAppStore } from '@/store'
import { header, tags } from '~/setting'

const appStore = useAppStore()
</script>
