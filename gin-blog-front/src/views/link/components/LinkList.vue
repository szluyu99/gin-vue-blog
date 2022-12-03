<script setup lang="ts">
interface Props { linkList: any }
const { linkList = [] } = defineProps<Props>()
</script>

<template>
  <div>
    <!-- 标题 -->
    <p flex items-center mb-30>
      <i-mdi:link-variant text-28 text-blue mr-5 />
      <span text-20 font-bold color="#344c67"> 友情链接 </span>
    </p>
    <!-- 友链数量为 0 -->
    <template v-if="!linkList.length">
      <n-result
        status="418"
        title="暂时还没有添加友情链接"
        description="快去后台添加吧..."
      />
    </template>
    <!-- 友链数量不为 0 -->
    <template v-else>
      <!-- 链接列表 -->
      <n-grid x-gap="12" y-gap="12" :cols="3">
        <n-gi
          v-for="link of linkList" :key="link.id"
          class="link-wrapper"
        >
          <a
            :href="link.address" target="_blank"
            flex flex-row p-5 hover:text-white
          >
            <div mr-5 flex items-center z-100>
              <n-image :src="link.avatar" width="65" class="avatar" rounded-full transition-600 />
            </div>
            <div flex-1 text-center w-290 h-95 flex-col z-100>
              <p text-18 font-bold> {{ link.name }} </p>
              <p flex-1 f-c-c> {{ link.intro }} </p>
            </div>
          </a>
        </n-gi>
      </n-grid>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.link-wrapper {
  position: relative;
  transition: all 0.3;
  border-radius: 8px;

  &:hover {
     box-shadow: 0 2px 20px #49b1f5;
  }
  &:hover::before {
    content: "";
    transform: scale(1);
  }
  &::before {
    content: "";
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    background: #49b1f5;
    border-radius: 8px;
    transition-timing-function: ease-out;
    transition-duration: 0.3s;
    transition-property: transform;
    transform: scale(0);
  }
  &:hover .avatar {
    transform: rotate(360deg);
  }
}
</style>
