<script setup>
import dayjs from 'dayjs'
import { convertImgUrl } from '@/utils'

const { articleList } = defineProps({
  articleList: Array,
})
</script>

<template>
  <Transition name="slide-fade" appear>
    <div class="card-view space-y-3">
      <div class="flex items-center">
        <span class="i-akar-icons:clock" />
        <span class="ml-2"> 最新文章 </span>
      </div>
      <ul class="space-y-1">
        <li v-for="item of articleList" :key="item.title" class="p-1 hover:bg-gray-200">
          <RouterLink :to="`/article/${item.id}`" class="border-b-1 border-dashed last:border-none">
            <div class="flex items-center">
              <img :src="convertImgUrl(item.img)" class="h-15 w-15 overflow-hidden">
              <div class="flex-1 overflow-hidden break-all pl-2">
                <p> {{ item.title }} </p>
                <p class="text-sm text-blueGray">
                  {{ dayjs(item.created_at).format('YYYY-MM-DD') }}
                </p>
              </div>
            </div>
          </RouterLink>
        </li>
      </ul>
    </div>
  </Transition>
</template>
