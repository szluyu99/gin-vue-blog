<script setup>
import {
  NAvatar,
  NCard,
  NEmpty,
  NGi,
  NGradientText,
  NGrid,
  NProgress,
  NSkeleton,
  NStatistic,
  NTag,
} from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import api from '@/api'
import AppPage from '@/components/common/AppPage.vue'
import { useUserStore } from '@/store'
import { formatDate } from '@/utils'

// 不解构 store: 解构后失去响应性, getUserInfo() 回来后这里的昵称和头像不会更新
const userStore = useUserStore()
const router = useRouter()

const homeInfo = ref({
  view_count: 0,
  user_count: 0,
  article_count: 0,
  message_count: 0,
})

// 待处理事项: 数字都来自各自列表接口的 total, 不需要新增后端接口
const todo = ref({ comment: 0, message: 0, recycle: 0, online: 0 })
const latestArticles = ref([])
const categoryStat = ref([])
const loginLogs = ref([])
const loading = ref(true)

const STATUS_MAP = {
  1: { text: '公开', type: 'success' },
  2: { text: '私密', type: 'warning' },
  3: { text: '草稿', type: 'default' },
}

// 分类分布用条形图展示占比, 按最大值归一化(不是按总数), 差异看得更清楚
const maxCategoryCount = computed(() =>
  Math.max(1, ...categoryStat.value.map(c => c.article_count ?? 0)),
)

const todoItems = computed(() => [
  { label: '待审核评论', value: todo.value.comment, path: '/message/comment' },
  { label: '待审核留言', value: todo.value.message, path: '/message/leave-msg' },
  { label: '回收站文章', value: todo.value.recycle, path: '/article/list' },
  { label: '当前在线用户', value: todo.value.online, path: '/user/online' },
])

// 每块数据独立取, 用 allSettled: 某个接口挂了不该让整个面板空着
async function fetchDashboard() {
  loading.value = true
  const [comment, message, recycle, articles, categories, logs, online] = await Promise.allSettled([
    api.getComments({ page_num: 1, page_size: 1, is_review: false }),
    api.getMessages({ page_num: 1, page_size: 1, is_review: false }),
    api.getArticles({ page_num: 1, page_size: 1, is_delete: true }),
    api.getArticles({ page_num: 1, page_size: 5, is_delete: false }),
    api.getCategorys({ page_num: 1, page_size: 100 }),
    api.getLoginLogs({ page_num: 1, page_size: 5 }),
    api.getOnlineUsers({ keyword: '' }),
  ])

  const total = r => (r.status === 'fulfilled' ? r.value?.data?.total ?? 0 : 0)
  const list = r => (r.status === 'fulfilled' ? r.value?.data?.page_data ?? [] : [])

  todo.value = {
    comment: total(comment),
    message: total(message),
    recycle: total(recycle),
    online: online.status === 'fulfilled' ? (online.value?.data ?? []).length : 0,
  }
  latestArticles.value = list(articles)
  loginLogs.value = list(logs)
  categoryStat.value = list(categories)
    .slice()
    .sort((a, b) => (b.article_count ?? 0) - (a.article_count ?? 0))
    .slice(0, 5)
  loading.value = false
}

onMounted(async () => {
  getOneSentence()
  fetchDashboard()
  // 裸 await 会在接口失败时留下未捕获的 rejection;
  // data 为空时也不能让 homeInfo 变成 null, 模板里要取它的字段
  try {
    const res = await api.getHomeInfo()
    if (res.data) {
      homeInfo.value = res.data
    }
  }
  catch (err) {
    console.error(err)
  }
})

