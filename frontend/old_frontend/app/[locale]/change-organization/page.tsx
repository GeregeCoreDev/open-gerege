/**
 * 🏢 Change Organization Page (/[locale]/change-organization/page.tsx)
 * 
 * Энэ нь байгууллага солих хуудас юм.
 * Зорилго: Хэрэглэгч харьяалагддаг байгууллагуудаас сонгож, контекст солих
 * 
 * Үйл ажиллагаа:
 * - Хэрэглэгчийн харьяалагддаг байгууллагуудыг жагсаана
 * - Байгууллага сонгоход API дуудаж, session солино
 * - Амжилттай бол хуудсыг дахин ачаална (window.location.reload)
 * - Одоо идэвхтэй байгууллагыг visual тэмдэглэнэ
 * 
 * UI Features:
 * - Grid layout (3 багана)
 * - Card дарагдах боломжтой
 * - Active байгууллагад ring border
 * - Building icon, нэр, регистрийн дугаар харуулна
 * 
 * Session Management:
 * - Backend API: POST /organization/change
 * - Response амжилттай бол хуудас reload хийнэ
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import * as React from 'react'
import { useTranslations } from 'next-intl'
import { useRouter } from '@/i18n/navigation'
import { useOrgStore } from '@/lib/stores/org'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { ArrowLeft, Building2, CheckCircle2 } from 'lucide-react'
import { cn } from '@/lib/utils'
import api from '@/lib/api'

export default function ChangeOrganizationPage() {
  const t = useTranslations()
  const router = useRouter()

  // 🔹 Байгууллагын store-оос мэдээлэл авах
  const organizations = useOrgStore((s) => s.organizations) as App.Organization[]
  const selectedOrganization = useOrgStore((s) => s.selectedOrganization) as
    | App.Organization
    | undefined
  const getOrganization = useOrgStore((s) => s.getOrganizations)
  const selectOrganization = useOrgStore((s) => s.selectOrg)

  /**
   * 🔄 Байгууллагуудын жагсаалт ачаалах
   * Хэрэв store-д байхгүй бол серверээс татна
   */
  React.useEffect(() => {
    if (!organizations || organizations.length === 0) {
      getOrganization?.().catch(() => {})
    }
  }, [organizations.length, getOrganization, organizations])

  // 🔹 Одоогийн идэвхтэй байгууллагын ID
  const activeId = selectedOrganization?.id

  /**
   * 🏢 Байгууллага солих handler
   * @param org - Сонгогдсон байгууллага
   */
  const changeOrganization = async (org: App.Organization) => {
    // ✅ Store-д байгууллагыг хадгална
    selectOrganization(org)

    // 🌐 Backend API дуудаж session солино
    const res = await api.post('/organization/change', {
      org_id: org.id,
    })

    if (res) {
      // 🔄 Амжилттай бол хуудсыг reload хийж, шинэ context-тэй ажиллана
      window.location.reload()
    }
  }

  /**
   * 🔙 Буцах функц
   */
  function goBack() {
    router.back()
  }

  return (
    <div className="bg-muted/20 h-screen w-screen">
      <div className="mx-auto flex h-full w-full max-w-5xl flex-col">
        {/* 🔙 Буцах товч */}
        <div className="flex justify-between pt-6">
          <Button variant="ghost" onClick={() => goBack()}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            {t('back')}
          </Button>
        </div>

        {/* 🏢 Байгууллагуудын жагсаалт */}
        <div className="flex h-full w-full flex-col items-center justify-center gap-6 pb-40">
          {/* 📝 Гарчиг */}
          <div className="space-y-2 text-center">
            <h1 className="text-4xl">{t('my_organization')}</h1>
            <p className="text-muted-foreground">
              Та ямар байгууллагаар хандахыг хүсч байгаагаа сонгоно уу
            </p>
          </div>

          {/* 🎴 Байгууллагуудын grid */}
          <div className="grid w-full grid-cols-3 gap-6">
            {(organizations ?? []).map((org) => {
              const isActive = activeId === org.id

              return (
                <Card
                  key={org.id}
                  onClick={() => changeOrganization(org)}
                  className={cn(
                    'group cursor-pointer border transition-all',
                    'hover:border-foreground/20 hover:shadow-md',
                    isActive && 'border-primary ring-primary/40 ring-2',
                  )}
                >
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      {/* 🔹 Байгууллагын мэдээлэл */}
                      <div className="flex items-center gap-3">
                        {/* 🎨 Building icon */}
                        <div
                          className={`flex h-12 w-12 items-center justify-center rounded-lg bg-yellow-100 dark:bg-yellow-900/20`}
                        >
                          <Building2 className={`h-6 w-6 text-yellow-600`} />
                        </div>

                        {/* 📝 Нэр, регистр */}
                        <div>
                          <CardTitle className="text-base leading-tight font-medium">
                            {org.name}
                          </CardTitle>
                          <CardDescription className="mt-0.5">
                            {t('reg_no', { defaultMessage: 'Reg. No' })}: {org.reg_no}
                          </CardDescription>
                        </div>
                      </div>
                      
                      {/* ✅ Идэвхтэй байгууллагад checkmark */}
                      {isActive && <CheckCircle2 className="text-primary h-5 w-5" />}
                    </div>
                  </CardHeader>
                </Card>
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
