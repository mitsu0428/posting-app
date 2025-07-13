import React, { useState } from "react";
import { subscriptionAPI } from "../utils/api";
import { loadStripe } from "@stripe/stripe-js";

const stripePromise = loadStripe(
  process.env.REACT_APP_STRIPE_PUBLISHABLE_KEY || ""
);

const Subscription: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Basic styles
  const containerStyles = "max-w-2xl mx-auto my-12 p-5 text-center";
  const titleStyles = "text-3xl font-bold mb-6";
  const descriptionStyles = "mb-8 leading-relaxed";
  const planCardStyles = "border border-gray-300 rounded-lg p-8 bg-gray-50 mb-8";
  const planTitleStyles = "text-2xl font-semibold mb-4";
  const priceStyles = "text-4xl font-bold text-blue-600 my-4";
  const featureListStyles = "text-left mb-8";
  const errorStyles = "text-red-600 mb-4";
  const subscribeButtonStyles = "py-4 px-8 text-lg bg-blue-600 text-white border-none rounded-lg cursor-pointer min-w-48 disabled:bg-gray-400 disabled:cursor-not-allowed";
  const footerTextStyles = "mt-8 text-sm text-gray-600";

  const handleSubscribe = async () => {
    setLoading(true);
    setError("");

    try {
      const response = await subscriptionAPI.createCheckoutSession();
      const stripe = await stripePromise;

      if (!stripe) {
        throw new Error("Stripeの読み込みに失敗しました");
      }

      const { error } = await stripe.redirectToCheckout({
        sessionId: response.session_id,
      });

      if (error) {
        throw new Error(error.message);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "決済処理に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={containerStyles}>
      <h1 className={titleStyles}>サブスクリプション登録</h1>
      <p className={descriptionStyles}>
        掲示板アプリをご利用いただくには、サブスクリプション登録が必要です。
        <br />
        月額料金をお支払いいただくことで、すべての機能をご利用いただけます。
      </p>

      <div className={planCardStyles}>
        <h2 className={planTitleStyles}>プレミアムプラン</h2>
        <div className={priceStyles}>
          ¥500/月
        </div>
        <ul className={featureListStyles}>
          <li>投稿の作成・閲覧</li>
          <li>返信機能</li>
          <li>匿名投稿対応</li>
          <li>24時間サポート</li>
        </ul>
      </div>

      {error && (
        <div className={errorStyles}>{error}</div>
      )}

      <button
        onClick={handleSubscribe}
        disabled={loading}
        className={subscribeButtonStyles}
      >
        {loading ? "処理中..." : "サブスクリプション登録"}
      </button>

      <p className={footerTextStyles}>
        Stripeによる安全な決済処理を使用しています
      </p>
    </div>
  );
};

export default Subscription;
