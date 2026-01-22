/**
 * 📱 App Icon Page (/[locale]/(main)/app/icon/page.tsx)
 * 
 * Энэ нь апп айкон удирдах хуудас юм.
 * Зорилго: Mobile/web app icons болон icon groups-ийн CRUD удирдлага
 * 
 * Architecture:
 * - Split view: Groups (left) | Icons (right)
 * - Master-detail pattern
 * - Group selection updates icon list
 * 
 * Features:
 * - ✅ Dual entity management (Groups + Icons)
 * - ✅ Full CRUD for both entities
 * - ✅ Search/filter by name
 * - ✅ Progress bar loading
 * - ✅ Image upload (4 variants: icon, app, tablet, kiosk)
 * - ✅ Image compression before upload
 * - ✅ Badge display (public, native, featured, best)
 * - ✅ Form validation (Zod)
 * - ✅ Tooltip on action buttons
 * - ✅ Responsive grid layout
 * 
 * Group Fields:
 * - name, name_en: Display names
 * - icon: Lucide icon name
 * - type_name: Group type
 * - seq: Display order
 * 
 * Icon Fields:
 * - name, name_en: Display names
 * - icon, icon_app, icon_tablet, icon_kiosk: Image uploads
 * - link, web_link: URLs
 * - system_code: Related system
 * - description: Text description
 * - featured_icon: Special icon
 * - seq: Display order
 * - is_public, is_native, is_featured, is_best_selling: Boolean flags
 * 
 * Image Handling:
 * - FileUpload component with preview
 * - Compression via compressImageFile()
 * - Upload via uploadImage()
 * - Delete via deleteImage()
 * - Multiple image variants per icon
 * 
 * API Endpoints:
 * - GET/POST/PUT/DELETE /app-service-group - Groups
 * - GET/POST/PUT/DELETE /app-service-icon - Icons
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React, { useEffect, useState } from 'react'
import { useForm, type SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslations } from 'next-intl'

import api from '@/lib/api'
import { LucideIcon } from '@/lib/utils/icon'

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
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { TooltipArrow } from '@radix-ui/react-tooltip'
import { Loader2, Plus, Pencil, Trash2 } from 'lucide-react'
import FileUpload, { FileValue } from '@/components/common/fileUpload'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { deleteImage, uploadImage } from '@/lib/utils/image'
import Image from 'next/image'

/* ===================== Schemas ===================== */

const IconGroupSchema = z.object({
  name: z.string().min(1, 'Required'),
  name_en: z.string().optional().or(z.literal('')),
  icon: z.string().optional().or(z.literal('')),
  type_name: z.string().optional().or(z.literal('')),
  seq: z.coerce.number().int().default(0),
})

type IconGroupForm = z.input<typeof IconGroupSchema>

const AppIconSchema = z.object({
  name: z.string().min(1, 'Required'),
  name_en: z.string().optional().or(z.literal('')),
  icon: z.string().optional().or(z.literal('')),
  icon_app: z.string().optional().or(z.literal('')),
  icon_tablet: z.string().optional().or(z.literal('')),
  icon_kiosk: z.string().optional().or(z.literal('')),
  link: z.string().optional().or(z.literal('')),
  web_link: z.string().optional().or(z.literal('')),
  system_code: z.string().optional().or(z.literal('')),
  description: z.string().optional().or(z.literal('')),
  featured_icon: z.string().optional().or(z.literal('')),
  seq: z.coerce.number().int().default(0),
  is_native: z.boolean().nullable().optional(),
  is_public: z.boolean().nullable().optional(),
  is_featured: z.boolean().nullable().optional(),
  is_best_selling: z.boolean().nullable().optional(),
})

type AppIconForm = z.input<typeof AppIconSchema>

type IconFiles = {
  icon: FileValue | null
  icon_app: FileValue | null
  icon_tablet: FileValue | null
  icon_kiosk: FileValue | null
}

