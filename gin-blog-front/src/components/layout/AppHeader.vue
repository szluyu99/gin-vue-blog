<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { Icon } from '@iconify/vue'

import MobileSideBar from './MobileSideBar.vue'
import { convertImgUrl } from '@/utils'
import { useAppStore, useUserStore } from '@/store'
import api from '@/api'

const appStore = useAppStore()
const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const menuOptions = [
  { text: '首页', icon: 'mdi:home', path: '/' },
  {
    text: '发现',
    icon: 'mdi:apple-safari',
    subMenu: [
      { text: '归档', icon: 'mdi:archive', path: '/archives' },
      { text: '分类', icon: 'mdi:menu', path: '/categories' },
      { text: '标签', icon: 'mdi:tag', path: '/tags' },
    ],
  },
  {
    text: '娱乐',
    icon: 'mdi:gamepad-circle',
    subMenu: [
      { text: '相册', icon: 'mdi:view-gallery', path: '/albums' },
    ],
  },
  { text: '友链', icon: 'mdi:vector-link', path: '/links' },
  { text: '关于', icon: 'mdi:information-outline', path: '/about' },
  { text: '留言', icon: 'mdi:forum', path: '/message' },
]

const navClass = ref('nav')
const barShow = ref(true)

// * 监听 y 效果比添加 scroll 监听器效果更好
// * 节流操作, 效果很好
const { y } = useWindowScroll()
const preY = ref(0) // 记录上一次的 y 滚动距离
watchThrottled(y, () => {
  if (Math.abs(preY.value - y.value) >= 50) { // 小幅度滚动不进行操作
    barShow.value = (y.value < preY.value)
    navClass.value = (y.value > 60) ? 'nav-fixed' : 'nav'
    preY.value = y.value
  }
}, { throttle: 100 })

async function logout() {
  await userStore.logout()
  if (route.name === 'User') {
    router.push('/')
  }
  window.$notify.success('退出登录成功!')
}

// const blogTitle = import.meta.env.VITE_APP_TITLE
</script>

<template>
  <!-- 移动端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div v-if="barShow" :class="navClass" class="fixed inset-x-0 top-0 z-11 h-[60px] flex items-center justify-between px-4 py-2 lg:hidden">
      <!-- 左上角标题 -->
      <RouterLink to="/" class="text-[18px] font-bold">
        {{ appStore.blogConfig.website_author }}
      </RouterLink>
      <!-- 右上角图标 -->
      <div class="flex items-center gap-2 text-2xl">
        <button @click="appStore.setSearchFlag(true)">
          <Icon icon="ic:round-search" />
        </button>
        <button @click="appStore.setCollapsed(true)">
          <Icon icon="ic:sharp-menu" />
        </button>
      </div>
    </div>
  </Transition>
  <!-- 侧边栏 -->
  <MobileSideBar />
  <!-- 电脑端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div v-if="barShow" :class="navClass" class="fixed inset-x-0 top-0 z-11 hidden h-[60px] lg:block">
      <div class="h-full flex items-center justify-between px-9">
        <!-- 左上角标题 -->
        <RouterLink to="/" class="text-xl font-bold">
          {{ appStore.blogConfig.website_author }}
        </RouterLink>
        <!-- 右上角菜单 -->
        <div class="flex items-center space-x-4">
          <!-- 搜索 -->
          <div class="menus-item">
            <a class="menu-btn flex items-center" @click="appStore.setSearchFlag(true)">
              <Icon icon="mdi:magnify" class="text-xl" />
              <span class="ml-1"> 搜索 </span>
            </a>
          </div>
          <div v-for="item of menuOptions" :key="item.text" class="menus-item">
            <!-- 不包含子菜单 -->
            <RouterLink v-if="!item.subMenu" :to="item.path" class="menu-btn flex items-center">
              <Icon :icon="item.icon" class="text-xl" />
              <span class="ml-1"> {{ item.text }} </span>
            </RouterLink>
            <!-- 包含子菜单 -->
            <div v-else class="menu-btn">
              <div class="flex items-center">
                <Icon :icon="item.icon" class="text-xl" />
                <span class="mx-1"> {{ item.text }} </span>
                <Icon icon="ep:arrow-down-bold" class="text-xl" />
              </div>
              <ul class="menus-submenu">
                <RouterLink v-for="sub of item.subMenu" :key="sub.text" :to="sub.path">
                  <div class="flex items-center">
                    <Icon :icon="sub.icon" class="text-xl" />
                    <span class="ml-1"> {{ sub.text }} </span>
                  </div>
                </RouterLink>
              </ul>
            </div>
          </div>
          <!-- 登录 -->
          <div class="menus-item">
            <a v-if="!userStore.userId" class="menu-btn" @click="appStore.setLoginFlag(true)">
              <div class="flex items-center">
                <Icon icon="ph:user-bold" class="text-xl" />
                <span class="ml-1"> 登录 </span>
              </div>
            </a>
            <template v-else>
              <img :src="convertImgUrl(userStore.avatar)" class="w-8 cursor-pointer rounded-full">
              <ul class="menus-submenu">
                <RouterLink to="/user">
                  <div class="flex items-center">
                    <Icon icon="mdi:account-circle" class="mr-1 text-xl" /> 个人中心
                  </div>
                </RouterLink>
                <a @click="logout">
                  <div class="flex items-center">
                    <Icon icon="mdi:logout" class="mr-1 text-xl" /> 退出登录
                  </div>
                </a>
              </ul>
            </template>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped lang="scss">
.nav {
  transition: all 0.8s;
  color: #fff;
  background: rgba(0, 0, 0, 0) !important;
}

.nav-fixed {
  transition: all 0.8s;
  color: #000;
  background: rgba(255, 255, 255, 0.8) !important;
  box-shadow: 0 5px 6px -5px rgba(133, 133, 133, 0.6);
  & .menu-btn:hover {
    color: #49b1f5 !important;
  }
}

.menus-item {
  position: relative;
  display: inline-block;
  // margin: 0 0 0 1rem;
  a {
    transition: all 0.2s;
  }
  a::after {
    position: absolute;
    bottom: -5px;
    left: 0;
    z-index: -1;
    width: 0;
    height: 3px;
    background-color: #80c8f8;
    content: "";
    transition: all 0.3s ease-in-out;
  }
  .menu-btn {
    cursor: pointer;
    &:hover::after {
      width: 100%;
    }
  }
}

.menus-item:hover .menus-submenu {
  display: block;
}

.menus-submenu {
  position: absolute;
  display: none;
  right: 0;
  width: max-content;
  margin-top: 8px;
  box-shadow: 0 5px 20px -4px rgba(0, 0, 0, 0.5);
  background-color: #fff;
  animation: submenu 0.3s 0.1s ease both;

  &::before {
    position: absolute;
    top: -8px;
    left: 0;
    width: 100%;
    height: 20px;
    content: "";
  }
  a {
    line-height: 2;
    color: #4c4948 !important;
    text-shadow: none;
    display: block;
    padding: 6px 14px;
  }
  a:hover {
    background: #4ab1f4;
  }
}

@keyframes submenu {
  0% {
    opacity: 0;
    filter: alpha(opacity=0);
    transform: translateY(10px);
  }

  100% {
    opacity: 1;
    filter: none;
    transform: translateY(0);
  }
}
</style>
