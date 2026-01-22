'use client'

import { Dialog, DialogContent, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { LogIn, Loader2 } from 'lucide-react'
import { VisuallyHidden } from '@radix-ui/react-visually-hidden'
import { useTranslations } from 'next-intl'
import { useSSOLogin, getSSORedirectUrl } from '@/domains/auth'

/**
 * SSO Popup Sign In Component
 */
export function SignInPopUp() {
  const t = useTranslations()
  const { isOpen, isLoading, openPopup, closePopup, ssoEmbedUrl } = useSSOLogin()

  return (
    <>
      <Button onClick={openPopup} className="h-12 gap-2" disabled={isLoading}>
        {isLoading ? (
          <Loader2 className="size-6 animate-spin" />
        ) : (
          <LogIn className="size-6" />
        )}
        <span>{t('sign_in')}</span>
      </Button>

      <Dialog open={isOpen} onOpenChange={(open) => !open && closePopup()}>
        <DialogContent className="h-full max-h-[70vh] overflow-hidden rounded-xl border-0 p-0 sm:max-w-4xl">
          <VisuallyHidden>
            <DialogTitle>{t('sign_in')}</DialogTitle>
            <DialogDescription>{t('sign_in_description')}</DialogDescription>
          </VisuallyHidden>

          <iframe
            src={ssoEmbedUrl}
            title="Gerege SSO Login"
            width="100%"
            height="100%"
            className="rounded-xl border-0"
            sandbox="allow-scripts allow-forms allow-same-origin allow-popups"
          />
        </DialogContent>
      </Dialog>
    </>
  )
}

/**
 * SSO Direct Link Component (redirect instead of popup)
 */
export function SignInLink() {
  const t = useTranslations()
  const { redirectToSSO, isLoading } = useSSOLogin()

  return (
    <Button onClick={redirectToSSO} className="h-12 gap-2" disabled={isLoading}>
      {isLoading ? (
        <Loader2 className="size-6 animate-spin" />
      ) : (
        <LogIn className="size-6" />
      )}
      <span>{t('sign_in')}</span>
    </Button>
  )
}

/**
 * Simple Sign In Button (direct redirect)
 */
export function SignInButton({ className }: { className?: string }) {
  const t = useTranslations()

  const handleClick = () => {
    window.location.href = getSSORedirectUrl()
  }

  return (
    <Button onClick={handleClick} className={className}>
      <LogIn className="mr-2 size-4" />
      {t('sign_in')}
    </Button>
  )
}
