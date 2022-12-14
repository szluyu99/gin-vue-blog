<script setup lang="ts">
import api from '@/api'

let loading = $ref(true)
let categoryList = $ref<any>([])

onMounted(async () => {
  const res = await api.getCategorys()
  categoryList = res.data.list
  loading = false
})
</script>

<template>
  <BannerPage :loading="loading" title="分类" label="category" card>
    <h2
      text-center leading-20 mb-6rem text-25 mt-10
      lg="text-36"
    >
      分类 - {{ categoryList.length }}
    </h2>
    <ul class="category-list" px-20>
      <li
        v-for="c of categoryList" :key="c.id"
        class="category-list-item"
        cursor-pointer hover:text-violet
      >
        <router-link :to="`categories/${c.id}?name=${c.name}`">
          {{ c.name }}
        </router-link>
        <span ml-5 text-4xl text-gray> ({{ c.article_count ?? 0 }}) </span>
      </li>
    </ul>
  </BannerPage>
</template>

<style lang="scss" scoped>
.category-list {
  margin: 0 1.8rem;
  list-style: none;

  .category-list-item {
    font-size: 4rem;
    padding: 8px 1.8rem 8px 0;
  }

  .category-list-item:before {
    display: inline-block;
    position: relative;
    left: -4rem;
    width: 12px;
    height: 12px;
    border: 0.9rem solid #49b1f5;
    border-radius: 50%;
    background: #fff;
    content: "";
    transition-duration: 0.3s;
  }

  .category-list-item:hover:before {
    border: 0.9rem solid #ff7242;
  }

  .category-list-item a:hover {
    transition: all 0.3s;
    color: #8e8cd8;
  }

  .category-list-item a:not(:hover) {
    transition: all 0.3s;
  }
}
</style>
