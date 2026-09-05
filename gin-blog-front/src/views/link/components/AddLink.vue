<template>
  <div class="space-y-5">
    <p class="flex items-center text-xl">
      <span class="i-mdi:link-variant-plus mr-4 text-blue" />
      <span class="color-main font-bold"> 申请友链 </span>
    </p>
    <!-- 添加友链格式 -->
    <blockquote class="border-l-3 border-primary rounded-l-5 bg-#ecf7fe px-4 py-3 leading-7 dark:bg-#1d2a33">
      <p>名称：阵、雨的个人博客</p>
      <p>简介：往事随风而去</p>
      <!-- break-all: 这串 URL 已经快占满整行, 窄屏下不断行会撑破容器 -->
      <p class="break-all">
        头像：https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images/common/header.jpeg
      </p>
    </blockquote>
    <p>
      需要交换友链的在下方留言💖:
    </p>
    <blockquote class="border-l-3 border-primary rounded-l-5 bg-#ecf7fe px-4 py-3 leading-7 dark:bg-#1d2a33">
      友链信息展示需要，您的信息格式要包含：名称、介绍、链接、头像
    </blockquote>
  </div>
</template>
