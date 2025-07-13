import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
// import { css } from '../styled-system/css';
// import { center } from '../styled-system/patterns';

const Login: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await login(email, password);
      navigate("/home");
    } catch (err) {
      setError(err instanceof Error ? err.message : "ログインに失敗しました");
    } finally {
      setLoading(false);
    }
  };

  const containerStyles =
    "min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4";

  const cardStyles = "max-w-md w-full bg-white shadow-lg rounded-xl p-8";

  const titleStyles = "text-2xl font-bold text-center text-gray-900 mb-8";

  const formStyles = "";

  const fieldStyles = "mb-6";

  const labelStyles = "block text-sm font-medium text-gray-700 mb-2";

  const inputStyles =
    "w-full px-3 py-2 border border-gray-300 rounded-md text-sm transition-all focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100";

  const errorStyles =
    "bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4 text-sm";

  const buttonStyles =
    "w-full bg-blue-600 text-white py-3 px-4 rounded-md text-sm font-medium border-none cursor-pointer transition-all hover:bg-blue-700 hover:-translate-y-0.5 hover:shadow-lg disabled:bg-gray-300 disabled:cursor-not-allowed disabled:transform-none";

  const linkContainerStyles = "text-center mt-6";

  const linkItemStyles = "mb-4";

  const linkStyles =
    "text-blue-600 no-underline text-sm hover:text-blue-700 hover:underline";

  return (
    <div className={containerStyles}>
      <div className={cardStyles}>
        <h1 className={titleStyles}>ログイン</h1>
        <form
          onSubmit={handleSubmit}
          className={formStyles}
        >
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
              placeholder="パスワードを入力"
            />
          </div>
          {error && <div className={errorStyles}>{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className={buttonStyles}
          >
            {loading ? "ログイン中..." : "ログイン"}
          </button>
        </form>
        <div className={linkContainerStyles}>
          <div className={linkItemStyles}>
            <Link
              to="/register"
              className={linkStyles}
            >
              新規登録はこちら
            </Link>
          </div>
          <div className={linkItemStyles}>
            <Link
              to="/forgot-password"
              className={linkStyles}
            >
              パスワードを忘れた方
            </Link>
          </div>
          <div className={linkItemStyles}>
            <Link
              to="/admin-login-page"
              className={linkStyles}
            >
              管理者ログイン
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
