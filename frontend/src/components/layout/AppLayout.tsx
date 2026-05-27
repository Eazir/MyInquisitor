import { Outlet, useLocation } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Header } from './Header';
import { useLanguage } from '../../contexts/LanguageContext';

export function AppLayout() {
  const { t } = useLanguage();
  const location = useLocation();

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
      <Sidebar />
      <Header title={title} />
      <main
        className="pt-16"
        style={{ marginLeft: 'var(--sidebar-width)' }}
      >
        <Outlet />
      </main>
    </div>
  );
}
