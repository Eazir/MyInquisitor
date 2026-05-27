import { NavLink } from 'react-router-dom';
import { cn } from '../../utils/cn';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage } from '../../contexts/LanguageContext';

interface SidebarProps {
  open: boolean;
  onClose: () => void;
}

export function Sidebar({ open, onClose }: SidebarProps) {
  const { isAdmin, user, logout } = useAuth();
  const { t } = useLanguage();

  const navItems = [
    { to: '/dashboard', label: t('sidebar.dashboard'), icon: '📊' },
    { to: '/debts', label: t('sidebar.debts'), icon: '💳' },
    { to: '/expenses', label: t('sidebar.expenses'), icon: '📄' },
    { to: '/accounting', label: t('sidebar.accounting'), icon: '💰' },
    { to: '/planning', label: t('sidebar.planning'), icon: '📈' },
    { to: '/admin', label: t('sidebar.admin'), icon: '⚙️', adminOnly: true },
  ];

  const handleLogout = () => {
    logout();
    window.location.href = '/login';
  };

  return (
    <aside
      className={cn(
        'fixed inset-y-0 left-0 z-40 flex flex-col',
        'transform transition-transform duration-300 ease-in-out',
        'w-[var(--sidebar-width)]',
        'lg:translate-x-0',
        open ? 'translate-x-0' : '-translate-x-full'
      )}
      style={{
        backgroundColor: 'var(--color-sidebar-bg)',
        color: 'var(--color-sidebar-text)',
      }}
    >
      <div className="flex items-center gap-3 px-6 py-5 border-b border-white/10">
        <span className="text-2xl">{'🕵️'}</span>
        <span className="text-lg font-bold text-white">{t('sidebar.myInquisitor')}</span>
      </div>

      <nav className="flex-1 overflow-y-auto py-4 px-3 space-y-0.5">
        {navItems.map(item => {
          if (item.adminOnly && !isAdmin) return null;
          return (
            <NavLink
              key={item.to}
              to={item.to}
              onClick={onClose}
              className={({ isActive }) => cn(
                'flex items-center gap-3 px-3 py-2.5 rounded-[var(--radius-md)] text-sm font-medium transition-all duration-150',
                isActive
                  ? 'bg-[var(--color-sidebar-active)] text-white shadow-sm'
                  : 'text-[var(--color-sidebar-text)] hover:bg-[var(--color-sidebar-hover)] hover:text-white'
              )}
            >
              <span className="text-lg">{item.icon}</span>
              <span>{item.label}</span>
            </NavLink>
          );
        })}
      </nav>

      <div className="px-4 py-4 space-y-1 border-t border-white/10">
        <NavLink
          to="/settings"
          onClick={onClose}
          className={({ isActive }) => cn(
            'flex items-center gap-3 px-3 py-2.5 rounded-[var(--radius-md)] text-sm font-medium transition-all duration-150',
            isActive
              ? 'bg-[var(--color-sidebar-active)] text-white shadow-sm'
              : 'text-[var(--color-sidebar-text)] hover:bg-[var(--color-sidebar-hover)] hover:text-white'
          )}
        >
          <span className="text-lg">{'⚙️'}</span>
          <span>{t('sidebar.settings')}</span>
        </NavLink>
        {user && (
          <div className="px-3 py-2 text-xs text-[var(--color-sidebar-text)] truncate">
            {user.full_name}
          </div>
        )}
        <button
          onClick={handleLogout}
          className="flex items-center gap-3 w-full px-3 py-2.5 rounded-[var(--radius-md)] text-sm font-medium text-[var(--color-sidebar-text)] hover:bg-[var(--color-sidebar-hover)] hover:text-white transition-all duration-150"
        >
          <span className="text-lg">{'🚪'}</span>
          <span>{t('sidebar.logout')}</span>
        </button>
      </div>
    </aside>
  );
}
