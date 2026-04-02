<template>
  <div class="reports-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>
              <AlertTriangle :size="18" />
              举报列表
            </h3>
            <span v-if="pendingCount > 0" class="pending-badge">
              <Bell :size="12" />
              {{ pendingCount }} 条待处理
            </span>
          </div>
          <div class="header-right">
            <el-select v-model="filterStatus" placeholder="选择状态" clearable style="width: 140px">
              <el-option label="全部状态" value="" />
              <el-option label="待处理" value="0" />
              <el-option label="已通过" value="1" />
              <el-option label="已驳回" value="2" />
            </el-select>
            <el-button type="primary" @click="loadReports">
              <RefreshCw :size="16" />
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="reports" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="举报人" width="140">
          <template #default="{ row }">
            <div class="reporter-cell">
              <el-avatar :size="28" :src="row.reporter?.avatar">
                <User :size="14" />
              </el-avatar>
              <span>{{ row.reporter?.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeType(row.target_type)" size="small">
              {{ getTargetType(row.target_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="举报原因">
          <template #default="{ row }">
            <el-tooltip :content="row.reason" placement="top" :disabled="row.reason.length < 50">
              <span class="reason-text">{{ row.reason }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="举报时间" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 0">
              <el-button link type="success" @click="handleReport(row, true)">
                <Check :size="14" />
                通过
              </el-button>
              <el-button link type="danger" @click="handleReport(row, false)">
                <X :size="14" />
                驳回
              </el-button>
            </template>
            <span v-else class="handled-text">已处理</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { AlertTriangle, Bell, User, Check, X, RefreshCw } from 'lucide-vue-next'

const reports = ref([])
const filterStatus = ref('')
const loading = ref(false)

const pendingCount = computed(() => reports.value.filter(r => r.status === 0).length)

function getTargetType(type) {
  const types = { topic: '帖子', post: '评论', message: '私信' }
  return types[type] || type
}

function getTypeType(type) {
  const types = { topic: 'primary', post: 'info', message: 'warning' }
  return types[type] || 'info'
}

function getStatusName(status) {
  const statuses = { 0: '待处理', 1: '已通过', 2: '已驳回' }
  return statuses[status] || '未知'
}

function getStatusType(status) {
  const types = { 0: 'warning', 1: 'success', 2: 'info' }
  return types[status] || 'info'
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadReports() {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value !== '') {
      params.status = filterStatus.value
    }
    const res = await api.get('/admin/reports', { params })
    reports.value = res || []
  } catch (e) {
    console.error('加载举报列表失败', e)
    ElMessage.error('加载举报列表失败')
  } finally {
    loading.value = false
  }
}

async function handleReport(report, approved) {
  const action = approved ? '通过' : '驳回'
  try {
    await ElMessageBox.confirm(
      approved ? '通过后相关内容将被删除，是否继续？' : '确定要驳回这条举报吗？',
      `${action}举报`,
      { confirmButtonText: action, cancelButtonText: '取消', type: approved ? 'warning' : 'info' }
    )

    await api.put(`/admin/reports/${report.id}/handle`, { approved })
    report.status = approved ? 1 : 2
    ElMessage.success(approved ? '举报已通过，相关内容已删除' : '举报已驳回')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('处理举报失败', e)
      ElMessage.error('操作失败')
    }
  }
}

onMounted(() => {
  loadReports()
})
</script>

<style scoped>
.reports-page {
  max-width: 1400px;
}

.main-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.pending-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.1);
  padding: 4px 10px;
  border-radius: 12px;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.reporter-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.reason-text {
  font-size: 13px;
  color: #374151;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
}

.handled-text {
  font-size: 12px;
  color: #9ca3af;
}

:deep(.el-table) {
  --el-table-border-color: #f3f4f6;
  --el-table-header-bg-color: #f9fafb;
}

:deep(.el-table th) {
  font-weight: 600;
  color: #374151;
}

:deep(.el-button.is-link) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>
