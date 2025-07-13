import axios from 'axios';
import type { User, Post, Reply, CreatePostRequest, CreateReplyRequest } from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  register: async (username: string, email: string, password: string): Promise<User> => {
    const response = await api.post('/auth/register', { username, email, password });
    return response.data;
  },

  login: async (email: string, password: string): Promise<{ token: string; user: User }> => {
    const response = await api.post('/auth/login', { email, password });
    return response.data;
  },

  adminLogin: async (email: string, password: string): Promise<{ token: string; user: User }> => {
    const response = await api.post('/admin/login', { email, password });
    return response.data;
  },

  logout: async (): Promise<void> => {
    await api.post('/auth/logout');
  },

  forgotPassword: async (email: string): Promise<{ message: string }> => {
    const response = await api.post('/auth/forgot-password', { email });
    return response.data;
  },

  resetPassword: async (token: string, newPassword: string): Promise<{ message: string }> => {
    const response = await api.post('/auth/reset-password', { token, new_password: newPassword });
    return response.data;
  },
};

export const postsAPI = {
  getPosts: async (page: number = 1, limit: number = 20): Promise<{ posts: Post[]; total: number }> => {
    const response = await api.get(`/posts?page=${page}&limit=${limit}`);
    return response.data;
  },

  getPost: async (id: number): Promise<Post> => {
    const response = await api.get(`/posts/${id}`);
    return response.data;
  },

  createPost: async (postData: CreatePostRequest): Promise<Post> => {
    const response = await api.post('/posts', postData);
    return response.data;
  },

  getUserPosts: async (): Promise<Post[]> => {
    const response = await api.get('/me/posts');
    return response.data;
  },

  getReplies: async (postId: number): Promise<Reply[]> => {
    const response = await api.get(`/posts/${postId}/replies`);
    return response.data;
  },

  createReply: async (postId: number, replyData: CreateReplyRequest): Promise<Reply> => {
    const response = await api.post(`/posts/${postId}/replies`, replyData);
    return response.data;
  },
};

export const adminAPI = {
  getPosts: async (status?: string): Promise<Post[]> => {
    const url = status ? `/admin/posts?status=${status}` : '/admin/posts';
    const response = await api.get(url);
    return response.data;
  },

  approvePost: async (id: number): Promise<void> => {
    await api.post(`/admin/posts/${id}/approve`);
  },

  rejectPost: async (id: number): Promise<void> => {
    await api.post(`/admin/posts/${id}/reject`);
  },

  deletePost: async (id: number): Promise<void> => {
    await api.delete(`/admin/posts/${id}`);
  },

  getUsers: async (): Promise<User[]> => {
    const response = await api.get('/admin/users');
    return response.data;
  },

  deactivateUser: async (id: number): Promise<void> => {
    await api.post(`/admin/users/${id}/deactivate`);
  },
};

export const subscriptionAPI = {
  createCheckoutSession: async (): Promise<{ session_id: string }> => {
    const response = await api.post('/subscription/create-checkout-session');
    return response.data;
  },
};