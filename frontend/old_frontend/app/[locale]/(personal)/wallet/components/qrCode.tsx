'use client'
/**
 * Энгийн placeholder: серверээс ирсэн QR payload-аа харуулна.
 * Хэрэв бодит QR зургаар харуулах бол `qrcode.react` эсвэл `next-qrcode` суулгаад энэ компонент дотроо ашиглаарай.
 */
export default function QrCode({ value }: { value: string }) {
  return (
    <div className="min-w-[220px] rounded-md border-4 border-gray-200 bg-white p-3 text-center">
      <div className="text-muted-foreground mb-2 text-xs">QR payload</div>
      <div className="max-w-[240px] font-mono text-[10px] break-all">{value || '—'}</div>
    </div>
  )
}
