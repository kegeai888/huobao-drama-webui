<template>
  <div class="app-layout">
    <!-- Global Header -->
    <header class="app-header">
      <div class="header-content">
        <div class="header-left">
          <router-link to="/" class="logo">
            <span class="logo-text">🎬 HuoBao Drama</span>
          </router-link>
        </div>
        <div class="header-right">
          <div class="contributor-info">
            <span class="contributor-text">ComfyUI 贡献者: 你们喜爱的老王</span>
            <a href="https://space.bilibili.com/97727630" target="_blank" class="bilibili-btn">
              <el-icon><VideoPlay /></el-icon>
              <span>B站</span>
            </a>
          </div>
          <LanguageSwitcher />
          <ThemeToggle />
          <el-button @click="showAIConfig = true" class="header-btn">
            <el-icon><Setting /></el-icon>
            <span class="btn-text">{{ $t('drama.aiConfig') }}</span>
          </el-button>
          <!-- <el-button :icon="Setting" circle @click="showAIConfig = true" :title="$t('aiConfig.title')" /> -->
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="app-main">
      <slot />
    </main>

    <!-- AI Config Dialog -->
    <AIConfigDialog v-model="showAIConfig" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Setting, VideoPlay } from '@element-plus/icons-vue'
import ThemeToggle from './ThemeToggle.vue'
import AIConfigDialog from './AIConfigDialog.vue'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'

const showAIConfig = ref(false)
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border-primary);
  backdrop-filter: blur(8px);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-4);
  height: 56px;
  max-width: 100%;
  margin: 0 auto;
}
.header-btn {
  border-radius: var(--radius-lg);
  font-weight: 500;
}

.header-btn.primary {
  background: linear-gradient(135deg, var(--accent) 0%, #0284c7 100%);
  border: none;
  box-shadow: 0 4px 14px rgba(14, 165, 233, 0.35);
}

.header-btn.primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(14, 165, 233, 0.45);
}
.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  text-decoration: none;
  color: var(--text-primary);
  font-weight: 700;
  font-size: 1.125rem;
  transition: opacity var(--transition-fast);
}

.logo:hover {
  opacity: 0.8;
}

.logo-text {
  background: linear-gradient(135deg, var(--accent) 0%, #06b6d4 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.contributor-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
  border: 1px solid rgba(14, 165, 233, 0.2);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.contributor-info:hover {
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.15) 0%, rgba(6, 182, 212, 0.15) 100%);
  border-color: rgba(14, 165, 233, 0.3);
}

.contributor-text {
  font-size: 14px;
  color: #333;
  font-weight: 500;
  white-space: nowrap;
}

.dark .contributor-text {
  color: #e5e7eb;
}

.bilibili-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: #00a1d6;
  color: white !important;
  border-radius: 6px;
  text-decoration: none;
  font-size: 12px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.bilibili-btn:hover {
  background: #0087b3;
  transform: scale(1.05);
}

.app-main {
  flex: 1;
}

/* Dark mode adjustments */
.dark .app-header {
  background: rgba(26, 33, 41, 0.95);
}

.dark .contributor-info {
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.15) 0%, rgba(6, 182, 212, 0.15) 100%);
  border-color: rgba(14, 165, 233, 0.3);
}

/* Responsive */
@media (max-width: 1024px) {
  .contributor-text {
    display: none;
  }
}

@media (max-width: 768px) {
  .contributor-info {
    padding: 4px 8px;
  }
  
  .btn-text {
    display: none;
  }
}
</style>
