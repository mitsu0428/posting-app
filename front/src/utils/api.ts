import { axiosInstance } from './api-mutator';

export interface LoginResponse {
  user: {
    id: number;
    email: string;
    display_name: string;
    bio?: string;
    role: 'user' | 'admin';
    subscription_status: 'active' | 'inactive' | 'past_due' | 'canceled';
    is_active: boolean;
    created_at: string;
    updated_at: string;
  };
  access_token: string;
}

export const authApi = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    const response = await axiosInstance.post('/auth/login', {
      email,
      password,
    });
    return response.data;
  },

  register: async (
    email: string,
    password: string,
    display_name: string
  ): Promise<any> => {
    const response = await axiosInstance.post('/auth/register', {
      email,
      password,
      display_name,
    });
    return response.data;
  },

  adminLogin: async (email: string, password: string): Promise<LoginResponse> => {
    const response = await axiosInstance.post('/admin/login', {
      email,
      password,
    });
    return response.data;
  },

  logout: async (): Promise<void> => {
    await axiosInstance.post('/auth/logout');
  },

  forgotPassword: async (email: string): Promise<any> => {
    const response = await axiosInstance.post('/auth/forgot-password', {
      email,
    });
    return response.data;
  },

  resetPassword: async (token: string, new_password: string): Promise<any> => {
    const response = await axiosInstance.post('/auth/reset-password', {
      token,
      new_password,
    });
    return response.data;
  },

  changePassword: async (
    current_password: string,
    new_password: string
  ): Promise<any> => {
    const response = await axiosInstance.post('/user/change-password', {
      current_password,
      new_password,
    });
    return response.data;
  },
};

export const userApi = {
  getProfile: async (): Promise<any> => {
    const response = await axiosInstance.get('/user/profile');
    return response.data;
  },

  updateProfile: async (
    display_name: string,
    bio?: string
  ): Promise<any> => {
    const response = await axiosInstance.put('/user/profile', {
      display_name,
      bio,
    });
    return response.data;
  },

  deactivateAccount: async (): Promise<any> => {
    const response = await axiosInstance.post('/user/deactivate');
    return response.data;
  },

  getUserPosts: async (page = 1, limit = 20): Promise<any> => {
    const response = await axiosInstance.get(
      `/user/posts?page=${page}&limit=${limit}`
    );
    return response.data;
  },
};

export const postApi = {
  getPosts: async (page = 1, limit = 20): Promise<any> => {
    const response = await axiosInstance.get(
      `/posts?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  getPost: async (id: number): Promise<any> => {
    const response = await axiosInstance.get(`/posts/${id}`);
    return response.data;
  },

  createPost: async (formData: FormData): Promise<any> => {
    const response = await axiosInstance.post('/posts', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  updatePost: async (id: number, formData: FormData): Promise<any> => {
    const response = await axiosInstance.put(`/posts/${id}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  deletePost: async (id: number): Promise<void> => {
    await axiosInstance.delete(`/posts/${id}`);
  },

  createReply: async (
    postId: number,
    content: string,
    isAnonymous: boolean
  ): Promise<any> => {
    const response = await axiosInstance.post(`/posts/${postId}/replies`, {
      content,
      is_anonymous: isAnonymous,
    });
    return response.data;
  },
};

export const subscriptionApi = {
  getStatus: async (): Promise<any> => {
    const response = await axiosInstance.get('/subscription/status');
    return response.data;
  },

  createCheckoutSession: async (): Promise<{ url: string }> => {
    const response = await axiosInstance.post(
      '/subscription/create-checkout-session'
    );
    return response.data;
  },
};

export const adminApi = {
  getPosts: async (
    page = 1,
    limit = 20,
    status?: string
  ): Promise<any> => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    if (status) {
      params.append('status', status);
    }
    const response = await axiosInstance.get(`/admin/posts?${params}`);
    return response.data;
  },

  approvePost: async (id: number): Promise<any> => {
    const response = await axiosInstance.post(`/admin/posts/${id}/approve`);
    return response.data;
  },

  rejectPost: async (id: number): Promise<any> => {
    const response = await axiosInstance.post(`/admin/posts/${id}/reject`);
    return response.data;
  },

  getUsers: async (page = 1, limit = 20): Promise<any> => {
    const response = await axiosInstance.get(
      `/admin/users?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  banUser: async (id: number): Promise<any> => {
    const response = await axiosInstance.post(`/admin/users/${id}/ban`);
    return response.data;
  },
};

export const groupApi = {
  deleteGroup: async (id: number): Promise<{ message: string }> => {
    const response = await axiosInstance.delete(`/groups/${id}`);
    return response.data;
  },
  
  getGroupMembers: async (id: number): Promise<Array<{ id: number; display_name: string; bio: string }>> => {
    const response = await axiosInstance.get(`/groups/${id}/members`);
    return response.data;
  },
  
  removeGroupMember: async (groupId: number, memberId: number): Promise<{ message: string }> => {
    const response = await axiosInstance.delete(`/groups/${groupId}/members/${memberId}`);
    return response.data;
  },
  
  leaveGroup: async (id: number): Promise<{ message: string }> => {
    const response = await axiosInstance.post(`/groups/${id}/leave`);
    return response.data;
  },
};