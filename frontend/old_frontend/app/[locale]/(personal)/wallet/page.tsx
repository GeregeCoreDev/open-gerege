/**
 * Wallet Page - Refactored
 *
 * Хэрэглэгчийн хэтэвч хуудас
 * - Дансны удирдлага
 * - Картын удирдлага
 * - Гүйлгээний түүх
 */

'use client'

import * as React from 'react'
import { useState } from 'react'
import { useTranslations } from 'next-intl'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/components/ui/form'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Wallet,
  CreditCard,
  ArrowLeftRight,
  QrCode,
  Plus,
  Loader2,
  FileText,
  Send,
  ArrowDownLeft,
  ArrowUpRight,
} from 'lucide-react'
import { cn } from '@/lib/utils'

// Local imports
import { useWalletAccounts, useWalletCards, useWalletTransactions } from './hooks'
import { BalanceCard, AccountCard, CreditCardItem } from './components'
import { formatMoney } from './utils'

// ============================================================
// Schemas
// ============================================================

const AddCardSchema = z.object({
  card_no: z.string().min(16, 'Card number must be 16 digits').max(19),
  expiry_date: z.string().min(5, 'Required (MM/YY)'),
  cvv: z.string().min(3, 'CVV must be 3 digits').max(4),
  card_holder_name: z.string().min(2, 'Required'),
})
type AddCardForm = z.infer<typeof AddCardSchema>

const ConfirmCardSchema = z.object({
  otp_amount: z.string().min(1, 'Required'),
})
type ConfirmCardForm = z.infer<typeof ConfirmCardSchema>

const QRPaySchema = z.object({
  qr_string: z.string().min(1, 'Required'),
  amount: z.number().optional(),
  pin: z.string().min(4, 'PIN required').max(6),
})
type QRPayForm = z.infer<typeof QRPaySchema>

const P2PSchema = z.object({
  to_account: z.string().min(1, 'Required'),
  amount: z.number().min(1, 'Amount required'),
  description: z.string().optional(),
  pin: z.string().min(4, 'PIN required').max(6),
})
type P2PForm = z.infer<typeof P2PSchema>

// ============================================================
// Main Component
// ============================================================

