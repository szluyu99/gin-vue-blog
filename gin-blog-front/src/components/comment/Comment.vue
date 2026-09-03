<script setup>
import dayjs from 'dayjs'
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'

import api from '@/api'

// 分页
import ULoading from '@/components/ui/ULoading.vue'

import { useAppStore, useUserStore } from '@/store'

import { convertImgUrl } from '@/utils'
import CommentField from './CommentField.vue'
// 评论 / 回复 框
import Paging from './Paging.vue'

const { type } = defineProps({
  // 评论类型: 1-文章, 2-友链, 3-说说
  type: Number,
})

const [userStore, appStore] = [useUserStore(), useAppStore()]

onMounted(() => {
  getComments()
})

// url 中存在 id 参数则为 topic_id, 否则为 0
const topicId = +(useRoute().params.id ?? 0)

// 加载评论
const commentList = ref([]) // 评论列表 (分页加载)
const commentCount = ref(0) // 评论总数量
const listLoading = ref(false) // 列表加载状态
const params = reactive({ type, page_size: 10, page_num: 1, topic_id: topicId }) // 加载评论的参数

// 延时定时器: 卸载时清掉, 否则会在组件销毁后改状态
let loadTimer = null
onUnmounted(() => clearTimeout(loadTimer))

async function getComments() {
  listLoading.value = true
  try {
    const resp = await api.getComments(params)

    // * 全局加载更多, 0.8s 延时
    loadTimer = setTimeout(() => {
      params.page_num === 1
        ? commentList.value = resp.data.page_data
        : commentList.value.push(...resp.data.page_data)
      commentCount.value = resp.data.total
      params.page_num++
      listLoading.value = false
    }, 800)
  }
  catch (err) {
    // 不复位 listLoading 的话, "点击加载更多" (v-if="!listLoading") 会永久消失
    listLoading.value = false
    console.error(err)
  }
}

// 重新加载评论(提交评论以后)
function reloadComments() {
  params.page_num = 1 // 页数重置
  getComments()
}

// 回复相关的状态一律按评论 id 记, 不再按 v-for 下标访问模板 ref:
// Vue 不保证模板 ref 数组的顺序与源数组一致, 新增评论后刷新列表就会点到别人的回复框
// (原来靠 watch commentList + nextTick 整体重建列表来打补丁, 现在不需要了)
const activeReply = ref(null) // 当前打开的回复框 { commentId, nickname, replyUserId, parentId }
const expandedIds = reactive(new Set()) // 已经点过"查看更多回复"的评论 id
const replyPages = reactive({}) // 评论 id -> 回复列表当前页

// 回复评论: target 可能是评论本身, 也可能是某条回复
function replyComment(comment, target) {
  activeReply.value = {
    commentId: comment.id,
    // 原来取的是 obj.nickname, 而接口返回的昵称在 user.info.nickname 下,
    // 所以回复框一直拿不到昵称, 既没有 "回复 @xxx" 提示也没有取消按钮
    nickname: target.user?.info?.nickname ?? '',
    replyUserId: target.user_id,
    parentId: comment.id,
  }
}

// 提交回复后, 重新加载该评论的回复
async function reloadReplies(comment) {
  try {
    const { data } = await api.getCommentReplies(
      comment.id,
      { page_size: 5, page_num: replyPages[comment.id] ?? 1 },
    )
    // * 局部更新某个评论的回复
    comment.reply_list = data
    comment.reply_count++ // 数量 + 1
    expandedIds.add(comment.id)
  }
  catch (err) {
    console.error(err)
  }
}

// "点击查看" 显示更多回复
async function checkReplies(comment) {
  try {
    // 查第一页 (5 条数据)
    const { data } = await api.getCommentReplies(comment.id, { page_num: 1, page_size: 5 })
    comment.reply_list = data
    replyPages[comment.id] = 1
    expandedIds.add(comment.id)
  }
  catch (err) {
    console.error(err)
  }
}

