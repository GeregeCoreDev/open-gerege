# Gerege Template - Project Master Prompt

**Author**: Sengum Soronzonbold, Developer  
**Company**: Gerege Core Team  
**Version**: 1.0.0  
**Last Updated**: December 7, 2025

---

## 📋 Project Overview

**Gerege Template** is a modern, enterprise-grade web application built with Next.js 15, designed to manage multi-organizational systems with advanced role-based access control (RBAC), modular architecture, and comprehensive internationalization support.

### Core Purpose
- Multi-tenant organization and user management
- Dynamic module and permission system
- Role-based access control with hierarchical systems
- Multi-language support (Mongolian/English)
- Integration with TPay, Business, and Admin subsystems

---

## 🏗️ Technology Stack

### Frontend Framework
- **Next.js**: 15.5.4 (App Router, React Server Components)
- **React**: 19.1.0
- **TypeScript**: 5.9.3
- **Build Tool**: Turbopack (for dev and production builds)

### Styling & UI
- **Tailwind CSS**: 4.x (latest)
- **Radix UI**: Headless UI components
- **Class Variance Authority (CVA)**: Component variant management
- **Lucide React**: Icon library
- **next-themes**: Dark/Light mode support

### State Management
- **Zustand**: 5.0.8 with middleware
- **zustand/middleware**: persist (localStorage sync)
- **Store Architecture**: Separate stores for user, org, role, system

### Forms & Validation
- **React Hook Form**: 7.65.0
- **Zod**: 4.1.12 (schema validation)
- **@hookform/resolvers**: Zod integration

### Internationalization
- **next-intl**: 4.3.12
- **Supported Languages**: Mongolian (mn), English (en)
- **Default Locale**: mn

### API & Data Fetching
- **Custom API Client**: `/lib/api.ts`
- **HTTP Library**: Native Fetch API with custom wrapper
- **Authentication**: Cookie-based sessions with credential inclusion

### UI Feedback
- **Sonner**: Toast notifications
- **Progress Components**: Custom loading states

### Development Tools
- **ESLint**: 9.38.0
- **Prettier**: 3.6.2 with Tailwind plugin
- **TypeScript Config**: Strict mode enabled

---

## 📁 Project Structure

