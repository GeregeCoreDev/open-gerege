/**
 * 👥🏢 Organization Users Page (/[locale]/(main)/admin/organization/users/page.tsx)
 * 
 * Энэ нь байгууллагын хэрэглэгчдийн хуудас юм.
 * Зорилго: Тухайн байгууллагын хэрэглэгчдийн жагсаалт, удирдлага
 * 
 * Features:
 * - ✅ List organization's users
 * - ✅ Pagination with server-side meta
 * - ✅ Search by user  name
 * - ✅ Add user to organization
 * - ✅ Remove user from organization
 * - ✅ View user roles
 * - ✅ Progress bar loading
 * - ✅ User search/find from system
 * - ✅ Form validation
 * 
 * Table Columns:
 * - Actions (Roles/Remove buttons)
 * - ID
 * - Registration Number
 * - Name (family + last + first)
 * - Phone
 * - Email
 * - Created Date
 * 
 * Add User Workflow:
 * 1. Search user by reg_no
 * 2. Find from Core system
 * 3. Add to organization
 * 
 * Related Components:
 * - UserRolesDialog: View/manage user roles
 * 
 * API Endpoints:
 * - GET /organization/:id/users - List org users
 * - POST /user/find-from-core - Search user
 * - POST /organization/user - Add user to org
 * - DELETE /organization/user - Remove user from org
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import * as React from 'react'
import { useTranslations } from 'next-intl'
import { Loader2, Plus, ShieldCheck, Trash2, Search } from 'lucide-react'
import api from '@/lib/api'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'

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
import { Separator } from '@/components/ui/separator'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationLink,
  PaginationEllipsis,
} from '@/components/ui/pagination'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Progress } from '@/components/ui/progress'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Root as VisuallyHidden } from '@radix-ui/react-visually-hidden'
import { cn } from '@/lib/utils'
import UserFind from '@/components/common/userFind'

/** ===== Types ===== */
type ListData<T> = {
  items: T[]
  meta: {
    page: number
    size: number
    pages: number
    total: number
    has_prev?: boolean
    has_next?: boolean
    start_idx?: number
    end_idx?: number
  }
}
type UserOrg = App.UserOrganization

/** ===== API helpers ===== */
async function getOrgUsers(params: {
  org_id: number | string | null
  page: number
  size: number
  name?: string
}) {
  return api.get<ListData<UserOrg>>('/orguser/users', {
    query: {
      page: params.page,
      size: params.size,
      name: params.name || undefined,
      org_id: params.org_id,
    },
    cache: 'no-store',
  })
}
async function addOrgUser(org_id: number | string, user_id: number) {
  return api.post('/orguser', { org_id, user_id } as Record<string, unknown>)
}
async function removeOrgUser(org_id: number | string, user_id: number) {
  return api.del('/orguser', { org_id, user_id } as Record<string, unknown>)
}
async function listRoles() {
  return api.get<App.ListData<App.Role> | App.Role[]>('/role', { query: { page: 1, size: 500 } })
}
function extractRoleId(x: unknown): number | undefined {
  if (typeof x === 'number') return x
  const obj = x as { role_id?: number; id?: number; name?: string; user_id?: number; role?: { id?: number } }
  if (obj?.role_id != null) return Number(obj.role_id)
  if (obj?.id != null && obj?.name && !obj?.user_id) return Number(obj.id)
  if (obj?.role?.id != null) return Number(obj.role.id)
  return undefined
}
async function getUserRoles(user_id: number, org_id: number | string) {
  type RoleResponse = unknown[] | { items?: unknown[]; data?: unknown[] }
  const res = await api.get<RoleResponse>('/role-matrix/roles', { query: { user_id, org_id } })
  const raw = Array.isArray(res) ? res : ((res as { items?: unknown[]; data?: unknown[] })?.items ?? (res as { items?: unknown[]; data?: unknown[] })?.data ?? [])
  const ids = (Array.isArray(raw) ? raw : [])
    .map(extractRoleId)
    .filter((n): n is number => Number.isFinite(n))
  return ids
}
async function saveUserRoles(user_id: number, org_id: number | string, role_ids: number[]) {
  return api.post('/role-matrix', { user_id, org_id, role_ids } as Record<string, unknown>)
}

