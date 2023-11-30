<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { NAvatar } from 'naive-ui'

import MobileSideBar from './MobileSideBar.vue'
import TheIcon from '@/components/icon/TheIcon.vue'
import { convertImgUrl } from '@/utils'
import { useAppStore, useUserStore } from '@/store'

const appStore = useAppStore()
const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

// 菜单项 TODO: 直接从路由中加载? 似乎不可行
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
const { y } = useWindowScroll() // 通过 $() 解构 ref
const preY = ref(0) // 记录上一次的 y 滚动距离
watchThrottled(y, () => {
  if (Math.abs(preY.value - y.value) >= 50) { // 小幅度滚动不进行操作
    barShow.value = (y.value < preY.value)
    navClass.value = (y.value > 60) ? 'nav-fixed' : 'nav'
    preY.value = y.value
  }
}, { throttle: 100 })

function logout() {
  userStore.logout()
  if (route.name === 'User') // 如果在个人信息页登出则回到首页
    router.push('/')
  window.$notification?.success({ title: '退出登录成功!', duration: 1500 })
}

// const blogTitle = import.meta.env.VITE_APP_TITLE
</script>

<template>
  <!-- 移动端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div
      v-if="barShow" :class="navClass"
      class="fixed inset-x-0 top-0 z-11 h-60 flex items-center justify-between px-16 py-10 lg:hidden"
    >
      <!-- 左上角标题 -->
      <RouterLink to="/" class="text-18 font-bold">
        {{ appStore.blogConfig.website_author }}
      </RouterLink>
      <!-- 右上角图标 -->
      <div class="flex items-center">
        <button class="mr-10" @click="appStore.setSearchFlag(true)">
          <TheIcon icon="ic:round-search" :size="22" />
        </button>
        <button @click="appStore.setCollapsed(true)">
          <TheIcon icon="ic:sharp-menu" :size="22" />
        </button>
      </div>
    </div>
  </Transition>
  <!-- 侧边栏 -->
  <MobileSideBar />
  <!-- 电脑端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div
      v-if="barShow"
      :class="navClass"
      class="fixed inset-x-0 top-0 z-11 hidden h-60 lg:block"
    >
      <div class="h-full flex items-center justify-between px-36">
        <!-- 左上角标题 -->
        <RouterLink to="/" class="text-18 font-bold">
          {{ appStore.blogConfig.website_author }}
        </RouterLink>
        <!-- 右上角菜单 -->
        <div class="flex items-center text-6xl space-x-18">
          <!-- 搜索 -->
          <div class="menus-item">
            <a class="menu-btn flex items-center" @click="appStore.setSearchFlag(true)">
              <TheIcon icon="mdi:magnify" :size="18" />
              <span class="ml-4"> 搜索 </span>
            </a>
          </div>
          <div v-for="item of menuOptions" :key="item.text" class="menus-item">
            <!-- 不包含子菜单 -->
            <RouterLink v-if="!item.subMenu" :to="item.path" class="menu-btn flex items-center">
              <TheIcon :icon="item.icon" :size="18" />
              <span class="ml-4"> {{ item.text }} </span>
            </RouterLink>
            <!-- 包含子菜单 -->
            <div v-else class="menu-btn">
              <div class="flex items-center">
                <TheIcon :icon="item.icon" :size="18" />
                <span class="mx-4"> {{ item.text }} </span>
                <TheIcon icon="ep:arrow-down-bold" :size="18" />
              </div>
              <ul class="menus-submenu">
                <RouterLink v-for="sub of item.subMenu" :key="sub.text" :to="sub.path">
                  <div class="flex items-center">
                    <TheIcon :icon="sub.icon" :size="18" />
                    <span class="ml-4"> {{ sub.text }} </span>
                  </div>
                </RouterLink>
              </ul>
            </div>
          </div>
          <!-- 登录 -->
          <div class="menus-item">
            <a v-if="!userStore.userId" class="menu-btn" @click="appStore.setLoginFlag(true)">
              <div class="flex items-center">
                <span class="i-ph:user-bold mr-4 text-18" /> 登录
              </div>
            </a>
            <template v-else>
              <NAvatar
                :size="28"
                :src="convertImgUrl(userStore.avatar)"
                class="round"
                fallback-src="https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png"
              />
              <ul class="menus-submenu">
                <RouterLink to="/user">
                  <div class="flex items-center">
                    <span class="i-mdi:account-circle mr-4 text-18" /> 个人中心
                  </div>
                </RouterLink>
                <a @click="logout">
                  <span class="i-mdi:logout text-18" /> 退出登录
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
