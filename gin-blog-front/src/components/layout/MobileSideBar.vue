<script setup>
import { storeToRefs } from 'pinia'

import { useRoute, useRouter } from 'vue-router'

import UDrawer from '@/components/ui/UDrawer.vue'
import { useAppStore, useUserStore } from '@/store'

const { collapsed, blogConfig, articleCount, categoryCount, tagCount } = storeToRefs(useAppStore())

const [route, router] = [useRoute(), useRouter()]
const [userStore, appStore] = [useUserStore(), useAppStore()]

const menuOptions = [
  { text: '首页', icon: 'i-mdi:home', path: '/' },
  { text: '归档', icon: 'i-mdi:archive', path: '/archives' },
  { text: '分类', icon: 'i-mdi:menu', path: '/categories' },
  { text: '标签', icon: 'i-mdi:tag', path: '/tags' },
  { text: '相册', icon: 'i-mdi:view-gallery', path: '/albums' },
  { text: '友链', icon: 'i-mdi:vector-link', path: '/links' },
  { text: '关于', icon: 'i-mdi:information-outline', path: '/about' },
  { text: '留言', icon: 'i-mdi:forum', path: '/message' },
]

async function logout() {
  await userStore.logout()
  if (route.name === 'User') { // 如果在个人信息页登出则回到首页
    router.push('/')
  }
  window.$notify?.success('退出登录成功！')
}
</script>

<template>
  <UDrawer v-model="collapsed" placement="right" :width="250">
    <div class="mx-5">
      <div class="pt-4 text-center space-y-3">
        <div class="flex justify-center">
          <img :src="blogConfig.website_avatar" class="h-20 rounded-full" alt="作者头像">
        </div>
        <!-- 头像和介绍 -->
        <div class="space-y-1">
          <p class="text-lg">
            {{ blogConfig.website_author }}
          </p>
          <p class="text-sm">
            {{ blogConfig.website_intro }}
          </p>
        </div>
        <!-- 博客信息 -->
        <div class="flex justify-center text-sm">
          <RouterLink to="/archives" class="flex-1" @click="appStore.setCollapsed(false)">
            <p> 文章 </p>
            <p> {{ articleCount }} </p>
          </RouterLink>
          <RouterLink to="/categories" class="flex-1" @click="appStore.setCollapsed(false)">
            <p> 分类 </p>
            <p> {{ categoryCount }} </p>
          </RouterLink>
          <RouterLink to="/tags" class="flex-1" @click="appStore.setCollapsed(false)">
            <p> 标签 </p>
            <p> {{ tagCount }} </p>
          </RouterLink>
        </div>
      </div>
      <!-- 分隔线 -->
      <hr class="my-4 border-2 border-color-divider border-dashed">
      <!-- 菜单 -->
      <div v-for="item of menuOptions" :key="item.text" class="m-2 p-1">
        <RouterLink :to="item.path" class="flex items-center" @click="appStore.setCollapsed(false)">
          <span :class="item.icon" class="text-lg" />
          <span class="ml-5"> {{ item.text }} </span>
        </RouterLink>
      </div>
      <!-- 登录 -->
      <div>
        <template v-if="!userStore.userId">
          <div class="m-2 flex items-center p-1" @click="appStore.setLoginFlag(true)">
            <span class="i-mdi:account text-lg" />
            <span class="ml-5"> 登录 </span>
          </div>
        </template>
        <template v-else>
          <RouterLink to="/user">
            <div class="m-2 flex items-center p-1" @click="appStore.setCollapsed(false)">
              <span class="i-mdi:account-circle text-lg" />
              <span class="ml-5"> 个人中心 </span>
            </div>
          </RouterLink>
          <div class="m-2 flex items-center p-1" @click="logout">
            <span class="i-mdi:logout text-lg" />
            <span class="ml-5"> 退出登录 </span>
          </div>
        </template>
      </div>
    </div>
  </UDrawer>
</template>
