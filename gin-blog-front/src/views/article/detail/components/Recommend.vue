<script setup>
import dayjs from 'dayjs'
import { convertImgUrl } from '@/utils'

defineProps({
  recommendList: Array,
})
</script>

<template>
  <div v-if="recommendList.length" class="space-y-3">
    <p class="flex items-center text-xl font-bold">
      <span class="i-mdi:thumb-up mr-3 text-orange" /> 相关推荐
    </p>
    <div class="grid grid-cols-12 gap-1.5">
      <div v-for="item of recommendList" :key="item.id" class="group col-span-12 lg:col-span-4">
        <RouterLink :to="`/article/${item.id}`">
          <div class="relative h-[150px] f-c-c overflow-hidden rounded bg-black lg:h-[200px]">
            <img :src="convertImgUrl(item.img)" class="h-full w-full object-cover opacity-40 transition-600 group-hover:scale-110 group-hover:opacity-80">
            <div class="absolute text-center text-white">
              <p class="f-c-c">
                <span class="i-mdi:calendar mr-1" /> {{ dayjs(item.created_at).format('YYYY-MM-DD') }}
              </p>
              <p> {{ item.title }} </p>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>
