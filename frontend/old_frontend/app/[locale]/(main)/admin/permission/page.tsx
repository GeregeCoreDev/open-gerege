/**
 * 🔐 Module Permission Page (/[locale]/(main)/admin/module-permission/page.tsx)
 *
 * Энэ нь модулийн эрх удирдах админ хуудас юм.
 * Зорилго: Модуль тус бүрийн CRUD permissions удирдлага
 *
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Search by name filter
 * - ✅ Filter by module
 * - ✅ Progress bar loading
 * - ✅ Active/Inactive toggle
 * - ✅ Badge display for code
 * - ✅ Form validation (Zod)
 * - ✅ Module selection dropdown
 *
 * Table Columns:
 * - Actions (Edit/Delete)
 * - ID
 * - Code (with Badge)
 * - Name
 * - Module (related)
 * - Description
 * - Active status
 * - Created Date
 *
 * Form Fields:
 * - code: Permission code (required, e.g. 'user.create')
 * - name: Display name (required)
 * - description: Optional
 * - module_id: Parent module (dropdown)
 * - is_active: Enable/Disable
 *
 * Permission Naming Convention:
 * - {module}.{action}
 * - Examples: 'user.create', 'user.read', 'user.update', 'user.delete'
 *
 * Related Entities:
 * - Module: Parent entity
 * - Role: Roles are granted permissions
 *
 * Authorization Flow:
 * User → Role → Permission → Module
 *
 * API Endpoints:
 * - GET /module-permission - List permissions
 * - GET /module - Get modules for dropdown
 * - POST /module-permission - Create permission
 * - PUT /module-permission - Update permission
 * - DELETE /module-permission/:id - Delete permission
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React, { useEffect, useMemo, useState } from 'react'
import { Badge } from '@/components/ui/badge'
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
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationLink,
  PaginationEllipsis,
} from '@/components/ui/pagination'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'
import { useTranslations } from 'next-intl'
import { Loader2 } from 'lucide-react'

type Row = App.Permission

export default function PermissionsPage() {
  const t = useTranslations()

  // list state
  const [rows, setRows] = useState<Row[]>([])
  const [loading, setLoading] = useState(false)
  const [deleting, _setDeleting] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)

  // meta/pagination from ListData
  const [meta, setMeta] = useState<App.ApiMeta | null>(null)
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  // filters
  const [searchQuery, setSearchQuery] = useState('')
  const [searchType, setSearchType] = useState<'name' | 'code' | 'both'>('both')
  const [filterModule, setFilterModule] = useState<number | 'all'>('all')

  // modules for filter dropdown (all modules from /module endpoint)
  const [modules, setModules] = useState<App.Module[]>([])

  // progress bar - simplified with hook
  const progress = useLoadingProgress(loading || deleting)

  /** ---------- Load list ---------- */
  async function load(page = pageNumber, size = pageSize, moduleFilter = filterModule) {
    setLoading(true)
    setFetchError(null)
    try {
      // Build query parameter based on search type
      const queryParams: Record<string, string | number> = {
        page,
        size,
      }

      if (searchQuery.trim()) {
        if (searchType === 'both') {
          // Format: q=name:value,code:value
          queryParams.q = `name:${searchQuery.trim()},code:${searchQuery.trim()}`
        } else if (searchType === 'name') {
          queryParams.q = `name:${searchQuery.trim()}`
        } else if (searchType === 'code') {
          queryParams.q = `code:${searchQuery.trim()}`
        }
      }

      // Add module filter if selected
      if (moduleFilter !== 'all') {
        queryParams.module_id = moduleFilter
      }

      const data = await api.get<App.ListData<Row>>('/permission', {
        query: queryParams,
      })
      const m = data.meta
      setRows(data.items ?? [])
      setMeta(m)
      setPageNumber(m?.page ?? page)
      setPageSize(m?.size ?? size)
      setTotalPage(m?.pages ?? 1)
      setTotalRow(m?.total ?? 0)
    } catch {
      setFetchError("Error occurred")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load(1, pageSize)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize])

  /** ---------- Load all modules for filter ---------- */
  useEffect(() => {
    async function loadAllModules() {
      try {
        const data = await api.get<App.ListData<App.Module>>('/module', {
          query: { page: 1, size: 50 },
        })
        setModules(data.items ?? [])
      } catch {
        setModules([])
      }
    }
    loadAllModules()
  }, [])

  /** ---------- Pagination helpers ---------- */
  const canPrev = meta?.has_prev ?? pageNumber > 1
  const canNext = meta?.has_next ?? pageNumber < totalPage
  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) load(target, pageSize)
  }

  const pageLinks = useMemo(() => {
    const links: (number | 'ellipsis')[] = []
    const tp = totalPage
    if (tp <= 7) {
      for (let i = 1; i <= tp; i++) links.push(i)
      return links
    }
    const w = 2
    links.push(1)
    if (pageNumber > 1 + w + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - w)
    const end = Math.min(tp - 1, pageNumber + w)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < tp - w - 1) links.push('ellipsis')
    links.push(tp)
    return links
  }, [pageNumber, totalPage])

  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? 0) + 1
  const showingTo = totalRow === 0 ? 0 : (meta?.end_idx ?? -1) + 1

  const ActiveBadge = ({ value }: { value: boolean | null | undefined }) => {
    const on = value === true
    const off = value === false
    return (
      <span
        className={[
          'inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs',
          on
            ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
            : off
              ? 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300'
              : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-300',
        ].join(' ')}
      >
        <span
          className={[
            'h-1.5 w-1.5 rounded-full',
            on ? 'bg-emerald-500' : off ? 'bg-rose-500' : 'bg-gray-400',
          ].join(' ')}
        />
        {on ? t('active') : off ? t('inactive') : ''}
      </span>
    )
  }

  /** ---------- JSX ---------- */
  return (
    <>
      <div className="h-full w-full">
        {/* Loading bar */}
        {progress > 0 && (
          <div className="absolute inset-x-0 top-0">
            <Progress value={progress} className="h-1 rounded-none" aria-label="Уншиж байна" />
          </div>
        )}

        <div className="flex flex-col overflow-hidden px-4 pb-4">
          <div className="flex flex-col gap-3 pt-4 md:flex-row md:items-center md:justify-between dark:border-gray-800">
            <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
              {t('permission')}
            </h1>
          </div>

          {/* Filters */}
          <div>
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2 py-2">
                <div className="flex items-center gap-2">
                  <Select
                    value={searchType}
                    onValueChange={(v) => setSearchType(v as 'name' | 'code' | 'both')}
                  >
                    <SelectTrigger className="h-9 w-35">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="both">
                        {t('name')} & {t('code')}
                      </SelectItem>
                      <SelectItem value="name">{t('name')}</SelectItem>
                      <SelectItem value="code">{t('code')}</SelectItem>
                    </SelectContent>
                  </Select>
                  <Input
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && load(1, pageSize)}
                    placeholder={
                      searchType === 'both'
                        ? `${t('search')}...`
                        : searchType === 'name'
                          ? `${t('search_by_name')}`
                          : `${t('search_by')} ${t('code')}`
                    }
                    className="h-9 sm:w-72"
                  />
                  <div className="flex items-center gap-2">
                    <p className="text-muted-foreground text-sm">{t('module')}:</p>
                    <Select
                      value={filterModule === 'all' ? 'all' : String(filterModule)}
                      onValueChange={(v) => {
                        const newFilter = v === 'all' ? 'all' : Number(v)
                        setFilterModule(newFilter)
                        load(1, pageSize, newFilter)
                      }}
                    >
                      <SelectTrigger className="h-9 w-56">
                        <SelectValue placeholder={t('all')} />
                      </SelectTrigger>
                      <SelectContent className="h-[60vh]">
                        <SelectItem value="all">{t('all')}</SelectItem>
                        {modules.map((m) => (
                          <SelectItem key={m.id} value={String(m.id)}>
                            <Badge>{m.system?.code}</Badge>- {m.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </div>

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
          </div>

          {/* Table Content */}
          <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
            {fetchError ? (
              <div className="p-6 text-sm text-red-600">{fetchError}</div>
            ) : (
              <div className="flex min-h-0 flex-1 flex-col">
                {/* Header */}
                <div className="min-w-full overflow-x-auto">
                  <Table className="w-full table-fixed">
                    <colgroup>
                      <col style={{ width: '80px' }} />
                      <col style={{ width: '250px' }} />
                      <col style={{ width: '200px' }} />
                      <col style={{ width: '300px' }} />
                      <col style={{ width: '120px' }} />
                      <col style={{ width: '100px' }} />
                    </colgroup>
                    <TableHeader>
                      <TableRow className="[&>th]:z-20 [&>th]:bg-gray-50 dark:[&>th]:bg-gray-800">
                        <TableHead>ID</TableHead>
                        <TableHead>{t('code')}</TableHead>
                        <TableHead>{t('name')}</TableHead>
                        <TableHead>{t('description')}</TableHead>
                        <TableHead>{t('module')}</TableHead>
                        <TableHead>{t('is_active')}</TableHead>
                      </TableRow>
                    </TableHeader>
                  </Table>
                </div>

                {loading ? (
                  <div className="flex h-32 w-full flex-col items-center justify-center gap-y-6">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <p className="text-gray-700 dark:text-gray-200">{t('loading')}</p>
                  </div>
                ) : rows.length === 0 ? (
                  <div className="text-muted-foreground flex h-48 w-full flex-col items-center justify-center gap-y-6">
                    <p className="text-gray-700 dark:text-gray-200">
                      {t('no_information_available')}
                    </p>
                  </div>
                ) : (
                  <div className="min-h-0 flex-1 overflow-auto">
                    <div className="min-w-full overflow-x-auto">
                      <Table className="w-full table-fixed">
                        <colgroup>
                          <col style={{ width: '80px' }} />
                          <col style={{ width: '250px' }} />
                          <col style={{ width: '200px' }} />
                          <col style={{ width: '300px' }} />
                          <col style={{ width: '120px' }} />
                          <col style={{ width: '100px' }} />
                        </colgroup>
                        <TableBody>
                          {rows.map((row) => {
                            return (
                              <TableRow key={row.id} className="[&>td]:align-center">
                                <TableCell className="text-muted-foreground">{row.id}</TableCell>
                                <TableCell className="font-medium">
                                  <Badge variant="outline">{row.code}</Badge>
                                </TableCell>
                                <TableCell className="font-medium">{row.name}</TableCell>
                                <TableCell className="text-muted-foreground">
                                  {row.description || '-'}
                                </TableCell>
                                <TableCell className="text-muted-foreground">
                                  {row.module
                                    ? `${row.module.id} - ${row.module.name}`
                                    : row.module_id}
                                </TableCell>
                                <TableCell>
                                  <ActiveBadge value={row.is_active} />
                                </TableCell>
                              </TableRow>
                            )
                          })}
                        </TableBody>
                      </Table>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Footer / Pagination */}
          <div className="flex items-center justify-between gap-3 rounded-b-md border-r border-b border-l border-gray-200 px-6 py-2 dark:border-gray-800">
            <div className="text-muted-foreground w-auto min-w-72 text-sm">
              {t.rich('showing', {
                from: () => <span className="font-medium">{showingFrom || 0}</span>,
                to: () => <span className="font-medium">{showingTo || 0}</span>,
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
                      <PaginationLink
                        className="cursor-pointer"
                        isActive={p === pageNumber}
                        onClick={() => goPage(p)}
                      >
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
        </div>
      </div>
    </>
  )
}
