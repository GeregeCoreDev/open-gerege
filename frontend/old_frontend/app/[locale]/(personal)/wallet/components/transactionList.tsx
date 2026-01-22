import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

export function TransactionList({
  data,
}: {
  data: {
    id: number
    description: string
    amount: number
    tran_type: 'C' | 'D'
    created_at: string
  }[]
}) {
  if (!data.length)
    return (
      <div className="text-muted-foreground flex flex-col items-center justify-center py-10 text-sm">
        Гүйлгээ олдсонгүй
      </div>
    )

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Тайлбар</TableHead>
          <TableHead>Огноо</TableHead>
          <TableHead className="text-right">Дүн</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.map((t) => (
          <TableRow key={t.id}>
            <TableCell>{t.description}</TableCell>
            <TableCell>{t.created_at}</TableCell>
            <TableCell
              className={`text-right font-medium ${
                t.tran_type === 'C' ? 'text-red-500' : 'text-green-500'
              }`}
            >
              {t.tran_type === 'C' ? '-' : '+'}
              {t.amount.toLocaleString()} ₮
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
