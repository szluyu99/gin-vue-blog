<script setup lang="ts">
import EasyTyper from 'easy-typer-js'
import { useAppStore } from '@/store'
const { pageList, blogConfig } = storeToRefs(useAppStore())

// 打字机特效配置
const typeObj = reactive({
  output: '',
  isEnd: false, // 全局控制是否终止
  speed: 300, // 打字速度
  singleBack: false, // 单次的回滚
  sleep: 0, // 完整输出一句话后, 睡眠一定时候后触发回滚事件
  type: 'rollback', // rollback, normal
  backSpeed: 80, // 回滚速度
  sentencePause: false, // 运行完毕后, 句子是否暂停显示
})

onMounted(() => {
  getOneSentence()
  setInterval(getOneSentence, 10000) // TODO: 让退格与新一言无缝衔接
})

function getOneSentence() {
  // 一言 + 打字机特效
  fetch('https://v1.hitokoto.cn?c=i')
    .then(res => res.json())
    .then(data => new EasyTyper(typeObj, data.hitokoto, () => {}, () => {}))
    .catch(() => new EasyTyper(typeObj, '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。', () => {}, () => {}))
}

function scrollDown() {
  window.scrollTo({
    behavior: 'smooth',
    top: document.documentElement.clientHeight,
  })
}

// 根据后端配置动态获取封面
const coverStyle = computed(() => {
  const page = pageList.value.find(e => e.label === 'home')
  return page
    ? `background: url('${page.cover}') center center / cover no-repeat;`
    : 'background: grey center center / cover no-repeat;'
})
</script>

<template>
  <div
    absolute-lrt h-screen text-center text-white
    class="banner-fade-down"
    :style="coverStyle"
  >
    <div absolute mt-43vh inset-x-0 text-center>
      <h1 text-40 font-bold animate-zoom-in>
        {{ blogConfig.website_name }}
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
