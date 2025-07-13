import React, { useState, useEffect, useCallback } from "react";
import { useAuth } from "../context/AuthContext";
import { adminAPI } from "../utils/api";
import type { Post, User } from "../types";

const AdminDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<"posts" | "users">("posts");
  const [posts, setPosts] = useState<Post[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("pending");
  const { logout } = useAuth();

  const loadPosts = useCallback(async () => {
    setLoading(true);
    try {
      const postsData = await adminAPI.getPosts(statusFilter || undefined);
      setPosts(postsData);
    } catch (err) {
      setError("投稿の読み込みに失敗しました");
    } finally {
      setLoading(false);
    }
  }, [statusFilter]);

  const loadUsers = useCallback(async () => {
    setLoading(true);
    try {
      const usersData = await adminAPI.getUsers();
      setUsers(usersData);
    } catch (err) {
      setError("ユーザーの読み込みに失敗しました");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (activeTab === "posts") {
      loadPosts();
    } else {
      loadUsers();
    }
  }, [activeTab, statusFilter, loadPosts, loadUsers]);

  const handleApprovePost = async (id: number) => {
    try {
      await adminAPI.approvePost(id);
      alert("投稿を承認しました");
      loadPosts();
    } catch (err) {
      alert("承認に失敗しました");
    }
  };

  const handleRejectPost = async (id: number) => {
    if (!window.confirm("本当に却下しますか？")) return;
    try {
      await adminAPI.rejectPost(id);
      alert("投稿を却下しました");
      loadPosts();
    } catch (err) {
      alert("却下に失敗しました");
    }
  };

  const handleDeletePost = async (id: number) => {
    if (!window.confirm("本当に削除しますか？この操作は取り消せません。"))
      return;
    try {
      await adminAPI.deletePost(id);
      alert("投稿を削除しました");
      loadPosts();
    } catch (err) {
      alert("削除に失敗しました");
    }
  };

  const handleDeactivateUser = async (id: number) => {
    if (!window.confirm("本当にこのユーザーを無効化しますか？")) return;
    try {
      await adminAPI.deactivateUser(id);
      alert("ユーザーを無効化しました");
      loadUsers();
    } catch (err) {
      alert("無効化に失敗しました");
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("ja-JP", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case "pending":
        return "承認待ち";
      case "approved":
        return "承認済み";
      case "rejected":
        return "却下";
      default:
        return status;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "pending":
        return "#ffc107";
      case "approved":
        return "#28a745";
      case "rejected":
        return "#dc3545";
      default:
        return "#6c757d";
    }
  };

  return (
    <div style={{ minHeight: "100vh", backgroundColor: "#f8f9fa" }}>
      <nav
        style={{
          backgroundColor: "#dc3545",
          padding: "1rem",
          color: "white",
        }}
      >
        <div
          style={{
            maxWidth: "1200px",
            margin: "0 auto",
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <h1 style={{ margin: 0 }}>管理者ダッシュボード</h1>
          <button
            onClick={logout}
            style={{
              padding: "0.5rem 1rem",
              backgroundColor: "#fff",
              color: "#dc3545",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            }}
          >
            ログアウト
          </button>
        </div>
      </nav>

      <div style={{ maxWidth: "1200px", margin: "0 auto", padding: "2rem" }}>
        <div style={{ marginBottom: "2rem" }}>
          <button
            onClick={() => setActiveTab("posts")}
            style={{
              padding: "0.75rem 1.5rem",
              backgroundColor: activeTab === "posts" ? "#007bff" : "#6c757d",
              color: "white",
              border: "none",
              borderRadius: "4px 0 0 4px",
              cursor: "pointer",
            }}
          >
            投稿管理
          </button>
          <button
            onClick={() => setActiveTab("users")}
            style={{
              padding: "0.75rem 1.5rem",
              backgroundColor: activeTab === "users" ? "#007bff" : "#6c757d",
              color: "white",
              border: "none",
              borderRadius: "0 4px 4px 0",
              cursor: "pointer",
            }}
          >
            ユーザー管理
          </button>
        </div>

        {error && (
          <div
            style={{
              color: "red",
              marginBottom: "1rem",
              padding: "1rem",
              backgroundColor: "#f8d7da",
              borderRadius: "4px",
            }}
          >
            {error}
          </div>
        )}

        {activeTab === "posts" && (
          <div>
            <div style={{ marginBottom: "1rem" }}>
              <label htmlFor="statusFilter">ステータスフィルター:</label>
              <select
                id="statusFilter"
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                style={{ marginLeft: "0.5rem", padding: "0.5rem" }}
              >
                <option value="">すべて</option>
                <option value="pending">承認待ち</option>
                <option value="approved">承認済み</option>
                <option value="rejected">却下</option>
              </select>
            </div>

            {loading ? (
              <div>読み込み中...</div>
            ) : (
              <div>
                <h2>投稿一覧 ({posts.length}件)</h2>
                {posts.length === 0 ? (
                  <p>該当する投稿がありません。</p>
                ) : (
                  <div>
                    {posts.map((post) => (
                      <div
                        key={post.id}
                        style={{
                          border: "1px solid #ddd",
                          borderRadius: "8px",
                          padding: "1rem",
                          marginBottom: "1rem",
                          backgroundColor: "white",
                        }}
                      >
                        <div
                          style={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "flex-start",
                            marginBottom: "1rem",
                          }}
                        >
                          <div>
                            <h3 style={{ margin: "0 0 0.5rem 0" }}>
                              {post.title}
                            </h3>
                            <p
                              style={{ color: "#666", margin: "0 0 0.5rem 0" }}
                            >
                              投稿者ID: {post.user_id} | 投稿日:{" "}
                              {post.created_at ? formatDate(post.created_at) : 'N/A'}
                            </p>
                            <span
                              style={{
                                padding: "0.25rem 0.5rem",
                                borderRadius: "4px",
                                color: "white",
                                backgroundColor: post.status ? getStatusColor(post.status) : '#gray',
                                fontSize: "0.8rem",
                              }}
                            >
                              {post.status ? getStatusText(post.status) : 'Unknown'}
                            </span>
                          </div>
                          <div style={{ display: "flex", gap: "0.5rem" }}>
                            {post.status === "pending" && (
                              <>
                                <button
                                  onClick={() => post.id && handleApprovePost(post.id)}
                                  style={{
                                    padding: "0.5rem 1rem",
                                    backgroundColor: "#28a745",
                                    color: "white",
                                    border: "none",
                                    borderRadius: "4px",
                                    cursor: "pointer",
                                  }}
                                >
                                  承認
                                </button>
                                <button
                                  onClick={() => post.id && handleRejectPost(post.id)}
                                  style={{
                                    padding: "0.5rem 1rem",
                                    backgroundColor: "#ffc107",
                                    color: "black",
                                    border: "none",
                                    borderRadius: "4px",
                                    cursor: "pointer",
                                  }}
                                >
                                  却下
                                </button>
                              </>
                            )}
                            <button
                              onClick={() => post.id && handleDeletePost(post.id)}
                              style={{
                                padding: "0.5rem 1rem",
                                backgroundColor: "#dc3545",
                                color: "white",
                                border: "none",
                                borderRadius: "4px",
                                cursor: "pointer",
                              }}
                            >
                              削除
                            </button>
                          </div>
                        </div>
                        <p
                          style={{
                            margin: 0,
                            overflow: "hidden",
                            textOverflow: "ellipsis",
                            display: "-webkit-box",
                            WebkitLineClamp: 3,
                            WebkitBoxOrient: "vertical",
                          }}
                        >
                          {post.content}
                        </p>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {activeTab === "users" && (
          <div>
            {loading ? (
              <div>読み込み中...</div>
            ) : (
              <div>
                <h2>ユーザー一覧 ({users.length}件)</h2>
                {users.length === 0 ? (
                  <p>ユーザーがいません。</p>
                ) : (
                  <div style={{ overflowX: "auto" }}>
                    <table
                      style={{
                        width: "100%",
                        borderCollapse: "collapse",
                        backgroundColor: "white",
                      }}
                    >
                      <thead>
                        <tr style={{ backgroundColor: "#f8f9fa" }}>
                          <th
                            style={{
                              padding: "1rem",
                              textAlign: "left",
                              border: "1px solid #ddd",
                            }}
                          >
                            ID
                          </th>
                          <th
                            style={{
                              padding: "1rem",
                              textAlign: "left",
                              border: "1px solid #ddd",
                            }}
                          >
                            ユーザー名
                          </th>
                          <th
                            style={{
                              padding: "1rem",
                              textAlign: "left",
                              border: "1px solid #ddd",
                            }}
                          >
                            メール
                          </th>
                          <th
                            style={{
                              padding: "1rem",
                              textAlign: "left",
                              border: "1px solid #ddd",
                            }}
                          >
                            サブスク状態
                          </th>
                          <th
                            style={{
                              padding: "1rem",
                              textAlign: "left",
                              border: "1px solid #ddd",
                            }}
                          >
                            管理者
                          </th>
                          <th
                            style={{
                              padding: "1rem",
                              textAlign: "left",
                              border: "1px solid #ddd",
                            }}
                          >
                            登録日
                          </th>
                          <th
                            style={{
                              padding: "1rem",
                              textAlign: "left",
                              border: "1px solid #ddd",
                            }}
                          >
                            操作
                          </th>
                        </tr>
                      </thead>
                      <tbody>
                        {users.map((user) => (
                          <tr key={user.id}>
                            <td
                              style={{
                                padding: "1rem",
                                border: "1px solid #ddd",
                              }}
                            >
                              {user.id}
                            </td>
                            <td
                              style={{
                                padding: "1rem",
                                border: "1px solid #ddd",
                              }}
                            >
                              {user.username}
                            </td>
                            <td
                              style={{
                                padding: "1rem",
                                border: "1px solid #ddd",
                              }}
                            >
                              {user.email}
                            </td>
                            <td
                              style={{
                                padding: "1rem",
                                border: "1px solid #ddd",
                              }}
                            >
                              <span
                                style={{
                                  padding: "0.25rem 0.5rem",
                                  borderRadius: "4px",
                                  color: "white",
                                  backgroundColor:
                                    user.subscription_status === "active"
                                      ? "#28a745"
                                      : "#dc3545",
                                  fontSize: "0.8rem",
                                }}
                              >
                                {user.subscription_status === "active"
                                  ? "有効"
                                  : "無効"}
                              </span>
                            </td>
                            <td
                              style={{
                                padding: "1rem",
                                border: "1px solid #ddd",
                              }}
                            >
                              {user.is_admin ? "管理者" : "一般"}
                            </td>
                            <td
                              style={{
                                padding: "1rem",
                                border: "1px solid #ddd",
                              }}
                            >
                              {user.created_at ? formatDate(user.created_at) : 'N/A'}
                            </td>
                            <td
                              style={{
                                padding: "1rem",
                                border: "1px solid #ddd",
                              }}
                            >
                              {!user.is_admin && (
                                <button
                                  onClick={() => user.id && handleDeactivateUser(user.id)}
                                  style={{
                                    padding: "0.5rem 1rem",
                                    backgroundColor: "#dc3545",
                                    color: "white",
                                    border: "none",
                                    borderRadius: "4px",
                                    cursor: "pointer",
                                  }}
                                >
                                  無効化
                                </button>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default AdminDashboard;
