<script setup lang="ts">
import api from '@/api'
import { useAppStore } from '@/store'
const appStore = useAppStore()

const searchFlag = $computed({
  get: () => appStore.searchFlag,
  set: val => appStore.setSearchFlag(val),
})

// 搜索关键字
const keyword = $ref('')
// 搜索结果
let articleList = $ref<any>([])
// 防抖 watch, 节流: throttledWatch
debouncedWatch(
  $$(keyword),
  () => keyword ? handleSearch() : articleList = [],
  { debounce: 300 },
)

async function handleSearch() {
  const res = await api.searchArticles({ keyword })
  articleList = res.data
}
</script>

<template>
  <n-modal
    v-model:show="searchFlag"
    display-directive="show"
    preset="card"
    title="本地搜索"
    transform-origin="center"
    :block-scroll="appStore.isMobile"
    px-10 h-full max-w-360
    lg="max-w-600"
  >
    <n-input
      v-model:value="keyword"
      round
      placeholder="输入文章标题或内容..."
      size="large"
      clearable
    >
      <template #prefix>
        <i-mdi:flash text-20 text-yellow />
      </template>
    </n-input>
    <hr my-15 border-dashed border-2px border-color="#d2ebfd">
    <div h-400>
      <ul v-if="articleList.length">
        <li v-for="item of articleList" :key="item.id">
          <router-link :to="`/article/${item.id}`">
            <span
              text-17 border-b-1 border-solid border="#999"
              @click="searchFlag = false"
              v-html="item.title"
            />
          </router-link>
          <p
            cursor-pointer color="#555" mt-5
            class="ell-3"
            v-html="item.content"
          />
          <hr my-10 border-1 border-dashed border="#d2ebfd">
        </li>
      </ul>
      <div v-else-if="keyword">
        找不到您查询的内容：{{ keyword }}
      </div>
    </div>
  </n-modal>
</template>

<style scoped lang="scss">
// 省略文字最多 4 行
.ell-3 {
  display: -webkit-box;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}
</style>

