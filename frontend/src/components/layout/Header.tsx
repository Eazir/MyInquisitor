import { useTheme } from '../../contexts/ThemeContext';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage } from '../../contexts/LanguageContext';

interface HeaderProps {
  title: string;
}

export function Header({ title }: HeaderProps) {
  const { theme, toggleTheme } = useTheme();
  const { user } = useAuth();
  const { t } = useLanguage();

  return (
    <header
      className="fixed top-0 right-0 z-30 flex items-center justify-between px-8 border-b border-[var(--color-border)] bg-[var(--color-bg-primary)] shadow-[var(--shadow-sm)]"
      style={{
        left: 'var(--sidebar-width)',
        height: 'var(--header-height)',
      }}
    >
      <h1 className="text-xl font-semibold text-[var(--color-text-primary)]">{title}</h1>

      <div className="flex items-center gap-4">
        {user && (
          <span className="text-sm text-[var(--color-text-secondary)] hidden sm:block">
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
