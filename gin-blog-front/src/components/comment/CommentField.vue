<script setup>
import { computed, reactive, ref } from 'vue'

import { convertImgUrl } from '@/utils'
import { useAppStore, useUserStore } from '@/store'
import api from '@/api'

const props = defineProps({
  // 评论类型 1-文章 2-友链
  type: Number,
  // 默认是否显示
  show: Boolean,
  // 主题 id
  topicId: Number,
})

const emit = defineEmits(['afterSubmit']) // 调用父方法
const [userStore, appStore] = [useUserStore(), useAppStore()]

const show = ref(props.show) // 是否显示
const data = reactive({
  nickname: '', // * 回复用户, 不为空则说明是回复框
  content: '', // 回复内容
  topic_id: props.topicId ?? 0, // 主题 id
  reply_user_id: 0, // 回复用户 id
  parent_id: 0, // 父评论 id
  type: props.type,
})

// 判断是回复还是评论: 存在 nickname 则是回复
const isReply = computed(() => !!data.nickname)

// 取消评论
function setReply(flag) {
  show.value = flag
}

// 提交评论
async function submitComment() {
  // 判断是否登录
  if (!userStore.userId) {
    appStore.setLoginFlag(true)
    return
  }
  // 判断内容不为空
  if (!data.content.trim()) {
    window.$message?.error('评论内容不能为空')
    return
  }

  // TODO: 解析表情

  // 调用接口
  try {
    await api.saveComment(data)
    window.$message?.info('评论成功')
    data.content = ''

    isReply.value && setReply(false)
    emit('afterSubmit') // 提交后调用刷新方法
  }
  catch (err) {
    console.error(err)
  }
}

// TODO: 表情框
function chooseEmoji() {
  window.$message?.info('表情选择正在开发中...')
}

// 输入框提示语
const placeholderText = computed(() =>
  (data.nickname ? `回复 @${data.nickname}：` : '留下点什么吧...'),
)

// 暴露给父组件, 可以在父组件中访问并修改
defineExpose({ data, setReply })
</script>

<template>
  <div v-if="show" class="mt-4 flex border-1 border-color-#90939950 border-rounded-1rem border-solid p-2">
    <img class="h-9 w-9" :src="convertImgUrl(userStore.avatar)">
    <div class="my-1 ml-3 w-full">
      <textarea
        v-model="data.content"
        :placeholder="placeholderText"
        rows="5"
        class="w-full rounded bg-light-400 p-2 outline-none"
      />
      <div class="flex justify-between">
        <!-- TODO: 表情框 -->
        <span class="i-mdi:emoticon-happy-outline cursor-pointer text-xl text-orange" @click="chooseEmoji" />
        <div>
          <span
            v-if="data.nickname"
            class="the-button mr-4 bg-bluegray hover:bg-bluegray"
            @click="setReply(false)"
          >
            取消
          </span>
          <span class="the-button" @click="submitComment"> 提交 </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
textarea {
  /* background: url(https://w.wallhaven.cc/full/1p/wallhaven-1poo61.jpg) 100% 100% no-repeat; */
  background-color:ghostwhite;
  resize: none;
}
</style>
