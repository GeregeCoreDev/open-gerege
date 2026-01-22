'use client'
import * as React from 'react'
import { toast } from 'sonner'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import api from '@/lib/api'
import { encryptCardNumber } from '../actions/encrypt'

type NegdiCard = { tranid: string; checkid: string; maskedpan: string }
type UserCard = { id: number; card_number: string }

const MONTHS = Array.from({ length: 12 }, (_, i) => String(i + 1).padStart(2, '0'))
const YEARS = Array.from({ length: 10 }, (_, i) => {
  const full = new Date().getFullYear() + i
  return { label: String(full), value: String(full).slice(2) }
})

export default function AddCard({ onClose }: { onClose: (selected?: UserCard) => void }) {
  const [step, setStep] = React.useState<1 | 2>(1)
  const [num, setNum] = React.useState('') // plain digits (16)
  const [numUi, setNumUi] = React.useState('') // **** **** ... (19)
  const [mm, setMM] = React.useState<string>('')
  const [yy, setYY] = React.useState<string>('')
  const [cvv, setCVV] = React.useState('')
  const [negdi, setNegdi] = React.useState<NegdiCard | null>(null)
  const [loading, setLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  const onCardInput = (v: string) => {
    const d = v.replace(/\D/g, '').slice(0, 16)
    setNum(d)
    setNumUi(d.replace(/(.{4})/g, '$1 ').trim())
    setError(null)
  }

  const check = async () => {
    if (num.length !== 16 || !mm || !yy) {
      setError('Картын мэдээлэл бүрэн бөглөнө үү')
      return
    }

    // Validate expiry date
    const currentYear = new Date().getFullYear() % 100
    const currentMonth = new Date().getMonth() + 1
    const expYear = parseInt(yy, 10)
    const expMonth = parseInt(mm, 10)

    if (expYear < currentYear || (expYear === currentYear && expMonth < currentMonth)) {
      setError('Картын хугацаа дууссан байна')
      return
    }

    setLoading(true)
    setError(null)

    try {
      // Encrypt card number on server side
      const encryptedCard = await encryptCardNumber(num)

      const res = await api.post<NegdiCard>('/card/negdi/create', {
        card_number: encryptedCard,
        card_exp: mm + yy,
      })
      setNegdi(res)
      setStep(2)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Алдаа гарлаа'
      setError(message)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  const confirm = async () => {
    if (!negdi) return

    if (!cvv || cvv.length !== 3) {
      setError('CVV код 3 оронтой байх ёстой')
      return
    }

    setLoading(true)
    setError(null)

    try {
      const list = await api.post<UserCard[]>('/card/negdi/confirm', {
        checkid: negdi.checkid,
        tranid: negdi.tranid,
        cardcvv: parseInt(cvv, 10),
      })
      toast.success('Карт амжилттай нэмэгдлээ')
      onClose(list?.[0])
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Алдаа гарлаа'
      setError(message)
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      {error && (
        <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">{error}</div>
      )}

      {step === 1 && (
        <>
          <div className="space-y-2">
            <label htmlFor="card-number" className="text-sm font-medium">
              Картын дугаар
            </label>
            <Input
              id="card-number"
              value={numUi}
              onChange={(e) => onCardInput(e.target.value)}
              placeholder="**** **** **** ****"
              maxLength={19}
              aria-invalid={!!error}
              autoComplete="cc-number"
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label htmlFor="card-month" className="sr-only">
                Сар
              </label>
              <select
                id="card-month"
                className="h-10 w-full rounded-md border px-3 text-sm"
                value={mm}
                onChange={(e) => setMM(e.target.value)}
                aria-label="Дуусах сар"
              >
                <option value="">MM</option>
                {MONTHS.map((m) => (
                  <option key={m} value={m}>
                    {m}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label htmlFor="card-year" className="sr-only">
                Жил
              </label>
              <select
                id="card-year"
                className="h-10 w-full rounded-md border px-3 text-sm"
                value={yy}
                onChange={(e) => setYY(e.target.value)}
                aria-label="Дуусах жил"
              >
                <option value="">YYYY</option>
                {YEARS.map((y) => (
                  <option key={y.value} value={y.value}>
                    {y.label}
                  </option>
                ))}
              </select>
            </div>
          </div>
          <Button className="w-full" onClick={check} disabled={loading}>
            {loading && <Loader />} Үргэлжлүүлэх
          </Button>
        </>
      )}

      {step === 2 && (
        <>
          <div className="rounded-md border p-3 text-sm">
            {negdi?.maskedpan?.replace(/(.{4})/g, '$1 ').trim()}
          </div>
          <div className="space-y-2">
            <label htmlFor="card-cvv" className="text-sm font-medium">
              CVV
            </label>
            <Input
              id="card-cvv"
              type="password"
              inputMode="numeric"
              value={cvv}
              onChange={(e) => setCVV(e.target.value.replace(/\D/g, ''))}
              placeholder="***"
              maxLength={3}
              aria-invalid={!!error}
              autoComplete="cc-csc"
            />
          </div>
          <Button className="w-full" onClick={confirm} disabled={loading || cvv.length !== 3}>
            {loading && <Loader />} Үргэлжлүүлэх
          </Button>
        </>
      )}
    </div>
  )
}

function Loader() {
  return (
    <span className="mr-2 inline-flex" role="status" aria-label="Уншиж байна">
      <svg className="h-4 w-4 animate-spin" viewBox="0 0 24 24" aria-hidden="true">
        <circle
          className="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          strokeWidth="4"
          fill="none"
        />
        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
      </svg>
    </span>
  )
}
