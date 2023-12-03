<script setup>
import { onMounted, ref } from 'vue'

import BannerPage from '@/components/page/BannerPage.vue'
import api from '@/api'

const loading = ref(true)
const categoryList = ref([])

onMounted(async () => {
  const res = await api.getCategorys()
  categoryList.value = res.data.list
  loading.value = false
})
</script>

<template>
  <BannerPage :loading="loading" title="分类" label="category" card>
    <h2 class="mb-6rem mt-10 text-center text-25 leading-20 lg:text-36">
      分类 - {{ categoryList.length }}
    </h2>
    <ul class="category-list px-20">
      <li
        v-for="c of categoryList" :key="c.id"
        class="group mx-1.8 cursor-pointer py-4 duration-300 hover:text-violet"
      >
        <RouterLink :to="`categories/${c.id}?name=${c.name}`">
          <div class="flex items-center">
            <span class="mr-8 inline-block h-15 w-15 rounded-full bg-[#49b1f5] group-hover:bg-[#ff7242]" />
            <span>
              <span class="text-16"> {{ c.name }} </span>
              <span class="ml-5 text-5xl text-gray"> ({{ c.article_count ?? 0 }}) </span>
            </span>
          </div>
        </RouterLink>
      </li>
    </ul>
  </BannerPage>
</template>