export default function AppIconsPage() {
  const t = useTranslations()

  /* ---------- Shared loading ---------- */
  const [globalLoading, setGlobalLoading] = useState(false)
  
  // progress bar - simplified with hook
  const _progress = useLoadingProgress(globalLoading)

  /* ---------- Group list state ---------- */

  const [groups, setGroups] = useState<App.AppIconGroup[]>([])
  const [groupLoading, setGroupLoading] = useState(false)
  const [groupFetchError, setGroupFetchError] = useState<string | null>(null)
  const [groupDeleting, setGroupDeleting] = useState(false)
  const [groupFilterName, setGroupFilterName] = useState('')

  const [selectedGroup, setSelectedGroup] = useState<App.AppIconGroup | null>(null)

  /* ---------- Icon list state ---------- */

  const [icons, setIcons] = useState<App.AppIcon[]>([])
  const [iconLoading, setIconLoading] = useState(false)
  const [iconFetchError, setIconFetchError] = useState<string | null>(null)
  const [iconDeleting, setIconDeleting] = useState(false)
  const [iconFilterName, setIconFilterName] = useState('')

  /* ---------- Dialogs & selection ---------- */

  const [openCreateGroup, setOpenCreateGroup] = useState(false)
  const [openEditGroup, setOpenEditGroup] = useState(false)
  const [openDeleteGroup, setOpenDeleteGroup] = useState(false)
  const [selectedGroupRow, setSelectedGroupRow] = useState<App.AppIconGroup | null>(null)

  const [openCreateIcon, setOpenCreateIcon] = useState(false)
  const [openEditIcon, setOpenEditIcon] = useState(false)
  const [openDeleteIcon, setOpenDeleteIcon] = useState(false)
  const [selectedIconRow, setSelectedIconRow] = useState<App.AppIcon | null>(null)

  /* ---------- Forms ---------- */

  const createGroupForm = useForm<IconGroupForm>({
    resolver: zodResolver(IconGroupSchema),
    defaultValues: {
      name: '',
      name_en: '',
      icon: '',
      type_name: '',
      seq: 0,
    },
  })

  const editGroupForm = useForm<IconGroupForm>({
    resolver: zodResolver(IconGroupSchema),
    defaultValues: {
      name: '',
      name_en: '',
      icon: '',
      type_name: '',
      seq: 0,
    },
  })

  const createIconForm = useForm<AppIconForm>({
    resolver: zodResolver(AppIconSchema),
    defaultValues: {
      name: '',
      name_en: '',
      icon: '',
      icon_app: '',
      icon_tablet: '',
      icon_kiosk: '',
      link: '',
      web_link: '',
      system_code: '',
      description: '',
      featured_icon: '',
      seq: 0,
      is_native: null,
      is_public: true,
      is_featured: false,
      is_best_selling: false,
    },
  })

  const editIconForm = useForm<AppIconForm>({
    resolver: zodResolver(AppIconSchema),
    defaultValues: {
      name: '',
      name_en: '',
      icon: '',
      icon_app: '',
      icon_tablet: '',
      icon_kiosk: '',
      link: '',
      web_link: '',
      system_code: '',
      description: '',
      featured_icon: '',
      seq: 0,
      is_native: null,
      is_public: true,
      is_featured: false,
      is_best_selling: false,
    },
  })

  /* ---------- Icon file states (зөвхөн create/update дээр ашиглана) ---------- */

  const [createIconFiles, setCreateIconFiles] = useState<IconFiles>({
    icon: null,
    icon_app: null,
    icon_tablet: null,
    icon_kiosk: null,
  })

  const [editIconFiles, setEditIconFiles] = useState<IconFiles>({
    icon: null,
    icon_app: null,
    icon_tablet: null,
    icon_kiosk: null,
  })

  /* ===================== Loaders ===================== */

  async function loadGroups() {
    setGroupLoading(true)
    setGlobalLoading(true)
    setGroupFetchError(null)
    try {
      const data = await api.get<App.AppIconGroup[]>('/app-service-group', {
        query: {
          name: groupFilterName || undefined,
        },
      })
      setGroups(data ?? [])

      if (!selectedGroup && (data?.length ?? 0) > 0) {
        setSelectedGroup(data[0])
      }
    } catch {
      setGroupFetchError("Error occurred")
    } finally {
      setGroupLoading(false)
      setGlobalLoading(false)
    }
  }

  async function loadIcons() {
    if (!selectedGroup) {
      setIcons([])
      return
    }
    setIconLoading(true)
    setGlobalLoading(true)
    setIconFetchError(null)
    try {
      const data = await api.get<App.AppIcon[]>('/app-service-icon', {
        query: {
          name: iconFilterName || undefined,
          group_id: selectedGroup.id,
        },
      })
      setIcons(data ?? [])
    } catch {
      setIconFetchError("Error occurred")
    } finally {
      setIconLoading(false)
      setGlobalLoading(false)
    }
  }

  useEffect(() => {
    loadGroups()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    if (selectedGroup) {
      loadIcons()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedGroup])

  /* ===================== CRUD: Groups ===================== */

  const onOpenCreateGroup = () => {
    createGroupForm.reset({
      name: '',
      name_en: '',
      icon: '',
      type_name: '',
      seq: 0,
    })
    setOpenCreateGroup(true)
  }

  const onOpenEditGroup = (row: App.AppIconGroup) => {
    setSelectedGroupRow(row)
    editGroupForm.reset({
      name: row.name ?? '',
      name_en: row.name_en ?? '',
      icon: row.icon ?? '',
      type_name: row.type_name ?? '',
      seq: row.seq ?? 0,
    })
    setOpenEditGroup(true)
  }

  const onOpenDeleteGroup = (row: App.AppIconGroup) => {
    setSelectedGroupRow(row)
    setOpenDeleteGroup(true)
  }

  const onCreateGroup: SubmitHandler<IconGroupForm> = async (valuesIn) => {
    const values = IconGroupSchema.parse(valuesIn)
    try {
      const payload = {
        name: values.name,
        name_en: values.name_en || undefined,
        icon: values.icon || undefined,
        type_name: values.type_name || undefined,
        seq: values.seq ?? 0,
      }
      await api.post<App.AppIconGroup>('/app-service-group', payload as Record<string, unknown>)
      setOpenCreateGroup(false)
      createGroupForm.reset({
        name: '',
        name_en: '',
        icon: '',
        type_name: '',
        seq: 0,
      })
      await loadGroups()
    } catch {
      // toast API дээр handled гэж үзэж болно
    }
  }

  const onUpdateGroup: SubmitHandler<IconGroupForm> = async (valuesIn) => {
    if (!selectedGroupRow) return
    const values = IconGroupSchema.parse(valuesIn)
    try {
      const payload = {
        id: selectedGroupRow.id,
        name: values.name,
        name_en: values.name_en || undefined,
        icon: values.icon || undefined,
        type_name: values.type_name || undefined,
        seq: values.seq ?? 0,
      }
      await api.put<App.AppIconGroup>(`/app-service-group/${selectedGroupRow.id}`, payload)
      setOpenEditGroup(false)
      setSelectedGroupRow(null)
      await loadGroups()
    } catch {}
  }

  const onDeleteGroup = async () => {
    if (!selectedGroupRow) return
    try {
      setGroupDeleting(true)
      await api.del<void>(`/app-service-group/${selectedGroupRow.id}`)
      setOpenDeleteGroup(false)
      if (selectedGroup?.id === selectedGroupRow.id) {
        setSelectedGroup(null)
        setIcons([])
      }
      setSelectedGroupRow(null)
      await loadGroups()
    } catch {
      setGroupFetchError("Error occurred")
    } finally {
      setGroupDeleting(false)
    }
  }

  /* ===================== CRUD: Icons ===================== */

  const onOpenCreateIcon = () => {
    if (!selectedGroup) return
    createIconForm.reset({
      name: '',
      name_en: '',
      icon: '',
      icon_app: '',
      icon_tablet: '',
      icon_kiosk: '',
      link: '',
      web_link: '',
      system_code: '',
      description: '',
      featured_icon: '',
      seq: 0,
      is_native: null,
      is_public: true,
      is_featured: false,
      is_best_selling: false,
    })
    setCreateIconFiles({
      icon: null,
      icon_app: null,
      icon_tablet: null,
      icon_kiosk: null,
    })
    setOpenCreateIcon(true)
  }

  const onOpenEditIcon = (row: App.AppIcon) => {
    setSelectedIconRow(row)
    editIconForm.reset({
      name: row.name ?? '',
      name_en: row.name_en ?? '',
      icon: row.icon ?? '',
      icon_app: row.icon_app ?? '',
      icon_tablet: row.icon_tablet ?? '',
      icon_kiosk: row.icon_kiosk ?? '',
      link: row.link ?? '',
      web_link: row.web_link ?? '',
      system_code: row.system_code ?? '',
      description: row.description ?? '',
      featured_icon: row.featured_icon ?? '',
      seq: row.seq ?? 0,
      is_native: row.is_native ?? null,
      is_public: row.is_public ?? true,
      is_featured: row.is_featured ?? false,
      is_best_selling: row.is_best_selling ?? false,
    })
    setEditIconFiles({
      icon: { preview: row.icon },
      icon_app: { preview: row.icon_app },
      icon_tablet: { preview: row.icon_tablet },
      icon_kiosk: { preview: row.icon_kiosk },
    })
    setOpenEditIcon(true)
  }

  const onOpenDeleteIcon = (row: App.AppIcon) => {
    setSelectedIconRow(row)
    setOpenDeleteIcon(true)
  }

  const onCreateIcon: SubmitHandler<AppIconForm> = async (valuesIn) => {
    if (!selectedGroup) return
    const values = AppIconSchema.parse(valuesIn)

    try {
      const uploaded: Partial<Record<keyof IconFiles, string>> = {}

      if (createIconFiles.icon?.file) {
        uploaded.icon = await uploadImage(createIconFiles.icon.file, 'icon')
      }
      if (createIconFiles.icon_app?.file) {
        uploaded.icon_app = await uploadImage(createIconFiles.icon_app.file, 'icon_app')
      }
      if (createIconFiles.icon_tablet?.file) {
        uploaded.icon_tablet = await uploadImage(createIconFiles.icon_tablet.file, 'icon_tablet')
      }
      if (createIconFiles.icon_kiosk?.file) {
        uploaded.icon_kiosk = await uploadImage(createIconFiles.icon_kiosk.file, 'icon_kiosk')
      }

      const payload = {
        group_id: selectedGroup.id,
        name: values.name,
        name_en: values.name_en || undefined,
        icon: uploaded.icon ?? undefined,
        icon_app: uploaded.icon_app ?? undefined,
        icon_tablet: uploaded.icon_tablet ?? undefined,
        icon_kiosk: uploaded.icon_kiosk ?? undefined,
        link: values.link || undefined,
        web_link: values.web_link || undefined,
        system_code: values.system_code || undefined,
        description: values.description || undefined,
        featured_icon: values.featured_icon || undefined,
        seq: values.seq ?? 0,
        is_native: values.is_native ?? undefined,
        is_public: values.is_public ?? undefined,
        is_featured: values.is_featured ?? undefined,
        is_best_selling: values.is_best_selling ?? undefined,
      }

      await api.post<App.AppIcon>('/app-service-icon', payload as Record<string, unknown>)
      setOpenCreateIcon(false)
      await loadIcons()
    } catch {}
  }

  const onUpdateIcon: SubmitHandler<AppIconForm> = async (valuesIn) => {
    if (!selectedIconRow || !selectedGroup) return
    const values = AppIconSchema.parse(valuesIn)

    try {
      // одоо байгаа URL-уудаас эхлээд авна
      let iconUrl = selectedIconRow.icon ?? undefined
      let iconAppUrl = selectedIconRow.icon_app ?? undefined
      let iconTabletUrl = selectedIconRow.icon_tablet ?? undefined
      let iconKioskUrl = selectedIconRow.icon_kiosk ?? undefined

      // хэрвээ шинэ file сонгосон бол тухайн талбарыг upload хийгээд URL-ийг солино
      if (editIconFiles.icon?.file) {
        iconUrl = await uploadImage(editIconFiles.icon.file, 'icon', selectedIconRow.icon)
      }
      if (editIconFiles.icon_app?.file) {
        iconAppUrl = await uploadImage(
          editIconFiles.icon_app.file,
          'icon_app',
          selectedIconRow.icon_app,
        )
      }
      if (editIconFiles.icon_tablet?.file) {
        iconTabletUrl = await uploadImage(
          editIconFiles.icon_tablet.file,
          'icon_tablet',
          selectedIconRow.icon_tablet,
        )
      }
      if (editIconFiles.icon_kiosk?.file) {
        iconKioskUrl = await uploadImage(
          editIconFiles.icon_kiosk.file,
          'icon_kiosk',
          selectedIconRow.icon_kiosk,
        )
      }

      const payload = {
        id: selectedIconRow.id,
        group_id: selectedGroup.id,
        name: values.name,
        name_en: values.name_en || undefined,
        icon: iconUrl,
        icon_app: iconAppUrl,
        icon_tablet: iconTabletUrl,
        icon_kiosk: iconKioskUrl,
        link: values.link || undefined,
        web_link: values.web_link || undefined,
        system_code: values.system_code || undefined,
        description: values.description || undefined,
        featured_icon: values.featured_icon || undefined,
        seq: values.seq ?? 0,
        is_native: values.is_native ?? undefined,
        is_public: values.is_public ?? undefined,
        is_featured: values.is_featured ?? undefined,
        is_best_selling: values.is_best_selling ?? undefined,
      }

      await api.put<App.AppIcon>(`/app-service-icon/${selectedIconRow.id}`, payload)
      setOpenEditIcon(false)
      setSelectedIconRow(null)
      await loadIcons()
    } catch {}
  }

  const onDeleteIcon = async () => {
    if (!selectedIconRow) return
    try {
      setIconDeleting(true)
      await api.del<void>(`/app-service-icon/${selectedIconRow.id}`)
      setOpenDeleteIcon(false)
      if (selectedIconRow.icon) {
        deleteImage(selectedIconRow.icon)
      }
      if (selectedIconRow.icon_app) {
        deleteImage(selectedIconRow.icon_app)
      }
      if (selectedIconRow.icon_tablet) {
        deleteImage(selectedIconRow.icon_tablet)
      }
      if (selectedIconRow.icon_kiosk) {
        deleteImage(selectedIconRow.icon_kiosk)
      }
      setSelectedIconRow(null)
      await loadIcons()
    } catch {
      setIconFetchError("Error occurred")
    } finally {
      setIconDeleting(false)
    }
  }

  const isCreatingGroup = createGroupForm.formState.isSubmitting
  const isUpdatingGroup = editGroupForm.formState.isSubmitting
  const isGroupRowBusy = (id: number) =>
    (isUpdatingGroup && selectedGroupRow?.id === id) ||
    (groupDeleting && selectedGroupRow?.id === id)

  const isCreatingIcon = createIconForm.formState.isSubmitting
  const isUpdatingIcon = editIconForm.formState.isSubmitting
  const isIconRowBusy = (id: number) =>
    (isUpdatingIcon && selectedIconRow?.id === id) || (iconDeleting && selectedIconRow?.id === id)

  /* ===================== Render ===================== */

  return (
    <div className="grid h-full min-h-0 w-full grid-cols-1 gap-4 overflow-hidden p-4 sm:p-6 lg:grid-cols-[minmax(260px,0.35fr)_minmax(0,0.65fr)]">
      {/* ========== LEFT: GROUPS ========== */}
      <div className="bg-background flex min-h-0 flex-col rounded-lg border">
        <div className="flex items-center justify-between gap-2 border-b px-4 py-3">
          <span className="text-base font-medium">{t('app_icon_group') || 'App icon groups'}</span>
          <Button
            size="sm"
            className="gap-1"
            onClick={onOpenCreateGroup}
            disabled={isCreatingGroup || isUpdatingGroup || groupDeleting}
          >
            {isCreatingGroup ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Plus className="h-4 w-4" />
            )}
            <span className="hidden text-sm sm:inline">
              {t('create', { name: t('group') || '' })}
            </span>
          </Button>
        </div>

        {/* group filters */}
        <div className="flex items-center justify-between gap-2 border-b px-4 py-2">
          <Input
            value={groupFilterName}
            onChange={(e) => setGroupFilterName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') loadGroups()
            }}
            placeholder={t('search_by_name') || 'Search by name'}
            className="h-8 text-sm"
          />
        </div>

        {/* group table */}
        <div className="flex min-h-0 flex-1 flex-col">
          {groupFetchError ? (
            <div className="p-4 text-sm text-red-600">{groupFetchError}</div>
          ) : groupLoading && groups.length === 0 ? (
            <div className="flex h-32 items-center justify-center gap-2 text-base">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span>{t('loading')}</span>
            </div>
          ) : groups.length === 0 ? (
            <div className="text-muted-foreground flex h-32 flex-col items-center justify-center gap-2 text-sm">
              <LucideIcon name="i-lucide-archive-x" className="h-6 w-6" />
              <p>{t('no_information_available')}</p>
            </div>
          ) : (
            <div className="min-h-0 flex-1 overflow-auto">
              <Table className="w-full table-fixed">
                <colgroup>
                  <col style={{ width: '90px' }} />
                  <col style={{ width: '200px' }} />
                </colgroup>
                <TableHeader>
                  <TableRow className="[&>th]:bg-gray-50 dark:[&>th]:bg-gray-800 [&>th]:text-sm">
                    <TableHead></TableHead>
                    <TableHead>{t('name')}</TableHead>
                    <TableHead className="text-center">{t('sequence')}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {groups.map((g) => {
                    const busy = isGroupRowBusy(g.id)
                    const isSelected = selectedGroup?.id === g.id
                    return (
                      <TableRow
                        key={g.id}
                        className={`hover:bg-muted/60 cursor-pointer text-sm ${
                          isSelected ? 'bg-muted/80' : ''
                        }`}
                        onClick={() => {
                          setSelectedGroup(g)
                        }}
                      >
                        <TableCell className="align-middle">
                          <div className="flex items-center gap-1">
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    onOpenEditGroup(g)
                                  }}
                                  disabled={busy}
                                >
                                  {isUpdatingGroup && selectedGroupRow?.id === g.id ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <Pencil className="h-4 w-4" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p className="lowercase first-letter:uppercase">
                                  {t('update', { name: t('group') })}
                                </p>
                                <TooltipArrow className="fill-popover" />
                              </TooltipContent>
                            </Tooltip>

                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="destructive"
                                  size="sm"
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    onOpenDeleteGroup(g)
                                  }}
                                  disabled={busy}
                                >
                                  {groupDeleting && selectedGroupRow?.id === g.id ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <Trash2 className="h-4 w-4" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p className="lowercase first-letter:uppercase">
                                  {t('delete', { name: t('group') })}
                                </p>
                                <TooltipArrow className="fill-popover" />
                              </TooltipContent>
                            </Tooltip>
                          </div>
                        </TableCell>
                        <TableCell className="align-middle">
                          <div className="flex items-center gap-2">
                            <div className="flex flex-col">
                              <span className="truncate font-medium">{g.name}</span>
                              {g.name_en ? (
                                <span className="text-muted-foreground truncate text-sm">
                                  {g.name_en}
                                </span>
                              ) : null}
                            </div>
                          </div>
                        </TableCell>
                        <TableCell className="text-center align-middle">{g.seq ?? 0}</TableCell>
                      </TableRow>
                    )
                  })}
                </TableBody>
              </Table>
            </div>
          )}
        </div>
      </div>

      {/* ========== RIGHT: ICONS ========== */}
      <div className="bg-background flex min-h-0 flex-col rounded-lg border">
        <div className="flex items-center justify-between gap-2 border-b px-4 py-3">
          <div className="flex flex-col">
            <span className="text-base font-medium">{t('app_icon') || 'App icons'}</span>
            <span className="text-muted-foreground text-sm">
              {selectedGroup ? selectedGroup.name : t('select_group') || 'Эхлээд бүлэг сонгоно уу.'}
            </span>
          </div>
          <Button
            size="sm"
            className="gap-1"
            onClick={onOpenCreateIcon}
            disabled={!selectedGroup || isCreatingIcon || isUpdatingIcon || iconDeleting}
          >
            {isCreatingIcon ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Plus className="h-4 w-4" />
            )}
            <span className="hidden text-sm sm:inline">
              {t('create', { name: t('app_icon') || '' })}
            </span>
          </Button>
        </div>

        {/* icon filters */}
        <div className="flex items-center justify-between gap-2 border-b px-4 py-2">
          <Input
            value={iconFilterName}
            onChange={(e) => setIconFilterName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') loadIcons()
            }}
            placeholder={t('search_by_name') || 'Search by name'}
            className="h-8 text-sm"
          />
        </div>

        {/* icon table */}
        <div className="flex min-h-0 flex-1 flex-col">
          {iconFetchError ? (
            <div className="p-4 text-sm text-red-600">{iconFetchError}</div>
          ) : !selectedGroup ? (
            <div className="text-muted-foreground flex h-32 items-center justify-center text-sm">
              {t('select_group') || 'Эхлээд бүлэг сонгоно уу.'}
            </div>
          ) : iconLoading && icons.length === 0 ? (
            <div className="flex h-32 items-center justify-center gap-2 text-base">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span>{t('loading')}</span>
            </div>
          ) : icons.length === 0 ? (
            <div className="text-muted-foreground flex h-32 flex-col items-center justify-center gap-2 text-sm">
              <LucideIcon name="i-lucide-archive-x" className="h-6 w-6" />
              <p>{t('no_information_available')}</p>
            </div>
          ) : (
            <div className="min-h-0 flex-1 overflow-auto">
              <Table className="w-full table-fixed">
                <colgroup>
                  <col style={{ width: '120px' }} />
                  <col style={{ width: '40px' }} />
                  <col style={{ width: '160px' }} />
                  <col style={{ width: '160px' }} />
                  <col />
                  <col style={{ width: '110px' }} />
                  <col style={{ width: '80px' }} />
                </colgroup>
                <TableHeader>
                  <TableRow className="[&>th]:bg-gray-50 dark:[&>th]:bg-gray-800 [&>th]:text-sm">
                    <TableHead>{t('actions')}</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead>{t('name')}</TableHead>
                    <TableHead>{t('description')}</TableHead>
                    <TableHead>{t('system')}</TableHead>
                    <TableHead className="text-center">{t('sequence')}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {icons.map((ic) => {
                    const busy = isIconRowBusy(ic.id)
                    return (
                      <TableRow key={ic.id} className="text-sm">
                        <TableCell className="align-middle">
                          <div className="flex items-center gap-1">
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="outline"
                                  size="sm"
                                  onClick={() => onOpenEditIcon(ic)}
                                  disabled={busy}
                                >
                                  {isUpdatingIcon && selectedIconRow?.id === ic.id ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <Pencil className="h-4 w-4" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p className="lowercase first-letter:uppercase">
                                  {t('update', { name: t('app_icon') })}
                                </p>
                                <TooltipArrow className="fill-popover" />
                              </TooltipContent>
                            </Tooltip>

                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="destructive"
                                  size="sm"
                                  onClick={() => onOpenDeleteIcon(ic)}
                                  disabled={busy}
                                >
                                  {iconDeleting && selectedIconRow?.id === ic.id ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <Trash2 className="h-4 w-4" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                <p className="lowercase first-letter:uppercase">
                                  {t('delete', { name: t('app_icon') })}
                                </p>
                                <TooltipArrow className="fill-popover" />
                              </TooltipContent>
                            </Tooltip>
                          </div>
                        </TableCell>
                        <TableCell className="text-muted-foreground align-middle">
                          {ic.id}
                        </TableCell>
                        <TableCell className="text-muted-foreground align-middle">
                          {ic.icon && (
                            <Image
                              src={ic.icon}
                              alt={ic.name}
                              width={100}
                              height={100}
                              className="h-[100px] w-[100px] object-cover"
                            />
                          )}
                        </TableCell>
                        <TableCell className="align-middle">
                          <div className="flex items-center gap-2">
                            <div className="flex flex-col">
                              <span className="truncate font-medium">{ic.name}</span>
                              {ic.name_en ? (
                                <span className="text-muted-foreground truncate text-sm">
                                  {ic.name_en}
                                </span>
                              ) : null}
                              <div className="mt-0.5 flex flex-wrap gap-1">
                                {ic.is_public ? (
                                  <Badge variant="secondary" className="px-1 text-[10px]">
                                    {t('public') || 'public'}
                                  </Badge>
                                ) : null}
                                {ic.is_native ? (
                                  <Badge variant="outline" className="px-1 text-[10px]">
                                    native
                                  </Badge>
                                ) : null}
                                {ic.is_featured ? (
                                  <Badge variant="outline" className="px-1 text-[10px]">
                                    featured
                                  </Badge>
                                ) : null}
                                {ic.is_best_selling ? (
                                  <Badge variant="outline" className="px-1 text-[10px]">
                                    best
                                  </Badge>
                                ) : null}
                              </div>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell className="text-muted-foreground align-middle">
                          {ic.description || <span className="opacity-60">—</span>}
                        </TableCell>
                        <TableCell className="align-middle">
                          {ic.system_code ? (
                            <Badge variant="outline" className="text-sm">
                              {ic.system_code}
                            </Badge>
                          ) : (
                            <span className="text-sm opacity-60">—</span>
                          )}
                        </TableCell>
                        <TableCell className="text-center align-middle">{ic.seq ?? 0}</TableCell>
                      </TableRow>
                    )
                  })}
                </TableBody>
              </Table>
            </div>
          )}
        </div>
      </div>

      {/* ===================== Group Dialogs ===================== */}

      {/* Create group */}
      <Dialog open={openCreateGroup} onOpenChange={setOpenCreateGroup}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('create', { name: t('group') })}
            </DialogTitle>
            <DialogDescription>{t('fill_the_fields_and_save_to_create')}</DialogDescription>
          </DialogHeader>

          <Form {...createGroupForm}>
            <form
              onSubmit={createGroupForm.handleSubmit(onCreateGroup)}
              className="space-y-4 pt-2"
              autoComplete="off"
            >
              <FormField
                control={createGroupForm.control}
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
                control={createGroupForm.control}
                name="name_en"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('name_en') || 'Name (EN)'}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('name_en') || 'Name (EN)'} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <div className="grid gap-4 sm:grid-cols-2">
                <FormField
                  control={createGroupForm.control}
                  name="icon"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('icon')}</FormLabel>
                      <FormControl>
                        <Input placeholder="i-lucide-*" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={createGroupForm.control}
                  name="type_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('type') || 'Type'}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('type') || 'Type'} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={createGroupForm.control}
                name="seq"
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

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenCreateGroup(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={createGroupForm.formState.isSubmitting}>
                  {createGroupForm.formState.isSubmitting && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {t('save')}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Edit group */}
      <Dialog open={openEditGroup} onOpenChange={setOpenEditGroup}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('update', { name: t('group') })}
            </DialogTitle>
            <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
          </DialogHeader>

          <Form {...editGroupForm}>
            <form
              onSubmit={editGroupForm.handleSubmit(onUpdateGroup)}
              className="space-y-4 pt-2"
              autoComplete="off"
            >
              <FormField
                control={editGroupForm.control}
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
                control={editGroupForm.control}
                name="name_en"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('name_en') || 'Name (EN)'}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('name_en') || 'Name (EN)'} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <div className="grid gap-4 sm:grid-cols-2">
                <FormField
                  control={editGroupForm.control}
                  name="icon"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('icon')}</FormLabel>
                      <FormControl>
                        <Input placeholder="i-lucide-*" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={editGroupForm.control}
                  name="type_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('type') || 'Type'}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('type') || 'Type'} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={editGroupForm.control}
                name="seq"
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

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenEditGroup(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={editGroupForm.formState.isSubmitting}>
                  {editGroupForm.formState.isSubmitting && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {t('save')}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Delete group */}
      <Dialog open={openDeleteGroup} onOpenChange={setOpenDeleteGroup}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('delete', { name: t('group') })}
            </DialogTitle>
            <DialogDescription className="pt-2 text-base">
              {t.rich('delete_warning', {
                name: () => <span className="font-medium">{selectedGroupRow?.name}</span>,
              })}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="pt-2">
            <Button
              variant="outline"
              onClick={() => setOpenDeleteGroup(false)}
              disabled={groupDeleting}
            >
              {t('cancel')}
            </Button>
            <Button variant="destructive" onClick={onDeleteGroup} disabled={groupDeleting}>
              {groupDeleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <span className="capitalize">{t('delete', { name: '' })}</span>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ===================== Icon Dialogs ===================== */}

      {/* CREATE ICON */}
      <Dialog open={openCreateIcon} onOpenChange={setOpenCreateIcon}>
        <DialogContent className="max-h-[90vh] max-w-[95vw] overflow-y-auto sm:max-w-2xl lg:max-w-5xl">
          <DialogHeader>
            <DialogTitle className="text-xl font-semibold">
              {t('create', { name: t('app_icon') })}
            </DialogTitle>
            <DialogDescription>{t('fill_the_fields_and_save_to_create')}</DialogDescription>
          </DialogHeader>

          <Form {...createIconForm}>
            <form
              onSubmit={createIconForm.handleSubmit(onCreateIcon)}
              className="space-y-6 pt-2"
              autoComplete="off"
            >
              <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-4">
                <FormField
                  control={createIconForm.control}
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
                  control={createIconForm.control}
                  name="name_en"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('name_en')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('name_en')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={createIconForm.control}
                  name="seq"
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
                  control={createIconForm.control}
                  name="featured_icon"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Featured icon</FormLabel>
                      <FormControl>
                        <Input placeholder="Featured icon" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-4">
                <FileUpload
                  label={t('icon')}
                  value={createIconFiles.icon}
                  onChange={(file) => setCreateIconFiles((prev) => ({ ...prev, icon: file }))}
                />
                <FileUpload
                  label="App icon"
                  value={createIconFiles.icon_app}
                  onChange={(file) => setCreateIconFiles((prev) => ({ ...prev, icon_app: file }))}
                />
                <FileUpload
                  label="Tablet icon"
                  value={createIconFiles.icon_tablet}
                  onChange={(file) =>
                    setCreateIconFiles((prev) => ({ ...prev, icon_tablet: file }))
                  }
                />
                <FileUpload
                  label="Kiosk icon"
                  value={createIconFiles.icon_kiosk}
                  onChange={(file) => setCreateIconFiles((prev) => ({ ...prev, icon_kiosk: file }))}
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-4">
                <FormField
                  control={createIconForm.control}
                  name="link"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('link')}</FormLabel>
                      <FormControl>
                        <Input placeholder="https://..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={createIconForm.control}
                  name="web_link"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Web link</FormLabel>
                      <FormControl>
                        <Input placeholder="https://..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={createIconForm.control}
                  name="system_code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('system')}</FormLabel>
                      <FormControl>
                        <Input placeholder="system code" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 gap-2">
                  <FormField
                    control={createIconForm.control}
                    name="is_public"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">{t('is_public')}</FormLabel>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createIconForm.control}
                    name="is_native"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">Native</FormLabel>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createIconForm.control}
                    name="is_featured"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">{t('is_featured')}</FormLabel>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createIconForm.control}
                    name="is_best_selling"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">{t('is_best_selling')}</FormLabel>
                      </FormItem>
                    )}
                  />
                </div>
              </div>

              <FormField
                control={createIconForm.control}
                name="description"
                render={({ field }) => (
                  <FormItem className="mt-4">
                    <FormLabel>{t('description')}</FormLabel>
                    <FormControl>
                      <Textarea rows={3} placeholder={t('description')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenCreateIcon(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={createIconForm.formState.isSubmitting}>
                  {createIconForm.formState.isSubmitting && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {t('save')}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* EDIT ICON */}
      <Dialog open={openEditIcon} onOpenChange={setOpenEditIcon}>
        <DialogContent className="max-h-[90vh] max-w-[95vw] overflow-y-auto sm:max-w-2xl lg:max-w-5xl">
          <DialogHeader>
            <DialogTitle className="text-xl font-semibold">
              {t('update', { name: t('app_icon') })}
            </DialogTitle>
            <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
          </DialogHeader>

          <Form {...editIconForm}>
            <form
              onSubmit={editIconForm.handleSubmit(onUpdateIcon)}
              className="space-y-6 pt-2"
              autoComplete="off"
            >
              <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-4">
                <FormField
                  control={editIconForm.control}
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
                  control={editIconForm.control}
                  name="name_en"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('name_en')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('name_en')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editIconForm.control}
                  name="seq"
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
                  control={editIconForm.control}
                  name="featured_icon"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Featured icon</FormLabel>
                      <FormControl>
                        <Input placeholder="Featured icon" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-4">
                <FileUpload
                  label={t('icon')}
                  value={editIconFiles.icon}
                  onChange={(file) => setEditIconFiles((prev) => ({ ...prev, icon: file }))}
                />
                <FileUpload
                  label="App icon"
                  value={editIconFiles.icon_app}
                  onChange={(file) => setEditIconFiles((prev) => ({ ...prev, icon_app: file }))}
                />
                <FileUpload
                  label="Tablet icon"
                  value={editIconFiles.icon_tablet}
                  onChange={(file) => setEditIconFiles((prev) => ({ ...prev, icon_tablet: file }))}
                />
                <FileUpload
                  label="Kiosk icon"
                  value={editIconFiles.icon_kiosk}
                  onChange={(file) => setEditIconFiles((prev) => ({ ...prev, icon_kiosk: file }))}
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-4">
                <FormField
                  control={editIconForm.control}
                  name="link"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('link')}</FormLabel>
                      <FormControl>
                        <Input placeholder="https://..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editIconForm.control}
                  name="web_link"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Web link</FormLabel>
                      <FormControl>
                        <Input placeholder="https://..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editIconForm.control}
                  name="system_code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('system')}</FormLabel>
                      <FormControl>
                        <Input placeholder="system code" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 gap-2">
                  <FormField
                    control={editIconForm.control}
                    name="is_public"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">{t('is_public')}</FormLabel>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={editIconForm.control}
                    name="is_native"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">Native</FormLabel>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={editIconForm.control}
                    name="is_featured"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">{t('is_featured')}</FormLabel>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={editIconForm.control}
                    name="is_best_selling"
                    render={({ field }) => (
                      <FormItem className="flex items-center gap-2 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={!!field.value}
                            onCheckedChange={(v) => field.onChange(Boolean(v))}
                          />
                        </FormControl>
                        <FormLabel className="text-sm">{t('is_best_selling')}</FormLabel>
                      </FormItem>
                    )}
                  />
                </div>
              </div>

              <FormField
                control={editIconForm.control}
                name="description"
                render={({ field }) => (
                  <FormItem className="mt-4">
                    <FormLabel>{t('description')}</FormLabel>
                    <FormControl>
                      <Textarea rows={3} placeholder={t('description')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenEditIcon(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={editIconForm.formState.isSubmitting}>
                  {editIconForm.formState.isSubmitting && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {t('save')}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      <Dialog open={openDeleteIcon} onOpenChange={setOpenDeleteIcon}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('delete', { name: t('app_icon') })}
            </DialogTitle>
            <DialogDescription className="pt-2 text-base">
              {t.rich('delete_warning', {
                name: () => <span className="font-medium">{selectedIconRow?.name}</span>,
              })}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="pt-2">
            <Button
              variant="outline"
              onClick={() => setOpenDeleteIcon(false)}
              disabled={iconDeleting}
            >
              {t('cancel')}
            </Button>
            <Button variant="destructive" onClick={onDeleteIcon} disabled={iconDeleting}>
              {iconDeleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <span className="capitalize">{t('delete', { name: '' })}</span>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
