export interface User {
  id: number;
  email: string;
  display_name: string;
  bio?: string;
  role: 'user' | 'admin';
  subscription_status: 'active' | 'inactive' | 'past_due' | 'canceled';
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  thumbnail_url?: string;
  author: User;
  status: 'pending' | 'approved' | 'rejected';
  replies?: Reply[];
  created_at: string;
  updated_at: string;
}

export interface Reply {
  id: number;
  content: string;
  post_id: number;
  author?: User;
  is_anonymous: boolean;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface SubscriptionStatus {
  status: 'active' | 'inactive' | 'past_due' | 'canceled';
  current_period_end?: string;
}