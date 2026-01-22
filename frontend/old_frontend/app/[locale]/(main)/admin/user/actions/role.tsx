/**
 * 🛡️👤 User Roles Dialog Component
 * (/[locale]/(main)/admin/user/actions/role.tsx)
 * 
 * Энэ нь хэрэглэгчийн дүр удирдах dialog component юм.
 * Зорилго: Тухайн хэрэглэгчид дүр (role) оноох, хасах
 * 
 * Features:
 * - ✅ Two-level dialog system:
 *   1. Main dialog: View assigned roles (no pagination)
 *   2. Add dialog: Select roles from catalog (card grid)
 * - ✅ Remove assigned roles
 * - ✅ Multi-select with checkboxes
 * - ✅ Card-based selection UI
 * - ✅ Bulk select/deselect
 * - ✅ Already assigned detection
 * - ✅ Loading states
 * - ✅ Error handling
 * 
 * Props:
 * @param open - Dialog visibility
 * @param onOpenChange - Toggle dialog
 * @param user - User object (UserLite)
 * @param onChanged - Callback after changes
 * 
 * User Workflow:
 * 1. Open dialog → See assigned roles
 * 2. Click "Add" → Open catalog dialog
 * 3. Select roles (card grid with checkboxes)
 * 4. Save → Assign roles to user
 * 5. Remove role → Confirm and delete
 * 
 * API Endpoints:
 * - GET /role-matrix/roles?user_id=... - Get assigned roles
 * - GET /role?page=1&size=500 - Get all roles
 * - POST /role-matrix { user_id, role_ids: [] } - Assign roles
 * - DELETE /role-matrix { user_id, role_id } - Remove role
 * 
 * UI Pattern:
 * - Main: Table view of assigned roles
 * - Add: Card grid with multi-select
 * - No pagination in main view (all roles shown)
 * - Large catalog load (500 items)
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
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Checkbox } from '@/components/ui/checkbox'
import { Card as UiCard, CardContent } from '@/components/ui/card'
import { Loader2, Plus, Trash2 } from 'lucide-react'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'

/**
 * 🧩 UserLite type - Энгийн хэрэглэгчийн мэдээлэл
 */
export type UserLite = {
  id: number
  first_name?: string | null
  last_name?: string | null
  reg_no?: string | null
}
type RoleRow = App.UserRole

/**
 * 🧩 Component Props
 */
type Props = {
  open: boolean
  onOpenChange: (v: boolean) => void
  user: UserLite | null
  onChanged?: () => void
}

