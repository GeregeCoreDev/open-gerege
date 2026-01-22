'use client'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Wallet, Star, Eye, EyeOff } from 'lucide-react'
import { formatMoney } from '../utils'

type BalanceCardProps = {
  totalBalance: number
  defaultAccountName?: string
  showBalance: boolean
  onToggleBalance: () => void
}

export function BalanceCard({
  totalBalance,
  defaultAccountName,
  showBalance,
  onToggleBalance,
}: BalanceCardProps) {
  return (
    <Card className="group relative overflow-hidden border-0 bg-gradient-to-br from-primary via-primary/90 to-primary/80 shadow-2xl shadow-primary/25 transition-all hover:shadow-primary/30">
      {/* Animated background pattern */}
      <div className="absolute inset-0 bg-[linear-gradient(45deg,transparent_25%,rgba(255,255,255,.1)_50%,transparent_75%,transparent_100%)] bg-[length:20px_20px] opacity-20" />
      <div className="absolute inset-0 bg-grid-white/10 [mask-image:linear-gradient(0deg,white,transparent)]" />

      {/* Decorative circles */}
      <div className="absolute -right-12 -top-12 h-40 w-40 rounded-full bg-white/10 blur-3xl transition-all group-hover:scale-150" />
      <div className="absolute -bottom-8 -left-8 h-32 w-32 rounded-full bg-white/5 blur-2xl transition-all group-hover:scale-125" />

      <CardContent className="relative p-8 text-white">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <p className="text-sm font-medium uppercase tracking-wider opacity-90">
              Нийт үлдэгдэл
            </p>
            <div className="mt-4 flex items-center gap-3">
              <p className="text-4xl font-bold tracking-tight lg:text-5xl">
                {showBalance ? formatMoney(totalBalance) : '••••••••'}
              </p>
              <Button
                variant="ghost"
                size="icon"
                className="h-9 w-9 text-white transition-all hover:scale-110 hover:bg-white/20"
                onClick={onToggleBalance}
                aria-label={showBalance ? 'Үлдэгдэл нуух' : 'Үлдэгдэл харуулах'}
              >
                {showBalance ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
              </Button>
            </div>
            {defaultAccountName && (
              <div className="mt-5 flex items-center gap-2">
                <Badge
                  variant="secondary"
                  className="border-white/30 bg-white/20 text-white backdrop-blur-sm hover:bg-white/30"
                >
                  <Star className="mr-1.5 h-3.5 w-3.5 fill-white" />
                  {defaultAccountName}
                </Badge>
              </div>
            )}
          </div>
          <div className="hidden sm:block">
            <div className="relative">
              <Wallet className="h-24 w-24 opacity-20 transition-transform group-hover:scale-110 group-hover:opacity-30" />
              <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent blur-3xl" />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
