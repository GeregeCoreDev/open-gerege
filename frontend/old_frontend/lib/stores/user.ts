'use client'

import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { useOrgStore } from './org'
import api from '../api'

/**
 * 🧩 UserState төрөл
 * Энэ store нь хэрэглэгчийн болон байгууллагын профайлын мэдээллийг удирдах зориулалттай.
 * @property user_info - Хэрэглэгчийн дэлгэрэнгүй мэдээлэл
 * @property org_info - Байгууллагын мэдээлэл (хэрэв хэрэглэгч байгууллагын төлөөлөгч бол)
 * @property user_name - Хэрэглэгчийн нэр (эсвэл байгууллагын нэр)
 * @property profile_image - Профайл зурагны URL
 * @property status - Ачааллын төлөв ('idle' | 'loading' | 'succeeded' | 'failed')
 * @property error - Алдааны мессеж
 * @property loadProfile - Хэрэглэгчийн профайлын мэдээлэл серверээс татах функц
 * @property clearAll - Хадгалагдсан бүх state-г цэвэрлэх функц
 */
type UserState = {
  user_info?: App.UserDetail
  org_info?: App.Organization
  is_org: boolean
  user_name?: string
  profile_image?: string
  status: 'idle' | 'loading' | 'succeeded' | 'failed'
  error?: string

  loadProfile: () => Promise<void>
  clearAll: () => void
}

/**
 * 🧱 useUserStore — zustand store
 * Энэ store нь хэрэглэгчийн мэдээлэл болон профайлын төлөвийг удирдаж, LocalStorage-д хадгалдаг.
 * persist middleware болон createJSONStorage ашиглан JSON хэлбэрээр хадгална.
 */
export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      // 🔹 Хэрэглэгчийн болон байгууллагын мэдээллийн анхны төлөв
      user_info: undefined,
      user_name: undefined,
      profile_image: undefined,
      is_org: false,
      org_info: undefined,
      status: 'idle',
      error: undefined,

      /**
       * ⚙️ Хэрэглэгчийн профайлыг серверээс ачаалах функц
       * - Session cookie байхгүй бол API дуудахгүй
       * - Давхар ачаалал (status === 'loading') эсвэл өмнө нь ачаалагдсан тохиолдолд дахин ажиллуулахгүй.
       */
      loadProfile: async () => {
        const currentStatus = get().status
        // Зөвхөн loading үед л дахин дуудахаас татгалзах
        if (currentStatus === 'loading') return

        // Session cookie байгаа эсэхийг шалгах (sid эсвэл session)
        if (typeof document !== 'undefined') {
          const hasSid = document.cookie.split('; ').some((row) =>
            row.startsWith('sid=') || row.startsWith('session=')
          )
          if (!hasSid) {
            // Session байхгүй бол API дуудахгүй
            return
          }
        }

        // Dev орчинд localStorage-д user_info байвал дахин дуудахаас татгалзах
        const isProduction = process.env.NODE_ENV === 'production'
        if (!isProduction) {
          const { user_info, org_info } = get()
          if (user_info || org_info) {
            if (currentStatus !== 'succeeded') {
              set({ status: 'succeeded' })
            }
            return
          }
        }

        set({ status: 'loading', error: undefined })
        try {
          const res = await api.get<App.UserProfileRes>('/me/profile', { hasToast: false })
          if (res) {
            // 🧍 Хувь хүн хэрэглэгчийн профайл
            if (res.is_org == false) {
              set({
                user_info: res.user,
                user_name: res.user.last_name[0] + '.' + res.user.first_name,
                profile_image: res.user.profile_img_url,
                status: 'succeeded',
                is_org: false,
              })

              // 🏢 Хэрэглэгчийн байгууллагуудыг ачаална
              await useOrgStore.getState().getOrganizations()
            } else {
              // 🏢 Байгууллагаар нэвтэрсэн тохиолдол
              set({
                org_info: res.org,
                user_name: res.org.name,
                profile_image: res.org.logo_image_url,
                user_info: undefined,
                status: 'succeeded',
                is_org: true,
              })
              // 🏢 Хэрэглэгчийн байгууллагуудыг store-оос устгана
              useOrgStore.getState().clear()

              // 🏢 Байгууллагын системүүдийг ачаална
            }
          }
        } catch (error) {
          // ❌ Алдаа гарсан үед төлөвийг "failed" болгож, алдааны мессеж хадгална
          const message = error instanceof Error ? error.message : 'Unknown error'
          set({ status: 'failed', error: message })
        }
      },

      /**
       * 🧹 Бүх state-г цэвэрлэх функц
       * Хэрэглэгчийн болон байгууллагын мэдээлэл, нэр, зураг, статус, алдааг анхны төлөвт буцаана.
       */
      clearAll: () =>
        set({
          user_info: undefined,
          org_info: undefined,
          user_name: undefined,
          profile_image: undefined,
          status: 'idle',
          error: undefined,
        }),
    }),
    {
      // 🗂 LocalStorage-д хадгалах нэр
      name: 'user-store',
      // 💾 JSON хэлбэрийн storage ашиглана
      storage: createJSONStorage(() => localStorage),
      // 🎯 Persist хийхдээ зөвхөн дараах талбаруудыг хадгална
      partialize: (s) => ({
        user_info: s.user_info,
        org_info: s.org_info,
        user_name: s.user_name,
        profile_image: s.profile_image,
        is_org: s.is_org,
      }),
    },
  ),
)
