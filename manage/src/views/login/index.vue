<template>
  <div class="login-page">
    <el-card class="login-card" shadow="hover">
      <h2 class="title">管理端登录</h2>
      <el-form label-width="72px" @submit.prevent="onSubmit">
        <el-form-item label="用户名">
          <el-input v-model="form.username" autocomplete="username" placeholder="用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            placeholder="密码"
            show-password
            @keyup.enter="onSubmit"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" native-type="submit" style="width: 100%">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '@/api/auth'
import type { LoginUser } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const form = reactive({
  username: '',
  password: '',
})

async function onSubmit() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = (await login({
      username: form.username,
      password: form.password,
    })) as { data?: { token?: string; user_info?: LoginUser } }
    const token = res?.data?.token
    const info = res?.data?.user_info
    if (!token) {
      ElMessage.error('登录响应缺少 token')
      return
    }
    userStore.setSession(token, info ?? null)
    const redirect = (route.query.redirect as string) || '/'
    await router.replace(redirect)
  } catch {
    /* 错误已在 request 拦截器提示 */
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(160deg, #0f172a 0%, #1e293b 45%, #312e81 100%);
  padding: 24px;
}
.login-card {
  width: 100%;
  max-width: 400px;
  border-radius: 12px;
}
.title {
  margin: 0 0 20px;
  text-align: center;
  font-size: 1.25rem;
  color: #0f172a;
}
</style>
