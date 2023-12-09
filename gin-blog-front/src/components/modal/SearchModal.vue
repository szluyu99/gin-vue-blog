<script setup>
import { computed, ref } from 'vue'
import { debouncedWatch } from '@vueuse/core'

import UModal from '@/components/ui/UModal.vue'
import api from '@/api'
import { useAppStore } from '@/store'

const appStore = useAppStore()

const searchFlag = computed({
  get: () => appStore.searchFlag,
  set: val => appStore.setSearchFlag(val),
})

// 搜索关键字
const keyword = ref('')
// 搜索结果
const articleList = ref([])

// 防抖 watch, 节流: throttledWatch
debouncedWatch(
  keyword,
  () => keyword.value ? handleSearch() : articleList.value = [],
  { debounce: 300 },
)

async function handleSearch() {
  const resp = await api.searchArticles({ keyword: keyword.value })
  articleList.value = resp.data
}
</script>

<template>
  <UModal v-model="searchFlag" :width="600">
    <div class="m-0">
      <div class="mb-4 text-xl font-bold">
        本地搜索
      </div>
      <div>
        <div class="relative rounded-md shadow-sm">
          <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-2">
            <div class="i-mdi:flash text-xl text-yellow" />
          </div>
          <input
            v-model="keyword"
            class="block w-full border-0 rounded-full py-2 pl-8 pr-5 text-gray-900 outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-green-300"
            placeholder="输入文章标题或内容..."
          >
        </div>
        <hr class="my-4 border-1.5 border-color-#d2ebfd border-dashed">
        <div class="h-[420px] overflow-y-auto">
          <ul v-if="articleList.length">
            <li v-for="item of articleList" :key="item.id" class="text-sm">
              <RouterLink :to="`/article/${item.id}`">
                <span class="border-b-1 border-#999 border-solid text-lg" @click="searchFlag = false" v-html="item.title" />
              </RouterLink>
              <div class="ell-4 mt-1">
                <p clsas="color-#555 cursor-pointer" v-html="item.content" />
              </div>
              <hr class="my-3 border-1 border-#d2ebfd border-dashed">
            </li>
          </ul>
          <div v-else-if="keyword">
            找不到您查询的内容：{{ keyword }}
          </div>
        </div>
      </div>
    </div>
  </UModal>
</template>

<style scoped>
/* 省略文字最多 N 行 */
.ell-4 {
  display: -webkit-box;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
}
</style>
