// lib/stores/system.ts
'use client'

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import api from '@/lib/api'
import { useUserStore } from './user'

/**
 * 🧩 RoleState төрөл
 * Энэ store нь хэрэглэгчийн систем дэх дүрийн (role) мэдээллийг удирдана.
 * @property roleList - Хэрэглэгчийн бүх боломжит дүрүүдийн жагсаалт
 * @property selectedRole - Одоогоор сонгогдсон дүр
 * @property getRoleList - Серверээс дүрийн жагсаалт татах функц
 * @property selectRole - Сонгосон дүрийг солих функц
 * @property clear - Дүрийн мэдээллийг цэвэрлэх функц
 */
type RoleState = {
  roleList: App.UserRole[]
  selectedRole?: App.UserRole
  getRoleList: () => Promise<void>
  selectRole: (sys: App.UserRole) => Promise<void>
  clear: () => void
}

/**
 * 🧱 useRoleStore — zustand store
 * Хэрэглэгчийн дүрийн жагсаалт болон сонголтыг хадгалах, сэргээх зорилготой.
 * persist middleware ашиглаж байгаа тул LocalStorage-д хадгална.
 */
export const useRoleStore = create<RoleState>()(
  persist(
    (set, get) => ({
      // 🔹 Дүрийн жагсаалт
      roleList: [],
      // 🔹 Сонгогдсон дүр
      selectedRole: undefined,

      /**
       * ⚙️ Серверээс дүрийн жагсаалтыг татах
       * Хэрэв хэрэглэгчийн мэдээлэл байгаа бол API дуудна.
       * Алдаа гарвал консолд log хийнэ.
       */
      getRoleList: async () => {
        try {
          const { user_info } = useUserStore.getState()
          if (!user_info) return

          const res = await api.get<App.ListData<App.UserRole>>(
            `/role-matrix/roles?user_id=${user_info.id}`,
            { hasToast: false },
          )

          const items = res.items
          set({
            roleList: items,
          })

          // 🧠 Хэрэв ямар ч role сонгогдоогүй бол эхний role-г автоматаар сонгоно
          if (items.length > 0) {
            if (get().selectedRole === undefined) {
              await get().selectRole(items[0])
            }
          }
        } catch (e) {
          console.error('getRoleList failed:', e)
        }
      },

      /**
       * 🎯 Дүр сонгох
       * @param sys - Сонгож буй хэрэглэгчийн дүр
       */
      selectRole: async (sys) => {
        if (sys) {
          set({ selectedRole: sys })
        }
      },

      /**
       * 🧹 Store-г цэвэрлэх
       * roleList болон selectedRole-г анхны төлөвт нь буцаана
       */
      clear: () => set({ roleList: [], selectedRole: undefined }),
    }),
    {
      // 🗂 LocalStorage-д хадгалах нэр
      name: 'role-store',
      // 💾 Хадгалах үед зөвхөн тодорхой талбаруудыг хадгална
      partialize: (s) => ({ roleList: s.roleList, selectedRole: s.selectedRole }),
    },
  ),
)