```
next-template-v25/
├── app/                              # Next.js App Router
│   ├── [locale]/                    # Internationalized routes
│   │   ├── (home)/                  # Public home pages
│   │   │   └── home/
│   │   │       ├── components/
│   │   │       └── page.tsx
│   │   ├── (main)/                  # Protected main application
│   │   │   ├── admin/               # Admin system modules
│   │   │   │   ├── dashboard/
│   │   │   │   ├── user/
│   │   │   │   ├── role/
│   │   │   │   ├── organization/
│   │   │   │   ├── system/
│   │   │   │   └── ...
│   │   │   ├── app/                 # App system modules
│   │   │   ├── business/            # Business system modules
│   │   │   ├── tpay/                # TPay system modules
│   │   │   └── layout.tsx           # Main app layout (sidebar + header)
│   │   ├── (personal)/              # Personal user pages
│   │   │   ├── profile/
│   │   │   ├── settings/
│   │   │   ├── wallet/
│   │   │   └── layout.tsx
│   │   ├── callback/
│   │   ├── change-organization/
│   │   ├── change-system/
│   │   └── layout.tsx               # Locale-level layout
│   ├── favicon.ico
│   ├── globals.css
│   ├── layout.tsx                   # Root layout
│   └── page.tsx                     # Root redirect page
│
├── components/
│   ├── ui/                          # Radix-based UI primitives
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── dialog.tsx
│   │   ├── form.tsx
│   │   ├── input.tsx
│   │   ├── table.tsx
│   │   ├── pagination.tsx
│   │   ├── select.tsx
│   │   ├── sidebar.tsx
│   │   └── ... (30+ components)
│   ├── layout/                      # Layout components
│   │   ├── mainHeader.tsx
│   │   ├── mainSidebar.tsx
│   │   ├── homeHeader.tsx
│   │   ├── profileDropDown.tsx
│   │   ├── profileSidebar.tsx
│   │   ├── headerOrgSystemChanger.tsx
│   ├── common/                      # Shared components
│   │   ├── fileUpload.tsx
│   │   ├── mic.tsx
│   │   ├── userFind.tsx
│   │   └── subSystemRolePage.tsx
│   └── flag-icon/
│       └── flagIcon.tsx
│
├── lib/
│   ├── api.ts                       # Centralized API client
│   ├── logout.ts                    # Logout utility
│   ├── utils.ts                     # Common utilities (cn function)
│   ├── stores/                      # Zustand state stores
│   │   ├── user.ts                  # User/Organization state
│   │   ├── org.ts                   # Organization list state
│   │   ├── role.ts                  # User roles state
│   │   └── system.ts                # System modules state
│   ├── bootstrap/                   # App initialization
│   │   ├── UserBootstrap.tsx        # Load user data on mount
│   │   └── AuthBootstrap.tsx        # Setup auth handlers
│   └── utils/
│       ├── icon.tsx                 # Lucide icon wrapper
│       └── image.tsx                # Image utilities
│
├── i18n/
│   ├── config.ts                    # Locale configuration
│   ├── mn.json                      # Mongolian translations
│   ├── en.json                      # English translations
│   ├── navigation.ts                # Typed navigation
│   ├── request.ts                   # Server-side i18n setup
│   └── routing.ts                   # Routing configuration
│
├── types/
│   ├── global.d.ts                  # App namespace types
│   └── system.d.ts                  # System namespace types
│
├── hooks/
│   ├── use-mobile.ts                # Mobile detection hook
│   └── useSpeechRecognition.ts      # Speech recognition hook
│
├── public/
│   ├── logo/
│   ├── flag/
│   └── images/
│
├── middleware.ts                    # next-intl middleware
├── next.config.ts                   # Next.js configuration
├── tsconfig.json                    # TypeScript configuration
├── tailwind.config.ts               # Tailwind configuration
├── components.json                  # shadcn/ui configuration
├── ecosystem.config.cjs             # PM2 deployment config
├── package.json
└── README.md
```

---

## 🎯 Architecture Principles

### 1. **Server Components First**
- Use React Server Components by default
- Mark components as `'use client'` only when necessary (hooks, interactivity)
- Keep server-side data fetching in page components

### 2. **Type Safety**
- All API responses must have TypeScript interfaces in `types/global.d.ts`
- Use Zod schemas for form validation
- Leverage TypeScript strict mode

### 3. **Component Organization**
- **Atomic Design**: ui → common → layout → pages
- **Co-location**: Keep related files together (page + components folder)
- **Separation of Concerns**: Logic in hooks, UI in components

### 4. **State Management Strategy**
- **Server State**: Fetch in Server Components or use React Query
- **Client State**: Zustand stores for cross-component state
- **Form State**: React Hook Form
- **URL State**: Next.js searchParams for filters/pagination

### 5. **Internationalization**
- All user-facing text must use `useTranslations()` hook
- Never hardcode strings in Mongolian or English
- Add new translations to both `mn.json` and `en.json`

---

## 🔐 Authentication & Authorization

### Authentication Flow
1. User logs in → Server sets HTTP-only cookie
2. `UserBootstrap` runs on app mount → calls `loadProfile()`
3. User/Organization data stored in `useUserStore`
4. Roles fetched → `useRoleStore`
5. Systems fetched based on role → `useSystemStore`
6. Sidebar menu dynamically generated

### Authorization Levels
1. **System Level**: User can access specific systems (Admin, App, Business, TPay)
2. **Module Level**: Within a system, user has access to specific modules
3. **Permission Level**: Within modules, specific CRUD permissions

### Handling Unauthorized (401)
- `AuthBootstrap` sets up `setUnauthorizedHandler()`
- On 401 response → Clear stores → Redirect to `/home`
- `logout()` function: POST to `/auth/logout`, clear storage, redirect

---

