import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
// import { css } from '../../styled-system/css';
// import { center } from '../../styled-system/patterns';

const AdminLogin: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { adminLogin } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await adminLogin(email, password);
      navigate("/admin");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "管理者ログインに失敗しました"
      );
    } finally {
      setLoading(false);
    }
  };

  // Temporary basic styles
  const containerStyles =
    "min-h-screen bg-gray-100 flex items-center justify-center py-12 px-4";
  const cardStyles =
    "max-w-md w-full bg-white shadow-xl rounded-xl p-8 border-2 border-orange-200";
  const titleStyles = "text-2xl font-bold text-center text-orange-700 mb-8";
  const fieldStyles = "mb-6";
  const labelStyles = "block text-sm font-semibold text-gray-700 mb-2";
  const inputStyles =
    "w-full px-3 py-2 border border-gray-300 rounded-md text-sm";
  const errorStyles =
    "bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4 text-sm";
  const buttonStyles =
    "w-full bg-orange-600 text-white py-3 px-4 rounded-md text-sm font-semibold border-none cursor-pointer";
  const linkContainerStyles = "text-center mt-6";
  const linkStyles = "text-orange-600 no-underline text-sm";

  return (
    <div className={containerStyles}>
      <div className={cardStyles}>
        <h1 className={titleStyles}>管理者ログイン</h1>
        <form onSubmit={handleSubmit}>
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
              placeholder="admin@example.com"
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
              placeholder="管理者パスワードを入力"
            />
          </div>
          {error && <div className={errorStyles}>{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className={buttonStyles}
          >
            {loading ? "ログイン中..." : "管理者ログイン"}
          </button>
        </form>
        <div className={linkContainerStyles}>
          <Link
            to="/login"
            className={linkStyles}
          >
            一般ユーザーログインはこちら
          </Link>
        </div>
      </div>
    </div>
  );
};

export default AdminLogin;
