module.exports = {
  'posting-app-api': {
    input: '../api/schema.yaml',
    output: {
      target: 'src/generated/api.ts',
      client: 'axios',
      mode: 'split',
      schemas: 'src/generated/models',
      override: {
        mutator: {
          path: './src/utils/api-mutator.ts',
          name: 'apiMutator',
        },
        operations: {
          // 認証が必要なエンドポイントにはデフォルトでBearer tokenを使用
          Auth: {
            'Auth.login': {
              mutator: './src/utils/api-mutator.ts',
            },
            'Auth.register': {
              mutator: './src/utils/api-mutator.ts',
            },
          },
        },
      },
    },
    hooks: {
      afterAllFilesWrite: 'prettier --write',
    },
  },
};