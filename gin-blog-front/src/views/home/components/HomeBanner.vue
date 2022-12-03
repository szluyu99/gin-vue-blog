<script setup lang="ts">
import EasyTyper from 'easy-typer-js'

// 打字机特效配置
const typeObj = reactive({
  output: '',
  isEnd: false, // 全局控制是否终止
  speed: 300, // 打字速度
  singleBack: false, // 单次的回滚
  sleep: 0, // 完整输出一句话后, 睡眠一定时候后触发回滚事件
  type: 'normal', // rollback
  backSpeed: 40, // 回滚速度
  sentencePause: false, // 运行完毕后, 句子是否暂停显示
})

onMounted(() => {
  // 一言 + 打字机特效
  fetch('https://v1.hitokoto.cn?c=i')
    .then(res => res.json())
    .then(data => new EasyTyper(typeObj, data.hitokoto, () => {}, () => {}))
    .catch(() => new EasyTyper(typeObj, '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。', () => {}, () => {}))
})

// TODO:  后端动态配置
const cover = computed(() => 'background: url("https://static.talkxj.com/config/0bee7ba5ac70155766648e14ae2a821f.jpg") center center /cover no-repeat')

function scrollDown() {
  window.scrollTo({
    behavior: 'smooth',
    top: document.documentElement.clientHeight,
  })
}

const blogTitle = import.meta.env.VITE_APP_TITLE
</script>

<template>
  <div
    absolute-lrt h-screen text-center text-white
    :style="cover" class="banner-fade-down"
  >
    <div absolute mt-43vh inset-x-0 text-center>
      <h1 text-40 font-bold animate-zoom-in>
        {{ `${blogTitle}的个人博客` }}
      </h1>
      <div text-22>
        {{ typeObj.output }}
        <span animate-ping> | </span>
      </div>
    </div>
    <!-- 向下滚动 -->
    <div absolute bottom-10 w-full cursor-pointer @click="scrollDown">
      <i-ep:arrow-down-bold animate-bounce text-white text-25 />
    </div>
  </div>
</template>

<style lang="scss" scoped>
</style>
