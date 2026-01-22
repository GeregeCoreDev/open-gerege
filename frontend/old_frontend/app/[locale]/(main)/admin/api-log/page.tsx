/**
 * 📋 API Log Page (/[locale]/(main)/admin/api-log/page.tsx)
 *
 * Энэ нь API log удирдах админ хуудас юм.
 * Зорилго: API хүсэлтийн логуудыг харах, шүүх
 *
 * Features:
 * - ✅ List API logs with pagination
 * - ✅ Filter by method, path, status_code, user_id, org_id, ip
 * - ✅ Progress bar loading
 * - ✅ Status code badge coloring
 *
 * Table Columns:
 * - ID
 * - Method
 * - Path
 * - Status Code
 * - User ID
 * - Username
 * - Org ID
 * - IP
 * - Created Date
 *
 * API Endpoints:
 * - GET /api-logs - List API logs with pagination and filters
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React, { useEffect, useMemo, useRef, useState } from 'react'
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
import { Input } from '@/components/ui/input'
import { Loader2 } from 'lucide-react'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'
import { Badge } from '@/components/ui/badge'

type APILogRow = {
  id: number
  org_id?: number | null
  user_id?: number | null
  username?: string
  path: string
  method: string
  status_code: number
  ip?: string
  created_date: string
}

/** ===== Schemas ===== */

