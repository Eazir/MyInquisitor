import type { ReactNode } from 'react';
import { cn } from '../../utils/cn';

export interface Column<T = any> {
  key: string;
  header: string;
  render?: (item: T) => ReactNode;
  sortable?: boolean;
  className?: string;
}

interface TableProps<T = any> {
  columns: Column<T>[];
  data: T[];
  loading?: boolean;
  variant?: 'default' | 'compact' | 'striped';
  onRowClick?: (item: T) => void;
  emptyMessage?: string;
}

export function Table<T extends Record<string, any>>({
  columns,
  data,
  loading,
  variant = 'default',
  onRowClick,
  emptyMessage = 'No data available',
}: TableProps<T>) {
  return (
    <div className="overflow-x-auto rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-bg-primary)]">
      <table className="w-full text-[var(--font-size-sm)]">
        <thead>
          <tr className="bg-[var(--color-bg-secondary)] border-b border-[var(--color-border)]">
            {columns.map(col => (
              <th
                key={col.key}
                className={cn(
                  'px-6 py-4 text-left font-semibold text-[var(--color-text-secondary)] text-xs uppercase tracking-wider',
                  col.className
                )}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-[var(--color-border)]">
          {loading ? (
            <tr>
              <td colSpan={columns.length} className="px-6 py-12 text-center text-[var(--color-text-muted)]">
                Loading...
              </td>
            </tr>
          ) : data.length === 0 ? (
            <tr>
              <td colSpan={columns.length} className="px-6 py-12 text-center text-[var(--color-text-muted)]">
                {emptyMessage}
              </td>
            </tr>
          ) : (
            data.map((item, i) => (
              <tr
                key={item.id || i}
                className={cn(
                  'transition-colors hover:bg-[var(--color-bg-secondary)]/50',
                  variant === 'striped' && i % 2 === 1 && 'bg-[var(--color-bg-secondary)]/30',
                  onRowClick && 'cursor-pointer'
                )}
                onClick={() => onRowClick?.(item)}
              >
                {columns.map(col => (
                  <td key={col.key} className={cn('px-6 py-4 text-[var(--color-text-primary)]', col.className)}>
                    {col.render ? col.render(item) : item[col.key]}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
