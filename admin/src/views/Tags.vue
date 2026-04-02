<template>
  <div class="tags-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>
              <Tag :size="18" />
              话题标签管理
            </h3>
          </div>
          <el-button type="primary" @click="openAddModal">
            <Plus :size="16" />
            添加官方标签
          </el-button>
        </div>
      </template>

      <div class="info-tip">
        <Info :size="16" />
        <span>用户发帖时自动创建标签，管理员仅需处理违规标签</span>
      </div>

      <el-table :data="tags" stripe style="width: 100%" v-loading="loading" :row-class-name="getRowClass">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="名称" min-width="180">
          <template #default="{ row }">
            <div class="tag-name">
              <span v-if="row.icon" class="tag-icon">{{ row.icon }}</span>
              <span :class="{ 'line-through text-gray-400': row.is_banned }">{{ row.name }}</span>
              <el-tag v-if="row.is_official" type="primary" size="small">官方</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="usage_count" label="使用次数" width="120" sortable />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_banned ? 'danger' : 'success'" size="small">
              {{ row.is_banned ? '已禁用' : '正常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editTag(row)">
              <Edit :size="14" />
              编辑
            </el-button>
            <el-button link :type="row.is_official ? 'info' : 'success'" @click="toggleOfficial(row)">
              <Star :size="14" />
              {{ row.is_official ? '取消官方' : '设为官方' }}
            </el-button>
            <el-button link :type="row.is_banned ? 'success' : 'warning'" @click="toggleBan(row)">
              <Ban :size="14" />
              {{ row.is_banned ? '解禁' : '禁用' }}
            </el-button>
            <el-button link type="info" @click="openMergeModal(row)">
              <GitMerge :size="14" />
              合并
            </el-button>
            <el-button link type="danger" @click="deleteTag(row)">
              <Trash2 :size="14" />
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editingTag ? '编辑标签' : '添加官方标签'" width="420px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="form" label-position="top">
        <el-form-item label="名称" prop="name" :rules="[{ required: true, message: '请输入标签名称', trigger: 'blur' }]">
          <el-input v-model="form.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="图标（Emoji）" prop="icon">
          <el-input v-model="form.icon" placeholder="🌐">
            <template #append>
              <span class="emoji-preview">{{ form.icon || '👁' }}</span>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item v-if="!editingTag">
          <el-checkbox v-model="form.is_official">设为官方推荐标签</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTag" :loading="saving">{{ editingTag ? '更新' : '创建' }}</el-button>
      </template>
    </el-dialog>

    <!-- 合并弹窗 -->
    <el-dialog v-model="mergeDialogVisible" title="合并标签" width="420px">
      <div class="merge-tip">
        <AlertCircle :size="16" />
        <span>将标签 <strong>"{{ mergeSource?.name }}"</strong> 合并到：</span>
      </div>
      <el-select v-model="mergeTargetId" placeholder="请选择目标标签" style="width: 100%; margin-top: 12px">
        <el-option v-for="tag in availableTags" :key="tag.id" :label="`${tag.name} (使用${tag.usage_count}次)`" :value="tag.id" />
      </el-select>
      <template #footer>
        <el-button @click="mergeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="mergeTags" :disabled="!mergeTargetId" :loading="saving">确认合并</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Tag, Plus, Edit, Star, Ban, GitMerge, Trash2, Info, AlertCircle } from 'lucide-vue-next'

const tags = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const mergeDialogVisible = ref(false)
const saving = ref(false)
const editingTag = ref(null)
const mergeSource = ref(null)
const mergeTargetId = ref('')
const formRef = ref(null)

const form = ref({
  name: '',
  icon: '',
  is_official: true
})

const availableTags = computed(() => tags.value.filter(t => t.id !== mergeSource.value?.id))

function getRowClass({ row }) {
  return row.is_banned ? 'banned-row' : ''
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadTags() {
  loading.value = true
  try {
    const res = await api.get('/admin/tags')
    tags.value = res || []
  } catch (e) {
    console.error('加载标签失败', e)
    ElMessage.error('加载标签失败')
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  editingTag.value = null
  form.value = { name: '', icon: '', is_official: true }
  dialogVisible.value = true
}

function editTag(tag) {
  editingTag.value = tag
  form.value = { name: tag.name, icon: tag.icon || '', is_official: tag.is_official }
  dialogVisible.value = true
}

async function saveTag() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    saving.value = true
    try {
      if (editingTag.value) {
        await api.put(`/admin/tags/${editingTag.value.id}`, form.value)
        Object.assign(editingTag.value, form.value)
        ElMessage.success('标签已更新')
      } else {
        await api.post('/admin/tags', form.value)
        ElMessage.success('标签已创建')
        loadTags()
      }
      dialogVisible.value = false
    } catch (e) {
      console.error('保存失败', e)
      ElMessage.error('保存失败')
    } finally {
      saving.value = false
    }
  })
}

async function toggleOfficial(tag) {
  try {
    await api.put(`/admin/tags/${tag.id}`, { ...tag, is_official: !tag.is_official })
    tag.is_official = !tag.is_official
    ElMessage.success(tag.is_official ? '已设为官方推荐标签' : '已取消官方推荐')
  } catch (e) {
    console.error('操作失败', e)
    ElMessage.error('操作失败')
  }
}

async function toggleBan(tag) {
  const action = tag.is_banned ? '解禁' : '禁用'
  try {
    await ElMessageBox.confirm(`确定要${action}标签 "${tag.name}" 吗？`, `${action}标签`, {
      confirmButtonText: action,
      cancelButtonText: '取消',
      type: 'warning'
    })

    await api.put(`/admin/tags/${tag.id}`, { ...tag, is_banned: !tag.is_banned })
    tag.is_banned = !tag.is_banned
    ElMessage.success(`标签已${action}`)
  } catch (e) {
    if (e !== 'cancel') {
      console.error('操作失败', e)
      ElMessage.error('操作失败')
    }
  }
}

async function deleteTag(tag) {
  try {
    await ElMessageBox.confirm('删除后无法恢复，确定要删除吗？', '删除标签', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await api.delete(`/admin/tags/${tag.id}`)
    tags.value = tags.value.filter(t => t.id !== tag.id)
    ElMessage.success('标签已删除')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除失败', e)
      ElMessage.error('删除失败')
    }
  }
}

function openMergeModal(tag) {
  mergeSource.value = tag
  mergeTargetId.value = ''
  mergeDialogVisible.value = true
}

async function mergeTags() {
  if (!mergeTargetId.value) return

  saving.value = true
  try {
    await api.post('/admin/tags/merge', {
      source_id: mergeSource.value.id,
      target_id: parseInt(mergeTargetId.value)
    })
    ElMessage.success('标签合并成功')
    mergeDialogVisible.value = false
    loadTags()
  } catch (e) {
    console.error('合并失败', e)
    ElMessage.error('合并失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.tags-page {
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

.info-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(102, 126, 234, 0.08);
  color: #667eea;
  border-radius: 10px;
  margin-bottom: 16px;
  font-size: 13px;
}

.tag-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-icon {
  font-size: 18px;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
}

.merge-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 14px;
}

.emoji-preview {
  font-size: 16px;
}

:deep(.el-table) {
  --el-table-border-color: #f3f4f6;
  --el-table-header-bg-color: #f9fafb;
}

:deep(.el-table th) {
  font-weight: 600;
  color: #374151;
}

:deep(.el-table .banned-row) {
  background-color: rgba(254, 226, 226, 0.3);
}

:deep(.el-button.is-link) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
}

:deep(.el-button--purple) {
  color: #8b5cf6;
}
</style>