export default function UserRolesDialog({ open, onOpenChange, user, onChanged }: Props) {
  const t = useTranslations()

  // ========================================
  // 📋 Main Dialog: Assigned Roles
  // ========================================
  
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [rows, setRows] = React.useState<RoleRow[]>([])
  const [removingId, setRemovingId] = React.useState<number | null>(null)

  /**
   * 📥 Хэрэглэгчид оноогдсон дүрүүдийг ачаалах
   */
  const loadAssigned = React.useCallback(async () => {
    if (!user) return
    setLoading(true)
    setError(null)
    try {
      // pagination-гүй, бүх roles буцаана
      const res = await api.get<RoleRow[] | App.ListData<RoleRow>>('/role-matrix/roles', {
        query: { user_id: user.id },
        cache: 'no-store',
      })
      const items = Array.isArray(res) ? res : (res.items ?? [])
      setRows(items)
    } catch {
      setError("Error occurred")
    } finally {
      setLoading(false)
    }
  }, [user])

  React.useEffect(() => {
    if (open) loadAssigned()
    else {
      setRows([])
      setError(null)
    }
  }, [open, loadAssigned])

  /**
   * 🗑️ Дүр хасах handler
   * @param roleId - Хасах дүрийн ID
   */
  const removeOne = async (roleId: number) => {
    if (!user) return
    try {
      setRemovingId(roleId)
      await api.del('/role-matrix', { user_id: user.id, role_id: roleId } as Record<string, unknown>)
      setRows((prev) => prev.filter((r) => r.role_id !== roleId))
      onChanged?.()
    } catch {
      setError("Error occurred")
    } finally {
      setRemovingId(null)
    }
  }

  // ========================================
  // ➕ Add Roles Dialog
  // ========================================
  
  const [openAdd, setOpenAdd] = React.useState(false)
  const [catalogLoading, setCatalogLoading] = React.useState(false)
  const [catalogErr, setCatalogErr] = React.useState<string | null>(null)
  const [catalog, setCatalog] = React.useState<App.Role[]>([])
  const [picked, setPicked] = React.useState<Set<number>>(new Set())
  const [saving, setSaving] = React.useState(false)

  const assignedIdSet = React.useMemo(() => new Set(rows.map((r) => r.role_id)), [rows])

  /**
   * 📥 Бүх дүрүүдийн catalog ачаалах
   */
  const loadCatalog = React.useCallback(async () => {
    setCatalogErr(null)
    setCatalog([])
    setPicked(new Set())
    setCatalogLoading(true)
    try {
      // бүх roles-ыг том size-р аваад дуусгая
      const res = await api.get<App.ListData<App.Role>>('/role', {
        query: { page: 1, size: 500 },
        cache: 'no-store',
      })
      const items = Array.isArray(res) ? res : (res.items ?? [])
      setCatalog(items)
    } catch {
      setCatalogErr("Error occurred")
    } finally {
      setCatalogLoading(false)
    }
  }, [])

  /**
   * ☑️ Checkbox toggle handler
   */
  const togglePick = (id: number, on?: boolean) => {
    setPicked((prev) => {
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

  /**
   * ☑️ Бүх харагдаж буй дүрүүдийг сонгох
   */
  const selectAllVisible = () => {
    const ids = catalog.filter((c) => !assignedIdSet.has(c.id)).map((c) => c.id)
    setPicked(new Set(ids))
  }
  
  /**
   * 🧹 Бүх сонголтыг цэвэрлэх
   */
  const clearAll = () => setPicked(new Set())

  /**
   * 💾 Сонгогдсон дүрүүдийг хадгалах
   */
  const savePicked = async () => {
    if (!user) return
    const role_ids = Array.from(picked).filter((id) => !assignedIdSet.has(id))
    if (role_ids.length === 0) {
      setOpenAdd(false)
      return
    }
    try {
      setSaving(true)
      await api.post('/role-matrix', { user_id: user.id, role_ids } as Record<string, unknown>)
      await loadAssigned()
      onChanged?.()
      setOpenAdd(false)
    } catch {
      setCatalogErr("Error occurred")
    } finally {
      setSaving(false)
    }
  }

  const fullUserName =
    [user?.last_name, user?.first_name].filter(Boolean).join(' ') ||
    user?.reg_no ||
    `User #${user?.id ?? ''}`

  return (
    <>
      {/* ===== Main dialog: Assigned roles (no pagination) ===== */}
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-3xl">
          <DialogHeader>
            <DialogTitle>
              {t('role')} — {fullUserName}
            </DialogTitle>
          </DialogHeader>

          <div className="mb-2 flex items-center justify-between gap-3">
            <Button
              size="sm"
              onClick={() => {
                setOpenAdd(true)
                loadCatalog()
              }}
            >
              <Plus className="mr-1 h-4 w-4" /> {t('add')}
            </Button>
          </div>

          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-16">{t('actions')}</TableHead>
                  <TableHead className="w-20">ID</TableHead>
                  <TableHead>{t('name')}</TableHead>
                  <TableHead>{t('description')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {error ? (
                  <TableRow>
                    <TableCell colSpan={4} className="py-8 text-center text-sm text-red-600">
                      {error}
                    </TableCell>
                  </TableRow>
                ) : loading ? (
                  <TableRow>
                    <TableCell colSpan={4} className="py-8 text-center text-sm opacity-70">
                      <Loader2 className="mr-2 inline h-4 w-4 animate-spin" /> {t('loading')}
                    </TableCell>
                  </TableRow>
                ) : rows.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} className="py-8 text-center text-sm opacity-70">
                      {t('no_information_available')}
                    </TableCell>
                  </TableRow>
                ) : (
                  rows.map((r) => (
                    <TableRow key={r.role_id}>
                      <TableCell className="text-right">
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={() => removeOne(r.role_id)}
                          disabled={removingId === r.role_id}
                          aria-label={t('delete', { name: t('role') })}
                        >
                          {removingId === r.role_id ? (
                            <Loader2 className="h-4 w-4 animate-spin" />
                          ) : (
                            <Trash2 className="h-4 w-4" />
                          )}
                        </Button>
                      </TableCell>
                      <TableCell className="text-muted-foreground text-xs">
                        {String(r.role_id)}
                      </TableCell>
                      <TableCell className="font-medium">{r.role.name}</TableCell>
                      <TableCell className="text-muted-foreground">
                        {r.role.description || '—'}
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              {t('close')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ===== Add roles dialog: Card grid, multi-select, no search ===== */}
      <Dialog
        open={openAdd}
        onOpenChange={(v) => {
          setOpenAdd(v)
          if (!v) {
            setPicked(new Set())
            setCatalogErr(null)
          }
        }}
      >
        <DialogContent className="sm:max-w-5xl">
          <DialogHeader>
            <DialogTitle>
              {t('add')} — {t('role')}
            </DialogTitle>
          </DialogHeader>

          {/* Bulk controls */}
          <div className="mb-3 flex items-center justify-between gap-2">
            <div className="text-muted-foreground text-sm">
              {catalogLoading ? t('loading') : t('total') + ': ' + (catalog?.length ?? 0)}
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={selectAllVisible}
                disabled={catalogLoading}
              >
                {t('select_all')}
              </Button>
              <Button variant="outline" size="sm" onClick={clearAll} disabled={!picked.size}>
                {t('clear')}
              </Button>
              {picked.size > 0 && (
                <span className="rounded bg-emerald-500/10 px-2 py-1 text-sm text-emerald-700">
                  {t('selected')}: <b>{picked.size}</b>
                </span>
              )}
            </div>
          </div>

          {/* Card grid */}
          {catalogErr && <div className="text-sm text-red-600">{catalogErr}</div>}

          {catalogLoading ? (
            <div className="py-10 text-center text-sm opacity-70">
              <Loader2 className="mx-auto mb-2 h-4 w-4 animate-spin" />
              {t('loading')}
            </div>
          ) : catalog.length === 0 ? (
            <div className="py-10 text-center text-sm opacity-70">
              {t('no_information_available')}
            </div>
          ) : (
            <div className="max-h-[60vh] overflow-y-auto rounded-md border p-3">
              <div className="grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
                {catalog.map((r) => {
                  const already = assignedIdSet.has(r.id)
                  const checked = already ? true : picked.has(r.id)
                  return (
                    <UiCard
                      key={r.id}
                      className={`cursor-pointer transition-all ${
                        already
                          ? 'opacity-50 cursor-not-allowed'
                          : checked
                            ? 'border-primary bg-primary/5 shadow-sm'
                            : 'hover:border-primary/50 hover:shadow-sm'
                      }`}
                      onClick={() => !already && togglePick(r.id)}
                      role="button"
                      tabIndex={already ? -1 : 0}
                    >
                      <CardContent className="flex items-start gap-2 p-2.5">
                        <Checkbox
                          checked={checked}
                          disabled={already}
                          onCheckedChange={(v) => !already && togglePick(r.id, Boolean(v))}
                          aria-label={`pick-${r.id}`}
                          className="mt-0.5 shrink-0"
                          onClick={(e) => e.stopPropagation()}
                        />
                        <div className="min-w-0 flex-1">
                          <div className="truncate text-sm font-medium leading-tight">{r.name}</div>
                          {r.description && (
                            <div className="text-muted-foreground mt-0.5 line-clamp-1 text-xs">
                              {r.description}
                            </div>
                          )}
                          <div className="text-muted-foreground mt-1 flex items-center gap-2 text-[10px]">
                            <span>ID: {r.id}</span>
                            {already && (
                              <span className="rounded bg-muted px-1 py-0.5 text-[10px]">
                                {t('assigned')}
                              </span>
                            )}
                          </div>
                        </div>
                      </CardContent>
                    </UiCard>
                  )
                })}
              </div>
            </div>
          )}

          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setOpenAdd(false)} disabled={saving}>
              {t('cancel')}
            </Button>
            <Button onClick={savePicked} disabled={saving || picked.size === 0}>
              {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {t('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
