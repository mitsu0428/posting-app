import React, { useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { authAPI } from "../utils/api";

const ResetPassword: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const token = searchParams.get("token");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (newPassword !== confirmPassword) {
      setError("パスワードが一致しません");
      return;
    }

    if (!token) {
      setError("無効なリセットトークンです");
      return;
    }

    setLoading(true);

    try {
      await authAPI.resetPassword(token, newPassword);
      alert("パスワードがリセットされました。新しいパスワードでログインしてください。");
      navigate("/login");
    } catch (err) {
      setError("パスワードリセットに失敗しました。トークンが無効または期限切れの可能性があります。");
    } finally {
      setLoading(false);
    }
  };

  if (!token) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full bg-white shadow-lg rounded-xl p-8 text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">無効なリンク</h1>
          <p className="text-gray-600 mb-6">
            パスワードリセットリンクが無効です。
          </p>
          <a href="/login" className="text-blue-600 hover:text-blue-700">
            ログインページに戻る
          </a>
        </div>
      </div>
    );
  }

  // Basic styles
  const containerStyles = "min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4";
  const cardStyles = "max-w-md w-full bg-white shadow-lg rounded-xl p-8";
  const titleStyles = "text-2xl font-bold text-center text-gray-900 mb-8";
  const fieldStyles = "mb-6";
  const labelStyles = "block text-sm font-medium text-gray-700 mb-2";
  const inputStyles = "w-full px-3 py-2 border border-gray-300 rounded-md text-sm transition-all focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100";
  const errorStyles = "bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4 text-sm";
  const buttonStyles = "w-full bg-blue-600 text-white py-3 px-4 rounded-md text-sm font-medium border-none cursor-pointer transition-all hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed";

  return (
    <div className={containerStyles}>
      <div className={cardStyles}>
        <h1 className={titleStyles}>新しいパスワードを設定</h1>
        <form onSubmit={handleSubmit}>
          <div className={fieldStyles}>
            <label htmlFor="newPassword" className={labelStyles}>
              新しいパスワード
            </label>
            <input
              id="newPassword"
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              required
              minLength={8}
              className={inputStyles}
              placeholder="8文字以上のパスワードを入力"
            />
          </div>
          <div className={fieldStyles}>
            <label htmlFor="confirmPassword" className={labelStyles}>
              パスワード確認
            </label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              className={inputStyles}
              placeholder="パスワードを再入力"
            />
          </div>
          {error && <div className={errorStyles}>{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className={buttonStyles}
          >
            {loading ? "更新中..." : "パスワードを更新"}
          </button>
        </form>
      </div>
    </div>
  );
};

export default ResetPassword;