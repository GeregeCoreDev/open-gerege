'use client'
import * as React from 'react'

export default function PinPad({
  value,
  onChange,
}: {
  value: string
  onChange: (v: string) => void
}) {
  const enter = (d: string) => {
    if (value.length < 4) onChange(value + d)
  }
  const back = () => onChange(value.slice(0, -1))
  const Dot = ({ i }: { i: number }) => (
    <div
      className={`h-3 w-3 rounded-full ${value.length > i ? 'bg-primary' : 'bg-muted-foreground/30'}`}
    />
  )
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-center gap-3">
        <Dot i={0} />
        <Dot i={1} />
        <Dot i={2} />
        <Dot i={3} />
      </div>
      {[
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9],
        ['', 0, '←'],
      ].map((row, idx) => (
        <div key={idx} className="grid grid-cols-3 gap-1">
          {row.map((n, i) => (
            <button
              key={i}
              className="h-10 rounded-md border text-lg font-medium"
              onClick={() => (n === '←' ? back() : n !== '' ? enter(String(n)) : undefined)}
            >
              {n}
            </button>
          ))}
        </div>
      ))}
    </div>
  )
}
