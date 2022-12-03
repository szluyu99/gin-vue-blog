<script setup lang="ts">
import dayjs from 'dayjs'

interface Props { articleList: any }
const { articleList = [] } = defineProps<Props>()
</script>

<template>
  <n-card hoverable rounded-2rem>
    <div flex items-center>
      <i-akar-icons:clock text-6xl /> <span text-16 ml-10>最新文章</span>
    </div>
    <ul>
      <router-link
        v-for="item of articleList"
        :key="item.title" :to="`/article/${item.id}`"
        my-3 border-b-1 border-dashed last:border-none
      >
        <div flex items-center py-6>
          <img :src="item.img" overflow-hidden w-58 h-58>
          <div flex-1 pl-10 break-all overflow-hidden>
            <p text-14>
              {{ item.title }}
            </p>
            <p text-13 text-blueGray>
              {{ dayjs(item.created_at).format('YYYY-MM-DD') }}
            </p>
          </div>
        </div>
      </router-link>
    </ul>
  </n-card>
</template>
