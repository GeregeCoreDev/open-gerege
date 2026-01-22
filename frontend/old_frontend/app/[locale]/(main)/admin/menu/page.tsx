/**
 * 🖥️ Menu Page (/[locale]/(main)/admin/menu/page.tsx)
 *
 * Энэ нь цэс удирдах админ хуудас юм.
 * Зорилго: Мод бүтэцтэй цэсүүдийн CRUD удирдлага
 *
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Tree structure display
 * - ✅ Expand/Collapse nodes
 * - ✅ Search filter
 * - ✅ Progress bar loading
 * - ✅ Icon selection (Lucide icons)
 * - ✅ Active/Inactive toggle
 * - ✅ Sequence ordering
 * - ✅ Parent-Child relationship
 * - ✅ Form validation (Zod)
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React, { useEffect, useMemo, useState, useCallback } from 'react'
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
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import {
  Plus,
  Pencil,
  Trash2,
  Loader2,
  ChevronRight,
  ChevronDown,
  Check,
  ChevronsUpDown,
} from 'lucide-react'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { LucideIcon } from '@/lib/utils/icon'
import { cn } from '@/lib/utils'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { TooltipArrow } from '@radix-ui/react-tooltip'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'

const MenuSchema = z.object({
  key: z.string().min(1, 'Required').max(255),
  code: z.string().min(1, 'Required').max(255),
  name: z.string().min(1, 'Required').max(255),
  description: z.string().max(255).optional().or(z.literal('')),
  icon: z.string().max(255).optional().or(z.literal('')),
  path: z.string().max(255).optional().or(z.literal('')),
  sequence: z.coerce.number().int().default(0),
  parent_id: z.coerce.number().int().optional().nullable(),
  is_active: z.boolean().nullable().optional(),
  permission_id: z.coerce.number().int().optional().nullable(),
})

type MenuForm = z.input<typeof MenuSchema>

export default function MenuPage() {
  const t = useTranslations()

  const [tree, setTree] = useState<App.Menu[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)
  const [expanded, setExpanded] = useState<Set<number>>(new Set())
  const [permissions, setPermissions] = useState<App.Permission[]>([])
  const [_openPermissionCreate, _setOpenPermissionCreate] = useState(false)
  const [openPermissionCreateChild, setOpenPermissionCreateChild] = useState(false)
  const [openPermissionEdit, setOpenPermissionEdit] = useState(false)

  // progress bar
  const progress = useLoadingProgress(loading)

  // filters
  const [filterName, _setFilterName] = useState('')

  // Build flat list from tree for filtering
  const flattenTree = useCallback((nodes: App.Menu[]): App.Menu[] => {
    const result: App.Menu[] = []
    const traverse = (items: App.Menu[]) => {
      for (const item of items) {
        result.push(item)
        if (item.children && item.children.length > 0) {
          traverse(item.children)
        }
      }
    }
    traverse(nodes)
    return result
  }, [])

  // Filter tree based on search
  const filteredTree = useMemo(() => {
    if (!filterName) return tree
    const flat = flattenTree(tree)
    const matchedIds = new Set(
      flat
        .filter(
          (m) =>
            m.name?.toLowerCase().includes(filterName.toLowerCase()) ||
            m.key?.toLowerCase().includes(filterName.toLowerCase()) ||
            m.code?.toLowerCase().includes(filterName.toLowerCase()),
        )
        .map((m) => m.id),
    )
    const filterNode = (nodes: App.Menu[]): App.Menu[] => {
      return nodes
        .filter((n) => matchedIds.has(n.id) || (n.children && n.children.length > 0))
        .map((n) => ({
          ...n,
          children: n.children ? filterNode(n.children) : undefined,
        }))
    }
    return filterNode(tree)
  }, [tree, filterName, flattenTree])

  async function load() {
    setLoading(true)
    setFetchError(null)
    try {
      const data = await api.get<App.Menu[]>('/menu')
      setTree(data)
    } catch {
      setFetchError("Error occurred")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  // Load permissions for dropdown
  useEffect(() => {
    async function loadPermissions() {
      try {
        const data = await api.get<App.ListData<App.Permission>>('/permission', {
          query: { page: 1, size: 500 },
          cache: 'no-store',
        })
        setPermissions(data.items ?? [])
      } catch {
        setPermissions([])
      }
    }
    loadPermissions()
  }, [])

  // Expand all nodes by default when tree is loaded
  useEffect(() => {
    if (tree.length > 0) {
      const allNodeIds = new Set<number>()
      const traverse = (nodes: App.Menu[]) => {
        for (const node of nodes) {
          if (node.children && node.children.length > 0) {
            allNodeIds.add(node.id)
            traverse(node.children)
          }
        }
      }
      traverse(tree)
      setExpanded(allNodeIds)
    }
  }, [tree])

  // dialogs
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [openCreateChild, setOpenCreateChild] = useState(false)
  const [selected, setSelected] = useState<App.Menu | null>(null)
  const [parentMenu, setParentMenu] = useState<App.Menu | null>(null)

  // forms
  const createForm = useForm<MenuForm>({
    resolver: zodResolver(MenuSchema),
    defaultValues: {
      key: '',
      code: '',
      name: '',
      description: '',
      icon: '',
      path: '',
      sequence: 0,
      parent_id: null,
      is_active: true,
      permission_id: null,
    },
  })
  const editForm = useForm<MenuForm>({
    resolver: zodResolver(MenuSchema),
    defaultValues: {
      key: '',
      code: '',
      name: '',
      description: '',
      icon: '',
      path: '',
      sequence: 0,
      parent_id: null,
      is_active: true,
      permission_id: null,
    },
  })
  const createChildForm = useForm<MenuForm>({
    resolver: zodResolver(MenuSchema),
    defaultValues: {
      key: '',
      code: '',
      name: '',
      description: '',
      icon: '',
      path: '',
      sequence: 0,
      parent_id: 0,
      is_active: true,
      permission_id: 0,
    },
  })

  const onOpenEdit = (row: App.Menu) => {
    setSelected(row)
    editForm.reset({
      key: row.key ?? '',
      code: row.code ?? '',
      name: row.name ?? '',
      description: row.description ?? '',
      icon: row.icon ?? '',
      path: row.path ?? '',
      sequence: row.sequence ?? 0,
      parent_id: row.parent_id ?? 0,
      is_active: row.is_active ?? false,
      permission_id: row.permission_id ?? 0,
    })

    setOpenEdit(true)
  }

  const onOpenDelete = (row: App.Menu) => {
    setSelected(row)
    setOpenDelete(true)
  }

  const onOpenCreateChild = (parent: App.Menu) => {
    setParentMenu(parent)
    createChildForm.reset({
      key: '',
      code: '',
      name: '',
      description: '',
      icon: '',
      path: '',
      sequence: 0,
      parent_id: parent.id ?? null,
      is_active: true,
      permission_id: null,
    })
    setOpenCreateChild(true)
  }

  // create / update / delete
  const onCreate: SubmitHandler<MenuForm> = async (valuesIn) => {
    const values = MenuSchema.parse(valuesIn)
    try {
      const payload = {
        key: values.key,
        code: values.code,
        name: values.name,
        description: values.description,
        icon: values.icon,
        path: values.path,
        sequence: values.sequence ?? 0,
        parent_id: values.parent_id ?? 0,
        is_active: values.is_active,
        permission_id: values.permission_id ?? 0,
      }
      await api.post<App.Menu>('/menu', payload as Record<string, unknown>)
      setOpenCreate(false)
      createForm.reset({
        key: '',
        code: '',
        name: '',
        description: '',
        icon: '',
        path: '',
        sequence: 0,
        parent_id: null,
        is_active: true,
        permission_id: null,
      })
      await load()
    } catch {}
  }

  const onCreateChild: SubmitHandler<MenuForm> = async (valuesIn) => {
    const values = MenuSchema.parse(valuesIn)
    try {
      const payload = {
        key: values.key,
        code: values.code,
        name: values.name,
        description: values.description,
        icon: values.icon,
        path: values.path,
        sequence: values.sequence ?? 0,
        parent_id: values.parent_id ?? 0,
        is_active: values.is_active,
        permission_id: values.permission_id,
      }
      await api.post<App.Menu>('/menu', payload as Record<string, unknown>)
      setOpenCreateChild(false)
      setParentMenu(null)
      createChildForm.reset({
        key: '',
        code: '',
        name: '',
        description: '',
        icon: '',
        path: '',
        sequence: 0,
        parent_id: null,
        is_active: true,
        permission_id: null,
      })
      await load()
    } catch {}
  }

  const onUpdate: SubmitHandler<MenuForm> = async (valuesIn) => {
    if (!selected) return
    const values = MenuSchema.parse(valuesIn)
    try {
      const payload = {
        id: selected.id,
        key: values.key,
        code: values.code,
        name: values.name,
        description: values.description,
        icon: values.icon,
        path: values.path,
        sequence: values.sequence ?? 0,
        parent_id: values.parent_id ?? 0,
        is_active: values.is_active,
        permission_id: values.permission_id ?? 0,
      }
      await api.put<App.Menu>(`/menu/${selected.id}`, payload as Record<string, unknown>)
      setOpenEdit(false)
      setSelected(null)
      await load()
    } catch {}
  }

  const onDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>(`/menu/${selected.id}`)
      setOpenDelete(false)
      setSelected(null)
      await load()
    } catch {
    } finally {
      setDeleting(false)
    }
  }

  const toggleExpand = (id: number) => {
    setExpanded((prev) => {
      const n = new Set(prev)
      if (n.has(id)) {
        n.delete(id)
      } else {
        n.add(id)
      }
      return n
    })
  }

  // helpers
  const isCreating = createForm.formState.isSubmitting
  const isUpdating = editForm.formState.isSubmitting
  const _isCreatingChild = createChildForm.formState.isSubmitting
  const isRowBusy = (rowId: number) =>
    (isUpdating && selected?.id === rowId) || (deleting && selected?.id === rowId)

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
            <div className="flex flex-col gap-4 py-4 md:flex-row md:items-center md:justify-between">
              <div>
                <h1 className="text-foreground text-2xl font-bold tracking-tight">{t('menu')}</h1>
                <p className="text-muted-foreground mt-1 text-sm">
                  {t('menu_management_description')}
                </p>
              </div>
            </div>

            {/* Content */}
            <Card className="flex min-h-0 flex-1 flex-col overflow-hidden border">
              {loading ? (
                <div className="flex h-32 w-full flex-col items-center justify-center gap-y-6">
                  <Loader2 className="text-muted-foreground h-6 w-6 animate-spin" />
                  <p className="text-muted-foreground">{t('loading')}</p>
                </div>
              ) : fetchError ? (
                <CardContent>
                  <div className="text-destructive bg-destructive/10 border-destructive/20 rounded-lg border p-8 text-sm">
                    {fetchError}
                  </div>
                </CardContent>
              ) : filteredTree.length === 0 ? (
                <CardContent>
                  <div className="text-muted-foreground flex h-48 w-full flex-col items-center justify-center gap-y-6">
                    <LucideIcon name="i-lucide-archive-x" className="h-12 w-12 opacity-50" />
                    <p className="text-muted-foreground">
                      {filterName ? t('no_results_found') : t('no_information_available')}
                    </p>
                  </div>
                </CardContent>
              ) : (
                <CardContent className="min-h-0 flex-1 overflow-auto p-4">
                  <ul className="space-y-1">
                    {filteredTree
                      .sort((a, b) => (a.sequence ?? 0) - (b.sequence ?? 0))
                      .map((node) => (
                        <MenuTreeRow
                          key={node.id}
                          node={node}
                          level={0}
                          expanded={expanded}
                          onToggle={toggleExpand}
                          onAddChild={onOpenCreateChild}
                          onEdit={onOpenEdit}
                          onDelete={onOpenDelete}
                          isRowBusy={isRowBusy}
                          isUpdating={isUpdating}
                          deleting={deleting}
                          isCreating={isCreating}
                          selected={selected}
                          t={t}
                          ActiveBadge={ActiveBadge}
                        />
                      ))}
                  </ul>
                </CardContent>
              )}
            </Card>
          </div>

          {/* ---------- Create Dialog ---------- */}
          <Dialog open={openCreate} onOpenChange={setOpenCreate}>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle className="lowercase first-letter:uppercase">
                  {t('create', { name: t('menu') })}
                </DialogTitle>
                <DialogDescription>{t('fill_the_fields_and_save_to_create')}</DialogDescription>
              </DialogHeader>

              <Form {...createForm}>
                <form
                  onSubmit={createForm.handleSubmit(onCreate)}
                  className="space-y-4 pt-2"
                  autoComplete="off"
                >
                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormField
                      control={createForm.control}
                      name="key"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('translation_key')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('translation_key')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createForm.control}
                      name="code"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('code')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('code')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
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
                    <FormField
                      control={createForm.control}
                      name="icon"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('icon')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('icon')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createForm.control}
                      name="path"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('path')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('path')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createForm.control}
                      name="sequence"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('sequence')}</FormLabel>
                          <FormControl>
                            <Input
                              type="number"
                              placeholder={t('sequence')}
                              value={
                                field.value === null || field.value === undefined
                                  ? ''
                                  : (field.value as number | string)
                              }
                              onChange={(e) => {
                                const v = e.target.value
                                field.onChange(v === '' ? undefined : Number(v))
                              }}
                              onBlur={field.onBlur}
                              name={field.name}
                              ref={field.ref}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createForm.control}
                      name="is_active"
                      render={({ field }) => (
                        <FormItem className="w-full">
                          <FormLabel>{t('status')}</FormLabel>
                          <FormControl className="w-full">
                            <Select
                              value={
                                field.value === null || field.value === undefined
                                  ? 'true'
                                  : field.value
                                    ? 'true'
                                    : 'false'
                              }
                              onValueChange={(v) => field.onChange(v === 'true')}
                            >
                              <SelectTrigger className="h-9 w-full">
                                <SelectValue placeholder={t('status')} />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="true">{t('active')}</SelectItem>
                                <SelectItem value="false">{t('inactive')}</SelectItem>
                              </SelectContent>
                            </Select>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createForm.control}
                      name="permission_id"
                      render={({ field }) => {
                        const selectedPermission = permissions.find((p) => p.id === field.value)
                        return (
                          <FormItem>
                            <FormLabel>{t('permission')}</FormLabel>
                            <FormControl>
                              <Popover>
                                <PopoverTrigger asChild>
                                  <Button
                                    variant="outline"
                                    role="combobox"
                                    className={cn(
                                      'h-9 w-full justify-between',
                                      !field.value && 'text-muted-foreground',
                                    )}
                                  >
                                    <div className="flex min-w-0 flex-1 items-center gap-2">
                                      {selectedPermission ? (
                                        <div className="flex min-w-0 flex-col truncate text-left">
                                          <span className="truncate text-sm">
                                            {selectedPermission.name}
                                          </span>
                                          <span className="text-muted-foreground truncate text-xs">
                                            {selectedPermission.code}
                                          </span>
                                        </div>
                                      ) : (
                                        <span className="text-muted-foreground">
                                          {t('optional')}
                                        </span>
                                      )}
                                    </div>
                                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                  </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-[400px] p-0" align="start">
                                  <Command>
                                    <CommandInput
                                      placeholder={t('search_permission') || 'Search permission...'}
                                    />
                                    <CommandList>
                                      <CommandEmpty>
                                        {t('no_permission_found') || 'No permission found.'}
                                      </CommandEmpty>
                                      <CommandGroup className="">
                                        <CommandItem
                                          value="none"
                                          onSelect={() => {
                                            field.onChange(null)
                                          }}
                                        >
                                          <Check
                                            className={cn(
                                              'mr-2 h-4 w-4',
                                              !field.value ? 'opacity-100' : 'opacity-0',
                                            )}
                                          />
                                          {t('optional')}
                                        </CommandItem>
                                        {permissions.map((perm) => (
                                          <CommandItem
                                            key={perm.id}
                                            value={`${perm.name} ${perm.code}`}
                                            onSelect={() => {
                                              field.onChange(perm.id)
                                            }}
                                          >
                                            <Check
                                              className={cn(
                                                'mr-2 h-4 w-4',
                                                field.value === perm.id
                                                  ? 'opacity-100'
                                                  : 'opacity-0',
                                              )}
                                            />
                                            <div className="flex min-w-0 flex-col">
                                              <span className="truncate text-sm">{perm.name}</span>
                                              <span className="text-muted-foreground truncate text-xs">
                                                {perm.code}
                                              </span>
                                            </div>
                                          </CommandItem>
                                        ))}
                                      </CommandGroup>
                                    </CommandList>
                                  </Command>
                                </PopoverContent>
                              </Popover>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )
                      }}
                    />
                  </div>

                  <FormField
                    control={createForm.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('description')}</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder={`${t('description')}, ${t('optional')}`}
                            rows={4}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

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

          {/* ---------- Create Child Dialog ---------- */}
          <Dialog open={openCreateChild} onOpenChange={setOpenCreateChild}>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle className="lowercase first-letter:uppercase">
                  {t('create', { name: t('menu') })}
                </DialogTitle>
                <DialogDescription>
                  {t('create_child_menu')} - {parentMenu?.name}
                </DialogDescription>
              </DialogHeader>

              <Form {...createChildForm}>
                <form
                  onSubmit={createChildForm.handleSubmit(onCreateChild)}
                  className="space-y-4 pt-2"
                  autoComplete="off"
                >
                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormField
                      control={createChildForm.control}
                      name="key"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('translation_key')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('translation_key')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createChildForm.control}
                      name="code"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('code')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('code')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createChildForm.control}
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
                    <FormField
                      control={createChildForm.control}
                      name="icon"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('icon')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('icon')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createChildForm.control}
                      name="path"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('path')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('path')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createChildForm.control}
                      name="sequence"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('sequence')}</FormLabel>
                          <FormControl>
                            <Input
                              type="number"
                              placeholder={t('sequence')}
                              value={
                                field.value === null || field.value === undefined
                                  ? ''
                                  : (field.value as number | string)
                              }
                              onChange={(e) => {
                                const v = e.target.value
                                field.onChange(v === '' ? undefined : Number(v))
                              }}
                              onBlur={field.onBlur}
                              name={field.name}
                              ref={field.ref}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createChildForm.control}
                      name="is_active"
                      render={({ field }) => (
                        <FormItem className="w-full">
                          <FormLabel>{t('status')}</FormLabel>
                          <FormControl className="w-full">
                            <Select
                              value={
                                field.value === null || field.value === undefined
                                  ? 'true'
                                  : field.value
                                    ? 'true'
                                    : 'false'
                              }
                              onValueChange={(v) => field.onChange(v === 'true')}
                            >
                              <SelectTrigger className="h-9 w-full">
                                <SelectValue placeholder={t('status')} />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="true">{t('active')}</SelectItem>
                                <SelectItem value="false">{t('inactive')}</SelectItem>
                              </SelectContent>
                            </Select>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={createChildForm.control}
                      name="permission_id"
                      render={({ field }) => {
                        const selectedPermission = permissions.find((p) => p.id === field.value)
                        return (
                          <FormItem>
                            <FormLabel>{t('permission')}</FormLabel>
                            <FormControl>
                              <Popover
                                open={openPermissionCreateChild}
                                onOpenChange={setOpenPermissionCreateChild}
                              >
                                <PopoverTrigger asChild>
                                  <Button
                                    variant="outline"
                                    role="combobox"
                                    className={cn(
                                      'h-9 w-full justify-between',
                                      !field.value && 'text-muted-foreground',
                                    )}
                                  >
                                    <div className="flex min-w-0 flex-1 items-center gap-2">
                                      {selectedPermission ? (
                                        <div className="flex min-w-0 flex-col truncate text-left">
                                          <span className="truncate text-sm">
                                            {selectedPermission.name}
                                          </span>
                                          <span className="text-muted-foreground truncate text-xs">
                                            {selectedPermission.code}
                                          </span>
                                        </div>
                                      ) : (
                                        <span className="text-muted-foreground">
                                          {t('optional')}
                                        </span>
                                      )}
                                    </div>
                                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                  </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-[400px] p-0" align="start">
                                  <Command>
                                    <CommandInput
                                      placeholder={t('search_permission') || 'Search permission...'}
                                    />
                                    <CommandList>
                                      <CommandEmpty>
                                        {t('no_permission_found') || 'No permission found.'}
                                      </CommandEmpty>
                                      <CommandGroup>
                                        <CommandItem
                                          value="none"
                                          onSelect={() => {
                                            field.onChange(null)
                                            setOpenPermissionCreateChild(false)
                                          }}
                                        >
                                          <Check
                                            className={cn(
                                              'mr-2 h-4 w-4',
                                              !field.value ? 'opacity-100' : 'opacity-0',
                                            )}
                                          />
                                          {t('optional')}
                                        </CommandItem>
                                        {permissions.map((perm) => (
                                          <CommandItem
                                            key={perm.id}
                                            value={`${perm.name} ${perm.code}`}
                                            onSelect={() => {
                                              field.onChange(perm.id)
                                              setOpenPermissionCreateChild(false)
                                            }}
                                          >
                                            <Check
                                              className={cn(
                                                'mr-2 h-4 w-4',
                                                field.value === perm.id
                                                  ? 'opacity-100'
                                                  : 'opacity-0',
                                              )}
                                            />
                                            <div className="flex min-w-0 flex-col">
                                              <span className="truncate text-sm">{perm.name}</span>
                                              <span className="text-muted-foreground truncate text-xs">
                                                {perm.code}
                                              </span>
                                            </div>
                                          </CommandItem>
                                        ))}
                                      </CommandGroup>
                                    </CommandList>
                                  </Command>
                                </PopoverContent>
                              </Popover>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )
                      }}
                    />
                  </div>

                  <FormField
                    control={createChildForm.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('description')}</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder={`${t('description')}, ${t('optional')}`}
                            rows={4}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => {
                        setOpenCreateChild(false)
                        setParentMenu(null)
                      }}
                    >
                      {t('cancel')}
                    </Button>
                    <Button type="submit" disabled={createChildForm.formState.isSubmitting}>
                      {createChildForm.formState.isSubmitting && (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      {t('save')}
                    </Button>
                  </DialogFooter>
                </form>
              </Form>
            </DialogContent>
          </Dialog>

          {/* ---------- Edit Dialog ---------- */}
          <Dialog open={openEdit} onOpenChange={setOpenEdit}>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle className="lowercase first-letter:uppercase">
                  {t('update', { name: t('menu') })}
                </DialogTitle>
                <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
              </DialogHeader>

              <Form {...editForm}>
                <form
                  onSubmit={editForm.handleSubmit(onUpdate)}
                  className="space-y-4 pt-2"
                  autoComplete="off"
                >
                  <div className="grid w-full gap-4 sm:grid-cols-2">
                    <FormField
                      control={editForm.control}
                      name="key"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('translation_key')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('translation_key')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={editForm.control}
                      name="code"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('code')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('code')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={editForm.control}
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
                    <FormField
                      control={editForm.control}
                      name="icon"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('icon')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('icon')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={editForm.control}
                      name="path"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('path')}</FormLabel>
                          <FormControl>
                            <Input placeholder={t('path')} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={editForm.control}
                      name="sequence"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('sequence')}</FormLabel>
                          <FormControl>
                            <Input
                              type="number"
                              placeholder={t('sequence')}
                              value={
                                field.value === null || field.value === undefined
                                  ? ''
                                  : (field.value as number | string)
                              }
                              onChange={(e) => {
                                const v = e.target.value
                                field.onChange(v === '' ? undefined : Number(v))
                              }}
                              onBlur={field.onBlur}
                              name={field.name}
                              ref={field.ref}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={editForm.control}
                      name="is_active"
                      render={({ field }) => (
                        <FormItem className="w-full">
                          <FormLabel>{t('status')}</FormLabel>
                          <FormControl className="w-full">
                            <Select
                              value={
                                field.value === null || field.value === undefined
                                  ? 'null'
                                  : field.value
                                    ? 'true'
                                    : 'false'
                              }
                              onValueChange={(v) =>
                                field.onChange(v === 'null' ? null : v === 'true' ? true : false)
                              }
                            >
                              <SelectTrigger className="h-9 w-full">
                                <SelectValue placeholder={t('status')} />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="true">{t('active')}</SelectItem>
                                <SelectItem value="false">{t('inactive')}</SelectItem>
                              </SelectContent>
                            </Select>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={editForm.control}
                      name="permission_id"
                      render={({ field }) => {
                        const selectedPermission = permissions.find((p) => p.id === field.value)
                        return (
                          <FormItem>
                            <FormLabel>{t('permission')}</FormLabel>
                            <FormControl>
                              <Popover
                                open={openPermissionEdit}
                                onOpenChange={setOpenPermissionEdit}
                              >
                                <PopoverTrigger asChild>
                                  <Button
                                    variant="outline"
                                    role="combobox"
                                    className={cn(
                                      'h-9 w-full justify-between',
                                      !field.value && 'text-muted-foreground',
                                    )}
                                  >
                                    <div className="flex min-w-0 flex-1 items-center gap-2">
                                      {selectedPermission ? (
                                        <div className="flex min-w-0 flex-col truncate text-left">
                                          <span className="truncate text-sm">
                                            {selectedPermission.name}
                                          </span>
                                          <span className="text-muted-foreground truncate text-xs">
                                            {selectedPermission.code}
                                          </span>
                                        </div>
                                      ) : (
                                        <span className="text-muted-foreground">
                                          {t('optional')}
                                        </span>
                                      )}
                                    </div>
                                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                  </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-[400px] p-0" align="start">
                                  <Command>
                                    <CommandInput
                                      placeholder={t('search_permission') || 'Search permission...'}
                                    />
                                    <CommandList>
                                      <CommandEmpty>
                                        {t('no_permission_found') || 'No permission found.'}
                                      </CommandEmpty>
                                      <CommandGroup>
                                        <CommandItem
                                          value="none"
                                          onSelect={() => {
                                            field.onChange(null)
                                            setOpenPermissionEdit(false)
                                          }}
                                        >
                                          <Check
                                            className={cn(
                                              'mr-2 h-4 w-4',
                                              !field.value ? 'opacity-100' : 'opacity-0',
                                            )}
                                          />
                                          {t('optional')}
                                        </CommandItem>
                                        {permissions.map((perm) => (
                                          <CommandItem
                                            key={perm.id}
                                            value={`${perm.name} ${perm.code}`}
                                            onSelect={() => {
                                              field.onChange(perm.id)
                                              setOpenPermissionEdit(false)
                                            }}
                                          >
                                            <Check
                                              className={cn(
                                                'mr-2 h-4 w-4',
                                                field.value === perm.id
                                                  ? 'opacity-100'
                                                  : 'opacity-0',
                                              )}
                                            />
                                            <div className="flex min-w-0 flex-col">
                                              <span className="truncate text-sm">{perm.name}</span>
                                              <span className="text-muted-foreground truncate text-xs">
                                                {perm.code}
                                              </span>
                                            </div>
                                          </CommandItem>
                                        ))}
                                      </CommandGroup>
                                    </CommandList>
                                  </Command>
                                </PopoverContent>
                              </Popover>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )
                      }}
                    />
                  </div>

                  <FormField
                    control={editForm.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('description')}</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder={`${t('description')}, ${t('optional')}`}
                            rows={4}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

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
            </DialogContent>
          </Dialog>

          {/* ---------- Delete Dialog ---------- */}
          <Dialog open={openDelete} onOpenChange={setOpenDelete}>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle className="lowercase first-letter:uppercase">
                  {t('delete', { name: t('menu') })}
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
                  <p className="capitalize">{t('delete', { name: '' })}</p>
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </>
  )
}

/** ===== Tree Row Component (Recursive) ===== */
function MenuTreeRow(props: {
  node: App.Menu
  level: number
  expanded: Set<number>
  onToggle: (id: number) => void
  onAddChild: (parent: App.Menu) => void
  onEdit: (node: App.Menu) => void
  onDelete: (node: App.Menu) => void
  isRowBusy: (rowId: number) => boolean
  isUpdating: boolean
  deleting: boolean
  isCreating: boolean
  selected: App.Menu | null
  t: ReturnType<typeof useTranslations>
  ActiveBadge: ({ value }: { value: boolean | null | undefined }) => React.ReactElement
}) {
  const {
    node,
    level,
    expanded,
    onToggle,
    onAddChild,
    onEdit,
    onDelete,
    isRowBusy,
    isUpdating,
    deleting,
    isCreating,
    selected,
    t,
    ActiveBadge,
  } = props
  const hasChildren = (node.children?.length ?? 0) > 0
  const isOpen = expanded.has(node.id)
  const busy = isRowBusy(node.id)
  const levelBadgeColors = [
    'bg-primary',
    'bg-emerald-500 dark:bg-emerald-600',
    'bg-amber-500 dark:bg-amber-600',
    'bg-slate-400 dark:bg-slate-500',
  ]
  const levelColor = levelBadgeColors[level] || levelBadgeColors[levelBadgeColors.length - 1]
  const rowBg =
    level === 0
      ? 'bg-card hover:bg-muted/30 dark:hover:bg-muted/30'
      : level === 1
        ? 'bg-muted/20 hover:bg-muted/30 dark:bg-muted/30 dark:hover:bg-muted/40'
        : level === 2
          ? 'bg-muted/10 hover:bg-muted/20 dark:bg-muted/20 dark:hover:bg-muted/30'
          : 'bg-muted/5 hover:bg-muted/10 dark:bg-muted/10 dark:hover:bg-muted/20'
  return (
    <li>
      <div
        className={cn(
          'group flex items-center gap-2 rounded-lg border px-3 py-2.5',
          'border-border/60 hover:border-border',
          'transition-all duration-150 ease-in-out',
          rowBg,
        )}
        style={{ paddingLeft: 20 + level * 12 }}
      >
        {/* level bullet */}
        <span className={cn('h-2 w-2 rounded-full opacity-80', levelColor)} aria-hidden />

        {hasChildren && (
          <button
            className="hover:bg-muted/50 flex h-6 w-6 items-center justify-center rounded transition-colors"
            onClick={() => onToggle(node.id)}
            aria-label="toggle"
          >
            {isOpen ? (
              <ChevronDown className="text-muted-foreground h-4 w-4" />
            ) : (
              <ChevronRight className="text-muted-foreground h-4 w-4" />
            )}
          </button>
        )}

        {node.icon && <LucideIcon name={node.icon} className="text-muted-foreground h-5 w-5" />}

        {node.code && level === 0 && <Badge>{node.code}</Badge>}

        <span className="text-foreground flex-1 font-medium">{t(node.key.toLowerCase())}</span>

        {node.path && (
          <span className="text-muted-foreground hidden font-mono text-xs sm:inline">
            {node.path}
          </span>
        )}

        <ActiveBadge value={node.is_active} />

        <div className="ml-auto flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100">
          {level < 2 ? (
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  size="icon"
                  variant="ghost"
                  onClick={() => onAddChild(node)}
                  aria-label="add child"
                  disabled={busy || deleting || isCreating || isUpdating}
                  className="h-8 w-8"
                >
                  <Plus className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p className="lowercase first-letter:uppercase">
                  {t('create', { name: t('child_menu') })}
                </p>
                <TooltipArrow className="fill-popover" />
              </TooltipContent>
            </Tooltip>
          ) : (
            <div className="w-8"></div>
          )}

          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                size="icon"
                variant="ghost"
                onClick={() => onEdit(node)}
                aria-label="edit"
                disabled={busy || deleting || isCreating}
                className="h-8 w-8"
              >
                {isUpdating && selected?.id === node.id ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Pencil className="h-4 w-4" />
                )}
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p className="lowercase first-letter:uppercase">{t('update', { name: t('menu') })}</p>
              <TooltipArrow className="fill-popover" />
            </TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                size="icon"
                variant="ghost"
                onClick={() => onDelete(node)}
                aria-label="delete"
                disabled={busy || isUpdating || isCreating}
                className="text-destructive hover:text-destructive hover:bg-destructive/10 h-8 w-8"
              >
                {deleting && selected?.id === node.id ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Trash2 className="h-4 w-4" />
                )}
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p className="lowercase first-letter:uppercase">{t('delete', { name: t('menu') })}</p>
              <TooltipArrow className="fill-popover" />
            </TooltipContent>
          </Tooltip>
        </div>
      </div>

      {hasChildren && isOpen ? (
        <ul className="mt-1 ml-4 space-y-1">
          {node
            .children!.slice()
            .sort((a, b) => (a.sequence ?? 0) - (b.sequence ?? 0))
            .map((child) => (
              <MenuTreeRow
                key={child.id}
                node={child}
                level={level + 1}
                expanded={expanded}
                onToggle={onToggle}
                onAddChild={onAddChild}
                onEdit={onEdit}
                onDelete={onDelete}
                isRowBusy={isRowBusy}
                isUpdating={isUpdating}
                deleting={deleting}
                isCreating={isCreating}
                selected={selected}
                t={t}
                ActiveBadge={ActiveBadge}
              />
            ))}
        </ul>
      ) : null}
    </li>
  )
}
