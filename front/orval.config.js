module.exports = {
  'posting-app': {
    input: '../api/schema.yaml',
    output: {
      mode: 'split',
      target: 'src/generated/api.ts',
      schemas: 'src/generated/models',
      client: 'react-query',
      mock: false,
      clean: true,
      prettier: true,
      tsconfig: 'tsconfig.json',
      override: {
        mutator: {
          path: './src/utils/api-mutator.ts',
          name: 'customInstance',
        },
        query: {
          useQuery: true,
          useMutation: true,
          signal: true,
        },
      },
    },
    hooks: {
      afterAllFilesWrite: [
        {
          command: 'npx prettier --write src/generated/**/*.ts',
        },
      ],
    },
  },
};