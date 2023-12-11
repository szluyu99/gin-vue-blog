<script setup>
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'

import UploadOne from './UploadOne.vue'
import BannerPage from '@/components/BannerPage.vue'

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

onMounted(async () => {
  await userStore.getUserInfo()
  if (!userStore.userId) {
    router.push('/')
  }
})

async function updateUserInfo() {
  try {
    await api.updateUser(form)
    window.$message?.success('修改成功!')
    userStore.getUserInfo()
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <BannerPage label="user" title="个人中心" card>
    <p class="mb-6 text-xl font-bold">
      基本信息
    </p>
    <div class="grid grid-cols-12 gap-4">
      <div class="col-span-4 f-c-c">
        <UploadOne v-model:preview="form.avatar" />
      </div>
      <div class="col-span-8 lg:col-span-7">
        <div class="my-6 space-y-3">
          <div
            v-for="item of [
              { label: '昵称', key: 'nickname' },
              { label: '个人网站', key: 'website' },
              { label: '简介', key: 'intro' },
              { label: '邮箱', key: 'email' },
            ]" :key="item.label"
          >
            <div class="mb-2">
              {{ item.label }}
            </div>
            <input
              v-model="form[item.key]" required :placeholder="`请输入${item.label}`"
              class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
            >
          </div>
        </div>
        <button class="the-button mt-2" @click="updateUserInfo">
          修改
        </button>
      </div>
      <div class="col-span-0 lg:col-span-1" />
    </div>
  </BannerPage>
</template>
