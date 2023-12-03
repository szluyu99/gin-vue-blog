<script setup>
import { computed, ref } from 'vue'
import { debouncedWatch } from '@vueuse/core'
import { NInput } from 'naive-ui'

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
  const res = await api.searchArticles({ keyword: keyword.value })
  articleList.value = res.data
}
</script>

<template>
  <UModal v-model="searchFlag" :width="600">
    <div class="m-0">
      <div class="mb-15 text-18 font-bold">
        本地搜索
      </div>
      <NInput
        v-model:value="keyword"
        round
        placeholder="输入文章标题或内容..."
        size="large"
        clearable
      >
        <template #prefix>
          <div class="i-mdi:flash text-20 text-yellow" />
        </template>
      </NInput>
      <hr class="my-15 border-2px border-color-#d2ebfd border-dashed">
      <div class="h-400">
        <ul v-if="articleList.length">
          <li v-for="item of articleList" :key="item.id">
            <RouterLink :to="`/article/${item.id}`">
              <span
                class="border-b-1 border-#999 border-solid text-17"
                @click="searchFlag = false"
                v-html="item.title"
              />
            </RouterLink>
            <p clsas="color-#555 mt-5 cursor-pointer ell-3" v-html="item.content" />
            <hr class="my-10 border-1 border-#d2ebfd border-dashed">
          </li>
        </ul>
        <div v-else-if="keyword">
          找不到您查询的内容：{{ keyword }}
        </div>
      </div>
    </div>
  </UModal>
</template>

<style scoped lang="scss">
// 省略文字最多 N 行
.ell-3 {
  display: -webkit-box;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}
</style>