export default function APILogPage() {
  const t = useTranslations()

  /** ========= state ========= */
  const [rows, setRows] = useState<APILogRow[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)

  // Filters
  const [filterMethod, setFilterMethod] = useState<string>('')
  const [filterPath, setFilterPath] = useState<string>('')
  const [filterStatusCode, setFilterStatusCode] = useState<string>('')
  const [filterUserId, setFilterUserId] = useState<string>('')
  const [filterOrgId, setFilterOrgId] = useState<string>('')
  const [filterIP, setFilterIP] = useState<string>('')

  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  const progress = useLoadingProgress(loading)

  const lastReqId = useRef(0)

  const load = async (
    page = pageNumber,
    size = pageSize,
    method = filterMethod,
    path = filterPath,
    statusCode = filterStatusCode,
    userId = filterUserId,
    orgId = filterOrgId,
    ip = filterIP,
  ) => {
    const reqId = ++lastReqId.current
    setLoading(true)
    setFetchError(null)

    try {
      const query: Record<string, string | number | undefined> = {
        page,
        size,
      }

      if (method) query.method = method
      if (path) query.path = path
      if (statusCode) query.status_code = Number(statusCode)
      if (userId) query.user_id = Number(userId)
      if (orgId) query.org_id = Number(orgId)
      if (ip) query.ip = ip

      const data = await api.get<App.ListData<APILogRow>>('/api-logs', {
        query,
      })
      if (reqId !== lastReqId.current) return

      const m = data.meta
      setRows(data.items ?? [])
      setPageNumber(m?.page ?? page)
      setPageSize(m?.size ?? size)
      setTotalPage(m?.pages ?? 1)
      setTotalRow(m?.total ?? 0)
    } catch {
      if (reqId === lastReqId.current) {
        setFetchError("Error occurred")
      }
    } finally {
      if (reqId === lastReqId.current) setLoading(false)
    }
  }

  useEffect(() => {
    load(
      pageNumber,
      pageSize,
      filterMethod,
      filterPath,
      filterStatusCode,
      filterUserId,
      filterOrgId,
      filterIP,
    )
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    pageNumber,
    pageSize,
    filterMethod,
    filterPath,
    filterStatusCode,
    filterUserId,
    filterOrgId,
    filterIP,
  ])

  /** ========= pagination / filters ========= */
  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) setPageNumber(target)
  }

  const _changePageSize = (val: number) => {
    const s = Math.max(5, Math.min(200, val || 50))
    setPageSize(s)
    setPageNumber(1)
  }

  const clearFilters = () => {
    setFilterMethod('')
    setFilterPath('')
    setFilterStatusCode('')
    setFilterUserId('')
    setFilterOrgId('')
    setFilterIP('')
    setPageNumber(1)
  }

  /** ========= computed ========= */
  const pageLinks = useMemo(() => {
    const links: (number | 'ellipsis')[] = []
    if (totalPage <= 7) {
      for (let i = 1; i <= totalPage; i++) links.push(i)
      return links
    }
    const windowSize = 2
    links.push(1)
    if (pageNumber > 1 + windowSize + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - windowSize)
    const end = Math.min(totalPage - 1, pageNumber + windowSize)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < totalPage - windowSize - 1) links.push('ellipsis')
    links.push(totalPage)
    return links
  }, [pageNumber, totalPage])

  const showingFrom = (pageNumber - 1) * pageSize + (rows.length ? 1 : 0)
  const showingTo = (pageNumber - 1) * pageSize + rows.length
  const canPrev = pageNumber > 1
  const canNext = pageNumber < totalPage

  const getStatusBadge = (status: number) => {
    if (status >= 200 && status < 300) {
      return (
        <Badge
          variant="secondary"
          className="bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
        >
          {status}
        </Badge>
      )
    }
    if (status >= 300 && status < 400) {
      return (
        <Badge
          variant="secondary"
          className="bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
        >
          {status}
        </Badge>
      )
    }
    if (status >= 400 && status < 500) {
      return (
        <Badge
          variant="secondary"
          className="bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400"
        >
          {status}
        </Badge>
      )
    }
    if (status >= 500) {
      return (
        <Badge
          variant="secondary"
          className="bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
        >
          {status}
        </Badge>
      )
    }
    return <Badge variant="secondary">{status}</Badge>
  }

  const hasActiveFilters =
    filterMethod || filterPath || filterStatusCode || filterUserId || filterOrgId || filterIP

  /** ========= render ========= */
  return (
    <>
      <div className="h-full w-full">
        <div className="relative flex h-full flex-col overflow-hidden">
          {progress > 0 && (
            <div className="absolute inset-x-0 top-0 z-10">
              <Progress value={progress} className="h-1 rounded-none" aria-label="Loading" />
            </div>
          )}

          <div className="flex flex-col overflow-hidden px-4 pb-4">
            {/* Header */}
            <div className="flex flex-col gap-3 pt-4 md:flex-row md:items-center md:justify-between">
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                {t('api_log')}
              </h1>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-4 dark:border-gray-800 dark:bg-gray-900/50">
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-4 lg:grid-cols-6">
                <div className="space-y-1">
                  <label className="text-muted-foreground text-xs">{t('method')}</label>
                  <Select
                    value={filterMethod || 'all'}
                    onValueChange={(v) => setFilterMethod(v === 'all' ? '' : v)}
                  >
                    <SelectTrigger className="h-9 w-full">
                      <SelectValue placeholder={t('all')} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">{t('all')}</SelectItem>
                      <SelectItem value="GET">GET</SelectItem>
                      <SelectItem value="POST">POST</SelectItem>
                      <SelectItem value="PUT">PUT</SelectItem>
                      <SelectItem value="PATCH">PATCH</SelectItem>
                      <SelectItem value="DELETE">DELETE</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-1">
                  <label className="text-muted-foreground text-xs">{t('path')}</label>
                  <Input
                    value={filterPath}
                    onChange={(e) => {
                      setFilterPath(e.target.value)
                      setPageNumber(1)
                    }}
                    placeholder={t('path')}
                    className="h-9"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-muted-foreground text-xs">{t('status_code')}</label>
                  <Input
                    type="number"
                    value={filterStatusCode}
                    onChange={(e) => {
                      setFilterStatusCode(e.target.value)
                      setPageNumber(1)
                    }}
                    placeholder={t('status_code')}
                    className="h-9"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-muted-foreground text-xs">{t('user_id')}</label>
                  <Input
                    type="number"
                    value={filterUserId}
                    onChange={(e) => {
                      setFilterUserId(e.target.value)
                      setPageNumber(1)
                    }}
                    placeholder={t('user_id')}
                    className="h-9"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-muted-foreground text-xs">{t('org_id')}</label>
                  <Input
                    type="number"
                    value={filterOrgId}
                    onChange={(e) => {
                      setFilterOrgId(e.target.value)
                      setPageNumber(1)
                    }}
                    placeholder={t('org_id')}
                    className="h-9"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-muted-foreground text-xs">{t('ip')}</label>
                  <Input
                    value={filterIP}
                    onChange={(e) => {
                      setFilterIP(e.target.value)
                      setPageNumber(1)
                    }}
                    placeholder={t('ip')}
                    className="h-9"
                  />
                </div>
              </div>

              {hasActiveFilters && (
                <div className="flex justify-end">
                  <button
                    onClick={clearFilters}
                    className="text-primary-600 hover:text-primary-700 dark:text-primary-400 text-sm font-medium"
                  >
                    {t('clear_filters')}
                  </button>
                </div>
              )}
            </div>

            {/* Table */}
            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
              {fetchError ? (
                <div className="p-6 text-sm text-red-600">{fetchError}</div>
              ) : (
                <div className="flex min-h-0 flex-1 flex-col">
                  <div className="min-w-full overflow-x-auto">
                    <Table className="w-full table-fixed">
                      <colgroup>
                        <col style={{ width: '80px' }} /> {/* id */}
                        <col style={{ width: '80px' }} /> {/* method */}
                        <col style={{ width: '250px' }} /> {/* path */}
                        <col style={{ width: '100px' }} /> {/* status_code */}
                        <col style={{ width: '100px' }} /> {/* user_id */}
                        <col style={{ width: '150px' }} /> {/* username */}
                        <col style={{ width: '100px' }} /> {/* org_id */}
                        <col style={{ width: '120px' }} /> {/* ip */}
                        <col style={{ width: '160px' }} /> {/* created_date */}
                      </colgroup>
                      <TableHeader>
                        <TableRow className="[&>th]:sticky [&>th]:top-0 [&>th]:z-20 [&>th]:bg-gray-50 dark:[&>th]:bg-gray-800">
                          <TableHead>ID</TableHead>
                          <TableHead>{t('method')}</TableHead>
                          <TableHead>{t('path')}</TableHead>
                          <TableHead>{t('status_code')}</TableHead>
                          <TableHead>{t('user_id')}</TableHead>
                          <TableHead>{t('username')}</TableHead>
                          <TableHead>{t('org_id')}</TableHead>
                          <TableHead>{t('ip')}</TableHead>
                          <TableHead>{t('created_date')}</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {loading ? (
                          <TableRow>
                            <TableCell colSpan={9} className="py-6 text-center">
                              <Loader2 className="mx-auto h-4 w-4 animate-spin" />
                            </TableCell>
                          </TableRow>
                        ) : rows.length === 0 ? (
                          <TableRow>
                            <TableCell
                              colSpan={9}
                              className="text-muted-foreground py-6 text-center"
                            >
                              {t('no_information_available')}
                            </TableCell>
                          </TableRow>
                        ) : (
                          rows.map((r) => (
                            <TableRow key={r.id} className="[&>td]:align-middle">
                              <TableCell className="text-muted-foreground">{r.id}</TableCell>
                              <TableCell>
                                <Badge variant="outline" className="font-mono text-xs">
                                  {r.method}
                                </Badge>
                              </TableCell>
                              <TableCell className="font-mono text-xs">{r.path}</TableCell>
                              <TableCell>{getStatusBadge(r.status_code)}</TableCell>
                              <TableCell className="text-muted-foreground">
                                {r.user_id || <span className="opacity-60">—</span>}
                              </TableCell>
                              <TableCell className="text-muted-foreground">
                                {r.username || <span className="opacity-60">—</span>}
                              </TableCell>
                              <TableCell className="text-muted-foreground">
                                {r.org_id || <span className="opacity-60">—</span>}
                              </TableCell>
                              <TableCell className="text-muted-foreground font-mono text-xs">
                                {r.ip || <span className="opacity-60">—</span>}
                              </TableCell>
                              <TableCell className="text-muted-foreground">
                                {r.created_date ? new Date(r.created_date).toLocaleString() : '-'}
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

            {/* Pagination */}
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
          </div>
        </div>
      </div>
    </>
  )
}
