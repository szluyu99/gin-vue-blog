<template>
  <div flex items-center>
    <!-- 侧边栏折叠按钮 -->
    <MenuCollapse />
    <!-- 面包屑导航 -->
    <BreadCrumb ml-15 />
  </div>
  <div flex ml-auto items-center>
    <!-- Github 网址图标 -->
    <GithubSite />
    <!-- 全屏按钮 -->
    <FullScreen />
    <!-- 用户头像 -->
    <UserAvatar />
  </div>
</template>

<script setup>
import MenuCollapse from './components/MenuCollapse.vue'
import BreadCrumb from './components/BreadCrumb.vue'
import GithubSite from './components/GithubSite.vue'
import FullScreen from './components/FullScreen.vue'
import UserAvatar from './components/UserAvatar.vue'
</script>
