import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from './context/AuthContext';
import { Layout } from './components/Layout';
import { PrivateRoute } from './components/PrivateRoute';
import { AdminRoute } from './components/AdminRoute';

// Pages
import { Home } from './pages/Home';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { AdminLogin } from './pages/AdminLogin';
import { ForgotPassword } from './pages/ForgotPassword';
import { ResetPassword } from './pages/ResetPassword';
import { CreatePost } from './pages/CreatePost';
import { PostDetail } from './pages/PostDetail';
import { MyPage } from './pages/MyPage';
import { Subscription } from './pages/Subscription';
import { Groups } from './pages/Groups';
import { AdminDashboard } from './pages/AdminDashboard';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Router>
          <div className="App">
            <Routes>
              {/* Public routes without layout */}
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/admin-login-page" element={<AdminLogin />} />
              <Route path="/forgot-password" element={<ForgotPassword />} />
              <Route path="/reset-password" element={<ResetPassword />} />

              {/* Routes with layout */}
              <Route
                path="/*"
                element={
                  <Layout>
                    <Routes>
                      {/* Public routes with layout */}
                      <Route path="/" element={<PrivateRoute><Home /></PrivateRoute>} />
                      
                      {/* Protected routes */}
                      <Route path="/create-post" element={<PrivateRoute><CreatePost /></PrivateRoute>} />
                      <Route path="/posts/:id" element={<PrivateRoute><PostDetail /></PrivateRoute>} />
                      <Route path="/my-page" element={<PrivateRoute><MyPage /></PrivateRoute>} />
                      <Route path="/subscription" element={<PrivateRoute><Subscription /></PrivateRoute>} />
                      <Route path="/groups" element={<PrivateRoute><Groups /></PrivateRoute>} />

                      {/* Admin routes */}
                      <Route path="/admin" element={<AdminRoute><AdminDashboard /></AdminRoute>} />

                      {/* Fallback */}
                      <Route path="*" element={<Navigate to="/" replace />} />
                    </Routes>
                  </Layout>
                }
              />
            </Routes>
          </div>
        </Router>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;