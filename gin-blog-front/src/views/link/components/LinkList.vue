<script setup>
import { NGi, NGrid, NImage, NResult } from 'naive-ui'
import { convertImgUrl } from '@/utils'

const { linkList } = defineProps({
  linkList: {
    type: Array,
    default: () => [],
  },
})
</script>

<template>
  <div>
    <!-- 标题 -->
    <p class="mb-30 flex items-center">
      <span class="i-mdi:link-variant mr-5 text-28 text-blue" />
      <span class="text-20 font-bold color-#344c67"> 友情链接 </span>
    </p>
    <!-- 友链数量为 0 -->
    <template v-if="!linkList.length">
      <NResult
        status="418"
        title="暂无友情链接"
        description="可以在后台添加"
      />
    </template>
    <!-- 友链数量不为 0 -->
    <template v-else>
      <!-- 链接列表 -->
      <NGrid x-gap="6" y-gap="12" cols="1 m:3" item-responsive responsive="screen">
        <NGi
          v-for="link of linkList" :key="link.id"
          :span="1"
          class="link-wrapper relative rounded-8 transition-300"
        >
          <a :href="link.address" target="_blank" class="flex flex-row p-5 hover:text-white">
            <div class="z-10 mr-5 w-120 f-c-c">
              <NImage :src="convertImgUrl(link.avatar)" width="65" class="avatar rounded-full transition-600" />
            </div>
            <div class="z-10 h-95 w-260 f-c-c text-center">
              <div>
                <p class="text-18 font-bold"> {{ link.name }} </p>
                <p class="flex-1"> {{ link.intro }} </p>
              </div>
            </div>
          </a>
        </NGi>
      </NGrid>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.link-wrapper {
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
