<template>
  <div class="config-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <h3>
            <Settings :size="18" />
            网站配置
          </h3>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本设置" name="basic">
          <el-form :model="config" label-position="top">
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="网站名称">
                  <el-input v-model="config.site_name" placeholder="请输入网站名称">
                    <template #prefix>
                      <Globe :size="16" />
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="网站 Logo URL">
                  <el-input v-model="config.site_logo" placeholder="请输入 Logo 图片 URL">
                    <template #prefix>
                      <Image :size="16" />
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="网站 Icon URL">
                  <el-input v-model="config.site_icon" placeholder="请输入 Icon 图片 URL">
                    <template #prefix>
                      <Image :size="16" />
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="网站描述">
                  <el-input v-model="config.site_description" type="textarea" :rows="2" placeholder="请输入网站描述" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="邮件配置" name="email">
          <el-alert
            title="邮件服务说明"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px"
          >
            <template #default>
              <p style="margin: 0">邮件服务用于用户注册时发送邮箱验证码，确保用户邮箱真实有效。</p>
              <p style="margin: 8px 0 0 0; color: #909399; font-size: 13px">
                启用邮件服务后，用户注册时需要输入邮箱验证码才能完成注册。
              </p>
            </template>
          </el-alert>

          <el-form :model="config" label-position="top">
            <el-form-item label="启用邮件服务">
              <el-switch v-model="config.email_enabled" />
              <div style="color: #909399; font-size: 12px; margin-top: 4px">
                启用后，用户注册时需要邮箱验证码；禁用则跳过邮箱验证
              </div>
            </el-form-item>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="SMTP 服务器">
                  <el-input v-model="config.email_host" placeholder="smtp.example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="SMTP 端口">
                  <el-input v-model="config.email_port" placeholder="465" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="邮箱账号">
                  <el-input v-model="config.email_user" placeholder="noreply@example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="邮箱密码">
                  <el-input v-model="config.email_password" type="password" placeholder="请输入邮箱密码" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="发件人地址">
                  <el-input v-model="config.email_from" placeholder="noreply@example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="发件人名称">
                  <el-input v-model="config.email_from_name" placeholder="网站名称" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="七牛云配置" name="qiniu">
          <el-form :model="config" label-position="top">
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="Access Key">
                  <el-input v-model="config.qiniu_access_key" placeholder="请输入 Access Key" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Secret Key">
                  <el-input v-model="config.qiniu_secret_key" type="password" placeholder="请输入 Secret Key" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="Bucket 名称">
                  <el-input v-model="config.qiniu_bucket" placeholder="请输入 Bucket 名称" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="CDN 域名">
                  <el-input v-model="config.qiniu_domain" placeholder="cdn.example.com" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="安全设置" name="security">
          <el-form :model="config" label-position="top">
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="JWT Secret">
                  <el-input v-model="config.jwt_secret" type="password" placeholder="请输入 JWT Secret" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Token 过期天数">
                  <el-input v-model="config.jwt_expire_days" placeholder="7" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div class="form-footer">
        <el-button @click="loadConfig">重置</el-button>
        <el-button type="primary" @click="saveConfig" :loading="saving">
          <Save :size="16" />
          保存配置
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { Settings, Globe, Image, Save } from 'lucide-vue-next'

const activeTab = ref('basic')
const config = ref({
  site_name: '',
  site_logo: '',
  site_icon: '',
  site_description: '',
  email_enabled: false,
  email_host: '',
  email_port: '465',
  email_user: '',
  email_password: '',
  email_from: '',
  email_from_name: '',
  qiniu_access_key: '',
  qiniu_secret_key: '',
  qiniu_bucket: '',
  qiniu_domain: '',
  jwt_secret: '',
  jwt_expire_days: '7'
})

const saving = ref(false)

async function loadConfig() {
  try {
    const res = await api.get('/config')
    if (res) {
      config.value = {
        site_name: res.site_name || '',
        site_logo: res.site_logo || '',
        site_icon: res.site_icon || '',
        site_description: res.site_description || '',
        email_enabled: res.email_enabled === 'true',
        email_host: res.email_host || '',
        email_port: res.email_port || '465',
        email_user: res.email_user || '',
        email_password: res.email_password || '',
        email_from: res.email_from || '',
        email_from_name: res.email_from_name || '',
        qiniu_access_key: res.qiniu_access_key || '',
        qiniu_secret_key: res.qiniu_secret_key || '',
        qiniu_bucket: res.qiniu_bucket || '',
        qiniu_domain: res.qiniu_domain || '',
        jwt_secret: res.jwt_secret || '',
        jwt_expire_days: res.jwt_expire_days || '7'
      }
    }
  } catch (e) {
    console.error('加载配置失败', e)
    ElMessage.error('加载配置失败')
  }
}

async function saveConfig() {
  saving.value = true
  try {
    const data = {
      ...config.value,
      email_enabled: config.value.email_enabled ? 'true' : 'false'
    }
    await api.put('/admin/config', data)
    ElMessage.success('保存成功')
    loadConfig()
  } catch (e) {
    console.error('保存失败', e)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.config-page {
  max-width: 1000px;
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

.form-footer {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #f3f4f6;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
}

:deep(.el-input__prefix) {
  color: #9ca3af;
}

:deep(.el-tabs__item) {
  font-weight: 500;
}
</style>
