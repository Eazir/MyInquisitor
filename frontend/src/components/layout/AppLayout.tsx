import { useState } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Header } from './Header';
import { useLanguage } from '../../contexts/LanguageContext';

export function AppLayout() {
  const { t } = useLanguage();
  const location = useLocation();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const toggleSidebar = () => setSidebarOpen(prev => !prev);
  const closeSidebar = () => setSidebarOpen(false);

  const routeTitles: Record<string, string> = {
    '/dashboard': t('dashboard.title'),
    '/debts': t('debts.title'),
    '/expenses': t('expenses.title'),
    '/accounting': t('accounting.title'),
    '/planning': t('planning.title'),
    '/admin': t('admin.title'),
    '/settings': t('settings.title'),
  };

  const title = routeTitles[location.pathname] || t('sidebar.myInquisitor');

  return (
    <div className="min-h-screen bg-[var(--color-bg-secondary)]">
      <Sidebar open={sidebarOpen} onClose={closeSidebar} />
      <Header title={title} onToggleSidebar={toggleSidebar} />
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-30 lg:hidden"
          onClick={closeSidebar}
        />
      )}
      <main
        className="pt-16 transition-all duration-300 lg:ml-[var(--sidebar-width)]"
      >
        <Outlet />
      </main>
    </div>
  );
}
