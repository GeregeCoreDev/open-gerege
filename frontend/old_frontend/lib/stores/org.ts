'use client'

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { api } from '@/lib/api'
import { createLogger } from '@/lib/logger'

const logger = createLogger('OrgStore')

/**
 * 🧩 OrgState төрөл
 * Энэ store нь хэрэглэгчийн байгууллагын (organization) жагсаалт болон сонгосон байгууллагыг удирддаг.
 * @property organizations - Хэрэглэгчийн харьяалагддаг байгууллагуудын жагсаалт
 * @property selectedOrganization - Одоогоор сонгогдсон байгууллага
 * @property getOrganizations - Серверээс байгууллагуудын жагсаалт татах функц
 * @property clear - Байгууллагын мэдээллийг цэвэрлэх функц
 * @property selectOrg - Сонгогдсон байгууллагыг солих функц
 */
type OrgState = {
  organizations: App.Organization[]
  selectedOrganization?: App.Organization
  status: 'idle' | 'loading' | 'succeeded' | 'failed'
  error?: string

  getOrganizations: () => Promise<void>
  clear: () => void
  selectOrg: (org: App.Organization) => void
}

/**
 * 🧱 useOrgStore — zustand store
 * Хэрэглэгчийн байгууллагын жагсаалт болон сонголтыг хадгалж, persist middleware ашиглан LocalStorage-д хадгалдаг.
 */
export const useOrgStore = create<OrgState>()(
  persist(
    (set, get) => ({
      // 🔹 Байгууллагуудын жагсаалт
      organizations: [],
      // 🔹 Сонгогдсон байгууллага
      selectedOrganization: undefined,
      status: 'idle',
      error: undefined,

      /**
       * ⚙️ Серверээс байгууллагуудын мэдээлэл татах
       * Хэрэв хүсэлт амжилттай бол байгууллагуудын жагсаалт болон сонгогдсон байгууллагыг state-д хадгална.
       */
      getOrganizations: async () => {
        // Session cookie байгаа эсэхийг шалгах
        if (typeof document !== 'undefined') {
          const hasSid = document.cookie.split('; ').some((row) =>
            row.startsWith('sid=') || row.startsWith('session=')
          )
          if (!hasSid) {
            // Session байхгүй бол API дуудахгүй
            return
          }
        }
        try {
          const res = await api.get<App.UserOrganizationRes>('/me/organizations')
          if (res) {
            set({ organizations: res.items, status: 'succeeded' })
            if (res.org) {
              set({ selectedOrganization: res.org })
            } else {
              if (res.items.length > 0) {
                get().selectOrg(res.items[0])
              }
            }
          }
        } catch (error) {
          logger.error('Failed to fetch organizations', { error })
          set({
            status: 'failed',
            error: error instanceof Error ? error.message : 'Unknown error'
          })
        }
      },

      /**
       * 🧹 Store-г цэвэрлэх функц
       * Байгууллагын жагсаалт болон сонгогдсон байгууллагыг reset хийнэ.
       */
      clear: () => set({ organizations: [], selectedOrganization: undefined, status: 'idle', error: undefined }),

      /**
       * 🏢 Байгууллага сонгох функц
       * @param org - Сонгож буй байгууллага
       * Хэрэв шинэ байгууллага өмнөхөөсөө ялгаатай бол state-г шинэчилнэ.
       */
      selectOrg: async (org) => {
        if (org && org.id != get().selectedOrganization?.id) {
          set({ selectedOrganization: org })
        }
      },
    }),
    {
      // 🗂 LocalStorage-д хадгалах нэр
      name: 'orgs-store',
      // 💾 Persist хийх үед зөвхөн organizations талбарыг хадгална
      partialize: (s) => ({ organizations: s.organizations }),
    },
  ),
)
