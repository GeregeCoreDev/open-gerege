'use client'

import * as React from 'react'
import { Loader2, Search, X } from 'lucide-react'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'

type Props = {
  onChange?: (u: App.User | null) => void
  defaultValue?: string
  autoFocus?: boolean
  className?: string
}

/** Монгол РД: 2 монгол үсэг + 8 тоо */
function normalizeRegNo(input: string) {
  const up = (input || '').trim().toUpperCase()
  return up.replace('Ө', 'Ө').replace('Ү', 'Ү')
}
function isValidRegNo(reg: string) {
  return /^[А-ЯӨҮ]{2}\d{8}$/.test(reg)
}

export default function UserFind({ onChange, defaultValue = '', autoFocus, className }: Props) {
  const t = useTranslations()

  const [regNo, setRegNo] = React.useState(normalizeRegNo(defaultValue))
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [user, setUser] = React.useState<App.User | null>(null)

  const canSearch = isValidRegNo(regNo)

  const doSearch = async () => {
    if (!canSearch) {
      setError(t('invalid') || 'Invalid')
      return
    }
    setLoading(true)
    setError(null)
    setUser(null)
    try {
      const res = await api.post<App.User>('/user/find-from-core', { search_text: regNo })
      // Шууд яг тухайн хүний объект гэж үзсэн
      setUser(res)
      onChange?.(res)
    } catch (error) {
      const e = error as Error
      setError(e?.message || t('not_found') || 'Not found')
      setUser(null)
      onChange?.(null)
    } finally {
      setLoading(false)
    }
  }

  const onInput = (v: string) => {
    const up = normalizeRegNo(v)
    setRegNo(up)
    setError(null)
    setUser(null)
    onChange?.(null)
  }

  const clearAll = () => {
    setRegNo('')
    setError(null)
    setUser(null)
    onChange?.(null)
  }

  return (
    <div className={className}>
      <div className="flex items-center gap-2">
        <div className="relative flex-1">
          <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
          <Input
            value={regNo}
            onChange={(e) => onInput(e.target.value)}
            placeholder="AA12345678"
            className="h-9 pl-8"
            autoFocus={autoFocus}
            inputMode="text"
            maxLength={10}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && canSearch && !loading) doSearch()
            }}
          />
        </div>
        <Button onClick={doSearch} disabled={!canSearch || loading}>
          {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Search className="h-4 w-4" />}
        </Button>
        {regNo && (
          <Button
            variant="ghost"
            size="icon"
            aria-label="clear"
            onClick={clearAll}
            className="h-9 w-9"
          >
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>

      {/* Хэрэглэгчийн preview */}
      {error && (
        <div className="mt-2 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
          {error}
        </div>
      )}

      {user && (
        <Card className="mt-3">
          <CardContent className="p-4">
            <div className="grid gap-2 sm:grid-cols-2">
              <div>
                <p className="text-muted-foreground text-xs">{t('reg_no')}</p>
                <p className="font-medium">{user.reg_no || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground text-xs">{t('name')}</p>
                <p className="font-medium capitalize">
                  {[user.family_name, user.last_name, user.first_name].filter(Boolean).join(' ') ||
                    '—'}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground text-xs">{t('email')}</p>
                <p className="font-medium">{user.email || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground text-xs">{t('phone_no')}</p>
                <p className="font-medium">{user.phone_no || '—'}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Жижиг зөвлөмж */}
      {!error && !user && regNo && !canSearch && (
        <div className="mt-2 text-xs text-amber-600">
          РД 10 оронтой: эхний 2 нь монгол үсэг, дараагийн 8 нь тоо (жишээ: АБ12345678)
        </div>
      )}
    </div>
  )
}
