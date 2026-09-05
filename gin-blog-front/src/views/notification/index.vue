<script setup>
import dayjs from 'dayjs'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'
import UPagination from '@/components/ui/UPagination.vue'
import { useAppStore, useNotificationStore, useUserStore } from '@/store'
import { convertImgUrl } from '@/utils'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()
const notificationStore = useNotificationStore()

// 头部铃铛只放最近 10 条, 这一页才是全部, 所以自己分页而不是复用 store 里的列表
const PAGE_SIZE = 10
const loading = ref(true)
const list = ref([])
const total = ref(0)
const current = ref(1)
const pageCount = computed(() => Math.ceil(total.value / PAGE_SIZE))

watch(current, () => {
  getNotifications()
  window.scrollTo({ behavior: 'smooth', top: 0 })
})

async function getNotifications() {
  if (!userStore.userId) {
    loading.value = false
    return
  }
  loading.value = true
  try {
    const resp = await api.getNotifications({ page_num: current.value, page_size: PAGE_SIZE })
    list.value = resp.data?.page_data ?? []
    total.value = resp.data?.total ?? 0
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
}

// 点一条: 先标已读再跳到文章
async function open(item) {
  if (!item.is_read) {
    await notificationStore.read([item.id])
    item.is_read = true
  }
  if (item.article_id) {
    router.push(`/article/${item.article_id}`)
  }
}

async function readAll() {
  await notificationStore.read()
  for (const item of list.value) {
    item.is_read = true
  }
}

onMounted(getNotifications)

// 未登录时进这个页面: 直接弹登录框, 登录成功后再拉数据
watch(() => userStore.userId, (id) => {
  if (id) {
    getNotifications()
  }
}, { immediate: false })
</script>

<template>
  <BannerPage title="站内通知" label="message" :loading="loading" card>
    <template v-if="!userStore.userId">
      <div class="py-10 text-center space-y-3">
        <p class="color-muted">
          登录后才能查看站内通知
        </p>
        <button class="rounded-md bg-primary px-4 py-1.5 text-white" @click="appStore.setLoginFlag(true)">
          去登录
        </button>
      </div>
    </template>
    <template v-else>
      <div class="flex items-center justify-between pb-4">
        <p class="text-lg">
          共 {{ total }} 条通知
        </p>
        <button
          v-if="notificationStore.unreadCount"
          class="text-sm text-primary" @click="readAll"
        >
          全部标为已读
        </button>
      </div>

      <div v-if="list.length" class="space-y-2">
        <div
          v-for="item of list" :key="item.id"
          class="flex cursor-pointer items-start gap-3 rounded-xl bg-surface-soft p-3 transition-300 hover:shadow-md hover:-translate-y-0.5"
          :class="item.is_read ? 'op-70' : ''"
          @click="open(item)"
        >
          <img
            :src="convertImgUrl(item.from_avatar)" :alt="item.from_nickname"
            class="h-10 w-10 shrink-0 rounded-full bg-surface object-cover"
          >
          <div class="min-w-0 flex-1">
            <p class="text-sm">
              <span class="font-bold">{{ item.from_nickname || '有人' }}</span>
              {{ item.type === 1 ? '回复了你' : '评论了你的文章' }}
              <span class="ml-1 text-xs color-muted">{{ dayjs(item.created_at).format('YYYY-MM-DD HH:mm') }}</span>
            </p>
            <p class="mt-1 text-sm color-muted">
              {{ item.content }}
            </p>
            <p v-if="item.article_title" class="mt-1 truncate text-xs color-muted">
              来自《{{ item.article_title }}》
            </p>
          </div>
          <span v-if="!item.is_read" class="mt-2 h-2 w-2 shrink-0 rounded-full bg-accent" />
        </div>
      </div>
      <div v-else-if="!loading" class="py-10 text-center color-muted">
        还没有通知
      </div>

      <div v-if="pageCount > 1" class="mt-8 flex justify-center">
        <UPagination v-model:page="current" :page-count="pageCount" />
      </div>
    </template>
  </BannerPage>
</template>
