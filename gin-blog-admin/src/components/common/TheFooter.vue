<script setup>
</script>

<template>
  <footer class="flex flex-col items-center justify-center text-gray-700">
    <p>
      Copyright©{{ new Date().getFullYear() }}
      <!-- 原来写的是 target="__blank": 不是保留关键字, 会开一个名叫 __blank 的窗口并被后续点击复用 -->
      <a href="https://github.com/szluyu99" target="_blank" rel="noopener noreferrer" hover="decoration-underline color-primary">
        szluyu99
      </a>
    </p>
  </footer>
</template>
