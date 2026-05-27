import type { ReactNode } from 'react';

interface PageContainerProps {
  children: ReactNode;
}

export function PageContainer({ children }: PageContainerProps) {
  return (
    <div className="p-4 md:p-6 lg:p-8 max-w-7xl mx-auto">
      {children}
    </div>
  );
}
