import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Layout from "../components/Layout";
import { postsAPI } from "../utils/api";
import type { Post } from "../types";
// import { css } from '../../styled-system/css';
// import { flex } from '../../styled-system/patterns';

const MyPage: React.FC = () => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadUserPosts();
  }, []);

  const loadUserPosts = async () => {
    try {
      const userPosts = await postsAPI.getUserPosts();
      setPosts(userPosts);
    } catch (err) {
      setError("投稿の読み込みに失敗しました");
    } finally {
      setLoading(false);
    }
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
        return "warning.500";
      case "approved":
        return "success.500";
      case "rejected":
        return "danger.500";
      default:
        return "gray.500";
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

  // Temporary basic styles
  const loadingStyles = "text-center py-12 text-gray-600 text-lg";
  const errorStyles =
    "bg-red-50 border border-red-200 text-red-700 px-6 py-4 rounded-lg text-center";
  const titleStyles = "text-3xl font-bold text-gray-900 mb-8";
  const headerStyles = "flex justify-between items-center mb-8";
  const createButtonStyles =
    "bg-blue-600 text-white px-6 py-3 rounded-lg no-underline text-sm font-semibold";
  const sectionTitleStyles = "text-2xl font-semibold text-gray-800 mb-6";
  const emptyStateStyles =
    "text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300";
  const postCardStyles =
    "bg-white border border-gray-200 rounded-xl p-6 mb-4 shadow-sm";
  const postHeaderStyles = "flex justify-between items-start mb-4";
  const postTitleStyles = "text-xl font-semibold text-gray-900 no-underline";
  const statusBadgeStyles =
    "px-3 py-1 rounded-full text-white text-xs font-bold whitespace-nowrap";
  const postMetaStyles = "text-gray-500 text-sm mb-3";
  const postContentStyles =
    "text-gray-700 leading-relaxed overflow-hidden mb-3";
  const thumbnailStyles = "max-w-40 max-h-24 object-cover rounded-md mt-3";

  if (loading) {
    return (
      <Layout>
        <div className={loadingStyles}>読み込み中...</div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className={errorStyles}>{error}</div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div>
        <div className={headerStyles}>
          <h1 className={titleStyles}>マイページ</h1>
          <Link
            to="/create-post"
            className={createButtonStyles}
          >
            新規投稿
          </Link>
        </div>

        <h2 className={sectionTitleStyles}>あなたの投稿一覧</h2>
        {posts.length === 0 ? (
          <div className={emptyStateStyles}>
            <p className="mb-4 text-gray-600">まだ投稿がありません。</p>
            <Link
              to="/create-post"
              className="text-blue-600 no-underline"
            >
              最初の投稿を作成する
            </Link>
          </div>
        ) : (
          <div>
            {posts.map((post) => (
              <div
                key={post.id}
                className={postCardStyles}
              >
                <div className={postHeaderStyles}>
                  <h3 className="m-0 flex-1 pr-4">
                    {post.status === "approved" ? (
                      <Link
                        to={`/posts/${post.id}`}
                        className={postTitleStyles}
                      >
                        {post.title}
                      </Link>
                    ) : (
                      <span className="text-gray-700">{post.title}</span>
                    )}
                  </h3>
                  <span
                    className={statusBadgeStyles}
                    style={{
                      backgroundColor:
                        post.status && getStatusColor(post.status) === "warning.500"
                          ? "#f59e0b"
                          : post.status && getStatusColor(post.status) === "success.500"
                            ? "#10b981"
                            : "#ef4444",
                    }}
                  >
                    {post.status ? getStatusText(post.status) : 'Unknown'}
                  </span>
                </div>
                <p className={postMetaStyles}>
                  投稿日: {post.created_at ? formatDate(post.created_at) : 'N/A'}
                  {post.updated_at && post.created_at && post.updated_at !== post.created_at && (
                    <span> (更新日: {formatDate(post.updated_at)})</span>
                  )}
                </p>
                <p className={postContentStyles}>{post.content}</p>
                {post.thumbnail_url && (
                  <img
                    src={post.thumbnail_url}
                    alt="サムネイル"
                    className={thumbnailStyles}
                  />
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
};

export default MyPage;
