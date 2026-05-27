import type { InputHTMLAttributes } from 'react';
import { forwardRef } from 'react';
import { cn } from '../../utils/cn';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, helperText, className, id, ...props }, ref) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, '-');

    return (
      <div className="flex flex-col gap-1.5">
        {label && (
          <label htmlFor={inputId} className="text-[var(--font-size-sm)] font-medium text-[var(--color-text-primary)]">
            {label}
          </label>
        )}
        <input
          ref={ref}
          id={inputId}
          className={cn(
            'w-full px-4 py-3 rounded-[var(--radius-md)] border text-[var(--font-size-base)]',
            'bg-[var(--color-bg-primary)] text-[var(--color-text-primary)]',
            'border-[var(--color-border)] placeholder:text-[var(--color-text-muted)]',
            'focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent',
            'transition-colors duration-200',
            error && 'border-[var(--color-danger)] focus:ring-[var(--color-danger)]',
            className
          )}
          {...props}
        />
        {error && <p className="text-[var(--font-size-xs)] text-[var(--color-danger)]">{error}</p>}
        {helperText && !error && <p className="text-[var(--font-size-xs)] text-[var(--color-text-muted)]">{helperText}</p>}
      </div>
    );
  }
);

Input.displayName = 'Input';
