import React, { ReactNode } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
// import { css } from '../../styled-system/css';
// import { flex, container } from '../../styled-system/patterns';

interface LayoutProps {
  children: ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  // Temporary basic styles
  const navStyles = 'bg-white border-b border-gray-200 shadow-sm sticky top-0 z-50 mb-8';
  const navContainerStyles = 'max-w-7xl mx-auto py-4 px-6';
  const navContentStyles = 'flex justify-between items-center';
  const logoStyles = 'text-xl font-bold text-blue-600 no-underline';
  const navLinksStyles = 'flex gap-6 items-center';
  const navLinkStyles = 'text-gray-600 no-underline px-3 py-2 rounded-md';
  const userInfoStyles = 'text-gray-700 text-sm px-3';
  const logoutButtonStyles = 'bg-red-500 text-white px-4 py-2 rounded-md border-none cursor-pointer text-sm font-medium';
  const mainContentStyles = 'max-w-7xl mx-auto px-6 py-4';

  return (
    <div>
      <nav className={navStyles}>
        <div className={navContainerStyles}>
          <div className={navContentStyles}>
            <div>
              <Link to="/home" className={logoStyles}>
    掲示板アプリ
              </Link>
            </div>
            <div className={navLinksStyles}>
              <Link to="/home" className={navLinkStyles}>
ホーム
              </Link>
              <Link to="/create-post" className={navLinkStyles}>
投稿作成
              </Link>
              <Link to="/my-page" className={navLinkStyles}>
マイページ
              </Link>
              {user?.is_admin && (
                <Link to="/admin" className={navLinkStyles}>
管理者画面
                </Link>
              )}
              <span className={userInfoStyles}>
                こんにちは、{user?.username}さん
              </span>
              <button onClick={handleLogout} className={logoutButtonStyles}>
                ログアウト
              </button>
            </div>
          </div>
        </div>
      </nav>
      <div className={mainContentStyles}>
        {children}
      </div>
    </div>
  );
};

export default Layout;