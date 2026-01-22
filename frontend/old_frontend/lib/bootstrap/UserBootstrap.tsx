'use client'

import { useEffect, useRef } from 'react'
import { useUserStore } from '@/lib/stores/user'
import { useRoleStore } from '../stores/role'
import { useMenuStore } from '../stores/menu'

export default function UserBootstrap() {
  const { loadProfile } = useUserStore()
  const { getRoleList, roleList } = useRoleStore()
  const { getMenuList, selectRoot } = useMenuStore()
  const called = useRef(false)

  useEffect(() => {
    if (called.current) return
    called.current = true
    ;(window.requestIdleCallback ?? ((fn: () => void) => setTimeout(fn, 0)))(async () => {
      // Session cookie байгаа эсэхийг шалгах
      const hasSid = document.cookie.split('; ').some((row) =>
        row.startsWith('sid=') || row.startsWith('session=')
      )
      if (!hasSid) {
        // Session байхгүй бол API дуудахгүй
        return
      }

      // Profile load хийх
      await loadProfile().catch(() => {})

      // Profile амжилттай болсон эсэхийг шалгах - зөвхөн status-г шалгах (localStorage-ийн хуучин утгыг ашиглахгүй)
      const { status } = useUserStore.getState()
      if (status !== 'succeeded') {
        // Profile амжилтгүй бол бусад API дуудахгүй
        return
      }

      // Role list load хийх
      if (roleList.length === 0) {
        await getRoleList().catch(() => {})
      }

      // Menu жагсаалт авна
      const menus = await getMenuList()
      if (menus.length > 0) {
        // ✅ Эхний root menu-г сонгоно
        const firstMenu = menus[0]
        selectRoot(firstMenu.id)
      }
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return null
}
