<script setup>
import { useUserStore } from '@/store'
import api from '@/api'
const { nickname, avatar } = useUserStore()

let homeInfo = $ref({})
onMounted(async () => {
  getOneSentence()
  const res = await api.getHomeInfo()
  homeInfo = res.data
})

// 一言
let sentence = $ref('')
async function getOneSentence() {
  fetch('https://v1.hitokoto.cn?c=i')
    .then(res => res.json())
    .then(data => sentence = data.hitokoto)
    .catch(() => sentence = '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。')
}
</script>

<template>
  <AppPage :show-footer="true">
    <div flex-1>
      <n-card rounded-10>
        <div flex items-center>
          <n-avatar round :size="60" :src="avatar" />
          <div ml-20>
            <p text-16>
              Hello, {{ nickname }}
            </p>
            <n-gradient-text
              mt-5 text-12 op-60
              gradient="linear-gradient(90deg, red 0%, green 50%, blue 100%)"
            >
              {{ sentence }}
            </n-gradient-text>
          </div>
          <div ml-auto flex items-center>
            <!-- <n-statistic label="待办" :value="4">
              <template #suffix>
                / 10
              </template>
            </n-statistic> -->
            <n-statistic label="Stars" w-80 ml-80>
              <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                <img
                  alt="stars"
                  src="https://badgen.net/github/stars/szluyu99/gin-vue-blog"
                >
              </a>
            </n-statistic>
            <n-statistic label="Forks" w-100 ml-80>
              <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                <img
                  alt="forks"
                  src="https://badgen.net/github/forks/szluyu99/gin-vue-blog"
                >
              </a>
            </n-statistic>
          </div>
        </div>
      </n-card>

      <n-grid mt-15 x-gap="12" :cols="4">
        <n-gi>
          <n-card>
            <TheIcon icon="fa6-solid:users" color="#40C9C6" :size="60" />
            <n-statistic float-right label="访问量">
              {{ homeInfo.view_count ?? 'unknown' }}
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <TheIcon icon="heroicons:users-solid" color="#34BFA3" :size="60" />
            <n-statistic float-right label="用户量">
              {{ homeInfo.user_count ?? 'unknown' }}
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <TheIcon icon="material-symbols:article" color="#F4516C" :size="60" />
            <n-statistic float-right label="文章量">
              {{ homeInfo.article_count ?? 'unknown' }}
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <TheIcon icon="bxs:comment-dots" color="#36A3F7" :size="60" />
            <n-statistic float-right label="留言量">
              {{ homeInfo.message_count ?? 'unknown' }}
            </n-statistic>
          </n-card>
        </n-gi>
      </n-grid>

      <!-- TODO: 完善首页设计 -->
      <n-card title="项目" size="small" mt-15 rounded-10>
        <template #header-extra>
          <n-button text type="primary">
            更多
          </n-button>
        </template>
        <div flex flex-wrap justify-between>
          <n-card
            v-for="i in 10"
            :key="i"
            class="w-300 flex-shrink-0 my-10 cursor-pointer"
            hover:card-shadow
            title="Gin Blog Admin"
            size="small"
          >
            <p op-60>
              这是个基于 gin 开发的博客管理后台
            </p>
          </n-card>
          <div w-300 h-0 />
          <div w-300 h-0 />
          <div w-300 h-0 />
          <div w-300 h-0 />
        </div>
      </n-card>
    </div>
  </AppPage>
</template>
