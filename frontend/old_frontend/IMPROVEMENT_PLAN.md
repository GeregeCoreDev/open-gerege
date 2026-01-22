# Дэлхийн Түвшний Код Сайжруулалтын Төлөвлөгөө

> Next.js 15 Enterprise Application - World-Class Standards

---

## Одоогийн Байдлын Үнэлгээ

| Категори | Одоогийн | Зорилго | Статус |
|----------|----------|---------|--------|
| Type Safety | 75% | 95% | Сайжруулах |
| Test Coverage | ~20% | 80% | Шаардлагатай |
| Error Handling | 60% | 95% | Сайжруулах |
| Performance | 70% | 90% | Сайжруулах |
| Security | 85% | 95% | Бага сайжруулалт |
| Documentation | 30% | 80% | Шаардлагатай |
| Accessibility | 50% | 90% | Сайжруулах |

---

## 1. АЛДАА БОЛОВСРУУЛАЛТ (Error Handling)

### 1.1 Silent Error Handling Засах

**Асуудал:** `org.ts` болон бусад store-уудад алдааг чимээгүй алгасаж байна.

```typescript
// ❌ Одоогийн - lib/stores/org.ts:54-56
} catch {
  // Silent fail
}

// ✅ Сайжруулсан
} catch (error) {
  console.error('[OrgStore] Failed to fetch organizations:', error);
  set({ status: 'failed', error: error instanceof Error ? error.message : 'Unknown error' });
}
```

### 1.2 Centralized Error Logger нэмэх

**Шинэ файл:** `lib/logger.ts`

```typescript
type LogLevel = 'debug' | 'info' | 'warn' | 'error';

interface LogEntry {
  level: LogLevel;
  message: string;
  context?: Record<string, unknown>;
  timestamp: Date;
}

export const logger = {
  debug: (msg: string, ctx?: object) => log('debug', msg, ctx),
  info: (msg: string, ctx?: object) => log('info', msg, ctx),
  warn: (msg: string, ctx?: object) => log('warn', msg, ctx),
  error: (msg: string, ctx?: object) => log('error', msg, ctx),
};

// Production-д Sentry/LogRocket-руу илгээх
function log(level: LogLevel, message: string, context?: object) {
  const entry: LogEntry = { level, message, context, timestamp: new Date() };

  if (process.env.NODE_ENV === 'development') {
    console[level](message, context);
  } else {
    // Send to monitoring service
    sendToMonitoring(entry);
  }
}
```

### 1.3 API Error Messages Локализаци

```typescript
// lib/api.ts дээр
import { getTranslations } from 'next-intl/server';

// Error messages-ийг локалаас унших
const errorMessages = {
  network: { mn: 'Сүлжээний алдаа', en: 'Network error' },
  timeout: { mn: 'Хугацаа дууссан', en: 'Request timeout' },
  unauthorized: { mn: 'Нэвтрэх шаардлагатай', en: 'Authentication required' },
};
```

---

## 2. ТЕСТ ХАМРАХ ХҮРЭЭ (Test Coverage)

### 2.1 Component Testing

**Шаардлагатай тестүүд:**

| Component | Priority | Status |
|-----------|----------|--------|
| `UserBootstrap` | P0 | Хийх |
| `AuthBootstrap` | P0 | Хийх |
| `useApiQuery` | P0 | Хийх |
| `useApiMutation` | P0 | Хийх |
| `useUserStore` | P1 | Хийх |
| `useMenuStore` | P1 | Хийх |
| `Sidebar` | P1 | Хийх |
| `ErrorBoundary` | P2 | Хийх |

### 2.2 Integration Tests нэмэх

```typescript
// __tests__/integration/auth-flow.test.tsx
describe('Authentication Flow', () => {
  it('should redirect to login when unauthorized', async () => {
    // Test 401 handling
  });

  it('should load user profile after login', async () => {
    // Test bootstrap sequence
  });

  it('should clear stores on logout', async () => {
    // Test cleanup
  });
});
```

### 2.3 E2E Testing (Playwright)

```bash
npm install -D @playwright/test
```

