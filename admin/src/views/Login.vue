<template>
  <div class="login-page">
    <div class="login-left">
      <div class="brand">
        <div class="logo">
          <LayoutDashboard :size="40" />
        </div>
        <h1>BBS Go</h1>
        <p>现代化社区管理系统</p>
      </div>
      <div class="features">
        <div class="feature-item">
          <Shield :size="20" />
          <span>安全可靠的后台管理</span>
        </div>
        <div class="feature-item">
          <Zap :size="20" />
          <span>高效便捷的操作体验</span>
        </div>
        <div class="feature-item">
          <BarChart3 :size="20" />
          <span>实时数据统计与分析</span>
        </div>
      </div>
    </div>

    <div class="login-right">
      <div class="login-card">
        <h2>管理员登录</h2>
        <p class="subtitle">欢迎回来，请登录您的账户</p>

        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleLogin">
          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名" size="large" prefix-icon="User">
              <template #prefix>
                <User :size="18" />
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" size="large" show-password>
              <template #prefix>
                <Lock :size="18" />
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" class="login-btn" :loading="loading" native-type="submit" @click="handleLogin">
              登录
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { ElMessage } from 'element-plus'
import { LayoutDashboard, Shield, Zap, BarChart3, User, Lock } from 'lucide-vue-next'

const router = useRouter()
const adminStore = useAdminStore()
const loading = ref(false)
const formRef = ref(null)

const form = ref({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function handleLogin() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await adminStore.login(form.value)
      router.push('/')
    } catch (e) {
      ElMessage.error(e.message || '登录失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-page {
  display: flex;
  min-height: 100vh;
}

.login-left {
  flex: 1;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  color: #fff;
  padding: 60px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.brand {
  margin-bottom: 60px;
}

.logo {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-bottom: 24px;
}

.brand h1 {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.brand p {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.6);
}

.features {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
}

.feature-item svg {
  color: #667eea;
}

.login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8fafc;
  padding: 40px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  padding: 48px 40px;
  border-radius: 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
}

.login-card h2 {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  color: #9ca3af;
  margin-bottom: 32px;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 12px;
}

:deep(.el-input__prefix) {
  color: #9ca3af;
}

:deep(.el-input__wrapper) {
  border-radius: 12px;
  padding: 12px 16px;
}
</style>
