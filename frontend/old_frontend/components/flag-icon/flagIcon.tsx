'use client'
import Image from 'next/image'

export default function FlagIcon({ currentLanguage }: { currentLanguage: { code: string } }) {
  const src = `/flag/${currentLanguage.code}.png`

  return <Image src={src} alt={`${currentLanguage.code} flag`} width={16} height={16} unoptimized/>
}
