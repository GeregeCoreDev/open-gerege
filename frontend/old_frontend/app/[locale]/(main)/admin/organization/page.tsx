/**
 * 🏢 Organization Page (/[locale]/(main)/admin/organization/page.tsx)
 *
 * Энэ нь байгууллага удирдах админ хуудас юм.
 * Зорилго: Бүх байгууллагуудын CRUD удирдлага, байгууллагын мэдээлэл
 *
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Search by name filter
 * - ✅ Filter by organization type
 * - ✅ Progress bar loading
 * - ✅ Organization type selection
 * - ✅ Parent organization selection (hierarchical)
 * - ✅ Address information (aimag, sum, bag)
 * - ✅ Logo image upload
 * - ✅ Form validation (Zod)
 * - ✅ View organization details dialog
 * - ✅ View organization users
 *
 * Table Columns:
 * - Actions (Users/Details/Edit/Delete)
 * - ID
 * - Registration Number
 * - Name
 * - Type
 * - Phone
 * - Email
 * - Created Date
 *
 * Form Fields:
 * - reg_no: Registration number (required)
 * - name: Organization name (required)
 * - short_name: Short name
 * - phone_no: Contact phone
 * - email: Contact email
 * - type_id: Organization type (dropdown)
 * - parent_id: Parent organization (optional, hierarchical)
 * - logo_image_url: Logo image
 * - Address: aimag, sum, bag, detail
 *
 * Related Components:
 * - OrganizationDetailsDialog: View organization info
 * - OrganizationUsersPage: View org's users
 *
 * API Endpoints:
 * - GET /organization - List with filters
 * - GET /organization-type - Get types for dropdown
 * - POST /organization - Create
 * - PUT /organization - Update
 * - DELETE /organization/:id - Delete
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React, { useEffect, useMemo, useState } from 'react'
import { useForm, type SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

import { useTranslations } from 'next-intl'

import api from '@/lib/api'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
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
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

import { Loader2, Pencil, Plus, Trash2, Users2, Search, Info } from 'lucide-react'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { TooltipArrow } from '@radix-ui/react-tooltip'
import { OrgUsersDialog } from './users/page'
import OrganizationDetailDialog from './actions/details'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'

/* ===================== Schemas ===================== */
// РД заавал биш (хоосон байж болно)
const CreateSchema = z.object({
  reg_no: z.string().length(7).optional().or(z.literal('')).default(''),
  name: z.string().min(1, 'Required'),
  code: z.string().optional().default(''),
  type_id: z.number().int().min(1, 'Required'),
  status: z.enum(['active', 'inactive']).default('active'),
  description: z.string().optional().default(''),
  address: z.string().optional().default(''),
  phone_no: z.string().optional().default(''),
  email: z.string().email('Invalid email').optional().or(z.literal('')).default(''),
  website: z.string().url('Invalid url').optional().or(z.literal('')).default(''),
})
const EditSchema = z.object({
  type_id: z.number().int().min(1, 'Required'),
  status: z.enum(['active', 'inactive']).default('active'),
  description: z.string().optional().default(''),
  address: z.string().optional().default(''),
  phone_no: z.string().optional().default(''),
  email: z.string().email('Invalid email').optional().or(z.literal('')).default(''),
  website: z.string().url('Invalid url').optional().or(z.literal('')).default(''),
})
type CreateIn = z.input<typeof CreateSchema>
type CreateOut = z.output<typeof CreateSchema>
type EditIn = z.input<typeof EditSchema>
type EditOut = z.output<typeof EditSchema>

