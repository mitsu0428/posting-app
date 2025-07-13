import React, { useState, useEffect, useCallback } from "react";
import { useParams, Link } from "react-router-dom";
import Layout from "../components/Layout";
import { postsAPI } from "../utils/api";
import type { Post, Reply, CreateReplyRequest } from "../types";

const PostDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [post, setPost] = useState<Post | null>(null);
  const [replies, setReplies] = useState<Reply[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [replyContent, setReplyContent] = useState("");
  const [isAnonymous, setIsAnonymous] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const loadPost = useCallback(async () => {
    try {
      const postData = await postsAPI.getPost(Number(id));
      setPost(postData);
    } catch (err) {
      setError("投稿の読み込みに失敗しました");
    }
  }, [id]);

  const loadReplies = useCallback(async () => {
    try {
      const repliesData = await postsAPI.getReplies(Number(id));
      setReplies(repliesData);
    } catch (err) {
      console.error("返信の読み込みに失敗しました");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    if (id) {
      loadPost();
      loadReplies();
    }
  }, [id, loadPost, loadReplies]);

  const handleSubmitReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!replyContent.trim()) return;

    setSubmitting(true);
    try {
      const replyData: CreateReplyRequest = {
        content: replyContent,
        is_anonymous: isAnonymous,
      };

      await postsAPI.createReply(Number(id), replyData);
      setReplyContent("");
      setIsAnonymous(false);
      await loadReplies();
    } catch (err) {
      alert("返信の投稿に失敗しました");
    } finally {
      setSubmitting(false);
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

  // Basic styles
  const loadingStyles = "text-center py-12 text-gray-600 text-lg";
  const errorStyles = "text-red-600";
  const backLinkStyles = "text-blue-600 no-underline mb-4 inline-block";
  const postCardStyles = "border border-gray-300 rounded-lg p-8 mt-4 bg-white";
  const postTitleStyles = "text-3xl font-bold mb-4";
  const postMetaStyles = "text-gray-600 mb-4";
  const postImageStyles = "max-w-full h-auto mb-4";
  const postContentStyles = "leading-relaxed whitespace-pre-wrap";
  const replyFormContainerStyles = "mt-8";
  const replyFormStyles = "border border-gray-300 rounded-lg p-4 bg-gray-50";
  const textareaStyles = "w-full p-2 border border-gray-300 rounded border-solid resize-vertical";
  const checkboxContainerStyles = "mb-4";
  const checkboxStyles = "mr-2";
  const submitButtonStyles = "py-2 px-4 bg-blue-600 text-white border-none rounded cursor-pointer disabled:bg-gray-400 disabled:cursor-not-allowed";
  const repliesContainerStyles = "mt-8";
  const repliesTitleStyles = "text-2xl font-semibold mb-4";
  const replyCardStyles = "border border-gray-300 rounded-lg p-4 mb-4 bg-white";
  const replyHeaderStyles = "flex justify-between items-center mb-2";
  const replyAuthorStyles = "font-bold text-gray-600";
  const replyDateStyles = "text-gray-500 text-sm";
  const replyContentStyles = "leading-relaxed whitespace-pre-wrap";

  if (loading) {
    return (
      <Layout>
        <div className={loadingStyles}>読み込み中...</div>
      </Layout>
    );
  }

  if (error || !post) {
    return (
      <Layout>
        <div className={errorStyles}>{error || "投稿が見つかりません"}</div>
        <Link to="/home" className={backLinkStyles}>ホームに戻る</Link>
      </Layout>
    );
  }

  return (
    <Layout>
      <div>
        <Link to="/home" className={backLinkStyles}>
          ← 投稿一覧に戻る
        </Link>

        <div className={postCardStyles}>
          <h1 className={postTitleStyles}>{post.title}</h1>
          <p className={postMetaStyles}>
            投稿日: {post.created_at ? formatDate(post.created_at) : 'N/A'}
          </p>
          {post.thumbnail_url && (
            <img
              src={post.thumbnail_url}
              alt="サムネイル"
              className={postImageStyles}
            />
          )}
          <div className={postContentStyles}>
            {post.content}
          </div>
        </div>

        <div className={replyFormContainerStyles}>
          <h2>返信を投稿</h2>
          <form onSubmit={handleSubmitReply} className={replyFormStyles}>
            <div className="mb-4">
              <textarea
                value={replyContent}
                onChange={(e) => setReplyContent(e.target.value)}
                placeholder="返信内容を入力してください"
                rows={4}
                className={textareaStyles}
                required
              />
            </div>
            <div className={checkboxContainerStyles}>
              <label>
                <input
                  type="checkbox"
                  checked={isAnonymous}
                  onChange={(e) => setIsAnonymous(e.target.checked)}
                  className={checkboxStyles}
                />
                匿名で投稿する
              </label>
            </div>
            <button
              type="submit"
              disabled={submitting || !replyContent.trim()}
              className={submitButtonStyles}
            >
              {submitting ? "投稿中..." : "返信を投稿"}
            </button>
          </form>
        </div>

        <div className={repliesContainerStyles}>
          <h2 className={repliesTitleStyles}>返信一覧 ({replies.length}件)</h2>
          {replies.length === 0 ? (
            <p>まだ返信がありません。</p>
          ) : (
            <div>
              {replies.map((reply) => (
                <div key={reply.id} className={replyCardStyles}>
                  <div className={replyHeaderStyles}>
                    <span className={replyAuthorStyles}>
                      {reply.is_anonymous
                        ? "匿名ユーザー"
                        : `ユーザー#${reply.user_id}`}
                    </span>
                    <span className={replyDateStyles}>
                      {reply.created_at ? formatDate(reply.created_at) : 'N/A'}
                    </span>
                  </div>
                  <div className={replyContentStyles}>
                    {reply.content}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default PostDetail;
