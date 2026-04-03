<template>
  <el-dialog v-model="visible" title="举报内容" width="500px">
    <el-form :model="form" label-width="80px">
      <el-form-item label="举报原因">
        <el-radio-group v-model="form.reason">
          <el-radio label="垃圾广告">垃圾广告</el-radio>
          <el-radio label="恶意灌水">恶意灌水</el-radio>
          <el-radio label="违法违规">违法违规</el-radio>
          <el-radio label="人身攻击">人身攻击</el-radio>
          <el-radio label="其他">其他</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="详细说明" v-if="form.reason === '其他'">
        <el-input v-model="form.detail" type="textarea" :rows="3" placeholder="请详细说明举报原因" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="submitReport" :loading="loading">提交举报</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/api'

const props = defineProps({
  modelValue: Boolean,
  targetType: {
    type: String,
    required: true,
    validator: (val) => ['topic', 'comment'].includes(val)
  },
  targetId: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const form = ref({
  reason: '垃圾广告',
  detail: ''
})

const loading = ref(false)

async function submitReport() {
  loading.value = true
  try {
    const reason = form.value.reason === '其他' ? form.value.detail : form.value.reason
    await api.post('/reports', {
      target_type: props.targetType,
      target_id: props.targetId,
      reason: reason
    })
    ElMessage.success('举报成功，我们会尽快处理')
    visible.value = false
    emit('success')
    form.value = { reason: '垃圾广告', detail: '' }
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '举报失败')
  } finally {
    loading.value = false
  }
}
</script>