export default function WalletPage() {
  const t = useTranslations()

  // Custom hooks
  const {
    accounts,
    loading: accountsLoading,
    defaultAccount,
    totalBalance,
    setDefault,
    generateQR,
    loadStatement,
  } = useWalletAccounts()

  const { cards, loading: cardsLoading, addCard, sendOtp, confirmCard } = useWalletCards()

  const { qrPay, p2pTransfer } = useWalletTransactions()

  // UI State
  const [showBalance, setShowBalance] = useState(true)
  const [openAddCard, setOpenAddCard] = useState(false)
  const [openConfirmCard, setOpenConfirmCard] = useState(false)
  const [openQRPay, setOpenQRPay] = useState(false)
  const [openP2P, setOpenP2P] = useState(false)
  const [openStatement, setOpenStatement] = useState(false)
  const [openQRCode, setOpenQRCode] = useState(false)

  // Selected items
  const [selectedCard, setSelectedCard] = useState<App.TpayCard | null>(null)
  const [selectedAccount, setSelectedAccount] = useState<App.TpayAccount | null>(null)
  const [qrCodeData, setQrCodeData] = useState<App.TpayQRCode | null>(null)
  const [statement, setStatement] = useState<App.TpayStatement[]>([])

  // Forms
  const addCardForm = useForm<AddCardForm>({
    resolver: zodResolver(AddCardSchema),
    defaultValues: { card_no: '', expiry_date: '', cvv: '', card_holder_name: '' },
  })

  const confirmCardForm = useForm<ConfirmCardForm>({
    resolver: zodResolver(ConfirmCardSchema),
    defaultValues: { otp_amount: '' },
  })

  const qrPayForm = useForm<QRPayForm>({
    resolver: zodResolver(QRPaySchema),
    defaultValues: { qr_string: '', pin: '' },
  })

  const p2pForm = useForm<P2PForm>({
    resolver: zodResolver(P2PSchema),
    defaultValues: { to_account: '', amount: 0, description: '', pin: '' },
  })

  // Handlers
  const handleViewStatement = async (account: App.TpayAccount) => {
    setSelectedAccount(account)
    const data = await loadStatement(account.id)
    setStatement(data)
    setOpenStatement(true)
  }

  const handleGenerateQR = async (accountId: number) => {
    const data = await generateQR(accountId)
    if (data) {
      setQrCodeData(data)
      setOpenQRCode(true)
    }
  }

  const handleAddCard = async (values: AddCardForm) => {
    const card = await addCard(values)
    if (card) {
      setSelectedCard(card)
      setOpenAddCard(false)
      addCardForm.reset()
      setOpenConfirmCard(true)
    }
  }

  const handleVerifyCard = async (card: App.TpayCard) => {
    setSelectedCard(card)
    await sendOtp(card.id)
    setOpenConfirmCard(true)
  }

  const handleConfirmCard = async (values: ConfirmCardForm) => {
    if (!selectedCard) return
    const success = await confirmCard(selectedCard.id, values.otp_amount)
    if (success) {
      setOpenConfirmCard(false)
      confirmCardForm.reset()
    }
  }

  const handleQRPay = async (values: QRPayForm) => {
    const success = await qrPay(values)
    if (success) {
      setOpenQRPay(false)
      qrPayForm.reset()
    }
  }

  const handleP2P = async (values: P2PForm) => {
    const success = await p2pTransfer(values)
    if (success) {
      setOpenP2P(false)
      p2pForm.reset()
    }
  }

  const loading = accountsLoading || cardsLoading

  return (
    <div className="flex h-full w-full flex-col overflow-auto p-4 sm:p-6 lg:p-8">
      <div className="mx-auto w-full max-w-7xl space-y-8">
        {/* Header */}
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight lg:text-4xl">{t('wallet_title')}</h1>
            <p className="mt-2 text-sm text-muted-foreground lg:text-base">{t('wallet_subtitle')}</p>
          </div>
          <div className="flex gap-3">
            <Button
              variant="outline"
              size="lg"
              onClick={() => setOpenQRPay(true)}
              className="gap-2 border-2 transition-all hover:border-primary/50 hover:bg-primary/5"
            >
              <QrCode className="h-4 w-4" />
              {t('qr_pay')}
            </Button>
            <Button
              size="lg"
              onClick={() => setOpenP2P(true)}
              className="gap-2 shadow-lg transition-all hover:shadow-xl"
            >
              <Send className="h-4 w-4" />
              {t('p2p_transfer')}
            </Button>
          </div>
        </div>

        {/* Balance Card */}
        <BalanceCard
          totalBalance={totalBalance}
          defaultAccountName={defaultAccount?.account_name}
          showBalance={showBalance}
          onToggleBalance={() => setShowBalance(!showBalance)}
        />

        {loading ? (
          <div className="flex h-40 items-center justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-primary" />
          </div>
        ) : (
          <Tabs defaultValue="accounts" className="w-full">
            <TabsList className="inline-flex h-12 w-full items-center justify-start rounded-xl bg-muted p-1.5">
              <TabsTrigger
                value="accounts"
                className="group relative flex flex-1 items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all data-[state=active]:bg-background data-[state=active]:shadow-sm"
              >
                <Wallet className="h-4 w-4 transition-transform group-data-[state=active]:scale-110" />
                <span>{t('accounts')}</span>
              </TabsTrigger>
              <TabsTrigger
                value="cards"
                className="group relative flex flex-1 items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all data-[state=active]:bg-background data-[state=active]:shadow-sm"
              >
                <CreditCard className="h-4 w-4 transition-transform group-data-[state=active]:scale-110" />
                <span>{t('cards')}</span>
              </TabsTrigger>
              <TabsTrigger
                value="transactions"
                className="group relative flex flex-1 items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all data-[state=active]:bg-background data-[state=active]:shadow-sm"
              >
                <ArrowLeftRight className="h-4 w-4 transition-transform group-data-[state=active]:scale-110" />
                <span>{t('transactions')}</span>
              </TabsTrigger>
            </TabsList>

            {/* Accounts Tab */}
            <TabsContent value="accounts" className="mt-6">
              <Card className="border-0 shadow-lg">
                <CardHeader className="pb-4">
                  <CardTitle className="text-xl">{t('my_accounts')}</CardTitle>
                  <CardDescription>{t('manage_your_accounts')}</CardDescription>
                </CardHeader>
                <CardContent>
                  {accounts.length === 0 ? (
                    <EmptyState icon={Wallet} title={t('no_accounts')} description={t('add_account_description')} />
                  ) : (
                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                      {accounts.map((account) => (
                        <AccountCard
                          key={account.id}
                          account={account}
                          showBalance={showBalance}
                          onSetDefault={setDefault}
                          onViewStatement={handleViewStatement}
                          onGenerateQR={handleGenerateQR}
                        />
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </TabsContent>

            {/* Cards Tab */}
            <TabsContent value="cards" className="mt-6">
              <Card className="border-0 shadow-lg">
                <CardHeader className="flex flex-row items-center justify-between pb-4">
                  <div>
                    <CardTitle className="text-xl">{t('my_cards')}</CardTitle>
                    <CardDescription>{t('manage_your_cards')}</CardDescription>
                  </div>
                  <Button size="lg" onClick={() => setOpenAddCard(true)} className="gap-2">
                    <Plus className="h-4 w-4" />
                    {t('add_new_card')}
                  </Button>
                </CardHeader>
                <CardContent>
                  {cards.length === 0 ? (
                    <EmptyState
                      icon={CreditCard}
                      title={t('no_cards')}
                      description={t('add_card_description')}
                      action={
                        <Button variant="outline" className="mt-6 gap-2" onClick={() => setOpenAddCard(true)}>
                          <Plus className="h-4 w-4" />
                          {t('add_new_card')}
                        </Button>
                      }
                    />
                  ) : (
                    <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                      {cards.map((card) => (
                        <CreditCardItem key={card.id} card={card} onVerify={handleVerifyCard} />
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </TabsContent>

            {/* Transactions Tab */}
            <TabsContent value="transactions" className="mt-6">
              <Card className="border-0 shadow-lg">
                <CardHeader className="pb-4">
                  <CardTitle className="text-xl">{t('transaction_history')}</CardTitle>
                  <CardDescription>{t('view_all_transactions')}</CardDescription>
                </CardHeader>
                <CardContent>
                  <EmptyState
                    icon={ArrowLeftRight}
                    title={t('no_transactions')}
                    description={t('transaction_empty_description')}
                  />
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        )}
      </div>

      {/* Dialogs */}
      <AddCardDialog
        open={openAddCard}
        onOpenChange={setOpenAddCard}
        form={addCardForm}
        onSubmit={handleAddCard}
      />

      <ConfirmCardDialog
        open={openConfirmCard}
        onOpenChange={setOpenConfirmCard}
        form={confirmCardForm}
        onSubmit={handleConfirmCard}
      />

      <QRPayDialog
        open={openQRPay}
        onOpenChange={setOpenQRPay}
        form={qrPayForm}
        onSubmit={handleQRPay}
      />

      <P2PDialog open={openP2P} onOpenChange={setOpenP2P} form={p2pForm} onSubmit={handleP2P} />

      <StatementDialog
        open={openStatement}
        onOpenChange={setOpenStatement}
        account={selectedAccount}
        statement={statement}
      />

      <QRCodeDialog open={openQRCode} onOpenChange={setOpenQRCode} qrCodeData={qrCodeData} />
    </div>
  )
}

// ============================================================
// Sub-components
// ============================================================

function EmptyState({
  icon: Icon,
  title,
  description,
  action,
}: {
  icon: React.ComponentType<{ className?: string }>
  title: string
  description: string
  action?: React.ReactNode
}) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="mb-4 rounded-full bg-muted p-4">
        <Icon className="h-8 w-8 text-muted-foreground" />
      </div>
      <p className="text-lg font-medium text-muted-foreground">{title}</p>
      <p className="mt-1 text-sm text-muted-foreground/80">{description}</p>
      {action}
    </div>
  )
}

function AddCardDialog({
  open,
  onOpenChange,
  form,
  onSubmit,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  form: ReturnType<typeof useForm<AddCardForm>>
  onSubmit: (values: AddCardForm) => void
}) {
  const t = useTranslations()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('add_new_card')}</DialogTitle>
          <DialogDescription>{t('fill_the_fields_and_save_to_create')}</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="card_no"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('card_number')}</FormLabel>
                  <FormControl>
                    <Input placeholder="0000 0000 0000 0000" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="expiry_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('expiration_date')}</FormLabel>
                    <FormControl>
                      <Input placeholder="MM/YY" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="cvv"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('cvv')}</FormLabel>
                    <FormControl>
                      <Input placeholder="123" type="password" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <FormField
              control={form.control}
              name="card_holder_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('card_holder_name')}</FormLabel>
                  <FormControl>
                    <Input placeholder="JOHN DOE" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                {t('cancel')}
              </Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {t('save')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

function ConfirmCardDialog({
  open,
  onOpenChange,
  form,
  onSubmit,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  form: ReturnType<typeof useForm<ConfirmCardForm>>
  onSubmit: (values: ConfirmCardForm) => void
}) {
  const t = useTranslations()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('confirm_card')}</DialogTitle>
          <DialogDescription>{t('enter_otp_amount')}</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="otp_amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('amount')}</FormLabel>
                  <FormControl>
                    <Input placeholder="25.15" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                {t('cancel')}
              </Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {t('verify')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

function QRPayDialog({
  open,
  onOpenChange,
  form,
  onSubmit,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  form: ReturnType<typeof useForm<QRPayForm>>
  onSubmit: (values: QRPayForm) => void
}) {
  const t = useTranslations()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('qr_pay')}</DialogTitle>
          <DialogDescription>{t('enter_qr_string')}</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="qr_string"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('enter_qr_string')}</FormLabel>
                  <FormControl>
                    <Input placeholder="QR code string..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="pin"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('enter_pin')}</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="••••" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                {t('cancel')}
              </Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {t('pay_now')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

function P2PDialog({
  open,
  onOpenChange,
  form,
  onSubmit,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  form: ReturnType<typeof useForm<P2PForm>>
  onSubmit: (values: P2PForm) => void
}) {
  const t = useTranslations()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('p2p_transfer')}</DialogTitle>
          <DialogDescription>{t('transfer_description')}</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="to_account"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('transfer_to')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('to_account')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('transfer_amount')}</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      placeholder="0"
                      {...field}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('description')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('transfer_description')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="pin"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('enter_pin')}</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="••••" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                {t('cancel')}
              </Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {t('transfer_now')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}

function StatementDialog({
  open,
  onOpenChange,
  account,
  statement,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  account: App.TpayAccount | null
  statement: App.TpayStatement[]
}) {
  const t = useTranslations()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>{t('account_statement')}</DialogTitle>
          <DialogDescription>
            {account?.account_name} - {account?.account_no}
          </DialogDescription>
        </DialogHeader>
        <div className="max-h-[60vh] overflow-auto">
          {statement.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <FileText className="h-12 w-12 text-gray-400" />
              <p className="mt-2 text-gray-500">{t('no_transactions')}</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t('type')}</TableHead>
                  <TableHead>{t('amount')}</TableHead>
                  <TableHead>{t('balance')}</TableHead>
                  <TableHead>{t('description')}</TableHead>
                  <TableHead>{t('created_date')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {statement.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell className="flex items-center gap-2">
                      {item.type === 'credit' ? (
                        <ArrowDownLeft className="h-4 w-4 text-emerald-500" />
                      ) : (
                        <ArrowUpRight className="h-4 w-4 text-rose-500" />
                      )}
                      {t(item.type)}
                    </TableCell>
                    <TableCell
                      className={cn(
                        'font-medium',
                        item.type === 'credit' ? 'text-emerald-600' : 'text-rose-600'
                      )}
                    >
                      {item.type === 'credit' ? '+' : '-'}
                      {formatMoney(item.amount)}
                    </TableCell>
                    <TableCell>{formatMoney(item.balance_after)}</TableCell>
                    <TableCell>{item.description || '—'}</TableCell>
                    <TableCell>{item.created_date?.slice(0, 10)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t('close')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function QRCodeDialog({
  open,
  onOpenChange,
  qrCodeData,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  qrCodeData: App.TpayQRCode | null
}) {
  const t = useTranslations()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('generate_qr')}</DialogTitle>
        </DialogHeader>
        <div className="flex flex-col items-center justify-center py-6">
          {qrCodeData?.qr_image_url ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img src={qrCodeData.qr_image_url} alt="QR Code" className="h-64 w-64" />
          ) : (
            <div className="flex h-64 w-64 items-center justify-center rounded-lg border-2 border-dashed">
              <QrCode className="h-24 w-24 text-gray-400" />
            </div>
          )}
          {qrCodeData?.qr_string && (
            <p className="mt-4 max-w-full truncate rounded bg-gray-100 px-4 py-2 text-xs dark:bg-gray-800">
              {qrCodeData.qr_string}
            </p>
          )}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t('close')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
