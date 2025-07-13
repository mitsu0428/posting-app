import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { postApi } from '../utils/api';
import { Post, Reply } from '../types';

export const PostDetail: React.FC = () => {
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [replyContent, setReplyContent] = useState('');
  const [isAnonymous, setIsAnonymous] = useState(false);
  const [submittingReply, setSubmittingReply] = useState(false);
  const [replyError, setReplyError] = useState('');

  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (id) {
      fetchPost();
    }
  }, [id]);

  const fetchPost = async () => {
    try {
      setLoading(true);
      const response = await postApi.getPost(Number(id));
      setPost(response);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to fetch post');
    } finally {
      setLoading(false);
    }
  };

  const handleReplySubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!replyContent.trim()) {
      setReplyError('Reply content is required');
      return;
    }

    if (replyContent.length > 2000) {
      setReplyError('Reply must be less than 2000 characters');
      return;
    }

    if (user?.subscription_status !== 'active') {
      setReplyError('Active subscription required to create replies');
      return;
    }

    try {
      setSubmittingReply(true);
      setReplyError('');

      await postApi.createReply(Number(id), replyContent.trim(), isAnonymous);
      setReplyContent('');
      setIsAnonymous(false);
      
      // Refresh post to show new reply
      await fetchPost();
    } catch (err: any) {
      setReplyError(err.response?.data?.message || 'Failed to create reply');
    } finally {
      setSubmittingReply(false);
    }
  };

  const canEdit = user && post && user.id === post.author.id && post.status === 'pending';
  const canDelete = user && post && user.id === post.author.id;

  const handleEdit = () => {
    navigate(`/posts/${id}/edit`);
  };

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this post?')) {
      return;
    }

    try {
      await postApi.deletePost(Number(id));
      navigate('/my-page');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete post');
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <div>Loading post...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <div style={{ color: '#dc2626', marginBottom: '1rem' }}>
          Error: {error}
        </div>
        <Link to="/" style={{ color: '#2563eb', textDecoration: 'none' }}>
          ← Back to home
        </Link>
      </div>
    );
  }

  if (!post) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <div style={{ marginBottom: '1rem' }}>Post not found</div>
        <Link to="/" style={{ color: '#2563eb', textDecoration: 'none' }}>
          ← Back to home
        </Link>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ marginBottom: '1rem' }}>
        <Link to="/" style={{ color: '#2563eb', textDecoration: 'none' }}>
          ← Back to posts
        </Link>
      </div>

      <article style={{ backgroundColor: 'white', border: '1px solid #e5e7eb', borderRadius: '0.5rem', padding: '2rem', marginBottom: '2rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '1.5rem' }}>
          <div>
            <h1 style={{ fontSize: '1.875rem', fontWeight: '700', marginBottom: '1rem' }}>
              {post.title}
            </h1>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', color: '#6b7280', fontSize: '0.875rem' }}>
              <span>By {post.author.display_name}</span>
              <span>•</span>
              <span>{new Date(post.created_at).toLocaleDateString()}</span>
              <span>•</span>
              <span style={{ 
                color: post.status === 'approved' ? '#059669' : post.status === 'pending' ? '#d97706' : '#dc2626',
                fontWeight: '500'
              }}>
                {post.status}
              </span>
            </div>
          </div>
          
          {(canEdit || canDelete) && (
            <div style={{ display: 'flex', gap: '0.5rem' }}>
              {canEdit && (
                <button
                  onClick={handleEdit}
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
                  Edit
                </button>
              )}
              {canDelete && (
                <button
                  onClick={handleDelete}
                  style={{
                    padding: '0.5rem 1rem',
                    backgroundColor: '#dc2626',
                    color: 'white',
                    border: 'none',
                    borderRadius: '0.375rem',
                    fontSize: '0.875rem',
                    cursor: 'pointer',
                  }}
                >
                  Delete
                </button>
              )}
            </div>
          )}
        </div>

        {post.thumbnail_url && (
          <div style={{ marginBottom: '1.5rem' }}>
            <img
              src={post.thumbnail_url}
              alt={post.title}
              style={{
                width: '100%',
                maxHeight: '400px',
                objectFit: 'cover',
                borderRadius: '0.375rem',
              }}
            />
          </div>
        )}

        <div style={{ 
          lineHeight: '1.7', 
          fontSize: '1rem',
          whiteSpace: 'pre-wrap',
          wordBreak: 'break-word'
        }}>
          {post.content}
        </div>
      </article>

      {/* Replies Section */}
      <div style={{ backgroundColor: 'white', border: '1px solid #e5e7eb', borderRadius: '0.5rem', padding: '2rem' }}>
        <h2 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '1.5rem' }}>
          Replies ({post.replies?.length || 0})
        </h2>

        {/* Reply Form */}
        {user && user.subscription_status === 'active' && (
          <form onSubmit={handleReplySubmit} style={{ marginBottom: '2rem', padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '0.5rem' }}>
            <h3 style={{ fontSize: '1rem', fontWeight: '500', marginBottom: '1rem' }}>
              Add a reply
            </h3>
            
            {replyError && (
              <div style={{ backgroundColor: '#fef2f2', border: '1px solid #fecaca', color: '#b91c1c', padding: '0.75rem', borderRadius: '0.375rem', marginBottom: '1rem', fontSize: '0.875rem' }}>
                {replyError}
              </div>
            )}

            <textarea
              value={replyContent}
              onChange={(e) => setReplyContent(e.target.value)}
              placeholder="Write your reply..."
              maxLength={2000}
              rows={4}
              style={{
                width: '100%',
                padding: '0.75rem',
                border: '1px solid #d1d5db',
                borderRadius: '0.375rem',
                fontSize: '0.875rem',
                fontFamily: 'inherit',
                resize: 'vertical',
                boxSizing: 'border-box',
                marginBottom: '1rem',
              }}
            />
            
            <div style={{ fontSize: '0.75rem', color: '#6b7280', marginBottom: '1rem' }}>
              {replyContent.length}/2000 characters
            </div>

            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem' }}>
                <input
                  type="checkbox"
                  checked={isAnonymous}
                  onChange={(e) => setIsAnonymous(e.target.checked)}
                />
                Post anonymously
              </label>

              <button
                type="submit"
                disabled={submittingReply || !replyContent.trim()}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: (submittingReply || !replyContent.trim()) ? '#9ca3af' : '#2563eb',
                  color: 'white',
                  border: 'none',
                  borderRadius: '0.375rem',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  cursor: (submittingReply || !replyContent.trim()) ? 'not-allowed' : 'pointer',
                }}
              >
                {submittingReply ? 'Posting...' : 'Post Reply'}
              </button>
            </div>
          </form>
        )}

        {user && user.subscription_status !== 'active' && (
          <div style={{ backgroundColor: '#fef3c7', border: '1px solid #fbbf24', borderRadius: '0.5rem', padding: '1rem', marginBottom: '2rem' }}>
            <p style={{ color: '#92400e', fontSize: '0.875rem' }}>
              Active subscription required to post replies.{' '}
              <Link to="/subscription" style={{ color: '#92400e', fontWeight: '500' }}>
                Manage subscription
              </Link>
            </p>
          </div>
        )}

        {/* Replies List */}
        {post.replies && post.replies.length > 0 ? (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            {post.replies.map((reply: Reply) => (
              <div
                key={reply.id}
                style={{
                  padding: '1rem',
                  border: '1px solid #e5e7eb',
                  borderRadius: '0.375rem',
                  backgroundColor: '#fafafa',
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '0.5rem' }}>
                  <div style={{ fontSize: '0.875rem', color: '#6b7280' }}>
                    {reply.is_anonymous ? (
                      <span style={{ fontStyle: 'italic' }}>Anonymous</span>
                    ) : (
                      <span>{reply.author?.display_name || 'Unknown User'}</span>
                    )}
                    <span style={{ margin: '0 0.5rem' }}>•</span>
                    <span>{new Date(reply.created_at).toLocaleDateString()}</span>
                  </div>
                </div>
                <div style={{ 
                  lineHeight: '1.6',
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-word'
                }}>
                  {reply.content}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div style={{ textAlign: 'center', color: '#6b7280', padding: '2rem' }}>
            No replies yet. Be the first to reply!
          </div>
        )}
      </div>
    </div>
  );
};