import React, { useState, useEffect, useCallback } from 'react';
import { DataGrid, GridColDef, GridActionsCellItem } from '@mui/x-data-grid';
import {
  Button,
  Tabs,
  Tab,
  Box,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import { CheckCircle, Cancel, Visibility, Block } from '@mui/icons-material';
import { adminApi } from '../utils/api';
import { Post, User, PaginatedResponse } from '../types';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`admin-tabpanel-${index}`}
      aria-labelledby={`admin-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

export const AdminDashboard: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [posts, setPosts] = useState<Post[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [postsLoading, setPostsLoading] = useState(false);
  const [usersLoading, setUsersLoading] = useState(false);
  const [postFilter, setPostFilter] = useState<string>('all');
  const [selectedPost, setSelectedPost] = useState<Post | null>(null);
  const [viewDialogOpen, setViewDialogOpen] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const fetchPosts = useCallback(async () => {
    try {
      setPostsLoading(true);
      setError('');
      const statusParam = postFilter === 'all' ? undefined : postFilter;
      const response: PaginatedResponse<Post> = await adminApi.getPosts(
        1,
        100,
        statusParam
      );
      setPosts(Array.isArray(response?.data) ? response.data : []);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to fetch posts');
      setPosts([]);
    } finally {
      setPostsLoading(false);
    }
  }, [postFilter]);

  const fetchUsers = useCallback(async () => {
    try {
      setUsersLoading(true);
      setError('');
      const response: PaginatedResponse<User> = await adminApi.getUsers(1, 100);
      setUsers(Array.isArray(response?.data) ? response.data : []);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to fetch users');
      setUsers([]);
    } finally {
      setUsersLoading(false);
    }
  }, []);

  useEffect(() => {
    if (tabValue === 0) {
      fetchPosts();
    } else if (tabValue === 1) {
      fetchUsers();
    }
  }, [tabValue, postFilter, fetchPosts, fetchUsers]);

  const handleApprovePost = async (postId: number) => {
    try {
      setActionLoading(true);
      await adminApi.approvePost(postId);
      setSuccess('Post approved successfully');
      fetchPosts();
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to approve post');
    } finally {
      setActionLoading(false);
    }
  };

  const handleRejectPost = async (postId: number) => {
    try {
      setActionLoading(true);
      await adminApi.rejectPost(postId);
      setSuccess('Post rejected successfully');
      fetchPosts();
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to reject post');
    } finally {
      setActionLoading(false);
    }
  };

  const handleBanUser = async (userId: number) => {
    if (!window.confirm('Are you sure you want to ban this user?')) {
      return;
    }

    try {
      setActionLoading(true);
      await adminApi.banUser(userId);
      setSuccess('User banned successfully');
      fetchUsers();
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to ban user');
    } finally {
      setActionLoading(false);
    }
  };

  const handleViewPost = (post: Post) => {
    setSelectedPost(post);
    setViewDialogOpen(true);
  };

  const getStatusChip = (status: string) => {
    const statusConfig = {
      pending: { color: 'warning' as const, label: 'Pending' },
      approved: { color: 'success' as const, label: 'Approved' },
      rejected: { color: 'error' as const, label: 'Rejected' },
    };

    const config = statusConfig[status as keyof typeof statusConfig] || {
      color: 'default' as const,
      label: status,
    };
    return <Chip label={config.label} color={config.color} size="small" />;
  };

  const getSubscriptionChip = (status: string) => {
    const statusConfig = {
      active: { color: 'success' as const, label: 'Active' },
      inactive: { color: 'default' as const, label: 'Inactive' },
      past_due: { color: 'warning' as const, label: 'Past Due' },
      canceled: { color: 'error' as const, label: 'Canceled' },
    };

    const config = statusConfig[status as keyof typeof statusConfig] || {
      color: 'default' as const,
      label: status,
    };
    return <Chip label={config.label} color={config.color} size="small" />;
  };

  const postColumns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 70 },
    { field: 'title', headerName: 'Title', width: 250 },
    {
      field: 'author',
      headerName: 'Author',
      width: 150,
      valueGetter: (params) => params.row.author?.display_name || 'Unknown',
    },
    {
      field: 'status',
      headerName: 'Status',
      width: 120,
      renderCell: (params) => getStatusChip(params.value),
    },
    {
      field: 'created_at',
      headerName: 'Created',
      width: 150,
      valueGetter: (params) => new Date(params.value).toLocaleDateString(),
    },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Actions',
      width: 150,
      getActions: (params) => [
        <GridActionsCellItem
          icon={<Visibility />}
          label="View"
          onClick={() => handleViewPost(params.row)}
        />,
        ...(params.row.status === 'pending'
          ? [
              <GridActionsCellItem
                icon={<CheckCircle />}
                label="Approve"
                onClick={() => handleApprovePost(params.row.id)}
                disabled={actionLoading}
              />,
              <GridActionsCellItem
                icon={<Cancel />}
                label="Reject"
                onClick={() => handleRejectPost(params.row.id)}
                disabled={actionLoading}
              />,
            ]
          : []),
      ],
    },
  ];

  const userColumns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 70 },
    { field: 'email', headerName: 'Email', width: 200 },
    { field: 'display_name', headerName: 'Display Name', width: 150 },
    { field: 'role', headerName: 'Role', width: 100 },
    {
      field: 'subscription_status',
      headerName: 'Subscription',
      width: 120,
      renderCell: (params) => getSubscriptionChip(params.value),
    },
    {
      field: 'is_active',
      headerName: 'Status',
      width: 100,
      renderCell: (params) => (
        <Chip
          label={params.value ? 'Active' : 'Banned'}
          color={params.value ? 'success' : 'error'}
          size="small"
        />
      ),
    },
    {
      field: 'created_at',
      headerName: 'Joined',
      width: 150,
      valueGetter: (params) => new Date(params.value).toLocaleDateString(),
    },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Actions',
      width: 100,
      getActions: (params) => [
        ...(params.row.is_active && params.row.role !== 'admin'
          ? [
              <GridActionsCellItem
                icon={<Block />}
                label="Ban"
                onClick={() => handleBanUser(params.row.id)}
                disabled={actionLoading}
              />,
            ]
          : []),
      ],
    },
  ];

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
      <h1 style={{ fontSize: '2rem', fontWeight: '700', marginBottom: '2rem' }}>
        管理者ダッシュボード
      </h1>

      {error && (
        <div
          style={{
            backgroundColor: '#fef2f2',
            border: '1px solid #fecaca',
            color: '#b91c1c',
            padding: '0.75rem',
            borderRadius: '0.375rem',
            marginBottom: '1rem',
          }}
        >
          {error}
        </div>
      )}

      {success && (
        <div
          style={{
            backgroundColor: '#f0fdf4',
            border: '1px solid #bbf7d0',
            color: '#15803d',
            padding: '0.75rem',
            borderRadius: '0.375rem',
            marginBottom: '1rem',
          }}
        >
          {success}
        </div>
      )}

      <Box sx={{ borderBottom: 1, borderColor: 'divider', marginBottom: 2 }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="投稿管理" />
          <Tab label="ユーザー管理" />
        </Tabs>
      </Box>

      <TabPanel value={tabValue} index={0}>
        <div style={{ marginBottom: '1rem' }}>
          <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
            <Button
              variant={postFilter === 'all' ? 'contained' : 'outlined'}
              onClick={() => setPostFilter('all')}
            >
              全ての投稿
            </Button>
            <Button
              variant={postFilter === 'pending' ? 'contained' : 'outlined'}
              onClick={() => setPostFilter('pending')}
            >
              保留中の投稿
            </Button>
            <Button
              variant={postFilter === 'approved' ? 'contained' : 'outlined'}
              onClick={() => setPostFilter('approved')}
            >
              承認済みの投稿
            </Button>
            <Button
              variant={postFilter === 'rejected' ? 'contained' : 'outlined'}
              onClick={() => setPostFilter('rejected')}
            >
              拒否された投稿
            </Button>
          </div>
        </div>

        <div style={{ height: 600, width: '100%', backgroundColor: 'white' }}>
          <DataGrid
            rows={
              Array.isArray(posts)
                ? posts.filter((post) => post && post.id)
                : []
            }
            columns={postColumns}
            loading={postsLoading}
            initialState={{
              pagination: {
                paginationModel: { page: 0, pageSize: 10 },
              },
            }}
            pageSizeOptions={[10, 25, 50]}
            disableRowSelectionOnClick
          />
        </div>
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        <div style={{ height: 600, width: '100%', backgroundColor: 'white' }}>
          <DataGrid
            rows={
              Array.isArray(users)
                ? users.filter((user) => user && user.id)
                : []
            }
            columns={userColumns}
            loading={usersLoading}
            initialState={{
              pagination: {
                paginationModel: { page: 0, pageSize: 10 },
              },
            }}
            pageSizeOptions={[10, 25, 50]}
            disableRowSelectionOnClick
          />
        </div>
      </TabPanel>

      {/* Post View Dialog */}
      <Dialog
        open={viewDialogOpen}
        onClose={() => setViewDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>Post Details</DialogTitle>
        <DialogContent>
          {selectedPost && (
            <div>
              <div style={{ marginBottom: '1rem' }}>
                <strong>Title:</strong> {selectedPost.title}
              </div>
              <div style={{ marginBottom: '1rem' }}>
                <strong>Author:</strong> {selectedPost.author?.display_name}
              </div>
              <div style={{ marginBottom: '1rem' }}>
                <strong>Status:</strong> {getStatusChip(selectedPost.status)}
              </div>
              <div style={{ marginBottom: '1rem' }}>
                <strong>Created:</strong>{' '}
                {new Date(selectedPost.created_at).toLocaleString()}
              </div>
              {selectedPost.thumbnail_url && (
                <div style={{ marginBottom: '1rem' }}>
                  <strong>Thumbnail:</strong>
                  <br />
                  <img
                    src={selectedPost.thumbnail_url}
                    alt="Post thumbnail"
                    style={{
                      maxWidth: '200px',
                      maxHeight: '150px',
                      objectFit: 'cover',
                      borderRadius: '4px',
                      marginTop: '0.5rem',
                    }}
                  />
                </div>
              )}
              <div style={{ marginBottom: '1rem' }}>
                <strong>Content:</strong>
                <div
                  style={{
                    padding: '1rem',
                    backgroundColor: '#f9fafb',
                    borderRadius: '4px',
                    marginTop: '0.5rem',
                    whiteSpace: 'pre-wrap',
                    wordBreak: 'break-word',
                  }}
                >
                  {selectedPost.content}
                </div>
              </div>
              {selectedPost.replies && selectedPost.replies.length > 0 && (
                <div>
                  <strong>Replies ({selectedPost.replies.length}):</strong>
                  <div style={{ marginTop: '0.5rem' }}>
                    {selectedPost.replies.map((reply) => (
                      <div
                        key={reply.id}
                        style={{
                          padding: '0.5rem',
                          backgroundColor: '#f3f4f6',
                          borderRadius: '4px',
                          marginBottom: '0.5rem',
                          fontSize: '0.875rem',
                        }}
                      >
                        <div
                          style={{ fontWeight: '500', marginBottom: '0.25rem' }}
                        >
                          {reply.is_anonymous
                            ? 'Anonymous'
                            : reply.author?.display_name || 'Unknown'}{' '}
                          • {new Date(reply.created_at).toLocaleDateString()}
                        </div>
                        <div
                          style={{
                            whiteSpace: 'pre-wrap',
                            wordBreak: 'break-word',
                          }}
                        >
                          {reply.content}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </DialogContent>
        <DialogActions>
          {selectedPost?.status === 'pending' && (
            <>
              <Button
                onClick={() => {
                  handleApprovePost(selectedPost.id);
                  setViewDialogOpen(false);
                }}
                color="success"
                disabled={actionLoading}
              >
                Approve
              </Button>
              <Button
                onClick={() => {
                  handleRejectPost(selectedPost.id);
                  setViewDialogOpen(false);
                }}
                color="error"
                disabled={actionLoading}
              >
                Reject
              </Button>
            </>
          )}
          <Button onClick={() => setViewDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
