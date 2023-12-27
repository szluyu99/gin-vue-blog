<script setup>
import { onMounted, ref } from 'vue'

import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'

const loading = ref(true)
const categoryList = ref([])

onMounted(async () => {
  const resp = await api.getCategorys()
  categoryList.value = resp.data
  loading.value = false
})
</script>

<template>
  <BannerPage :loading="loading" title="分类" label="category" card>
    <h2 class="text-center text-2xl leading-8 lg:text-3xl">
      分类 - {{ categoryList.length }}
    </h2>
    <ul class="mt-4 px-5 space-y-2">
      <li
        v-for="c of categoryList" :key="c.id"
        class="group cursor-pointer duration-300 hover:text-violet"
      >
        <RouterLink :to="`categories/${c.id}?name=${c.name}`">
          <div class="flex items-center">
            <span class="mr-2 inline-block h-4 w-4 rounded-full bg-[#49b1f5] group-hover:bg-[#ff7242]" />
            <div>
              <span class="text-lg"> {{ c.name }} </span>
              <span class="ml-1 text-sm text-gray"> ({{ c.article_count ?? 0 }}) </span>
            </div>
          </div>
        </RouterLink>
      </li>
    </ul>
  </BannerPage>
</template>
