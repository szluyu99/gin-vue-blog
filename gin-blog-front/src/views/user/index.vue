<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NForm, NFormItem, NGi, NGrid, NInput } from 'naive-ui'

import BannerPage from '@/components/page/BannerPage.vue'
import UploadOne from '@/components/upload/UploadOne.vue'

import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()
const router = useRouter()

const infoForm = ref({
  avatar: userStore.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
  email: userStore.email,
})

onMounted(() => {
  // 如果未登录, 退回首页
  if (!userStore.userId)
    router.push('/')
})

async function updateUserInfo() {
  try {
    await api.updateUser(infoForm)
    window.$message?.success('修改成功!')
    userStore.getUserInfo() // 重新获取用户信息
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <BannerPage label="user" title="个人中心" card>
    <p class="text-24 font-bold">
      基本信息
    </p>
    <NGrid x-gap="15" cols="12">
      <NGi span="4" class="f-c-c">
        <UploadOne v-model:preview="infoForm.avatar" :width="140" />
      </NGi>
      <NGi span="7">
        <NForm
          ref="infoFormRef"
          label-align="left"
          :label-width="80"
          :model="infoForm"
        >
          <NFormItem label="昵称" path="nickname">
            <NInput v-model:value="infoForm.nickname" placeholder="输入您的昵称" />
          </NFormItem>
          <NFormItem label="个人网站" path="website">
            <NInput v-model:value="infoForm.website" placeholder="请输入个人网站" />
          </NFormItem>
          <NFormItem label="简介" path="intro">
            <NInput v-model:value="infoForm.intro" placeholder="介绍一下自己吧" />
          </NFormItem>
          <NFormItem label="邮箱" path="email">
            <NInput v-model:value="infoForm.email" placeholder="请输入邮箱号" disabled />
          </NFormItem>
        </NForm>
        <NButton @click="updateUserInfo">
          修改
        </NButton>
      </NGi>
    </NGrid>
  </BannerPage>
</template>
