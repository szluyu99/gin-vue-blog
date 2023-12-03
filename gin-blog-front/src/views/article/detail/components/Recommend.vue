<script setup>
import { convertImgUrl, formatDate } from '@/utils'

defineProps({
  recommendList: Array,
})
</script>

<template>
  <div v-if="recommendList.length">
    <p class="mb-15 flex items-center text-20 font-bold">
      <span class="i-mdi:thumb-up mr-6 text-24 text-orange" /> 相关推荐
    </p>
    <div class="grid grid-cols-12 gap-6">
      <template v-for="item of recommendList" :key="item.id">
        <div class="group col-span-12 lg:col-span-4">
          <RouterLink :to="`/article/${item.id}`">
            <div class="relative h-150 f-c-c overflow-hidden rounded-2 bg-black lg:h-200">
              <img :src="convertImgUrl(item.img)" class="h-full w-full object-cover opacity-40 transition-600 group-hover:scale-110 group-hover:opacity-80">
              <div class="absolute text-center text-white">
                <p>
                  <span class="i-mdi:calendar text-16" />
                  {{ formatDate(item.created_at) }}
                </p>
                <p> {{ item.title }} </p>
              </div>
            </div>
          </RouterLink>
        </div>
      </template>
    </div>
  </div>
</template>
