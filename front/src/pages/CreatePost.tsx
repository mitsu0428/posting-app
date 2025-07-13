import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Layout from "../components/Layout";
import { postsAPI } from "../utils/api";
import type { CreatePostRequest } from "../types";
// import { css } from '../../styled-system/css';
// import { flex } from '../../styled-system/patterns';

const CreatePost: React.FC = () => {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [thumbnailUrl, setThumbnailUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const postData: CreatePostRequest = {
        title,
        content,
        ...(thumbnailUrl && { thumbnail_url: thumbnailUrl }),
      };

      await postsAPI.createPost(postData);
      alert("投稿が作成されました。管理者の承認をお待ちください。");
      navigate("/my-page");
    } catch (err) {
      setError("投稿の作成に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  // Temporary basic styles
  const titleStyles = "text-2xl font-bold text-gray-900 mb-8 text-center";
  const formContainerStyles =
    "max-w-4xl mx-auto bg-white rounded-xl shadow-lg p-8";
  const fieldStyles = "mb-6";
  const labelStyles = "block text-sm font-semibold text-gray-700 mb-2";
  const inputStyles =
    "w-full px-4 py-3 border border-gray-300 rounded-lg text-sm";
  const textareaStyles =
    "w-full px-4 py-3 border border-gray-300 rounded-lg text-sm resize-vertical min-h-48";
  const errorStyles =
    "bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-6 text-sm";
  const buttonGroupStyles = "flex gap-4 justify-center mt-8";
  const primaryButtonStyles =
    "bg-blue-600 text-white px-6 py-3 rounded-lg border-none cursor-pointer text-sm font-semibold";
  const secondaryButtonStyles =
    "bg-gray-500 text-white px-6 py-3 rounded-lg border-none cursor-pointer text-sm font-semibold";
  const previewStyles = "mt-4 p-4 bg-gray-50 rounded-md border border-gray-200";

  return (
    <Layout>
      <div>
        <h1 className={titleStyles}>新規投稿作成</h1>
        <div className={formContainerStyles}>
          <form onSubmit={handleSubmit}>
            <div className={fieldStyles}>
              <label
                htmlFor="title"
                className={labelStyles}
              >
                タイトル
              </label>
              <input
                id="title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                className={inputStyles}
                placeholder="魅力的なタイトルを入力してください"
              />
            </div>

            <div className={fieldStyles}>
              <label
                htmlFor="thumbnailUrl"
                className={labelStyles}
              >
                サムネイル画像URL（任意）
              </label>
              <input
                id="thumbnailUrl"
                type="url"
                value={thumbnailUrl}
                onChange={(e) => setThumbnailUrl(e.target.value)}
                placeholder="https://example.com/image.jpg"
                className={inputStyles}
              />
              {thumbnailUrl && (
                <div className={previewStyles}>
                  <p className="text-sm text-gray-600 mb-2">プレビュー:</p>
                  <img
                    src={thumbnailUrl}
                    alt="プレビュー"
                    className="max-w-64 max-h-48 object-cover rounded-md"
                    onError={(e) => {
                      e.currentTarget.style.display = "none";
                    }}
                  />
                </div>
              )}
            </div>

            <div className={fieldStyles}>
              <label
                htmlFor="content"
                className={labelStyles}
              >
                内容
              </label>
              <textarea
                id="content"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                required
                className={textareaStyles}
                placeholder="投稿の内容を詳しく書いてください..."
              />
              <p className="text-xs text-gray-500 mt-1">
                {content.length} 文字
              </p>
            </div>

            {error && <div className={errorStyles}>{error}</div>}

            <div className={buttonGroupStyles}>
              <button
                type="submit"
                disabled={loading}
                className={primaryButtonStyles}
              >
                {loading ? "投稿中..." : "投稿する"}
              </button>
              <button
                type="button"
                onClick={() => navigate("/home")}
                className={secondaryButtonStyles}
              >
                キャンセル
              </button>
            </div>
          </form>
        </div>
      </div>
    </Layout>
  );
};

export default CreatePost;
