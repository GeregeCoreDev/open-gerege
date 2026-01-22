'use client'

import { useRef } from 'react'
import Image from 'next/image'
import { X } from 'lucide-react'
import { cn } from '@/lib/utils'

export type FileValue = {
  file?: File
  preview?: string | null
}

type Props = {
  value: FileValue | null
  onChange: (v: FileValue | null) => void
  label?: string
  className?: string
}

export default function FileUpload({ value, onChange, label, className }: Props) {
  const ref = useRef<HTMLInputElement>(null)

  function pickFile() {
    ref.current?.click()
  }

  function onFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return

    const preview = URL.createObjectURL(file)
    onChange({ file, preview })
  }

  function remove() {
    onChange(null)
  }

  return (
    <div className={cn('space-y-1', className)}>
      {label && <p className="text-sm font-medium">{label}</p>}

      <div
        className="bg-muted hover:bg-muted/80 relative aspect-square w-full cursor-pointer overflow-hidden rounded-xl border"
        onClick={pickFile}
      >
        {value?.preview ? (
          <Image src={value.preview} alt="preview" fill className="object-cover" />
        ) : (
          <div className="text-muted-foreground flex h-full w-full items-center justify-center text-xs">
            No image
          </div>
        )}

        {value && (
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation()
              remove()
            }}
            className="absolute top-1 right-1 rounded-full bg-black/60 p-1 text-white hover:bg-black/90"
          >
            <X size={14} />
          </button>
        )}
      </div>

      <input type="file" accept="image/*" className="hidden" ref={ref} onChange={onFile} />
    </div>
  )
}
