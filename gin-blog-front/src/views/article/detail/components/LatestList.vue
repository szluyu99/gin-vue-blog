<script setup lang="ts">
import { convertImgUrl, formatDate } from '@/utils'
const { articleList = [] } = defineProps<{ articleList: any }>()
</script>

<template>
  <Transition name="slide-fade" appear>
    <div card-view>
      <div flex items-center>
        <i-akar-icons:clock text-6xl />
        <span text-16 ml-10> 最新文章 </span>
      </div>
      <ul>
        <router-link
          v-for="item of articleList"
          :key="item.title" :to="`/article/${item.id}`"
          my-3 border-b-1 border-dashed last:border-none
        >
          <div flex items-center py-6>
            <img :src="convertImgUrl(item.img)" overflow-hidden w-58 h-58>
            <div flex-1 pl-10 break-all overflow-hidden>
              <p text-14>
                {{ item.title }}
              </p>
              <p text-13 text-blueGray>
                {{ formatDate(item.created_at) }}
              </p>
            </div>
          </div>
        </router-link>
      </ul>
    </div>
  </Transition>
</template>
