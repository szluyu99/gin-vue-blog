<script setup lang="ts">
import UploadOne from '@/components/upload/UploadOne.vue'

import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()
const router = useRouter()

const infoForm = $ref<any>({
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
    <p text-24 font-bold>
      基本信息
    </p>
    <n-grid x-gap="15" cols="12">
      <n-gi span="4" f-c-c>
        <UploadOne
          v-model:preview="infoForm.avatar"
          :width="140"
        />
      </n-gi>
      <n-gi span="7">
        <n-form
          ref="infoFormRef"
          label-align="left"
          :label-width="80"
          :model="infoForm"
        >
          <n-form-item label="昵称" path="nickname">
            <n-input v-model:value="infoForm.nickname" placeholder="输入您的昵称" />
          </n-form-item>
          <n-form-item label="个人网站" path="website">
            <n-input v-model:value="infoForm.website" placeholder="请输入个人网站" />
          </n-form-item>
          <n-form-item label="简介" path="intro">
            <n-input v-model:value="infoForm.intro" placeholder="介绍一下自己吧" />
          </n-form-item>
          <n-form-item label="邮箱" path="email">
            <n-input v-model:value="infoForm.email" placeholder="请输入邮箱号" disabled />
          </n-form-item>
        </n-form>
        <n-button @click="updateUserInfo">
          修改
        </n-button>
      </n-gi>
    </n-grid>
  </BannerPage>
</template>
