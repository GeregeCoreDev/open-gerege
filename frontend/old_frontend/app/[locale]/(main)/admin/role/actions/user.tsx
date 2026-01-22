/**
 * 👥🛡️ Role Users Dialog Component
 * (/[locale]/(main)/admin/role/actions/user.tsx)
 *
 * Энэ нь дүрийн хэрэглэгчид удирдах dialog component юм.
 * Зорилго: Тухайн дүрт хэрэглэгч нэмэх, хасах
 *
 * Features:
 * - ✅ Two-level dialog system:
 *   1. Main dialog: View assigned users (WITH pagination)
 *   2. Add dialog: Search user from Core system
 * - ✅ Paginated user list
 * - ✅ Remove user from role
 * - ✅ Core system integration
 * - ✅ Search user by name/reg_no/phone/email
 * - ✅ User preview before adding
 * - ✅ Duplicate detection (already assigned)
 * - ✅ Loading states
 *
 * Props:
 * @param open - Dialog visibility
 * @param onOpenChange - Toggle dialog
 * @param role - Role object
 * @param onChanged - Callback after changes
 *
 * User Workflow:
 * 1. Open dialog → See users with this role (paginated)
 * 2. Click "Add" → Open search dialog
 * 3. Search user in Core system
 * 4. Preview found user
 * 5. Confirm and add
 *
 * Pagination:
 * - Page size: 20 (configurable 5/10/20/50/100)
 * - Server-side pagination
 * - Smart page navigation (ellipsis)
 * - Auto adjust page when removing last item
 *
 * API Endpoints:
 * - GET /role-matrix/users?role_id=... - Get users with role
 * - POST /user/find-from-core - Search user
 * - POST /role-matrix { role_id, user_ids: [] } - Assign user
 * - DELETE /role-matrix { role_id, user_id } - Remove user
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import * as React from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import { Loader2, Plus, Trash2, Search, X } from 'lucide-react'
import api from '@/lib/api'
import { appConfig } from '@/config/app.config'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationLink,
  PaginationEllipsis,
} from '@/components/ui/pagination'
import { useTranslations } from 'next-intl'

/**
 * 🧱 Mono font cell helper
 */
function CellMono({ children }: { children?: React.ReactNode }) {
  return <span className="font-mono text-xs">{children ?? '-'}</span>
}

/**
 * 🧩 Component Props
 */
type Props = {
  open: boolean
  onOpenChange: (v: boolean) => void
  role: App.Role | null
  onChanged?: () => void
}

