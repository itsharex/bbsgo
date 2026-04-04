<template>
  <el-dialog v-model="visible" :title="t('report.title')" width="500px">
    <el-form :model="form" label-width="80px">
      <el-form-item :label="t('report.reportReason')">
        <el-radio-group v-model="form.reason">
          <el-radio label="spam">{{ t('report.reasonSpam') }}</el-radio>
          <el-radio label="flood">{{ t('report.reasonFlood') }}</el-radio>
          <el-radio label="illegal">{{ t('report.reasonIllegal') }}</el-radio>
          <el-radio label="attack">{{ t('report.reasonAttack') }}</el-radio>
          <el-radio label="other">{{ t('report.reasonOther') }}</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="t('report.detail')" v-if="form.reason === 'other'">
        <el-input v-model="form.detail" type="textarea" :rows="3" :placeholder="t('report.reasonPlaceholder')" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">{{ t('report.cancel') }}</el-button>
      <el-button type="primary" @click="submitReport" :loading="loading">{{ t('report.submit') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'

const { t } = useI18n()

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
  reason: 'spam',
  detail: ''
})

const loading = ref(false)

async function submitReport() {
  loading.value = true
  try {
    const reason = form.value.reason === 'other' ? form.value.detail : t(`report.reason${form.value.reason.charAt(0).toUpperCase() + form.value.reason.slice(1)}`)
    await api.post('/reports', {
      target_type: props.targetType,
      target_id: props.targetId,
      reason: reason
    })
    ElMessage.success(t('report.reportSuccess'))
    visible.value = false
    emit('success')
    form.value = { reason: 'spam', detail: '' }
  } catch (e) {
    ElMessage.error(e.response?.data?.message || t('report.failed'))
  } finally {
    loading.value = false
  }
}
</script>
