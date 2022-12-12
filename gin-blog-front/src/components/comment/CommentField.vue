<script setup lang="ts">
import { convertImgUrl } from '@/utils'
import api from '@/api'

import { useAppStore, useUserStore } from '@/store'

const props = defineProps<{
  type: number // 评论类型 1-文章 2-友链
  show: boolean // 默认是否显示
  topicId?: number // 主题 id
}>()

const emit = defineEmits(['afterSubmit']) // 调用父方法
const [userStore, appStore] = [useUserStore(), useAppStore()]

let show = $ref(props.show) // 是否显示
const data = reactive({
  nickname: '', // * 回复用户, 不为空则说明是回复框
  content: '', // 回复内容
  topic_id: props.topicId ?? 0, // 主题 id
  reply_user_id: 0, // 回复用户 id
  parent_id: 0, // 父评论 id
  type: props.type,
})

// 判断是回复还是评论: 存在 nickname 则是回复
const isReply = $computed(() => !!data.nickname)

// 取消评论
function setReply(flag: boolean) {
  show = flag
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

    isReply && setReply(false)
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
  <div v-if="show" flex p-10 mt-15 border-1px border-solid border-rounded-1rem border-color="#90939950">
    <n-avatar
      round :size="36"
      :src="convertImgUrl(userStore.avatar)"
    />
    <div w-full my-5 ml-12>
      <textarea
        v-model="data.content"
        :placeholder="placeholderText"
        rows="5"
        text-16 w-full
        bg-light-400
      />
      <div flex justify-between>
        <!-- TODO: 表情框 -->
        <i-mdi:emoticon-happy-outline text-24 text-orange cursor-pointer @click="chooseEmoji" />
        <div>
          <span
            v-if="data.nickname"
            btn bg-bluegray hover:bg-bluegray mr-15
            @click="setReply(false)"
          > 取消 </span>
          <span btn @click="submitComment"> 提交 </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
textarea {
  background: url(https://static.talkxj.com/config/commentBack.webp) 100% 100% no-repeat;
  resize: none;
}
</style>
