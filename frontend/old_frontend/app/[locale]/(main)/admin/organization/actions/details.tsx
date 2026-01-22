/**
 * 🏢📋 Organization Details Dialog Component
 * (/[locale]/(main)/admin/organization/actions/details.tsx)
 * 
 * Энэ нь байгууллагын дэлгэрэнгүй мэдээлэл харуулах dialog юм.
 * Зорилго: Байгууллагын бүрэн мэдээллийг харуулах (read-only)
 * 
 * Features:
 * - ✅ Full organization details display
 * - ✅ Contact information (phone, email, address, website)
 * - ✅ Organization type badge
 * - ✅ Active/Inactive status
 * - ✅ Allowed systems list
 * - ✅ Skeleton loading states
 * - ✅ Responsive layout
 * - ✅ Icon-based sections
 * - ✅ Clickable links (email, website)
 * 
 * Props:
 * @param open - Dialog visibility
 * @param onOpenChange - Toggle dialog
 * @param orgId - Organization ID to load
 * 
 * Display Sections:
 * - Header: Name, ID, Type badge, Status
 * - Description: Organization description
 * - Contact: Phone, Email, Address, Website
 * - Systems: Allowed/access systems grid
 * - Statistics: Placeholder for future metrics
 * 
 * Data Loading:
 * - Fetches on open + orgId change
 * - Auto cleanup on unmount
 * - Flexible field mapping (backend variations)
 * 
 * API:
 * - GET /organization/:id - Fetch org details
 * 
 * Related Types:
 * - App.Organization (base)
 * - OrgDetail (extended with optional fields)
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

// app/[locale]/organization/actions/detail.tsx
'use client'

import * as React from 'react'
import { useEffect, useMemo, useState } from 'react'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'

import { Building2, Phone, Mail, MapPin, Globe, Database } from 'lucide-react'

/**
 * 🧩 Component Props
 */
type Props = {
  open: boolean
  onOpenChange: (v: boolean) => void
  /** Хүссэн ID-гаа дамжуул; нээгдэхэд /organization/:id руугаас татна */
  orgId: number | null | undefined
}

/**
 * 🧩 Extended Organization Detail Type
 * Backend-аас ирж болох өөр өөр field нэршлүүдийг хамруулна
 */
type OrgDetail = App.Organization & {
  systems?: App.System[]
  access_systems?: App.System[]
  allowed_systems?: App.System[]
  website: string
  description: string
  // Backend альтернатив field нэршлүүд
  phone?: string
  contact_email?: string
  web_site?: string
  about?: string
}

