<script setup>
import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'
import { storeToRefs } from 'pinia'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useAppStore } from '@/store'
import { convertImgUrl } from '@/utils'

/*
侧边信息卡

原来是三张独立卡片(AuthorInfo / Announcement / WebsiteInfo), 而且都写着
`hidden lg:block` —— 窄屏下整列消失, 手机上看不到公告和网站资讯。
合成一张卡之后信息密度够高, 窄屏也放得下, 由首页决定放在文章流下方。
*/
dayjs.extend(duration)

const { blogConfig, blogInfo, articleCount, categoryCount, tagCount, viewCount } = storeToRefs(useAppStore())

const notice = computed(() => blogInfo.value.blog_config?.website_notice)

const counts = computed(() => [
  { label: '文章', value: articleCount.value, path: '/archives' },
  { label: '分类', value: categoryCount.value, path: '/categories' },
  { label: '标签', value: tagCount.value, path: '/tags' },
])

const runTime = ref('-')

// 不能用 duration.format('D 天'): D 是"天"这个分量, 建站超过一个月就会归零重算
// (1340 天会显示成 "1 天"), 总天数要用 asDays
function refreshRunTime() {
  const createTime = dayjs(blogConfig.value.website_createtime)
  if (!createTime.isValid()) { // 配置还没拉回来
    runTime.value = '-'
    return
  }
  const diff = dayjs.duration(dayjs().diff(createTime))
  runTime.value = `${Math.floor(diff.asDays())} 天 ${diff.hours()} 时 ${diff.minutes()} 分`
}

let timer = null

onMounted(() => {
  refreshRunTime()
  timer = setInterval(refreshRunTime, 30 * 1000)
})

// 建站时间来自接口, 拿到之后立刻重算, 不用等下一个 30 秒
watch(() => blogConfig.value.website_createtime, refreshRunTime)

onUnmounted(() => clearInterval(timer))

function addToFavorites() {
  window.$message?.info('按 CTRL + D 将本页加入书签')
}
</script>

<template>
  <div class="card-view card-enter space-y-4">
    <!-- 博主 -->
    <div class="flex items-center gap-3 lg:flex-col lg:text-center">
      <img
        class="h-16 w-16 shrink-0 rounded-full bg-surface-soft object-cover duration-600 lg:h-[105px] lg:w-[105px] hover:rotate-360"
        :src="convertImgUrl(blogConfig.website_avatar)" alt="博主头像"
      >
      <div class="min-w-0">
        <p class="truncate text-xl lg:text-2xl">
          {{ blogConfig.website_author }}
        </p>
        <p class="text-sm color-muted lg:text-base">
          {{ blogConfig.website_intro }}
        </p>
      </div>
    </div>

    <!-- 文章 / 分类 / 标签 -->
    <div class="flex border-y border-color-divider py-2 text-center">
      <RouterLink
        v-for="item of counts" :key="item.label"
        :to="item.path" class="flex-1 transition-300 hover:text-primary"
      >
        <p class="text-sm color-muted">
          {{ item.label }}
        </p>
        <p class="text-xl">
          {{ item.value }}
        </p>
      </RouterLink>
    </div>

    <!-- 公告: 没配就整段不渲染, 留一个空标题没有意义 -->
    <div v-if="notice" class="space-y-1">
      <p class="flex items-center">
        <span class="i-fluent-emoji-flat:bell mr-1.5 inline-block" />
        <span>公告</span>
      </p>
      <p class="text-sm color-muted leading-6">
        {{ notice }}
      </p>
    </div>

    <!-- 网站资讯 -->
    <div class="space-y-1">
      <p class="flex items-center">
        <span class="i-icon-park:analysis mr-1.5 inline-block" />
        <span>网站资讯</span>
      </p>
      <p class="flex justify-between text-sm">
        <span class="color-muted">运行时间</span>
        <span>{{ runTime }}</span>
      </p>
      <p class="flex justify-between text-sm">
        <span class="color-muted">总访问量</span>
        <span>{{ viewCount }}</span>
      </p>
    </div>

    <!-- 社交图标: 窄屏不重复渲染, 封面图上已经有一组 -->
    <div class="hidden items-center justify-center gap-4 text-xl lg:flex">
      <a
        :href="`http://wpa.qq.com/msgrd?v=3&uin=${blogConfig.qq}&site=qq&menu=yes`"
        target="_blank" rel="noopener noreferrer" title="QQ"
      >
        <span class="i-ant-design:qq-circle-filled inline-block transition-300 hover:text-accent" />
      </a>
      <a :href="blogConfig.github" target="_blank" rel="noopener noreferrer" title="GitHub">
        <span class="i-mdi:github inline-block transition-300 hover:text-accent" />
      </a>
      <a :href="blogConfig.gitee" target="_blank" rel="noopener noreferrer" title="Gitee">
        <span class="i-simple-icons:gitee inline-block transition-300 hover:text-accent" />
      </a>
    </div>

    <button
      class="h-9 w-full f-c-c rounded bg-primary text-white transition-200 hover:bg-accent"
      @click="addToFavorites"
    >
      <span class="i-mdi:bookmark mr-1 inline-block text-xl" /> 加入书签
    </button>
  </div>
</template>