// 修改回复分页中当前页数
async function changeReplyCurrent(comment, pageNum) {
  try {
    const { data } = await api.getCommentReplies(comment.id, { page_num: pageNum, page_size: 5 })
    comment.reply_list = data
    replyPages[comment.id] = pageNum
  }
  catch (err) {
    console.error(err)
  }
}

// TODO: 点赞
async function likeComment(comment) {
  // 判断是否登录
  if (!userStore.userId) {
    appStore.setLoginFlag(true)
    return
  }

  try {
    await api.saveLikeComment(comment.id)
    // 判断是否点赞
    if (userStore.commentLikeSet.includes(comment.id)) {
      comment.like_count--
      window.$message?.info('已取消')
    }
    else {
      comment.like_count++
      window.$message?.success('已点赞')
    }
    // 维护全局状态中的点赞 Set
    userStore.commentLike(comment.id)
  }
  catch (err) {
    console.error(err)
  }
}

// 判断当前用户是否点赞过该评论
const isLike = computed(() => id => userStore.commentLikeSet.includes(id))
</script>

<template>
  <div>
    <p class="flex items-center text-xl font-bold">
      <span class="i-fa:comments mr-3 text-blue" /> 评论
    </p>
    <!-- 评论框 -->
    <CommentField
      :show="true"
      :type="type"
      :topic-id="topicId"
      @after-submit="reloadComments"
    />
    <!-- 评论详情 -->
    <div v-if="commentCount">
      <!-- 评论数量 -->
      <p class="mb-4 mt-7 flex items-center text-xl font-bold">
        <span> {{ commentCount }} 评论 </span>
        <span
          class="i-uiw:reload ml-4 cursor-pointer text-base"
          :class="listLoading ? 'animate-spin' : ''"
          @click="reloadComments"
        />
      </p>
      <!-- 评论列表 -->
      <div v-for="(comment, idx) of commentList" :key="comment.id" class="my-1 flex">
        <img :src="convertImgUrl(comment.user?.info?.avatar)" class="h-[40px] w-[40px] duration-600 hover:rotate-360" loading="lazy">
        <div class="ml-3 flex flex-1 flex-col">
          <!-- 评论人名称: 根据是否有 website 显示不同效果 -->
          <div>
            <span v-if="!comment.user?.info?.website" class="text-sm">
              {{ comment.user?.info?.nickname }}
            </span>
            <a v-else :href="comment.user?.info?.website" target="_blank" class="color-[#1abc9c] font-500 transition-300">
              {{ comment.user?.info?.nickname }}
            </a>
            <!-- TODO: 博主标记 -->
            <!-- <span v-if="comment.user_id === 10" class="ml-2 inline-block rounded-3 bg-#ffa51e px-6 py-1 text-xs color-#fff">
              博主
            </span> -->
          </div>
          <!-- 楼层 + 时间 + 点赞 + 回复按钮 -->
          <div class="flex justify-between text-sm">
            <div class="flex items-center gap-2 py-1 color-muted">
              <span> {{ commentCount - idx }}楼 </span>
              <span> {{ dayjs(comment.created_at).format('YYYY-MM-DD') }} </span>
              <button
                class="i-mdi:thumb-up hover-bg-red"
                :class="isLike(comment.id) ? 'bg-red' : ''"
                @click="likeComment(comment)"
              />
              <span v-show="comment.like_count"> {{ comment.like_count }} </span>
            </div>
            <button class="color-#ef2f11" @click="replyComment(comment, comment)">
              回复
            </button>
          </div>
          <!-- 评论内容 -->
          <div class="my-1" v-html="comment.content" />
          <!-- 评论回复 start -->
          <div v-for="reply of comment.reply_list" :key="reply.id" class="mt-2 flex">
            <img :src="convertImgUrl(reply.user?.info?.avatar)" class="h-[40px] w-[40px] duration-600 hover:rotate-360" loading="lazy">
            <div class="ml-2 flex flex-1 flex-col">
              <!-- 回复人名称 -->
              <div>
                <!-- 根据是否有 website 显示不同效果 -->
                <span v-if="!reply.user?.info?.website" class="text-sm">
                  {{ reply.user?.info.nickname }}
                </span>
                <a v-else :href="reply.user?.info?.website" target="_blank" class="color-#1abc9c font-500 transition-300">
                  {{ reply.user?.info.nickname }}
                </a>
                <!-- TODO: 博主标记 -->
                <!-- <span v-if="reply.user_id === 10" class="ml-6 inline-block rounded-3 bg-#ffa51e px-6 py-1 text-sm color-#fff">
                  博主
                </span> -->
              </div>
              <!-- 时间 + 点赞 + 回复按钮 -->
              <div class="flex justify-between text-sm">
                <div class="flex items-center gap-2 py-1 color-muted">
                  <span> {{ dayjs(reply.created_at).format('YYYY-MM-DD') }} </span>
                  <button
                    class="i-mdi:thumb-up hover-bg-red"
                    :class="isLike(reply.id) ? 'bg-red' : ''"
                    @click="likeComment(reply)"
                  />
                  <span v-show="reply.like_count"> {{ reply.like_count }} </span>
                </div>
                <button class="color-#ef2f11" @click="replyComment(comment, reply)">
                  回复
                </button>
              </div>
              <!-- 回复内容 -->
              <div>
                <!-- 回复用户名: 自己回复自己不显示 "@名称" -->
                <template v-if="reply.user_id !== comment.user_id">
                  <a v-if="reply.user?.info?.website" :href="reply.reply_website" target="_blank">
                    @{{ reply.user?.info?.nickname }}
                  </a>
                  <span v-else>
                    @{{ reply.user?.info?.nickname }}
                  </span>，
                </template>
                <span class="my-3" v-html="reply.content" />
              </div>
            </div>
          </div>
          <!-- 评论回复 end -->

          <!-- 回复数量 -->
          <div
            v-show="comment.reply_count > 3 && !expandedIds.has(comment.id)"
            class="mt-4 text-[13px] color-muted"
          >
            共 <b> {{ comment.reply_count }} </b>  条回复
            <button class="color-#00a1d6" @click="checkReplies(comment)">
              ，点击查看
            </button>
          </div>
          <!-- 回复分页: 展开后且回复超过 5 条才有分页 -->
          <Paging
            v-if="expandedIds.has(comment.id) && comment.reply_count > 5"
            :page-total="Math.ceil(comment.reply_count / 5)"
            :current="replyPages[comment.id] ?? 1"
            @change-current="page => changeReplyCurrent(comment, page)"
          />
          <!-- 回复框: 同一时刻只打开一个 -->
          <CommentField
            v-if="activeReply?.commentId === comment.id"
            :show="true"
            :type="type"
            :topic-id="topicId"
            :nickname="activeReply.nickname"
            :reply-user-id="activeReply.replyUserId"
            :parent-id="activeReply.parentId"
            @cancel="activeReply = null"
            @after-submit="reloadReplies(comment)"
          />
          <!-- 分隔线: 注意最后一个评论没有线 (比的是已加载条数, 不是总数) -->
          <div v-if="(idx + 1) !== commentList.length" class="my-2.5 h-0.5 bg-light-500" />
        </div>
      </div>
      <!-- 加载更多 -->
      <div class="m-4 f-c-c">
        <button
          v-if="commentCount > commentList.length && !listLoading"
          text @click="getComments"
        >
          点击加载更多...
        </button>
        <ULoading :show="listLoading" />
      </div>
    </div>
    <!-- 没有评论的提示 -->
    <div v-else class="mb-10 mt-30 text-center text-zinc">
      暂无评论，来发评论吧~
    </div>
  </div>
</template>