## 🌐 API Client (`lib/api.ts`)

### Core Features
- **Centralized**: All API calls go through `api()` function
- **Type-Safe**: Generic type parameter `api.get<T>(...)`
- **Timeout**: Default 15 seconds, configurable
- **Auto Unwrap**: Extracts `data` field from response
- **Error Handling**: Automatic toast notifications
- **Dev Proxy**: Rewrites `/api/*` in development mode

### Usage Pattern

```typescript
// GET request
const users = await api.get<App.User[]>('/user', {
  query: { page: 1, size: 50 }
})

// POST request
const newUser = await api.post<App.User>('/user', {
  reg_no: 'XX12345678',
  phone_no: '99999999'
})

// PUT request
await api.put('/user', { id: 123, email: 'new@example.com' })

// DELETE request
await api.del(`/user/${userId}`)

// Disable toast on specific call
const data = await api.get('/some-endpoint', { hasToast: false })
```

### API Options
```typescript
{
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  query?: Record<string, string | number | boolean>
  headers?: HeadersInit
  body?: BodyInit | Record<string, unknown>
  baseURL?: string
  timeoutMs?: number
  cache?: RequestCache
  signal?: AbortSignal
  unwrapData?: boolean        // Default: true
  hasToast?: boolean | 'success' | 'error'  // Default: true
  onRequest?(url, init): void
  onResponse?(res): void
  onError?(err): void
}
```

---

## 📦 State Stores (Zustand)

### User Store (`lib/stores/user.ts`)

**Purpose**: Manage user/organization profile and authentication state

```typescript
interface UserState {
  user_info?: App.UserDetail           // Individual user data
  org_info?: App.Organization          // Organization data (if logged in as org)
  is_org: boolean                      // Flag: user vs organization
  user_name?: string                   // Display name
  profile_image?: string               // Avatar URL
  status: 'idle' | 'loading' | 'succeeded' | 'failed'
  error?: string
  
  loadProfile: () => Promise<void>     // Fetch profile from server
  clearAll: () => void                 // Reset store
}
```

**Persistence**: `user_info`, `user_name`, `profile_image` → localStorage

### Organization Store (`lib/stores/org.ts`)

**Purpose**: Manage user's organization list and selection

```typescript
interface OrgState {
  organizations: App.Organization[]    // List of orgs user belongs to
  selectedOrganization?: App.Organization
  
  getOrganizations: () => Promise<void>
  selectOrg: (org: App.Organization) => void  // Switch org + reload page
  clear: () => void
}
```

### Role Store (`lib/stores/role.ts`)

**Purpose**: Manage user roles within systems

```typescript
interface RoleState {
  roleList: App.UserRole[]             // User's available roles
  selectedRole?: App.UserRole          // Currently active role
  
  getRoleList: () => Promise<void>
  selectRole: (role: App.UserRole) => Promise<void>
  clear: () => void
}
```

### System Store (`lib/stores/system.ts`)

**Purpose**: Manage systems accessible by current role

```typescript
interface SystemState {
  systemList: App.System[]             // Systems for current role
  selectedSystem?: App.System          // Currently active system
  
  selectSystem: (sys?: App.System) => void
  changeSystemList: (sList: App.System[]) => void
  clear: () => void
}
```

**Note**: Each `System` contains `groups[]` → `ModuleGroup` contains `modules[]` → Used to build sidebar navigation

---

## 🧭 Navigation & Routing

### Route Structure

```
/[locale]/                           # mn or en
├── /home                            # Public home page
├── /admin/                          # Admin system
│   ├── /dashboard
│   ├── /user
│   ├── /role
│   ├── /organization
│   ├── /system
│   ├── /module
│   └── ...
├── /app/                            # App system
│   ├── /dashboard
│   ├── /icon
│   └── /role
├── /business/                       # Business system
│   ├── /agent
│   ├── /dashboard
│   └── ...
├── /tpay/                           # TPay system
│   ├── /balance
│   ├── /transaction
│   └── /wallet
├── /profile                         # Personal profile
├── /settings                        # User settings
└── /wallet                          # Personal wallet
```

