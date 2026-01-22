/**
 * 🔔 Notification Page (/[locale]/(main)/admin/notification/page.tsx)
 * 
 * Энэ нь мэдэгдэл удирдах админ хуудас юм.
 * Зорилго: Мэдэгдлийн жагсаалт, илгээх (нэг эсвэл олон хүн рүү)
 * 
 * Features:
 * - ✅ List notifications with pagination
 * - ✅ Filter by read/unread status
 * - ✅ Send notification (single user or broadcast)
 * - ✅ User find from core system
 * - ✅ Progress bar loading
 * - ✅ Form validation (Zod)
 * 
 * Table Columns:
 * - ID
 * - Title
 * - Content
 * - User ID
 * - Group ID
 * - Read Status
 * - Created Date
 * 
 * API Endpoints:
 * - GET /notification - List notifications with pagination
 * - POST /notification - Send notification
 * - POST /user/find-from-core - Find user from core system
 * 
 * @author Gerege Core Team
 */

'use client'

import React, { useEffect, useMemo, useRef, useState } from 'react'
import { useForm, type SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { cn } from '@/lib/utils'
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
import { Plus, Loader2, Search } from 'lucide-react'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'
import { Badge } from '@/components/ui/badge'

type NotificationRow = App.Notification
type UserRow = App.User

/** ===== Schemas ===== */
const FindUserSchema = z.object({
  reg_no: z.string().min(1, 'Required'),
})

const SendNotificationSchema = z.object({
  tenant: z.string().min(1, 'Required'),
  user_id: z.number().int().min(0).optional().default(0), // 0 = broadcast_all
  title: z.string().min(1, 'Required'),
  content: z.string().min(1, 'Required'),
  idempotency_key: z.string().optional().default(''),
})
type SendNotificationIn = z.input<typeof SendNotificationSchema>

export default function NotificationPage() {
  const t = useTranslations()

  /** ========= state ========= */
  const [rows, setRows] = useState<NotificationRow[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)

  const [filterRead, setFilterRead] = useState<'all' | 'read' | 'unread'>('all')
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  const progress = useLoadingProgress(loading)

  const lastReqId = useRef(0)

  const load = async (page = pageNumber, size = pageSize) => {
    const reqId = ++lastReqId.current
    setLoading(true)
    setFetchError(null)

    try {
      const data = await api.get<App.ListData<NotificationRow>>('/notification', {
        query: { page, size },
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
    load(pageNumber, pageSize)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageNumber, pageSize])

  // Filter rows by read status
  const filteredRows = useMemo(() => {
    if (filterRead === 'all') return rows
    if (filterRead === 'read') return rows.filter((r) => r.is_read)
    return rows.filter((r) => !r.is_read)
  }, [rows, filterRead])

  const [openSend, setOpenSend] = useState(false)

  const findUserForm = useForm<{ reg_no: string }>({
    resolver: zodResolver(FindUserSchema),
    defaultValues: { reg_no: '' },
  })

  const sendForm = useForm<SendNotificationIn>({
    resolver: zodResolver(SendNotificationSchema),
    defaultValues: {
      tenant: 'template.gerege.mn',
      user_id: 0,
      title: '',
      content: '',
      idempotency_key: '',
    },
  })

  /** Core-оос хайлт */
  const [finding, setFinding] = useState(false)
  const [foundUser, setFoundUser] = useState<Partial<UserRow> | null>(null)
  const [findError, setFindError] = useState<string | null>(null)

  const doFindFromCore = async () => {
    setFindError(null)
    setFoundUser(null)
    const { reg_no } = findUserForm.getValues()
    const parsed = FindUserSchema.safeParse({ reg_no })
    if (!parsed.success) {
      findUserForm.trigger('reg_no')
      return
    }
    setFinding(true)
    try {
      const data = await api.post<UserRow>('/user/find-from-core', { search_text: reg_no })
      setFoundUser(data)
      sendForm.setValue('user_id', data.id || 0)
    } catch {
      setFindError("Error occurred")
      setFoundUser(null)
      sendForm.setValue('user_id', 0)
    } finally {
      setFinding(false)
    }
  }

  const clearUserFind = () => {
    setFoundUser(null)
    setFindError(null)
    findUserForm.reset({ reg_no: '' })
    sendForm.setValue('user_id', 0)
  }

  /** CRUD */
  const onSend: SubmitHandler<SendNotificationIn> = async (values) => {
    try {
      await api.post<NotificationRow>('/notification', values)
      setOpenSend(false)
      sendForm.reset({
        tenant: 'template.gerege.mn',
        user_id: 0,
        title: '',
        content: '',
        idempotency_key: '',
      })
      clearUserFind()
      setPageNumber(1)
      load(1, pageSize)
    } catch {
      setFetchError("Error occurred")
    }
  }

  /** ========= pagination / filters ========= */
  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) setPageNumber(target)
  }

  const changePageSize = (val: number) => {
    const s = Math.max(5, Math.min(200, val || 50))
    setPageSize(s)
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

  const showingFrom = (pageNumber - 1) * pageSize + (filteredRows.length ? 1 : 0)
  const showingTo = (pageNumber - 1) * pageSize + filteredRows.length
  const canPrev = pageNumber > 1
  const canNext = pageNumber < totalPage

  const unreadCount = rows.filter((r) => !r.is_read).length
  const _isBroadcast = !foundUser || sendForm.watch('user_id') === 0

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
              <div className="flex items-center gap-3">
                <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                  {t('notifications')}
                </h1>
                {unreadCount > 0 && (
                  <Badge variant="default" className="bg-primary-500">
                    {unreadCount} {t('unread')}
                  </Badge>
                )}
              </div>
              <Button
                onClick={() => {
                  setOpenSend(true)
                  sendForm.reset({
                    tenant: 'template.gerege.mn',
                    user_id: 0,
                    title: '',
                    content: '',
                    idempotency_key: '',
                  })
                  clearUserFind()
                }}
                className="gap-2"
                disabled={loading}
              >
                <Plus className="h-4 w-4" />
                <span className="lowercase first-letter:uppercase">{t('send_notification')}</span>
              </Button>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <Select
                  value={filterRead}
                  onValueChange={(v) => setFilterRead(v as 'all' | 'read' | 'unread')}
                >
                  <SelectTrigger className="h-9 w-40">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">{t('all')}</SelectItem>
                    <SelectItem value="read">{t('read')}</SelectItem>
                    <SelectItem value="unread">{t('unread')}</SelectItem>
                  </SelectContent>
                </Select>
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
                        <col style={{ width: '200px' }} /> {/* title */}
                        <col style={{ width: '300px' }} /> {/* content */}
                        <col style={{ width: '100px' }} /> {/* user_id */}
                        <col style={{ width: '100px' }} /> {/* group_id */}
                        <col style={{ width: '120px' }} /> {/* is_read */}
                        <col style={{ width: '160px' }} /> {/* created_date */}
                      </colgroup>
                      <TableHeader>
                        <TableRow className="[&>th]:bg-gray-50 dark:[&>th]:bg-gray-800 [&>th]:sticky [&>th]:top-0 [&>th]:z-20">
                          <TableHead>ID</TableHead>
                          <TableHead>{t('title')}</TableHead>
                          <TableHead>{t('content')}</TableHead>
                          <TableHead>{t('user_id')}</TableHead>
                          <TableHead>{t('group_id')}</TableHead>
                          <TableHead>{t('status')}</TableHead>
                          <TableHead>{t('created_date')}</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {loading ? (
                          <TableRow>
                            <TableCell colSpan={7} className="py-6 text-center">
                              <Loader2 className="mx-auto h-4 w-4 animate-spin" />
                            </TableCell>
                          </TableRow>
                        ) : filteredRows.length === 0 ? (
                          <TableRow>
                            <TableCell colSpan={7} className="text-muted-foreground py-6 text-center">
                              {t('no_information_available')}
                            </TableCell>
                          </TableRow>
                        ) : (
                          filteredRows.map((r) => (
                            <TableRow
                              key={r.id}
                              className={cn(
                                '[&>td]:align-middle',
                                !r.is_read && 'bg-primary/10 dark:bg-primary/20',
                              )}
                            >
                              <TableCell className="text-muted-foreground">{r.id}</TableCell>
                              <TableCell className="font-medium">{r.title}</TableCell>
                              <TableCell className="text-muted-foreground">
                                <div className="line-clamp-2">{r.content}</div>
                              </TableCell>
                              <TableCell className="text-muted-foreground">
                                {r.user_id || <span className="opacity-60">—</span>}
                              </TableCell>
                              <TableCell className="text-muted-foreground">{r.group_id}</TableCell>
                              <TableCell>
                                {r.is_read ? (
                                  <Badge variant="secondary" className="bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400">
                                    {t('read')}
                                  </Badge>
                                ) : (
                                  <Badge variant="default" className="bg-primary-500">
                                    {t('unread')}
                                  </Badge>
                                )}
                              </TableCell>
                              <TableCell className="text-muted-foreground">
                                {r.created_date?.slice(0, 10) || '-'}
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

        {/* ---------- Send Notification Dialog ---------- */}
        <Dialog open={openSend} onOpenChange={setOpenSend}>
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle className="lowercase first-letter:uppercase">
                {t('send_notification')}
              </DialogTitle>
              <DialogDescription>{t('fill_the_fields_and_save_to_send')}</DialogDescription>
            </DialogHeader>

            <Form {...sendForm}>
              <form
                onSubmit={sendForm.handleSubmit(onSend)}
                className="space-y-4 pt-2"
                autoComplete="off"
              >
                <FormField
                  control={sendForm.control}
                  name="tenant"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('tenant')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('tenant')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* User Find Section */}
                <div className="space-y-2">
                  <FormLabel>{t('user')} ({t('optional')})</FormLabel>
                  <Form {...findUserForm}>
                    <div className="grid grid-cols-1 gap-3 sm:grid-cols-[1fr_auto]">
                      <FormField
                        control={findUserForm.control}
                        name="reg_no"
                        render={({ field }) => (
                          <FormItem>
                            <FormControl>
                              <Input placeholder={t('reg_no')} {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <div className="flex items-end gap-2">
                        <Button
                          type="button"
                          variant="secondary"
                          onClick={doFindFromCore}
                          disabled={finding}
                          className="w-full sm:w-auto"
                        >
                          {finding ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          ) : (
                            <Search className="mr-2 h-4 w-4" />
                          )}
                          {t('search')}
                        </Button>
                        {foundUser && (
                          <Button
                            type="button"
                            variant="outline"
                            onClick={clearUserFind}
                            className="w-full sm:w-auto"
                          >
                            {t('clear')}
                          </Button>
                        )}
                      </div>
                    </div>
                  </Form>

                  {findError && (
                    <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
                      {findError}
                    </div>
                  )}

                  {foundUser && (
                    <div className="rounded-lg border p-3">
                      <div className="grid gap-2 sm:grid-cols-2">
                        <div>
                          <p className="text-muted-foreground text-xs">{t('reg_no')}</p>
                          <p className="font-medium">{foundUser?.reg_no || '—'}</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground text-xs">{t('name')}</p>
                          <p className="font-medium">
                            {[
                              foundUser?.family_name,
                              foundUser?.last_name,
                              foundUser?.first_name,
                            ]
                              .filter(Boolean)
                              .join(' ') || '—'}
                          </p>
                        </div>
                        <div>
                          <p className="text-muted-foreground text-xs">{t('user_id')}</p>
                          <p className="font-medium">{foundUser?.id || '—'}</p>
                        </div>
                      </div>
                    </div>
                  )}

                  {!foundUser && !findError && (
                    <p className="text-muted-foreground text-sm">
                      {t('leave_empty_to_broadcast_to_all')}
                    </p>
                  )}
                </div>

                <FormField
                  control={sendForm.control}
                  name="title"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('title')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('title')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={sendForm.control}
                  name="content"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('content')}</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t('content')}
                          rows={4}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={sendForm.control}
                  name="idempotency_key"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('idempotency_key')} ({t('optional')})</FormLabel>
                      <FormControl>
                        <Input placeholder={t('idempotency_key')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <DialogFooter>
                  <Button type="button" variant="outline" onClick={() => setOpenSend(false)}>
                    {t('cancel')}
                  </Button>
                  <Button type="submit" disabled={sendForm.formState.isSubmitting}>
                    {sendForm.formState.isSubmitting && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    {t('send')}
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          </DialogContent>
        </Dialog>
      </div>
    </>
  )
}
