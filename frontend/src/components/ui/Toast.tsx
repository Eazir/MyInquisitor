import { useState, useEffect, useCallback } from 'react';
import { cn } from '../../utils/cn';

export type ToastVariant = 'success' | 'error' | 'info';

interface ToastData {
  id: number;
  message: string;
  variant: ToastVariant;
}

let toastId = 0;
let addToastFn: ((message: string, variant: ToastVariant) => void) | null = null;

export function toast(message: string, variant: ToastVariant = 'info') {
  addToastFn?.(message, variant);
}

export function ToastContainer() {
  const [toasts, setToasts] = useState<ToastData[]>([]);

  const addToast = useCallback((message: string, variant: ToastVariant) => {
    const id = ++toastId;
    setToasts(prev => [...prev, { id, message, variant }]);
    setTimeout(() => {
      setToasts(prev => prev.filter(t => t.id !== id));
    }, 4000);
  }, []);

  useEffect(() => {
    addToastFn = addToast;
    return () => { addToastFn = null; };
  }, [addToast]);

  if (toasts.length === 0) return null;

  const variants: Record<ToastVariant, string> = {
    success: 'bg-green-600 text-white',
    error: 'bg-[var(--color-danger)] text-white',
    info: 'bg-[var(--color-accent)] text-white',
  };

  return (
    <div className="fixed top-4 right-4 z-[100] flex flex-col gap-2 pointer-events-none">
      {toasts.map(t => (
        <div
          key={t.id}
          className={cn(
            'pointer-events-auto px-5 py-3 rounded-[var(--radius-md)] shadow-lg',
            'text-sm font-medium animate-in slide-in-from-right-2',
            'transition-all duration-300',
            variants[t.variant]
          )}
        >
          {t.message}
        </div>
      ))}
    </div>
  );
}
