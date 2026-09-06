<script setup>
import dayjs from 'dayjs'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'
import Comment from '@/components/comment/Comment.vue'
import { convertImgUrl } from '@/utils'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const talk = ref(null)

onMounted(async () => {
  try {
    const resp = await api.getTalk(route.params.id)
    talk.value = resp.data
  }
  catch (err) {
    // 私密说说在后端就当作不存在, 这里跟文章详情一样跳 404
    console.error(err)
    router.replace('/404')
  }
  finally {
    loading.value = false
  }
})
</script>

<template>
  <BannerPage label="talk" title="说说" card :loading="loading">
    <template v-if="talk">
      <div class="flex items-center gap-2">
        <img
          :src="convertImgUrl(talk.avatar)" :alt="talk.nickname"
          class="h-11 w-11 shrink-0 rounded-full bg-surface-soft object-cover"
        >
        <div class="min-w-0 flex-1">
          <p class="flex items-center gap-2">
            <span class="font-bold"> {{ talk.nickname || '博主' }} </span>
            <span v-if="talk.is_top" class="rounded bg-accent px-1.5 py-0.5 text-xs text-white">置顶</span>
          </p>
          <p class="text-xs color-muted">
            {{ dayjs(talk.created_at).format('YYYY-MM-DD HH:mm') }}
          </p>
        </div>
      </div>

      <!-- 说说是纯文本, 保留用户敲的换行 -->
      <p class="my-5 whitespace-pre-wrap leading-relaxed">
        {{ talk.content }}
      </p>

      <RouterLink to="/talks" class="text-sm text-primary">
        ← 返回说说列表
      </RouterLink>

      <!-- 评论: type 3 就是说说, topic_id 由路由参数带过去 -->
      <Comment class="mt-16" :type="3" />
    </template>
  </BannerPage>
</template>