// 一言: 每次进后台看到一句新的文案
// 接口在部分网络下不通, 兜底从内置文案里随机取一句, 而不是固定显示同一句
const FALLBACK_SENTENCES = [
  '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。',
  '书山有路勤为径，学海无涯苦作舟。',
  '纸上得来终觉浅，绝知此事要躬行。',
  '不积跬步，无以至千里；不积小流，无以成江海。',
  '路漫漫其修远兮，吾将上下而求索。',
  '业不可不勤，勤则百弊自去。',
  '博观而约取，厚积而薄发。',
]
function randomSentence() {
  return FALLBACK_SENTENCES[Math.floor(Math.random() * FALLBACK_SENTENCES.length)]
}

// 先用内置文案填上, 一言回来了再替换
// 实测 v1.hitokoto.cn 要 0.9~1.5s, 之前这段时间这行是空的
const sentence = ref(randomSentence())
async function getOneSentence() {
  try {
    // 超时就放弃, 不然接口挂着时这个 promise 一直悬着
    const resp = await fetch('https://v1.hitokoto.cn?c=i', { signal: AbortSignal.timeout(2000) })
    const data = await resp.json()
    if (data?.hitokoto) {
      sentence.value = data.hitokoto
    }
  }
  catch {
    // 保持内置文案
  }
}
</script>

