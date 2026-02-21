<template>
  <div class="workflow-selector">
    <el-select
      :model-value="modelValue"
      @update:model-value="handleSelect"
      :placeholder="placeholder"
      filterable
      allow-create
      default-first-option
      style="width: 100%"
    >
      <el-option-group label="预设模型">
        <el-option
          v-for="preset in presetModels"
          :key="preset"
          :label="preset"
          :value="preset"
        />
      </el-option-group>
      <el-option-group label="自定义 Workflow" v-if="workflows.length > 0">
        <el-option
          v-for="workflow in workflows"
          :key="workflow.path"
          :label="workflow.name"
          :value="workflow.name"
        >
          <span>{{ workflow.name }}</span>
          <el-button
            type="danger"
            size="small"
            text
            @click.stop="handleDelete(workflow)"
            style="float: right"
          >
            删除
          </el-button>
        </el-option>
      </el-option-group>
    </el-select>

    <div class="workflow-actions">
      <el-upload
        :action="uploadUrl"
        :headers="uploadHeaders"
        :show-file-list="false"
        :before-upload="beforeUpload"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        accept=".json"
      >
        <el-button type="primary" size="small" :loading="uploading">
          <el-icon><Upload /></el-icon>
          上传 Workflow
        </el-button>
      </el-upload>
      <el-button size="small" @click="loadWorkflows" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新列表
      </el-button>
    </div>

    <div class="workflow-tip">
      支持上传 ComfyUI workflow JSON 文件，文件中可使用 %prompt% 等占位符
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Upload, Refresh } from '@element-plus/icons-vue';
import { workflowAPI, type Workflow } from '@/api/workflow';

const props = defineProps<{
  modelValue: string | string[];
  serviceType: 'image' | 'video';
  presetModels: string[];
  placeholder?: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string | string[]];
}>();

const workflows = ref<Workflow[]>([]);
const loading = ref(false);
const uploading = ref(false);

const uploadUrl = computed(() => {
  return `${import.meta.env.VITE_API_BASE_URL || ''}/api/v1/workflows/upload`;
});

const uploadHeaders = computed(() => {
  // 如果需要认证，在这里添加 token
  return {};
});

const loadWorkflows = async () => {
  loading.value = true;
  try {
    const response = await workflowAPI.list(props.serviceType);
    console.log('[WorkflowSelector] Loaded workflows:', response);
    workflows.value = response.workflows || [];
    if (workflows.value.length === 0) {
      console.log('[WorkflowSelector] No workflows found');
    }
  } catch (error: any) {
    console.error('[WorkflowSelector] Failed to load workflows:', error);
    ElMessage.error(error.message || '加载 workflow 列表失败');
  } finally {
    loading.value = false;
  }
};

const handleSelect = (value: string | string[]) => {
  emit('update:modelValue', value);
};

const beforeUpload = (file: File) => {
  const isJSON = file.name.endsWith('.json');
  const isLt10M = file.size / 1024 / 1024 < 10;

  if (!isJSON) {
    ElMessage.error('只能上传 JSON 文件！');
    return false;
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB！');
    return false;
  }

  uploading.value = true;
  return true;
};

const handleUploadSuccess = (response: any) => {
  uploading.value = false;
  ElMessage.success('Workflow 上传成功！');
  loadWorkflows();
  // 自动选择刚上传的文件（使用 filename 而不是 path）
  if (response.filename) {
    emit('update:modelValue', response.filename);
  }
};

const handleUploadError = (error: any) => {
  uploading.value = false;
  ElMessage.error('上传失败：' + (error.message || '未知错误'));
};

const handleDelete = async (workflow: Workflow) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 workflow "${workflow.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    );

    await workflowAPI.delete(workflow.name);
    ElMessage.success('删除成功');
    loadWorkflows();

    // 如果删除的是当前选中的，清空选择
    if (props.modelValue === workflow.path) {
      emit('update:modelValue', '');
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败');
    }
  }
};

onMounted(() => {
  loadWorkflows();
});
</script>

<style scoped>
.workflow-selector {
  width: 100%;
}

.workflow-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.workflow-tip {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 8px;
  line-height: 1.5;
}

:deep(.el-select-dropdown__item) {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
