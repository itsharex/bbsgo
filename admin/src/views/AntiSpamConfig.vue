<template>
  <div class="p-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">防刷系统配置</h1>
      <p class="text-gray-600 mt-1">配置防刷策略和参数，系统将自动执行所有规则</p>
    </div>

    <el-tabs v-model="activeTab" class="mb-6">
      <el-tab-pane label="频率限制" name="rate">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">操作频率限制</h3>
          <el-form :model="rateConfig" label-width="200px">
            <el-form-item label="发帖最小间隔（秒）">
              <el-input-number v-model="rateConfig.topic_min_interval" :min="10" :max="600" />
            </el-form-item>
            <el-form-item label="评论最小间隔（秒）">
              <el-input-number v-model="rateConfig.comment_min_interval" :min="5" :max="300" />
            </el-form-item>
            <el-form-item label="单日最大发帖数">
              <el-input-number v-model="rateConfig.max_topics_per_day" :min="1" :max="100" />
            </el-form-item>
            <el-form-item label="单日最大评论数">
              <el-input-number v-model="rateConfig.max_comments_per_day" :min="1" :max="200" />
            </el-form-item>
            <el-divider />
            <h4 class="text-md font-medium mb-3">新用户限制</h4>
            <el-form-item label="新用户判定时长（小时）">
              <el-input-number v-model="rateConfig.new_user_hours" :min="1" :max="168" />
            </el-form-item>
            <el-form-item label="新用户单日发帖上限">
              <el-input-number v-model="rateConfig.new_user_max_topics_per_day" :min="1" :max="20" />
            </el-form-item>
            <el-form-item label="新用户单日评论上限">
              <el-input-number v-model="rateConfig.new_user_max_comments_per_day" :min="1" :max="50" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveRateConfig">保存配置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane label="内容质量" name="quality">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">内容质量检测</h3>
          <el-form :model="qualityConfig" label-width="200px">
            <el-form-item label="最小内容长度">
              <el-input-number v-model="qualityConfig.min_content_length" :min="5" :max="100" />
              <span class="ml-2 text-gray-500">汉字数量</span>
            </el-form-item>
            <el-form-item label="重复字符阈值">
              <el-input-number v-model="qualityConfig.repeat_char_threshold" :min="3" :max="20" />
              <span class="ml-2 text-gray-500">连续重复字符数</span>
            </el-form-item>
            <el-form-item label="内容相似度阈值">
              <el-slider v-model="qualityConfig.similarity_threshold" :min="0.5" :max="1" :step="0.1" show-input />
              <span class="ml-2 text-gray-500">超过此阈值判定为重复内容</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveQualityConfig">保存配置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane label="信誉分系统" name="reputation">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">信誉分配置</h3>
          <el-form :model="reputationConfig" label-width="200px">
            <el-form-item label="低信誉分阈值">
              <el-input-number v-model="reputationConfig.low_reputation_threshold" :min="20" :max="80" />
              <span class="ml-2 text-gray-500">低于此分数将降低内容热度</span>
            </el-form-item>
            <el-form-item label="禁言信誉分阈值">
              <el-input-number v-model="reputationConfig.ban_reputation_threshold" :min="0" :max="50" />
              <span class="ml-2 text-gray-500">低于此分数将自动禁言</span>
            </el-form-item>
            <el-form-item label="启用低信誉分禁言">
              <el-switch v-model="reputationConfig.ban_low_reputation" />
            </el-form-item>
            <el-divider />
            <h4 class="text-md font-medium mb-3">热度权重</h4>
            <el-form-item label="低质量内容热度系数">
              <el-slider v-model="reputationConfig.low_quality_hot_multiplier" :min="0" :max="1" :step="0.1" show-input />
            </el-form-item>
            <el-form-item label="低信誉分内容热度系数">
              <el-slider v-model="reputationConfig.low_reputation_hot_multiplier" :min="0" :max="1" :step="0.1" show-input />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveReputationConfig">保存配置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane label="举报处理" name="report">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">举报自动处理</h3>
          <el-form :model="reportConfig" label-width="200px">
            <el-form-item label="自动隐藏举报阈值">
              <el-input-number v-model="reportConfig.report_threshold" :min="1" :max="10" />
              <span class="ml-2 text-gray-500">达到此举报数自动隐藏内容</span>
            </el-form-item>
            <el-form-item label="自动禁言违规阈值">
              <el-input-number v-model="reportConfig.report_ban_threshold" :min="1" :max="20" />
              <span class="ml-2 text-gray-500">7天内被隐藏内容达到此数自动禁言</span>
            </el-form-item>
            <el-form-item label="自动禁言天数">
              <el-input-number v-model="reportConfig.report_ban_days" :min="1" :max="30" />
            </el-form-item>
            <el-form-item label="每日举报上限">
              <el-input-number v-model="reportConfig.max_reports_per_day" :min="1" :max="50" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveReportConfig">保存配置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane label="用户信誉" name="users">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium">用户信誉分管理</h3>
            <el-input v-model="userSearch" placeholder="搜索用户" style="width: 300px" clearable>
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
          
          <el-table :data="filteredUsers" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="username" label="用户名" width="150" />
            <el-table-column prop="nickname" label="昵称" width="150" />
            <el-table-column label="信誉分" width="200">
              <template #default="{ row }">
                <el-progress 
                  :percentage="row.reputation" 
                  :color="getReputationColor(row.reputation)"
                  :stroke-width="20"
                  :text-inside="true"
                />
              </template>
            </el-table-column>
            <el-table-column label="信誉等级" width="120">
              <template #default="{ row }">
                <el-tag :type="getReputationTagType(row.reputation)">
                  {{ getReputationLevel(row.reputation) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="注册时间" width="180" />
            <el-table-column label="操作" fixed="right" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="showReputationDialog(row)">调整信誉分</el-button>
                <el-button size="small" type="danger" @click="banUser(row)" v-if="!isBanned(row)">禁言</el-button>
                <el-button size="small" type="success" @click="unbanUser(row)" v-else>解禁</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="flex justify-center mt-4">
            <el-pagination
              v-model:current-page="userPage"
              :page-size="20"
              :total="userTotal"
              layout="prev, pager, next"
              @current-change="loadUsers"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="reputationDialogVisible" title="调整用户信誉分" width="400px">
      <el-form :model="reputationForm" label-width="100px">
        <el-form-item label="当前信誉分">
          <el-input :value="currentUser?.reputation" disabled />
        </el-form-item>
        <el-form-item label="调整分数">
          <el-input-number v-model="reputationForm.change" :min="-100" :max="100" />
        </el-form-item>
        <el-form-item label="调整原因">
          <el-input v-model="reputationForm.reason" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reputationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="adjustReputation">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import api from '@/api'

const activeTab = ref('rate')
const userSearch = ref('')
const userPage = ref(1)
const userTotal = ref(0)
const users = ref([])
const reputationDialogVisible = ref(false)
const currentUser = ref(null)

const rateConfig = ref({
  topic_min_interval: 60,
  comment_min_interval: 30,
  max_topics_per_day: 10,
  max_comments_per_day: 50,
  new_user_hours: 24,
  new_user_max_topics_per_day: 3,
  new_user_max_comments_per_day: 10
})

const qualityConfig = ref({
  min_content_length: 10,
  repeat_char_threshold: 5,
  similarity_threshold: 0.8
})

const reputationConfig = ref({
  low_reputation_threshold: 60,
  ban_reputation_threshold: 20,
  ban_low_reputation: true,
  low_quality_hot_multiplier: 0.3,
  low_reputation_hot_multiplier: 0.5
})

const reportConfig = ref({
  report_threshold: 3,
  report_ban_threshold: 5,
  report_ban_days: 3,
  max_reports_per_day: 10
})

const reputationForm = ref({
  change: 0,
  reason: ''
})

const filteredUsers = computed(() => {
  if (!userSearch.value) return users.value
  const search = userSearch.value.toLowerCase()
  return users.value.filter(u => 
    u.username?.toLowerCase().includes(search) ||
    u.nickname?.toLowerCase().includes(search)
  )
})

onMounted(() => {
  loadConfigs()
  loadUsers()
})

async function loadConfigs() {
  try {
    const res = await api.get('/admin/antispam/config')
    if (res) {
      Object.keys(res).forEach(key => {
        if (rateConfig.value.hasOwnProperty(key)) {
          rateConfig.value[key] = parseInt(res[key])
        } else if (qualityConfig.value.hasOwnProperty(key)) {
          const val = res[key]
          qualityConfig.value[key] = key === 'similarity_threshold' ? parseFloat(val) : parseInt(val)
        } else if (reputationConfig.value.hasOwnProperty(key)) {
          const val = res[key]
          if (key === 'ban_low_reputation') {
            reputationConfig.value[key] = val === 'true'
          } else if (key.includes('multiplier')) {
            reputationConfig.value[key] = parseFloat(val)
          } else {
            reputationConfig.value[key] = parseInt(val)
          }
        } else if (reportConfig.value.hasOwnProperty(key)) {
          reportConfig.value[key] = parseInt(res[key])
        }
      })
    }
  } catch (e) {
    console.error('加载配置失败', e)
  }
}

async function saveRateConfig() {
  await saveConfig(rateConfig.value)
}

async function saveQualityConfig() {
  await saveConfig(qualityConfig.value)
}

async function saveReputationConfig() {
  await saveConfig(reputationConfig.value)
}

async function saveReportConfig() {
  await saveConfig(reportConfig.value)
}

async function saveConfig(config) {
  try {
    await api.post('/admin/antispam/config', config)
    ElMessage.success('配置保存成功')
  } catch (e) {
    ElMessage.error('配置保存失败')
    console.error(e)
  }
}

async function loadUsers() {
  try {
    const res = await api.get('/admin/users', {
      params: { page: userPage.value }
    })
    users.value = res.list || []
    userTotal.value = res.total || 0
  } catch (e) {
    console.error('加载用户失败', e)
  }
}

function getReputationColor(reputation) {
  if (reputation >= 80) return '#67c23a'
  if (reputation >= 60) return '#409eff'
  if (reputation >= 40) return '#e6a23c'
  if (reputation >= 20) return '#f56c6c'
  return '#909399'
}

function getReputationTagType(reputation) {
  if (reputation >= 80) return 'success'
  if (reputation >= 60) return ''
  if (reputation >= 40) return 'warning'
  return 'danger'
}

function getReputationLevel(reputation) {
  if (reputation >= 80) return '正常'
  if (reputation >= 60) return '需验证'
  if (reputation >= 40) return '受限'
  if (reputation >= 20) return '严重'
  return '禁言'
}

function isBanned(user) {
  return user.reputation < 20
}

function showReputationDialog(user) {
  currentUser.value = user
  reputationForm.value = { change: 0, reason: '' }
  reputationDialogVisible.value = true
}

async function adjustReputation() {
  try {
    await api.post(`/admin/users/${currentUser.value.id}/reputation`, reputationForm.value)
    ElMessage.success('信誉分调整成功')
    reputationDialogVisible.value = false
    loadUsers()
  } catch (e) {
    ElMessage.error('调整失败')
    console.error(e)
  }
}

async function banUser(user) {
  try {
    await ElMessageBox.confirm('确定要禁言该用户吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await api.post(`/admin/users/${user.id}/ban`, { reason: '管理员手动禁言' })
    ElMessage.success('用户已被禁言')
    loadUsers()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败')
      console.error(e)
    }
  }
}

async function unbanUser(user) {
  try {
    await api.post(`/admin/users/${user.id}/unban`)
    ElMessage.success('用户已解禁')
    loadUsers()
  } catch (e) {
    ElMessage.error('操作失败')
    console.error(e)
  }
}
</script>
