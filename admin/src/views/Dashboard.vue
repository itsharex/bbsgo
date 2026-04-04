<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card stat-users">
        <div class="stat-icon">
          <UserPlus :size="24" />
        </div>
        <div class="stat-content">
          <span class="stat-label">{{ t('dashboard.users') }}</span>
          <span class="stat-value">{{ stats.users }}</span>
        </div>
        <div class="stat-trend up">
          <TrendingUp :size="14" />
          <span>+12%</span>
        </div>
      </div>

      <div class="stat-card stat-topics">
        <div class="stat-icon">
          <FileText :size="24" />
        </div>
        <div class="stat-content">
          <span class="stat-label">{{ t('dashboard.topics') }}</span>
          <span class="stat-value">{{ stats.topics }}</span>
        </div>
        <div class="stat-trend up">
          <TrendingUp :size="14" />
          <span>+8%</span>
        </div>
      </div>

      <div class="stat-card stat-comments">
        <div class="stat-icon">
          <MessageCircle :size="24" />
        </div>
        <div class="stat-content">
          <span class="stat-label">{{ t('dashboard.comments') }}</span>
          <span class="stat-value">{{ stats.comments }}</span>
        </div>
        <div class="stat-trend up">
          <TrendingUp :size="14" />
          <span>+23%</span>
        </div>
      </div>

      <div class="stat-card stat-reports">
        <div class="stat-icon">
          <AlertTriangle :size="24" />
        </div>
        <div class="stat-content">
          <span class="stat-label">{{ t('dashboard.pendingReports') }}</span>
          <span class="stat-value">{{ stats.reports }}</span>
        </div>
        <div class="stat-trend down">
          <TrendingDown :size="14" />
          <span>-5%</span>
        </div>
      </div>
    </div>

    <!-- 系统信息 -->
    <div class="info-section">
      <div class="section-header">
        <h3 class="section-title">
          <Info :size="18" />
          {{ t('dashboard.systemInfo') }}
        </h3>
      </div>
      <div class="info-grid">
        <div class="info-item">
          <div class="info-icon blue">
            <Codepen :size="18" />
          </div>
          <div class="info-content">
            <span class="info-label">{{ t('dashboard.version') }}</span>
            <span class="info-value">v1.0.0</span>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon green">
            <Box :size="18" />
          </div>
          <div class="info-content">
            <span class="info-label">{{ t('dashboard.goVersion') }}</span>
            <span class="info-value">go1.21+</span>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon purple">
            <Database :size="18" />
          </div>
          <div class="info-content">
            <span class="info-label">{{ t('dashboard.database') }}</span>
            <span class="info-value">SQLite</span>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon orange">
            <Zap :size="18" />
          </div>
          <div class="info-content">
            <span class="info-label">{{ t('dashboard.cache') }}</span>
            <span class="info-value">Ristretto</span>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon cyan">
            <Cpu :size="18" />
          </div>
          <div class="info-content">
            <span class="info-label">{{ t('dashboard.uptime') }}</span>
            <span class="info-value">{{ uptime }}</span>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon pink">
            <Activity :size="18" />
          </div>
          <div class="info-content">
            <span class="info-label">{{ t('dashboard.status') }}</span>
            <span class="info-value status-online">{{ t('dashboard.online') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="quick-actions">
      <div class="section-header">
        <h3 class="section-title">
          <Rocket :size="18" />
          {{ t('dashboard.quickActions') }}
        </h3>
      </div>
      <div class="actions-grid">
        <router-link to="/users" class="action-card">
          <div class="action-icon purple">
            <UserPlus :size="22" />
          </div>
          <span class="action-text">{{ t('dashboard.userManagement') }}</span>
        </router-link>
        <router-link to="/topics" class="action-card">
          <div class="action-icon green">
            <FilePlus :size="22" />
          </div>
          <span class="action-text">{{ t('dashboard.postReview') }}</span>
        </router-link>
        <router-link to="/reports" class="action-card">
          <div class="action-icon red">
            <ShieldAlert :size="22" />
          </div>
          <span class="action-text">{{ t('dashboard.handleReports') }}</span>
        </router-link>
        <router-link to="/announcements" class="action-card">
          <div class="action-icon yellow">
            <Megaphone :size="22" />
          </div>
          <span class="action-text">{{ t('dashboard.publishAnnouncement') }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  UserPlus, FileText, MessageCircle, AlertTriangle, TrendingUp, TrendingDown,
  Info, Codepen, Box, Database, Zap, Cpu, Activity, Rocket, FilePlus,
  ShieldAlert, Megaphone
} from 'lucide-vue-next'

const { t } = useI18n()
const stats = ref({
  users: 0,
  topics: 0,
  comments: 0,
  reports: 0
})

const uptime = ref('0 days 0 hours')

onMounted(() => {
  stats.value = {
    users: 156,
    topics: 1024,
    comments: 3567,
    reports: 5
  }
  uptime.value = '15 days 8 hours'
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
}

.stat-users::before { background: linear-gradient(90deg, #667eea, #764ba2); }
.stat-topics::before { background: linear-gradient(90deg, #34d399, #10b981); }
.stat-comments::before { background: linear-gradient(90deg, #22d3ee, #06b6d4); }
.stat-reports::before { background: linear-gradient(90deg, #f87171, #ef4444); }

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-users .stat-icon { background: rgba(102, 126, 234, 0.1); color: #667eea; }
.stat-topics .stat-icon { background: rgba(52, 211, 153, 0.1); color: #34d399; }
.stat-comments .stat-icon { background: rgba(34, 211, 238, 0.1); color: #22d3ee; }
.stat-reports .stat-icon { background: rgba(248, 113, 113, 0.1); color: #f87171; }

.stat-content {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.stat-label {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 8px;
  position: absolute;
  top: 16px;
  right: 16px;
}

.stat-trend.up {
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
}

.stat-trend.down {
  background: rgba(248, 113, 113, 0.1);
  color: #f87171;
}

/* 信息区域 */
.info-section, .quick-actions {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
}

.section-header {
  margin-bottom: 20px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
}

.info-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.info-icon.blue { background: rgba(102, 126, 234, 0.1); color: #667eea; }
.info-icon.green { background: rgba(52, 211, 153, 0.1); color: #34d399; }
.info-icon.purple { background: rgba(171, 122, 224, 0.1); color: #c084fc; }
.info-icon.orange { background: rgba(251, 146, 60, 0.1); color: #fb923c; }
.info-icon.cyan { background: rgba(34, 211, 238, 0.1); color: #22d3ee; }
.info-icon.pink { background: rgba(244, 114, 182, 0.1); color: #f472b6; }

.info-content {
  display: flex;
  flex-direction: column;
}

.info-label {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 2px;
}

.info-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.info-value.status-online {
  color: #34d399;
}

/* 快捷操作 */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 16px;
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px 16px;
  background: #f9fafb;
  border-radius: 12px;
  text-decoration: none;
  transition: all 0.2s;
}

.action-card:hover {
  background: #f3f4f6;
  transform: translateY(-2px);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-icon.purple { background: rgba(171, 122, 224, 0.15); color: #c084fc; }
.action-icon.green { background: rgba(52, 211, 153, 0.15); color: #34d399; }
.action-icon.red { background: rgba(248, 113, 113, 0.15); color: #f87171; }
.action-icon.yellow { background: rgba(251, 191, 36, 0.15); color: #fbbf24; }

.action-text {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}
</style>
