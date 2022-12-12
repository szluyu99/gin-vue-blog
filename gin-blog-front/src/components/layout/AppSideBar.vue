<script setup lang="ts">
import { useAppStore, useUserStore } from '@/store'
const { collapsed, blogConfig, articleCount, categoryCount, tagCount } = storeToRefs(useAppStore())

const [route, router] = [useRoute(), useRouter()]
const [userStore, appStore] = [useUserStore(), useAppStore()]

const menuOptions = [
  {
    text: '首页',
    icon: 'mdi:home',
    path: '/',
  },

  {
    text: '归档',
    icon: 'mdi:archive',
    path: '/archives',
  },
  {
    text: '分类',
    icon: 'mdi:menu',
    path: '/categories',
  },
  {
    text: '标签',
    icon: 'mdi:tag',
    path: '/tags',
  },

  {
    text: '相册',
    icon: 'mdi:view-gallery',
    path: '/albums',
  },
  {
    text: '友链',
    icon: 'mdi:vector-link',
    path: '/links',
  },
  {
    text: '关于',
    icon: 'mdi:information-outline',
    path: '/about',
  },
  {
    text: '留言',
    icon: 'mdi:forum',
    path: '/message',
  },
]

function logout() {
  userStore.logout()
  if (route.name === 'User') // 如果在个人信息页登出则回到首页
    router.push('/')
  window.$notification?.success({ title: '退出登录成功!', duration: 1500 })
}
</script>

<template>
  <n-drawer
    v-model:show="collapsed"
    :width="250"
    opacity-98
  >
    <n-drawer-content>
      <div text-center pt-15>
        <n-avatar :size="80" round :src="blogConfig.website_avatar" />
        <p text-20>
          {{ blogConfig.website_author }}
        </p>
        <p text-13 my-3>
          {{ blogConfig.website_intro }}
        </p>
        <!-- 博客信息 -->
        <div flex justify-center py-10 text-14>
          <router-link
            to="/archives" flex-1
            @click="appStore.setCollapsed(false)"
          >
            <p>文章</p>
            <p text-15>
              {{ articleCount }}
            </p>
          </router-link>
          <router-link
            to="/categories" flex-1
            @click="appStore.setCollapsed(false)"
          >
            <p>分类</p>
            <p text-15>
              {{ categoryCount }}
            </p>
          </router-link>
          <router-link
            to="/tags" flex-1
            @click="appStore.setCollapsed(false)"
          >
            <p>标签</p>
            <p text-15>
              {{ tagCount }}
            </p>
          </router-link>
        </div>
        <hr mt-5 mb-10 border-dashed border-2px border-color="#d2ebfd">
      </div>
      <!-- 菜单 -->
      <div v-for="item of menuOptions" :key="item.text" p-3 m-8>
        <router-link
          :to="item.path"
          flex items-center @click="appStore.setCollapsed(false)"
        >
          <TheIcon :icon="item.icon" :size="20" />
          <span ml-25 text-15> {{ item.text }} </span>
        </router-link>
      </div>
      <!-- 登录 -->
      <div>
        <template v-if="!userStore.userId">
          <div
            flex items-center p-3 m-8
            @click="appStore.setLoginFlag(true)"
          >
            <TheIcon icon="ph:user-bold" :size="20" />
            <span ml-25 text-15> 登录 </span>
          </div>
        </template>
        <template v-else>
          <router-link to="/user">
            <div
              flex items-center p-5 m-8
              @click="appStore.setCollapsed(false)"
            >
              <TheIcon icon="mdi:account-circle" :size="20" />
              <span ml-25 text-15> 个人中心 </span>
            </div>
          </router-link>
          <div
            flex items-center p-5 m-8
            @click="logout"
          >
            <TheIcon icon="mdi:logout" :size="20" />
            <span ml-25 text-15> 退出登录 </span>
          </div>
        </template>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>
