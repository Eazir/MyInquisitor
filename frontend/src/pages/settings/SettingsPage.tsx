import type { FormEvent } from 'react';
import { useState } from 'react';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { Select } from '../../components/ui/Select';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage, type Language } from '../../contexts/LanguageContext';
import { updateProfile, changePassword } from '../../services/profile';

export function SettingsPage() {
  const { user, setUser } = useAuth();
  const { t, language, setLanguage } = useLanguage();
  const [profileForm, setProfileForm] = useState({
    full_name: user?.full_name || '',
    email: user?.email || '',
  });
  const [profileMsg, setProfileMsg] = useState('');
  const [profileError, setProfileError] = useState('');
  const [profileLoading, setProfileLoading] = useState(false);

  const [passwordForm, setPasswordForm] = useState({
    current_password: '',
    new_password: '',
  });
  const [passwordMsg, setPasswordMsg] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [passwordLoading, setPasswordLoading] = useState(false);

  const handleProfileSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setProfileMsg('');
    setProfileError('');
    setProfileLoading(true);
    try {
      const result = await updateProfile(profileForm);
      setUser(result.user);
      setProfileMsg(t('settings.profileUpdated'));
    } catch {
      setProfileError('Failed to update profile');
    } finally {
      setProfileLoading(false);
    }
  };

  const handlePasswordSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setPasswordMsg('');
    setPasswordError('');
    if (passwordForm.new_password.length < 8) {
      setPasswordError(t('auth.passwordMinLength'));
      return;
    }
    setPasswordLoading(true);
    try {
      await changePassword(passwordForm);
      setPasswordMsg(t('settings.passwordUpdated'));
      setPasswordForm({ current_password: '', new_password: '' });
    } catch {
      setPasswordError(t('settings.currentPasswordIncorrect'));
    } finally {
      setPasswordLoading(false);
    }
  };

  return (
    <PageContainer>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('settings.title')}</h2>
        <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('settings.description')}</p>
      </div>

      <Card title={t('settings.profile')} subtitle={t('settings.profileDescription')} className="mb-8">
        <form onSubmit={handleProfileSubmit} className="space-y-5">
          {profileMsg && (
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-success)]/10 text-[var(--color-success)] text-sm">
              {profileMsg}
            </div>
          )}
          {profileError && (
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-danger)]/10 text-[var(--color-danger)] text-sm">
              {profileError}
            </div>
          )}
          <Input
            label={t('settings.fullName')}
            value={profileForm.full_name}
            onChange={e => setProfileForm(p => ({ ...p, full_name: e.target.value }))}
            required
          />
          <Input
            label={t('settings.email')}
            type="email"
            value={profileForm.email}
            onChange={e => setProfileForm(p => ({ ...p, email: e.target.value }))}
            required
          />
          <Button type="submit" loading={profileLoading}>{t('common.save')}</Button>
        </form>
      </Card>

      <Card title={t('settings.password')} subtitle={t('settings.passwordDescription')} className="mb-8">
        <form onSubmit={handlePasswordSubmit} className="space-y-5">
          {passwordMsg && (
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-success)]/10 text-[var(--color-success)] text-sm">
              {passwordMsg}
            </div>
          )}
          {passwordError && (
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-danger)]/10 text-[var(--color-danger)] text-sm">
              {passwordError}
            </div>
          )}
          <Input
            label={t('settings.currentPassword')}
            type="password"
            value={passwordForm.current_password}
            onChange={e => setPasswordForm(p => ({ ...p, current_password: e.target.value }))}
            required
          />
          <Input
            label={t('settings.newPassword')}
            type="password"
            value={passwordForm.new_password}
            onChange={e => setPasswordForm(p => ({ ...p, new_password: e.target.value }))}
            required
          />
          <Button type="submit" loading={passwordLoading}>{t('settings.changePassword')}</Button>
        </form>
      </Card>

      <Card title={t('settings.language')} subtitle={t('settings.languageDescription')}>
        <Select
          label={t('settings.language')}
          value={language}
          options={[
            { value: 'es', label: 'Español' },
            { value: 'en', label: 'English' },
          ]}
          onChange={e => setLanguage(e.target.value as Language)}
        />
      </Card>
    </PageContainer>
  );
}