```typescript
// e2e/login.spec.ts
test('user can login and see dashboard', async ({ page }) => {
  await page.goto('/mn/login');
  await page.fill('[name=email]', 'test@example.com');
  await page.fill('[name=password]', 'password');
  await page.click('button[type=submit]');
  await expect(page).toHaveURL('/mn/app/dashboard');
});
```

---

## 3. PERFORMANCE OPTIMIZATION

### 3.1 Request Deduplication

```typescript
// hooks/useApiQuery.ts - Сайжруулсан
const inFlightRequests = new Map<string, Promise<unknown>>();

export function useApiQuery<T>(options: UseApiQueryOptions<T>) {
  const fetchData = useCallback(async () => {
    const key = `${options.path}:${JSON.stringify(options.query)}`;

    // Check for in-flight request
    if (inFlightRequests.has(key)) {
      return inFlightRequests.get(key) as Promise<T>;
    }

    const promise = api.get<T>(options.path, { query: options.query });
    inFlightRequests.set(key, promise);

    try {
      return await promise;
    } finally {
      inFlightRequests.delete(key);
    }
  }, [options.path, options.query]);

  // ... rest of hook
}
```

### 3.2 Component Memoization

```typescript
// components/layout/sidebar/index.tsx
import { memo, useMemo } from 'react';

export const Sidebar = memo(function Sidebar() {
  const menus = useMenuStore((s) => s.menus);

  const menuTree = useMemo(() => buildMenuTree(menus), [menus]);

  return <nav>{/* ... */}</nav>;
});
```

### 3.3 Bundle Size Optimization

```typescript
// next.config.ts
const nextConfig = {
  experimental: {
    optimizePackageImports: ['lucide-react', '@radix-ui/react-icons'],
  },
  // Dynamic imports for heavy components
};
```

### 3.4 Image Optimization

```typescript
// components/common/OptimizedImage.tsx
import Image from 'next/image';

export function OptimizedImage({ src, alt, ...props }) {
  return (
    <Image
      src={src}
      alt={alt}
      loading="lazy"
      placeholder="blur"
      blurDataURL="data:image/jpeg;base64,..."
      {...props}
    />
  );
}
```

---

## 4. TYPE SAFETY САЙЖРУУЛАЛТ

### 4.1 Strict Type Assertions Хасах

```typescript
// ❌ Одоогийн - lib/stores/user.ts:116
const e = error as Error;

// ✅ Сайжруулсан
function isError(error: unknown): error is Error {
  return error instanceof Error;
}

const message = isError(error) ? error.message : 'Unknown error';
```

### 4.2 API Response Types Сайжруулах

```typescript
// types/api.d.ts
declare namespace Api {
  interface SuccessResponse<T> {
    data: T;
    meta?: {
      total: number;
      page: number;
      limit: number;
    };
  }

  interface ErrorResponse {
    error: {
      code: string;
      message: string;
      details?: Record<string, string[]>;
    };
  }

  type Response<T> = SuccessResponse<T> | ErrorResponse;
}
```

### 4.3 Zod Schema Validation

```bash
npm install zod
```

```typescript
// lib/schemas/user.ts
import { z } from 'zod';

export const userSchema = z.object({
  id: z.number(),
  email: z.string().email(),
  name: z.string().min(1),
  role: z.enum(['admin', 'user', 'guest']),
});

export type User = z.infer<typeof userSchema>;

// API response validation
const response = await api.get('/users/me');
const user = userSchema.parse(response); // Runtime validation
```

---

## 5. STATE MANAGEMENT САЙЖРУУЛАЛТ

### 5.1 Store Selectors

```typescript
// lib/stores/user.ts
export const useUserStore = create<UserState>()(/* ... */);

// Selectors - Performance optimization
export const selectUser = (state: UserState) => state.user_info;
export const selectOrg = (state: UserState) => state.org_info;
export const selectIsLoading = (state: UserState) => state.status === 'loading';

// Usage
const user = useUserStore(selectUser); // Only re-renders when user changes
```

### 5.2 Immer Integration

```bash
npm install immer
```