export default function OrganizationDetailDialog({ open, onOpenChange, orgId }: Props) {
  const t = useTranslations()

  // 📊 State management
  const [loading, setLoading] = useState(false)
  const [_error, setError] = useState<string | null>(null)
  const [data, setData] = useState<OrgDetail | null>(null)

  /**
   * 📥 Dialog нээгдэх эсвэл ID солигдоход мэдээлэл ачаална
   */
  useEffect(() => {
    if (!open || !orgId) {
      setData(null)
      setError(null)
      return
    }

    let canceled = false
    ;(async () => {
      setLoading(true)
      setError(null)
      try {
        const res = await api.get<OrgDetail>(`/organization/${orgId}`)
        if (canceled) return

        // талбаруудыг жигдрүүлж авъя
        const d: OrgDetail = {
          ...res,
          phone_no: res.phone_no ?? res.phone ?? '',
          email: res.email ?? res.contact_email ?? '',
          website: res.website ?? res.web_site ?? '',
          description: res.description ?? res.about ?? '',
        }
        setData(d)
      } catch {
        setError("Error occurred")
      } finally {
        if (!canceled) setLoading(false)
      }
    })()

    return () => {
      canceled = true
    }
  }, [open, orgId])

  /**
   * 🖥️ Системүүдийн жагсаалт (flexible field mapping)
   */
  const systems: App.System[] = useMemo(() => {
    if (!data) return []
    return data.systems ?? data.access_systems ?? data.allowed_systems ?? ([] as App.System[])
  }, [data])

  /**
   * 🎨 Status badge component
   */
  const StatusChip = ({ active }: { active?: boolean }) => (
    <Badge variant={active ? 'default' : 'secondary'}>{active ? t('active') : t('inactive')}</Badge>
  )

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[85vh] overflow-hidden p-0 sm:max-w-2xl">
        {/* Header area */}
        <div className="flex items-start gap-3 p-4 sm:p-5">
          <div className="bg-muted text-muted-foreground flex h-12 w-12 items-center justify-center rounded-xl">
            <Building2 className="h-6 w-6" />
          </div>
          <div className="min-w-0 flex-1">
            <DialogHeader className="p-0">
              <DialogTitle className="truncate">
                {data?.name ?? <Skeleton className="h-6 w-40" />}
              </DialogTitle>
              <DialogDescription className="mt-1 flex flex-wrap items-center gap-2">
                {/* Код */}
                {loading ? (
                  <Skeleton className="h-5 w-20" />
                ) : (
                  data?.id && <Badge variant="outline">ID: {data.id}</Badge>
                )}
                {/* Төрөл */}
                {loading ? (
                  <Skeleton className="h-5 w-24" />
                ) : data?.type?.name ? (
                  <Badge variant="secondary">{data.type.name}</Badge>
                ) : null}
                {/* Статус */}
                {!loading && <StatusChip active={(data)?.is_active !== false} />}
              </DialogDescription>
            </DialogHeader>
          </div>
        </div>

        <Separator />

        {/* Body scrollable */}
        <div className="max-h-[65vh] overflow-auto px-4 pb-4 sm:px-5 sm:pb-5">
          {/* Description */}
          <section className="py-4">
            <div className="mb-2 flex items-center gap-2">
              <span className="text-muted-foreground">
                {/* simple icon via emoji-style for light weight */}
                📝
              </span>
              <h3 className="text-base font-medium">{t('description')}</h3>
            </div>
            {loading ? (
              <div className="space-y-2">
                <Skeleton className="h-4 w-5/6" />
                <Skeleton className="h-4 w-2/3" />
              </div>
            ) : (
              <p className="text-muted-foreground">{data?.description || '—'}</p>
            )}
          </section>

          <Separator />

          {/* Contacts */}
          <section className="py-4">
            <h3 className="mb-3 text-base font-medium">{t('contact')}</h3>
            <div className="grid gap-3 sm:grid-cols-2">
              <ContactRow
                icon={<Phone className="h-4 w-4" />}
                label={t('phone_no')}
                value={loading ? null : data?.phone_no || '—'}
              />
              <ContactRow
                icon={<Mail className="h-4 w-4" />}
                label={t('email')}
                value={
                  loading ? null : data?.email ? (
                    <a
                      href={`mailto:${data.email}`}
                      className="text-primary underline-offset-2 hover:underline"
                    >
                      {data.email}
                    </a>
                  ) : (
                    '—'
                  )
                }
              />
              <ContactRow
                icon={<MapPin className="h-4 w-4" />}
                label={t('address')}
                value={loading ? null : (data)?.address || '—'}
              />
              <ContactRow
                icon={<Globe className="h-4 w-4" />}
                label="Website"
                value={
                  loading ? null : data?.website ? (
                    <a
                      href={data.website}
                      target="_blank"
                      rel="noreferrer"
                      className="text-primary break-all underline-offset-2 hover:underline"
                    >
                      {data.website}
                    </a>
                  ) : (
                    '—'
                  )
                }
              />
            </div>
          </section>

          <Separator />

          {/* Allowed systems */}
          <section className="py-4">
            <h3 className="mb-3 text-base font-medium">{t('systems') || 'Системүүд'}</h3>

            {loading ? (
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                {Array.from({ length: 3 }).map((_, i) => (
                  <Skeleton key={i} className="h-20 w-full rounded-xl" />
                ))}
              </div>
            ) : systems.length === 0 ? (
              <p className="text-muted-foreground">{t('no_information_available')}</p>
            ) : (
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                {systems.map((s) => (
                  <div key={s.id} className="hover:bg-muted/40 rounded-xl border p-3 transition">
                    <div className="mb-2 flex items-center gap-2">
                      <div className="bg-muted text-muted-foreground flex h-8 w-8 items-center justify-center rounded-lg">
                        <Database className="h-4 w-4" />
                      </div>
                      <div className="min-w-0">
                        <div className="truncate font-medium">{s.name}</div>
                        <div className="text-muted-foreground mt-0.5 text-xs">
                          {t('code')}: <span className="font-mono">{s.code ?? s.id}</span>
                        </div>
                      </div>
                    </div>
                    {s.description ? (
                      <p className="text-muted-foreground line-clamp-2 text-sm">{s.description}</p>
                    ) : null}
                  </div>
                ))}
              </div>
            )}
          </section>

          {/* (optional) Statistics placeholder */}
          <Separator />
          <section className="py-4">
            <h3 className="mb-2 text-base font-medium">'Статистик'</h3>
            <p className="text-muted-foreground text-sm">{t('no_information_available')}</p>
          </section>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-2 border-t p-3">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t('close')}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}

/**
 * 🧱 Contact Row Helper Component
 * Icon + Label + Value харуулах reusable component
 */
function ContactRow({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode
  label: string
  value: React.ReactNode | null
}) {
  return (
    <div className="rounded-lg border p-3">
      <div className="mb-1 flex items-center gap-2 text-sm">
        <span className="text-muted-foreground">{icon}</span>
        <span className="opacity-70">{label}</span>
      </div>
      {value === null ? <Skeleton className="h-4 w-40" /> : <div className="text-sm">{value}</div>}
    </div>
  )
}
