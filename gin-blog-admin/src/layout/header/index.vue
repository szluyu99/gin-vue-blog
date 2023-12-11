<script setup>
import MenuCollapse from './components/MenuCollapse.vue'
import BreadCrumb from './components/BreadCrumb.vue'
import ThemeMode from './components/ThemeMode.vue'
import Watermark from './components/Watermark.vue'
import GithubSite from './components/GithubSite.vue'
import FullScreen from './components/FullScreen.vue'
import UserAvatar from './components/UserAvatar.vue'
</script>

<template>
  <div class="flex items-center space-x-3">
    <!-- 侧边栏折叠按钮 -->
    <MenuCollapse />
    <!-- 面包屑导航 -->
    <BreadCrumb class="mt-1.5" />
  </div>
  <div class="ml-auto flex items-center space-x-5">
    <!-- Github 网址图标 -->
    <GithubSite />
    <!-- 切换模式 -->
    <ThemeMode />
    <!-- 水印 -->
    <Watermark />
    <!-- 全屏按钮 -->
    <FullScreen />
    <!-- 用户头像 -->
    <UserAvatar />
  </div>
</template>