### Navigation Utilities

```typescript
import { Link, useRouter, usePathname } from '@/i18n/navigation'

// These are locale-aware versions from next-intl
<Link href="/admin/user">Users</Link>

const router = useRouter()
router.push('/admin/dashboard')

const pathname = usePathname()  // Without locale prefix
```

---

## 🎨 UI Component Guidelines

### Using UI Components

All UI components in `components/ui/` follow these patterns:

```typescript
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card'

// Button variants
<Button variant="default">Primary Action</Button>
<Button variant="outline">Secondary</Button>
<Button variant="destructive">Delete</Button>
<Button variant="ghost">Subtle</Button>

// Button sizes
<Button size="sm">Small</Button>
<Button size="default">Default</Button>
<Button size="lg">Large</Button>
<Button size="icon"><Icon /></Button>
```

### Form Pattern

```typescript
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from '@/components/ui/form'

const schema = z.object({
  name: z.string().min(2, 'Name too short'),
  email: z.string().email('Invalid email'),
})

type FormData = z.infer<typeof schema>

const form = useForm<FormData>({
  resolver: zodResolver(schema),
  defaultValues: { name: '', email: '' }
})

const onSubmit = async (data: FormData) => {
  await api.post('/endpoint', data)
}

<Form {...form}>
  <form onSubmit={form.handleSubmit(onSubmit)}>
    <FormField
      control={form.control}
      name="name"
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t('name')}</FormLabel>
          <FormControl>
            <Input {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
    <Button type="submit">Submit</Button>
  </form>
</Form>
```

### Dialog Pattern

```typescript
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'

const [open, setOpen] = useState(false)

<Dialog open={open} onOpenChange={setOpen}>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Dialog Title</DialogTitle>
      <DialogDescription>Description text</DialogDescription>
    </DialogHeader>
    
    {/* Content */}
    
    <DialogFooter>
      <Button variant="outline" onClick={() => setOpen(false)}>Cancel</Button>
      <Button onClick={handleSubmit}>Confirm</Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

### Table Pattern

```typescript
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'

<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
      <TableHead>Email</TableHead>
      <TableHead>Actions</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    {items.map(item => (
      <TableRow key={item.id}>
        <TableCell>{item.name}</TableCell>
        <TableCell>{item.email}</TableCell>
        <TableCell>
          <Button size="sm" onClick={() => edit(item)}>Edit</Button>
        </TableCell>
      </TableRow>
    ))}
  </TableBody>
</Table>
```

---

## 📄 Page Component Pattern

### Standard CRUD Page Structure

```typescript
'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card'
// ... other imports

type Item = App.SomeType

const createSchema = z.object({
  name: z.string().min(1, 'Required'),
  // ... fields
})

const editSchema = z.object({
  // ... fields
})

