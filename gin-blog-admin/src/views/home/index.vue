<script setup>
import { onMounted, ref } from 'vue'
import { NAvatar, NButton, NCard, NGi, NGradientText, NGrid, NStatistic } from 'naive-ui'

import AppPage from '@/components/page/AppPage.vue'
import { useUserStore } from '@/store'
import api from '@/api'

const { nickname, avatar } = useUserStore()

const homeInfo = ref({})
onMounted(async () => {
  getOneSentence()
  const res = await api.getHomeInfo()
  homeInfo.value = res.data
})

// 一言
const sentence = ref('')
async function getOneSentence() {
  fetch('https://v1.hitokoto.cn?c=i')
    .then(res => res.json())
    .then(data => sentence.value = data.hitokoto)
    .catch(() => sentence.value = '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。')
}
</script>

<template>
  <AppPage :show-footer="true">
    <div class="flex-1">
      <NCard class="rounded-10">
        <div class="flex items-center">
          <NAvatar round :size="60" :src="avatar" />
          <div class="ml-20">
            <p class="text-16">
              Hello, {{ nickname }}
            </p>
            <NGradientText
              class="mt-5 text-12 op-60"
              gradient="linear-gradient(90deg, red 0%, green 50%, blue 100%)"
            >
              {{ sentence }}
            </NGradientText>
          </div>
          <div class="ml-auto flex items-center">
            <NStatistic label="Stars" class="ml-80 w-80">
              <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                <img
                  alt="stars"
                  src="https://badgen.net/github/stars/szluyu99/gin-vue-blog"
                >
              </a>
            </NStatistic>
            <NStatistic label="Forks" class="ml-80 w-100">
              <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                <img
                  alt="forks"
                  src="https://badgen.net/github/forks/szluyu99/gin-vue-blog"
                >
              </a>
            </NStatistic>
          </div>
        </div>
      </NCard>

      <NGrid class="mt-15" x-gap="12" :cols="4">
        <NGi>
          <NCard>
            <span class="i-fa6-solid:users text-60 text-#40C9C6" />
            <NStatistic class="float-right" label="访问量">
              {{ homeInfo.view_count ?? 'unknown' }}
            </NStatistic>
          </NCard>
        </NGi>
        <NGi>
          <NCard>
            <span class="i-heroicons:users-solid text-60 text-#34BFA3" />
            <NStatistic class="float-right" label="用户量">
              {{ homeInfo.user_count ?? 'unknown' }}
            </NStatistic>
          </NCard>
        </NGi>
        <NGi>
          <NCard>
            <span class="i-material-symbols:article text-60 text-#F4516C" />
            <NStatistic class="float-right" label="文章量">
              {{ homeInfo.article_count ?? 'unknown' }}
            </NStatistic>
          </NCard>
        </NGi>
        <NGi>
          <NCard>
            <span class="i-bxs:comment-dots text-60 text-#36A3F7" />
            <NStatistic class="float-right" label="留言量">
              {{ homeInfo.message_count ?? 'unknown' }}
            </NStatistic>
          </NCard>
        </NGi>
      </NGrid>

      <!-- TODO: 完善首页设计 -->
      <NCard title="项目" size="small" class="mt-15 rounded-10">
        <template #header-extra>
          <NButton text type="primary">
            更多
          </NButton>
        </template>
        <NCard
          v-for="i in 5" :key="i"
          class="my-10 w-300 flex-shrink-0 cursor-pointer hover:shadow-lg"
          title="Gin Blog Admin"
          size="small"
        >
          <p class="op-60">
            这是个基于 gin 开发的博客管理后台
          </p>
        </NCard>
      </NCard>
    </div>
  </AppPage>
</template>
