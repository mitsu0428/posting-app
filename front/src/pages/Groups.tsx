import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { getGroups, postGroups, postGroupsIdMembers, putGroupsId, postGroupsIdMembersByName, getUsersSearch } from '../generated/api';
import { groupApi } from '../utils/api';
import { Group } from '../generated/models/group';

export const Groups: React.FC = () => {
  const { user, isAuthenticated } = useAuth();
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newGroupName, setNewGroupName] = useState('');
  const [newGroupDescription, setNewGroupDescription] = useState('');
  const [creating, setCreating] = useState(false);
  const [joinLoading, setJoinLoading] = useState<number | null>(null);
  const [editingGroup, setEditingGroup] = useState<Group | null>(null);
  const [editName, setEditName] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [showAddMemberForm, setShowAddMemberForm] = useState<number | null>(null);
  const [memberDisplayName, setMemberDisplayName] = useState('');
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [deletingGroup, setDeletingGroup] = useState<Group | null>(null);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [showMembersFor, setShowMembersFor] = useState<number | null>(null);
  const [members, setMembers] = useState<Array<{ id: number; display_name: string; bio: string }>>([]);
  const [membersLoading, setMembersLoading] = useState(false);
  const [memberActionLoading, setMemberActionLoading] = useState<number | null>(null);

  useEffect(() => {
    if (isAuthenticated) {
      fetchGroups();
    }
  }, [isAuthenticated]);

  const fetchGroups = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await getGroups();
      setGroups(response || []);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch groups');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newGroupName.trim() || !newGroupDescription.trim()) {
      setError('グループ名と説明は必須です');
      return;
    }

    if (newGroupName.length > 100) {
      setError('グループ名は100文字以内で入力してください');
      return;
    }

    if (newGroupDescription.length > 500) {
      setError('説明は500文字以内で入力してください');
      return;
    }

    try {
      setCreating(true);
      setError('');
      
      await postGroups({
        name: newGroupName.trim(),
        description: newGroupDescription.trim(),
      });

      setNewGroupName('');
      setNewGroupDescription('');
      setShowCreateForm(false);
      fetchGroups();
    } catch (err: any) {
      setError(err.message || 'グループ作成に失敗しました');
    } finally {
      setCreating(false);
    }
  };

  const handleJoinGroup = async (groupId: number) => {
    if (!user) return;

    try {
      setJoinLoading(groupId);
      setError('');
      
      await postGroupsIdMembers(groupId, { user_id: user.id });
      fetchGroups();
    } catch (err: any) {
      setError(err.message || 'グループ参加に失敗しました');
    } finally {
      setJoinLoading(null);
    }
  };

  const handleEditGroup = (group: Group) => {
    setEditingGroup(group);
    setEditName(group.name || '');
    setEditDescription(group.description || '');
  };

  const handleUpdateGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingGroup) return;

    if (!editName.trim() || !editDescription.trim()) {
      setError('グループ名と説明は必須です');
      return;
    }

    try {
      setCreating(true);
      setError('');
      
      await putGroupsId(editingGroup.id, {
        name: editName.trim(),
        description: editDescription.trim(),
      });
      
      setEditingGroup(null);
      setEditName('');
      setEditDescription('');
      fetchGroups();
    } catch (err: any) {
      setError(err.message || 'グループの更新に失敗しました');
    } finally {
      setCreating(false);
    }
  };

  const handleSearchUsers = async (query: string) => {
    if (query.length < 2) {
      setSearchResults([]);
      return;
    }

    try {
      const results = await getUsersSearch({ q: query });
      setSearchResults(results || []);
    } catch (err: any) {
      console.error('Failed to search users:', err);
      setSearchResults([]);
    }
  };

  const handleAddMemberByName = async (groupId: number) => {
    if (!memberDisplayName.trim()) {
      setError('ユーザー名を入力してください');
      return;
    }

    try {
      setJoinLoading(groupId);
      setError('');
      
      await postGroupsIdMembersByName(groupId, {
        display_name: memberDisplayName.trim(),
      });
      
      setShowAddMemberForm(null);
      setMemberDisplayName('');
      setSearchResults([]);
      fetchGroups();
    } catch (err: any) {
      setError(err.message || 'メンバー追加に失敗しました');
    } finally {
      setJoinLoading(null);
    }
  };

  const canJoinGroup = (group: Group): boolean => {
    if (!user) return false;
    if (group.owner_id === user.id) return false;
    // グループメンバーの確認は、実際のAPIレスポンスに依存
    return true;
  };

  const handleDeleteGroup = async () => {
    if (!deletingGroup) return;

    try {
      setDeleteLoading(true);
      setError('');
      
      await groupApi.deleteGroup(deletingGroup.id);
      
      setDeletingGroup(null);
      fetchGroups();
    } catch (err: any) {
      setError(err.message || 'グループ削除に失敗しました');
    } finally {
      setDeleteLoading(false);
    }
  };

  const fetchGroupMembers = async (groupId: number) => {
    try {
      setMembersLoading(true);
      setError('');
      
      const membersData = await groupApi.getGroupMembers(groupId);
      setMembers(membersData);
      setShowMembersFor(groupId);
    } catch (err: any) {
      setError(err.message || 'メンバー一覧の取得に失敗しました');
    } finally {
      setMembersLoading(false);
    }
  };

  const handleRemoveMember = async (groupId: number, memberId: number) => {
    try {
      setMemberActionLoading(memberId);
      setError('');
      
      await groupApi.removeGroupMember(groupId, memberId);
      
      // Update members list
      setMembers(prev => prev.filter(member => member.id !== memberId));
      fetchGroups(); // Refresh groups to update member count
    } catch (err: any) {
      setError(err.message || 'メンバー除名に失敗しました');
    } finally {
      setMemberActionLoading(null);
    }
  };

  const handleLeaveGroup = async (groupId: number) => {
    try {
      setMemberActionLoading(groupId);
      setError('');
      
      await groupApi.leaveGroup(groupId);
      
      setShowMembersFor(null);
      setMembers([]);
      fetchGroups(); // Refresh groups list
    } catch (err: any) {
      setError(err.message || 'グループ退会に失敗しました');
    } finally {
      setMemberActionLoading(null);
    }
  };

  if (!isAuthenticated) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <h2>ログインが必要です</h2>
        <p>グループ機能を使用するには、ログインしてください。</p>
      </div>
    );
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <div>Loading groups...</div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: '700' }}>グループ</h1>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          style={{
            backgroundColor: '#2563eb',
            color: 'white',
            padding: '0.75rem 1.5rem',
            border: 'none',
            borderRadius: '0.375rem',
            fontWeight: '500',
            cursor: 'pointer',
          }}
        >
          {showCreateForm ? 'キャンセル' : '新しいグループを作成'}
        </button>
      </div>

      {error && (
        <div style={{ backgroundColor: '#fef2f2', border: '1px solid #fecaca', color: '#b91c1c', padding: '0.75rem', borderRadius: '0.375rem', marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {showCreateForm && (
        <div style={{ backgroundColor: 'white', border: '1px solid #e5e7eb', borderRadius: '0.5rem', padding: '1.5rem', marginBottom: '2rem' }}>
          <h2 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '1rem' }}>新しいグループを作成</h2>
          <form onSubmit={handleCreateGroup}>
            <div style={{ marginBottom: '1rem' }}>
              <label htmlFor="groupName" style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' }}>
                グループ名 <span style={{ color: '#dc2626' }}>*</span>
              </label>
              <input
                id="groupName"
                type="text"
                value={newGroupName}
                onChange={(e) => setNewGroupName(e.target.value)}
                maxLength={100}
                required
                style={{
                  width: '100%',
                  padding: '0.75rem',
                  border: '1px solid #d1d5db',
                  borderRadius: '0.375rem',
                  fontSize: '1rem',
                  boxSizing: 'border-box',
                }}
                placeholder="グループ名を入力"
              />
              <div style={{ fontSize: '0.75rem', color: '#6b7280', marginTop: '0.25rem' }}>
                {newGroupName.length}/100 文字
              </div>
            </div>

            <div style={{ marginBottom: '1.5rem' }}>
              <label htmlFor="groupDescription" style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' }}>
                説明 <span style={{ color: '#dc2626' }}>*</span>
              </label>
              <textarea
                id="groupDescription"
                value={newGroupDescription}
                onChange={(e) => setNewGroupDescription(e.target.value)}
                maxLength={500}
                required
                rows={4}
                style={{
                  width: '100%',
                  padding: '0.75rem',
                  border: '1px solid #d1d5db',
                  borderRadius: '0.375rem',
                  fontSize: '1rem',
                  fontFamily: 'inherit',
                  resize: 'vertical',
                  boxSizing: 'border-box',
                }}
                placeholder="グループの説明を入力"
              />
              <div style={{ fontSize: '0.75rem', color: '#6b7280', marginTop: '0.25rem' }}>
                {newGroupDescription.length}/500 文字
              </div>
            </div>

            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
              <button
                type="button"
                onClick={() => setShowCreateForm(false)}
                style={{
                  padding: '0.75rem 1.5rem',
                  border: '1px solid #d1d5db',
                  backgroundColor: 'white',
                  color: '#374151',
                  borderRadius: '0.375rem',
                  fontWeight: '500',
                  cursor: 'pointer',
                }}
              >
                キャンセル
              </button>
              <button
                type="submit"
                disabled={creating}
                style={{
                  padding: '0.75rem 1.5rem',
                  border: 'none',
                  backgroundColor: creating ? '#9ca3af' : '#2563eb',
                  color: 'white',
                  borderRadius: '0.375rem',
                  fontWeight: '500',
                  cursor: creating ? 'not-allowed' : 'pointer',
                }}
              >
                {creating ? '作成中...' : 'グループを作成'}
              </button>
            </div>
          </form>
        </div>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
        {groups.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '3rem', backgroundColor: '#f9fafb', borderRadius: '0.5rem' }}>
            <h3>グループがありません</h3>
            <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
              最初のグループを作成してみましょう！
            </p>
            <button
              onClick={() => setShowCreateForm(true)}
              style={{
                backgroundColor: '#2563eb',
                color: 'white',
                padding: '0.75rem 1.5rem',
                border: 'none',
                borderRadius: '0.375rem',
                fontWeight: '500',
                cursor: 'pointer',
              }}
            >
              グループを作成
            </button>
          </div>
        ) : (
          groups.map((group) => (
            <div
              key={group.id}
              style={{
                border: '1px solid #e5e7eb',
                borderRadius: '0.5rem',
                padding: '1.5rem',
                backgroundColor: 'white',
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '1rem' }}>
                <div>
                  <h3 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '0.5rem' }}>
                    {group.name}
                  </h3>
                  <p style={{ color: '#6b7280', fontSize: '0.875rem', marginBottom: '0.5rem' }}>
                    {group.description}
                  </p>
                  <div style={{ display: 'flex', gap: '1rem', fontSize: '0.75rem', color: '#6b7280' }}>
                    <span>メンバー: {group.member_count || 0}人</span>
                    <span>作成日: {new Date(group.created_at).toLocaleDateString()}</span>
                  </div>
                </div>
                
                <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
                  {group.owner_id === user?.id ? (
                    <>
                      <span style={{
                        backgroundColor: '#fbbf24',
                        color: 'white',
                        padding: '0.25rem 0.75rem',
                        borderRadius: '1rem',
                        fontSize: '0.75rem',
                        fontWeight: '500',
                      }}>
                        オーナー
                      </span>
                      <button
                        onClick={() => handleEditGroup(group)}
                        style={{
                          backgroundColor: '#6366f1',
                          color: 'white',
                          border: 'none',
                          padding: '0.25rem 0.75rem',
                          borderRadius: '0.25rem',
                          fontSize: '0.75rem',
                          cursor: 'pointer',
                          fontWeight: '500',
                        }}
                      >
                        編集
                      </button>
                      <button
                        onClick={() => setShowAddMemberForm(showAddMemberForm === group.id ? null : group.id)}
                        style={{
                          backgroundColor: '#10b981',
                          color: 'white',
                          border: 'none',
                          padding: '0.25rem 0.75rem',
                          borderRadius: '0.25rem',
                          fontSize: '0.75rem',
                          cursor: 'pointer',
                          fontWeight: '500',
                        }}
                      >
                        メンバー追加
                      </button>
                      <button
                        onClick={() => setDeletingGroup(group)}
                        style={{
                          backgroundColor: '#dc2626',
                          color: 'white',
                          border: 'none',
                          padding: '0.25rem 0.75rem',
                          borderRadius: '0.25rem',
                          fontSize: '0.75rem',
                          cursor: 'pointer',
                          fontWeight: '500',
                        }}
                      >
                        削除
                      </button>
                      <button
                        onClick={() => fetchGroupMembers(group.id)}
                        style={{
                          backgroundColor: '#3b82f6',
                          color: 'white',
                          border: 'none',
                          padding: '0.25rem 0.75rem',
                          borderRadius: '0.25rem',
                          fontSize: '0.75rem',
                          cursor: 'pointer',
                          fontWeight: '500',
                        }}
                      >
                        メンバー一覧
                      </button>
                    </>
                  ) : canJoinGroup(group) ? (
                    <button
                      onClick={() => handleJoinGroup(group.id)}
                      disabled={joinLoading === group.id}
                      style={{
                        backgroundColor: joinLoading === group.id ? '#9ca3af' : '#10b981',
                        color: 'white',
                        border: 'none',
                        padding: '0.5rem 1rem',
                        borderRadius: '0.25rem',
                        fontSize: '0.875rem',
                        cursor: joinLoading === group.id ? 'not-allowed' : 'pointer',
                        fontWeight: '500',
                      }}
                    >
                      {joinLoading === group.id ? '参加中...' : '参加'}
                    </button>
                  ) : (
                    <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                      <span style={{
                        backgroundColor: '#10b981',
                        color: 'white',
                        padding: '0.25rem 0.75rem',
                        borderRadius: '1rem',
                        fontSize: '0.75rem',
                        fontWeight: '500',
                      }}>
                        メンバー
                      </span>
                      <button
                        onClick={() => fetchGroupMembers(group.id)}
                        style={{
                          backgroundColor: '#3b82f6',
                          color: 'white',
                          border: 'none',
                          padding: '0.25rem 0.75rem',
                          borderRadius: '0.25rem',
                          fontSize: '0.75rem',
                          cursor: 'pointer',
                          fontWeight: '500',
                        }}
                      >
                        メンバー一覧
                      </button>
                      <button
                        onClick={() => handleLeaveGroup(group.id)}
                        disabled={memberActionLoading === group.id}
                        style={{
                          backgroundColor: memberActionLoading === group.id ? '#9ca3af' : '#f59e0b',
                          color: 'white',
                          border: 'none',
                          padding: '0.25rem 0.75rem',
                          borderRadius: '0.25rem',
                          fontSize: '0.75rem',
                          cursor: memberActionLoading === group.id ? 'not-allowed' : 'pointer',
                          fontWeight: '500',
                        }}
                      >
                        {memberActionLoading === group.id ? '退会中...' : '退会'}
                      </button>
                    </div>
                  )}
                </div>
              </div>
              
              {/* メンバー追加フォーム */}
              {showAddMemberForm === group.id && (
                <div style={{
                  marginTop: '1rem',
                  padding: '1rem',
                  backgroundColor: '#f9fafb',
                  borderRadius: '0.375rem',
                  border: '1px solid #e5e7eb'
                }}>
                  <h4 style={{ fontSize: '1rem', fontWeight: '600', marginBottom: '1rem' }}>
                    メンバーを追加
                  </h4>
                  <div style={{ marginBottom: '1rem' }}>
                    <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' }}>
                      ユーザー名で検索
                    </label>
                    <input
                      type="text"
                      value={memberDisplayName}
                      onChange={(e) => {
                        setMemberDisplayName(e.target.value);
                        handleSearchUsers(e.target.value);
                      }}
                      style={{
                        width: '100%',
                        padding: '0.75rem',
                        border: '1px solid #d1d5db',
                        borderRadius: '0.375rem',
                        fontSize: '1rem',
                        boxSizing: 'border-box',
                      }}
                      placeholder="表示ユーザー名を入力"
                    />
                  </div>
                  
                  {/* 検索結果 */}
                  {searchResults.length > 0 && (
                    <div style={{
                      marginBottom: '1rem',
                      maxHeight: '200px',
                      overflowY: 'auto',
                      border: '1px solid #d1d5db',
                      borderRadius: '0.375rem',
                      backgroundColor: 'white'
                    }}>
                      {searchResults.map((searchUser) => (
                        <div
                          key={searchUser.id}
                          onClick={() => setMemberDisplayName(searchUser.display_name)}
                          style={{
                            padding: '0.75rem',
                            borderBottom: '1px solid #e5e7eb',
                            cursor: 'pointer',
                          }}
                          onMouseEnter={(e) => {
                            e.currentTarget.style.backgroundColor = '#f3f4f6';
                          }}
                          onMouseLeave={(e) => {
                            e.currentTarget.style.backgroundColor = 'white';
                          }}
                        >
                          <div style={{ fontWeight: '500' }}>{searchUser.display_name}</div>
                          {searchUser.bio && (
                            <div style={{ fontSize: '0.875rem', color: '#6b7280' }}>{searchUser.bio}</div>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                  
                  <div style={{ display: 'flex', gap: '0.5rem' }}>
                    <button
                      onClick={() => handleAddMemberByName(group.id)}
                      disabled={!memberDisplayName.trim() || joinLoading === group.id}
                      style={{
                        padding: '0.5rem 1rem',
                        border: 'none',
                        backgroundColor: (!memberDisplayName.trim() || joinLoading === group.id) ? '#9ca3af' : '#10b981',
                        color: 'white',
                        borderRadius: '0.375rem',
                        fontSize: '0.875rem',
                        fontWeight: '500',
                        cursor: (!memberDisplayName.trim() || joinLoading === group.id) ? 'not-allowed' : 'pointer',
                      }}
                    >
                      {joinLoading === group.id ? '追加中...' : '追加'}
                    </button>
                    <button
                      onClick={() => {
                        setShowAddMemberForm(null);
                        setMemberDisplayName('');
                        setSearchResults([]);
                      }}
                      style={{
                        padding: '0.5rem 1rem',
                        border: '1px solid #d1d5db',
                        backgroundColor: 'white',
                        color: '#374151',
                        borderRadius: '0.375rem',
                        fontSize: '0.875rem',
                        cursor: 'pointer',
                      }}
                    >
                      キャンセル
                    </button>
                  </div>
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {/* グループ編集モーダル */}
      {editingGroup && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 1000,
        }}>
          <div style={{
            backgroundColor: 'white',
            padding: '2rem',
            borderRadius: '0.5rem',
            width: '100%',
            maxWidth: '500px',
            margin: '1rem',
          }}>
            <h2 style={{ fontSize: '1.5rem', fontWeight: '700', marginBottom: '1.5rem' }}>
              グループを編集
            </h2>
            
            <form onSubmit={handleUpdateGroup}>
              <div style={{ marginBottom: '1.5rem' }}>
                <label htmlFor="editGroupName" style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' }}>
                  グループ名 <span style={{ color: '#dc2626' }}>*</span>
                </label>
                <input
                  id="editGroupName"
                  type="text"
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                  maxLength={100}
                  required
                  style={{
                    width: '100%',
                    padding: '0.75rem',
                    border: '1px solid #d1d5db',
                    borderRadius: '0.375rem',
                    fontSize: '1rem',
                    boxSizing: 'border-box',
                  }}
                  placeholder="グループ名を入力"
                />
                <div style={{ fontSize: '0.75rem', color: '#6b7280', marginTop: '0.25rem' }}>
                  {editName.length}/100 文字
                </div>
              </div>

              <div style={{ marginBottom: '1.5rem' }}>
                <label htmlFor="editGroupDescription" style={{ display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' }}>
                  説明 <span style={{ color: '#dc2626' }}>*</span>
                </label>
                <textarea
                  id="editGroupDescription"
                  value={editDescription}
                  onChange={(e) => setEditDescription(e.target.value)}
                  maxLength={500}
                  required
                  rows={4}
                  style={{
                    width: '100%',
                    padding: '0.75rem',
                    border: '1px solid #d1d5db',
                    borderRadius: '0.375rem',
                    fontSize: '1rem',
                    fontFamily: 'inherit',
                    resize: 'vertical',
                    boxSizing: 'border-box',
                  }}
                  placeholder="グループの説明を入力"
                />
                <div style={{ fontSize: '0.75rem', color: '#6b7280', marginTop: '0.25rem' }}>
                  {editDescription.length}/500 文字
                </div>
              </div>

              <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
                <button
                  type="button"
                  onClick={() => {
                    setEditingGroup(null);
                    setEditName('');
                    setEditDescription('');
                  }}
                  style={{
                    padding: '0.75rem 1.5rem',
                    border: '1px solid #d1d5db',
                    backgroundColor: 'white',
                    color: '#374151',
                    borderRadius: '0.375rem',
                    fontWeight: '500',
                    cursor: 'pointer',
                  }}
                >
                  キャンセル
                </button>
                <button
                  type="submit"
                  disabled={creating}
                  style={{
                    padding: '0.75rem 1.5rem',
                    border: 'none',
                    backgroundColor: creating ? '#9ca3af' : '#2563eb',
                    color: 'white',
                    borderRadius: '0.375rem',
                    fontWeight: '500',
                    cursor: creating ? 'not-allowed' : 'pointer',
                  }}
                >
                  {creating ? '更新中...' : 'グループを更新'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* グループ削除確認モーダル */}
      {deletingGroup && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 1000,
        }}>
          <div style={{
            backgroundColor: 'white',
            padding: '2rem',
            borderRadius: '0.5rem',
            width: '100%',
            maxWidth: '400px',
            margin: '1rem',
          }}>
            <h2 style={{ fontSize: '1.5rem', fontWeight: '700', marginBottom: '1rem', color: '#dc2626' }}>
              グループ削除の確認
            </h2>
            
            <p style={{ marginBottom: '1rem', color: '#374151' }}>
              「{deletingGroup.name}」を削除しますか？
            </p>
            
            <p style={{ marginBottom: '1.5rem', fontSize: '0.875rem', color: '#6b7280' }}>
              この操作は取り消せません。グループに関連する投稿はグループから削除され、メンバーもすべて削除されます。
            </p>
            
            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
              <button
                type="button"
                onClick={() => setDeletingGroup(null)}
                disabled={deleteLoading}
                style={{
                  padding: '0.75rem 1.5rem',
                  border: '1px solid #d1d5db',
                  backgroundColor: 'white',
                  color: '#374151',
                  borderRadius: '0.375rem',
                  fontWeight: '500',
                  cursor: deleteLoading ? 'not-allowed' : 'pointer',
                }}
              >
                キャンセル
              </button>
              <button
                type="button"
                onClick={handleDeleteGroup}
                disabled={deleteLoading}
                style={{
                  padding: '0.75rem 1.5rem',
                  border: 'none',
                  backgroundColor: deleteLoading ? '#9ca3af' : '#dc2626',
                  color: 'white',
                  borderRadius: '0.375rem',
                  fontWeight: '500',
                  cursor: deleteLoading ? 'not-allowed' : 'pointer',
                }}
              >
                {deleteLoading ? '削除中...' : '削除する'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* メンバー一覧モーダル */}
      {showMembersFor && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 1000,
        }}>
          <div style={{
            backgroundColor: 'white',
            padding: '2rem',
            borderRadius: '0.5rem',
            width: '100%',
            maxWidth: '600px',
            margin: '1rem',
            maxHeight: '80vh',
            overflowY: 'auto',
          }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
              <h2 style={{ fontSize: '1.5rem', fontWeight: '700' }}>
                グループメンバー
              </h2>
              <button
                onClick={() => {
                  setShowMembersFor(null);
                  setMembers([]);
                }}
                style={{
                  backgroundColor: '#6b7280',
                  color: 'white',
                  border: 'none',
                  padding: '0.5rem 1rem',
                  borderRadius: '0.25rem',
                  fontSize: '0.875rem',
                  cursor: 'pointer',
                  fontWeight: '500',
                }}
              >
                閉じる
              </button>
            </div>
            
            {membersLoading ? (
              <div style={{ textAlign: 'center', padding: '2rem' }}>
                <div>メンバー一覧を読み込み中...</div>
              </div>
            ) : members.length === 0 ? (
              <div style={{ textAlign: 'center', padding: '2rem', color: '#6b7280' }}>
                メンバーがいません
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                {members.map((member) => {
                  const currentGroup = groups.find(g => g.id === showMembersFor);
                  const isOwner = currentGroup?.owner_id === user?.id;
                  const isSelf = member.id === user?.id;
                  
                  return (
                    <div
                      key={member.id}
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        padding: '1rem',
                        border: '1px solid #e5e7eb',
                        borderRadius: '0.375rem',
                        backgroundColor: '#f9fafb',
                      }}
                    >
                      <div>
                        <div style={{ fontWeight: '600', marginBottom: '0.25rem' }}>
                          {member.display_name}
                          {currentGroup?.owner_id === member.id && (
                            <span style={{
                              marginLeft: '0.5rem',
                              backgroundColor: '#fbbf24',
                              color: 'white',
                              padding: '0.125rem 0.5rem',
                              borderRadius: '0.75rem',
                              fontSize: '0.75rem',
                              fontWeight: '500',
                            }}>
                              オーナー
                            </span>
                          )}
                          {isSelf && (
                            <span style={{
                              marginLeft: '0.5rem',
                              backgroundColor: '#10b981',
                              color: 'white',
                              padding: '0.125rem 0.5rem',
                              borderRadius: '0.75rem',
                              fontSize: '0.75rem',
                              fontWeight: '500',
                            }}>
                              あなた
                            </span>
                          )}
                        </div>
                        {member.bio && (
                          <div style={{ fontSize: '0.875rem', color: '#6b7280' }}>
                            {member.bio}
                          </div>
                        )}
                      </div>
                      
                      {isOwner && !isSelf && currentGroup?.owner_id !== member.id && (
                        <button
                          onClick={() => handleRemoveMember(showMembersFor, member.id)}
                          disabled={memberActionLoading === member.id}
                          style={{
                            backgroundColor: memberActionLoading === member.id ? '#9ca3af' : '#dc2626',
                            color: 'white',
                            border: 'none',
                            padding: '0.5rem 1rem',
                            borderRadius: '0.25rem',
                            fontSize: '0.875rem',
                            cursor: memberActionLoading === member.id ? 'not-allowed' : 'pointer',
                            fontWeight: '500',
                          }}
                        >
                          {memberActionLoading === member.id ? '除名中...' : '除名'}
                        </button>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};