import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { userApi, authApi } from '../utils/api';
import { Post } from '../types';

export const MyPage: React.FC = () => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showEditProfile, setShowEditProfile] = useState(false);
  const [displayName, setDisplayName] = useState('');
  const [bio, setBio] = useState('');
  const [updating, setUpdating] = useState(false);
  const [updateError, setUpdateError] = useState('');
  const [updateSuccess, setUpdateSuccess] = useState('');
  const [showChangePassword, setShowChangePassword] = useState(false);
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordLoading, setPasswordLoading] = useState(false);
  const [passwordError, setPasswordError] = useState('');
  const [passwordSuccess, setPasswordSuccess] = useState('');

  const limit = 10;
  const { user, updateUser } = useAuth();

  const fetchPosts = useCallback(async () => {
    try {
      setLoading(true);
      setError('');
      const response = await userApi.getUserPosts(page, limit);
      
      // Handle the response structure safely
      if (response && typeof response === 'object') {
        const postsData = response.data || [];
        const totalCount = response.total || 0;
        
        setPosts(Array.isArray(postsData) ? postsData : []);
        setTotal(totalCount);
      } else {
        setPosts([]);
        setTotal(0);
      }
    } catch (err: any) {
      console.error('Error fetching user posts:', err);
      setError(err.message || 'Failed to fetch posts');
      setPosts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  useEffect(() => {
    fetchPosts();
    if (user) {
      setDisplayName(user.display_name);
      setBio(user.bio || '');
    }
  }, [page, user, fetchPosts]);

  const handleProfileUpdate = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!displayName.trim()) {
      setUpdateError('Display name is required');
      return;
    }

    if (displayName.length > 100) {
      setUpdateError('Display name must be less than 100 characters');
      return;
    }

    if (bio.length > 500) {
      setUpdateError('Bio must be less than 500 characters');
      return;
    }

    try {
      setUpdating(true);
      setUpdateError('');

      const updatedUser = await userApi.updateProfile(displayName.trim(), bio.trim() || undefined);
      updateUser(updatedUser);
      setUpdateSuccess('Profile updated successfully');
      setShowEditProfile(false);
      
      setTimeout(() => setUpdateSuccess(''), 3000);
    } catch (err: any) {
      setUpdateError(err.response?.data?.message || 'Failed to update profile');
    } finally {
      setUpdating(false);
    }
  };

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!currentPassword || !newPassword || !confirmPassword) {
      setPasswordError('All fields are required');
      return;
    }

    if (newPassword !== confirmPassword) {
      setPasswordError('New passwords do not match');
      return;
    }

    if (newPassword.length < 8) {
      setPasswordError('New password must be at least 8 characters long');
      return;
    }

    try {
      setPasswordLoading(true);
      setPasswordError('');

      await authApi.changePassword(currentPassword, newPassword);
      setPasswordSuccess('Password changed successfully');
      setShowChangePassword(false);
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
      
      setTimeout(() => setPasswordSuccess(''), 3000);
    } catch (err: any) {
      setPasswordError(err.response?.data?.message || 'Failed to change password');
    } finally {
      setPasswordLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'approved': return '#059669';
      case 'pending': return '#d97706';
      case 'rejected': return '#dc2626';
      default: return '#6b7280';
    }
  };

  const getStatusBadge = (status: string) => {
    return (
      <span style={{
        backgroundColor: getStatusColor(status),
        color: 'white',
        padding: '0.25rem 0.5rem',
        borderRadius: '0.25rem',
        fontSize: '0.75rem',
        fontWeight: '500',
        textTransform: 'uppercase',
      }}>
        {status}
      </span>
    );
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: '700' }}>My Page</h1>
        <Link
          to="/create-post"
          style={{
            backgroundColor: '#2563eb',
            color: 'white',
            padding: '0.75rem 1.5rem',
            textDecoration: 'none',
            borderRadius: '0.375rem',
            fontWeight: '500',
          }}
        >
          Create New Post
        </Link>
      </div>

      {/* Success Messages */}
      {updateSuccess && (
        <div style={{ backgroundColor: '#f0fdf4', border: '1px solid #bbf7d0', color: '#15803d', padding: '0.75rem', borderRadius: '0.375rem', marginBottom: '1rem' }}>
          {updateSuccess}
        </div>
      )}

      {passwordSuccess && (
        <div style={{ backgroundColor: '#f0fdf4', border: '1px solid #bbf7d0', color: '#15803d', padding: '0.75rem', borderRadius: '0.375rem', marginBottom: '1rem' }}>
          {passwordSuccess}
        </div>
      )}

      {/* Profile Section */}
      <div style={{ backgroundColor: 'white', border: '1px solid #e5e7eb', borderRadius: '0.5rem', padding: '1.5rem', marginBottom: '2rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '1rem' }}>
          <div>
            <h2 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '0.5rem' }}>
              Profile Information
            </h2>
            <div style={{ color: '#6b7280', fontSize: '0.875rem' }}>
              <p><strong>Email:</strong> {user?.email}</p>
              <p><strong>Display Name:</strong> {user?.display_name}</p>
              <p><strong>Subscription:</strong> {user?.subscription_status}</p>
              {user?.bio && <p><strong>Bio:</strong> {user.bio}</p>}
            </div>
          </div>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <button
              onClick={() => setShowEditProfile(!showEditProfile)}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#2563eb',
                color: 'white',
                border: 'none',
                borderRadius: '0.375rem',
                fontSize: '0.875rem',
                cursor: 'pointer',
              }}
            >
              Edit Profile
            </button>
            <button
              onClick={() => setShowChangePassword(!showChangePassword)}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#059669',
                color: 'white',
                border: 'none',
                borderRadius: '0.375rem',
                fontSize: '0.875rem',
                cursor: 'pointer',
              }}
            >
              Change Password
            </button>
          </div>
        </div>

        {/* Edit Profile Form */}
        {showEditProfile && (
          <form onSubmit={handleProfileUpdate} style={{ marginTop: '1rem', padding: '1rem', backgroundColor: '#f9fafb', borderRadius: '0.375rem' }}>
            <h3 style={{ fontSize: '1rem', fontWeight: '500', marginBottom: '1rem' }}>
              Edit Profile
            </h3>
            
            {updateError && (
              <div style={{ backgroundColor: '#fef2f2', border: '1px solid #fecaca', color: '#b91c1c', padding: '0.75rem', borderRadius: '0.375rem', marginBottom: '1rem', fontSize: '0.875rem' }}>
                {updateError}
              </div>
            )}

            <div style={{ marginBottom: '1rem' }}>
              <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.25rem' }}>
                Display Name
              </label>
              <input
                type="text"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                maxLength={100}
                required
                style={{
                  width: '100%',
                  padding: '0.5rem',
                  border: '1px solid #d1d5db',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  boxSizing: 'border-box',
                }}
              />
              <div style={{ fontSize: '0.75rem', color: '#6b7280', marginTop: '0.25rem' }}>
                {displayName.length}/100 characters
              </div>
            </div>

            <div style={{ marginBottom: '1rem' }}>
              <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.25rem' }}>
                Bio (optional)
              </label>
              <textarea
                value={bio}
                onChange={(e) => setBio(e.target.value)}
                maxLength={500}
                rows={3}
                style={{
                  width: '100%',
                  padding: '0.5rem',
                  border: '1px solid #d1d5db',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  fontFamily: 'inherit',
                  resize: 'vertical',
                  boxSizing: 'border-box',
                }}
              />
              <div style={{ fontSize: '0.75rem', color: '#6b7280', marginTop: '0.25rem' }}>
                {bio.length}/500 characters
              </div>
            </div>

            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <button
                type="submit"
                disabled={updating}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: updating ? '#9ca3af' : '#2563eb',
                  color: 'white',
                  border: 'none',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  cursor: updating ? 'not-allowed' : 'pointer',
                }}
              >
                {updating ? 'Updating...' : 'Update Profile'}
              </button>
              <button
                type="button"
                onClick={() => setShowEditProfile(false)}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: '#6b7280',
                  color: 'white',
                  border: 'none',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  cursor: 'pointer',
                }}
              >
                Cancel
              </button>
            </div>
          </form>
        )}

        {/* Change Password Form */}
        {showChangePassword && (
          <form onSubmit={handlePasswordChange} style={{ marginTop: '1rem', padding: '1rem', backgroundColor: '#f9fafb', borderRadius: '0.375rem' }}>
            <h3 style={{ fontSize: '1rem', fontWeight: '500', marginBottom: '1rem' }}>
              Change Password
            </h3>
            
            {passwordError && (
              <div style={{ backgroundColor: '#fef2f2', border: '1px solid #fecaca', color: '#b91c1c', padding: '0.75rem', borderRadius: '0.375rem', marginBottom: '1rem', fontSize: '0.875rem' }}>
                {passwordError}
              </div>
            )}

            <div style={{ marginBottom: '1rem' }}>
              <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.25rem' }}>
                Current Password
              </label>
              <input
                type="password"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                required
                style={{
                  width: '100%',
                  padding: '0.5rem',
                  border: '1px solid #d1d5db',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  boxSizing: 'border-box',
                }}
              />
            </div>

            <div style={{ marginBottom: '1rem' }}>
              <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.25rem' }}>
                New Password
              </label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                minLength={8}
                required
                style={{
                  width: '100%',
                  padding: '0.5rem',
                  border: '1px solid #d1d5db',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  boxSizing: 'border-box',
                }}
              />
            </div>

            <div style={{ marginBottom: '1rem' }}>
              <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.25rem' }}>
                Confirm New Password
              </label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
                style={{
                  width: '100%',
                  padding: '0.5rem',
                  border: '1px solid #d1d5db',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  boxSizing: 'border-box',
                }}
              />
            </div>

            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <button
                type="submit"
                disabled={passwordLoading}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: passwordLoading ? '#9ca3af' : '#059669',
                  color: 'white',
                  border: 'none',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  cursor: passwordLoading ? 'not-allowed' : 'pointer',
                }}
              >
                {passwordLoading ? 'Changing...' : 'Change Password'}
              </button>
              <button
                type="button"
                onClick={() => setShowChangePassword(false)}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: '#6b7280',
                  color: 'white',
                  border: 'none',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  cursor: 'pointer',
                }}
              >
                Cancel
              </button>
            </div>
          </form>
        )}
      </div>

      {/* Posts Section */}
      <div style={{ backgroundColor: 'white', border: '1px solid #e5e7eb', borderRadius: '0.5rem', padding: '1.5rem' }}>
        <h2 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '1.5rem' }}>
          My Posts ({total})
        </h2>

        {loading ? (
          <div style={{ textAlign: 'center', padding: '2rem' }}>
            <div>Loading posts...</div>
          </div>
        ) : error ? (
          <div style={{ textAlign: 'center', padding: '2rem', color: '#dc2626' }}>
            Error: {error}
          </div>
        ) : posts && posts.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '3rem' }}>
            <h3>No posts yet</h3>
            <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
              You haven't created any posts yet.
            </p>
            <Link
              to="/create-post"
              style={{
                backgroundColor: '#2563eb',
                color: 'white',
                padding: '0.75rem 1.5rem',
                textDecoration: 'none',
                borderRadius: '0.375rem',
                fontWeight: '500',
              }}
            >
              Create Your First Post
            </Link>
          </div>
        ) : (
          <>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
              {posts && posts.map((post) => (
                <div
                  key={post.id}
                  style={{
                    border: '1px solid #e5e7eb',
                    borderRadius: '0.375rem',
                    padding: '1rem',
                    backgroundColor: '#fafafa',
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '0.5rem' }}>
                    <div>
                      <h3 style={{ fontSize: '1.125rem', fontWeight: '600', marginBottom: '0.25rem' }}>
                        <Link
                          to={`/posts/${post.id}`}
                          style={{ color: '#1f2937', textDecoration: 'none' }}
                        >
                          {post.title}
                        </Link>
                      </h3>
                      <div style={{ color: '#6b7280', fontSize: '0.875rem' }}>
                        {new Date(post.created_at).toLocaleDateString()}
                      </div>
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                      {getStatusBadge(post.status)}
                      {post.status === 'pending' && (
                        <Link
                          to={`/posts/${post.id}/edit`}
                          style={{
                            color: '#2563eb',
                            textDecoration: 'none',
                            fontSize: '0.875rem',
                            fontWeight: '500',
                          }}
                        >
                          Edit
                        </Link>
                      )}
                    </div>
                  </div>

                  <p style={{ color: '#4b5563', lineHeight: '1.5', marginBottom: '0.5rem' }}>
                    {post.content.substring(0, 150)}
                    {post.content.length > 150 && '...'}
                  </p>

                  <div style={{ color: '#6b7280', fontSize: '0.875rem' }}>
                    {post.replies?.length || 0} replies
                  </div>
                </div>
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '1rem', marginTop: '2rem' }}>
                <button
                  onClick={() => setPage(Math.max(1, page - 1))}
                  disabled={page === 1}
                  style={{
                    padding: '0.5rem 1rem',
                    border: '1px solid #d1d5db',
                    backgroundColor: page === 1 ? '#f9fafb' : 'white',
                    color: page === 1 ? '#9ca3af' : '#374151',
                    borderRadius: '0.375rem',
                    cursor: page === 1 ? 'not-allowed' : 'pointer',
                  }}
                >
                  Previous
                </button>

                <span style={{ color: '#6b7280' }}>
                  Page {page} of {totalPages}
                </span>

                <button
                  onClick={() => setPage(Math.min(totalPages, page + 1))}
                  disabled={page === totalPages}
                  style={{
                    padding: '0.5rem 1rem',
                    border: '1px solid #d1d5db',
                    backgroundColor: page === totalPages ? '#f9fafb' : 'white',
                    color: page === totalPages ? '#9ca3af' : '#374151',
                    borderRadius: '0.375rem',
                    cursor: page === totalPages ? 'not-allowed' : 'pointer',
                  }}
                >
                  Next
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};