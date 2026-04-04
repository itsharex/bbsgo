<template>
  <div class="reputation-badge" :class="levelClass">
    <el-tooltip :content="tooltipContent" placement="top">
      <div class="flex items-center gap-1">
        <el-icon><Shield /></el-icon>
        <span class="text-sm font-medium">{{ reputation }}</span>
      </div>
    </el-tooltip>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Shield } from '@element-plus/icons-vue'

const { t } = useI18n()

const props = defineProps({
  reputation: {
    type: Number,
    default: 100
  }
})

const levelClass = computed(() => {
  if (props.reputation >= 80) return 'level-normal'
  if (props.reputation >= 60) return 'level-warning'
  if (props.reputation >= 40) return 'level-limited'
  if (props.reputation >= 20) return 'level-restricted'
  return 'level-banned'
})

const levelName = computed(() => {
  if (props.reputation >= 80) return t('reputation.good')
  if (props.reputation >= 60) return t('reputation.normal')
  if (props.reputation >= 40) return t('reputation.poor')
  if (props.reputation >= 20) return t('reputation.veryPoor')
  return t('reputation.banned')
})

const tooltipContent = computed(() => {
  return `${levelName.value} (${props.reputation}/100)`
})
</script>

<style scoped>
.reputation-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.level-normal {
  background-color: #f0f9ff;
  color: #0284c7;
}

.level-warning {
  background-color: #fef3c7;
  color: #d97706;
}

.level-limited {
  background-color: #fee2e2;
  color: #dc2626;
}

.level-restricted {
  background-color: #fecaca;
  color: #b91c1c;
}

.level-banned {
  background-color: #f3f4f6;
  color: #6b7280;
}
</style>
