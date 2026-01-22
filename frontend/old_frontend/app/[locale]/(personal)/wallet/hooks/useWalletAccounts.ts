'use client'

import { useState, useCallback, useEffect } from 'react'
import api from '@/lib/api'
import { toast } from 'sonner'

type UseWalletAccountsReturn = {
  accounts: App.TpayAccount[]
  loading: boolean
  defaultAccount: App.TpayAccount | undefined
  totalBalance: number
  loadAccounts: () => Promise<void>
  setDefault: (accountId: number) => Promise<void>
  generateQR: (accountId: number) => Promise<App.TpayQRCode | null>
  loadStatement: (accountId: number) => Promise<App.TpayStatement[]>
}

export function useWalletAccounts(): UseWalletAccountsReturn {
  const [accounts, setAccounts] = useState<App.TpayAccount[]>([])
  const [loading, setLoading] = useState(false)

  const loadAccounts = useCallback(async () => {
    setLoading(true)
    try {
      const data = await api.get<App.TpayAccount[]>('/me/accounts', { hasToast: false })
      setAccounts(data || [])
    } catch {
      setAccounts([])
    } finally {
      setLoading(false)
    }
  }, [])

  const setDefault = useCallback(async (accountId: number) => {
    try {
      await api.put('/me/accounts/default', { account_id: accountId })
      await loadAccounts()
      toast.success('Үндсэн данс тохируулагдлаа')
    } catch (error) {
      console.error('Failed to set default account:', error)
    }
  }, [loadAccounts])

  const generateQR = useCallback(async (accountId: number): Promise<App.TpayQRCode | null> => {
    try {
      const data = await api.post<App.TpayQRCode>(`/me/accounts/${accountId}/qr`, {})
      return data
    } catch (error) {
      console.error('Failed to generate QR code:', error)
      return null
    }
  }, [])

  const loadStatement = useCallback(async (accountId: number): Promise<App.TpayStatement[]> => {
    try {
      const data = await api.get<App.TpayStatement[]>('/me/accounts/statement', {
        query: { account_id: accountId },
        hasToast: false,
      })
      return data || []
    } catch {
      return []
    }
  }, [])

  useEffect(() => {
    loadAccounts()
  }, [loadAccounts])

  const defaultAccount = accounts.find((a) => a.is_default) || accounts[0]
  const totalBalance = accounts.reduce((sum, a) => sum + a.balance, 0)

  return {
    accounts,
    loading,
    defaultAccount,
    totalBalance,
    loadAccounts,
    setDefault,
    generateQR,
    loadStatement,
  }
}
