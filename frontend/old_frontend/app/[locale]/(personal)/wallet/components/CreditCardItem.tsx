'use client'

import { useTranslations } from 'next-intl'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Check, ShieldCheck } from 'lucide-react'
import { cn } from '@/lib/utils'
import { maskCardNumber, getCardGradient } from '../utils'

type CreditCardItemProps = {
  card: App.TpayCard
  onVerify: (card: App.TpayCard) => void
}

export function CreditCardItem({ card, onVerify }: CreditCardItemProps) {
  const t = useTranslations()

  return (
    <div
      className={cn(
        'group relative overflow-hidden rounded-2xl p-6 text-white shadow-xl transition-all hover:scale-[1.02] hover:shadow-2xl',
        'bg-gradient-to-br',
        getCardGradient(card.card_type)
      )}
    >
      {/* Decorative elements */}
      <div className="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-white/10 blur-2xl" />
      <div className="absolute -bottom-8 -left-8 h-24 w-24 rounded-full bg-white/5 blur-xl" />

      <div className="relative z-10">
        <div className="mb-6 flex items-start justify-between">
          <div>
            <p className="text-xs font-medium uppercase tracking-wider opacity-80">
              {card.bank_name || t('card')}
            </p>
            <p className="mt-4 font-mono text-2xl tracking-widest">
              {maskCardNumber(card.card_no)}
            </p>
          </div>
          <div className="text-right">
            <p className="text-xs font-semibold uppercase tracking-wider opacity-80">
              {card.card_type}
            </p>
            {card.is_verified ? (
              <Badge className="mt-2 bg-white/20 text-white backdrop-blur-sm">
                <Check className="mr-1 h-3 w-3" />
                {t('verified')}
              </Badge>
            ) : (
              <Badge className="mt-2 bg-amber-400/30 text-white backdrop-blur-sm">
                {t('not_verified')}
              </Badge>
            )}
          </div>
        </div>

        <div className="mt-8 flex items-end justify-between border-t border-white/20 pt-4">
          <div>
            <p className="text-xs uppercase tracking-wider opacity-60">{t('card_holder')}</p>
            <p className="mt-1 font-semibold uppercase tracking-wide">{card.card_holder_name}</p>
          </div>
          <div className="text-right">
            <p className="text-xs uppercase tracking-wider opacity-60">{t('expiration_date')}</p>
            <p className="mt-1 font-semibold">{card.expiry_date}</p>
          </div>
        </div>

        {!card.is_verified && (
          <div className="mt-6">
            <Button
              size="sm"
              variant="secondary"
              className="w-full bg-white/20 text-white backdrop-blur-sm hover:bg-white/30"
              onClick={() => onVerify(card)}
            >
              <ShieldCheck className="mr-2 h-4 w-4" />
              {t('verify_card')}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
