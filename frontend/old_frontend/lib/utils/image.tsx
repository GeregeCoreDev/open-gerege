import api from '../api'

export async function compressImageFile(
  file: File,
  opts: { maxWidth?: number; maxHeight?: number; quality?: number } = {},
): Promise<File> {
  const { maxWidth = 1200, maxHeight = 1200, quality = 0.92 } = opts

  if (!file.type.startsWith('image/')) return file

  const arrayBuffer = await file.arrayBuffer()
  const blobURL = URL.createObjectURL(new Blob([arrayBuffer]))

  const img = await new Promise<HTMLImageElement>((resolve, reject) => {
    const image = new Image()
    image.onload = () => resolve(image)
    image.onerror = reject
    image.src = blobURL
  })

  let width = img.naturalWidth
  let height = img.naturalHeight

  // --------- Maintain aspect ratio ---------
  const ratio = Math.min(maxWidth / width, maxHeight / height, 1)
  width = Math.round(width * ratio)
  height = Math.round(height * ratio)

  const canvas = document.createElement('canvas')
  canvas.width = width
  canvas.height = height
  const ctx = canvas.getContext('2d')!
  ctx.drawImage(img, 0, 0, width, height)

  // -------------------------------------------------
  // Convert canvas to Blob (keep same mime type)
  // -------------------------------------------------
  const mime = file.type.includes('png') ? 'image/png' : 'image/jpeg'

  const blob = await new Promise<Blob | null>((resolve) =>
    canvas.toBlob((b) => resolve(b), mime, quality),
  )

  URL.revokeObjectURL(blobURL)

  if (!blob) return file

  return new File([blob], file.name, {
    type: mime,
    lastModified: Date.now(),
  })
}

export async function uploadImage(file: File, des: string, name?: string) {
  const base = (process.env.NEXT_PUBLIC_API_BASE ?? '').replace(/\/+$/, '')

  let cutName = name
  if (cutName && base && cutName.startsWith(base)) {
    cutName = cutName.slice(base.length) // '/file/xxx.png' үлдэнэ
  }

  if (cutName) {
    const pathParts = cutName.split('/')
    const fileName = pathParts[pathParts.length - 1] || ''
    cutName = fileName.split('.')[0] || ''
  } else {
    cutName = ''
  }

  const compressed = await compressImageFile(file, {
    maxWidth: 1200,
    maxHeight: 1200,
    quality: 0.8,
  })

  const fd = new FormData()
  fd.append('file', compressed)
  fd.append('description', des)
  fd.append('name', cutName)

  const url = await api.post<string>('/file/upload', fd)
  return url
}

export async function deleteImage(name: string) {
  const preview = process.env.NEXT_PUBLIC_API_BASE ?? ''
  let cuttedName = name?.includes(preview) ? name.replace(preview, '') : name

  cuttedName = cuttedName ? cuttedName.split('.')[0] : ''

  await api.del<string>('/file', { name: cuttedName })
}
