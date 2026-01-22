import { defineRouting } from 'next-intl/routing'

export const routing = defineRouting({
  locales: ['mn', 'en'] as const,
  defaultLocale: 'mn',
  localePrefix: 'always',
})