export default function SomePage() {
  const t = useTranslations()
  
  // State
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(false)
  const [selected, setSelected] = useState<Item | null>(null)
  
  // Modals
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  
  // Pagination
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(50)
  const [totalPages, setTotalPages] = useState(1)
  
  // Forms
  const createForm = useForm({ resolver: zodResolver(createSchema) })
  const editForm = useForm({ resolver: zodResolver(editSchema) })
  
  // Load data
  const loadItems = async () => {
    setLoading(true)
    try {
      const data = await api.get<App.ListData<Item>>('/endpoint', {
        query: { page, size: pageSize }
      })
      setItems(data.items)
      setTotalPages(data.meta.pages)
    } catch (e) {
      // Error handled by api client
    } finally {
      setLoading(false)
    }
  }
  
  useEffect(() => {
    loadItems()
  }, [page, pageSize])
  
  // CRUD handlers
  const onCreate = async (data: any) => {
    await api.post('/endpoint', data)
    setOpenCreate(false)
    loadItems()
  }
  
  const onUpdate = async (data: any) => {
    await api.put('/endpoint', { id: selected?.id, ...data })
    setOpenEdit(false)
    loadItems()
  }
  
  const onDelete = async () => {
    await api.del(`/endpoint/${selected?.id}`)
    setOpenDelete(false)
    loadItems()
  }
  
  return (
    <div className="h-full w-full overflow-hidden p-4 sm:p-6">
      <Card className="flex h-full flex-col">
        <CardHeader>
          <CardTitle>{t('page_title')}</CardTitle>
          <Button onClick={() => setOpenCreate(true)}>
            {t('create')}
          </Button>
        </CardHeader>
        
        <CardContent>
          {/* Filters, search, table */}
        </CardContent>
        
        <CardFooter>
          {/* Pagination */}
        </CardFooter>
      </Card>
      
      {/* Create/Edit/Delete Dialogs */}
    </div>
  )
}
```

---

## 🌍 Internationalization (i18n)

### Adding New Translations

1. Add key to `i18n/mn.json`:
```json
{
  "my_new_key": "Монгол орчуулга"
}
```

2. Add same key to `i18n/en.json`:
```json
{
  "my_new_key": "English translation"
}
```

3. Use in component:
```typescript
const t = useTranslations()
<p>{t('my_new_key')}</p>
```

### Rich Text Formatting

```typescript
// For dynamic content in translations
{
  "delete_warning": "Та <name></name>-г устгахдаа итгэлтэй байна уу?"
}

// Usage
t.rich('delete_warning', {
  name: () => <span className="font-medium">{userName}</span>
})
```

### Pluralization

```json
{
  "systemCount": "{count} систем"
}

// Usage
t('systemCount', { count: 5 })  // "5 систем"
```

---

## 🎨 Styling Guidelines

### Tailwind Best Practices

1. **Use Utility Classes**: Prefer Tailwind utilities over custom CSS
2. **Responsive Design**: Mobile-first (`sm:`, `md:`, `lg:`)
3. **Dark Mode**: Always provide dark variants
4. **Consistent Spacing**: Use spacing scale (p-4, gap-2, etc.)

### Common Patterns

```typescript
// Card container with full height
<div className="h-full w-full overflow-hidden p-4 sm:p-6">
  <Card className="flex h-full flex-col">
    {/* Content */}
  </Card>
</div>

// Flex layout with gap
<div className="flex items-center gap-2">
  {/* Items */}
</div>

// Grid layout
<div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
  {/* Items */}
</div>

// Text truncation
<p className="truncate">{longText}</p>

// Conditional classes
<div className={cn(
  "base classes",
  isActive && "active classes",
  variant === 'primary' && "primary classes"
)}>
```

### Dark Mode Support

```typescript
// Always provide dark variants
<div className="bg-white dark:bg-gray-900">
<p className="text-gray-900 dark:text-white">
<div className="border-gray-200 dark:border-gray-700">
```

---

## 🔧 Development Workflow

### Getting Started

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint code
npm run lint

# Format code
npm run format
```

### Environment Variables

Create `.env.local`:

```env
# API Configuration
NEXT_PUBLIC_API_BASE=/api

# Development Proxy (optional)
API_PROXY_TARGET=https://template.gerege.mn/api

# Development Session ID (optional)
NEXT_PUBLIC_DEV_SID=your-dev-session-token
```

### Code Quality Checks

1. **TypeScript**: No errors allowed
```bash
npx tsc --noEmit
```

2. **ESLint**: Must pass
```bash
npm run lint
```

3. **Prettier**: Auto-format before commit
```bash
npm run format
```

---

## 📝 Code Conventions

### Naming Conventions

- **Files**: camelCase for components (`userList.tsx`)
- **Components**: PascalCase (`UserList`)
- **Functions**: camelCase (`loadUserData`)
- **Constants**: UPPER_SNAKE_CASE (`API_BASE_URL`)
- **Types/Interfaces**: PascalCase (`UserProfile`, `App.User`)

### File Organization

