import React, { useState, useEffect, useCallback } from "react";
import { Link } from "react-router-dom";
import Layout from "../components/Layout";
import { postsAPI } from "../utils/api";
import type { Post } from "../types";
// import { css } from '../../styled-system/css';
// import { flex, grid } from '../../styled-system/patterns';

const Home: React.FC = () => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  const loadPosts = useCallback(async () => {
    try {
      setLoading(true);
      const response = await postsAPI.getPosts(page, 20);
      setPosts(response.posts);
      setTotal(response.total);
    } catch (err) {
      setError("投稿の読み込みに失敗しました");
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    loadPosts();
  }, [page, loadPosts]);

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
  const titleStyles = "text-3xl font-bold text-gray-900 mb-8 text-center";
  const emptyStyles =
    "text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300";
  const postsGridStyles =
    "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8";
  const postCardStyles = "bg-white rounded-xl shadow-md overflow-hidden";
  const postImageStyles = "w-full h-48 object-cover";
  const postContentStyles = "p-6";
  const postTitleStyles =
    "text-xl font-bold text-blue-600 no-underline mb-3 block";
  const postMetaStyles = "text-gray-500 text-sm mb-3";
  const postExcerptStyles = "text-gray-700 leading-relaxed overflow-hidden";
  const paginationStyles = "flex justify-center items-center gap-4 mt-8";
  const pageButtonStyles =
    "px-4 py-2 rounded-md border-none cursor-pointer text-sm font-medium";
  const pageButtonActiveStyles = "bg-blue-600 text-white";
  const pageInfoStyles = "text-gray-600 text-sm px-4";

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
        <h1 className={titleStyles}>投稿一覧</h1>
        {posts.length === 0 ? (
          <div className={emptyStyles}>
            <p>投稿がありません。</p>
            <p className="text-gray-500 text-sm mt-2">
              新しい投稿を作成してみましょう！
            </p>
          </div>
        ) : (
          <div className={postsGridStyles}>
            {posts.map((post) => (
              <article
                key={post.id}
                className={postCardStyles}
              >
                {post.thumbnail_url && (
                  <img
                    src={post.thumbnail_url}
                    alt={post.title}
                    className={postImageStyles}
                  />
                )}
                <div className={postContentStyles}>
                  <Link
                    to={`/posts/${post.id}`}
                    className={postTitleStyles}
                  >
                    {post.title}
                  </Link>
                  <div className={postMetaStyles}>
                    {post.created_at ? formatDate(post.created_at) : 'N/A'}
                  </div>
                  <p className={postExcerptStyles}>{post.content}</p>
                </div>
              </article>
            ))}
          </div>
        )}

        <div className={paginationStyles}>
          <button
            onClick={() => setPage(page - 1)}
            disabled={page <= 1}
            className={`${pageButtonStyles} ${page > 1 ? pageButtonActiveStyles : "bg-gray-300 text-gray-500"}`}
          >
            ← 前のページ
          </button>
          <span className={pageInfoStyles}>
            ページ {page} / {Math.ceil(total / 20)}
          </span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={page >= Math.ceil(total / 20)}
            className={`${pageButtonStyles} ${page < Math.ceil(total / 20) ? pageButtonActiveStyles : "bg-gray-300 text-gray-500"}`}
          >
            次のページ →
          </button>
        </div>
      </div>
    </Layout>
  );
};

export default Home;