/** =======================================================================
 * Core body
 * =======================================================================*/
function OrgUsersBody({
  orgIdProp,
  headerRight,
  showTitle = true,
  titleText,
}: {
  orgIdProp?: number | string
  headerRight?: React.ReactNode
  showTitle?: boolean
  titleText?: string
}) {
  const t = useTranslations()

  /** ------- org_id ------- */
  const [orgId, setOrgId] = React.useState<number | string | null>(orgIdProp ?? null)
  const [orgErr, setOrgErr] = React.useState<string | null>(null)
  React.useEffect(() => {
    let mounted = true
    if (orgIdProp != null) {
      setOrgId(orgIdProp)
      return
    }
    ;(async () => {
      try {
        setOrgErr(null)
        type SessionResponse = { org_id?: number; orgId?: number; organization_id?: number }
        const session = await api.get<SessionResponse>('/auth/session', { cache: 'no-store' })
        const id = session?.org_id ?? session?.orgId ?? session?.organization_id
        if (mounted) {
          if (id == null) setOrgErr('org_id not found in session')
          setOrgId(id ?? null)
        }
      } catch {
        if (mounted) setOrgErr("Error occurred")
      }
    })()
    return () => {
      mounted = false
    }
  }, [orgIdProp])

  /** ------- list state ------- */
  const [rows, setRows] = React.useState<UserOrg[]>([])
  const [loading, setLoading] = React.useState(false)
  const [fetchError, setFetchError] = React.useState<string | null>(null)

    const [filterName, setFilterName] = React.useState('')
    const [pageNumber, setPageNumber] = React.useState(1)
    const [pageSize, setPageSize] = React.useState<number>(appConfig.pagination.defaultPageSize)
    const [totalPage, setTotalPage] = React.useState(1)
    const [totalRow, setTotalRow] = React.useState(0)
    const [meta, setMeta] = React.useState<ListData<UserOrg>['meta'] | null>(null)

    // progress bar - simplified with hook
    const progress = useLoadingProgress(loading)

    /** ------- load ------- */
  const load = async (page = pageNumber, size = pageSize, name = filterName) => {
    setLoading(true)
    setFetchError(null)
    try {
      const data = await getOrgUsers({ org_id: orgId, page, size, name })
      setRows(data.items ?? [])
      setMeta(data.meta ?? null)
      setPageNumber(data.meta?.page ?? page)
      setPageSize(data.meta?.size ?? size)
      setTotalPage(data.meta?.pages ?? 1)
      setTotalRow(data.meta?.total ?? 0)
    } catch {
      setFetchError("Error occurred")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    load(1, pageSize, filterName)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize])

  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) load(target, pageSize, filterName)
  }
  const changePageSize = (s: number) => {
    const size = Math.max(5, Math.min(200, s || 50))
    if (size !== pageSize) load(1, size, filterName)
  }

  const pageLinks = React.useMemo(() => {
    const links: (number | 'ellipsis')[] = []
    if (totalPage <= 7) {
      for (let i = 1; i <= totalPage; i++) links.push(i)
      return links
    }
    const w = 2
    links.push(1)
    if (pageNumber > 1 + w + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - w)
    const end = Math.min(totalPage - 1, pageNumber + w)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < totalPage - w - 1) links.push('ellipsis')
    links.push(totalPage)
    return links
  }, [pageNumber, totalPage])

  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? (pageNumber - 1) * pageSize) + 1
  const showingTo =
    totalRow === 0 ? 0 : (meta?.end_idx ?? (pageNumber - 1) * pageSize + rows.length - 1) + 1
  const canPrev = meta?.has_prev ?? pageNumber > 1
  const canNext = meta?.has_next ?? pageNumber < totalPage

  /** ------- row actions ------- */
  const [openDelete, setOpenDelete] = React.useState(false)
  const [selectedRow, setSelectedRow] = React.useState<UserOrg | null>(null)
  const onOpenDelete = (row: UserOrg) => {
    setSelectedRow(row)
    setOpenDelete(true)
  }
  const onDelete = async () => {
    if (!selectedRow || orgId == null) return
    await removeOrgUser(orgId, selectedRow.user_id)
    setOpenDelete(false)
    setSelectedRow(null)
    load(pageNumber, pageSize, filterName)
  }

  /** ======================= Add user (SINGLE via <UserFind>) ======================= */
  const [openAdd, setOpenAdd] = React.useState(false)
  const [selectedUser, setSelectedUser] = React.useState<App.User | null>(null)
  const alreadyInOrg = React.useMemo(
    () => (selectedUser ? rows.some((r) => r.user_id === selectedUser.id) : false),
    [selectedUser, rows],
  )

  const addUserToOrg = async () => {
    if (!selectedUser || orgId == null) return
    await addOrgUser(orgId, selectedUser.id)
    await load(pageNumber, pageSize, filterName)
    setOpenAdd(false)
    setSelectedUser(null)
  }

  /** ======================= Role change modal ======================= */
  const [roleModalOpen, setRoleModalOpen] = React.useState(false)
  const [allRoles, setAllRoles] = React.useState<App.Role[]>([])
  const [roleLoading, setRoleLoading] = React.useState(false)
  const [roleError, setRoleError] = React.useState<string | null>(null)
  const [assigned, setAssigned] = React.useState<Set<number>>(new Set())

  const onOpenRole = async (row: UserOrg) => {
    setSelectedRow(row)
    setRoleModalOpen(true)
    setRoleError(null)
    setRoleLoading(true)
    try {
      const r = await listRoles()
      const list = Array.isArray(r) ? r : (r.items ?? [])
      setAllRoles(list)

      const a = await getUserRoles(row.user_id, orgId ?? 0)
      setAssigned(new Set(a))
    } catch {
      setRoleError("Error occurred")
    } finally {
      setRoleLoading(false)
    }
  }

  const toggleRole = (id: number, on?: boolean) => {
    setAssigned((prev) => {
      const n = new Set(prev)
      if (on === undefined) {
        if (n.has(id)) {
          n.delete(id)
        } else {
          n.add(id)
        }
      } else if (on) {
        n.add(id)
      } else {
        n.delete(id)
      }
      return n
    })
  }

  const saveRoles = async () => {
    if (!selectedRow || orgId == null) return
    setRoleLoading(true)
    try {
      await saveUserRoles(selectedRow.user_id, orgId, Array.from(assigned))
      setRoleModalOpen(false)
      setSelectedRow(null)
    } catch {
      setRoleError("Error occurred")
    } finally {
      setRoleLoading(false)
    }
  }

  /** ------- render ------- */
  return (
    <div className="relative flex h-full flex-col gap-0 overflow-hidden ">
      {progress > 0 && (
        <div className="absolute inset-x-0 top-0">
          <Progress value={progress} className="h-1 rounded-none" aria-label="Loading" />
        </div>
      )}

      <div className="flex flex-col gap-3 border-b border-gray-100 px-6 py-4 dark:border-gray-800 md:flex-row md:items-center md:justify-between">
        <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
          {showTitle ? (titleText ?? t('user')) : null}
          {!showTitle ? (
            <VisuallyHidden>
              <span>{titleText ?? t('user')}</span>
            </VisuallyHidden>
          ) : null}
        </h1>
        <div className="flex items-center gap-2">
          {headerRight}
          <Button
            onClick={() => {
              setOpenAdd(true)
              setSelectedUser(null)
            }}
            className="gap-2"
          >
            <Plus className="h-4 w-4" />
            <span className="lowercase first-letter:uppercase">{t('add_user')}</span>
          </Button>
        </div>
      </div>

      <Separator />

      <div>
        {orgErr ? (
          <div className="text-sm text-red-600">{orgErr}</div>
        ) : (
          <div className="flex items-center justify-between gap-x-4">
            <div className="flex gap-x-2 py-2">
              <div className="relative">
                <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
                <Input
                  value={filterName}
                  onChange={(e) => setFilterName(e.target.value)}
                  placeholder={t('search_by_name')}
                  className="h-9 pl-8 sm:w-64"
                />
              </div>
            </div>

            <div className="flex items-center gap-x-2">
              <p className="text-muted-foreground text-sm">{t('rows')}</p>
              <Select value={String(pageSize)} onValueChange={(v) => changePageSize(Number(v))}>
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
        )}
      </div>

      <Separator />

      <div className="flex min-h-0 flex-1 flex-col overflow-hidden p-0 py-0">
        {fetchError ? (
          <div className="p-6 text-sm text-red-600">{fetchError}</div>
        ) : (
          <div className="flex min-h-0 flex-1 flex-col">
            <div className="min-w-full overflow-x-auto">
              <Table className="w-full table-fixed">
                <colgroup>
                  <col style={{ width: '100px' }} />
                  <col style={{ width: '90px' }} />
                  <col style={{ width: '160px' }} />
                  <col style={{ width: '180px' }} />
                  <col style={{ width: '180px' }} />
                  <col />
                </colgroup>
                <TableHeader>
                  <TableRow className="[&>th]:bg-gray-50 dark:[&>th]:bg-gray-800 [&>th]:sticky [&>th]:top-0 [&>th]:z-20">
                    <TableHead></TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead>{t('reg_no')}</TableHead>
                    <TableHead>{t('last_name')}</TableHead>
                    <TableHead>{t('name')}</TableHead>
                    <TableHead>{t('created_date')}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {loading ? (
                    <TableRow>
                      <TableCell colSpan={8} className="py-6 text-center">
                        <Loader2 className="mx-auto h-4 w-4 animate-spin" />
                      </TableCell>
                    </TableRow>
                  ) : rows.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-muted-foreground py-6 text-center">
                        {t('no_information_available')}
                      </TableCell>
                    </TableRow>
                  ) : (
                    rows.map((row) => (
                      <TableRow key={row.user_id} className="[&>td]:align-middle">
                        <TableCell className="flex gap-2">
                          <Button
                            variant="default"
                            size="sm"
                            aria-label="Change role"
                            onClick={() => onOpenRole(row)}
                          >
                            <ShieldCheck className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="destructive"
                            size="sm"
                            aria-label="Delete"
                            onClick={() => onOpenDelete(row)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </TableCell>
                        <TableCell>{row.user_id}</TableCell>
                        <TableCell>{row.reg_no}</TableCell>
                        <TableCell className="capitalize">{row.last_name}</TableCell>
                        <TableCell className="font-medium capitalize">{row.first_name}</TableCell>
                        <TableCell className="text-muted-foreground">
                          {row.created_date
                            ? new Date(row.created_date).toLocaleString('sv').replace(',', '')
                            : '-'}
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </div>
        )}
      </div>

      <Separator />

      <div className="flex items-center justify-between gap-3 border-t border-gray-100 px-6 py-4 dark:border-gray-800">
        <div className="text-muted-foreground w-auto min-w-72 text-sm">
          {t.rich('showing', {
            from: () => <span className="font-medium">{showingFrom}</span>,
            to: () => <span className="font-medium">{showingTo}</span>,
            total: () => <span className="font-medium">{totalRow || 0}</span>,
          })}
        </div>
        <Pagination className="flex w-full justify-end">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                aria-disabled={!canPrev}
                className={!canPrev ? 'pointer-events-none opacity-50' : ''}
                onClick={() => canPrev && goPage(pageNumber - 1)}
              />
            </PaginationItem>
            {pageLinks.map((p, idx) =>
              p === 'ellipsis' ? (
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

      {/* -------- Delete confirm -------- */}
      <Dialog open={openDelete} onOpenChange={setOpenDelete}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('delete', { name: t('user') })}
            </DialogTitle>
            <DialogDescription className="pt-2 text-base">
              {t.rich('delete_warning', {
                name: () => (
                  <span className="font-medium capitalize">
                    {selectedRow?.last_name} {selectedRow?.first_name}
                  </span>
                ),
              })}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="pt-2">
            <Button variant="outline" onClick={() => setOpenDelete(false)}>
              {t('cancel')}
            </Button>
            <Button variant="destructive" onClick={onDelete}>
              <Trash2 className="mr-2 h-4 w-4" /> {t('delete', { name: '' })}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* -------- Add user (SINGLE via <UserFind>) -------- */}
      <Dialog open={openAdd} onOpenChange={setOpenAdd}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{t('add_user')}</DialogTitle>
            <DialogDescription>
              {t('search_by_name')} / {t('reg_no')} (core)
            </DialogDescription>
          </DialogHeader>

          {/* Хайлт + Preview бүгдийг UserFind хариуцна */}
          <UserFind onChange={setSelectedUser} autoFocus />

          <DialogFooter>
            <Button variant="outline" onClick={() => setOpenAdd(false)}>
              {t('cancel')}
            </Button>
            <Button
              onClick={addUserToOrg}
              disabled={!selectedUser || alreadyInOrg}
              title={alreadyInOrg ? t('already_exists') : t('add')}
            >
              {t('add')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* -------- Role change -------- */}
      <Dialog open={roleModalOpen} onOpenChange={setRoleModalOpen}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>{t('role')}</DialogTitle>
            <DialogDescription className="capitalize">
              {selectedRow
                ? `${t('change')} · ${selectedRow.last_name ?? ''} ${selectedRow.first_name ?? ''}`
                : t('change')}
            </DialogDescription>
          </DialogHeader>

          {roleError && (
            <div className="mb-2 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
              {roleError}
            </div>
          )}

          {roleLoading ? (
            <div className="text-muted-foreground py-8 text-center text-sm">
              <Loader2 className="mx-auto mb-2 h-4 w-4 animate-spin" />
              {t('loading')}
            </div>
          ) : allRoles.length === 0 ? (
            <div className="text-muted-foreground py-8 text-center text-sm">
              {t('no_information_available')}
            </div>
          ) : (
            <div className="max-h-[50vh] overflow-auto rounded-md border p-2">
              <ul className="space-y-2">
                {allRoles.map((r) => {
                  const on = assigned.has(r.id)
                  return (
                    <li
                      key={r.id}
                      onClick={() => toggleRole(r.id)}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter' || e.key === ' ') {
                          e.preventDefault()
                          toggleRole(r.id)
                        }
                      }}
                      role="button"
                      tabIndex={0}
                      aria-pressed={on}
                      className={cn(
                        'hover:bg-muted/40 flex cursor-pointer items-center justify-between gap-3 rounded-md p-2',
                        on && 'bg-muted/60',
                      )}
                    >
                      <div className="flex items-center gap-3">
                        <Checkbox
                          checked={on}
                          onCheckedChange={(v) => toggleRole(r.id, Boolean(v))}
                          onClick={(e) => e.stopPropagation()}
                        />
                        <div>
                          <div className="font-medium">{r.name}</div>
                          {r.description ? (
                            <div className="text-muted-foreground text-xs">{r.description}</div>
                          ) : null}
                        </div>
                      </div>
                    </li>
                  )
                })}
              </ul>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setRoleModalOpen(false)}>
              {t('cancel')}
            </Button>
            <Button onClick={saveRoles} disabled={roleLoading || !selectedRow}>
              {t('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

/** ===================== Page wrapper ===================== */
export default function OrgUsers({ orgId }: { orgId?: number | string }) {
  const t = useTranslations()
  return (
    <div className="h-full w-full overflow-hidden p-4 sm:p-6">
      <OrgUsersBody orgIdProp={orgId} headerRight={null} showTitle={true} titleText={t('user')} />
    </div>
  )
}

/** ===================== Dialog wrapper ===================== */
export function OrgUsersDialog({
  open,
  onOpenChange,
  orgId,
}: {
  open: boolean
  onOpenChange: (v: boolean) => void
  orgId?: number | string
}) {
  const t = useTranslations()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="rounded-2xl p-0 sm:max-w-5xl sm:p-0 [&>button.absolute]:hidden">
        <DialogHeader className="sr-only">
          <VisuallyHidden>
            <DialogTitle>{t('user')}</DialogTitle>
          </VisuallyHidden>
        </DialogHeader>
        <OrgUsersBody
          orgIdProp={orgId}
          headerRight={null}
          showTitle={false}
          titleText={t('user')}
        />
      </DialogContent>
    </Dialog>
  )
}
