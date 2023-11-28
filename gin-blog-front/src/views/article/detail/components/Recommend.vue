<script setup>
import { NGi, NGrid } from 'naive-ui'
import { convertImgUrl, formatDate } from '@/utils'

const { recommendList } = defineProps({
  recommendList: Array,
})
</script>

<template>
  <div v-if="recommendList.length">
    <p class="mb-15 flex items-center text-20 font-bold">
      <span class="i-mdi:thumb-up mr-5 text-24 text-orange" /> 相关推荐
    </p>
    <NGrid :cols="12" :x-gap="6" :y-gap="6" item-responsive responsive="screen">
      <template v-for="item of recommendList" :key="item.id">
        <NGi span="12 m:4" class="art-card">
          <RouterLink :to="`/article/${item.id}`">
            <div class="relative h-150 f-c-c overflow-hidden rounded-2 bg-black lg:h-200">
              <img :src="convertImgUrl(item.img)" class="art-img h-full w-full opacity-40 transition-600">
              <div class="absolute text-center text-white">
                <p>
                  <span class="i-mdi:calendar text-16" />
                  {{ formatDate(item.created_at) }}
                </p>
                <p> {{ item.title }} </p>
              </div>
            </div>
          </RouterLink>
        </NGi>
      </template>
    </NGrid>
  </div>
</template>

<style scoped>
.art-img {
  object-fit: cover;
}

.art-card:hover .art-img {
  opacity: 0.8;
  transform: scale(1.1);
}
</style>

