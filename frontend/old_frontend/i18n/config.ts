export const locales = ['mn', 'en'] as const
export type Locale = (typeof locales)[number]
export const defaultLocale: Locale = 'mn'
export const isLocale = (l: string): l is Locale => (locales as readonly string[]).includes(l)
