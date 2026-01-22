'use client'

import { useCallback } from 'react'
import api from '@/lib/api'
import { toast } from 'sonner'

type QRPayData = {
  qr_string: string
  amount?: number
  pin: string
}

type P2PData = {
  to_account: string
  amount: number
  description?: string
  pin: string
}

type UseWalletTransactionsReturn = {
  qrPay: (data: QRPayData) => Promise<boolean>
  p2pTransfer: (data: P2PData) => Promise<boolean>
}

export function useWalletTransactions(): UseWalletTransactionsReturn {
  const qrPay = useCallback(async (data: QRPayData): Promise<boolean> => {
    try {
      await api.post<App.TpayTransactionRes>('/me/tpay/transaction/qr-pay', data)
      toast.success('QR төлбөр амжилттай')
      return true
    } catch (error) {
      console.error('Failed to process QR payment:', error)
      return false
    }
  }, [])

  const p2pTransfer = useCallback(async (data: P2PData): Promise<boolean> => {
    try {
      await api.post<App.TpayTransactionRes>('/me/tpay/transaction/p2p', data)
      toast.success('Шилжүүлэг амжилттай')
      return true
    } catch (error) {
      console.error('Failed to process P2P transfer:', error)
      return false
    }
  }, [])

  return {
    qrPay,
    p2pTransfer,
  }
}
