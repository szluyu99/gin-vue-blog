<script setup lang="ts">
interface Article {
  id: number
  created_at: string
  title: string
  content: string
  img: string
  category: any
  tags: any
}

interface Props {
  title?: string
  bannerImg?: string
  showFooter?: boolean
  article?: Article
}

const {
  bannerImg = 'https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg',
  showFooter = true,
} = defineProps<Props>()

const article = ref({
  content: `
  > 学习视频：[Go语言 Gin + Vue 前后端分离实战 - OceanLearn](https://www.bilibili.com/video/BV1CE411H7bQ)
  >
  > 本课程的核心思想：**快速开发，逐步优化**

  项目代码：[https://gitee.com/szluyu99/gin-demo](https://gitee.com/szluyu99/gin-demo)

  # 后端

  Go 语言学习资料：

  * [首页 - Go语言中文网 - Golang中文社区 (studygolang.com)](https://studygolang.com/)
  * [Go 语言教程 | 菜鸟教程 (runoob.com)](https://www.runoob.com/go/go-tutorial.html)
  * [Go语言圣经 (studygolang.com)](https://books.studygolang.com/gopl-zh/)
  * [Learn Go with tests - learn-go-with-tests (gitbook.io)](https://studygolang.gitbook.io/learn-go-with-tests/)
  * [Go by Example](https://gobyexample.com/)（英）



  接口测试：[ApiPost - API 文档、调试、Mock、测试一体化协作平台](https://www.apipost.cn/)

  ## Goland 操作

  **Goland 设置保存时自动格式化代码**：Preferences > Tools > File Watchers > **+** > go fmt

  **Goland 设置自动引入模块**：Preferences > Go > Go Modules > Enable Go moudules integration

  **Goland 调试小技巧**：调试模式下，右键 > Run to Cursoer，可以运行到指定行

  **定义好接口后，让某个结构体实现它**：光标在结构体上，按 ctrl + i

  创建项目：

  ![](https://szluyu99-1259132563.cos.ap-nanjing.myqcloud.com/noteImg/20220201191727.png)

  ## Golang 基础语法

  ## 项目搭建

  ## 数据库创建

  ## 项目部署
  `,
})

const styleVal = ref(`background: url('${bannerImg}') center center / cover no-repeat;`)

// 锚点目录
const preview = ref<any>(null)
const titles = ref<any>([])
onMounted(() => {
  const anchors = preview.value.$el.querySelectorAll('h1,h2,h3,h4,h5,h6')
  const titleList = Array.from(anchors).filter((t: any) => !!t.innerText.trim())
  if (!titleList.length)
    return
  const hTags = Array.from(new Set(titleList.map((t: any) => t.tagName))).sort()
  titles.value = titleList.map((el: any) => ({
    title: el.innerText,
    lineIndex: el.getAttribute('data-v-md-line'),
    indent: hTags.indexOf(el.tagName),
  }))
})
// 点击锚点目录
function handleAnchorClick(anchor: any) {
  const { lineIndex } = anchor
  const heading = preview.value.$el.querySelector(`[data-v-md-line="${lineIndex}"]`)
  if (heading) {
    preview.value.scrollToTarget({
      target: heading,
      scrollContainer: window,
      top: 60,
    })
  }
}

// 最新文章
const latestArticles = ref([
  {
    title: '为什么使用spring框架中的IOC来创建对象？它和使用new关键字来创建对象的区别是什么？',
    cover: 'https://static.talkxj.com/articles/ceb6402f4d22c463638533577f54619b.png',
    createTime: '2021-12-31',
  },
  {
    title: 'https 配置',
    cover: 'https://static.talkxj.com/articles/ceb6402f4d22c463638533577f54619b.png',
    createTime: '2021-12-31',
  },
  {
    title: 'Vue 源码解析',
    cover: 'https://static.talkxj.com/articles/ceb6402f4d22c463638533577f54619b.png',
    createTime: '2021-12-31',
  },
  {
    title: '下一站启程',
    cover: 'https://static.talkxj.com/articles/ceb6402f4d22c463638533577f54619b.png',
    createTime: '2021-12-31',
  },
  {
    title: 'Python 脚本开发',
    cover: 'https://static.talkxj.com/articles/ceb6402f4d22c463638533577f54619b.png',
    createTime: '2021-12-31',
  },
])
</script>

<template>
  <!-- 顶部图片 -->
  <div
    :style="styleVal"
    absolute inset-x-0 top-0 h-400
    flex items-center justify-center
    class="banner-fade-down"
  >
    <div text-center text-light mt-50>
      <h1 text-34 font-bold>
        文章标题
      </h1>
      <p py-8>
        发布于 2022-01-01 | 更新于 2022-01-03 | 生活随笔
      </p>
      <p>
        字数统计：1.3k |
        阅读时长：3 分钟 |
        阅读量：146 |
        评论数：45
      </p>
    </div>
  </div>
  <!-- 主体内容 -->
  <main flex-1>
    <n-grid
      x-gap="15" cols="12"
      max-w-1200 mt-475 mb-50 mx-auto px-5
      class="card-fade-up"
    >
      <n-grid-item span="9">
        <n-card
          hoverable
          rounded-2rem
        >
          <!-- 文章内容 -->
          <v-md-preview ref="preview" :text="article.content" />
        </n-card>
      </n-grid-item>
      <n-grid-item span="3">
        <div sticky top-20>
          <!-- 目录 -->
          <n-card
            hoverable
            rounded-2rem
            mb-15
          >
            <div flex items-center>
              <i-bi:list-ul text-6xl /> <span text-16 ml-10>目录</span>
            </div>
            <div
              v-for="anchor in titles"
              :key="anchor.title"
              :style="{ padding: `10px 0 10px ${anchor.indent * 20}px` }"
              @click="handleAnchorClick(anchor)"
            >
              <a style="cursor: pointer">{{ anchor.title }}</a>
            </div>
          </n-card>
          <!-- 最新文章 -->
          <n-card
            hoverable
            rounded-2rem
          >
            <div flex items-center>
              <i-akar-icons:clock text-6xl /> <span text-16 ml-10>最新文章</span>
            </div>
            <ul>
              <li
                v-for="item of latestArticles" :key="item.title"
                my-3 border-b-1 border-dashed last:border-none
              >
                <div flex items-center py-6>
                  <img :src="item.cover" overflow-hidden w-58 h-58>
                  <div flex-1 pl-10 break-all overflow-hidden>
                    <p text-14>
                      {{ item.title }}
                    </p>
                    <p text-13 text-blueGray>
                      {{ item.createTime }}
                    </p>
                  </div>
                </div>
                <!-- <n-divider dashed /> -->
              </li>
            </ul>
          </n-card>
        </div>
      </n-grid-item>
    </n-grid>
  </main>
  <!-- 底部(可选) -->
  <footer v-if="showFooter">
    <AppFooter />
  </footer>
</template>

<style lang="scss" scoped>
</style>
