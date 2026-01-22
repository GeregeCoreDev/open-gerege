// 'use client'

// import { FloatingMic } from '@/components/common/mic'
// import { profileMenuList } from '@/components/layout/profileSidebar'
// import { useRouter } from '@/i18n/navigation'
// import { useSystemStore } from '@/lib/stores/system'
// import { useCallback, useMemo } from 'react'
// import { toast } from 'sonner'

// type CommandRoute = {
//   keywords: string[]
//   path: string
// }

// const STATIC_ROUTES: CommandRoute[] = [{ keywords: ['home', 'нүүр', 'эхлэл'], path: '/home' }]

// export default function LayoutMic() {
//   const router = useRouter()
//   const selectedSystem = useSystemStore((s) => s.selectedSystem)

//   const modules = useMemo(
//     () => selectedSystem?.groups?.flatMap((g) => g.modules || []) ?? [],
//     [selectedSystem],
//   )

//   const normalizeText = (raw: string) =>
//     raw
//       .toLowerCase()
//       .replace(/[.,!?'"`]/g, ' ')
//       .replace(/\s+/g, ' ')
//       .trim()

//   const findStaticRoute = useCallback((text: string) => {
//     return STATIC_ROUTES.find((route) => route.keywords.some((k) => text.includes(k)))?.path
//   }, [])

//   const findModuleRoute = useCallback(
//     (text: string) => {
//       for (const mod of modules) {
//         const name = mod?.name?.toLowerCase?.() ?? ''
//         if (!name) continue
//         if (text.includes(name)) return mod.path || '/'
//         const tokens = name.split(/\s+/).filter(Boolean)
//         if (tokens.some((t) => text.includes(t))) return mod.path || '/'
//       }

//       for (const mod of profileMenuList) {
//         const name = mod?.name?.toLowerCase?.() ?? ''
//         if (!name) continue
//         if (text.includes(name)) return mod.path || '/'
//         const tokens = name.split(/\s+/).filter(Boolean)
//         if (tokens.some((t) => text.includes(t))) return mod.path || '/'
//       }
//       return null
//     },
//     [modules],
//   )

//   const resolveRoute = useCallback(
//     (rawText: string) => {
//       const text = normalizeText(rawText)
//       if (!text) return null
//       const staticMatch = findStaticRoute(text)
//       if (staticMatch) return staticMatch
//       const moduleMatch = findModuleRoute(text)
//       if (moduleMatch) return moduleMatch.startsWith('/') ? moduleMatch : `/${moduleMatch}`
//       return null
//     },
//     [findModuleRoute, findStaticRoute],
//   )

//   const sendMessage = useCallback(
//     (spoken: string) => {
//       const target = resolveRoute(spoken)
//       if (!target) {
//         toast.error('Тушаал ойлгосонгүй. Дахин хэлээд үзээрэй.')
//         return
//       }
//       toast.success('Замчилж байна...')
//       router.push(target)
//     },
//     [resolveRoute, router],
//   )

//   return <FloatingMic onSend={sendMessage} />
// }
