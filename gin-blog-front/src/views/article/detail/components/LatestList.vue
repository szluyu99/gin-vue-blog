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
        <!-- 分隔线要挂在 li 上: 原来写在 RouterLink 上, 而它是 li 的唯一子元素,
             last:border-none 恒命中, 虚线分隔从来没显示过(实测 border 是 none/0px) -->
        <li
          v-for="item of articleList" :key="item.id"
          class="border-b-1 border-color-divider border-dashed p-1 last:border-none hover:bg-surface-soft"
        >
          <RouterLink :to="`/article/${item.id}`">
            <div class="flex items-center">
              <!-- object-cover: 原来是默认的 fill, 非正方形封面会被拉伸变形 -->
              <img
                :src="convertImgUrl(item.img)"
                class="h-15 w-15 shrink-0 rounded bg-surface-soft object-cover"
                loading="lazy" :alt="item.title"
              >
              <div class="min-w-0 flex-1 pl-2">
                <p class="line-clamp-2">
                  {{ item.title }}
                </p>
                <p class="text-sm text-muted">
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
