<script setup lang="ts">
import { formatDate } from '@/utils'
import api from '@/api'

const router = useRouter()

let loading = $ref(true)
let total = $ref(0)
let archiveList = $ref<any>([])

const current = $ref(1) // 当前页数
watch($$(current), () => getArchives())

async function getArchives() {
  const res = await api.getArchives({
    page_num: current,
    page_size: 10,
  })
  archiveList = res.data.pageData
  total = res.data.total
  loading = false
}

onMounted(() => {
  getArchives()
})
</script>

<template>
  <BannerPage title="归档" label="archive" :loading="loading" card>
    <n-timeline :icon-size="18" size="large">
      <n-timeline-item>
        <template #header>
          <p text-18 pb-20 lg="text-24">
            目前共计 {{ archiveList.length }} 篇文章，继续加油！
          </p>
        </template>
      </n-timeline-item>
      <template v-for="(item, idx) of archiveList" :key="item.id">
        <n-timeline-item type="info" :content="item.desc">
          <template #header>
            <span color="#666" mr-15 lg="text-15">
              {{ formatDate(item.created_at) }}
            </span>
            <a
              color="#666" lg="text-16"
              @click="router.push(`/article/${item.id}`)"
            >
              {{ item.title }}
            </a>
          </template>
          <template #icon>
            <i-mdi:circle />
          </template>
          <template #footer>
            <hr
              v-if="idx !== archiveList.length - 1"
              mt-15 mb-15 border-dashed border-1px border-color="#d2ebfd"
            >
          </template>
        </n-timeline-item>
      </template>
    </n-timeline>
    <!-- 分页按钮 -->
    <div f-c-c mt-20 my-15>
      <n-pagination
        v-model:page="current"
        :page-count="Math.ceil(total / 10)"
      />
    </div>
  </BannerPage>
</template>
