<script setup>
import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

// import EmojiList from '@/assets/js/emoji'
import CommentField from './CommentField.vue' // 评论 / 回复 框
import Paging from './Paging.vue' // 分页
import ULoading from '@/components/ui/ULoading.vue'

import { convertImgUrl, formatDate } from '@/utils'
import { useAppStore, useUserStore } from '@/store'
import api from '@/api'

const { type } = defineProps({
  // 评论类型: 1-文章, 2-友链
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

async function getComments() {
  listLoading.value = true
  try {
    const res = await api.getComments(params)
    // * 全局加载更多, 0.8s 延时
    setTimeout(() => {
      params.page_num === 1
        ? commentList.value = res.data.pageData
        : commentList.value.push(...res.data.pageData)
      commentCount.value = res.data.total
      params.page_num++
      listLoading.value = false
    }, 800)
  }
  catch (err) {
    console.error(err)
  }
}
// 重新加载评论(提交评论以后)
function reloadComments() {
  params.page_num = 1 // 页数重置
  getComments()
}

// * 解决新增评论后刷新数据, 点击回复的顺序错乱问题
const refresh = ref(true) // 重新刷新整个评论列表
watch(commentList, () => {
  refresh.value = false
  nextTick(() => {
    refresh.value = true
  })
}, { deep: false }) // deep = false 防止 "查看更多" 时刷新整个数据

// 回复相关
// ! 可以获取 v-for 循环中的 DOM 数组
const replyFieldRefs = ref(null)
// 回复评论
function replyComment(idx, obj) {
  // 关闭所有回复框
  replyFieldRefs.value.forEach(e => e.setReply(false))
  // 打开当前点击的回复框
  const curRef = replyFieldRefs.value[idx]
  if (curRef) {
    curRef.setReply(true)
    // * 将值传给回复框
    curRef.data.nickname = obj.nickname // 用户昵称
    curRef.data.reply_user_id = obj.user_id // 回复用户 id
    curRef.data.parent_id = commentList.value[idx].id // 父评论 id
  }
}

// 提交回复后, 重新加载评论回复
const pageRefs = ref([]) // 分页
const checkRefs = ref([]) // 查看
async function reloadReplies(idx) {
  const { data } = await api.getCommentReplies(
    commentList.value[idx].id, { page_size: 5, page_num: pageRefs.value[idx].current },
  )
  // * 局部更新某个评论的回复
  commentList.value[idx].reply_vo_list = data
  commentList.value[idx].reply_count++ // 数量 + 1
  // 回复大于 5 条展示评论分页
  commentList.value[idx].reply_count > 5 && (pageRefs.value[idx].setShow(true))
  // 直接隐藏查看
  checkRefs.value[idx].style.display = 'none' // * dom 操作隐藏 "查看"
}

// "点击查看" 显示更多回复
async function checkReplies(idx, obj) {
  // 查第一页 (5 条数据)
  const { data } = await api.getCommentReplies(
    obj.id, { page_num: 1, page_size: 5 })
  // 更新对应楼评论的回复列表
  obj.reply_vo_list = data
  // 超过 5 条数据显示分页
  obj.reply_count > 5 && (pageRefs.value[idx].setShow(true))
  // 隐藏 "点击查看"
  checkRefs.value[idx].style.display = 'none' // * dom 操作隐藏 "查看"
}

// 修改回复分页中当前页数
async function changeReplyCurrent(pageNum, idx, commentId) {
  const { data } = await api.getCommentReplies(
    commentId, { page_num: pageNum, page_size: 5 })
  commentList[idx].reply_vo_list = data
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
    <p class="flex items-center text-20 font-bold">
      <span class="i-fa:comments mr-5 text-24 text-blue" /> 评论
    </p>
    <!-- 评论框 -->
    <CommentField
      :show="true"
      :type="type"
      :topic-id="topicId"
      @after-submit="reloadComments"
    />
    <!-- 评论详情 -->
    <div v-if="commentCount > 0 && refresh">
      <!-- 评论数量 -->
      <p class="mb-15 mt-30 flex items-center text-20 font-bold">
        <span> {{ commentCount }} 评论 </span>
        <span
          class="i-uiw:reload ml-15 cursor-pointer text-18"
          :class="listLoading ? 'animate-spin' : ''"
          @click="reloadComments"
        />
      </p>
      <!-- 评论列表 -->
      <div v-for="(comment, idx) of commentList" :key="comment.id" class="my-5 flex">
        <img :src="convertImgUrl(comment.avatar)" class="h-40 w-40 duration-600 hover:rotate-360">
        <div class="ml-10 flex flex-1 flex-col">
          <!-- 评论人名称 -->
          <p>
            <!-- 根据是否有 website 显示不同效果 -->
            <span v-if="!comment.website" class="text-14">
              {{ comment.nickname }}
            </span>
            <a
              v-else :href="comment.website" target="_blank"
              class="font-500 color-#1abc9c transition-300"
            >
              {{ comment.nickname }}
            </a>
            <!-- 博主标记 -->
            <span
              v-if="comment.user_id === 10"
              class="ml-6 inline-block rounded-3 bg-#ffa51e px-6 py-1 text-12 color-#fff"
            >
              博主
            </span>
          </p>
          <!-- 楼层 + 时间 + 点赞 + 回复按钮 -->
          <div class="flex justify-between text-12">
            <div class="flex items-center py-5 color-#b3b3b3">
              <span> {{ commentCount - idx }}楼 </span>
              <span class="mx-10"> {{ formatDate(comment.created_at) }} </span>
              <button
                class="i-mdi:thumb-up hover-bg-red"
                :class="isLike(comment.id) ? 'bg-red' : ''"
                @click="likeComment(comment)"
              />
              <span v-show="comment.like_count"> {{ comment.like_count }} </span>
            </div>
            <button class="color-#ef2f11" @click="replyComment(idx, comment)">
              回复
            </button>
          </div>
          <!-- 评论内容 -->
          <div class="my-3" v-html="comment.content" />

          <!-- 评论回复 start -->
          <div v-for="reply of comment.reply_vo_list" :key="reply.id" class="mt-10 flex">
            <img :src="convertImgUrl(reply.avatar)" class="h-40 w-40 duration-600 hover:rotate-360">
            <div class="ml-10 flex flex-1 flex-col">
              <!-- 回复人名称 -->
              <div>
                <!-- 根据是否有 website 显示不同效果 -->
                <span v-if="!reply.website" class="text-14">
                  {{ reply.nickname }}
                </span>
                <a
                  v-else :href="reply.website" target="_blank"
                  class="font-500 color-#1abc9c transition-300"
                >
                  {{ reply.nickname }}
                </a>
                <!-- 博主标记 -->
                <span
                  v-if="reply.user_id === 10"
                  class="ml-6 inline-block rounded-3 bg-#ffa51e px-6 py-1 text-12 color-#fff"
                >
                  博主
                </span>
              </div>
              <!-- 时间 + 点赞 + 回复按钮 -->
              <div class="flex justify-between text-12">
                <div class="flex items-center py-5 color-#b3b3b3">
                  <span class="mr-10"> {{ formatDate(reply.created_at) }} </span>
                  <button
                    class="i-mdi:thumb-up hover-bg-red"
                    :class="isLike(reply.id) ? 'bg-red' : ''"
                    @click="likeComment(reply)"
                  />
                  <span v-show="reply.like_count"> {{ reply.like_count }} </span>
                </div>
                <button class="color-#ef2f11" @click="replyComment(idx, reply)">
                  回复
                </button>
              </div>
              <!-- 回复内容 -->
              <div>
                <!-- 回复用户名: 自己回复自己不显示 "@名称" -->
                <template v-if="reply.reply_user_id !== comment.user_id">
                  <span v-if="!reply.reply_website">
                    @{{ reply.reply_nickname }}
                  </span>
                  <a v-else :href="reply.reply_website" target="_blank">
                    @{{ reply.reply_nickname }}
                  </a> ，
                </template>
                <span class="my-3" v-html="reply.content" />
              </div>
            </div>
          </div>
          <!-- 评论回复 end -->

          <!-- 回复数量 -->
          <div
            v-show="comment.reply_count > 3"
            ref="checkRefs"
            class="mt-15 text-13 color-#6d757a"
          >
            共 <b> {{ comment.reply_count }} </b>  条回复
            <button class="color-#00a1d6" @click="checkReplies(idx, comment)">
              ，点击查看
            </button>
          </div>
          <!-- 回复分页 -->
          <Paging
            ref="pageRefs"
            :page-total="Math.ceil(comment.reply_count / 5)"
            :index="idx"
            :comment-id="comment.id"
            @change-current="changeReplyCurrent"
          />
          <!-- 回复框 -->
          <CommentField
            ref="replyFieldRefs"
            :show="false"
            :type="type"
            :topic-id="topicId"
            @after-submit="reloadReplies(idx)"
          />
          <!-- 分隔线: 注意最后一个评论没有线 -->
          <div v-if="(idx + 1) !== commentCount" class="my-10 h-1 bg-light-500" />
        </div>
      </div>
      <!-- 加载更多 -->
      <div class="m-15 f-c-c">
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
    <div v-else class="mb-10 mt-30 text-center text-16 text-zinc">
      暂无评论，来发评论吧~
    </div>
  </div>
</template>