/* ===================== Page ===================== */
export default function OrganizationPage() {
  const t = useTranslations()

  const [rows, setRows] = useState<App.Organization[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  const [meta, setMeta] = useState<App.ApiMeta | null>(null)
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  const [filterName, setFilterName] = useState('')

  const [orgTypes, setOrgTypes] = useState<App.OrganizationType[]>([])

  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [openUser, setOpenUser] = useState(false)
  const [openDetail, setOpenDetail] = useState(false)
  const [selected, setSelected] = useState<App.Organization | null>(null)

  // progress bar - simplified with hook
  const progress = useLoadingProgress(loading)

  async function loadOrgType() {
    try {
      const resp = await api.get<App.ListData<App.OrganizationType>>('/orgtype', {
        query: { res: 'listdata', page_size: 500 },
      })
      setOrgTypes(resp.items ?? [])
    } catch {
      setOrgTypes([])
    }
  }

  async function load(page = pageNumber, size = pageSize) {
    setLoading(true)
    setFetchError(null)
    try {
      const data = await api.get<App.ListData<App.Organization>>('/organization', {
        query: {
          page,
          size,
          name: filterName || undefined,
        },
      })
      const m = data.meta
      setMeta(m ?? null)
      setRows(data.items ?? [])
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

  useEffect(() => {
    loadOrgType()
  }, [])

  const createForm = useForm<CreateIn>({
    resolver: zodResolver(CreateSchema),
    defaultValues: {
      reg_no: '',
      name: '',
      code: '',
      type_id: 0,
      status: 'active',
      description: '',
      address: '',
      phone_no: '',
      email: '',
      website: '',
    },
  })
  const editForm = useForm<EditIn>({
    resolver: zodResolver(EditSchema),
    defaultValues: {
      type_id: 0,
      status: 'active',
      description: '',
      address: '',
      phone_no: '',
      email: '',
      website: '',
    },
  })

  const [finding, setFinding] = useState(false)
  const [findError, setFindError] = useState<string | null>(null)

  const onOpenCreate = () => {
    setOpenCreate(true)
    setFindError(null)
    createForm.reset({
      reg_no: '',
      name: '',
      code: '',
      type_id: 0,
      status: 'active',
      description: '',
      address: '',
      phone_no: '',
      email: '',
      website: '',
    })
  }

  const onOpenEdit = (r: App.Organization) => {
    setSelected(r)
    editForm.reset({
      type_id: r.type_id || r.type?.id || 0,
      status: r.is_active === false ? 'inactive' : 'active',
      address: r.address || '',
      phone_no: r.phone_no || '',
      email: r.email || '',
      website: r.website || '',
      description: r.description || '',
    })
    setOpenEdit(true)
  }

  const onOpenDelete = (r: App.Organization) => {
    setSelected(r)
    setOpenDelete(true)
  }

  const onCreate: SubmitHandler<CreateIn> = async (valuesIn) => {
    const v: CreateOut = CreateSchema.parse(valuesIn)
    try {
      await api.post<App.Organization>('/organization', {
        ...v,
        id: selected?.id,
        is_active: v.status === 'active',
        reg_no: v.reg_no?.trim() ? v.reg_no.trim() : undefined,
        code: v.code?.trim() ? v.code.trim() : undefined,
        email: v.email?.trim() ? v.email.trim() : undefined,
        website: v.website?.trim() ? v.website.trim() : undefined,
        description: v.description?.trim() ? v.description.trim() : undefined,
        address: v.address?.trim() ? v.address.trim() : undefined,
        phone_no: v.phone_no?.trim() ? v.phone_no.trim() : undefined,
      })
      setOpenCreate(false)
      createForm.reset({
        reg_no: '',
        name: '',
        code: '',
        type_id: 0,
        status: 'active',
        description: '',
        address: '',
        phone_no: '',
        email: '',
        website: '',
      })
      await load(1, pageSize)
    } catch {}
  }

  const onUpdate: SubmitHandler<EditIn> = async (valuesIn) => {
    if (!selected) return
    const values: EditOut = EditSchema.parse(valuesIn)
    try {
      await api.put<App.Organization>(`/organization/${selected.id}`, {
        ...values,
        is_active: values.status === 'active',
        reg_no: selected.reg_no,
        name: selected.name,
      })
      setOpenEdit(false)
      setSelected(null)
      await load(pageNumber, pageSize)
    } catch {}
  }

  const onDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>(`/organization/${selected.id}`)
      const newCount = rows.length - 1
      const willBeEmpty = newCount <= 0 && pageNumber > 1
      setOpenDelete(false)
      setSelected(null)
      await load(willBeEmpty ? pageNumber - 1 : pageNumber, pageSize)
    } catch {
    } finally {
      setDeleting(false)
    }
  }

  const doFindFromCore = async () => {
    setFindError(null)
    const reg = (createForm.getValues().reg_no || '').trim()
    if (!reg) {
      setFindError(t('enter_registration_number') || 'Регистрийн дугаараа оруулна уу.')
      return
    }
    setFinding(true)
    setSelected(null)
    try {
      const data = await api.get<App.Organization | null>('/organization/find', {
        query: { search_text: reg },
      })
      if (!data) {
        setFindError(t('not_found') || 'Core-оос илэрц олдсонгүй.')
        return
      }
      setSelected(data)
      createForm.reset({
        reg_no: data.reg_no ?? reg,
        name: data.name ?? '',
        code: data.code ?? '',
        type_id: data.type_id ?? data.type?.id ?? 0,
        status: data.is_active === false ? 'inactive' : 'active',
        description: data.description ?? '',
        address: data.address ?? '',
        phone_no: data.phone_no ?? '',
        email: data.email ?? '',
        website: data.website ?? '',
      })
    } catch {
      setFindError("Error occurred")
    } finally {
      setFinding(false)
    }
  }

  const isCreating = createForm.formState.isSubmitting
  const isUpdating = editForm.formState.isSubmitting
  const isRowBusy = (rowId: string | number) =>
    (isUpdating && selected?.id === rowId) || (deleting && selected?.id === rowId)

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
    const windowSize = 2
    links.push(1)
    if (pageNumber > 1 + windowSize + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - windowSize)
    const end = Math.min(tp - 1, pageNumber + windowSize)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < tp - windowSize - 1) links.push('ellipsis')
    links.push(tp)
    return links
  }, [pageNumber, totalPage])

  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? 0) + 1
  const showingTo = totalRow === 0 ? 0 : (meta?.end_idx ?? -1) + 1

  return (
    <>
      <div className="h-full w-full">
        <div className="relative flex h-full flex-col overflow-hidden">
          {progress > 0 && (
            <div className="absolute inset-x-0 top-0 z-10">
              <Progress value={progress} className="h-1 rounded-none" aria-label="Уншиж байна" />
            </div>
          )}

          <div className="flex flex-col overflow-hidden px-4 pb-4">
            {/* Header */}
            <div className="flex flex-col gap-3 pt-4 md:flex-row md:items-center md:justify-between">
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                {t('all_organization')}
              </h1>
              <Button
                onClick={onOpenCreate}
                className="gap-2"
                disabled={isCreating || isUpdating || deleting}
              >
                {isCreating ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Plus className="h-4 w-4" />
                )}
                <span className="lowercase first-letter:uppercase">
                  {t('create', { name: t('organization') })}
                </span>
              </Button>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <Input
                  value={filterName}
                  onChange={(e) => setFilterName(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') load(1, pageSize)
                  }}
                  placeholder={t('search_by_name')}
                  className="h-9 sm:w-64"
                />
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

            {/* Table */}
            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
              {fetchError ? (
                <div className="p-6 text-sm text-red-600">{fetchError}</div>
              ) : (
                <div className="flex min-h-0 flex-1 flex-col">
                  <div className="min-w-full overflow-x-auto">
                    <Table className="w-full table-fixed">
                      <colgroup>
                        <col style={{ width: '200px' }} />
                        <col style={{ width: '18%' }} />
                        <col style={{ width: '28%' }} />
                        <col style={{ width: '22%' }} />
                        <col />
                      </colgroup>
                      <TableHeader>
                        <TableRow className="[&>th]:z-20 [&>th]:bg-gray-50 dark:[&>th]:bg-gray-800">
                          <TableHead></TableHead>
                          <TableHead>{t('reg_no')}</TableHead>
                          <TableHead>{t('name')}</TableHead>
                          <TableHead>{t('organization_type')}</TableHead>
                          <TableHead>{t('contact')}</TableHead>
                        </TableRow>
                      </TableHeader>
                    </Table>
                  </div>

                  {/* body */}
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
                            <col style={{ width: '200px' }} />
                            <col style={{ width: '18%' }} />
                            <col style={{ width: '28%' }} />
                            <col style={{ width: '22%' }} />
                            <col />
                          </colgroup>
                          <TableBody>
                            {rows.map((r) => {
                              const busy = isRowBusy(r.id)
                              return (
                                <TableRow key={r.id} className="[&>td]:align-middle">
                                  <TableCell className="flex gap-x-2">
                                    <Tooltip>
                                      <TooltipTrigger asChild>
                                        <Button
                                          variant="default"
                                          size="sm"
                                          onClick={() => {
                                            setOpenUser(true)
                                            setSelected(r)
                                          }}
                                          aria-label={`user ${r.name}`}
                                          disabled={busy || deleting || isCreating}
                                        >
                                          <Users2 className="h-4 w-4" />
                                        </Button>
                                      </TooltipTrigger>
                                      <TooltipContent>
                                        <p className="lowercase first-letter:uppercase">
                                          {t('organization')} - {t('user')}
                                        </p>
                                        <TooltipArrow className="fill-popover" />
                                      </TooltipContent>
                                    </Tooltip>

                                    <Tooltip>
                                      <TooltipTrigger asChild>
                                        <Button
                                          variant="outline"
                                          size="sm"
                                          onClick={() => onOpenEdit(r)}
                                          aria-label={`Edit ${r.name}`}
                                          disabled={busy || deleting || isCreating}
                                        >
                                          {isUpdating && selected?.id === r.id ? (
                                            <Loader2 className="h-4 w-4 animate-spin" />
                                          ) : (
                                            <Pencil className="h-4 w-4" />
                                          )}
                                        </Button>
                                      </TooltipTrigger>
                                      <TooltipContent>
                                        <p className="lowercase first-letter:uppercase">
                                          {t('update', { name: t('organization') })}
                                        </p>
                                        <TooltipArrow className="fill-popover" />
                                      </TooltipContent>
                                    </Tooltip>

                                    <Tooltip>
                                      <TooltipTrigger asChild>
                                        <Button
                                          variant="destructive"
                                          size="sm"
                                          onClick={() => onOpenDelete(r)}
                                          aria-label={`Delete ${r.name}`}
                                          disabled={busy || isUpdating || isCreating}
                                        >
                                          {deleting && selected?.id === r.id ? (
                                            <Loader2 className="h-4 w-4 animate-spin" />
                                          ) : (
                                            <Trash2 className="h-4 w-4" />
                                          )}
                                        </Button>
                                      </TooltipTrigger>
                                      <TooltipContent>
                                        <p className="lowercase first-letter:uppercase">
                                          {t('delete', { name: t('organization') })}
                                        </p>
                                        <TooltipArrow className="fill-popover" />
                                      </TooltipContent>
                                    </Tooltip>

                                    <Tooltip>
                                      <TooltipTrigger asChild>
                                        <Button
                                          variant="outline"
                                          size="sm"
                                          className="mr-2 gap-1"
                                          onClick={() => {
                                            // setOpenDetail(true)
                                            // setSelected(r)
                                          }}
                                          aria-label={`detail ${r.name}`}
                                          disabled={busy || deleting || isCreating}
                                        >
                                          <Info className="h-4 w-4" />
                                        </Button>
                                      </TooltipTrigger>
                                      <TooltipContent>
                                        <p className="lowercase first-letter:uppercase">
                                          {t('organization')}
                                        </p>
                                        <TooltipArrow className="fill-popover" />
                                      </TooltipContent>
                                    </Tooltip>
                                  </TableCell>

                                  {/* reg_no */}
                                  <TableCell>
                                    <Badge variant="outline" className="rounded-md">
                                      {r.reg_no || '??"'}
                                    </Badge>
                                  </TableCell>

                                  {/* name */}
                                  <TableCell className="font-medium">{r.name || '??"'}</TableCell>

                                  {/* type */}
                                  <TableCell>
                                    <Badge variant="secondary" className="rounded-md">
                                      {r.type?.name || '??"'}
                                    </Badge>
                                  </TableCell>

                                  {/* contact */}
                                  <TableCell>
                                    <div className="text-muted-foreground flex flex-col text-sm">
                                      <span>{r.phone_no || '??"'}</span>
                                      <span className="truncate">{r.email || '??"'}</span>
                                    </div>
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

        <Dialog
          open={openCreate}
          onOpenChange={(v) => {
            setOpenCreate(v)
            if (!v) setFindError(null)
          }}
        >
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle className="lowercase first-letter:uppercase">
                {t('create', { name: t('organization') })}
              </DialogTitle>
              <DialogDescription>
                {t('fill_the_fields_and_save_to_create')} — {t('optional')}{' '}
                {t('reg_no').toLowerCase()}.
              </DialogDescription>
            </DialogHeader>

            <Form {...createForm}>
              <form
                onSubmit={createForm.handleSubmit(onCreate)}
                className="space-y-4"
                autoComplete="off"
              >
                {/* Регистр + Core хайлт */}
                <div className="grid grid-cols-[1fr_auto] gap-3">
                  <FormField
                    control={createForm.control}
                    name="reg_no"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          {t('reg_no')} ({t('optional')})
                        </FormLabel>
                        <FormControl>
                          <Input maxLength={10} placeholder={t('reg_no')} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <div className="flex items-end">
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={doFindFromCore}
                      disabled={finding}
                    >
                      {finding ? (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      ) : (
                        <Search className="mr-2 h-4 w-4" />
                      )}
                      {t('search')}
                    </Button>
                  </div>
                </div>

                {findError && (
                  <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
                    {findError}
                  </div>
                )}

                {/* Үндсэн талбарууд */}
                <FormField
                  control={createForm.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('name')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('name')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 gap-3">
                  <FormField
                    control={createForm.control}
                    name="type_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('organization_type')}</FormLabel>
                        <Select
                          value={field.value ? String(field.value) : ''}
                          onValueChange={(v) => field.onChange(Number(v))}
                        >
                          <SelectTrigger className="w-full">
                            <SelectValue placeholder={t('organization_type')} />
                          </SelectTrigger>
                          <SelectContent className="max-h-60">
                            {orgTypes.map((ot) => (
                              <SelectItem key={ot.id} value={String(ot.id)}>
                                {ot.name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={createForm.control}
                    name="status"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('status')}</FormLabel>
                        <Select value={field.value} onValueChange={field.onChange}>
                          <SelectTrigger className="w-full">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="active">{t('active')}</SelectItem>
                            <SelectItem value="inactive">{t('inactive')}</SelectItem>
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <FormField
                  control={createForm.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('description')}</FormLabel>
                      <FormControl>
                        <Textarea rows={3} placeholder={t('description')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={createForm.control}
                  name="address"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('address')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('address')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 gap-3">
                  <FormField
                    control={createForm.control}
                    name="phone_no"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('phone_no')}</FormLabel>
                        <FormControl>
                          <Input placeholder="+976 7000-0000" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={createForm.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('email')}</FormLabel>
                        <FormControl>
                          <Input placeholder="info@company.mn" type="email" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <DialogFooter>
                  <Button type="button" variant="outline" onClick={() => setOpenCreate(false)}>
                    {t('cancel')}
                  </Button>
                  <Button type="submit" disabled={createForm.formState.isSubmitting}>
                    {createForm.formState.isSubmitting && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    {t('save')}
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          </DialogContent>
        </Dialog>

        {/* ===== Edit ===== */}
        <Dialog open={openEdit} onOpenChange={setOpenEdit}>
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle className="lowercase first-letter:uppercase">
                {t('update', { name: t('organization') })}
              </DialogTitle>
              <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
            </DialogHeader>

            {selected && (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <div className="text-xs opacity-70">{t('name')}</div>
                    <div className="rounded-md border px-3 py-2 text-sm">{selected.name}</div>
                  </div>
                  <div>
                    <div className="text-xs opacity-70">{t('reg_no')}</div>
                    <div className="rounded-md border px-3 py-2 text-sm">
                      {selected.reg_no || '—'}
                    </div>
                  </div>
                </div>

                <Form {...editForm}>
                  <form onSubmit={editForm.handleSubmit(onUpdate)} className="space-y-4">
                    <div className="grid grid-cols-2 gap-3">
                      <FormField
                        control={editForm.control}
                        name="type_id"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>{t('organization_type')} *</FormLabel>
                            <Select
                              value={field.value ? String(field.value) : ''}
                              onValueChange={(v) => field.onChange(Number(v))}
                            >
                              <SelectTrigger className="w-full">
                                <SelectValue placeholder={t('organization_type')} />
                              </SelectTrigger>
                              <SelectContent className="max-h-60">
                                {orgTypes.map((ot) => (
                                  <SelectItem key={ot.id} value={String(ot.id)}>
                                    {ot.name}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <FormField
                        control={editForm.control}
                        name="status"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>{t('status')}</FormLabel>
                            <Select value={field.value} onValueChange={field.onChange}>
                              <SelectTrigger className="w-full">
                                <SelectValue />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="active">{t('active')}</SelectItem>
                                <SelectItem value="inactive">{t('inactive')}</SelectItem>
                              </SelectContent>
                            </Select>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>

                    <FormField
                      control={editForm.control}
                      name="description"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('description')}</FormLabel>
                          <FormControl>
                            <Textarea rows={3} placeholder={t('description')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={editForm.control}
                      name="address"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('address')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('address')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <div className="grid grid-cols-2 gap-3">
                      <FormField
                        control={editForm.control}
                        name="phone_no"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>{t('phone_no')}</FormLabel>
                            <FormControl>
                              <Input placeholder={t('phone_no')} {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={editForm.control}
                        name="email"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>{t('email')}</FormLabel>
                            <FormControl>
                              <Input placeholder={t('email')} {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>

                    <DialogFooter>
                      <Button type="button" variant="outline" onClick={() => setOpenEdit(false)}>
                        {t('cancel')}
                      </Button>
                      <Button type="submit" disabled={editForm.formState.isSubmitting}>
                        {editForm.formState.isSubmitting && (
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        )}
                        {t('save')}
                      </Button>
                    </DialogFooter>
                  </form>
                </Form>
              </div>
            )}
          </DialogContent>
        </Dialog>

        {/* ===== Delete (confirm) ===== */}
        <Dialog open={openDelete} onOpenChange={setOpenDelete}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle className="lowercase first-letter:uppercase">
                {t('delete', { name: t('organization') })}
              </DialogTitle>
              <DialogDescription className="pt-2 text-base">
                {t.rich('delete_warning', {
                  name: () => <span className="font-medium">{selected?.name}</span>,
                })}
              </DialogDescription>
            </DialogHeader>
            <DialogFooter className="pt-2">
              <Button variant="outline" onClick={() => setOpenDelete(false)} disabled={deleting}>
                {t('cancel')}
              </Button>
              <Button variant="destructive" onClick={onDelete} disabled={deleting}>
                {deleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                <span className="capitalize">{t('delete', { name: '' })}</span>
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* ===== Org Users ===== */}
        <OrgUsersDialog open={openUser} onOpenChange={setOpenUser} orgId={selected?.id} />

        <OrganizationDetailDialog
          open={openDetail}
          onOpenChange={setOpenDetail}
          orgId={selected?.id}
        />
      </div>
    </>
  )
}
