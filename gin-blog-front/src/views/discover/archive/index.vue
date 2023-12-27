<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'

import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'

const router = useRouter()

const loading = ref(true)
const total = ref(0)
const archiveList = ref([])

const current = ref(1) // 当前页数
watch(current, () => getArchives())

async function getArchives() {
  const resp = await api.getArchives({
    page_num: current.value,
    page_size: 50,
  })
  archiveList.value = resp.data.page_data
  total.value = resp.data.total
  loading.value = false
}

onMounted(() => {
  getArchives()
})

// const activity = [
//   { id: 1, type: 'created', person: { name: 'Chelsea Hagon' }, date: '7d ago', dateTime: '2023-01-23T10:32' },
//   { id: 2, type: 'edited', person: { name: 'Chelsea Hagon' }, date: '6d ago', dateTime: '2023-01-23T11:03' },
//   { id: 3, type: 'sent', person: { name: 'Chelsea Hagon' }, date: '6d ago', dateTime: '2023-01-23T11:24' },
//   {
//     id: 4,
//     type: 'commented',
//     person: {
//       name: 'Chelsea Hagon',
//       imageUrl:
//         'https://images.unsplash.com/photo-1550525811-e5869dd03032?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
//     },
//     comment: 'Called client, they reassured me the invoice would be paid by the 25th.',
//     date: '3d ago',
//     dateTime: '2023-01-23T15:56',
//   },
//   { id: 5, type: 'viewed', person: { name: 'Alex Curren' }, date: '2d ago', dateTime: '2023-01-24T09:12' },
//   { id: 6, type: 'paid', person: { name: 'Alex Curren' }, date: '1d ago', dateTime: '2023-01-24T09:20' },
// ]
</script>

<template>
  <BannerPage title="归档" label="archive" :loading="loading" card>
    <p class="pb-5 text-lg lg:text-2xl">
      目前共计 {{ archiveList.length }} 篇文章，继续加油！
    </p>
    <template v-for="(item, idx) of archiveList" :key="item.id">
      <div class="flex items-center gap-2">
        <div class="i-mdi:circle bg-blue text-sm" />
        <span class="text-sm color-#666 lg:text-base">
          {{ dayjs(item.created_at).format('YYYY-MM-DD') }}
        </span>
        <a class="color-#666 lg:text-lg hover:text-orange" @click="router.push(`/article/${item.id}`)">
          {{ item.title }}
        </a>
      </div>
      <hr v-if="idx !== archiveList.length - 1" class="my-4 border-1 border-color-#d2ebfd border-dashed">
    </template>

    <!-- TODO: 分页 -->
    <!-- <div class="my-15 mt-20 f-c-c">
      <NPagination v-model:page="current" :page-count="Math.ceil(total / 10)" />
    </div> -->
  </BannerPage>
</template>
