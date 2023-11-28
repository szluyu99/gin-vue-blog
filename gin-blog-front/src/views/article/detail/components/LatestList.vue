<script setup>
import { convertImgUrl, formatDate } from '@/utils'

const { articleList } = defineProps({
  articleList: Array,
})
</script>

<template>
  <Transition name="slide-fade" appear>
    <div class="card-view">
      <div class="flex items-center">
        <span class="i-akar-icons:clock text-6xl" />
        <span class="ml-10 text-16"> 最新文章 </span>
      </div>
      <ul>
        <li v-for="item of articleList" :key="item.title">
          <RouterLink :to="`/article/${item.id}`" class="my-3 border-b-1 border-dashed last:border-none">
            <div class="flex items-center py-6">
              <img :src="convertImgUrl(item.img)" class="h-58 w-58 overflow-hidden">
              <div class="flex-1 overflow-hidden break-all pl-10">
                <p class="text-14">
                  {{ item.title }}
                </p>
                <p class="text-13 text-blueGray">
                  {{ formatDate(item.created_at) }}
                </p>
              </div>
            </div>
          </RouterLink>
        </li>
      </ul>
    </div>
  </Transition>
</template>
