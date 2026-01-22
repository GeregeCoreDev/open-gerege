/**
 * 📱 Business Icon Page (/[locale]/(main)/business/icon/page.tsx)
 *
 * Энэ нь Business app icon удирдах хуудас юм.
 * Зорилго: Business app-ийн service icons удирдлага
 *
 * Architecture:
 * - Master-detail view: Select icon from list
 * - View children icons
 * - Hierarchical icon structure
 *
 * Features:
 * - ✅ Icon listing with search
 * - ✅ Parent-child relationship
 * - ✅ View details dialog
 * - ✅ Responsive table
 *
 * API: GET /business-app-service-icon
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

// app/[locale]/app-business-icon/page.tsx
'use client'

import * as React from 'react'
import { useState, useEffect, useMemo, useCallback } from 'react'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Checkbox } from '@/components/ui/checkbox'
import { Progress } from '@/components/ui/progress'

import {
  Loader2,
  Plus,
  Pencil,
  Trash2,
  ChevronRight,
  ChevronDown,
  MoveUp,
  MoveDown,
  Search,
} from 'lucide-react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import Image from 'next/image'

/** ===== Types ===== */
type IconNode = {
  id: number
  key: string
  name: string
  description?: string | null
  icon?: string | null
  link?: string | null
  web_link?: string | null
  is_native?: boolean | null
  is_public?: boolean | null
  sequence?: number
  parent_id?: number | null
  children?: IconNode[]
}

// DTO (create/update)
const IconDtoSchema = z.object({
  key: z.string().min(1).max(255),
  name: z.string().min(1).max(255),
  description: z.string().max(255).optional().or(z.literal('')),
  icon: z.string().max(255).optional().or(z.literal('')),
  link: z.string().max(255).optional().or(z.literal('')),
  web_link: z.string().max(255).optional().or(z.literal('')),
  is_native: z.boolean().optional().default(false),
  is_public: z.boolean().optional().default(true),
  sequence: z.coerce.number().int().min(0).default(0),
  parent_id: z.number().int().optional().nullable(),
})
type IconDtoIn = z.input<typeof IconDtoSchema>
type IconDtoOut = z.output<typeof IconDtoSchema>

/** ===== Page wrapper ===== */
export default function BusinessIconPage() {
  const t = useTranslations()
  return (
    <div className="h-full w-full p-6">
      <BusinessIconTree t={t} />
    </div>
  )
}

