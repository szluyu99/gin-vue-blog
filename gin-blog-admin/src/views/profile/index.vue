<template>
  <CommonPage show-footer title="个人设置">
    <n-form
      :model="formModel"
      label-placement="left"
      label-align="left"
      :label-width="100"
      :style="{ maxWidth: '500px' }"
    >
      <n-form-item label="作者名称" path="name">
        <n-input v-model:value="formModel.name" />
      </n-form-item>
      <n-form-item label="个人简介" path="desc">
        <n-input
          v-model:value="formModel.desc"
          type="textarea"
          :autosize="{
            minRows: 3,
            maxRows: 5,
          }"
        />
      </n-form-item>
      <n-form-item label="网站备案号" path="icp_record">
        <n-input v-model:value="formModel.icp_record" />
      </n-form-item>
      <n-form-item label="QQ号码" path="qq_chat">
        <n-input v-model:value="formModel.qq_chat" />
      </n-form-item>
      <n-form-item label="微信" path="wechat">
        <n-input v-model:value="formModel.wechat" />
      </n-form-item>
      <n-form-item label="微博" path="weibo">
        <n-input v-model:value="formModel.weibo" />
      </n-form-item>
      <n-form-item label="B站地址" path="bili">
        <n-input v-model:value="formModel.bili" />
      </n-form-item>
      <n-form-item label="Email" path="email">
        <n-input v-model:value="formModel.email" />
      </n-form-item>
      <n-form-item label="头像" path="avatar">
        <n-upload
          action="http://127.0.0.1:8080/api/upload"
          list-type="image-card"
          :default-file-list="avaterList"
          :max="1"
          @finish="handleAaterUpload"
        >
          点击上传
        </n-upload>
      </n-form-item>
      <n-form-item label="头像背景图" path="img">
        <n-upload
          action="http://127.0.0.1:8080/api/upload"
          :default-file-list="imgList"
          list-type="image-card"
          :max="1"
          @finish="handleImgUpload"
        >
          点击上传
        </n-upload>
      </n-form-item>
      <n-button type="primary" @click="handleUpdate"> 更新 </n-button>
    </n-form>
  </CommonPage>
</template>

<script setup>
import { useAppStore } from '@/store'
import { useProfileApi } from '@/api'
const { getProfile, updateProfile } = useProfileApi()
const appStore = useAppStore()

// const baseUrl = import.meta.env.VITE_BASE_API

onMounted(async () => {
  await fetchData()
})

const formModel = ref({})
const avaterList = ref([])
const imgList = ref([])

async function handleUpdate() {
  await updateProfile(formModel.value)
  $message.success('个人信息更新成功')
  await appStore.reloadPage()
}

async function fetchData() {
  const res = await getProfile(1)
  formModel.value = res.data
  res.data.avatar &&
    (avaterList.value[0] = {
      id: res.data.avatar,
      name: 'avatar',
      status: 'finished',
      // url: 'http://127.0.0.1:8080/' + res.data.avatar,
      url: res.data.avatar,
    })
  res.data.img &&
    (imgList.value[0] = {
      id: res.data.img,
      name: 'img',
      status: 'finished',
      url: res.data.img,
    })
}

function handleAaterUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  formModel.value.avatar = res.data
}

function handleImgUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  formModel.value.img = res.data
}
</script>
