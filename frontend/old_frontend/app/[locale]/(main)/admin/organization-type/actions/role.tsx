/**
 * 🏛️👤 Organization Type Roles Dialog Component
 * (/[locale]/(main)/admin/organization-type/actions/role.tsx)
 * 
 * Энэ нь байгууллагын төрөлд дүр оноох dialog юм.
 * Зорилго: Тухайн төрлийн байгууллагууд ямар дүрүүдтэй болохыг тохируулах
 * 
 * Features:
 * - ✅ Multi-select card grid interface
 * - ✅ Search/filter roles
 * - ✅ Select all / deselect all
 * - ✅ Visual selection indicators
 * - ✅ Already assigned detection
 * - ✅ Save selected roles
 * - ✅ Beautiful card-based UI
 * - ✅ Loading states
 * 
 * Props:
 * @param open - Dialog visibility
 * @param onOpenChange - Toggle dialog
 * @param orgType - Organization type object
 * 
 * Workflow:
 * 1. Open dialog → Load all roles + assigned roles
 * 2. Display card grid with checkboxes
 * 3. Search/filter available
 * 4. Multi-select roles
 * 5. Save → Update assignments
 * 
 * UI Pattern:
 * - Card grid layout (2 columns on desktop)
 * - Checkbox + CheckCircle2 for selection
 * - Ring border on selected cards
 * - Search input for filtering
 * - "Select All" checkbox
 * 
 * API Endpoints:
 * - GET /role?page=1&size=500 - Get all roles
 * - GET /orgtype/role?type_id=... - Get assigned roles
 * - POST /orgtype/role { type_id, role_ids: [] } - Save assignments
 * 
 * Data Flow:
 * - Load all roles (500 max)
 * - Load assigned role IDs
 * - User selects/deselects
 * - Save sends full ID array
 * 
 * @author Gerege Core Team
 */

'use client'

import * as React from 'react'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Loader2, UserCog, CheckCircle2 } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

/**
 * 🧩 Component Props
 */
type Props = {
  open: boolean
  onOpenChange: (v: boolean) => void
  orgType: App.OrganizationType | null
}

type RoleRow = App.Role