export default function RoleUsersDialog({ open, onOpenChange, role, onChanged }: Props) {
  const t = useTranslations()

  // ========================================
  // 📋 Main Dialog: Assigned Users (WITH pagination)
  // ========================================

  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [rows, setRows] = React.useState<App.RoleUser[]>([])

  // 📄 Pagination state
  const [pageNumber, setPageNumber] = React.useState(1)
  const [pageSize, setPageSize] = React.useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = React.useState(1)
  const [totalRow, setTotalRow] = React.useState(0)
  const [meta, setMeta] = React.useState<App.ApiMeta | null>(null)

  // remove state
  const [removingId, setRemovingId] = React.useState<number | null>(null)

  const pickItems = <T,>(res: unknown): T[] => {
    const r = res as { items?: T[]; data?: T[] } | T[]
    if (Array.isArray(r)) return r
    if (Array.isArray(r?.items)) return r.items
    if (Array.isArray(r?.data)) return r.data
    return []
  }

  const loadUsers = React.useCallback(
    async (page = pageNumber, size = pageSize) => {
      if (!role) return
      setLoading(true)
      setError(null)
      try {
        const res = await api.get<App.ListData<App.RoleUser>>('/role-matrix/users', {
          query: { role_id: role.id, page, size },
          cache: 'no-store',
        })
        const items = pickItems<App.RoleUser>(res)
        const m = (res?.meta ?? null) as App.ApiMeta | null

        setRows(items)
        setMeta(m)
        setPageNumber(m?.page ?? page)
        setPageSize(m?.size ?? size)
        setTotalPage(m?.pages ?? 1)
        setTotalRow(m?.total ?? items.length ?? 0)
      } catch {
        setError("Error occurred")
      } finally {
        setLoading(false)
      }
    },
    [role, pageNumber, pageSize],
  )

  // open болгонд эхний хуудсыг татна
  React.useEffect(() => {
    if (open) {
      loadUsers(1, pageSize)
    } else {
      setRows([])
      setError(null)
      setPageNumber(1)
      setTotalPage(1)
      setTotalRow(0)
      setMeta(null)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, role?.id])

  // pageSize солигдоход 1-р хуудас руу
  React.useEffect(() => {
    if (open) loadUsers(1, pageSize)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize])

  const goPage = (p: number) => {
    const target = Math.min(Math.max(1, p), Math.max(1, totalPage))
    if (target !== pageNumber) loadUsers(target, pageSize)
  }

  const pageLinks = React.useMemo(() => {
    const tp = totalPage
    const cur = pageNumber
    const links: (number | '…')[] = []
    if (tp <= 7) {
      for (let i = 1; i <= tp; i++) links.push(i)
      return links
    }
    const w = 2
    links.push(1)
    if (cur > 1 + w + 1) links.push('…')
    const start = Math.max(2, cur - w)
    const end = Math.min(tp - 1, cur + w)
    for (let i = start; i <= end; i++) links.push(i)
    if (cur < tp - w - 1) links.push('…')
    links.push(tp)
    return links
  }, [pageNumber, totalPage])

  /**
   * 🗑️ Хэрэглэгчийг дүрээс хасах
   * Smart pagination: Хуудас хоосорвол өмнөх хуудас руу буцна
   */
  const removeOne = async (userId: number) => {
    if (!role) return
    try {
      setRemovingId(userId)
      await api.del('/role-matrix', { role_id: role.id, user_id: userId } as Record<string, unknown>)
      // энэ хуудсан дээр нэг мөрөөр багасна → хоосорвол өмнөх хуудсыг ачаал
      const nextCount = rows.length - 1
      const willBeEmpty = nextCount <= 0 && pageNumber > 1
      await loadUsers(willBeEmpty ? pageNumber - 1 : pageNumber, pageSize)
      onChanged?.()
    } catch {
      setError("Error occurred")
    } finally {
      setRemovingId(null)
    }
  }

  // ========================================
  // 🔍 Add User Dialog: CORE Search
  // ========================================

  const [openAdd, setOpenAdd] = React.useState(false)
  const [q, setQ] = React.useState('')
  const [searching, setSearching] = React.useState(false)
  const [searchErr, setSearchErr] = React.useState<string | null>(null)
  const [found, setFound] = React.useState<App.User | null>(null)
  const [adding, setAdding] = React.useState(false)

  const assignedIdSet = React.useMemo(() => new Set(rows.map((u) => u.user_id)), [rows])

  /**
   * 🔍 Core system-с хэрэглэгч хайх
   */
  const findFromCore = async () => {
    const term = q.trim()
    if (!term) return
    setSearching(true)
    setSearchErr(null)
    setFound(null)
    try {
      const res = await api.post<App.User | null>('/user/find-from-core', { search_text: term })
      if (res && typeof res === 'object') setFound(res)
      else {
        setFound(null)
        setSearchErr('Илэрц олдсонгүй.')
      }
    } catch {
      setSearchErr("Error occurred")
    } finally {
      setSearching(false)
    }
  }

  /**
   * ➕ Олдсон хэрэглэгчийг дүрт нэмэх
   */
  const addOne = async () => {
    if (!role || !found) return
    if (assignedIdSet.has(found.id)) return
    setAdding(true)
    try {
      await api.post('/role-matrix', { role_id: role.id, user_ids: [found.id] } as Record<string, unknown>)
      // нэмсэний дараа одоогийн хуудсаа рефреш
      await loadUsers(pageNumber, pageSize)
      onChanged?.()
      setOpenAdd(false)
      setQ('')
      setFound(null)
      setSearchErr(null)
    } catch {
      setSearchErr("Error occurred")
    } finally {
      setAdding(false)
    }
  }

  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? (pageNumber - 1) * pageSize) + 1
  const showingTo =
    totalRow === 0 ? 0 : (meta?.end_idx ?? Math.min(pageNumber * pageSize - 1, totalRow - 1)) + 1
  const canPrev = pageNumber > 1
  const canNext = pageNumber < totalPage

  /** ===== Render ===== */
  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-4xl">
          <DialogHeader>
            <DialogTitle>{role?.name ?? 'Role'} — Хэрэглэгчид</DialogTitle>
            <DialogDescription>Энэ role-д оноогдсон хэрэглэгчид.</DialogDescription>
          </DialogHeader>

          <div className="mb-2 flex items-center justify-between gap-3">
            <Button size="sm" onClick={() => setOpenAdd(true)}>
              <Plus className="mr-1 h-4 w-4" /> Нэмэх
            </Button>

            <div className="flex items-center gap-x-2">
              <p className="text-muted-foreground text-sm">{t('rows')}</p>
              <Select value={String(pageSize)} onValueChange={(v) => setPageSize(Number(v))}>
                <SelectTrigger className="h-9 w-[84px]">
                  <SelectValue placeholder={pageSize} />
                </SelectTrigger>
                <SelectContent>
                  {[5, 10, 20, 50, 100].map((s) => (
                    <SelectItem key={s} value={String(s)}>
                      {s}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-16"></TableHead>
                  <TableHead className="w-24">ID</TableHead>
                  <TableHead>Нэр</TableHead>
                  <TableHead>Регистер</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Утас</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {error ? (
                  <TableRow>
                    <TableCell colSpan={6} className="py-8 text-center text-sm text-red-600">
                      {error}
                    </TableCell>
                  </TableRow>
                ) : loading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="py-8 text-center text-sm opacity-70">
                      <Loader2 className="mr-2 inline h-4 w-4 animate-spin" />
                      Уншиж байна...
                    </TableCell>
                  </TableRow>
                ) : rows.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="py-8 text-center text-sm opacity-70">
                      Хоосон байна.
                    </TableCell>
                  </TableRow>
                ) : (
                  rows.map((u: App.RoleUser) => {
                    const name =
                      `${u.user.last_name ?? ''} ${u.user.first_name ?? ''}`
                        .replace(/\s+/g, ' ')
                        .trim() ||
                      u.user.first_name ||
                      '-'
                    return (
                      <TableRow key={u.user.id}>
                        <TableCell className="text-right">
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => removeOne(u.user.id)}
                            disabled={removingId === u.user.id}
                            aria-label="Хасах"
                          >
                            {removingId === u.user.id ? (
                              <Loader2 className="h-4 w-4 animate-spin" />
                            ) : (
                              <Trash2 className="h-4 w-4" />
                            )}
                          </Button>
                        </TableCell>
                        <TableCell className="text-muted-foreground text-xs">
                          {String(u.user.id)}
                        </TableCell>
                        <TableCell className="font-medium capitalize">{name}</TableCell>
                        <TableCell>
                          <CellMono>{u.user.reg_no ?? '-'}</CellMono>
                        </TableCell>
                        <TableCell className="text-muted-foreground text-xs">
                          {u.user.email ?? '-'}
                        </TableCell>
                        <TableCell className="text-muted-foreground text-xs">
                          {u.user.phone_no ?? '-'}
                        </TableCell>
                      </TableRow>
                    )
                  })
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination footer */}
          <div className="mt-3 flex items-center justify-between gap-3">
            <div className="text-muted-foreground w-auto min-w-72 text-sm">
              {t.rich('showing', {
                from: () => <span className="font-medium">{showingFrom || 0}</span>,
                to: () => <span className="font-medium">{showingTo || 0}</span>,
                total: () => <span className="font-medium">{totalRow || 0}</span>,
              })}
            </div>

            <Pagination className="justify-end">
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    aria-disabled={!canPrev}
                    className={!canPrev ? 'pointer-events-none opacity-50' : ''}
                    onClick={() => canPrev && goPage(pageNumber - 1)}
                  />
                </PaginationItem>

                {pageLinks.map((p, idx) =>
                  p === '…' ? (
                    <PaginationItem key={`e-${idx}`}>
                      <PaginationEllipsis />
                    </PaginationItem>
                  ) : (
                    <PaginationItem key={p}>
                      <PaginationLink isActive={p === pageNumber} onClick={() => goPage(p)}>
                        {p}
                      </PaginationLink>
                    </PaginationItem>
                  ),
                )}

                <PaginationItem>
                  <PaginationNext
                    aria-disabled={!canNext}
                    className={!canNext ? 'pointer-events-none opacity-50' : ''}
                    onClick={() => canNext && goPage(pageNumber + 1)}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>

          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Хаах
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ===== Nested: CORE хайлт (но pagination, но авто хайлт) ===== */}
      <Dialog
        open={openAdd}
        onOpenChange={(v) => {
          setOpenAdd(v)
          if (!v) {
            setQ('')
            setFound(null)
            setSearchErr(null)
            setSearching(false)
            setAdding(false)
          }
        }}
      >
        <DialogContent className="z-100 sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>CORE — Хэрэглэгч хайх</DialogTitle>
            <DialogDescription>Бичээд Enter/Хайх дарна. Олдсон бол дор харуулна.</DialogDescription>
          </DialogHeader>

          {/* Хайлтын мөр */}
          <div className="flex items-center gap-2">
            <div className="relative w-full max-w-xl">
              <Search className="text-muted-foreground pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2" />
              <Input
                value={q}
                onChange={(e) => setQ(e.target.value)}
                placeholder="Нэр / регистр / утас / имэйл"
                className="focus-visible:ring-primary/30 h-10 rounded-xl border-gray-200 pr-30 pl-10 shadow-sm focus-visible:ring-2 dark:border-gray-800"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    findFromCore()
                  }
                }}
              />
              {q && !searching ? (
                <button
                  type="button"
                  onClick={() => setQ('')}
                  aria-label="Цэвэрлэх"
                  className="hover:bg-muted absolute top-1/2 right-22 -translate-y-1/2 rounded-full p-1.5 transition"
                >
                  <X className="text-muted-foreground h-4 w-4" />
                </button>
              ) : null}
              <Button
                size="sm"
                className="absolute top-1/2 right-1 h-8 -translate-y-1/2 rounded-lg px-3 shadow-sm disabled:opacity-60"
                onClick={findFromCore}
                disabled={searching || !q.trim()}
              >
                {searching ? (
                  <>
                    <Loader2 className="mr-1 h-4 w-4 animate-spin" />
                    Хайж байна
                  </>
                ) : (
                  <>
                    <Search className="mr-1 h-4 w-4" />
                    Хайх
                  </>
                )}
              </Button>
            </div>
          </div>

          {searchErr && <div className="text-sm text-red-600">{searchErr}</div>}

          {found && (
            <div className="mt-3 rounded-md border p-3">
              <div className="flex items-start gap-3">
                {found.profile_img_url ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={found.profile_img_url}
                    alt="profile"
                    className="h-16 w-16 rounded-full object-cover"
                  />
                ) : (
                  <div className="bg-muted h-16 w-16 rounded-full" />
                )}
                <div className="grid flex-1 grid-cols-1 gap-1 sm:grid-cols-2">
                  <Field label="ID">
                    <CellMono>{String(found.id)}</CellMono>
                  </Field>
                  <Field label="Регистер">
                    <CellMono>{found?.reg_no ?? '-'}</CellMono>
                  </Field>
                  <Field label="Нэр">
                    {`${found?.last_name ?? ''} ${found?.first_name ?? ''}`
                      .replace(/\s+/g, ' ')
                      .trim() ||
                      found.first_name ||
                      '-'}
                  </Field>
                  <Field label="Утас">{found?.phone_no ?? '-'}</Field>
                  <Field label="И-мэйл">{found?.email ?? '-'}</Field>
                  <Field label="Төрсөн огноо">{found?.birth_date ?? '-'}</Field>
                </div>
              </div>
              <div className="mt-3 flex items-center justify-end gap-2">
                {assignedIdSet.has(found.id) && (
                  <span className="text-sm opacity-70">Энэ role-д аль хэдийн нэмэгдсэн</span>
                )}
                <Button onClick={addOne} disabled={adding || assignedIdSet.has(found.id)}>
                  {adding && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Нэмэх
                </Button>
              </div>
            </div>
          )}

          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setOpenAdd(false)}>
              Хаах
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="min-w-0">
      <div className="text-xs tracking-wide uppercase opacity-60">{label}</div>
      <div className="truncate text-sm">{children}</div>
    </div>
  )
}
