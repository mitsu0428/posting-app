import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { postApi } from '../utils/api';
import { Post } from '../generated/models/post';
import { useAuth } from '../context/AuthContext';
import { deletePostsId, postPostsIdLike } from '../generated/api';

export const Home: React.FC = () => {
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 10;

  const fetchPosts = useCallback(async () => {
    try {
      setLoading(true);
      setError('');
      const response = await postApi.getPosts(page, limit);

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
      console.error('Error fetching posts:', err);
      setError(err.message || 'Failed to fetch posts');
      setPosts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  const handleDeletePost = async (postId: number) => {
    if (!user) return;

    if (!window.confirm('この投稿を削除しますか？')) {
      return;
    }

    try {
      await deletePostsId(postId);
      // Remove deleted post from local state
      setPosts(posts.filter((post) => post.id !== postId));
      setTotal(total - 1);
    } catch (err: any) {
      console.error('Error deleting post:', err);
      setError(err.message || 'Failed to delete post');
    }
  };

  const canDeletePost = (post: Post): boolean => {
    if (!user) return false;
    // Admin can delete any post, users can delete their own posts
    return user.role === 'admin' || post.author.id === user.id;
  };

  const handleToggleLike = async (postId: number) => {
    if (!user) return;

    try {
      await postPostsIdLike(postId);
      // Update the post's like status in local state
      setPosts(
        posts.map((post) => {
          if (post.id === postId) {
            return {
              ...post,
              is_liked: !post.is_liked,
              likes_count: post.is_liked
                ? (post.likes_count || 0) - 1
                : (post.likes_count || 0) + 1,
            };
          }
          return post;
        })
      );
    } catch (err: any) {
      console.error('Error toggling like:', err);
      setError(err.message || 'Failed to toggle like');
    }
  };

  useEffect(() => {
    if (!authLoading) {
      if (isAuthenticated) {
        fetchPosts();
      } else {
        setLoading(false);
      }
    }
  }, [page, isAuthenticated, authLoading, fetchPosts]);

  const totalPages = Math.ceil(total / limit);

  if (authLoading || loading) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <div>Loading posts...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <h2>Welcome to Posting App</h2>
        <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
          Please log in to view and create posts.
        </p>
        <Link
          to="/login"
          style={{
            backgroundColor: '#2563eb',
            color: 'white',
            padding: '0.75rem 1.5rem',
            textDecoration: 'none',
            borderRadius: '0.375rem',
            fontWeight: '500',
            marginRight: '1rem',
          }}
        >
          Log In
        </Link>
        <Link
          to="/register"
          style={{
            backgroundColor: '#6b7280',
            color: 'white',
            padding: '0.75rem 1.5rem',
            textDecoration: 'none',
            borderRadius: '0.375rem',
            fontWeight: '500',
          }}
        >
          Sign Up
        </Link>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem', color: 'red' }}>
        Error: {error}
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '2rem',
        }}
      >
        <h1>最新の投稿</h1>
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
          新規作成
        </Link>
      </div>

      {posts && posts.length === 0 ? (
        <div
          style={{
            textAlign: 'center',
            padding: '3rem',
            backgroundColor: '#f9fafb',
            borderRadius: '0.5rem',
          }}
        >
          <h3>No posts yet</h3>
          <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
            Be the first to create a post!
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
            新規作成
          </Link>
        </div>
      ) : (
        <>
          <div
            style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}
          >
            {posts &&
              posts.map((post) => (
                <div
                  key={post.id}
                  style={{
                    border: '1px solid #e5e7eb',
                    borderRadius: '0.5rem',
                    padding: '1.5rem',
                    backgroundColor: 'white',
                  }}
                >
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'start',
                      marginBottom: '1rem',
                    }}
                  >
                    <div>
                      <h2
                        style={{
                          fontSize: '1.25rem',
                          fontWeight: '600',
                          marginBottom: '0.5rem',
                        }}
                      >
                        <Link
                          to={`/posts/${post.id}`}
                          style={{ color: '#1f2937', textDecoration: 'none' }}
                        >
                          {post.title}
                        </Link>
                      </h2>
                      <div style={{ color: '#6b7280', fontSize: '0.875rem' }}>
                        By {post.author.display_name} •{' '}
                        {new Date(post.created_at).toLocaleDateString()}
                      </div>
                    </div>
                    {post.thumbnail_url && (
                      <img
                        src={post.thumbnail_url}
                        alt={post.title}
                        style={{
                          width: '80px',
                          height: '80px',
                          objectFit: 'cover',
                          borderRadius: '0.375rem',
                        }}
                      />
                    )}
                  </div>

                  <p
                    style={{
                      color: '#4b5563',
                      lineHeight: '1.6',
                      marginBottom: '1rem',
                    }}
                  >
                    {post.content.substring(0, 200)}
                    {post.content.length > 200 && '...'}
                  </p>

                  {/* Categories */}
                  {post.categories && post.categories.length > 0 && (
                    <div
                      style={{
                        display: 'flex',
                        flexWrap: 'wrap',
                        gap: '0.5rem',
                        marginBottom: '1rem',
                      }}
                    >
                      {post.categories.map((category) => (
                        <span
                          key={category.id}
                          style={{
                            padding: '0.25rem 0.75rem',
                            backgroundColor: category.color,
                            color: 'white',
                            borderRadius: '1rem',
                            fontSize: '0.75rem',
                            fontWeight: '500',
                          }}
                        >
                          {category.name}
                        </span>
                      ))}
                    </div>
                  )}

                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '1rem',
                      }}
                    >
                      <Link
                        to={`/posts/${post.id}`}
                        style={{
                          color: '#2563eb',
                          textDecoration: 'none',
                          fontWeight: '500',
                        }}
                      >
                        詳細はこちら →
                      </Link>
                      <button
                        onClick={() => handleToggleLike(post.id)}
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: '0.25rem',
                          backgroundColor: 'transparent',
                          border: 'none',
                          color: post.is_liked ? '#dc2626' : '#6b7280',
                          cursor: 'pointer',
                          padding: '0.25rem',
                          borderRadius: '0.25rem',
                          fontSize: '0.875rem',
                        }}
                        onMouseOver={(e) =>
                          (e.currentTarget.style.backgroundColor = '#f3f4f6')
                        }
                        onMouseOut={(e) =>
                          (e.currentTarget.style.backgroundColor =
                            'transparent')
                        }
                      >
                        <span style={{ fontSize: '1rem' }}>
                          {post.is_liked ? '❤️' : '🤍'}
                        </span>
                        <span>{post.likes_count || 0}</span>
                      </button>
                      {canDeletePost(post) && (
                        <button
                          onClick={() => handleDeletePost(post.id)}
                          style={{
                            backgroundColor: '#dc2626',
                            color: 'white',
                            border: 'none',
                            padding: '0.375rem 0.75rem',
                            borderRadius: '0.25rem',
                            fontSize: '0.875rem',
                            cursor: 'pointer',
                            fontWeight: '500',
                          }}
                          onMouseOver={(e) =>
                            (e.currentTarget.style.backgroundColor = '#b91c1c')
                          }
                          onMouseOut={(e) =>
                            (e.currentTarget.style.backgroundColor = '#dc2626')
                          }
                        >
                          削除
                        </button>
                      )}
                    </div>
                    <div style={{ color: '#6b7280', fontSize: '0.875rem' }}>
                      {post.replies?.length || 0} replies
                    </div>
                  </div>
                </div>
              ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                gap: '1rem',
                marginTop: '2rem',
              }}
            >
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
  );
};
