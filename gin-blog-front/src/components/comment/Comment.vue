<script setup lang="ts">
// import EmojiList from '@/assets/js/emoji'
import CommentField from './CommentField.vue' // 评论 / 回复 框
import Paging from './Paging.vue' // 分页
import { convertImgUrl, formatDate } from '@/utils'

import { useAppStore, useUserStore } from '@/store'
import api from '@/api'

const { type } = defineProps<{
  type: number // 评论类型: 1-文章, 2-友链
}>()

const [userStore, appStore] = [useUserStore(), useAppStore()]

onMounted(() => {
  getComments()
})

// url 中存在 id 参数则为 topic_id, 否则为 0
const topicId = +(useRoute().params.id ?? 0)

// 加载评论
let commentList = $ref<any>([]) // 评论列表 (分页加载)
let commentCount = $ref(0) // 评论总数量
let listLoading = $ref(false) // 列表加载状态
const params = reactive({ type, page_size: 10, page_num: 1, topic_id: topicId }) // 加载评论的参数
async function getComments() {
  listLoading = true
  try {
    const res = await api.getComments(params)
    // * 全局加载更多, 0.8s 延时
    setTimeout(() => {
      params.page_num === 1
        ? commentList = res.data.pageData
        : commentList.push(...res.data.pageData)
      commentCount = res.data.total
      params.page_num++
      listLoading = false
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
let refresh = $ref(true) // 重新刷新整个评论列表
watch($$(commentList), () => {
  refresh = false
  nextTick(() => {
    refresh = true
  })
}, { deep: false }) // deep = false 防止 "查看更多" 时刷新整个数据

// 回复相关
// ! 可以获取 v-for 循环中的 DOM 数组
const replyFieldRefs = $ref<any>([])
// 回复评论
function replyComment(idx: number, obj: any) {
  // 关闭所有回复框
  replyFieldRefs.forEach((e: any) => e.setReply(false))
  // 打开当前点击的回复框
  const curRef = replyFieldRefs[idx]
  curRef.setReply(true)
  // * 将值传给回复框
  curRef.data.nickname = obj.nickname // 用户昵称
  curRef.data.reply_user_id = obj.user_id // 回复用户 id
  curRef.data.parent_id = commentList[idx].id // 父评论 id
}

// 提交回复后, 重新加载评论回复
const pageRefs = $ref<any>([]) // 分页
const checkRefs = $ref<any>([]) // 查看
async function reloadReplies(idx: number) {
  const { data } = await api.getCommentReplies(
    commentList[idx].id, { page_size: 5, page_num: pageRefs[idx].current },
  )
  // * 局部更新某个评论的回复
  commentList[idx].reply_vo_list = data
  commentList[idx].reply_count++ // 数量 + 1
  // 回复大于 5 条展示评论分页
  commentList[idx].reply_count > 5 && (pageRefs[idx].setShow(true))
  // 直接隐藏查看
  checkRefs[idx].style.display = 'none' // * dom 操作隐藏 "查看"
}

// "点击查看" 显示更多回复
async function checkReplies(idx: number, obj: any) {
  // 查第一页 (5 条数据)
  const { data } = await api.getCommentReplies(
    obj.id, { page_num: 1, page_size: 5 })
  // 更新对应楼评论的回复列表
  obj.reply_vo_list = data
  // 超过 5 条数据显示分页
  obj.reply_count > 5 && (pageRefs[idx].setShow(true))
  // 隐藏 "点击查看"
  checkRefs[idx].style.display = 'none' // * dom 操作隐藏 "查看"
}

// 修改回复分页中当前页数
async function changeReplyCurrent(pageNum: number, idx: number, commentId: number) {
  const { data } = await api.getCommentReplies(
    commentId, { page_num: pageNum, page_size: 5 })
  commentList[idx].reply_vo_list = data
}

// TODO: 点赞
async function likeComment(comment: any) {
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
const isLike = computed(() => (id: number) => userStore.commentLikeSet.includes(id))
</script>

<template>
  <div>
    <p text-20 font-bold flex items-center>
      <i-fa:comments mr-5 text-24 text-blue /> 评论
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
      <p mt-30 mb-15 text-20 font-bold flex items-center>
        <span> {{ commentCount }} 评论 </span>
        <i-uiw:reload
          ml-15 cursor-pointer text-18
          :class="listLoading ? 'animate-spin' : ''"
          @click="reloadComments"
        />
      </p>
      <!-- 评论列表 -->
      <div v-for="(comment, idx) of commentList" :key="comment.id" my-5 flex>
        <div>
          <n-image
            :src="convertImgUrl(comment.avatar)"
            width="40"
            duration-600 hover-rotate-360
          />
        </div>
        <div flex-col flex-1 ml-10>
          <!-- 评论人名称 -->
          <p>
            <!-- 根据是否有 website 显示不同效果 -->
            <span v-if="!comment.website" text-14>
              {{ comment.nickname }}
            </span>
            <a
              v-else :href="comment.website" target="_blank"
              color="#1abc9c" font-500 transition-300
            >
              {{ comment.nickname }}
            </a>
            <!-- 博主标记 -->
            <span
              v-if="comment.user_id === 10"
              bg="#ffa51e" color="#fff"
              px-6 py-1 ml-6 rounded-3 text-12 inline-block
            >
              博主
            </span>
          </p>
          <!-- 楼层 + 时间 + 点赞 + 回复按钮 -->
          <div flex justify-between text-12>
            <div color="#b3b3b3" py-5 flex items-center>
              <span> {{ commentCount - idx }}楼 </span>
              <span mx-10> {{ formatDate(comment.created_at) }} </span>
              <button
                i-mdi:thumb-up
                :class="isLike(comment.id) ? 'bg-red' : ''"
                hover-bg-red
                @click="likeComment(comment)"
              />
              <span v-show="comment.like_count"> {{ comment.like_count }} </span>
            </div>
            <button color="#ef2f11" @click="replyComment(idx, comment)">
              回复
            </button>
          </div>
          <!-- 评论内容 -->
          <div my-3 v-html="comment.content" />

          <!-- 评论回复 start -->
          <div v-for="reply of comment.reply_vo_list" :key="reply.id" mt-10 flex>
            <div>
              <n-image
                :src="convertImgUrl(reply.avatar)"
                width="40"
                duration-600 hover-rotate-360
              />
            </div>
            <div flex-col flex-1 ml-10>
              <!-- 回复人名称 -->
              <div>
                <!-- 根据是否有 website 显示不同效果 -->
                <span v-if="!reply.website" text-14>
                  {{ reply.nickname }}
                </span>
                <a
                  v-else :href="reply.website" target="_blank"
                  color="#1abc9c" font-500 transition-300
                >
                  {{ reply.nickname }}
                </a>
                <!-- 博主标记 -->
                <span
                  v-if="reply.user_id === 10"
                  bg="#ffa51e" color="#fff"
                  px-6 py-1 ml-6 rounded-3 text-12 inline-block
                >
                  博主
                </span>
              </div>
              <!-- 时间 + 点赞 + 回复按钮 -->
              <div flex justify-between text-12>
                <div color="#b3b3b3" py-5 flex items-center>
                  <span mr-10> {{ formatDate(reply.created_at) }} </span>
                  <button
                    i-mdi:thumb-up
                    :class="isLike(reply.id) ? 'bg-red' : ''"
                    hover-bg-red
                    @click="likeComment(reply)"
                  />
                  <span v-show="reply.like_count"> {{ reply.like_count }} </span>
                </div>
                <button color="#ef2f11" @click="replyComment(idx, reply)">
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
                <span my-3 v-html="reply.content" />
              </div>
            </div>
          </div>
          <!-- 评论回复 end -->

          <!-- 回复数量 -->
          <div
            v-show="comment.reply_count > 3"
            ref="checkRefs"
            mt-15 color="#6d757a" text-13
          >
            共 <b> {{ comment.reply_count }} </b>  条回复
            <button color="#00a1d6" @click="checkReplies(idx, comment)">
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
          <div v-if="(idx + 1) !== commentCount" h-1 bg-light-500 my-10 />
        </div>
      </div>
      <!-- 加载更多 -->
      <div f-c-c m-15>
        <n-button
          v-if="commentCount > commentList.length && !listLoading"
          text @click="getComments"
        >
          点击加载更多...
        </n-button>
        <div v-if="listLoading" animate-bounce>
          <n-spin size="small" />
        </div>
      </div>
    </div>
    <!-- 没有评论的提示 -->
    <div v-else text-center mt-30 mb-10 text-zinc text-16>
      暂无评论，来发评论吧~
    </div>
  </div>
</template>

