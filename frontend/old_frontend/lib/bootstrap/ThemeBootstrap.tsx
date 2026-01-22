'use client'

import { useEffect } from 'react'

/**
 * 🎨 Theme Bootstrap
 * 
 * App эхлэх үед default theme state-аас CSS variable-уудыг тохируулна.
 * Энэ нь settings хуудас руу ороогүй тохиолдолд ч primary color shade-ууд
 * зөв ажиллах боломжийг олгоно.
 */

const STORAGE_KEY = 'ui:theme'

const DEFAULT_PRIMARY_COLOR = '#3b82f6'

// Helper function to generate color shades (same as in settings/page.tsx)
function generateColorShades(hex: string) {
  // Convert hex to RGB
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)

  // Generate shades
  const shades: Record<string, string> = {}
  
  // Lighter shades (50-400)
  const lightShades = [
    { name: '50', mix: 0.95 },
    { name: '100', mix: 0.9 },
    { name: '200', mix: 0.8 },
    { name: '300', mix: 0.7 },
    { name: '400', mix: 0.6 },
  ]
  
  lightShades.forEach(({ name, mix }) => {
    const newR = Math.round(r + (255 - r) * mix)
    const newG = Math.round(g + (255 - g) * mix)
    const newB = Math.round(b + (255 - b) * mix)
    shades[name] = `#${newR.toString(16).padStart(2, '0')}${newG.toString(16).padStart(2, '0')}${newB.toString(16).padStart(2, '0')}`
  })
  
  // Base color (500)
  shades['500'] = hex
  
  // Darker shades (600-900)
  const darkShades = [
    { name: '600', mix: 0.8 },
    { name: '700', mix: 0.6 },
    { name: '800', mix: 0.4 },
    { name: '900', mix: 0.2 },
  ]
  
  darkShades.forEach(({ name, mix }) => {
    const newR = Math.round(r * mix)
    const newG = Math.round(g * mix)
    const newB = Math.round(b * mix)
    shades[name] = `#${newR.toString(16).padStart(2, '0')}${newG.toString(16).padStart(2, '0')}${newB.toString(16).padStart(2, '0')}`
  })
  
  return shades
}

export default function ThemeBootstrap() {
  useEffect(() => {
    // Read theme state from localStorage
    let primaryColor = DEFAULT_PRIMARY_COLOR
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      if (raw) {
        const themeState = JSON.parse(raw)
        if (themeState.primaryColor) {
          primaryColor = themeState.primaryColor
        }
      }
    } catch {
      // Use default if parsing fails
    }

    // Generate and apply primary color shades
    const shades = generateColorShades(primaryColor)
    const root = document.documentElement
    
    // Apply primary color
    root.style.setProperty('--primary-hex', primaryColor)
    root.style.setProperty('--primary', primaryColor)
    
    // Apply all shade variables
    Object.entries(shades).forEach(([shade, color]) => {
      root.style.setProperty(`--primary-${shade}`, color)
    })
  }, [])

  return null
}

