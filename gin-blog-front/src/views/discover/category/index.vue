<script setup>
import { onMounted, ref } from 'vue'

import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'

const loading = ref(true)
const categoryList = ref([])

onMounted(async () => {
  // 失败时 loading 也要复位, 否则页面卡在加载态;
  // data 为空不能让 categoryList 变成 null, 模板里直接取 .length
  try {
    const resp = await api.getCategorys()
    categoryList.value = resp.data ?? []
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
})
</script>

<template>
  <BannerPage :loading="loading" title="分类" label="category" card>
    <!-- 横幅上已经写着"分类"了, 这里只说明数量 -->
    <p class="text-center text-muted">
      共 {{ categoryList.length }} 个分类
    </p>
    <!-- 原来是单列列表: 970px 宽的卡片里放几行左对齐文字, 大片空白。
         改成响应式网格, 每项自成一张小卡片, 文章数用 badge 挂在右侧 -->
    <ul class="grid grid-cols-1 mt-6 gap-4 lg:grid-cols-3 sm:grid-cols-2">
      <li v-for="c of categoryList" :key="c.id">
        <RouterLink
          :to="`categories/${c.id}?name=${c.name}`"
          class="group flex items-center justify-between gap-3 rounded-xl bg-surface-soft px-4 py-3 shadow-sm transition-300 hover:shadow-md hover:-translate-y-0.5"
        >
          <span class="flex items-center gap-2 truncate">
            <span class="h-3 w-3 shrink-0 rounded-full bg-primary transition-300 group-hover:bg-#ff7242" />
            <span class="truncate text-lg group-hover:text-primary">{{ c.name }}</span>
          </span>
          <span class="shrink-0 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
            {{ c.article_count ?? 0 }}
          </span>
        </RouterLink>
      </li>
    </ul>

    <div v-if="!loading && !categoryList.length" class="py-10 text-center text-muted">
      还没有分类
    </div>
  </BannerPage>
</template>
