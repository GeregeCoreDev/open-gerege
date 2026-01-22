'use server'

import CryptoJS from 'crypto-js'

const ENCRYPTION_KEY = process.env.CARD_ENCRYPTION_KEY || ''
const ENCRYPTION_IV = process.env.CARD_ENCRYPTION_IV || ''

if (!ENCRYPTION_KEY || !ENCRYPTION_IV) {
  console.warn(
    'Warning: CARD_ENCRYPTION_KEY or CARD_ENCRYPTION_IV not set in environment variables'
  )
}

export async function encryptCardNumber(cardNumber: string): Promise<string> {
  if (!ENCRYPTION_KEY || !ENCRYPTION_IV) {
    throw new Error('Encryption configuration missing')
  }

  // Validate card number format (16 digits)
  const cleanNumber = cardNumber.replace(/\D/g, '')
  if (cleanNumber.length !== 16) {
    throw new Error('Invalid card number format')
  }

  // Luhn algorithm validation
  if (!isValidLuhn(cleanNumber)) {
    throw new Error('Invalid card number')
  }

  const key = CryptoJS.enc.Utf8.parse(ENCRYPTION_KEY)
  const iv = CryptoJS.enc.Utf8.parse(ENCRYPTION_IV)

  return CryptoJS.AES.encrypt(cleanNumber, key, {
    iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  }).toString()
}

/**
 * Luhn algorithm (mod 10) to validate card numbers
 */
function isValidLuhn(cardNumber: string): boolean {
  let sum = 0
  let isEven = false

  for (let i = cardNumber.length - 1; i >= 0; i--) {
    let digit = parseInt(cardNumber[i], 10)

    if (isEven) {
      digit *= 2
      if (digit > 9) {
        digit -= 9
      }
    }

    sum += digit
    isEven = !isEven
  }

  return sum % 10 === 0
}
