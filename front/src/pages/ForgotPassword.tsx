import React, { useState } from "react";
import { Link } from "react-router-dom";
import { authAPI } from "../utils/api";

const ForgotPassword: React.FC = () => {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setMessage("");
    setLoading(true);

    try {
      const response = await authAPI.forgotPassword(email);
      setMessage(response.message || "パスワードリセットのメールを送信しました。");
    } catch (err) {
      setError("メール送信に失敗しました。メールアドレスを確認してください。");
    } finally {
      setLoading(false);
    }
  };

  // Basic styles
  const containerStyles = "min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4";
  const cardStyles = "max-w-md w-full bg-white shadow-lg rounded-xl p-8";
  const titleStyles = "text-2xl font-bold text-center text-gray-900 mb-8";
  const fieldStyles = "mb-6";
  const labelStyles = "block text-sm font-medium text-gray-700 mb-2";
  const inputStyles = "w-full px-3 py-2 border border-gray-300 rounded-md text-sm transition-all focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100";
  const errorStyles = "bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4 text-sm";
  const successStyles = "bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md mb-4 text-sm";
  const buttonStyles = "w-full bg-blue-600 text-white py-3 px-4 rounded-md text-sm font-medium border-none cursor-pointer transition-all hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed";
  const linkContainerStyles = "text-center mt-6";
  const linkStyles = "text-blue-600 no-underline text-sm hover:text-blue-700 hover:underline";

  return (
    <div className={containerStyles}>
      <div className={cardStyles}>
        <h1 className={titleStyles}>パスワードを忘れた方</h1>
        <form onSubmit={handleSubmit}>
          <div className={fieldStyles}>
            <label htmlFor="email" className={labelStyles}>
              メールアドレス
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className={inputStyles}
              placeholder="登録時のメールアドレスを入力"
            />
          </div>
          {error && <div className={errorStyles}>{error}</div>}
          {message && <div className={successStyles}>{message}</div>}
          <button
            type="submit"
            disabled={loading}
            className={buttonStyles}
          >
            {loading ? "送信中..." : "パスワードリセットメールを送信"}
          </button>
        </form>
        <div className={linkContainerStyles}>
          <Link to="/login" className={linkStyles}>
            ログインページに戻る
          </Link>
        </div>
      </div>
    </div>
  );
};

export default ForgotPassword;