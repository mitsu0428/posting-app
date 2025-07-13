import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { authApi } from '../utils/api';

interface LayoutProps {
  children: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      logout();
      navigate('/login');
    }
  };

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <header
        style={{
          backgroundColor: '#2563eb',
          color: 'white',
          padding: '1rem 2rem',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '2rem' }}>
          <Link
            to="/"
            style={{
              color: 'white',
              textDecoration: 'none',
              fontSize: '1.5rem',
              fontWeight: 'bold',
            }}
          >
            Posting App
          </Link>
          
          {user && (
            <nav style={{ display: 'flex', gap: '1rem' }}>
              <Link
                to="/"
                style={{ color: 'white', textDecoration: 'none' }}
              >
                Home
              </Link>
              <Link
                to="/create-post"
                style={{ color: 'white', textDecoration: 'none' }}
              >
                Create Post
              </Link>
              <Link
                to="/my-page"
                style={{ color: 'white', textDecoration: 'none' }}
              >
                My Page
              </Link>
              <Link
                to="/subscription"
                style={{ color: 'white', textDecoration: 'none' }}
              >
                Subscription
              </Link>
              {user.role === 'admin' && (
                <Link
                  to="/admin"
                  style={{ color: 'white', textDecoration: 'none' }}
                >
                  Admin
                </Link>
              )}
            </nav>
          )}
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          {user ? (
            <>
              <span>Welcome, {user.display_name}</span>
              <button
                onClick={handleLogout}
                style={{
                  backgroundColor: 'transparent',
                  border: '1px solid white',
                  color: 'white',
                  padding: '0.5rem 1rem',
                  borderRadius: '0.25rem',
                  cursor: 'pointer',
                }}
              >
                Logout
              </button>
            </>
          ) : (
            <div style={{ display: 'flex', gap: '1rem' }}>
              <Link
                to="/login"
                style={{
                  color: 'white',
                  textDecoration: 'none',
                  padding: '0.5rem 1rem',
                  border: '1px solid white',
                  borderRadius: '0.25rem',
                }}
              >
                Login
              </Link>
              <Link
                to="/register"
                style={{
                  backgroundColor: 'white',
                  color: '#2563eb',
                  textDecoration: 'none',
                  padding: '0.5rem 1rem',
                  borderRadius: '0.25rem',
                }}
              >
                Register
              </Link>
            </div>
          )}
        </div>
      </header>

      <main style={{ flex: 1, padding: '2rem' }}>
        {children}
      </main>

      <footer
        style={{
          backgroundColor: '#f3f4f6',
          padding: '1rem 2rem',
          textAlign: 'center',
          borderTop: '1px solid #e5e7eb',
        }}
      >
        <p>&copy; 2024 Posting App. All rights reserved.</p>
      </footer>
    </div>
  );
};