```typescript
// 1. Imports - external libraries first
import { useState } from 'react'
import { useForm } from 'react-hook-form'

// 2. Imports - internal components/utils
import { Button } from '@/components/ui/button'
import api from '@/lib/api'

// 3. Types/Interfaces
type Props = { ... }
type FormData = { ... }

// 4. Component
export default function MyComponent() {
  // 4.1. Hooks
  const [state, setState] = useState()
  const form = useForm()
  
  // 4.2. Effects
  useEffect(() => { ... }, [])
  
  // 4.3. Handlers
  const handleClick = () => { ... }
  
  // 4.4. Render
  return ( ... )
}

// 5. Helper functions (if not exported)
function helperFunction() { ... }
```

### Comment Guidelines

- **Mongolian comments**: Use for complex logic explanation
- **JSDoc**: For exported functions/components
- **Emoji markers**: 
  - 🧩 Types/Interfaces
  - 🧱 Main functions
  - ⚙️ Configuration
  - 🔹 Important sections
  - ❌ Error handling
  - ✅ Success cases

Example:
```typescript
/**
 * 🧩 UserState төрөл
 * Энэ store нь хэрэглэгчийн болон байгууллагын профайлын мэдээллийг удирдах зориулалттай.
 */
```

---

## 🚀 Deployment

### PM2 Configuration

File: `ecosystem.config.cjs`

```javascript
module.exports = {
  apps: [{
    name: 'template',
    script: 'npm',
    args: 'start',
    env: {
      NODE_ENV: 'production',
      PORT: 3000
    }
  }]
}
```

### Deployment Commands

```bash
# Standard deployment
npm run deploy
# → git pull && npm run build && pm2 restart template

# With dependency install
npm run ideploy
# → git pull && npm i && npm run build && pm2 restart template
```

### Build Optimization

- Uses Turbopack for faster builds
- Static optimization for public pages
- Dynamic rendering for authenticated pages
- Image optimization via next/image

---

## 🧪 Testing Strategy

### Recommended Testing Approach

1. **Unit Tests**: Utils, pure functions
2. **Component Tests**: UI components in isolation
3. **Integration Tests**: Page flows, API interactions
4. **E2E Tests**: Critical user journeys

### Testing Libraries (Future)
- **Vitest**: Unit/integration testing
- **React Testing Library**: Component testing
- **Playwright**: E2E testing

---

## 🐛 Error Handling

### API Error Handling

All API errors are automatically caught and displayed via toast:

```typescript
try {
  const data = await api.get('/endpoint')
} catch (e) {
  // Toast already shown by api client
  // Additional handling if needed
  console.error(e)
}
```

### Custom Error Handling

```typescript
// Disable automatic toast
try {
  const data = await api.get('/endpoint', { hasToast: false })
} catch (e) {
  // Custom error handling
  if (e instanceof APIError) {
    if (e.status === 404) {
      // Handle 404
    }
  }
}
```

### Form Validation Errors

Handled automatically by React Hook Form + Zod:

```typescript
const schema = z.object({
  email: z.string().email('Invalid email format')
})

// Error message automatically shown in <FormMessage />
```

---

## 📚 Key Patterns & Best Practices

### 1. **Loading States**

```typescript
const [loading, setLoading] = useState(false)
const [progress, setProgress] = useState(0)

// Show progress bar
{progress > 0 && (
  <Progress value={progress} className="h-1" />
)}

// Show skeleton
{loading ? <Skeleton /> : <Content />}

// Show spinner
{loading && <Loader2 className="animate-spin" />}
```

### 2. **Pagination**

```typescript
const [page, setPage] = useState(1)
const [pageSize, setPageSize] = useState(50)
const [totalPages, setTotalPages] = useState(1)
const [totalItems, setTotalItems] = useState(0)

// Load with pagination
const data = await api.get('/endpoint', {
  query: { page, size: pageSize }
})

setTotalPages(data.meta.pages)
setTotalItems(data.meta.total)
```

### 3. **Filtering & Search**

```typescript
const [searchTerm, setSearchTerm] = useState('')
const [filters, setFilters] = useState({})

// Debounced search recommended
useEffect(() => {
  const timer = setTimeout(() => {
    loadData(searchTerm)
  }, 300)
  return () => clearTimeout(timer)
}, [searchTerm])
```

