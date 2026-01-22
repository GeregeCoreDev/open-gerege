'use client'

import { useState, useCallback, useEffect } from 'react'
import api from '@/lib/api'
import { toast } from 'sonner'

type AddCardData = {
  card_no: string
  expiry_date: string
  cvv: string
  card_holder_name: string
}

type UseWalletCardsReturn = {
  cards: App.TpayCard[]
  loading: boolean
  loadCards: () => Promise<void>
  addCard: (data: AddCardData) => Promise<App.TpayCard | null>
  sendOtp: (cardId: number) => Promise<boolean>
  confirmCard: (cardId: number, otpAmount: string) => Promise<boolean>
  verifyCard: (cardId: number) => Promise<boolean>
}

export function useWalletCards(): UseWalletCardsReturn {
  const [cards, setCards] = useState<App.TpayCard[]>([])
  const [loading, setLoading] = useState(false)

  const loadCards = useCallback(async () => {
    setLoading(true)
    try {
      const data = await api.get<App.TpayCard[]>('/me/card/list', { hasToast: false })
      setCards(data || [])
    } catch {
      setCards([])
    } finally {
      setLoading(false)
    }
  }, [])

  const addCard = useCallback(async (data: AddCardData): Promise<App.TpayCard | null> => {
    try {
      const card = await api.post<App.TpayCard>('/me/card/create', data)
      toast.success('Карт амжилттай нэмэгдлээ')
      return card
    } catch (error) {
      console.error('Failed to add card:', error)
      return null
    }
  }, [])

  const sendOtp = useCallback(async (cardId: number): Promise<boolean> => {
    try {
      await api.get('/me/card/otp', { query: { card_id: cardId } })
      toast.success('OTP код илгээгдлээ')
      return true
    } catch (error) {
      console.error('Failed to send OTP:', error)
      return false
    }
  }, [])

  const confirmCard = useCallback(async (cardId: number, otpAmount: string): Promise<boolean> => {
    try {
      await api.post('/me/card/confirm', {
        card_id: cardId,
        otp_amount: otpAmount,
      })
      await loadCards()
      toast.success('Карт баталгаажлаа')
      return true
    } catch (error) {
      console.error('Failed to confirm card:', error)
      return false
    }
  }, [loadCards])

  const verifyCard = useCallback(async (cardId: number): Promise<boolean> => {
    try {
      await api.post('/me/card/verify', { card_id: cardId })
      await loadCards()
      toast.success('Карт баталгаажлаа')
      return true
    } catch (error) {
      console.error('Failed to verify card:', error)
      return false
    }
  }, [loadCards])

  useEffect(() => {
    loadCards()
  }, [loadCards])

  return {
    cards,
    loading,
    loadCards,
    addCard,
    sendOtp,
    confirmCard,
    verifyCard,
  }
}
