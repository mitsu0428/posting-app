import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { subscriptionApi } from '../utils/api';
import { SubscriptionStatus } from '../types';

export const Subscription: React.FC = () => {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [subscriptionStatus, setSubscriptionStatus] =
    useState<SubscriptionStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      fetchSubscriptionStatus();
    }
  }, [isAuthenticated, authLoading]);

  const fetchSubscriptionStatus = async () => {
    try {
      setLoading(true);
      const response = await subscriptionApi.getStatus();
      setSubscriptionStatus(response);
    } catch (err: any) {
      setError(
        err.response?.data?.message || 'Failed to fetch subscription status'
      );
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCheckoutSession = async () => {
    try {
      setCreating(true);
      setError('');

      const response = await subscriptionApi.createCheckoutSession();

      // Redirect to Stripe checkout
      window.location.href = response.url;
    } catch (err: any) {
      setError(
        err.response?.data?.message || 'Failed to create checkout session'
      );
    } finally {
      setCreating(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return '#059669';
      case 'past_due':
        return '#d97706';
      case 'canceled':
        return '#dc2626';
      case 'inactive':
        return '#6b7280';
      default:
        return '#6b7280';
    }
  };

  const getStatusMessage = (status: string) => {
    switch (status) {
      case 'active':
        return 'サブスクリプションは有効です。投稿や返信が可能です。';
      case 'past_due':
        return 'お支払いが滞っています。サービスを継続利用するにはお支払い方法を更新してください。';
      case 'canceled':
        return 'サブスクリプションはキャンセルされました。再度登録すると投稿や返信が可能になります。';
      case 'inactive':
        return '現在有効なサブスクリプションがありません。投稿や返信を行うには登録が必要です。';
      default:
        return 'サブスクリプションのステータスが不明です。';
    }
  };

  if (authLoading) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <div>Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <h2>Authentication Required</h2>
        <p>Please log in to access subscription management.</p>
      </div>
    );
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem' }}>
        <div>Loading subscription information...</div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '600px', margin: '0 auto' }}>
      <h1 style={{ fontSize: '2rem', fontWeight: '700', marginBottom: '2rem' }}>
        サブスクリプション
      </h1>

      {error && (
        <div
          style={{
            backgroundColor: '#fef2f2',
            border: '1px solid #fecaca',
            color: '#b91c1c',
            padding: '0.75rem',
            borderRadius: '0.375rem',
            marginBottom: '2rem',
          }}
        >
          {error}
        </div>
      )}

      <div
        style={{
          backgroundColor: 'white',
          border: '1px solid #e5e7eb',
          borderRadius: '0.5rem',
          padding: '2rem',
          marginBottom: '2rem',
        }}
      >
        <h2
          style={{
            fontSize: '1.5rem',
            fontWeight: '600',
            marginBottom: '1rem',
          }}
        >
          現在のステータス
        </h2>

        <div style={{ marginBottom: '1.5rem' }}>
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '1rem',
              marginBottom: '1rem',
            }}
          >
            <span style={{ fontSize: '1.125rem', fontWeight: '500' }}>
              ステータス:
            </span>
            <span
              style={{
                backgroundColor: getStatusColor(
                  subscriptionStatus?.status || 'inactive'
                ),
                color: 'white',
                padding: '0.5rem 1rem',
                borderRadius: '0.375rem',
                fontSize: '0.875rem',
                fontWeight: '500',
                textTransform: 'uppercase',
              }}
            >
              {subscriptionStatus?.status || 'inactive'}
            </span>
          </div>

          <p style={{ color: '#6b7280', lineHeight: '1.6' }}>
            {getStatusMessage(subscriptionStatus?.status || 'inactive')}
          </p>

          {subscriptionStatus?.current_period_end && (
            <p
              style={{
                color: '#6b7280',
                fontSize: '0.875rem',
                marginTop: '0.5rem',
              }}
            >
              Current period ends:{' '}
              {new Date(
                subscriptionStatus.current_period_end
              ).toLocaleDateString()}
            </p>
          )}
        </div>

        {subscriptionStatus?.status === 'inactive' && (
          <div
            style={{
              padding: '1.5rem',
              backgroundColor: '#f0f9ff',
              border: '1px solid #bae6fd',
              borderRadius: '0.375rem',
              marginBottom: '1.5rem',
            }}
          >
            <h3
              style={{
                fontSize: '1.125rem',
                fontWeight: '600',
                color: '#1e40af',
                marginBottom: '1rem',
              }}
            >
              Subscribe to Posting App Premium
            </h3>
            <div style={{ marginBottom: '1.5rem' }}>
              <h4
                style={{
                  fontSize: '1rem',
                  fontWeight: '500',
                  marginBottom: '0.5rem',
                }}
              >
                Features included:
              </h4>
              <ul
                style={{
                  listStyle: 'disc',
                  paddingLeft: '1.5rem',
                  color: '#374151',
                }}
              >
                <li>Create unlimited posts</li>
                <li>Reply to posts (anonymous or with your name)</li>
                <li>Upload thumbnail images for your posts</li>
                <li>Access to all approved content</li>
                <li>Priority customer support</li>
              </ul>
            </div>
            <div style={{ marginBottom: '1.5rem' }}>
              <div
                style={{
                  fontSize: '2rem',
                  fontWeight: '700',
                  color: '#1e40af',
                }}
              >
                $9.99
                <span
                  style={{
                    fontSize: '1rem',
                    fontWeight: '400',
                    color: '#6b7280',
                  }}
                >
                  /month
                </span>
              </div>
            </div>
            <button
              onClick={handleCreateCheckoutSession}
              disabled={creating}
              style={{
                width: '100%',
                backgroundColor: creating ? '#9ca3af' : '#2563eb',
                color: 'white',
                padding: '0.75rem 1.5rem',
                border: 'none',
                borderRadius: '0.375rem',
                fontSize: '1rem',
                fontWeight: '500',
                cursor: creating ? 'not-allowed' : 'pointer',
              }}
            >
              {creating ? 'Creating checkout session...' : 'Subscribe Now'}
            </button>
          </div>
        )}

        {subscriptionStatus?.status === 'active' && (
          <div
            style={{
              padding: '1.5rem',
              backgroundColor: '#f0fdf4',
              border: '1px solid #bbf7d0',
              borderRadius: '0.375rem',
            }}
          >
            <h3
              style={{
                fontSize: '1.125rem',
                fontWeight: '600',
                color: '#15803d',
              }}
            >
              サブスクリプションは有効です
            </h3>
          </div>
        )}

        {subscriptionStatus?.status === 'past_due' && (
          <div
            style={{
              padding: '1.5rem',
              backgroundColor: '#fef3c7',
              border: '1px solid #fbbf24',
              borderRadius: '0.375rem',
            }}
          >
            <h3
              style={{
                fontSize: '1.125rem',
                fontWeight: '600',
                color: '#92400e',
                marginBottom: '1rem',
              }}
            >
              Payment Required
            </h3>
            <p style={{ color: '#92400e', marginBottom: '1rem' }}>
              Your payment is past due. Please update your payment method to
              continue using the service.
            </p>
            <button
              onClick={handleCreateCheckoutSession}
              disabled={creating}
              style={{
                backgroundColor: creating ? '#9ca3af' : '#d97706',
                color: 'white',
                padding: '0.75rem 1.5rem',
                border: 'none',
                borderRadius: '0.375rem',
                fontSize: '0.875rem',
                fontWeight: '500',
                cursor: creating ? 'not-allowed' : 'pointer',
              }}
            >
              {creating ? 'Processing...' : 'Update Payment Method'}
            </button>
          </div>
        )}

        {subscriptionStatus?.status === 'canceled' && (
          <div
            style={{
              padding: '1.5rem',
              backgroundColor: '#fef2f2',
              border: '1px solid #fecaca',
              borderRadius: '0.375rem',
            }}
          >
            <h3
              style={{
                fontSize: '1.125rem',
                fontWeight: '600',
                color: '#b91c1c',
                marginBottom: '1rem',
              }}
            >
              Subscription Canceled
            </h3>
            <p style={{ color: '#b91c1c', marginBottom: '1rem' }}>
              Your subscription has been canceled. You can resubscribe at any
              time to regain access to premium features.
            </p>
            <button
              onClick={handleCreateCheckoutSession}
              disabled={creating}
              style={{
                backgroundColor: creating ? '#9ca3af' : '#2563eb',
                color: 'white',
                padding: '0.75rem 1.5rem',
                border: 'none',
                borderRadius: '0.375rem',
                fontSize: '0.875rem',
                fontWeight: '500',
                cursor: creating ? 'not-allowed' : 'pointer',
              }}
            >
              {creating ? 'Processing...' : 'Resubscribe'}
            </button>
          </div>
        )}
      </div>

      <div
        style={{
          backgroundColor: 'white',
          border: '1px solid #e5e7eb',
          borderRadius: '0.5rem',
          padding: '2rem',
        }}
      >
        <h2
          style={{
            fontSize: '1.25rem',
            fontWeight: '600',
            marginBottom: '1rem',
          }}
        >
          請求情報の管理
        </h2>
        <p style={{ color: '#6b7280', marginBottom: '1rem' }}>
          請求情報の管理、支払い方法の更新、請求書の表示については、Stripeカスタマーポータルをご利用ください。
        </p>
        <p style={{ color: '#6b7280', fontSize: '0.875rem' }}>
          請求の更新に関するメールが届き、Stripeの安全なポータルを通じてサブスクリプションを管理できます。
        </p>
      </div>
    </div>
  );
};
