'use client'
import * as React from 'react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import api from '@/lib/api'

export default function VerifyCard({
  card,
  onClose,
}: {
  card: { id: number; card_number: string }
  onClose: () => void
}) {
  const [step, setStep] = React.useState<1 | 2>(1)
  const [otp, setOtp] = React.useState('')
  const [loading, setLoading] = React.useState(false)

  const sendOtp = async () => {
    setLoading(true)
    try {
      await api.post('/card/negdi/send_otp', { id: card.id })
      setStep(2)
    } finally {
      setLoading(false)
    }
  }
  const verify = async () => {
    setLoading(true)
    try {
      const formatted = `${otp.slice(0, 2)}.${otp.slice(2, 4)}`
      await api.post('/card/negdi/verify', { id: card.id, otp: formatted })
      onClose()
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      {step === 1 && (
        <>
          <p className="text-muted-foreground text-sm">
            Таны {card.card_number} картад ирэх нэг удаагийн кодоор баталгаажуулна.
          </p>
          <Button onClick={sendOtp} disabled={loading} className="w-full">
            Код илгээх
          </Button>
        </>
      )}
      {step === 2 && (
        <>
          <div className="text-sm font-medium">Нэг удаагийн код</div>
          <Input
            value={otp}
            onChange={(e) => setOtp(e.target.value)}
            maxLength={4}
            placeholder="0000"
          />
          <Button onClick={verify} disabled={loading || otp.length !== 4} className="w-full">
            Баталгаажуулах
          </Button>
        </>
      )}
    </div>
  )
}
