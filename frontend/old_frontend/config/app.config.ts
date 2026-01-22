/**
 * ⚙️ Application Configuration
 * 
 * Төвлөрсөн тохиргоо - бүх hardcoded values энд байна.
 * Environment-based configuration болон feature flags.
 * 
 * Benefits:
 * - ✅ Single source of truth
 * - ✅ Easy to modify
 * - ✅ Type-safe with TypeScript
 * - ✅ Environment-based values
 * 
 * Usage:
 * ```tsx
 * import { appConfig } from '@/config/app.config'
 * 
 * const [pageSize] = useState(appConfig.pagination.defaultPageSize)
 * ```
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

export const appConfig = {
  /**
   * 📄 Pagination Settings
   */
  pagination: {
    defaultPageSize: 50,
    pageSizeOptions: [10, 20, 50, 100, 200],
    maxPageSize: 500,
  },

  /**
   * 🌐 API Settings
   */
  api: {
    baseURL: process.env.NEXT_PUBLIC_API_BASE ?? '/api',
    timeout: 15000, // 15 seconds
    retries: 0, // No retries by default (can be enabled per request)
    retryDelay: 1000, // 1 second between retries
  },

  /**
   * 🎨 UI Settings
   */
  ui: {
    toastDuration: 3000, // 3 seconds
    progressBarUpdateInterval: 250, // ms
    animationDuration: 300, // ms
  },

  /**
   * 🌍 i18n Settings
   */
  i18n: {
    defaultLocale: 'mn' as const,
    supportedLocales: ['mn', 'en'] as const,
  },

  /**
   * 🎙️ Features (Feature Flags)
   */
  features: {
    enableVoiceRecognition: true,
    enableDarkMode: true,
    enableDevTools: process.env.NODE_ENV === 'development',
  },

  /**
   * 📊 Development Settings
   */
  dev: {
    enableDebugLogs: process.env.NODE_ENV === 'development',
    showQueryParams: process.env.NODE_ENV === 'development',
  },
} as const

/**
 * Type helper for config
 */
export type AppConfig = typeof appConfig

/**
 * Helper to check if feature is enabled
 */
export const isFeatureEnabled = (feature: keyof typeof appConfig.features): boolean => {
  return appConfig.features[feature]
}

