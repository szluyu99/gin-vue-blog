<script setup lang="ts">
import { convertImgUrl } from '@/utils'
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
        title="暂无友情链接"
        description="可以在后台添加"
      />
    </template>
    <!-- 友链数量不为 0 -->
    <template v-else>
      <!-- 链接列表 -->
      <n-grid
        x-gap="6" y-gap="12" cols="1 m:3"
        item-responsive responsive="screen"
      >
        <n-gi
          v-for="link of linkList" :key="link.id"
          :span="1"
          relative transition-300 rounded-8
          class="link-wrapper"
        >
          <a
            :href="link.address" target="_blank"
            flex flex-row p-5 hover:text-white
          >
            <div f-c-c mr-5 z-10 w-120>
              <n-image
                :src="convertImgUrl(link.avatar)"
                width="65"
                class="avatar"
                rounded-full transition-600
              />
            </div>
            <div f-c-c text-center w-260 h-95 z-10>
              <div>
                <p text-18 font-bold> {{ link.name }} </p>
                <p flex-1> {{ link.intro }} </p>
              </div>
            </div>
          </a>
        </n-gi>
      </n-grid>
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
