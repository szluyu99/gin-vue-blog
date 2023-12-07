<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import BannerPage from '@/components/page/BannerPage.vue'
import UploadOne from '@/components/upload/UploadOne.vue'

import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()
const router = useRouter()

const form = reactive({
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
    await api.updateUser(form)
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
    <p class="mb-6 text-22 font-bold">
      基本信息
    </p>
    <div class="grid grid-cols-12 gap-15">
      <div class="col-span-4 f-c-c">
        <UploadOne v-model:preview="form.avatar" />
      </div>
      <div class="col-span-7">
        <div class="my-25 space-y-12">
          <div
            v-for="item of [
              { label: '昵称', key: 'nickname' },
              { label: '个人网站', key: 'website' },
              { label: '简介', key: 'intro' },
              { label: '邮箱', key: 'email' },
            ]" :key="item.label"
          >
            <div class="mb-5">
              {{ item.label }}
            </div>
            <input
              v-model="form[item.key]" required :placeholder="`请输入${item.label}`"
              class="block w-full border-0 rounded-md px-8 py-6 text-15 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
            >
          </div>
        </div>
        <button class="the-button mt-8" @click="updateUserInfo">
          修改
        </button>
      </div>
      <div class="col-span-1" />
    </div>
  </BannerPage>
</template>
