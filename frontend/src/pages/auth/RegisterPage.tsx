import type { FormEvent } from 'react';
import { useState } from 'react';
import { Link, useNavigate, useParams, Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage } from '../../contexts/LanguageContext';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { Card } from '../../components/ui/Card';
import { Loading } from '../../components/ui/Loading';

export function RegisterPage() {
  const { token } = useParams<{ token: string }>();
  const navigate = useNavigate();
  const { register, isAuthenticated } = useAuth();
  const { t } = useLanguage();
  const [fullName, setFullName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  if (!token) {
    return (
    <div className="min-h-dvh flex items-center justify-center bg-[var(--color-bg-secondary)] p-8">
        <Card className="w-full max-w-md text-center">
          <h2 className="text-xl font-semibold text-[var(--color-text-primary)]">{t('auth.registrationClosed')}</h2>
          <p className="text-[var(--color-text-secondary)] mt-2">
            {t('auth.registrationClosedDescription')}
          </p>
          <Link to="/login" className="text-[var(--color-accent)] hover:underline mt-4 inline-block">
            {t('auth.goToLogin')}
          </Link>
        </Card>
      </div>
    );
  }

  if (loading && !fullName && !email && !password) {
    return <Loading text={t('auth.validatingInvite')} />;
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    if (password.length < 8) {
      setError(t('auth.passwordMinLength'));
      return;
    }
    setLoading(true);
    try {
      await register(email, password, fullName, token);
      navigate('/dashboard', { replace: true });
    } catch {
      setError(t('auth.registrationFailed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-dvh flex items-center justify-center bg-[var(--color-bg-secondary)] p-8">
      <Card className="w-full max-w-md">
        <div className="text-center mb-8">
          <span className="text-4xl">{'🕵️'}</span>
          <h1 className="text-2xl font-bold text-[var(--color-text-primary)] mt-2">{t('auth.createAccount')}</h1>
          <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('auth.createAccountDescription')}</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-danger)]/10 text-[var(--color-danger)] text-sm">
              {error}
            </div>
          )}

          <Input
            label={t('auth.fullName')}
            value={fullName}
            onChange={e => setFullName(e.target.value)}
            placeholder={t('auth.namePlaceholder')}
            required
          />

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
            placeholder={t('auth.atLeast8Chars')}
            required
          />

          <Button type="submit" className="w-full" loading={loading}>
            {t('auth.createAccount')}
          </Button>
        </form>

        <p className="text-center text-sm text-[var(--color-text-secondary)] mt-6">
          {t('auth.alreadyHaveAccount')}{' '}
          <Link to="/login" className="text-[var(--color-accent)] hover:underline">
            {t('auth.signIn')}
          </Link>
        </p>
      </Card>
    </div>
  );
}
