import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { postApi } from '../utils/api';
import { getCategories, getGroups } from '../generated/api';
import { Category } from '../generated/models/category';
import { Group } from '../generated/models/group';

export const CreatePost: React.FC = () => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [thumbnail, setThumbnail] = useState<File | null>(null);
  const [thumbnailPreview, setThumbnailPreview] = useState<string>('');
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategoryIds, setSelectedCategoryIds] = useState<number[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const { user } = useAuth();
  const navigate = useNavigate();
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [categoriesData, groupsData] = await Promise.all([
          getCategories(),
          getGroups(),
        ]);
        setCategories(categoriesData || []);
        setGroups(groupsData || []);
      } catch (err) {
        console.error('Failed to fetch data:', err);
      }
    };

    fetchData();
  }, []);

  const handleCategoryToggle = (categoryId: number) => {
    setSelectedCategoryIds((prev) => {
      if (prev.includes(categoryId)) {
        return prev.filter((id) => id !== categoryId);
      } else if (prev.length < 5) {
        return [...prev, categoryId];
      } else {
        setError('最大5個までカテゴリを選択できます');
        return prev;
      }
    });
  };

  const handleThumbnailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];

    if (!file) return;

    // Check file size (5MB)
    if (file.size > 5 * 1024 * 1024) {
      setError('File size must be less than 5MB');
      return;
    }

    // Check file type
    if (!file.type.startsWith('image/')) {
      setError('Please select an image file');
      return;
    }

    if (!['image/jpeg', 'image/png'].includes(file.type)) {
      setError('Only JPEG and PNG images are allowed');
      return;
    }

    setThumbnail(file);
    setError('');

    // Create preview
    const reader = new FileReader();
    reader.onload = (e) => {
      setThumbnailPreview(e.target?.result as string);
    };
    reader.readAsDataURL(file);
  };

  const removeThumbnail = () => {
    setThumbnail(null);
    setThumbnailPreview('');
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!title.trim() || !content.trim()) {
      setError('Title and content are required');
      return;
    }

    if (title.length > 200) {
      setError('Title must be less than 200 characters');
      return;
    }

    if (content.length > 5000) {
      setError('Content must be less than 5000 characters');
      return;
    }

    if (user?.subscription_status !== 'active') {
      setError('Active subscription required to create posts');
      return;
    }

    try {
      setLoading(true);
      setError('');

      const formData = new FormData();
      formData.append('title', title.trim());
      formData.append('content', content.trim());

      if (selectedCategoryIds.length > 0) {
        formData.append('category_ids', selectedCategoryIds.join(','));
      }

      if (selectedGroupId) {
        formData.append('group_id', selectedGroupId.toString());
      }

      if (thumbnail) {
        formData.append('thumbnail', thumbnail);
      }

      await postApi.createPost(formData);
      navigate('/my-page');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create post');
    } finally {
      setLoading(false);
    }
  };

  if (user?.subscription_status !== 'active') {
    return (
      <div style={{ maxWidth: '600px', margin: '0 auto', textAlign: 'center' }}>
        <div
          style={{
            backgroundColor: '#fef3c7',
            border: '1px solid #fbbf24',
            borderRadius: '0.5rem',
            padding: '2rem',
          }}
        >
          <h2 style={{ color: '#92400e', marginBottom: '1rem' }}>
            Subscription Required
          </h2>
          <p style={{ color: '#92400e', marginBottom: '1.5rem' }}>
            You need an active subscription to create posts.
          </p>
          <button
            onClick={() => navigate('/subscription')}
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
            Manage Subscription
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <h1 style={{ fontSize: '2rem', fontWeight: '700', marginBottom: '2rem' }}>
        新規投稿
      </h1>

      <form
        onSubmit={handleSubmit}
        style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}
      >
        {error && (
          <div
            style={{
              backgroundColor: '#fef2f2',
              border: '1px solid #fecaca',
              color: '#b91c1c',
              padding: '0.75rem',
              borderRadius: '0.375rem',
            }}
          >
            {error}
          </div>
        )}

        <div>
          <label
            htmlFor="title"
            style={{
              display: 'block',
              fontSize: '0.875rem',
              fontWeight: '500',
              color: '#374151',
              marginBottom: '0.5rem',
            }}
          >
            Title <span style={{ color: '#dc2626' }}>*</span>
          </label>
          <input
            id="title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            maxLength={200}
            required
            style={{
              width: '100%',
              padding: '0.75rem',
              border: '1px solid #d1d5db',
              borderRadius: '0.375rem',
              fontSize: '1rem',
              boxSizing: 'border-box',
            }}
            placeholder="Enter post title"
          />
          <div
            style={{
              fontSize: '0.75rem',
              color: '#6b7280',
              marginTop: '0.25rem',
            }}
          >
            {title.length}/200 字
          </div>
        </div>

        <div>
          <label
            style={{
              display: 'block',
              fontSize: '0.875rem',
              fontWeight: '500',
              color: '#374151',
              marginBottom: '0.5rem',
            }}
          >
            カテゴリ (最大5個まで)
          </label>
          <div
            style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: '0.5rem',
              marginBottom: '0.5rem',
            }}
          >
            {categories.map((category) => (
              <button
                key={category.id}
                type="button"
                onClick={() => handleCategoryToggle(category.id)}
                style={{
                  padding: '0.5rem 1rem',
                  border: selectedCategoryIds.includes(category.id)
                    ? `2px solid ${category.color}`
                    : '1px solid #d1d5db',
                  backgroundColor: selectedCategoryIds.includes(category.id)
                    ? category.color
                    : 'white',
                  color: selectedCategoryIds.includes(category.id)
                    ? 'white'
                    : '#374151',
                  borderRadius: '1rem',
                  fontSize: '0.875rem',
                  cursor: 'pointer',
                  fontWeight: selectedCategoryIds.includes(category.id)
                    ? '600'
                    : '400',
                  transition: 'all 0.2s',
                }}
              >
                {category.name}
              </button>
            ))}
          </div>
          <div style={{ fontSize: '0.75rem', color: '#6b7280' }}>
            選択中: {selectedCategoryIds.length}/5
          </div>
        </div>

        <div>
          <label
            style={{
              display: 'block',
              fontSize: '0.875rem',
              fontWeight: '500',
              color: '#374151',
              marginBottom: '0.5rem',
            }}
          >
            グループ (任意)
          </label>
          <div style={{ marginBottom: '0.5rem' }}>
            <select
              value={selectedGroupId || ''}
              onChange={(e) =>
                setSelectedGroupId(
                  e.target.value ? Number(e.target.value) : null
                )
              }
              style={{
                width: '100%',
                padding: '0.75rem',
                border: '1px solid #d1d5db',
                borderRadius: '0.375rem',
                fontSize: '1rem',
                backgroundColor: 'white',
                boxSizing: 'border-box',
              }}
            >
              <option value="">公開投稿（グループなし）</option>
              {groups.map((group) => (
                <option key={group.id} value={group.id}>
                  {group.name}
                </option>
              ))}
            </select>
          </div>
          <div style={{ fontSize: '0.75rem', color: '#6b7280' }}>
            グループを選択すると、そのグループのメンバーのみが投稿を見ることができます
          </div>
        </div>

        <div>
          <label
            htmlFor="thumbnail"
            style={{
              display: 'block',
              fontSize: '0.875rem',
              fontWeight: '500',
              color: '#374151',
              marginBottom: '0.5rem',
            }}
          >
            サムネイル (任意)
          </label>
          <input
            ref={fileInputRef}
            id="thumbnail"
            type="file"
            accept="image/jpeg,image/png"
            onChange={handleThumbnailChange}
            style={{
              width: '100%',
              padding: '0.75rem',
              border: '1px solid #d1d5db',
              borderRadius: '0.375rem',
              boxSizing: 'border-box',
            }}
          />
          <div
            style={{
              fontSize: '0.75rem',
              color: '#6b7280',
              marginTop: '0.25rem',
            }}
          >
            JPEGまたはPNGで最大5MBまで
          </div>

          {thumbnailPreview && (
            <div
              style={{
                marginTop: '1rem',
                position: 'relative',
                display: 'inline-block',
              }}
            >
              <img
                src={thumbnailPreview}
                alt="Thumbnail preview"
                style={{
                  width: '200px',
                  height: '150px',
                  objectFit: 'cover',
                  borderRadius: '0.375rem',
                  border: '1px solid #d1d5db',
                }}
              />
              <button
                type="button"
                onClick={removeThumbnail}
                style={{
                  position: 'absolute',
                  top: '0.25rem',
                  right: '0.25rem',
                  backgroundColor: '#dc2626',
                  color: 'white',
                  border: 'none',
                  borderRadius: '50%',
                  width: '24px',
                  height: '24px',
                  cursor: 'pointer',
                  fontSize: '0.75rem',
                }}
              >
                ×
              </button>
            </div>
          )}
        </div>

        <div>
          <label
            htmlFor="content"
            style={{
              display: 'block',
              fontSize: '0.875rem',
              fontWeight: '500',
              color: '#374151',
              marginBottom: '0.5rem',
            }}
          >
            内容 <span style={{ color: '#dc2626' }}>*</span>
          </label>
          <textarea
            id="content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            maxLength={5000}
            required
            rows={15}
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
            placeholder="投稿する内容を入力..."
          />
          <div
            style={{
              fontSize: '0.75rem',
              color: '#6b7280',
              marginTop: '0.25rem',
            }}
          >
            {content.length}/5000 字
          </div>
        </div>

        <div
          style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}
        >
          <button
            type="button"
            onClick={() => navigate(-1)}
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
            disabled={loading}
            style={{
              padding: '0.75rem 1.5rem',
              border: 'none',
              backgroundColor: loading ? '#9ca3af' : '#2563eb',
              color: 'white',
              borderRadius: '0.375rem',
              fontWeight: '500',
              cursor: loading ? 'not-allowed' : 'pointer',
            }}
          >
            {loading ? '作成中...' : '作成'}
          </button>
        </div>
      </form>
    </div>
  );
};