```typescript
import { immer } from 'zustand/middleware/immer';

export const useMenuStore = create<MenuState>()(
  persist(
    immer((set) => ({
      menus: [],
      toggleGroup: (id) =>
        set((state) => {
          const group = state.openGroups.find((g) => g === id);
          if (group) {
            state.openGroups = state.openGroups.filter((g) => g !== id);
          } else {
            state.openGroups.push(id);
          }
        }),
    })),
    { name: 'menu-store' }
  )
);
```

### 5.3 DevTools Integration

```typescript
import { devtools } from 'zustand/middleware';

export const useUserStore = create<UserState>()(
  devtools(
    persist(
      (set, get) => ({/* ... */}),
      { name: 'user-store' }
    ),
    { name: 'UserStore', enabled: process.env.NODE_ENV === 'development' }
  )
);
```

---

## 6. ACCESSIBILITY (a11y) САЙЖРУУЛАЛТ

### 6.1 ARIA Labels

```typescript
// components/ui/button.tsx
export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ children, isLoading, ...props }, ref) => (
    <button
      ref={ref}
      aria-busy={isLoading}
      aria-disabled={props.disabled}
      {...props}
    >
      {isLoading ? (
        <span aria-hidden="true">
          <Spinner />
        </span>
      ) : null}
      <span className={isLoading ? 'sr-only' : ''}>{children}</span>
    </button>
  )
);
```

### 6.2 Focus Management

```typescript
// hooks/useFocusTrap.ts
export function useFocusTrap(isActive: boolean) {
  const containerRef = useRef<HTMLElement>(null);

  useEffect(() => {
    if (!isActive || !containerRef.current) return;

    const focusableElements = containerRef.current.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );

    // Trap focus within container
  }, [isActive]);

  return containerRef;
}
```

### 6.3 Screen Reader Announcements

```typescript
// components/common/Announcer.tsx
export function Announcer({ message }: { message: string }) {
  return (
    <div
      role="status"
      aria-live="polite"
      aria-atomic="true"
      className="sr-only"
    >
      {message}
    </div>
  );
}
```

---

## 7. SECURITY САЙЖРУУЛАЛТ

### 7.1 CSP Strict Mode

```typescript
// next.config.ts
const cspHeader = `
  default-src 'self';
  script-src 'self' 'nonce-{nonce}';
  style-src 'self' 'unsafe-inline';
  img-src 'self' blob: data: https://cdn.gerege.mn;
  font-src 'self';
  connect-src 'self' https://api.gerege.mn;
  frame-ancestors 'none';
  base-uri 'self';
  form-action 'self';
`;
```

### 7.2 Input Sanitization

```typescript
// lib/sanitize.ts
import DOMPurify from 'dompurify';

export function sanitizeHtml(dirty: string): string {
  return DOMPurify.sanitize(dirty, {
    ALLOWED_TAGS: ['b', 'i', 'em', 'strong', 'a'],
    ALLOWED_ATTR: ['href'],
  });
}

export function sanitizeInput(input: string): string {
  return input.replace(/[<>'"&]/g, '');
}
```

### 7.3 Rate Limiting (Client-side)

```typescript
// lib/rateLimit.ts
const requestCounts = new Map<string, { count: number; resetTime: number }>();

export function checkRateLimit(key: string, limit: number, windowMs: number): boolean {
  const now = Date.now();
  const record = requestCounts.get(key);

  if (!record || now > record.resetTime) {
    requestCounts.set(key, { count: 1, resetTime: now + windowMs });
    return true;
  }

  if (record.count >= limit) {
    return false;
  }

  record.count++;
  return true;
}
```

---

## 8. MONITORING & OBSERVABILITY

### 8.1 Sentry Integration

```bash
npm install @sentry/nextjs
```

```typescript
// sentry.client.config.ts
import * as Sentry from '@sentry/nextjs';

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  tracesSampleRate: 0.1,
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1.0,
  integrations: [
    Sentry.replayIntegration(),
    Sentry.browserTracingIntegration(),
  ],
});
```

### 8.2 Performance Monitoring

```typescript
// lib/performance.ts
export function measurePageLoad() {
  if (typeof window === 'undefined') return;

  const observer = new PerformanceObserver((list) => {
    for (const entry of list.getEntries()) {
      // Send to analytics
      analytics.track('performance', {
        name: entry.name,
        duration: entry.duration,
        startTime: entry.startTime,
      });
    }
  });

  observer.observe({ entryTypes: ['navigation', 'resource', 'largest-contentful-paint'] });
}
```

