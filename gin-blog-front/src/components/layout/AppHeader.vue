<script setup>
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAppStore, useNotificationStore, useUserStore } from '@/store'
import { convertImgUrl } from '@/utils'
import MobileSideBar from './MobileSideBar.vue'

const appStore = useAppStore()
const userStore = useUserStore()
const notificationStore = useNotificationStore()
const router = useRouter()
const route = useRoute()

// 未读数只在登录后才有意义; 退出登录要连列表一起清掉, 否则下一个人登录
// 会先看到上一个人的通知
onMounted(() => {
  if (userStore.userId) {
    notificationStore.fetchUnreadCount()
  }
})

watch(() => userStore.userId, (id) => {
  if (id) {
    notificationStore.fetchUnreadCount()
  }
  else {
    notificationStore.reset()
  }
})

// 列表按需拉: 鼠标移上铃铛才请求, 不做轮询
function onBellEnter() {
  notificationStore.fetchList()
}

// 点一条通知: 先标已读再跳到对应文章
async function openNotification(item) {
  if (!item.is_read) {
    await notificationStore.read([item.id])
  }
  if (item.article_id) {
    router.push(`/article/${item.article_id}`)
  }
}

const menuOptions = [
  { text: '首页', icon: 'i-mdi:home', path: '/' },
  {
    text: '发现',
    icon: 'i-mdi:apple-safari',
    subMenu: [
      { text: '归档', icon: 'i-mdi:archive', path: '/archives' },
      { text: '分类', icon: 'i-mdi:menu', path: '/categories' },
      { text: '标签', icon: 'i-mdi:tag', path: '/tags' },
    ],
  },
  {
    text: '娱乐',
    icon: 'i-mdi:gamepad-circle',
    subMenu: [
      { text: '相册', icon: 'i-mdi:view-gallery', path: '/albums' },
    ],
  },
  { text: '友链', icon: 'i-mdi:vector-link', path: '/links' },
  { text: '关于', icon: 'i-mdi:information-outline', path: '/about' },
  { text: '留言', icon: 'i-mdi:forum', path: '/message' },
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
  window.$notify?.success('退出登录成功!')
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
        <button :title="appStore.isDark ? '切换浅色模式' : '切换深色模式'" @click="appStore.toggleTheme()">
          <span :class="appStore.isDark ? 'i-mdi:weather-sunny' : 'i-mdi:weather-night'" />
        </button>
        <button @click="appStore.setSearchFlag(true)">
          <span class="i-ic:round-search" />
        </button>
        <button @click="appStore.setCollapsed(true)">
          <span class="i-ic:sharp-menu" />
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
              <span class="i-mdi:magnify text-xl" />
              <span class="ml-1"> 搜索 </span>
            </a>
          </div>
          <div v-for="item of menuOptions" :key="item.text" class="menus-item">
            <!-- 不包含子菜单 -->
            <RouterLink v-if="!item.subMenu" :to="item.path" class="menu-btn flex items-center">
              <span :class="item.icon" class="text-xl" />
              <span class="ml-1"> {{ item.text }} </span>
            </RouterLink>
            <!-- 包含子菜单 -->
            <div v-else class="menu-btn">
              <div class="flex items-center">
                <span :class="item.icon" class="text-xl" />
                <span class="mx-1"> {{ item.text }} </span>
                <span class="i-ep:arrow-down-bold text-xl" />
              </div>
              <ul class="menus-submenu">
                <RouterLink v-for="sub of item.subMenu" :key="sub.text" :to="sub.path">
                  <div class="flex items-center">
                    <span :class="sub.icon" class="text-xl" />
                    <span class="ml-1"> {{ sub.text }} </span>
                  </div>
                </RouterLink>
              </ul>
            </div>
          </div>
          <!-- 主题切换 -->
          <div class="menus-item">
            <a class="menu-btn flex items-center" :title="appStore.isDark ? '切换浅色模式' : '切换深色模式'" @click="appStore.toggleTheme()">
              <span :class="appStore.isDark ? 'i-mdi:weather-sunny' : 'i-mdi:weather-night'" class="text-xl" />
            </a>
          </div>
          <!-- 站内通知: 未登录没有通知可看, 整块不渲染 -->
          <div v-if="userStore.userId" class="menus-item" @mouseenter="onBellEnter">
            <a class="menu-btn relative flex items-center" title="站内通知">
              <span class="i-mdi:bell-outline text-xl" />
              <!-- 未读红点: 数字超过 99 就显示 99+, 否则会把导航栏撑开 -->
              <span
                v-if="notificationStore.unreadCount"
                class="absolute h-4 min-w-4 flex items-center justify-center rounded-full bg-accent px-1 text-[10px] text-white -right-2 -top-1"
              >
                {{ notificationStore.unreadCount > 99 ? '99+' : notificationStore.unreadCount }}
              </span>
            </a>
            <ul class="menus-submenu w-[320px] text-left">
              <li class="flex items-center justify-between border-b border-color-divider px-3 py-2">
                <span class="text-sm font-bold">站内通知</span>
                <span
                  v-if="notificationStore.unreadCount"
                  class="cursor-pointer text-xs text-primary" @click="notificationStore.read()"
                >
                  全部已读
                </span>
              </li>
              <li v-if="notificationStore.loading" class="px-3 py-4 text-center text-sm color-muted">
                加载中...
              </li>
              <li v-else-if="!notificationStore.list.length" class="px-3 py-6 text-center text-sm color-muted">
                还没有通知
              </li>
              <template v-else>
                <li
                  v-for="item of notificationStore.list" :key="item.id"
                  class="flex cursor-pointer gap-2 px-3 py-2 transition-300 hover:bg-surface-soft"
                  :class="item.is_read ? 'op-60' : ''"
                  @click="openNotification(item)"
                >
                  <img
                    :src="convertImgUrl(item.from_avatar)" :alt="item.from_nickname"
                    class="h-8 w-8 shrink-0 rounded-full bg-surface-soft object-cover"
                  >
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm">
                      <span class="font-bold">{{ item.from_nickname || '有人' }}</span>
                      {{ item.type === 1 ? '回复了你' : '评论了你的文章' }}
                    </p>
                    <p class="truncate text-xs color-muted">
                      {{ item.content }}
                    </p>
                    <p class="mt-0.5 text-xs color-muted">
                      {{ item.article_title }}
                    </p>
                  </div>
                  <span v-if="!item.is_read" class="mt-1 h-2 w-2 shrink-0 rounded-full bg-accent" />
                </li>
              </template>
              <!-- 下拉只放最近 10 条, 更多的去通知页翻 -->
              <li class="border-t border-color-divider">
                <RouterLink to="/notifications" class="block px-3 py-2 text-center text-xs text-primary">
                  查看全部
                </RouterLink>
              </li>
            </ul>
          </div>
          <!-- 登录 -->
          <div class="menus-item">
            <a v-if="!userStore.userId" class="menu-btn" @click="appStore.setLoginFlag(true)">
              <div class="flex items-center">
                <span class="i-mdi:account text-xl" />
                <span class="ml-1"> 登录 </span>
              </div>
            </a>
            <template v-else>
              <img :src="convertImgUrl(userStore.avatar)" class="w-8 cursor-pointer rounded-full">
              <ul class="menus-submenu">
                <RouterLink to="/user">
                  <div class="flex items-center">
                    <span class="i-mdi:account-circle mr-1 text-xl" /> 个人中心
                  </div>
                </RouterLink>
                <a @click="logout">
                  <div class="flex items-center">
                    <span class="i-mdi:logout mr-1 text-xl" /> 退出登录
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

/* 滚动后导航栏是实底的, 深色模式下要跟着变暗, 否则一片白 */
html.dark .nav-fixed {
  color: var(--c-text);
  background: rgba(24, 26, 31, 0.85) !important;
  box-shadow: 0 5px 6px -5px rgba(0, 0, 0, 0.8);
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
  background-color: var(--c-surface);
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
    color: var(--c-text) !important;
    text-shadow: none;
    display: block;
    padding: 6px 14px;
  }
  a:hover {
    background: #4ab1f4;
    color: #fff !important;
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
