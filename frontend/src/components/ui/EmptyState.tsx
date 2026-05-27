import type { ReactNode } from 'react';
import { cn } from '../../utils/cn';

interface EmptyStateProps {
  icon?: ReactNode;
  title: string;
  description?: string;
  action?: ReactNode;
}

export function EmptyState({ icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-20 px-6 text-center border-2 border-dashed border-[var(--color-border)] rounded-[var(--radius-lg)] bg-[var(--color-bg-primary)]/50">
      {icon ? (
        <div className="mb-5 text-4xl">{icon}</div>
      ) : (
        <div className="mb-5 w-16 h-16 rounded-[var(--radius-xl)] bg-[var(--color-bg-secondary)] flex items-center justify-center text-3xl">
          {'📭'}
        </div>
      )}
      <h3 className="text-lg font-semibold text-[var(--color-text-primary)] mb-2">{title}</h3>
      {description && (
        <p className={cn('text-sm text-[var(--color-text-secondary)] mb-6 max-w-sm')}>
          {description}
        </p>
      )}
      {action && <div className="mt-2">{action}</div>}
    </div>
  );
}
