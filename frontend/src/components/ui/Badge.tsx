import { cn } from '../../utils/cn';

interface BadgeProps {
  variant?: 'info' | 'success' | 'warning' | 'danger';
  size?: 'sm' | 'md';
  removable?: boolean;
  onRemove?: () => void;
  children: string;
}

export function Badge({ variant = 'info', size = 'md', removable, onRemove, children }: BadgeProps) {
  const variants = {
    info: 'bg-[var(--color-info)]/10 text-[var(--color-info)] border-[var(--color-info)]/20',
    success: 'bg-[var(--color-success)]/10 text-[var(--color-success)] border-[var(--color-success)]/20',
    warning: 'bg-[var(--color-warning)]/10 text-[var(--color-warning)] border-[var(--color-warning)]/20',
    danger: 'bg-[var(--color-danger)]/10 text-[var(--color-danger)] border-[var(--color-danger)]/20',
  };

  const sizes = {
    sm: 'px-2 py-0.5 text-[var(--font-size-xs)]',
    md: 'px-2.5 py-1 text-[var(--font-size-sm)]',
  };

  return (
    <span className={cn(
      'inline-flex items-center gap-1 rounded-full border font-medium',
      variants[variant], sizes[size]
    )}>
      {children}
      {removable && onRemove && (
        <button onClick={onRemove} className="ml-1 hover:opacity-70">&times;</button>
      )}
    </span>
  );
}