/** ===== Tree Component ===== */
function BusinessIconTree({ t: _t }: { t: ReturnType<typeof useTranslations> }) {
  const [tree, setTree] = useState<IconNode[]>([])
  const [loading, setLoading] = useState(false)
  const [progress, setProgress] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [expanded, setExpanded] = useState<Set<number>>(new Set())
  const [filter, setFilter] = useState('')

  // progress bar
  useEffect(() => {
    let timer: ReturnType<typeof setInterval> | undefined
    if (loading) {
      setProgress(0)
      timer = setInterval(() => setProgress((p) => Math.min(90, p + Math.random() * 12 + 8)), 250)
    } else {
      setProgress(100)
      const id = setTimeout(() => setProgress(0), 300)
      return () => clearTimeout(id)
    }
    return () => timer && clearInterval(timer)
  }, [loading])

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await api.get<IconNode[]>('/app-business-icon', { cache: 'no-store' })
      setTree(Array.isArray(data) ? data : [])
    } catch {
      setError("Error occurred")
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const toggleExpand = (id: number) => {
    setExpanded((prev) => {
      const n = new Set(prev)
      if (n.has(id)) {
        n.delete(id)
      } else {
        n.add(id)
      }
      return n
    })
  }

  /** ------- Filter (name/key contains) ------- */
  const filteredTree = useMemo(() => {
    const q = filter.trim().toLowerCase()
    if (!q) return tree
    const dfs = (nodes: IconNode[]): IconNode[] => {
      const out: IconNode[] = []
      for (const n of nodes) {
        const hit =
          (n.name ?? '').toLowerCase().includes(q) ||
          (n.key ?? '').toLowerCase().includes(q) ||
          (n.description ?? '').toLowerCase().includes(q)
        const kids = n.children?.length ? dfs(n.children) : []
        if (hit || kids.length) {
          out.push({ ...n, children: kids })
          if (hit) expanded.add(n.id) // auto-expand hits
        }
      }
      return out
    }
    // mutate expanded intentionally to open hits
    return dfs(tree)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tree, filter])

  /** ------- Create / Edit dialog ------- */
  const [editing, setEditing] = useState<{
    mode: 'create' | 'edit'
    node?: IconNode
    parent?: IconNode | null
  } | null>(null)
  const form = useForm<IconDtoIn>({
    resolver: zodResolver(IconDtoSchema),
    defaultValues: {
      key: '',
      name: '',
      description: '',
      icon: '',
      link: '',
      web_link: '',
      is_native: false,
      is_public: true,
      sequence: 0,
      parent_id: undefined,
    },
  })

  const openCreateRoot = () => {
    setEditing({ mode: 'create', parent: null })
    form.reset({
      key: '',
      name: '',
      description: '',
      icon: '',
      link: '',
      web_link: '',
      is_native: false,
      is_public: true,
      sequence: 0,
      parent_id: undefined,
    })
  }
  const openCreateChild = (parent: IconNode) => {
    setEditing({ mode: 'create', parent })
    form.reset({
      key: '',
      name: '',
      description: '',
      icon: '',
      link: '',
      web_link: '',
      is_native: false,
      is_public: true,
      sequence: (parent.children?.length || 0) + 1,
      parent_id: parent.id,
    })
  }
  const openEdit = (node: IconNode) => {
    setEditing({ mode: 'edit', node })
    form.reset({
      key: node.key || '',
      name: node.name || '',
      description: node.description || '',
      icon: node.icon || '',
      link: node.link || '',
      web_link: node.web_link || '',
      is_native: !!node.is_native,
      is_public: node.is_public ?? true,
      sequence: node.sequence ?? 0,
      parent_id: node.parent_id ?? undefined,
    })
  }

  const closeDialog = () => setEditing(null)

  const submit = form.handleSubmit(async (values) => {
    const payload: IconDtoOut = IconDtoSchema.parse(values)
    try {
      if (editing?.mode === 'create') {
        await api.post('/app-business-icon', payload as Record<string, unknown>)
      } else if (editing?.mode === 'edit' && editing.node) {
        await api.put(`/app-business-icon/${editing.node.id}`, {
          id: editing.node.id,
          ...payload,
        } as Record<string, unknown>)
      }
      closeDialog()
      await load()
    } catch {}
  })

  /** ------- Delete ------- */
  const del = async (node: IconNode) => {
    if (!confirm(`Устгах уу? — ${node.name}`)) return
    await api.del(`/app-business-icon/${node.id}`)
    await load()
  }

  /** ------- Sequence up/down (optimistic then PUT) ------- */
  const move = async (node: IconNode, dir: 'up' | 'down') => {
    const siblings = findSiblings(tree, node.parent_id)
    const idx = siblings.findIndex((n) => n.id === node.id)
    const nxt = dir === 'up' ? idx - 1 : idx + 1
    if (nxt < 0 || nxt >= siblings.length) return
    // swap sequence
    const a = siblings[idx]
    const b = siblings[nxt]
    const aSeq = a.sequence ?? idx
    const bSeq = b.sequence ?? nxt
    // optimistic
    siblings[idx] = b
    siblings[nxt] = a
    setTree([...tree])
    try {
      await Promise.all([
        api.put(`/app-business-icon/${a.id}`, { ...a, sequence: bSeq } as Record<string, unknown>),
        api.put(`/app-business-icon/${b.id}`, { ...b, sequence: aSeq } as Record<string, unknown>),
      ])
      await load()
    } catch {
      // reload on error
      await load()
    }
  }

  return (
    <>
      {progress > 0 && (
        <div className="absolute inset-x-0 top-0">
          <Progress value={progress} className="h-1 rounded-none" />
        </div>
      )}

      <div className="flex flex-col gap-3 pb-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Business Icons</h1>
        <div className="flex items-center gap-2">
          <div className="relative">
            <Search className="absolute top-1/2 left-2 h-4 w-4 -translate-y-1/2 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Нэр, key..."
              className="h-9 pl-8 sm:w-64"
            />
          </div>
          <Button onClick={openCreateRoot} className="gap-2">
            <Plus className="h-4 w-4" /> Нэмэх (үндсэн)
          </Button>
        </div>
      </div>

      <Separator />

      <div className="min-h-0 flex-1 overflow-auto p-3">
        {error ? (
          <div className="text-sm text-red-600">{error}</div>
        ) : loading ? (
          <div className="flex h-32 items-center justify-center">
            <Loader2 className="h-5 w-5 animate-spin" />
          </div>
        ) : filteredTree.length === 0 ? (
          <div className="text-muted-foreground text-sm">Хоосон байна.</div>
        ) : (
          <ul className="space-y-1">
            {filteredTree
              .sort((a, b) => (a.sequence ?? 0) - (b.sequence ?? 0))
              .map((n) => (
                <TreeRow
                  key={n.id}
                  node={n}
                  level={0}
                  expanded={expanded}
                  onToggle={toggleExpand}
                  onAddChild={openCreateChild}
                  onEdit={openEdit}
                  onDelete={del}
                  onMove={move}
                />
              ))}
          </ul>
        )}
      </div>

      {/* Create/Edit Dialog */}
      <Dialog open={!!editing} onOpenChange={(v) => (v ? null : closeDialog())}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>{editing?.mode === 'create' ? 'Icon нэмэх' : 'Icon засах'}</DialogTitle>
            <DialogDescription>
              {editing?.mode === 'create'
                ? editing?.parent
                  ? `Эцэг: ${editing.parent.name}`
                  : 'Үндсэн түвшинд нэмнэ.'
                : 'Талбаруудаа засч хадгална уу.'}
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={submit} className="space-y-3">
            <div className="grid gap-3 sm:grid-cols-2">
              <div>
                <label className="mb-1 block text-sm">Key</label>
                <Input {...form.register('key')} />
                <FieldErr form={form} name="key" />
              </div>
              <div>
                <label className="mb-1 block text-sm">Нэр</label>
                <Input {...form.register('name')} />
                <FieldErr form={form} name="name" />
              </div>
              <div>
                <label className="mb-1 block text-sm">Icon (class/url)</label>
                <Input placeholder="e.g. lucide:home" {...form.register('icon')} />
              </div>
              <div>
                <label className="mb-1 block text-sm">App Link</label>
                <Input placeholder="/path эсвэл app://" {...form.register('link')} />
              </div>
              <div>
                <label className="mb-1 block text-sm">Web Link</label>
                <Input placeholder="https://..." {...form.register('web_link')} />
              </div>
              <div>
                <label className="mb-1 block text-sm">Sequence</label>
                <Input type="number" {...form.register('sequence', { valueAsNumber: true })} />
              </div>
              <div className="flex items-center gap-2">
                <Checkbox
                  checked={!!form.watch('is_native')}
                  onCheckedChange={(v) => form.setValue('is_native', Boolean(v))}
                />
                <span className="text-sm">Native</span>
              </div>
              <div className="flex items-center gap-2">
                <Checkbox
                  checked={form.watch('is_public') ?? true}
                  onCheckedChange={(v) => form.setValue('is_public', Boolean(v))}
                />
                <span className="text-sm">Public</span>
              </div>
            </div>

            <div>
              <label className="mb-1 block text-sm">Тайлбар</label>
              <Input {...form.register('description')} />
            </div>

            {/* parent_id (read-only when editing) */}
            <div className="text-xs opacity-70">
              Parent ID:{' '}
              {editing?.mode === 'edit'
                ? (editing.node?.parent_id ?? '—')
                : (editing?.parent?.id ?? '—')}
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={closeDialog}>
                Хаах
              </Button>
              <Button type="submit">Хадгалах</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  )
}

