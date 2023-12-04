<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import vueDanmaku from 'vue3-danmaku'
import USwitch from '@/components/ui/USwitch.vue'

import api from '@/api'
import { convertImgUrl } from '@/utils'
import { useAppStore, useUserStore } from '@/store'
const userStore = useUserStore()
const { pageList } = storeToRefs(useAppStore())

const content = ref('')
const showBtn = ref(false)

const dmRef = ref(null) // 弹幕 ref 对象
const isHide = ref(false) // 隐藏弹幕
const isLoop = ref(false) // 循环播放

// 弹幕列表
const danmus = ref([{
  avatar: 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png',
  content: '大家好，我是作者，欢迎给我点一颗 Star!',
  nickname: '阵、雨',
}])

onMounted(async () => {
  const res = await api.getMessages()
  danmus.value = [...danmus.value, ...res.data]
})

async function send() {
  content.value = content.value.trim()
  if (!content.value) {
    window?.$message?.info('消息不能为空!')
    return
  }
  const data = {
    avatar: userStore.avatar,
    nickname: userStore.nickname,
    content,
  }
  await api.saveMessage(data)
  dmRef.value.push(data)
  content.value = ''
}

watch(isHide, val => val ? dmRef.value.hide() : dmRef.value.show())

// 根据后端配置动态获取封面
const coverStyle = computed(() => {
  const page = pageList.value.find(e => e.label === 'message')
  return page
    ? `background: url('${page?.cover}') center center / cover no-repeat;`
    : 'background: url("https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg") center center / cover no-repeat;'
})
</script>

<template>
  <div :style="coverStyle" class="banner-fade-down absolute inset-x-0 h-screen overflow-hidden">
    <!-- 弹幕输入框 -->
    <div class="absolute inset-x-1 top-3/10 z-5 mx-auto w-350 animate-zoom-in border-1 rounded-10 px-5 py-15 text-center text-light shadow-2xl lg:w-420 lg:py-20">
      <h1 class="text-25 font-bold lg:text-30">
        留言板
      </h1>
      <div class="mt-20 h-36 flex justify-center lg:mt-25 lg:h-40">
        <input
          v-model="content"
          class="w-3/4 border-1 rounded-20 bg-transparent px-20 text-#eee"
          placeholder="说点什么吧？"
          @click.stop="showBtn = true"
          @keyup.enter="send"
        >
        <button
          v-if="showBtn"
          class="ml-12 animate-back-in-right border-1 rounded-20 px-18"
          @click="send"
        >
          发送
        </button>
      </div>
      <ul class="ml-20 text-left text-white">
        <li class="mt-25 flex items-center">
          循环播放： <USwitch v-model="isLoop" />
        </li>
        <li class="my-20">
          操作弹幕：
          <button class="ml-5 border-1 rounded-8 p-3" @click="dmRef.play">
            播放
          </button>
          <button class="mx-15 border-1 rounded-8 p-3" @click="dmRef.pause">
            暂停
          </button>
          <button class="border-1 rounded-8 p-3" @click="dmRef.stop">
            停止
          </button>
        </li>
        <li class="mb-15 flex items-center">
          隐藏弹幕： <USwitch v-model="isHide" class="ml-5" />
        </li>
      </ul>
    </div>
    <!-- 弹幕列表 -->
    <div class="absolute inset-0 top-60">
      <vue-danmaku
        ref="dmRef"
        v-model:danmus="danmus"
        class="h-full w-full"
        use-slot
        :loop="isLoop"
        :speeds="200"
        :channels="0"
        :top="5"
        :is-suspend="true"
      >
        <template #dm="{ danmu }">
          <div class="flex items-center rounded-20 bg-#00000060 px-6 py-4 text-15 text-white lg:px-8 lg:py-5 lg:text-16">
            <img class="h-28 rounded-full" :src="convertImgUrl(danmu.avatar)" alt="danmu avatar">
            <span class="ml-10"> {{ `${danmu.nickname} : ${danmu.content}` }}</span>
          </div>
        </template>
      </vue-danmaku>
    </div>
  </div>
</template>

<style scoped>
input::-webkit-input-placeholder {
  color: #eee;
}
</style>