<template>
  <AppPage>
    <div class="flex-1">
      <!-- 问候 -->
      <NCard>
        <div class="flex items-center">
          <NAvatar round :size="60" :src="userStore.avatar" />
          <div class="ml-5">
            <p> Hello, {{ userStore.nickname }} </p>
            <NGradientText class="mt-1 op-60" gradient="linear-gradient(90deg, red 0%, green 50%, blue 100%)">
              {{ sentence }}
            </NGradientText>
          </div>
          <!-- 原来这里是两张 badgen.net 的徽章图片, 实测要 3.3s 才回来, 而且没写宽高,
               外网不通时高度塌成 0、整块跟着抖。改成一个普通链接, 不依赖外部服务 -->
          <div class="ml-auto flex items-center">
            <a
              class="flex items-center gap-1 text-sm transition-300 hover:text-primary"
              href="https://github.com/szluyu99/gin-vue-blog"
              target="_blank" rel="noopener noreferrer"
            >
              <span class="i-mdi:github text-xl" />
              项目仓库
            </a>
          </div>
        </div>
      </NCard>

      <!-- 总量统计: 图标色改用语义色 token, 原来是四个硬编码 hex -->
      <NGrid class="mt-4" x-gap="12" y-gap="12" cols="2 s:4" responsive="screen">
        <NGi
          v-for="item of [
            { icon: 'i-fa6-solid:users', color: 'text-primary', label: '访问量', key: 'view_count' },
            { icon: 'i-heroicons:users-solid', color: 'text-success', label: '用户量', key: 'user_count' },
            { icon: 'i-material-symbols:article', color: 'text-info', label: '文章量', key: 'article_count' },
            { icon: 'i-bxs:comment-dots', color: 'text-warning', label: '留言量', key: 'message_count' },
          ]" :key="item.key"
        >
          <NCard>
            <span class="text-[52px]" :class="[item.icon, item.color]" />
            <NStatistic class="float-right" :label="item.label">
              {{ homeInfo[item.key] ?? '-' }}
            </NStatistic>
          </NCard>
        </NGi>
      </NGrid>

      <!-- 待处理 + 最新文章 -->
      <NGrid class="mt-4" x-gap="12" y-gap="12" cols="1 l:24" responsive="screen">
        <NGi :span="8">
          <NCard title="待处理" size="small" class="h-full">
            <div class="space-y-1">
              <div
                v-for="item of todoItems" :key="item.label"
                class="flex cursor-pointer items-center justify-between rounded px-2 py-2 transition-300 hover:bg-primary/8"
                @click="router.push(item.path)"
              >
                <span class="text-sm">{{ item.label }}</span>
                <NSkeleton v-if="loading" :width="24" text />
                <!-- 0 的时候压暗, 有待办才醒目 -->
                <span
                  v-else class="text-lg font-bold"
                  :class="item.value ? 'text-warning' : 'op-40'"
                >
                  {{ item.value }}
                </span>
              </div>
            </div>
          </NCard>
        </NGi>

        <NGi :span="16">
          <NCard title="最新文章" size="small" class="h-full">
            <template #header-extra>
              <span class="cursor-pointer text-sm op-60 hover:text-primary" @click="router.push('/article/list')">
                全部
              </span>
            </template>
            <NSkeleton v-if="loading" :repeat="5" text class="my-2" />
            <NEmpty v-else-if="!latestArticles.length" class="py-6" description="还没有文章" />
            <div v-else class="space-y-1">
              <div
                v-for="article of latestArticles" :key="article.id"
                class="flex cursor-pointer items-center gap-3 rounded px-2 py-2 transition-300 hover:bg-primary/8"
                @click="router.push(`/article/write/${article.id}`)"
              >
                <span class="flex-1 truncate text-sm">{{ article.title }}</span>
                <NTag v-if="article.category" size="small" :bordered="false">
                  {{ article.category.name }}
                </NTag>
                <NTag size="small" :bordered="false" :type="(STATUS_MAP[article.status] || {}).type">
                  {{ (STATUS_MAP[article.status] || {}).text || '未知' }}
                </NTag>
                <span class="w-[80px] text-right text-xs op-60">
                  {{ formatDate(article.created_at) }}
                </span>
              </div>
            </div>
          </NCard>
        </NGi>
      </NGrid>

      <!-- 分类分布 + 最近登录 -->
      <NGrid class="mt-4" x-gap="12" y-gap="12" cols="1 l:2" responsive="screen">
        <NGi>
          <NCard title="分类分布" size="small" class="h-full">
            <template #header-extra>
              <span class="cursor-pointer text-sm op-60 hover:text-primary" @click="router.push('/article/category')">
                全部
              </span>
            </template>
            <NSkeleton v-if="loading" :repeat="4" text class="my-2" />
            <NEmpty v-else-if="!categoryStat.length" class="py-6" description="还没有分类" />
            <div v-else class="space-y-3">
              <div v-for="c of categoryStat" :key="c.id">
                <div class="mb-1 flex justify-between text-sm">
                  <span class="truncate">{{ c.name }}</span>
                  <span class="op-60">{{ c.article_count ?? 0 }} 篇</span>
                </div>
                <!-- 按最大值归一化, 不是按总数: 分类多的时候按总数算每条都是细线 -->
                <NProgress
                  type="line" :height="6" :border-radius="3"
                  :percentage="Math.round((c.article_count ?? 0) / maxCategoryCount * 100)"
                  :show-indicator="false"
                />
              </div>
            </div>
          </NCard>
        </NGi>

        <NGi>
          <NCard title="最近登录" size="small" class="h-full">
            <template #header-extra>
              <span class="cursor-pointer text-sm op-60 hover:text-primary" @click="router.push('/log/login')">
                全部
              </span>
            </template>
            <NSkeleton v-if="loading" :repeat="5" text class="my-2" />
            <NEmpty v-else-if="!loginLogs.length" class="py-6" description="还没有登录记录" />
            <div v-else class="space-y-1">
              <div
                v-for="log of loginLogs" :key="log.id"
                class="flex items-center gap-2 px-2 py-1.5 text-sm"
              >
                <span class="w-[70px] truncate">{{ log.nickname || log.username || '未知' }}</span>
                <span class="flex-1 truncate text-xs op-60">{{ log.ip_address }} · {{ log.ip_source || '未知' }}</span>
                <NTag size="small" :bordered="false" :type="log.status === 1 ? 'success' : 'error'">
                  {{ log.status === 1 ? '成功' : '失败' }}
                </NTag>
                <span class="w-[125px] text-right text-xs op-60">
                  {{ formatDate(log.created_at, 'MM-DD HH:mm:ss') }}
                </span>
              </div>
            </div>
          </NCard>
        </NGi>
      </NGrid>
    </div>
  </AppPage>
</template>
