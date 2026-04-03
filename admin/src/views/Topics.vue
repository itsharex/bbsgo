<template>
  <div class="topics-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>
              <FileText :size="18" />
              帖子列表
            </h3>
            <span class="total-count">共 {{ total }} 条帖子</span>
          </div>
          <div class="header-right">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索帖子标题"
              clearable
              @clear="loadTopics"
              @keyup.enter="loadTopics"
              style="width: 220px"
            >
              <template #prefix>
                <Search :size="16" />
              </template>
            </el-input>
            <el-button type="primary" @click="loadTopics">
              <Search :size="16" />
              搜索
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="topics" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="标题" min-width="200">
          <template #default="{ row }">
            <div class="title-cell">
              <el-link :href="`/topic/${row.id}`" target="_blank" type="primary">
                {{ row.title }}
              </el-link>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="作者" width="120">
          <template #default="{ row }">
            <div class="author-cell">
              <el-avatar :size="24" :src="row.user?.avatar">
                <User :size="12" />
              </el-avatar>
              <span>{{ row.user?.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="版块" width="120">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.forum?.name || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="浏览/点赞/评论" width="150">
          <template #default="{ row }">
            <div class="stats-cell">
              <span class="stat-item">
                <Eye :size="12" />
                {{ row.view_count }}
              </span>
              <span class="stat-item">
                <Heart :size="12" />
                {{ row.like_count }}
              </span>
              <span class="stat-item">
                <MessageCircle :size="12" />
                {{ row.reply_count }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="150">
          <template #default="{ row }">
            <div class="status-cell">
              <el-tag v-if="row.is_pinned" type="danger" size="small">置顶</el-tag>
              <el-tag v-if="row.is_essence" type="warning" size="small">精华</el-tag>
              <el-tag v-if="row.is_locked" type="info" size="small">锁定</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="发布时间" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <div class=" inline-flex">
              <el-button link type="primary" @click="togglePin(row)">
                <Pin :size="14" />
                {{ row.is_pinned ? '取消置顶' : '置顶' }}
              </el-button>
              <el-button link type="primary" @click="viewTopic(row)">
                <ExternalLink :size="14" />
                查看
              </el-button>
              <el-button link type="danger" @click="deleteTopic(row)">
                <Trash2 :size="14" />
                删除
              </el-button>
            </div>

          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadTopics"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { FileText, Search, User, Eye, Heart, MessageCircle, ExternalLink, Trash2, Pin } from 'lucide-vue-next'

const topics = ref([])
const searchKeyword = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)

function formatDate(date) {
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadTopics() {
  loading.value = true
  try {
    const res = await api.get('/admin/topics', {
      params: {
        page: page.value,
        page_size: pageSize,
        keyword: searchKeyword.value
      }
    })
    topics.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('加载帖子失败', e)
    ElMessage.error('加载帖子失败')
  } finally {
    loading.value = false
  }
}

function viewTopic(topic) {
  window.open(`/topic/${topic.id}`, '_blank')
}

async function togglePin(topic) {
  try {
    await ElMessageBox.confirm(
      topic.is_pinned ? '确定要取消置顶这条帖子吗？' : '确定要置顶这条帖子吗？',
      topic.is_pinned ? '取消置顶' : '置顶帖子',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )

    await api.put(`/admin/topics/${topic.id}/pin`, {
      pinned: !topic.is_pinned
    })
    topic.is_pinned = !topic.is_pinned
    ElMessage.success(topic.is_pinned ? '帖子已置顶' : '帖子已取消置顶')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('操作失败', e)
      ElMessage.error('操作失败')
    }
  }
}

async function deleteTopic(topic) {
  try {
    await ElMessageBox.confirm(`确定要删除帖子 "${topic.title}" 吗？`, '删除帖子', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await api.delete(`/admin/topics/${topic.id}`)
    topics.value = topics.value.filter(t => t.id !== topic.id)
    total.value--
    ElMessage.success('帖子已删除')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除失败', e)
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadTopics()
})
</script>

<style scoped>
.topics-page {
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

.total-count {
  font-size: 13px;
  color: #6b7280;
  background: #f3f4f6;
  padding: 4px 10px;
  border-radius: 12px;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.title-cell {
  max-width: 300px;
}

.author-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.stats-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
}

.status-cell {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
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
