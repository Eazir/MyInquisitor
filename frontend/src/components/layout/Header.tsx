import { useTheme } from '../../contexts/ThemeContext';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage } from '../../contexts/LanguageContext';

interface HeaderProps {
  title: string;
  onToggleSidebar: () => void;
}

export function Header({ title, onToggleSidebar }: HeaderProps) {
  const { theme, toggleTheme } = useTheme();
  const { user } = useAuth();
  const { t } = useLanguage();

  return (
    <header
      className="fixed top-0 right-0 z-30 flex items-center justify-between px-4 md:px-6 lg:px-8 border-b border-[var(--color-border)] bg-[var(--color-bg-primary)] shadow-[var(--shadow-sm)] transition-all duration-300 lg:left-[var(--sidebar-width)]"
      style={{
        height: 'var(--header-height)',
      }}
    >
      <div className="flex items-center gap-3">
        <button
          onClick={onToggleSidebar}
          className="lg:hidden p-2 rounded-[var(--radius-md)] text-[var(--color-text-secondary)] hover:bg-[var(--color-bg-secondary)] transition-colors"
          aria-label={t('header.toggleMenu')}
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <h1 className="text-lg md:text-xl font-semibold text-[var(--color-text-primary)]">{title}</h1>
      </div>

      <div className="flex items-center gap-4">
        {user && (
          <span className="text-sm text-[var(--color-text-secondary)] hidden lg:block">
            {user.full_name}
          </span>
        )}

        <button
          onClick={toggleTheme}
          className="p-2 rounded-[var(--radius-md)] text-[var(--color-text-secondary)] hover:bg-[var(--color-bg-secondary)] transition-colors"
          title={t('header.switchTheme', { mode: theme === 'light' ? t('header.dark') : t('header.light') })}
        >
          {theme === 'light' ? '🌙' : '☀️'}
        </button>
      </div>
    </header>
  );
}
