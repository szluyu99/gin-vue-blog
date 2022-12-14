<script setup lang="ts">
let sentence = $ref('书山有路勤为径，学海无涯苦作舟。')

onMounted(() => {
  fetch('https://v1.hitokoto.cn?c=i')
    .then(res => res.json())
    .then(data => sentence = data.hitokoto)
})
</script>

<template>
  <div card-view animate-zoom-in cursor-pointer>
    <div flex text-center>
      <i-mdi-chat-outline text-20 />
      <div flex-1>
        {{ sentence }}
      </div>
      <i-mdi-chevron-double-right class="arrow" text-20 />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.arrow {
  animation: 1s passing infinite;
}

// 左 -> 右 闪的特效
@keyframes passing {
  0% {
    transform: translateX(-50%);
    opacity: 0;
  }

  50% {
    transform: translateX(0);
    opacity: 1;
  }

  100% {
    transform: translateX(50%);
    opacity: 0;
  }
}
</style>
