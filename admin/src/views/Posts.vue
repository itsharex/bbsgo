<template>
  <div class="posts-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>
              <MessageCircle :size="18" />
              评论列表
            </h3>
            <span class="total-count">共 {{ total }} 条评论</span>
          </div>
          <div class="header-right">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索评论内容"
              clearable
              @clear="loadPosts"
              @keyup.enter="loadPosts"
              style="width: 220px"
            >
              <template #prefix>
                <Search :size="16" />
              </template>
            </el-input>
            <el-button type="primary" @click="loadPosts">
              <Search :size="16" />
              搜索
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="posts" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="内容" min-width="250">
          <template #default="{ row }">
            <el-tooltip :content="row.content" placement="top" :disabled="row.content.length < 60">
              <span class="content-text">{{ row.content }}</span>
            </el-tooltip>
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
        <el-table-column label="所属帖子" width="120">
          <template #default="{ row }">
            <el-link :href="`/topic/${row.topic_id}`" target="_blank" type="primary">
              <ExternalLink :size="12" />
              查看
            </el-link>
          </template>
        </el-table-column>
        <el-table-column label="点赞" width="80">
          <template #default="{ row }">
            <span class="like-count">
              <Heart :size="12" />
              {{ row.like_count }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="发布时间" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="danger" @click="deletePost(row)">
              <Trash2 :size="14" />
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadPosts"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { MessageCircle, Search, User, ExternalLink, Heart, Trash2 } from 'lucide-vue-next'

const posts = ref([])
const searchKeyword = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)

function formatDate(date) {
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadPosts() {
  loading.value = true
  try {
    const res = await api.get('/admin/posts', {
      params: {
        page: page.value,
        page_size: pageSize,
        keyword: searchKeyword.value
      }
    })
    posts.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('加载评论失败', e)
    ElMessage.error('加载评论失败')
  } finally {
    loading.value = false
  }
}

async function deletePost(post) {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '删除评论', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await api.delete(`/admin/posts/${post.id}`)
    posts.value = posts.value.filter(p => p.id !== post.id)
    total.value--
    ElMessage.success('评论已删除')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除失败', e)
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadPosts()
})
</script>

<style scoped>
.posts-page {
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

.content-text {
  font-size: 13px;
  color: #374151;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.author-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.like-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #f472b6;
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