### 8.3 Error Boundary Integration

```typescript
// components/common/ErrorBoundary.tsx - Сайжруулсан
import * as Sentry from '@sentry/nextjs';

class ErrorBoundary extends Component {
  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    Sentry.captureException(error, {
      extra: {
        componentStack: errorInfo.componentStack,
      },
    });
  }
}
```

---

## 9. CODE QUALITY

### 9.1 Commented Code Цэвэрлэх

**Хасах файлууд:**
- `lib/stores/org.ts:74-84` - Organization change API
- `components/layout/UserBootstrap.tsx:107-111` - Module API
- Бусад commented code-ууд

### 9.2 ESLint Rules Чангатгах

```javascript
// eslint.config.mjs
export default [
  {
    rules: {
      '@typescript-eslint/no-explicit-any': 'error', // warn -> error
      '@typescript-eslint/explicit-function-return-type': 'warn',
      '@typescript-eslint/no-unused-vars': 'error',
      'no-console': ['warn', { allow: ['warn', 'error'] }],
    },
  },
];
```

### 9.3 Prettier + Husky + lint-staged

```bash
npm install -D husky lint-staged
npx husky init
```

```json
// package.json
{
  "lint-staged": {
    "*.{ts,tsx}": ["eslint --fix", "prettier --write"],
    "*.{json,md}": ["prettier --write"]
  }
}
```

---

## 10. DOCUMENTATION

### 10.1 Component Storybook

```bash
npx storybook@latest init
```

```typescript
// components/ui/button.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import { Button } from './button';

const meta: Meta<typeof Button> = {
  component: Button,
  tags: ['autodocs'],
};

export default meta;

export const Primary: StoryObj<typeof Button> = {
  args: {
    children: 'Click me',
    variant: 'default',
  },
};
```

### 10.2 API Documentation

```typescript
// lib/api.ts - JSDoc comments
/**
 * Makes an HTTP request to the API
 * @param path - API endpoint path
 * @param options - Request options
 * @returns Promise with typed response
 * @throws {APIError} When request fails
 * @example
 * const user = await api.get<User>('/users/me');
 */
export async function request<T>(path: string, options?: ApiOptions): Promise<T> {
  // ...
}
```

### 10.3 Architecture Decision Records (ADR)

```markdown
# ADR-001: State Management with Zustand

## Status
Accepted

## Context
We needed a state management solution for our Next.js application.

## Decision
We chose Zustand for its simplicity, small bundle size, and excellent TypeScript support.

## Consequences
- Simple API reduces boilerplate
- Easy to test with mock stores
- Persistence middleware handles localStorage
```

---

## ХЭРЭГЖҮҮЛЭХ ДАРААЛАЛ

### Phase 1: Foundation (Яаралтай)
1. [ ] Error logging system нэмэх
2. [ ] Silent error handling засах
3. [ ] Commented code цэвэрлэх
4. [ ] Type assertion-ууд засах

### Phase 2: Quality (Чухал)
5. [ ] Test coverage 80% хүргэх
6. [ ] Request deduplication нэмэх
7. [ ] Component memoization
8. [ ] ESLint rules чангатгах

### Phase 3: Observability (Дунд)
9. [ ] Sentry integration
10. [ ] Performance monitoring
11. [ ] Structured logging

### Phase 4: DX & Documentation (Хожим)
12. [ ] Storybook setup
13. [ ] API documentation
14. [ ] ADR бичиглэл

---

## ХҮЛЭЭГДЭЖ БУЙ ҮР ДҮН

| Metric | Before | After |
|--------|--------|-------|
| Type Safety | 75% | 95% |
| Test Coverage | 20% | 80% |
| Error Handling | 60% | 95% |
| Bundle Size | - | -20% |
| LCP | - | <2.5s |
| FID | - | <100ms |
| CLS | - | <0.1 |

---

*Энэ төлөвлөгөө нь Google, Meta, Stripe зэрэг дэлхийн шилдэг компаниудын frontend стандартад үндэслэсэн болно.*
