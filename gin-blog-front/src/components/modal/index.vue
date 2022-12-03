<script setup lang="ts">
import SearchModal from '@/components/modal/SearchModal.vue'
import LoginModal from '@/components/modal/LoginModal.vue'
import RegisterModal from '@/components/modal/RegisterModal.vue'
import ForgetModal from '@/components/modal/ForgetModal.vue'
</script>

<template>
  <!-- 搜索窗口 -->
  <SearchModal />
  <!-- 登录窗口 -->
  <LoginModal />
  <!-- 注册窗口 -->
  <RegisterModal />
  <!-- 忘记密码窗口 -->
  <ForgetModal />
</template>
