<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NPagination, NTimeline, NTimelineItem } from 'naive-ui'

import BannerPage from '@/components/page/BannerPage.vue'
import { formatDate } from '@/utils'
import api from '@/api'

const router = useRouter()

const loading = ref(true)
const total = ref(0)
const archiveList = ref([])

const current = ref(1) // 当前页数
watch(current, () => getArchives())

async function getArchives() {
  const res = await api.getArchives({
    page_num: current.value,
    page_size: 10,
  })
  archiveList.value = res.data.pageData
  total.value = res.data.total
  loading.value = false
}

onMounted(() => {
  getArchives()
})
</script>

<template>
  <BannerPage title="归档" label="archive" :loading="loading" card>
    <NTimeline :icon-size="18" size="large">
      <NTimelineItem>
        <template #header>
          <p class="pb-20 text-18 lg:text-24">
            目前共计 {{ archiveList.length }} 篇文章，继续加油！
          </p>
        </template>
      </NTimelineItem>
      <template v-for="(item, idx) of archiveList" :key="item.id">
        <NTimelineItem type="info" :content="item.desc">
          <template #header>
            <span class="mr-15 color-#666 lg:text-15">
              {{ formatDate(item.created_at) }}
            </span>
            <a class="color-#666 lg:text-16" @click="router.push(`/article/${item.id}`)">
              {{ item.title }}
            </a>
          </template>
          <template #icon>
            <div class="i-mdi:circle" />
          </template>
          <template #footer>
            <hr v-if="idx !== archiveList.length - 1" class="mb-15 mt-15 border-1px border-color-#d2ebfd border-dashed">
          </template>
        </NTimelineItem>
      </template>
    </NTimeline>
    <!-- 分页按钮 -->
    <div class="my-15 mt-20 f-c-c">
      <NPagination v-model:page="current" :page-count="Math.ceil(total / 10)" />
    </div>
  </BannerPage>
</template>
