# ComfyUI 界面显示修复

## 问题描述

在 AI 配置页面添加视频生成配置时，厂商下拉列表中看不到 ComfyUI 选项。

## 原因分析

原始代码中，`availableProviders` 计算属性只显示已有激活配置的厂商：

```typescript
// 原始逻辑（有问题）
const availableProviders = computed(() => {
  // 只返回已有激活配置的厂商
  const activeConfigs = configs.value.filter(
    (c) => c.service_type === form.service_type && c.is_active,
  );
  const activeProviderIds = new Set(activeConfigs.map((c) => c.provider));
  const allProviders = providerConfigs[form.service_type] || [];
  return allProviders.filter((p) => activeProviderIds.has(p.id));
});
```

这导致：
- 如果某个厂商（如 ComfyUI）还没有任何配置，就不会显示在下拉列表中
- 用户无法创建第一个 ComfyUI 配置

## 解决方案

修改 `availableProviders` 和 `availableModels` 的逻辑，让所有在 `providerConfigs` 中定义的厂商都显示出来。

### 修改文件

`web/src/views/settings/AIConfig.vue`

### 修改内容

#### 1. 修改 availableProviders

```typescript
// 新逻辑（已修复）
const availableProviders = computed(() => {
  // 返回当前 service_type 下的所有厂商
  return providerConfigs[form.service_type] || [];
});
```

#### 2. 修改 availableModels

```typescript
// 新逻辑（已修复）
const availableModels = computed(() => {
  if (!form.provider) return [];

  // 优先从 providerConfigs 中查找该 provider 的模型列表
  const allProviders = providerConfigs[form.service_type] || [];
  const providerConfig = allProviders.find((p) => p.id === form.provider);
  
  if (providerConfig) {
    return providerConfig.models;
  }

  // 如果在 providerConfigs 中找不到，则从已激活的配置中提取
  const activeConfigsForProvider = configs.value.filter(
    (c) =>
      c.provider === form.provider &&
      c.service_type === form.service_type &&
      c.is_active,
  );

  const models = new Set<string>();
  activeConfigsForProvider.forEach((config) => {
    config.model.forEach((m) => models.add(m));
  });

  return Array.from(models);
});
```

## 效果

修复后，在视频生成配置页面的厂商下拉列表中，现在会显示：

- ✅ 火山引擎
- ✅ Chatfire
- ✅ OpenAI
- ✅ **ComfyUI** ← 新增可见
- ✅ MiniMax（如果取消注释）

用户现在可以：
1. 选择 ComfyUI 作为厂商
2. 看到预定义的模型列表（svd, svd_xt, custom）
3. 成功创建第一个 ComfyUI 配置

## 测试步骤

1. 启动服务
2. 访问 **设置 → AI 服务配置 → 视频生成**
3. 点击 **添加配置**
4. 打开厂商下拉列表
5. 确认可以看到 ComfyUI 选项
6. 选择 ComfyUI
7. 确认模型列表显示：svd, svd_xt, custom
8. 填写配置并保存

## 相关文件

- `web/src/views/settings/AIConfig.vue` - 主要修改文件
- `web/dist/` - 重新构建的前端文件

## 版本信息

- 修复日期：2026-02-20
- 影响版本：v1.0.6
- 修复类型：UI 显示问题

## 注意事项

此修改不影响：
- 后端逻辑
- 现有配置
- 其他厂商的显示和功能

只是让所有预定义的厂商都能在界面上显示，方便用户创建配置。
