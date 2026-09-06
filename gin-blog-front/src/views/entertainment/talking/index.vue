<script setup>
import dayjs from 'dayjs'
import { computed, onMounted, ref, watch } from 'vue'

import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'
import UPagination from '@/components/ui/UPagination.vue'
import { convertImgUrl } from '@/utils'

const PAGE_SIZE = 10

const loading = ref(true)
const list = ref([])
const total = ref(0)
const current = ref(1)
const pageCount = computed(() => Math.ceil(total.value / PAGE_SIZE))

watch(current, () => {
  getTalks()
  window.scrollTo({ behavior: 'smooth', top: 0 })
})

async function getTalks() {
  loading.value = true
  try {
    const resp = await api.getTalks({ page_num: current.value, page_size: PAGE_SIZE })
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

onMounted(getTalks)
</script>

<template>
  <BannerPage label="talk" title="说说" card :loading="loading && !list.length">
    <p class="pb-4 text-lg">
      共 {{ total }} 条说说
    </p>

    <div v-if="list.length" class="space-y-4">
      <RouterLink
        v-for="(talk, idx) of list" :key="talk.id"
        :to="`/talk/${talk.id}`"
        class="card-enter block rounded-xl bg-surface-soft p-4 transition-300 hover:shadow-md hover:-translate-y-0.5"
        :style="{ '--i': idx }"
      >
        <div class="flex items-center gap-2">
          <img
            :src="convertImgUrl(talk.avatar)" :alt="talk.nickname"
            class="h-9 w-9 shrink-0 rounded-full bg-surface object-cover"
          >
          <div class="min-w-0 flex-1">
            <p class="flex items-center gap-2 text-sm">
              <span class="font-bold"> {{ talk.nickname || '博主' }} </span>
              <span v-if="talk.is_top" class="rounded bg-accent px-1.5 py-0.5 text-xs text-white">置顶</span>
            </p>
            <p class="text-xs color-muted">
              {{ dayjs(talk.created_at).format('YYYY-MM-DD HH:mm') }}
            </p>
          </div>
        </div>

        <!-- 说说是纯文本, 列表页最多显示 4 行 -->
        <p class="line-clamp-4 mt-3 whitespace-pre-wrap text-sm">
          {{ talk.content }}
        </p>

        <p class="mt-3 flex items-center gap-1 text-xs color-muted">
          <span class="i-mdi:comment-outline" />
          {{ talk.comment_count || 0 }} 条评论
        </p>
      </RouterLink>
    </div>
    <div v-else-if="!loading" class="py-10 text-center color-muted">
      还没有说说
    </div>

    <div v-if="pageCount > 1" class="mt-8 flex justify-center">
      <UPagination v-model:page="current" :page-count="pageCount" />
    </div>
  </BannerPage>
</template>