export default function OrgTypeRolesDialog({ open, onOpenChange, orgType }: Props) {
  const t = useTranslations()

  // 📊 State management
  const [loading, setLoading] = React.useState(false)
  const [saving, setSaving] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  const [roles, setRoles] = React.useState<RoleRow[]>([])
  const [selectedIds, setSelectedIds] = React.useState<Set<number>>(new Set())

  const [q, setQ] = React.useState('')

  /**
   * 📥 Бүх дүрүүдийг ачаалах
   */
  const fetchAllRoles = React.useCallback(async () => {
    const data = await api.get<App.ListData<RoleRow> | RoleRow[]>('/role', {
      query: { page: 1, size: 500 },
    })
    return Array.isArray(data) ? data : (data.items ?? [])
  }, [])

  /**
   * 📥 Байгууллагын төрөлд оноогдсон дүрүүдийг ачаалах
   */
  const fetchSelectedRoles = React.useCallback(async (orgTypeId: number) => {
    const data = await api.get<App.Role[]>('/orgtype/role', {
      query: { type_id: orgTypeId, page_size: 500 },
    })
    const ids = (data as App.Role[]).map((x) => x.id)

    return new Set<number>(ids)
  }, [])

  /**
   * 🔄 Dialog нээгдэх үед өгөгдөл ачаална
   */
  React.useEffect(() => {
    if (!open || !orgType?.id) return
    setError(null)
    setSaving(false)
    setSelectedIds(new Set())
    ;(async () => {
      setLoading(true)
      try {
        const [all, selected] = await Promise.all([
          fetchAllRoles(),
          fetchSelectedRoles(orgType.id),
        ])
        setRoles(all)
        setSelectedIds(selected)
      } catch {
        setError("Error occurred")
      } finally {
        setLoading(false)
      }
    })()
  }, [open, orgType?.id, fetchAllRoles, fetchSelectedRoles])

  /**
   * 🔍 Хайлтын filter
   */
  const filtered = React.useMemo(() => {
    const s = q.trim().toLowerCase()
    if (!s) return roles
    return roles.filter(
      (x) =>
        (x.name ?? '').toLowerCase().includes(s) ||
        (x.code ?? '').toLowerCase().includes(s) ||
        String(x.id).includes(s),
    )
  }, [q, roles])

  /**
   * ☑️ Дүрийг сонгох/цуцлах toggle
   */
  const toggle = (id: number) =>
    setSelectedIds((prev) => {
      const n = new Set(prev)
      if (n.has(id)) {
        n.delete(id)
      } else {
        n.add(id)
      }
      return n
    })

  /**
   * ☑️ Бүх харагдаж буй дүрүүдийг сонгох/цуцлах
   */
  const toggleAllFiltered = () =>
    setSelectedIds((prev) => {
      const n = new Set(prev)
      const ids = filtered.map((s) => s.id)
      const every = ids.every((id) => n.has(id))
      if (every) ids.forEach((id) => n.delete(id))
      else ids.forEach((id) => n.add(id))
      return n
    })

  /**
   * 💾 Сонголтыг хадгалах
   */
  const onSave = async () => {
    if (!orgType?.id) return
    setSaving(true)
    setError(null)
    try {
      await api.post('/orgtype/role', {
        type_id: orgType.id,
        role_ids: Array.from(selectedIds),
      } as Record<string, unknown>)
      onOpenChange(false)
    } catch {
      setError("Error occurred")
    } finally {
      setSaving(false)
    }
  }

  const allChecked = filtered.length > 0 && filtered.every((s) => selectedIds.has(s.id))
  const someChecked = !allChecked && filtered.some((s) => selectedIds.has(s.id))

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-4xl">
        <DialogHeader>
          <DialogTitle>
            {t('role')} — <span className="font-medium">{orgType?.name ?? ''}</span>
          </DialogTitle>
        </DialogHeader>

        {/* Controls */}
        <div className="mb-3 flex flex-wrap items-center gap-3">
          <Input
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder={`${t('search')}...`}
            className="h-9 w-full max-w-sm"
          />
          <label className="ml-auto flex items-center gap-2 text-sm select-none">
            <Checkbox
              checked={allChecked}
              aria-checked={someChecked ? 'mixed' : allChecked}
              onCheckedChange={toggleAllFiltered}
            />
            {t('all')}
          </label>
        </div>

        {/* Content */}
        {loading ? (
          <div className="flex items-center justify-center py-10 text-sm opacity-70">
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            {t('loading')}
          </div>
        ) : error ? (
          <div className="text-sm text-red-600">{error}</div>
        ) : filtered.length === 0 ? (
          <div className="py-10 text-center text-sm opacity-70">
            {t('no_information_available')}
          </div>
        ) : (
          <ul className="grid max-h-[50vh] grid-cols-1 gap-3 overflow-y-auto sm:grid-cols-2 lg:grid-cols-2">
            {filtered.map((r) => {
              const on = selectedIds.has(r.id)
              return (
                <div
                  key={r.id}
                  role="button"
                  onClick={() => toggle(r.id)}
                  className="w-full text-left"
                  aria-pressed={on}
                >
                  <Card
                    className={`group relative overflow-hidden rounded-xl border transition hover:shadow-md ${on ? 'border-primary/60 ring-2 ring-primary/20' : ''}`}
                  >
                    {/* selection badge */}
                    <div className="pointer-events-none absolute top-3 right-3 z-10 flex items-center justify-center rounded-full bg-white/90 p-1 shadow dark:bg-gray-800/90">
                      {on ? (
                        <CheckCircle2 className="h-5 w-5 text-primary" />
                      ) : (
                        <Checkbox checked={on} disabled className="h-5 w-5" />
                      )}
                    </div>

                    <div className="flex items-start gap-3 p-4">
                      <div className="bg-muted mt-0.5 rounded-lg p-2">
                        <UserCog className="h-5 w-5 opacity-70" />
                      </div>
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-2">
                          <div className="truncate font-medium">{r.name}</div>
                          {r.code ? (
                            <Badge variant="outline" className="text-[11px]">
                              {r.code}
                            </Badge>
                          ) : null}
                        </div>
                        <p className="text-muted-foreground mt-1 line-clamp-2 text-sm">
                          {r.description || '—'}
                        </p>
                        {r.system ? (
                          <p className="text-muted-foreground mt-1 text-xs">
                            {t('system')}: {r.system.name}
                          </p>
                        ) : null}
                      </div>
                    </div>
                  </Card>
                </div>
              )
            })}
          </ul>
        )}

        <DialogFooter className="mt-3">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t('cancel')}
          </Button>
          <Button onClick={onSave} disabled={saving}>
            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {t('save')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