/** ===== Tree row (recursive) ===== */
function TreeRow(props: {
  node: IconNode
  level: number
  expanded: Set<number>
  onToggle: (id: number) => void
  onAddChild: (parent: IconNode) => void
  onEdit: (node: IconNode) => void
  onDelete: (node: IconNode) => void
  onMove: (node: IconNode, dir: 'up' | 'down') => void
}) {
  const { node, level, expanded, onToggle, onAddChild, onEdit, onDelete, onMove } = props
  const hasChildren = (node.children?.length ?? 0) > 0
  const isOpen = expanded.has(node.id)

  return (
    <li>
      <div
        className="group hover:bg-muted/50 flex items-center gap-2 rounded-md px-2 py-1"
        style={{ paddingLeft: 8 + level * 20 }}
      >
        <button
          className="hover:bg-muted/70 flex h-6 w-6 items-center justify-center rounded"
          onClick={() => (hasChildren ? onToggle(node.id) : null)}
          aria-label="toggle"
        >
          {hasChildren ? (
            isOpen ? (
              <ChevronDown className="h-4 w-4" />
            ) : (
              <ChevronRight className="h-4 w-4" />
            )
          ) : (
            <span className="inline-block w-4" />
          )}
        </button>

        {node.icon ? <Image src={node.icon} alt={node.name || 'icon'} width={40} height={40} className="w-[40px]" /> : null}

        <span className="font-medium">{node.name}</span>

        <span className="text-muted-foreground ml-3 text-xs">seq: {node.sequence ?? 0}</span>

        {node.is_native ? (
          <Badge variant="outline" className="ml-2">
            native
          </Badge>
        ) : null}
        {node.is_public === false ? (
          <Badge className="ml-2" variant="destructive">
            private
          </Badge>
        ) : null}

        <div className="ml-auto flex items-center gap-1 opacity-0 transition group-hover:opacity-100">
          <Button
            size="icon"
            variant="ghost"
            onClick={() => onMove(node, 'up')}
            aria-label="move up"
          >
            <MoveUp className="h-4 w-4" />
          </Button>
          <Button
            size="icon"
            variant="ghost"
            onClick={() => onMove(node, 'down')}
            aria-label="move down"
          >
            <MoveDown className="h-4 w-4" />
          </Button>
          <Button
            size="icon"
            variant="ghost"
            onClick={() => onAddChild(node)}
            aria-label="add child"
          >
            <Plus className="h-4 w-4" />
          </Button>
          <Button size="icon" variant="ghost" onClick={() => onEdit(node)} aria-label="edit">
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            size="icon"
            variant="destructive"
            onClick={() => onDelete(node)}
            aria-label="delete"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {hasChildren && isOpen ? (
        <ul className="mt-1 space-y-1">
          {node
            .children!.slice()
            .sort((a, b) => (a.sequence ?? 0) - (b.sequence ?? 0))
            .map((c) => (
              <TreeRow
                key={c.id}
                node={c}
                level={level + 1}
                expanded={expanded}
                onToggle={onToggle}
                onAddChild={onAddChild}
                onEdit={onEdit}
                onDelete={onDelete}
                onMove={onMove}
              />
            ))}
        </ul>
      ) : null}
    </li>
  )
}

/** ===== helpers ===== */
function FieldErr({ form, name }: { form: { formState: { errors: Record<string, { message?: string }> } }; name: keyof IconDtoIn }) {
  const err = form.formState.errors[name]?.message as string | undefined
  return err ? <p className="mt-1 text-xs text-red-600">{err}</p> : null
}

function findSiblings(tree: IconNode[], parentId?: number | null): IconNode[] {
  if (!parentId) return tree
  const stack = [...tree]
  while (stack.length) {
    const n = stack.pop()!
    if (n.id === parentId) return n.children ?? (n.children = [])
    if (n.children?.length) stack.push(...n.children)
  }
  return tree
}
