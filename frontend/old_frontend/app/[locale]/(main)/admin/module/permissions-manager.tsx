'use client'

import React, { useEffect, useMemo, useState } from 'react'
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
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { Plus, Pencil, Trash2, Loader2 } from 'lucide-react'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'
import { appConfig } from '@/config/app.config'

interface PermissionsManagerProps {
  module: App.Module | null
}

export function PermissionsManager({ module }: PermissionsManagerProps) {
  const t = useTranslations()

  const [permissions, setPermissions] = useState<App.Permission[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)
  const [meta, setMeta] = useState<App.ApiMeta | null>(null)

  // actions for create form
  const [actions, setActions] = useState<App.Action[]>([])

  // pagination
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  // dialogs
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selected, setSelected] = useState<App.Permission | null>(null)

  // create form
  const [selectedActionIds, setSelectedActionIds] = useState<number[]>([])

  // edit form
  const [editCode, setEditCode] = useState('')
  const [editName, setEditName] = useState('')
  const [editDescription, setEditDescription] = useState('')
  const [editActionId, setEditActionId] = useState<number | null>(null)
  const [editIsActive, setEditIsActive] = useState(true)

  const [creating, setCreating] = useState(false)
  const [updating, setUpdating] = useState(false)

  // Load actions
  const loadActions = React.useCallback(async () => {
    try {
      const data = await api.get<App.ListData<App.Action>>('/action', {
        query: { size: 50 },
        hasToast: false,
      })
      setActions(data.items ?? [])
    } catch {
      setActions([])
    }
  }, [])

  // Load permissions
  const load = React.useCallback(
    async (page = 1, size = pageSize) => {
      if (!module) return
      setLoading(true)
      setFetchError(null)
      try {
        const data = await api.get<App.ListData<App.Permission>>('/permission', {
          query: {
            page,
            size,
            module_id: module.id,
          },
          hasToast: false,
        })
        const m = data.meta
        setMeta(m)
        setPermissions(data.items ?? [])
        setPageNumber(m?.page ?? page)
        setPageSize(m?.size ?? size)
        setTotalPage(m?.pages ?? 1)
        setTotalRow(m?.total ?? 0)
      } catch (e: unknown) {
        const error = e as Error
        setFetchError(error?.message || 'Failed to load')
      } finally {
        setLoading(false)
      }
    },
    [module, pageSize]
  )

  useEffect(() => {
    if (module) {
      loadActions()
      load(1, pageSize)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [module?.id, pageSize, loadActions, load])

  // Pagination helpers
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

  if (!module) return null

  // Create
  const handleCreate = async () => {
    if (!module || selectedActionIds.length === 0) return
    try {
      setCreating(true)
      const payload = {
        system_id: module.system_id || module.system?.id,
        module_id: module.id,
        action_ids: selectedActionIds,
      }
      await api.post<App.Permission>('/permission', payload as Record<string, unknown>)
      setOpenCreate(false)
      setSelectedActionIds([])
      await load(1, pageSize)
    } catch {
    } finally {
      setCreating(false)
    }
  }

  // Update
  const onOpenEdit = (permission: App.Permission) => {
    setSelected(permission)
    setEditCode(permission.code ?? '')
    setEditName(permission.name ?? '')
    setEditDescription(permission.description ?? '')
    setEditActionId(permission.action_id ?? null)
    setEditIsActive(permission.is_active ?? true)
    setOpenEdit(true)
  }

  const handleUpdate = async () => {
    if (!selected || !module) return
    try {
      setUpdating(true)
      const payload = {
        code: editCode,
        name: editName,
        description: editDescription || undefined,
        module_id: module.id,
        system_id: module.system_id || module.system?.id,
        action_id: editActionId || undefined,
        is_active: editIsActive,
      }
      await api.put<App.Permission>(
        `/permission/${selected.id}`,
        payload as Record<string, unknown>,
      )
      setOpenEdit(false)
      setSelected(null)
      await load(pageNumber, pageSize)
    } catch {
    } finally {
      setUpdating(false)
    }
  }

  // Delete
  const onOpenDelete = (permission: App.Permission) => {
    setSelected(permission)
    setOpenDelete(true)
  }

  const handleDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>(`/permission/${selected.id}`)
      const newCount = permissions.length - 1
      const willBeEmpty = newCount <= 0 && pageNumber > 1
      setOpenDelete(false)
      setSelected(null)
      await load(willBeEmpty ? pageNumber - 1 : pageNumber, pageSize)
    } catch {
    } finally {
      setDeleting(false)
    }
  }

  // Pagination helpers
  const canPrev = meta?.has_prev ?? pageNumber > 1
  const canNext = meta?.has_next ?? pageNumber < totalPage
  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) load(target, pageSize)
  }

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

  // col widths
  const headerCols = [90, 40, 240, 240, 200, 120]
  const bodyCols = headerCols

  return (
    <div className="flex flex-col">
      {/* Header with Create Button */}
      <div className="flex items-center justify-between py-4">
        <div>
          <h3 className="text-lg font-semibold">{t('permission')}</h3>
          <p className="text-muted-foreground text-sm">{t('manage_permissions_for_module')}</p>
        </div>
        <Button onClick={() => setOpenCreate(true)} className="gap-2" disabled={loading}>
          <Plus className="h-4 w-4" />
          {t('create', { name: t('permission') })}
        </Button>
      </div>

      {/* Table */}
      <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
        {fetchError ? (
          <div className="p-6 text-sm text-red-600">{fetchError}</div>
        ) : loading ? (
          <div className="flex h-32 w-full flex-col items-center justify-center gap-y-6">
            <Loader2 className="h-4 w-4 animate-spin" />
            <p>{t('loading')}</p>
          </div>
        ) : permissions.length === 0 ? (
          <div className="text-muted-foreground flex h-48 items-center justify-center">
            {t('no_permissions_found')}
          </div>
        ) : (
          <div className="min-h-0 flex-1 overflow-auto">
            <Table className="w-full table-fixed">
              <colgroup>
                {bodyCols.map((w, i) => (
                  <col key={i} style={{ width: typeof w === 'number' ? `${w}px` : String(w) }} />
                ))}
              </colgroup>

              <TableHeader className="bg-background sticky top-0 z-10">
                <TableRow>
                  <TableHead></TableHead>
                  <TableHead>ID</TableHead>
                  <TableHead>{t('code')}</TableHead>
                  <TableHead>{t('name')}</TableHead>
                  <TableHead>{t('description')}</TableHead>
                  <TableHead>{t('is_active')}</TableHead>
                </TableRow>
              </TableHeader>

              <TableBody>
                {permissions.map((permission) => (
                  <TableRow key={permission.id}>
                    <TableCell>
                      <Button
                        variant="outline"
                        size="sm"
                        className="mr-2"
                        onClick={() => onOpenEdit(permission)}
                        disabled={deleting || creating || updating}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => onOpenDelete(permission)}
                        disabled={deleting || creating || updating}
                      >
                        {deleting && selected?.id === permission.id ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <Trash2 className="h-4 w-4" />
                        )}
                      </Button>
                    </TableCell>
                    <TableCell>{permission.id}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className="font-mono text-xs">
                        {permission.code}
                      </Badge>
                    </TableCell>
                    <TableCell>{permission.name}</TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {permission.description || '—'}
                    </TableCell>
                    <TableCell>
                      <ActiveBadge value={permission.is_active} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between gap-3 rounded-b-md border-r border-b border-l border-gray-200 px-6 py-2 dark:border-gray-800">
        <div className="text-muted-foreground w-72 text-sm">
          {t.rich('showing', {
            from: () => <span className="font-medium">{showingFrom}</span>,
            to: () => <span className="font-medium">{showingTo}</span>,
            total: () => <span className="font-medium">{totalRow}</span>,
          })}
        </div>
        <Pagination className="flex justify-end">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                onClick={() => canPrev && goPage(pageNumber - 1)}
                className={!canPrev ? 'pointer-events-none opacity-50' : ''}
              />
            </PaginationItem>
            {pageLinks.map((p, i) =>
              p === 'ellipsis' ? (
                <PaginationItem key={`e-${i}`}>
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
                onClick={() => canNext && goPage(pageNumber + 1)}
                className={!canNext ? 'pointer-events-none opacity-50' : ''}
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>

      {/* Create Dialog */}
      <Dialog open={openCreate} onOpenChange={setOpenCreate}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>{t('create', { name: t('permission') })}</DialogTitle>
            <DialogDescription>{t('select_actions_to_create_permissions')}</DialogDescription>
          </DialogHeader>

          <div className="space-y-4">
            <div>
              <label className="text-sm font-medium">{t('action')}</label>
              <div className="mt-2 max-h-64 space-y-2 overflow-auto rounded-md border p-3">
                {actions.length === 0 ? (
                  <div className="text-muted-foreground text-center text-sm">
                    {t('no_actions_found')}
                  </div>
                ) : (
                  actions.map((action) => (
                    <div key={action.id} className="flex items-center space-x-2">
                      <Checkbox
                        id={`action-${action.id}`}
                        checked={selectedActionIds.includes(action.id)}
                        onCheckedChange={(checked) => {
                          if (checked) {
                            setSelectedActionIds([...selectedActionIds, action.id])
                          } else {
                            setSelectedActionIds(selectedActionIds.filter((id) => id !== action.id))
                          }
                        }}
                      />
                      <label
                        htmlFor={`action-${action.id}`}
                        className="text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                      >
                        <div className="flex items-center gap-2">
                          <span>{action.name}</span>
                        </div>
                      </label>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setOpenCreate(false)} disabled={creating}>
              {t('cancel')}
            </Button>
            <Button onClick={handleCreate} disabled={creating || selectedActionIds.length === 0}>
              {creating && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {t('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={openEdit} onOpenChange={setOpenEdit}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>{t('update', { name: t('permission') })}</DialogTitle>
            <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
          </DialogHeader>

          <div className="space-y-4">
            <div>
              <label className="text-sm font-medium">{t('code')}</label>
              <Input
                value={editCode}
                onChange={(e) => setEditCode(e.target.value)}
                className="mt-1 font-mono"
              />
            </div>

            <div>
              <label className="text-sm font-medium">{t('name')}</label>
              <Input
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                placeholder={t('name')}
                className="mt-1"
              />
            </div>

            <div>
              <label className="text-sm font-medium">{t('description')}</label>
              <Textarea
                value={editDescription}
                onChange={(e) => setEditDescription(e.target.value)}
                placeholder={t('description')}
                className="mt-1"
                rows={3}
              />
            </div>

            <div>
              <label className="text-sm font-medium">{t('action')}</label>
              <Select
                value={editActionId ? String(editActionId) : 'none'}
                onValueChange={(v) => setEditActionId(v === 'none' ? null : Number(v))}
              >
                <SelectTrigger className="mt-1 w-full">
                  <SelectValue placeholder={t('optional')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">{t('optional')}</SelectItem>
                  {actions.map((action) => (
                    <SelectItem key={action.id} value={String(action.id)}>
                      <div className="flex items-center gap-2">
                        <Badge variant="outline" className="font-mono text-xs">
                          {action.code}
                        </Badge>
                        <span>{action.name}</span>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div>
              <label className="text-sm font-medium">{t('status')}</label>
              <Select
                value={editIsActive ? 'true' : 'false'}
                onValueChange={(v) => setEditIsActive(v === 'true')}
              >
                <SelectTrigger className="mt-1 w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="true">{t('active')}</SelectItem>
                  <SelectItem value="false">{t('inactive')}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setOpenEdit(false)} disabled={updating}>
              {t('cancel')}
            </Button>
            <Button
              onClick={handleUpdate}
              disabled={updating || !editCode.trim() || !editName.trim()}
            >
              {updating && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {t('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <Dialog open={openDelete} onOpenChange={setOpenDelete}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{t('delete', { name: t('permission') })}</DialogTitle>
            <DialogDescription>
              {t.rich('delete_warning', {
                name: () => <span className="font-medium">{selected?.name}</span>,
              })}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpenDelete(false)} disabled={deleting}>
              {t('cancel')}
            </Button>
            <Button variant="destructive" onClick={handleDelete} disabled={deleting}>
              {deleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <span className="capitalize">{t('delete', { name: '' })}</span>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
