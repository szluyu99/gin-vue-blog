<script setup lang="ts">
import dayjs from 'dayjs'
import api from '@/api'

const router = useRouter()

const loading = ref(true)
const archiveList = ref<any>([])
const total = ref(0)

const current = ref(1) // 当前页数
watch(current, () => getArchives())

function getArchives() {
  api.getArchives({ page_num: current.value, page_size: 10 }).then((res) => {
    archiveList.value = res.data.pageData
    total.value = res.data.total
    loading.value = false
  })
}

onMounted(() => {
  getArchives()
})
</script>

<template>
  <BannerCard
    :loading="loading"
    title="归档"
    banner-img="https://static.talkxj.com/config/643f28683e1c59a80ccfc9cb19735a9c.jpg"
  >
    <n-timeline :icon-size="20" size="large">
      <n-timeline-item>
        <template #header>
          <p text-24 pb-20>
            目前共计 {{ archiveList.length }} 篇文章，继续加油！
          </p>
        </template>
      </n-timeline-item>
      <template v-for="(item, idx) of archiveList" :key="item.id">
        <n-timeline-item type="info" :content="item.desc">
          <template #header>
            <span color="#666" text-15 mr-15>
              {{ dayjs(item.created_at).format("YYYY-MM-DD") }}
            </span>
            <button color="#666" text-18 @click="router.push(`/article/${item.id}`)">
              {{ item.title }}
            </button>
          </template>
          <template #icon>
            <i-mdi:timeline />
          </template>
          <template #footer>
            <n-divider v-if="idx !== archiveList.length - 1" dashed />
          </template>
        </n-timeline-item>
      </template>
      <!-- 分页按钮 -->
      <v-pagination
        v-model="current"
        :length="Math.ceil(total / 10)"
        total-visible="6"
      />
    </n-timeline>
  </BannerCard>
</template>

<style lang="scss" scoped>
</style>
