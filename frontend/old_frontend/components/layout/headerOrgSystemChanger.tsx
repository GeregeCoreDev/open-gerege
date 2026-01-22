// // components/HeaderOrgSystem.tsx
// 'use client'

// import { memo, useMemo } from 'react'
// import { Building2, ChevronDown, Grid } from 'lucide-react'
// import { Button } from '@/components/ui/button'
// import {
//   DropdownMenu,
//   DropdownMenuTrigger,
//   DropdownMenuContent,
//   DropdownMenuLabel,
//   DropdownMenuSeparator,
//   DropdownMenuGroup,
//   DropdownMenuRadioGroup,
//   DropdownMenuRadioItem,
// } from '@/components/ui/dropdown-menu'
// import { useOrgStore } from '@/lib/stores/org'
// import { useSystemStore } from '@/lib/stores/system'
// import { useMenuStore } from '@/lib/stores/menu'
// import { cn } from '@/lib/utils'
// import { useLocale, useTranslations } from 'next-intl'
// import { useRouter } from 'next/navigation'

// function HeaderOrgSystemInner() {
//   const t = useTranslations()
//   const router = useRouter()
//   const locale = useLocale()

//   const { organizations, selectedOrganization, selectOrg } = useOrgStore()
//   const { systemList, selectedSystem, selectSystem } = useSystemStore()
//   const { getMenuList, getFirstChildPath, selectRoot } = useMenuStore()

//   const orgBtnText = selectedOrganization?.name ?? 'Organization'

//   const orgId = selectedOrganization?.id ? String(selectedOrganization.id) : undefined
//   const sysId = selectedSystem?.id ? String(selectedSystem.id) : undefined

//   const orgMap = useMemo(
//     () => new Map(organizations.map((o) => [String(o.id), o])),
//     [organizations],
//   )
//   const sysMap = useMemo(() => new Map(systemList.map((s) => [String(s.id), s])), [systemList])

//   const toLocaleHref = (rawPath?: string) => {
//     const p = (rawPath || '').startsWith('/') ? rawPath : `/${rawPath || ''}`
//     return `/${locale}${p}`.replace(/\/{2,}/g, '/')
//   }

//   async function onSelectSystem(systemId: string) {
//     const s = sysMap.get(systemId)
//     if (!s) return

//     selectSystem(s)

//     // Menu жагсаалт авна
//     const menus = await getMenuList()
//     // Энэ системтэй холбоотой эхний root menu-г олох
//     const systemRoot = menus.find((m) => m.system_id === s.id)
//     if (systemRoot) {
//       selectRoot(systemRoot.id)
//       // Эхний path-тай menu-г олох
//       const firstPath = getFirstChildPath(systemRoot.id)
//       if (firstPath) {
//         const href = toLocaleHref(firstPath)
//         router.push(href)
//         router.refresh()
//       } else {
//         router.push(`/${locale}/profile`)
//         router.refresh()
//       }
//     } else {
//       router.push(`/${locale}/profile`)
//       router.refresh()
//     }
//   }

//   return (
//     <div className="flex flex-wrap items-center gap-2">
//       <DropdownMenu>
//         <DropdownMenuTrigger asChild>
//           <Button variant="outline" className={cn('h-9 w-56 justify-between')}>
//             <span className="flex min-w-0 items-center gap-2">
//               <Building2 className="h-4 w-4 shrink-0" />
//               <span className="truncate">{orgBtnText}</span>
//             </span>
//             <ChevronDown className="h-4 w-4 shrink-0 opacity-70" />
//           </Button>
//         </DropdownMenuTrigger>
//         <DropdownMenuContent className="w-72" align="start">
//           <DropdownMenuLabel className="text-xs uppercase opacity-70">
//             Organization
//           </DropdownMenuLabel>
//           <DropdownMenuSeparator />
//           <DropdownMenuRadioGroup
//             value={orgId}
//             onValueChange={(v) => {
//               const o = orgMap.get(v)
//               if (o) selectOrg(o)
//             }}
//           >
//             <div className="max-h-[320px] overflow-auto pr-1">
//               {organizations.map((o) => (
//                 <DropdownMenuRadioItem key={o.id} value={String(o.id)} className="py-2">
//                   <div className="flex min-w-0 flex-col">
//                     <span className="truncate font-medium">{o.name}</span>
//                   </div>
//                 </DropdownMenuRadioItem>
//               ))}
//             </div>
//           </DropdownMenuRadioGroup>
//         </DropdownMenuContent>
//       </DropdownMenu>

//       <DropdownMenu>
//         <DropdownMenuTrigger asChild>
//           <Button variant="outline" className={cn('h-9 w-64 justify-between')}>
//             <span className="flex min-w-0 items-center gap-2">
//               <Grid className="h-4 w-4 shrink-0" />
//               {selectedSystem ? (
//                 <span className="truncate">
//                   {selectedSystem.name}{' '}
//                   <span className="text-gray-700 dark:text-gray-200">({selectedSystem.code})</span>
//                 </span>
//               ) : (
//                 <span className="truncate">{t('system')}</span>
//               )}
//             </span>
//             <ChevronDown className="h-4 w-4 shrink-0 opacity-70" />
//           </Button>
//         </DropdownMenuTrigger>
//         <DropdownMenuContent className="w-72" align="start">
//           <DropdownMenuLabel className="text-xs uppercase opacity-70">System</DropdownMenuLabel>
//           <DropdownMenuSeparator />
//           <DropdownMenuGroup>
//             <DropdownMenuRadioGroup value={sysId} onValueChange={(v) => onSelectSystem(v)}>
//               <div className="max-h-[320px] overflow-auto pr-1">
//                 {systemList.map((s) => (
//                   <DropdownMenuRadioItem key={s.id} value={String(s.id)} className="py-2">
//                     <div className="flex min-w-0 flex-col">
//                       <span className="truncate font-medium">{s.name}</span>
//                       {s.code ? (
//                         <span className="text-muted-foreground truncate text-xs">{s.code}</span>
//                       ) : null}
//                     </div>
//                   </DropdownMenuRadioItem>
//                 ))}
//               </div>
//             </DropdownMenuRadioGroup>
//           </DropdownMenuGroup>
//         </DropdownMenuContent>
//       </DropdownMenu>
//     </div>
//   )
// }

// const HeaderOrgSystem = memo(HeaderOrgSystemInner)
// export default HeaderOrgSystem
