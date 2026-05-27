import type { ReactNode } from 'react';
import { cn } from '../../utils/cn';

interface CardProps {
  title?: string;
  subtitle?: string;
  children: ReactNode;
  variant?: 'default' | 'stats' | 'highlight';
  className?: string;
}

export function Card({ title, subtitle, children, variant = 'default', className }: CardProps) {
  return (
    <div
      className={cn(
        'rounded-[var(--radius-lg)] bg-[var(--color-bg-primary)]',
        'border border-[var(--color-border)]',
        'shadow-[var(--shadow-card)]',
        variant === 'stats' && 'border-l-4 border-l-[var(--color-accent)]',
        variant === 'highlight' && 'bg-[var(--color-accent)] text-white border-none',
        className
      )}
      data-variant={variant}
    >
      {title && (
        <div className="px-6 md:px-8 pt-6 md:pt-8 pb-4 md:pb-5 border-b border-[var(--color-border)]">
          <h3 className={cn(
            'text-lg font-semibold',
            variant === 'highlight' ? 'text-white' : 'text-[var(--color-text-primary)]'
          )}>
            {title}
          </h3>
          {subtitle && (
            <p className={cn(
              'text-sm mt-1',
              variant === 'highlight' ? 'text-white/80' : 'text-[var(--color-text-secondary)]'
            )}>
              {subtitle}
            </p>
          )}
        </div>
      )}
      {!title && subtitle && (
        <div className="px-6 md:px-8 pt-6 md:pt-8 pb-4">
          <p className="text-sm text-[var(--color-text-secondary)]">{subtitle}</p>
        </div>
      )}
      <div className={cn('p-6 md:p-8', !title && 'pt-6 md:pt-8')}>{children}</div>
    </div>
  );
}
