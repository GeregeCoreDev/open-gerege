'use client'

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import api from '../api'

/**
 * 🧩 MenuState төрөл
 * Энэ store нь цэсийн жагсаалт болон сонгосон root menu-г удирдана.
 * @property menuList - Цэсийн жагсаалт (Menu tree)
 * @property selectedRootId - Сонгогдсон root menu ID
 * @property openGroups - Нээлттэй группууд
 * @property getMenuList - Цэсийн жагсаалт татах функц
 * @property selectRoot - Root сонгох функц
 * @property toggleGroup - Group нээх/хаах
 * @property clear - Store-г анхны төлөвт нь буцаах функц
 */
type MenuState = {
  menuList: App.Menu[]
  selectedRootId: number | null
  openGroups: Record<number, boolean>
  getMenuList: () => Promise<App.Menu[]>
  selectRoot: (rootId: number) => void
  findRootByPath: (path: string) => number | null
  getFirstChildPath: (rootId: number) => string | null
  toggleGroup: (groupId: number) => void
  setOpenGroups: (groups: Record<number, boolean>) => void
  clearSelectedRoot: () => void
  clear: () => void
}

/**
 * 🧱 useMenuStore — zustand store
 * Цэсийн жагсаалт болон сонгосон root menu-г хадгалах зориулалттай.
 * persist middleware ашиглан LocalStorage-д хадгална.
 */
export const useMenuStore = create<MenuState>()(
  persist(
    (set, get) => ({
      // 🔹 Цэсийн жагсаалт
      menuList: [],
      // 🔹 Сонгогдсон root menu ID
      selectedRootId: null,
      // 🔹 Нээлттэй группууд
      openGroups: {},

      /**
       * ⚙️ Цэсийн жагсаалт татах
       * API-аас menu-уудыг татаж, мод бүтэцтэй болгож буцаана.
       * @returns Promise<App.Menu[]> - цэсийн жагсаалт
       */
      getMenuList: async () => {
        // Session cookie байгаа эсэхийг шалгах
        if (typeof document !== 'undefined') {
          const hasSid = document.cookie.split('; ').some((row) =>
            row.startsWith('sid=') || row.startsWith('session=')
          )
          if (!hasSid) {
            // Session байхгүй бол API дуудахгүй, хоосон жагсаалт буцаана
            return []
          }
        }
        try {
          // const menus = menuList
          // set({ menuList: menus })
          const menus = await api.get<App.Menu[]>('/menu/my', { hasToast: false })
          set({ menuList: menus })

          // Хэрэв ямар ч root сонгогдоогүй бол эхний root-г автоматаар сонгоно
          const { selectedRootId } = get()
          if (selectedRootId == null && menus.length > 0) {
            const firstRoot = menus[0]
            const defaults: Record<number, boolean> = {}
            for (const child of firstRoot.children ?? []) {
              defaults[child.id] = true
            }
            set({
              selectedRootId: firstRoot.id,
              openGroups: defaults,
            })
          }

          return menus
        } catch {
          return []
        }
      },

      /**
       * 🔍 Path-аас root menu олох
       * @param path - Хайх path
       * @returns Root menu ID эсвэл null
       */
      findRootByPath: (path: string) => {
        const { menuList } = get()
        const normalizedPath =
          path.replace(/^\/[a-z]{2}(?:-[A-Z]{2})?(?=\/|$)/, '').replace(/\/+$/, '') || '/'

        // Menu tree-гээр хайх
        for (const root of menuList) {
          // Root-н path шалгах
          if (root.path && normalizedPath.startsWith(root.path)) {
            return root.id
          }

          // Children-д хайх
          for (const group of root.children ?? []) {
            for (const module of group.children ?? []) {
              if (module.path && normalizedPath.startsWith(module.path)) {
                return root.id
              }
            }
          }
        }

        return null
      },

      /**
       * 🎯 Root menu сонгох
       * @param rootId - Сонгож буй root menu ID
       */
      selectRoot: (rootId: number) => {
        const { menuList } = get()
        const root = menuList.find((r) => r.id === rootId)

        // Root-н children-г нээх
        const defaults: Record<number, boolean> = {}
        if (root) {
          for (const child of root.children ?? []) {
            defaults[child.id] = true
          }
        }

        set({
          selectedRootId: rootId,
          openGroups: defaults,
        })
      },

      /**
       * 🔍 Эхний child path олох (рекурсив)
       * Хамгийн доод level-н эхний child-н path-г олох
       * @param rootId - Root menu ID
       * @returns Хамгийн доод level-н эхний child-н path эсвэл null
       */
      getFirstChildPath: (rootId: number) => {
        const { menuList } = get()
        const root = menuList.find((r) => r.id === rootId)
        if (!root) return null

        // Рекурсив функц: хамгийн доод level-н эхний child-н path-г олох
        const findFirstPath = (menu: App.Menu): string | null => {
          // Хэрэв children байхгүй бол энэ menu-н path-г буцаана (хамгийн доод level)
          if (!menu.children || menu.children.length === 0) {
            return menu.path && menu.path !== '/' ? menu.path : null
          }
          // Хэрэв children байвал эхний child-д рекурсив хайна
          const firstChild = menu.children[0]
          const path = findFirstPath(firstChild)
          // Хэрэв child-д path олдохгүй бол энэ menu-н path-г буцаана
          return path || (menu.path && menu.path !== '/' ? menu.path : null)
        }

        return findFirstPath(root)
      },

      /**
       * 📂 Group нээх/хаах
       * @param groupId - Group ID
       */
      toggleGroup: (groupId: number) => {
        const { openGroups } = get()
        set({
          openGroups: {
            ...openGroups,
            [groupId]: !openGroups[groupId],
          },
        })
      },

      /**
       * 📂 Open groups тохируулах
       * @param groups - Groups object
       */
      setOpenGroups: (groups: Record<number, boolean>) => {
        set({ openGroups: groups })
      },

      /**
       * 🧹 Сонгогдсон root-г цэвэрлэх
       */
      clearSelectedRoot: () =>
        set({
          selectedRootId: null,
          openGroups: {},
        }),

      /**
       * 🧹 Store-г цэвэрлэх
       */
      clear: () =>
        set({
          menuList: [],
          selectedRootId: null,
          openGroups: {},
        }),
    }),
    {
      // 🗂 LocalStorage-д хадгалах нэр
      name: 'menu-store',
      // 💾 Persist хийхдээ зөвхөн menuList, selectedRootId, openGroups-г хадгална
      partialize: (s) => ({
        menuList: s.menuList,
        selectedRootId: s.selectedRootId,
        openGroups: s.openGroups,
      }),
    },
  ),
)
