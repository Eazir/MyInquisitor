import type { FormEvent } from 'react';
import { useState } from 'react';
import { Link, useNavigate, Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage } from '../../contexts/LanguageContext';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { Card } from '../../components/ui/Card';

export function LoginPage() {
  const navigate = useNavigate();
  const { login, isAuthenticated } = useAuth();
  const { t } = useLanguage();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(email, password);
      navigate('/dashboard', { replace: true });
    } catch {
      setError(t('auth.invalidCredentials'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-dvh flex items-center justify-center bg-[var(--color-bg-secondary)] p-8">
      <Card className="w-full max-w-md">
        <div className="text-center mb-8">
          <span className="text-4xl">🕵️</span>
          <h1 className="text-2xl font-bold text-[var(--color-text-primary)] mt-2">{t('sidebar.myInquisitor')}</h1>
          <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('auth.signInTitle')}</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-danger)]/10 text-[var(--color-danger)] text-sm">
              {error}
            </div>
          )}

          <Input
            label={t('auth.email')}
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            placeholder={t('auth.yourEmail')}
            required
          />

          <Input
            label={t('auth.password')}
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            placeholder={t('auth.enterPassword')}
            required
          />

          <Button type="submit" className="w-full" loading={loading}>
            {t('auth.signIn')}
          </Button>
        </form>

        <p className="text-center text-sm text-[var(--color-text-secondary)] mt-6">
          {t('auth.dontHaveAccount')}{' '}
          <Link to="/register" className="text-[var(--color-accent)] hover:underline">
            {t('auth.register')}
          </Link>
        </p>
      </Card>
    </div>
  );
}
