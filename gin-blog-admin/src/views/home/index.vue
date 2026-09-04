<script setup>
import { NAvatar, NButton, NCard, NGi, NGradientText, NGrid, NStatistic } from 'naive-ui'
import { onMounted, ref } from 'vue'

import api from '@/api'
import AppPage from '@/components/common/AppPage.vue'
import { useUserStore } from '@/store'

// 不解构 store: 解构后失去响应性, getUserInfo() 回来后这里的昵称和头像不会更新
const userStore = useUserStore()

const homeInfo = ref({
  view_count: 0,
  user_count: 0,
  article_count: 0,
  message_count: 0,
})

onMounted(async () => {
  getOneSentence()
  // 裸 await 会在接口失败时留下未捕获的 rejection;
  // data 为空时也不能让 homeInfo 变成 null, 模板里要取它的字段
  try {
    const res = await api.getHomeInfo()
    if (res.data) {
      homeInfo.value = res.data
    }
  }
  catch (err) {
    console.error(err)
  }
})

// 一言: 每次进后台看到一句新的文案
// 接口在部分网络下不通, 兜底从内置文案里随机取一句, 而不是固定显示同一句
const FALLBACK_SENTENCES = [
  '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。',
  '书山有路勤为径，学海无涯苦作舟。',
  '纸上得来终觉浅，绝知此事要躬行。',
  '不积跬步，无以至千里；不积小流，无以成江海。',
  '路漫漫其修远兮，吾将上下而求索。',
  '业不可不勤，勤则百弊自去。',
  '博观而约取，厚积而薄发。',
]
function randomSentence() {
  return FALLBACK_SENTENCES[Math.floor(Math.random() * FALLBACK_SENTENCES.length)]
}

const sentence = ref('')
async function getOneSentence() {
  try {
    const resp = await fetch('https://v1.hitokoto.cn?c=i')
    const data = await resp.json()
    sentence.value = data?.hitokoto || randomSentence()
  }
  catch {
    sentence.value = randomSentence()
  }
}
</script>

<template>
  <AppPage>
    <div class="flex-1">
      <NCard>
        <div class="flex items-center">
          <NAvatar round :size="60" :src="userStore.avatar" />
          <div class="ml-5">
            <p> Hello, {{ userStore.nickname }} </p>
            <NGradientText class="mt-1 op-60" gradient="linear-gradient(90deg, red 0%, green 50%, blue 100%)">
              {{ sentence }}
            </NGradientText>
          </div>
          <div class="ml-auto flex items-center">
            <NStatistic label="Stars" class="w-[80px]">
              <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                <img
                  alt="stars"
                  src="https://badgen.net/github/stars/szluyu99/gin-vue-blog"
                >
              </a>
            </NStatistic>
            <NStatistic label="Forks" class="ml-10 w-[100px]">
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

      <NGrid class="mt-4" x-gap="12" :cols="4">
        <template
          v-for="item of [
            { icon: 'i-fa6-solid:users', color: 'text-[#40C9C6]', label: '访问量', key: 'view_count' },
            { icon: 'i-heroicons:users-solid', color: 'text-[#34BFA3]', label: '用户量', key: 'user_count' },
            { icon: 'i-material-symbols:article', color: 'text-[#F4516C]', label: '文章量', key: 'article_count' },
            { icon: 'i-bxs:comment-dots', color: 'text-[#36A3F7]', label: '留言量', key: 'message_count' },
          ]" :key="item.key"
        >
          <NGi>
            <NCard>
              <span
                class="text-[60px]"
                :class="[item.icon, item.color]"
              />
              <NStatistic class="float-right" :label="item.label">
                {{ homeInfo[item.key] ?? 'unknown' }}
              </NStatistic>
            </NCard>
          </NGi>
        </template>
      </NGrid>

      <!-- TODO: 完善首页设计 -->
      <NCard title="项目" size="small" class="mt-4">
        <template #header-extra>
          <NButton text type="primary">
            更多
          </NButton>
        </template>
        <NCard
          v-for="i in 5" :key="i"
          class="my-2 w-[300px] flex-shrink-0 cursor-pointer hover:shadow-lg"
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
