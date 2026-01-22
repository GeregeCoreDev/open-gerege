import type { JSX, ReactNode } from 'react'
import { NextIntlClientProvider } from 'next-intl'
import { getMessages, setRequestLocale } from 'next-intl/server'
import UserBootstrap from '@/lib/bootstrap/UserBootstrap'
import AuthBootstrap from '@/lib/bootstrap/AuthBootstrap'
import { ErrorBoundary } from '@/components/error-boundary'

type LocaleLayoutProps = {
  children: ReactNode
  params: Promise<{ locale: string }> // <- async params
}

export default async function LocaleLayout({
  children,
  params,
}: LocaleLayoutProps): Promise<JSX.Element> {

  // Next.js 15 requires this:
  const { locale } = await params

  setRequestLocale(locale)
  const messages = await getMessages({ locale })

  return (
    <NextIntlClientProvider key={locale} locale={locale} messages={messages}>
      <ErrorBoundary>
        <UserBootstrap />
        <AuthBootstrap />
        {children}
      </ErrorBoundary>
    </NextIntlClientProvider>
  )
}

export const prerender = false
export const dynamic = 'force-dynamic'