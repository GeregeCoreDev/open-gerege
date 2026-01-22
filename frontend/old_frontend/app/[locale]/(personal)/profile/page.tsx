/**
 * 👤 Profile Page (/[locale]/(personal)/profile/page.tsx)
 * 
 * Энэ нь хэрэглэгчийн профайл хуудас юм.
 * Зорилго: Хэрэглэгчийн болон байгууллагын дэлгэрэнгүй мэдээлэл харуулах
 * 
 * Features:
 * - Dual mode: Хувь хүн болон байгууллагын профайл харуулна
 * - is_org flag-аар ялгана (хувь хүн эсвэл байгууллага)
 * - Profile зураг эсвэл initials харуулна
 * - Бүлэг бүрд мэдээлэл: Contact, Personal, Address
 * - Organizations: Хэрэглэгчийн бүх байгууллагуудын жагсаалт
 * - Responsive grid layout
 * - Read-only mode (засах функц байхгүй)
 * 
 * Store Dependencies:
 * - useUserStore: user_info, org_info, is_org
 * - useOrgStore: organizations list
 * 
 * Хувь хүний мэдээлэл:
 * - Регистр, нэр, утас, имэйл
 * - Хүйс, төрсөн огноо, иргэншил
 * - Хаяг (аймаг, сум, баг, дэлгэрэнгүй)
 * - Байгууллагуудын жагсаалт
 * 
 * Байгууллагын мэдээлэл:
 * - Регистр, нэр, утас, имэйл
 * - Байгууллагын төрөл
 * - Хаяг мэдээлэл
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import * as React from 'react'
import { useTranslations } from 'next-intl'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Loader2, Building2, ChevronRight, Check } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useUserStore } from '@/lib/stores/user'
import { useOrgStore } from '@/lib/stores/org'

export default function ProfilePage() {
  const t = useTranslations()
  
  // 🔹 User store-оос мэдээлэл авах
  const { user_info, org_info, is_org } = useUserStore()
  const { organizations, selectedOrganization, selectOrg } = useOrgStore()

  /**
   * 🚻 Хүйсний текст орчуулга
   */
  const genderText = (g?: number) =>
    g === 1
      ? t('male', { defaultMessage: 'Эр' })
      : g === 2
        ? t('female', { defaultMessage: 'Эм' })
        : t('unknown', { defaultMessage: 'Тодорхойгүй' })

  /**
   * 📅 Огноо форматлах
   */
  const fmtDate = (s?: string) => (s ? s.slice(0, 10) : '—')

  // 🏢 Байгууллагын профайл харуулах
  if (is_org) {
    return (
      <div className="flex h-full w-full items-center justify-center p-4 sm:p-6">
        <Card className="relative flex w-full max-w-5xl flex-col overflow-hidden border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900">
          <CardHeader className="flex flex-shrink-0 flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <CardTitle className="text-xl">
                {t('organization_profile', { defaultMessage: 'Байгууллагын мэдээлэл' })}
              </CardTitle>
              <CardDescription>
                {t('organization_details', { defaultMessage: 'Байгууллагын ерөнхий мэдээлэл' })}
              </CardDescription>
            </div>
          </CardHeader>

          <Separator className="flex-shrink-0" />

          <CardContent className="max-h-[70vh] overflow-y-auto p-4 sm:p-6">
            {!org_info ? (
              <div className="flex h-40 items-center justify-center">
                <Loader2 className="h-5 w-5 animate-spin" />
              </div>
            ) : (
              <div className="flex flex-col gap-6">
                {/* Header / Identity */}
                <section className="flex flex-col items-start gap-4 sm:flex-row sm:items-center">
                  <div
                    className={cn(
                      'h-24 w-24 overflow-hidden rounded-xl ring-1 ring-black/5',
                      'flex items-center justify-center bg-gray-100 dark:bg-gray-800',
                    )}
                  >
                    {org_info.logo_image_url ? (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img
                        src={org_info.logo_image_url}
                        alt="logo"
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <span className="text-2xl font-semibold capitalize">
                        {org_info.short_name?.[0] || org_info.name?.[0] || 'O'}
                      </span>
                    )}
                  </div>

                  <div className="space-y-1">
                    <h2 className="text-2xl font-semibold">{org_info.name}</h2>
                    <p className="text-muted-foreground">
                      {t('reg_no')}:{' '}
                      <span className="font-medium uppercase">{org_info.reg_no || '—'}</span>
                    </p>
                  </div>
                </section>

                <Separator />

                {/* Contact */}
                <section>
                  <h3 className="mb-3 text-base font-medium">
                    {t('contact_information') ?? 'Холбоо барих'}
                  </h3>
                  <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                    <Field label={t('phone_no')} value={org_info.phone_no || '—'} />
                    <Field label={t('email')} value={org_info.email || '—'} />
                  </div>
                </section>

                {/* Organization Type */}
                <section>
                  <h3 className="mb-3 text-base font-medium">
                    {t('organization_type', { defaultMessage: 'Байгууллагын төрөл' })}
                  </h3>
                  <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                    <Field
                      label={t('type', { defaultMessage: 'Төрөл' })}
                      value={org_info.type?.name || '—'}
                    />
                    <Field
                      label={t('parent_org', { defaultMessage: 'Эцэг байгууллага' })}
                      value={org_info.parent_id ? String(org_info.parent_id) : '—'}
                    />
                  </div>
                </section>

                {/* Address */}
                <section>
                  <h3 className="mb-3 text-base font-medium">{t('address') ?? 'Хаяг'}</h3>
                  <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                    <Field
                      label={t('aimag', { defaultMessage: 'Аймаг/Нийслэл' })}
                      value={org_info.aimag_name || '—'}
                    />
                    <Field
                      label={t('sum', { defaultMessage: 'Сум/Дүүрэг' })}
                      value={org_info.sum_name || '—'}
                    />
                    <Field
                      label={t('bag', { defaultMessage: 'Баг/Хороо' })}
                      value={org_info.bag_name || '—'}
                    />
                    <Field
                      label={t('address_detail', { defaultMessage: 'Дэлгэрэнгүй' })}
                      value={org_info.address_detail || '—'}
                      className="sm:col-span-3"
                    />
                  </div>
                </section>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    )
  }

  // 👤 Хувийн profile (өмнөх код чинь хэвээр)
  return (
    <div className="flex h-full w-full items-center justify-center p-4 sm:p-6">
      <Card className="relative flex w-full max-w-5xl flex-col overflow-hidden border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900">
        <CardHeader className="flex flex-shrink-0 flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle className="text-xl">{t('profile') ?? 'Profile'}</CardTitle>
            <CardDescription>
              {t('your_personal_information', { defaultMessage: 'Таны хувийн мэдээлэл' })}
            </CardDescription>
          </div>
        </CardHeader>

        <Separator className="flex-shrink-0" />

        <CardContent className="max-h-[70vh] overflow-y-auto p-4 sm:p-6">
          {!user_info ? (
            <div className="flex h-40 items-center justify-center">
              <Loader2 className="h-5 w-5 animate-spin" />
            </div>
          ) : (
            <div className="flex flex-col gap-6">
              {/* Header / Identity */}
              <section className="flex flex-col items-start gap-4 sm:flex-row sm:items-center">
                <div
                  className={cn(
                    'h-24 w-24 overflow-hidden rounded-xl ring-1 ring-black/5',
                    'flex items-center justify-center bg-gray-100 dark:bg-gray-800',
                  )}
                >
                  {user_info.profile_img_url ? (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img
                      src={user_info.profile_img_url}
                      alt="profile"
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <span className="text-2xl font-semibold capitalize">
                      {user_info.first_name?.[0]}
                      {user_info.last_name?.[0]}
                    </span>
                  )}
                </div>

                <div className="space-y-1">
                  <h2 className="text-2xl leading-tight font-semibold capitalize">
                    {user_info.last_name
                      ? `${user_info.last_name} ${user_info.first_name}`
                      : user_info.first_name}
                  </h2>
                  <p className="text-muted-foreground">
                    {t('reg_no')}:{' '}
                    <span className="font-medium uppercase">{user_info.reg_no || '—'}</span>
                  </p>
                </div>
              </section>

              <Separator />

              {/* Contact */}
              <section>
                <h3 className="mb-3 text-base font-medium">
                  {t('contact_information') ?? 'Холбоо барих'}
                </h3>
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                  <Field label={t('phone_no')} value={user_info.phone_no || '—'} />
                  <Field label={t('email')} value={user_info.email || '—'} />
                </div>
              </section>

              {/* Personal */}
              <section>
                <h3 className="mb-3 text-base font-medium">
                  {t('personal_information') ?? 'Хувийн мэдээлэл'}
                </h3>
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                  <Field label={t('gender')} value={genderText(user_info.gender)} />
                  <Field label={t('birth_date')} value={fmtDate(user_info.birth_date)} />
                  <Field
                    label={t('nationality', { defaultMessage: 'Үндэс угсаа' })}
                    value={user_info.nationality || '—'}
                  />
                  <Field
                    label={t('country', { defaultMessage: 'Улс' })}
                    value={user_info.country_name || '—'}
                  />
                  <Field label="Civil ID" value={String(user_info.civil_id ?? '—')} />
                </div>
              </section>

              {/* Address */}
              <section>
                <h3 className="mb-3 text-base font-medium">{t('address') ?? 'Хаяг'}</h3>
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                  <Field
                    label={t('aimag', { defaultMessage: 'Аймаг/Нийслэл' })}
                    value={user_info.aimag_name || '—'}
                  />
                  <Field
                    label={t('sum', { defaultMessage: 'Сум/Дүүрэг' })}
                    value={user_info.sum_name || '—'}
                  />
                  <Field
                    label={t('bag', { defaultMessage: 'Баг/Хороо' })}
                    value={user_info.bag_name || '—'}
                  />
                  <Field
                    label={t('address_detail', { defaultMessage: 'Дэлгэрэнгүй' })}
                    value={user_info.address_detail || '—'}
                    className="sm:col-span-3"
                  />
                </div>
              </section>

              {/* Organizations */}
              <section>
                <h3 className="mb-3 text-base font-medium">
                  {t('my_organizations')}
                </h3>
                {organizations.length === 0 ? (
                  <div className="flex flex-col items-center justify-center rounded-md border border-dashed p-6 text-center">
                    <Building2 className="h-10 w-10 text-gray-400" />
                    <p className="mt-2 text-sm text-gray-500">{t('no_organizations')}</p>
                  </div>
                ) : (
                  <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                    {organizations.map((org) => {
                      const isCurrent = selectedOrganization?.id === org.id
                      return (
                        <div
                          key={org.id}
                          className={cn(
                            'relative flex items-center gap-3 rounded-lg border p-4 transition-all',
                            isCurrent
                              ? 'border-primary-200 bg-primary-50 dark:border-primary-800 dark:bg-primary-900/20'
                              : 'border-gray-200 hover:border-gray-300 dark:border-gray-700 dark:hover:border-gray-600',
                          )}
                        >
                          {/* Logo / Initial */}
                          <div
                            className={cn(
                              'flex h-12 w-12 shrink-0 items-center justify-center rounded-lg',
                              'bg-gray-100 dark:bg-gray-800',
                            )}
                          >
                            {org.logo_image_url ? (
                              // eslint-disable-next-line @next/next/no-img-element
                              <img
                                src={org.logo_image_url}
                                alt={org.name}
                                className="h-full w-full rounded-lg object-cover"
                              />
                            ) : (
                              <Building2 className="h-6 w-6 text-gray-500" />
                            )}
                          </div>

                          {/* Info */}
                          <div className="min-w-0 flex-1">
                            <div className="flex items-center gap-2">
                              <p className="truncate font-medium">{org.name}</p>
                              {isCurrent && (
                                <Badge
                                  variant="secondary"
                                  className="bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300"
                                >
                                  <Check className="mr-1 h-3 w-3" />
                                  {t('current_organization')}
                                </Badge>
                              )}
                            </div>
                            <p className="text-muted-foreground text-sm">
                              {org.reg_no} • {org.type?.name || '—'}
                            </p>
                          </div>

                          {/* Switch Button */}
                          {!isCurrent && (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => selectOrg(org)}
                              className="shrink-0"
                            >
                              {t('switch_to_organization')}
                              <ChevronRight className="ml-1 h-4 w-4" />
                            </Button>
                          )}
                        </div>
                      )
                    })}
                  </div>
                )}
              </section>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

/**
 * 🧱 Reusable талбар компонент
 * Label болон value харуулах энгийн талбар
 * @param label - Талбарын гарчиг
 * @param value - Талбарын утга
 * @param className - Нэмэлт CSS class
 */
function Field({
  label,
  value,
  className,
}: {
  label: string
  value: React.ReactNode
  className?: string
}) {
  return (
    <div className={cn('min-w-0 rounded-md border p-3', className)}>
      <div className="text-muted-foreground text-xs">{label}</div>
      <div className="mt-1 text-sm break-words">{value}</div>
    </div>
  )
}
