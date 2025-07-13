// 生成された型定義を再エクスポート
export * from '../generated/models';

// 追加のコンテキスト型定義（生成された型を使用）
export interface AuthContextType {
  user: import('../generated/models').User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  adminLogin: (email: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

// 旧来の手動型定義は削除し、生成された型を使用
// 既存のコンポーネントとの互換性のため、必要に応じて型エイリアスを作成
export type { User, Post, Reply, CreatePostRequest, CreateReplyRequest } from '../generated/models';