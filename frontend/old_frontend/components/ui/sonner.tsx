"use client"

import { useTheme } from "next-themes"
import { Toaster as Sonner, ToasterProps } from "sonner"
import { useEffect, useState } from "react"

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = "system" } = useTheme()
  const [position, setPosition] = useState<ToasterProps["position"]>("bottom-right")
  const [richColors, setRichColors] = useState(true)
  const [maxVisibleToasts, setMaxVisibleToasts] = useState(2)

  useEffect(() => {
    // Read notification settings from localStorage
    const readNotificationSettings = () => {
      try {
        const raw = localStorage.getItem("ui:theme")
        if (raw) {
          const themeState = JSON.parse(raw)
          if (themeState.notificationPosition) {
            setPosition(themeState.notificationPosition as ToasterProps["position"])
          }
          if (themeState.notificationMaxCount) {
            setMaxVisibleToasts(themeState.notificationMaxCount)
          }
          // Group style affects richColors
          if (themeState.notificationStyle === "expanded") {
            setRichColors(false)
          } else {
            setRichColors(true)
          }
        }
      } catch {}
    }

    readNotificationSettings()

    // Listen for storage changes
    const handleStorageChange = () => {
      readNotificationSettings()
    }
    window.addEventListener("storage", handleStorageChange)
    
    // Also listen for custom event from settings page
    const handleThemeUpdate = () => {
      readNotificationSettings()
    }
    window.addEventListener("theme-updated", handleThemeUpdate)

    return () => {
      window.removeEventListener("storage", handleStorageChange)
      window.removeEventListener("theme-updated", handleThemeUpdate)
    }
  }, [])

  return (
    <Sonner
      theme={theme as ToasterProps["theme"]}
      position={position}
      richColors={richColors}
      visibleToasts={maxVisibleToasts}
      className="toaster group"
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
        } as React.CSSProperties
      }
      {...props}
    />
  )
}

export { Toaster }
