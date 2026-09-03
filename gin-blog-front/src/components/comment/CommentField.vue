<script setup>
import { computed, reactive, watch } from 'vue'

import api from '@/api'
import { useAppStore, useUserStore } from '@/store'
import { convertImgUrl } from '@/utils'

const props = defineProps({
  // 评论类型 1-文章 2-友链
  type: Number,
  // 是否显示
  show: Boolean,
  // 主题 id
  topicId: Number,
  // 以下三项不为空时是"回复框": 以前由父组件通过模板 ref 直接改内部的 data,
  // 按 v-for 下标取 ref 会改到别人的回复框, 改成入参传进来
  nickname: { type: String, default: '' },
  replyUserId: { type: Number, default: 0 },
  parentId: { type: Number, default: 0 },
})

const emit = defineEmits(['afterSubmit', 'cancel']) // 调用父方法
const [userStore, appStore] = [useUserStore(), useAppStore()]

const data = reactive({
  nickname: props.nickname, // * 回复用户, 不为空则说明是回复框
  content: '', // 回复内容
  topic_id: props.topicId ?? 0, // 主题 id
  reply_user_id: props.replyUserId, // 回复用户 id
  parent_id: props.parentId, // 父评论 id
  type: props.type,
})

// 判断是回复还是评论: 存在 nickname 则是回复
const isReply = computed(() => !!data.nickname)

// 同一条评论下切换回复对象时组件不会重建, 这里把新的回复对象同步进 data
// (只同步回复对象, 不清空已经写的内容)
watch(() => [props.nickname, props.replyUserId, props.parentId], ([nickname, replyUserId, parentId]) => {
  data.nickname = nickname
  data.reply_user_id = replyUserId
  data.parent_id = parentId
})

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

    // 先通知父组件刷新, 再让它关掉回复框 (关掉会卸载本组件)
    emit('afterSubmit')
    isReply.value && emit('cancel')
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
            @click="emit('cancel')"
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
