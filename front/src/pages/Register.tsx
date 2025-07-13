import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
// import { css } from '../../styled-system/css';
// import { center } from '../../styled-system/patterns';

const Register: React.FC = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password !== confirmPassword) {
      setError("パスワードが一致しません");
      return;
    }

    setLoading(true);

    try {
      await register(username, email, password);
      alert("ユーザー登録が完了しました。ログインしてください。");
      navigate("/login");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "ユーザー登録に失敗しました"
      );
    } finally {
      setLoading(false);
    }
  };

  // Temporary basic styles
  const containerStyles =
    "min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4";
  const cardStyles = "max-w-md w-full bg-white shadow-lg rounded-xl p-8";
  const titleStyles = "text-2xl font-bold text-center text-gray-900 mb-8";
  const fieldStyles = "mb-6";
  const labelStyles = "block text-sm font-medium text-gray-700 mb-2";
  const inputStyles =
    "w-full px-3 py-2 border border-gray-300 rounded-md text-sm";
  const errorStyles =
    "bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4 text-sm";
  const buttonStyles =
    "w-full bg-blue-600 text-white py-3 px-4 rounded-md text-sm font-medium border-none cursor-pointer";
  const linkContainerStyles = "text-center mt-6";
  const linkStyles = "text-blue-600 no-underline text-sm";

  return (
    <div className={containerStyles}>
      <div className={cardStyles}>
        <h1 className={titleStyles}>新規登録</h1>
        <form onSubmit={handleSubmit}>
          <div className={fieldStyles}>
            <label
              htmlFor="username"
              className={labelStyles}
            >
              ユーザー名
            </label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className={inputStyles}
              placeholder="お名前を入力してください"
            />
          </div>
          <div className={fieldStyles}>
            <label
              htmlFor="email"
              className={labelStyles}
            >
              メールアドレス
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className={inputStyles}
              placeholder="your-email@example.com"
            />
          </div>
          <div className={fieldStyles}>
            <label
              htmlFor="password"
              className={labelStyles}
            >
              パスワード
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className={inputStyles}
              placeholder="8文字以上のパスワードを入力"
            />
          </div>
          <div className={fieldStyles}>
            <label
              htmlFor="confirmPassword"
              className={labelStyles}
            >
              パスワード確認
            </label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              className={inputStyles}
              placeholder="パスワードを再入力してください"
            />
          </div>
          {error && <div className={errorStyles}>{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className={buttonStyles}
          >
            {loading ? "登録中..." : "新規登録"}
          </button>
        </form>
        <div className={linkContainerStyles}>
          <Link
            to="/login"
            className={linkStyles}
          >
            既にアカウントをお持ちの方はこちら
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Register;
