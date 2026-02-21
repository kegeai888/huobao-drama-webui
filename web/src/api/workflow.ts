import request from '@/utils/request';

export interface Workflow {
  name: string;
  path: string;
  size: number;
  modified_at: string;
}

export interface WorkflowListResponse {
  workflows: Workflow[];
  count: number;
}

export const workflowAPI = {
  // 获取 workflow 列表
  list: (type?: 'image' | 'video'): Promise<WorkflowListResponse> => {
    return request.get('/workflows', { params: { type } });
  },

  // 上传 workflow 文件
  upload: (file: File): Promise<{ message: string; filename: string; path: string }> => {
    const formData = new FormData();
    formData.append('file', file);
    return request.post('/workflows/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },

  // 获取 workflow 内容
  get: (filename: string): Promise<any> => {
    return request.get(`/workflows/${filename}`);
  },

  // 删除 workflow
  delete: (filename: string): Promise<{ message: string }> => {
    return request.delete(`/workflows/${filename}`);
  },
};
