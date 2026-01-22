'use client'

import { useTranslations } from 'next-intl'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Wallet, Star, FileText, QrCode } from 'lucide-react'
import { cn } from '@/lib/utils'
import { formatMoney, getStatusBadgeClass } from '../utils'

type AccountCardProps = {
  account: App.TpayAccount
  showBalance: boolean
  onSetDefault: (accountId: number) => void
  onViewStatement: (account: App.TpayAccount) => void
  onGenerateQR: (accountId: number) => void
}

export function AccountCard({
  account,
  showBalance,
  onSetDefault,
  onViewStatement,
  onGenerateQR,
}: AccountCardProps) {
  const t = useTranslations()

  return (
    <Card
      className={cn(
        'group relative overflow-hidden border-2 transition-all hover:shadow-lg',
        account.is_default
          ? 'border-primary/50 bg-primary/5 shadow-md'
          : 'border-border hover:border-primary/30'
      )}
    >
      <CardContent className="p-6">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="mb-3 flex items-center gap-2">
              <div
                className={cn(
                  'flex h-10 w-10 items-center justify-center rounded-xl',
                  account.is_default
                    ? 'bg-primary/10 text-primary'
                    : 'bg-muted text-muted-foreground'
                )}
              >
                <Wallet className="h-5 w-5" />
              </div>
              {account.is_default && (
                <Badge variant="secondary" className="gap-1 bg-primary/10 text-primary">
                  <Star className="h-3 w-3 fill-primary" />
                  {t('default')}
                </Badge>
              )}
            </div>
            <h3 className="mb-1 font-semibold">{account.account_name}</h3>
            <p className="mb-4 font-mono text-xs text-muted-foreground">{account.account_no}</p>
            <div className="mb-4">
              <p className="text-xs text-muted-foreground">{t('balance')}</p>
              <p className="text-xl font-bold">
                {showBalance ? formatMoney(account.balance, account.currency) : '••••••'}
              </p>
            </div>
            <Badge className={cn('text-xs', getStatusBadgeClass(account.status))}>
              {t(account.status)}
            </Badge>
          </div>
        </div>
        <div className="mt-4 flex gap-2 border-t pt-4">
          {!account.is_default && (
            <Button
              variant="ghost"
              size="sm"
              className="flex-1"
              onClick={() => onSetDefault(account.id)}
            >
              <Star className="mr-1.5 h-3.5 w-3.5" />
              {t('set_default')}
            </Button>
          )}
          <Button
            variant="ghost"
            size="sm"
            className="flex-1"
            onClick={() => onViewStatement(account)}
          >
            <FileText className="mr-1.5 h-3.5 w-3.5" />
            {t('statement')}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="flex-1"
            onClick={() => onGenerateQR(account.id)}
          >
            <QrCode className="mr-1.5 h-3.5 w-3.5" />
            {t('qr')}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
