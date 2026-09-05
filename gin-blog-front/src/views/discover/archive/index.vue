<script setup>
import dayjs from 'dayjs'
import { computed, onMounted, ref, watch } from 'vue'

import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'
import UPagination from '@/components/ui/UPagination.vue'

const loading = ref(true)
const total = ref(0)
const archiveList = ref([])

// 一页 50 条: 归档条目很轻(只有日期和标题), 一页多放一些翻页次数少
const PAGE_SIZE = 50
const current = ref(1)
const pageCount = computed(() => Math.ceil(total.value / PAGE_SIZE))

watch(current, () => {
  getArchives()
  window.scrollTo({ behavior: 'smooth', top: 0 })
})

async function getArchives() {
  // 以前失败时 loading 永远停在 true, 页面卡在加载态;
  // resp.data 为空时直接取 .page_data 还会抛
  try {
    const resp = await api.getArchives({
      page_num: current.value,
      page_size: PAGE_SIZE,
    })
    archiveList.value = resp.data?.page_data ?? []
    total.value = resp.data?.total ?? 0
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
}

// 按年月分组: 归档页的意义就是时间轴, 平铺一长串看不出节奏
// 接口已经按发布时间倒序返回, 所以直接按出现顺序建组, 不用再排一次
const monthGroups = computed(() => {
  const groups = []
  for (const item of archiveList.value) {
    const key = dayjs(item.created_at).format('YYYY-MM')
    const last = groups[groups.length - 1]
    if (last?.key === key) {
      last.items.push(item)
      continue
    }
    groups.push({
      key,
      label: dayjs(item.created_at).format('YYYY 年 M 月'),
      items: [item],
    })
  }
  return groups
})

onMounted(() => {
  getArchives()
})
</script>

<template>
  <BannerPage title="归档" label="archive" :loading="loading" card>
    <p class="pb-5 text-lg lg:text-2xl">
      <!-- 用 total 而不是 archiveList.length: 后者只是当页条数, 分页后会少报 -->
      目前共计 {{ total }} 篇文章，继续加油！
    </p>

    <section v-for="group of monthGroups" :key="group.key" class="mb-7 last:mb-0">
      <h3 class="mb-3 flex items-center gap-2 text-lg font-bold">
        <span class="i-mdi:calendar-blank-outline text-primary" />
        {{ group.label }}
        <span class="text-sm text-muted font-normal">{{ group.items.length }} 篇</span>
      </h3>
      <!-- 左侧竖线 + 圆点构成时间轴, 圆点用 ring 描一圈卡片底色, 压在竖线上更干净 -->
      <ul class="ml-1.5 border-l-2 border-color-divider pl-6 space-y-3">
        <li v-for="item of group.items" :key="item.id" class="relative">
          <span class="absolute left-[-31px] top-2.5 h-2.5 w-2.5 rounded-full bg-primary ring-3 ring-surface" />
          <!-- 原来是 <a @click="router.push()">, 没有 href, 中键新开和复制链接都用不了 -->
          <RouterLink
            :to="`/article/${item.id}`"
            class="flex flex-wrap items-baseline gap-x-3 transition-300 hover:text-primary"
          >
            <span class="text-sm text-muted">{{ dayjs(item.created_at).format('MM-DD') }}</span>
            <span class="lg:text-lg">{{ item.title }}</span>
          </RouterLink>
        </li>
      </ul>
    </section>

    <div v-if="!loading && !archiveList.length" class="py-10 text-center text-muted">
      还没有文章
    </div>

    <div v-if="pageCount > 1" class="mt-10 flex justify-center">
      <UPagination v-model:page="current" :page-count="pageCount" />
    </div>
  </BannerPage>
</template>
