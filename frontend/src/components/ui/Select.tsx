import type { SelectHTMLAttributes } from 'react';
import { cn } from '../../utils/cn';

interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  options: SelectOption[];
  error?: string;
}

export function Select({ label, options, error, className, id, ...props }: SelectProps) {
  const selectId = id || label?.toLowerCase().replace(/\s+/g, '-');

  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label htmlFor={selectId} className="text-[var(--font-size-sm)] font-medium text-[var(--color-text-primary)]">
          {label}
        </label>
      )}
      <select
        id={selectId}
        className={cn(
          'w-full px-4 py-3 rounded-[var(--radius-md)] border text-[var(--font-size-base)]',
          'bg-[var(--color-bg-primary)] text-[var(--color-text-primary)]',
          'border-[var(--color-border)]',
          'focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent',
          'transition-colors duration-200',
          error && 'border-[var(--color-danger)]',
          className
        )}
        {...props}
      >
        {options.map(opt => (
          <option key={opt.value} value={opt.value}>{opt.label}</option>
        ))}
      </select>
      {error && <p className="text-[var(--font-size-xs)] text-[var(--color-danger)]">{error}</p>}
    </div>
  );
}
