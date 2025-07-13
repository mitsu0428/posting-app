-- Seed data for development and testing

-- Insert sample users (password: testpass123)
INSERT INTO users (username, email, password_hash, subscription_status) VALUES
('山田太郎', 'yamada@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('佐藤花子', 'sato@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('田中一郎', 'tanaka@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active'),
('鈴木美咲', 'suzuki@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active')
ON CONFLICT (email) DO NOTHING;

-- Insert sample approved posts
INSERT INTO posts (title, content, status, user_id, thumbnail_url) VALUES
('初心者向けプログラミング学習法について', 
 'プログラミングを始めたばかりの方に向けて、効率的な学習方法をご紹介します。

1. 基礎をしっかりと学ぶ
まずは変数、条件分岐、ループなどの基本概念を理解することが重要です。

2. 実際にコードを書く
理論だけでなく、実際に手を動かしてコードを書くことで理解が深まります。

3. 小さなプロジェクトから始める
計算機アプリやToDoリストなど、簡単なアプリケーションから始めましょう。

皆さんの経験やおすすめの学習リソースがあれば、ぜひ共有してください！',
 'approved', 2, 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=400'),

('Web開発のトレンド2024', 
 '2024年のWeb開発トレンドについて議論しましょう。

現在注目されている技術：
- React Server Components
- Astro
- Tailwind CSS
- TypeScript
- Vite

特にReact Server Componentsは、Next.js 13+で導入され、サーバーサイドレンダリングの新しいアプローチとして注目されています。

皆さんが最近使ってみた技術や、今後学びたい技術について教えてください！',
 'approved', 3, 'https://images.unsplash.com/photo-1547658719-da2b51169166?w=400'),

('リモートワークでの生産性向上テクニック', 
 'リモートワークが当たり前になった今、生産性を向上させるためのテクニックを共有しませんか？

私が実践している方法：
- ポモドーロテクニック（25分作業、5分休憩）
- 朝の時間を最も重要なタスクに充てる
- Slackやメールのチェック時間を制限
- 作業環境の整備（デスク周りの整理整頓）

皆さんはどのような工夫をされていますか？効果的だった方法があれば、ぜひ教えてください！',
 'approved', 4, 'https://images.unsplash.com/photo-1522202176988-66273c2fd55f?w=400'),

('おすすめの技術書・学習リソース', 
 'プログラミングや技術分野で読んで良かった本、おすすめの学習リソースを紹介し合いましょう！

私のおすすめ：
📚 書籍
- 「リーダブルコード」- コードの品質向上に
- 「Clean Architecture」- アーキテクチャ設計に
- 「システム設計入門」- スケーラブルなシステム設計に

🌐 オンラインリソース
- freeCodeCamp - 無料でプログラミングを学べる
- MDN Web Docs - Web技術のリファレンス
- LeetCode - アルゴリズム問題の練習

皆さんのおすすめがあれば、ジャンル問わず教えてください！',
 'approved', 2, 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=400'),

('AI・機械学習の最新動向', 
 'ChatGPTをはじめとするAI技術の進歩が目覚ましいですね。

最近の注目トピック：
- GPT-4の活用事例
- Stable Diffusionによる画像生成
- プログラミング支援AI（GitHub Copilot等）
- AutoGPTなどの自律型AI

開発者として、これらの技術をどう活用していくか、皆さんの意見を聞かせてください。

実際に業務で使っている方がいれば、具体的な活用方法も教えていただけると嬉しいです！',
 'approved', 3, 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=400'),

('データベース設計のベストプラクティス', 
 'データベース設計で気をつけているポイントや、失敗談があれば共有しませんか？

基本的なポイント：
- 正規化の適切な適用
- インデックスの効果的な設計
- 命名規則の統一
- 制約の適切な設定

経験から学んだこと：
- 過度な正規化は性能問題を招くことがある
- 将来の拡張性を考慮した設計の重要性
- パフォーマンステストの必要性

皆さんの経験やノウハウを教えてください！',
 'approved', 4, 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=400');

-- Insert sample replies
INSERT INTO replies (post_id, content, user_id, is_anonymous) VALUES
(1, 'とても参考になりました！私も最初は基礎を疎かにして苦労しました。地道にコツコツが一番ですね。', 3, FALSE),
(1, '実際に手を動かすのが重要というのは本当にそうですね。最初は簡単なものでも達成感があります。', 4, FALSE),
(1, 'プロゲートやAtCoderなどのサービスも初心者におすすめです！', 2, TRUE),

(2, 'React Server Components、まだ触ったことがないので勉強してみます！', 2, FALSE),
(2, 'Astroが気になってます。静的サイト生成において、Next.jsとの使い分けはどう考えていますか？', 4, FALSE),
(2, 'TypeScriptはもう必須ですね。型安全性がもたらす開発体験の向上は計り知れません。', 3, TRUE),

(3, 'ポモドーロテクニック、私も実践してます！集中力が格段に上がりました。', 2, FALSE),
(3, '作業環境の整備、大事ですね。モニターのアームを導入してから肩こりが減りました。', 3, FALSE),

(4, '「リーダブルコード」は名著ですね！チーム開発では必読だと思います。', 3, FALSE),
(4, 'Udemyのコースもおすすめです。実践的な内容が多くて助かってます。', 4, TRUE),

(5, 'GitHub Copilotを使い始めたのですが、確かにコーディング効率が上がりますね。ただし、生成されたコードの理解は必要ですが。', 2, FALSE),
(5, 'AI技術の進歩は本当にすごいですが、基礎をしっかり学ぶことの重要性は変わらないと思います。', 4, FALSE),

(6, 'インデックスの張りすぎで更新性能が落ちた経験があります。バランスが大事ですね。', 2, TRUE),
(6, 'ER図をしっかり書くことで、チーム内での認識合わせもスムーズになりますね。', 3, FALSE);

-- Insert some pending posts (awaiting approval)
INSERT INTO posts (title, content, status, user_id) VALUES
('新しいフレームワークについて質問', 'Svelteについて教えてください。Reactとの違いは何でしょうか？', 'pending', 2),
('Docker環境での開発', 'Docker Composeを使った開発環境構築で困っています。', 'pending', 3);