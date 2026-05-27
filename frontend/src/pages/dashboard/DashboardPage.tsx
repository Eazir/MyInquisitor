import { Card } from '../../components/ui/Card';
import { PageContainer } from '../../components/layout/PageContainer';
import { useLanguage } from '../../contexts/LanguageContext';

export function DashboardPage() {
  const { t } = useLanguage();

  const stats = [
    { label: t('dashboard.totalManaged'), value: '$0.00' },
    { label: t('dashboard.available'), value: '$0.00' },
    { label: t('dashboard.reserved'), value: '$0.00' },
  ];

  return (
    <PageContainer>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('dashboard.title')}</h2>
        <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('dashboard.description')}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
        {stats.map(s => (
          <Card key={s.label} variant="stats">
            <p className="text-sm text-[var(--color-text-secondary)] mb-1">{s.label}</p>
            <p className="text-3xl font-bold text-[var(--color-text-primary)]">{s.value}</p>
          </Card>
        ))}
      </div>

      <Card title={t('dashboard.upcomingPayments')} className="mb-8">
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <span className="text-4xl mb-4">{'📅'}</span>
          <p className="text-[var(--color-text-muted)] text-sm">
            {t('dashboard.noUpcomingPayments')}
          </p>
        </div>
      </Card>

      <Card title={t('dashboard.balanceTrend')}>
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <span className="text-4xl mb-4">{'📈'}</span>
          <p className="text-[var(--color-text-muted)] text-sm">
            {t('dashboard.chartComingSoon')}
          </p>
        </div>
      </Card>
    </PageContainer>
  );
}
