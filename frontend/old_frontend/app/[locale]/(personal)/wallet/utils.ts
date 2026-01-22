/**
 * Wallet utility functions and constants
 */

/**
 * Format money amount with currency
 */
export function formatMoney(amount: number, currency = 'MNT'): string {
  return new Intl.NumberFormat('mn-MN', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
  }).format(amount)
}

/**
 * Format card number with spaces
 */
export function formatCardNumber(cardNo: string): string {
  return cardNo.replace(/(\d{4})/g, '$1 ').trim()
}

/**
 * Mask card number showing only last 4 digits
 */
export function maskCardNumber(cardNo: string): string {
  const last4 = cardNo.slice(-4)
  return `•••• •••• •••• ${last4}`
}

/**
 * Get status badge color classes
 */
export function getStatusBadgeClass(status: string): string {
  const colors: Record<string, string> = {
    active: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300',
    completed: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300',
    pending: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
    failed: 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300',
    blocked: 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300',
    expired: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
    inactive: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
  }
  return colors[status] || colors.inactive
}

/**
 * Get card gradient based on card type
 */
export function getCardGradient(cardType: string): string {
  switch (cardType?.toLowerCase()) {
    case 'visa':
      return 'from-blue-600 via-blue-500 to-blue-700'
    case 'mastercard':
      return 'from-orange-500 via-red-500 to-red-600'
    default:
      return 'from-green-500 via-teal-500 to-teal-600'
  }
}
