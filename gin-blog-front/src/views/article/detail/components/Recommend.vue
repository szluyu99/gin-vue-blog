<script setup lang="ts">
import { convertImgUrl, formatDate } from '@/utils'
const { recommendList = [] } = defineProps<{ recommendList: any }>()
</script>

<template>
  <div v-if="recommendList.length">
    <p mb-15 text-20 font-bold flex items-center>
      <i-mdi:thumb-up mr-5 text-24 text-orange /> 相关推荐
    </p>
    <n-grid
      :cols="12" :x-gap="6" :y-gap="6"
      item-responsive responsive="screen"
    >
      <template v-for="item of recommendList" :key="item.id">
        <n-gi span="12 m:4" class="art-card">
          <router-link :to="`/article/${item.id}`">
            <div
              relative f-c-c bg-black overflow-hidden rounded-2
              h-150 lg:h-200
            >
              <img
                :src="convertImgUrl(item.img)"
                class="art-img" transition-600 opacity-40 wh-full
              >
              <div absolute text-white text-center>
                <p>
                  <i-mdi:calendar text-16 />
                  {{ formatDate(item.created_at) }}
                </p>
                <p> {{ item.title }} </p>
              </div>
            </div>
          </router-link>
        </n-gi>
      </template>
    </n-grid>
  </div>
</template>

<style lang="scss" scoped>
.art-img {
  object-fit: cover;
}

.art-card:hover .art-img {
  opacity: 0.8;
  transform: scale(1.1);
}
</style>

