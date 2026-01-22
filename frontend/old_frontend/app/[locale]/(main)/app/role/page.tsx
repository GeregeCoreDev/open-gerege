/**
 * 🛡️ App Role Page (/[locale]/(main)/app/role/page.tsx)
 * 
 * Энэ нь App системийн дүр удирдах хуудас юм.
 * Зорилго: App системийн дүрүүдийн удирдлага
 * 
 * Implementation:
 * - Reusable component: RolesBySystem
 * - Automatically filters roles by current system
 * - Shares logic with other subsystems
 * 
 * Related:
 * - Business/Role page (same component)
 * - TPay/Role page (same component)
 * - Component: @/components/common/subSystemRolePage
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

import RolesBySystem from '@/components/common/subSystemRolePage'

export default async function TPayRolePage() {
  return <RolesBySystem />
}
