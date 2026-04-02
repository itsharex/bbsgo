<template>
  <div class="change-password-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <h3>
            <Key :size="18" />
            修改密码
          </h3>
        </div>
      </template>

      <div class="form-wrapper">
        <el-form ref="formRef" :model="form" :rules="rules" label-position="top" class="password-form">
          <el-form-item label="原密码" prop="old_password">
            <el-input v-model="form.old_password" type="password" placeholder="请输入原密码" show-password size="large">
              <template #prefix>
                <Lock :size="18" />
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="新密码" prop="new_password">
            <el-input v-model="form.new_password" type="password" placeholder="请输入新密码（至少6位）" show-password size="large">
              <template #prefix>
                <Key :size="18" />
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="确认新密码" prop="confirm_password">
            <el-input v-model="form.confirm_password" type="password" placeholder="请再次输入新密码" show-password size="large">
              <template #prefix>
                <KeyRound :size="18" />
              </template>
            </el-input>
          </el-form-item>

          <div class="form-actions">
            <el-button size="large" @click="resetForm">重置</el-button>
            <el-button type="primary" size="large" @click="handleSubmit" :loading="loading">
              <Save :size="16" />
              保存修改
            </el-button>
          </div>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { Key, Lock, KeyRound, Save } from 'lucide-vue-next'

const router = useRouter()
const adminStore = useAdminStore()
const formRef = ref(null)
const loading = ref(false)

const form = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirm = (rule, value, callback) => {
  if (value !== form.value.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  old_password: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' }
  ]
}

function resetForm() {
  form.value = {
    old_password: '',
    new_password: '',
    confirm_password: ''
  }
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await api.post('/admin/change-password', {
        old_password: form.value.old_password,
        new_password: form.value.new_password
      })
      ElMessage.success('密码修改成功，请重新登录')
      adminStore.logout()
      router.push('/login')
    } catch (e) {
      console.error('密码修改失败', e)
      ElMessage.error(e.response?.data?.message || '密码修改失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.change-password-page {
  max-width: 1400px;
}

.main-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.card-header h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.form-wrapper {
  padding: 8px 0;

}

.password-form {
  max-width: 420px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 32px;
}

.form-actions .el-button {
  flex: 1;
}

:deep(.el-input__prefix) {
  color: #9ca3af;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
}
</style>
