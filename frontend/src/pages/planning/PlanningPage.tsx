import { useState, useEffect } from 'react';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { Table, type Column } from '../../components/ui/Table';
import { Loading } from '../../components/ui/Loading';
import { useLanguage } from '../../contexts/LanguageContext';
import { accountingApi, type Projection } from '../../services/accounting';

export function PlanningPage() {
  const { t } = useLanguage();
  const [projections, setProjections] = useState<Projection[]>([]);
  const [loading, setLoading] = useState(true);
  const [months, setMonths] = useState(6);

  const columns: Column<Projection>[] = [
    { key: 'month', header: t('planning.month') },
    {
      key: 'projected_income',
      header: t('planning.income'),
      render: (p) => <span className="font-medium text-[var(--color-success)]">${p.projected_income.toFixed(2)}</span>,
    },
    {
      key: 'projected_expenses',
      header: t('planning.expenses'),
      render: (p) => <span className="font-medium text-[var(--color-danger)]">${p.projected_expenses.toFixed(2)}</span>,
    },
    {
      key: 'projected_debts',
      header: t('planning.debts'),
      render: (p) => <span className="font-medium text-[var(--color-warning)]">${p.projected_debts.toFixed(2)}</span>,
    },
    {
      key: 'projected_balance',
      header: t('planning.balance'),
      render: (p) => {
        const color = p.projected_balance >= 0 ? 'var(--color-success)' : 'var(--color-danger)';
        return <span className="font-bold" style={{ color }}>${p.projected_balance.toFixed(2)}</span>;
      },
    },
  ];

  const load = async () => {
    setLoading(true);
    try {
      const data = await accountingApi.getProjections(months);
      setProjections(data);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, [months]);

  return (
    <PageContainer>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('planning.title')}</h2>
          <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('planning.description')}</p>
        </div>
        <div className="flex items-center gap-3">
          <label className="text-sm font-medium text-[var(--color-text-secondary)]">{t('planning.months')}:</label>
          <Input
            type="number"
            value={String(months)}
            onChange={e => setMonths(Math.max(1, parseInt(e.target.value) || 1))}
            className="!w-20"
          />
          <Button onClick={load}>{t('planning.update')}</Button>
        </div>
      </div>

      {loading ? (
        <Loading text={t('planning.calculatingProjections')} />
      ) : (
        <>
          <Card title={t('planning.monthlyProjections')} className="mb-8">
            <Table columns={columns} data={projections} variant="striped" />
          </Card>

          {projections.length > 0 && (
            <Card title={t('planning.summary')}>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                <div>
                  <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.avgMonthlyIncome')}</p>
                  <p className="text-xl font-bold text-[var(--color-success)]">
                    ${(projections.reduce((s, p) => s + p.projected_income, 0) / projections.length).toFixed(2)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.avgMonthlyExpenses')}</p>
                  <p className="text-xl font-bold text-[var(--color-danger)]">
                    ${(projections.reduce((s, p) => s + p.projected_expenses, 0) / projections.length).toFixed(2)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.totalProjectedBalance')}</p>
                  <p className="text-xl font-bold text-[var(--color-text-primary)]">
                    ${projections.reduce((s, p) => s + p.projected_balance, 0).toFixed(2)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.monthsProjected')}</p>
                  <p className="text-xl font-bold text-[var(--color-text-primary)]">{projections.length}</p>
                </div>
              </div>
            </Card>
          )}
        </>
      )}
    </PageContainer>
  );
}
