// For more info, see https://github.com/storybookjs/eslint-plugin-storybook#configuration-flat-config-format
import storybook from "eslint-plugin-storybook";

import pluginNext from '@next/eslint-plugin-next'
import pluginReact from 'eslint-plugin-react'
import pluginReactHooks from 'eslint-plugin-react-hooks'
import pluginJsxA11y from 'eslint-plugin-jsx-a11y'
import tseslint from 'typescript-eslint'

export default tseslint.config({
  ignores: [
    'node_modules/**',
    '.next/**',
    'out/**',
    'build/**',
    'storybook-static/**',
    'next-env.d.ts',
  ],
}, // TypeScript configuration
...tseslint.configs.recommended, // React configuration
{
  files: ['**/*.{js,jsx,ts,tsx}'],
  plugins: {
    react: pluginReact,
    'react-hooks': pluginReactHooks,
    'jsx-a11y': pluginJsxA11y,
    '@next/next': pluginNext,
  },
  settings: {
    react: {
      version: 'detect',
    },
  },
  rules: {
    // React rules
    'react/react-in-jsx-scope': 'off',
    'react/prop-types': 'off',
    'react/jsx-no-target-blank': 'warn',

    // React Hooks rules
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn',

    // Next.js rules
    '@next/next/no-html-link-for-pages': 'error',
    '@next/next/no-img-element': 'warn',
    '@next/next/no-sync-scripts': 'error',

    // TypeScript rules - strengthened
    '@typescript-eslint/no-unused-vars': [
      'error',
      { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
    ],
    '@typescript-eslint/no-explicit-any': 'error',
    '@typescript-eslint/no-empty-object-type': 'off',

    // Accessibility rules
    'jsx-a11y/alt-text': 'warn',
    'jsx-a11y/anchor-is-valid': 'warn',
  },
}, storybook.configs["flat/recommended"]);