### 4. **Optimistic UI Updates**

```typescript
const onUpdate = async (item: Item) => {
  // Optimistically update UI
  setItems(prev => prev.map(i => 
    i.id === item.id ? { ...i, ...item } : i
  ))
  
  try {
    await api.put('/endpoint', item)
  } catch (e) {
    // Revert on error
    loadItems()
  }
}
```

### 5. **Preventing Duplicate Requests**

```typescript
const lastReqId = useRef(0)

const loadData = async () => {
  const reqId = ++lastReqId.current
  setLoading(true)
  
  try {
    const data = await api.get('/endpoint')
    
    // Only update if this is the latest request
    if (reqId === lastReqId.current) {
      setData(data)
    }
  } finally {
    if (reqId === lastReqId.current) {
      setLoading(false)
    }
  }
}
```

---

## 🔒 Security Guidelines

### 1. **Never Expose Sensitive Data**
- No API keys in client-side code
- Use environment variables for secrets
- Session tokens in HTTP-only cookies

### 2. **Input Validation**
- Always validate on both client (Zod) and server
- Sanitize user inputs
- Use TypeScript for type safety

### 3. **Authentication**
- Cookie-based sessions
- CSRF protection via SameSite cookies
- Automatic 401 handling

### 4. **Authorization**
- Check permissions on server-side
- Never rely only on UI hiding
- Validate role/permissions per request

---

## 📊 Performance Guidelines

### 1. **Code Splitting**
- Use dynamic imports for heavy components
- Lazy load routes
- Server Components by default

### 2. **Image Optimization**
```typescript
import Image from 'next/image'

<Image 
  src="/path/to/image.jpg"
  width={500}
  height={300}
  alt="Description"
  priority={false}  // true for above-fold images
/>
```

### 3. **Memoization**
```typescript
import { useMemo, useCallback } from 'react'

const expensiveValue = useMemo(() => {
  return computeExpensiveValue(data)
}, [data])

const handleClick = useCallback(() => {
  // handler logic
}, [dependencies])
```

### 4. **API Caching**
```typescript
// Cache for 60 seconds
const data = await api.get('/endpoint', {
  cache: 'force-cache',
  next: { revalidate: 60 }
})
```

---

## 🎯 Project-Specific Rules

### 1. **Always Use Translation Hook**
```typescript
// ❌ Wrong
<Button>Create User</Button>
<Button>Хэрэглэгч үүсгэх</Button>

// ✅ Correct
const t = useTranslations()
<Button>{t('create', { name: t('user') })}</Button>
```

### 2. **Consistent Error Handling**
```typescript
// ❌ Wrong - showing custom alerts
alert('Error occurred!')

// ✅ Correct - let API client handle it
await api.post('/endpoint', data)
// Toast automatically shown
```

### 3. **Use Centralized API Client**
```typescript
// ❌ Wrong
const res = await fetch('/api/user')
const data = await res.json()

// ✅ Correct
const data = await api.get<App.User[]>('/user')
```

### 4. **Store Usage**
```typescript
// ❌ Wrong - direct store mutation
useUserStore.getState().user_info = newData

// ✅ Correct - use store actions
useUserStore.getState().loadProfile()
```

### 5. **Type Definitions**
```typescript
// ❌ Wrong - inline types
const [users, setUsers] = useState<{id: number, name: string}[]>([])

// ✅ Correct - use App namespace
const [users, setUsers] = useState<App.User[]>([])
```

---

## 📞 Support & Contact

**Developer**: Sengum Soronzonbold  
**Team**: Gerege Core Team  
**Project**: Gerege Template v25

For questions, issues, or contributions, contact the Gerege Core Team.

---

## 📜 Version History

- **v1.0.0** (2025-12-07): Initial master prompt creation

---

**End of Master Prompt**

Use this document as the source of truth for all development decisions, architectural patterns, and coding conventions in the Gerege Template project.

