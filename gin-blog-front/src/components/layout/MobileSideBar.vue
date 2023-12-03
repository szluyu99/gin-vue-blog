<script setup>
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { NDrawer, NDrawerContent } from 'naive-ui'
import TheIcon from '@/components/icon/TheIcon.vue'

import { useAppStore, useUserStore } from '@/store'
const { collapsed, blogConfig, articleCount, categoryCount, tagCount } = storeToRefs(useAppStore())

const [route, router] = [useRoute(), useRouter()]
const [userStore, appStore] = [useUserStore(), useAppStore()]

const menuOptions = [
  { text: '首页', icon: 'mdi:home', path: '/' },
  { text: '归档', icon: 'mdi:archive', path: '/archives' },
  { text: '分类', icon: 'mdi:menu', path: '/categories' },
  { text: '标签', icon: 'mdi:tag', path: '/tags' },
  { text: '相册', icon: 'mdi:view-gallery', path: '/albums' },
  { text: '友链', icon: 'mdi:vector-link', path: '/links' },
  { text: '关于', icon: 'mdi:information-outline', path: '/about' },
  { text: '留言', icon: 'mdi:forum', path: '/message' },
]

function logout() {
  userStore.logout()
  if (route.name === 'User') // 如果在个人信息页登出则回到首页
    router.push('/')
  window.$notification?.success({ title: '退出登录成功!', duration: 1500 })
}
</script>

<template>
  <NDrawer v-model:show="collapsed" :width="250" class="opacity-98">
    <NDrawerContent>
      <div class="pt-15 text-center">
        <div class="flex justify-center">
          <img :src="blogConfig.website_avatar" class="h-80 rounded-full" alt="">
        </div>
        <div class="mb-3 mt-6">
          <p class="text-20">
            {{ blogConfig.website_author }}
          </p>
          <p class="mt-3 text-13">
            {{ blogConfig.website_intro }}
          </p>
        </div>
        <!-- 博客信息 -->
        <div class="flex justify-center py-10 text-14">
          <RouterLink to="/archives" class="flex-1" @click="appStore.setCollapsed(false)">
            <p>文章</p>
            <p class="text-15">
              {{ articleCount }}
            </p>
          </RouterLink>
          <RouterLink to="/categories" class="flex-1" @click="appStore.setCollapsed(false)">
            <p>分类</p>
            <p class="text-15">
              {{ categoryCount }}
            </p>
          </RouterLink>
          <RouterLink to="/tags" class="flex-1" @click="appStore.setCollapsed(false)">
            <p>标签</p>
            <p class="text-15">
              {{ tagCount }}
            </p>
          </RouterLink>
        </div>
        <hr class="mb-10 mt-5 border-2px border-color-#d2ebfd border-dashed">
      </div>
      <!-- 菜单 -->
      <div v-for="item of menuOptions" :key="item.text" class="m-8 p-3">
        <RouterLink :to="item.path" class="flex items-center" @click="appStore.setCollapsed(false)">
          <TheIcon :icon="item.icon" :size="20" />
          <span class="ml-25 text-15"> {{ item.text }} </span>
        </RouterLink>
      </div>
      <!-- 登录 -->
      <div>
        <template v-if="!userStore.userId">
          <div class="m-8 flex items-center p-3" @click="appStore.setLoginFlag(true)">
            <TheIcon icon="ph:user-bold" :size="20" />
            <span class="ml-25 text-15"> 登录 </span>
          </div>
        </template>
        <template v-else>
          <RouterLink to="/user">
            <div class="m-8 flex items-center p-5" @click="appStore.setCollapsed(false)">
              <TheIcon icon="mdi:account-circle" :size="20" />
              <span class="ml-25 text-15"> 个人中心 </span>
            </div>
          </RouterLink>
          <div class="m-8 flex items-center p-5" @click="logout">
            <TheIcon icon="mdi:logout" :size="20" />
            <span class="ml-25 text-15"> 退出登录 </span>
          </div>
        </template>
      </div>
    </NDrawerContent>
  </NDrawer>
</template>
