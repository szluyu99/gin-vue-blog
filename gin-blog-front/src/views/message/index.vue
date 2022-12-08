<script setup lang="ts">
import vueDanmaku from 'vue3-danmaku'

import api from '@/api'
import { convertImgUrl } from '@/utils'
import { useAppStore, useUserStore } from '@/store'
const userStore = useUserStore()
const { pageList } = storeToRefs(useAppStore())

let content = $ref('')
const showBtn = $ref(false)

const dmRef = $ref<any>(null) // 弹幕 ref 对象
const isHide = $ref(false) // 隐藏弹幕
const isLoop = $ref(false) // 循环播放

// 弹幕列表
let danmus = $ref([{
  avatar: 'https://static.talkxj.com/avatar/15bb3e8d4144163e006f4b2776aa7bac.jpg',
  content: '大家好，我是作者，欢迎给我点一颗 Star!',
  nickname: '阵、雨',
}])

onMounted(async () => {
  const res = await api.getMessages()
  danmus = [...danmus, ...res.data]
})

async function send() {
  content = content.trim()
  if (!content) {
    window?.$message?.info('消息不能为空!')
    return
  }
  const data = {
    avatar: userStore.avatar,
    nickname: userStore.nickname,
    content,
  }
  await api.saveMessage(data)
  dmRef.push(data)
  content = ''
}

// 注意 watch 第一个参数可以是函数, 由于用了 $ref, 直接写 isHide 不生效
watch(() => isHide, val => val ? dmRef.hide() : dmRef.show())

// 根据后端配置动态获取封面
const coverStyle = computed(() => {
  const page = pageList.value.find(e => e.label === 'message')
  return page
    ? `background: url('${page?.cover}') center center / cover no-repeat;`
    : 'background: url("https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg") center center / cover no-repeat;'
})
</script>

<template>
  <div
    :style="coverStyle"
    overflow-hidden
    absolute inset-x-0 h-screen
    class="banner-fade-down"
  >
    <!-- 弹幕输入框 -->
    <div
      animate-zoom-in
      top="3/10"
      border-1 rounded-10
      shadow-2xl
      absolute inset-x-1
      w-430
      text-center text-light
      px-5 py-25
      mx-auto z-5
    >
      <h1 text-32 font-bold>
        留言板
      </h1>
      <div flex justify-center mt-32 h-40>
        <input
          v-model="content"
          text-16 border-1 rounded-20 bg-transparent px-20 w="3/4" text="#eee"
          placeholder="说点什么吧？"
          @click.stop="showBtn = true"
          @keyup.enter="send"
        >
        <button
          v-if="showBtn"
          ml-12 text-16 rounded-20 px-18 border-1 animate-back-in-right
          @click="send"
        >
          发送
        </button>
      </div>
      <ul text-left text-white ml-20>
        <li mt-25 flex items-center>
          循环播放： <n-switch v-model:value="isLoop" ml-5 />
        </li>
        <li my-20>
          操作弹幕：
          <button ml-5 border-1 rounded-8 p-3 @click="dmRef.play">
            播放
          </button>
          <button mx-15 border-1 rounded-8 p-3 @click="dmRef.pause">
            暂停
          </button>
          <button border-1 rounded-8 p-3 @click="dmRef.stop">
            停止
          </button>
        </li>
        <li mb-15>
          隐藏弹幕：
          <n-switch v-model:value="isHide" ml-5 />
        </li>
      </ul>
    </div>
    <!-- 弹幕列表 -->
    <div absolute inset-0 top-60>
      <vue-danmaku
        ref="dmRef"
        v-model:danmus="danmus"
        wh-full
        use-slot
        :loop="isLoop"
        :speeds="200"
        :channels="0"
        :top="5"
        :is-suspend="true"
      >
        <template #dm="{ danmu }">
          <div bg="#00000080" rounded-20 text-white px-10 py-7 flex items-center>
            <n-avatar round size="small" :src="convertImgUrl(danmu.avatar)" mr-10 />
            <span vertical-middle> {{ `${danmu.nickname} : ${danmu.content}` }}</span>
          </div>
        </template>
      </vue-danmaku>
    </div>
  </div>
</template>

<style lang="scss" scoped>
input::-webkit-input-placeholder {
  color: #eee;
}
</style>
