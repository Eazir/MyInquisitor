import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Badge } from '../../components/ui/Badge';
import { Loading } from '../../components/ui/Loading';
import { Table, type Column } from '../../components/ui/Table';
import { useLanguage } from '../../contexts/LanguageContext';
import { debtsApi, type Debt, type DebtMonthlyStatus } from '../../services/debts';

export function DebtDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [debt, setDebt] = useState<Debt | null>(null);
  const [monthly, setMonthly] = useState<DebtMonthlyStatus[]>([]);
  const [loading, setLoading] = useState(true);

  const statusBadge = (status: string) => {
    const map: Record<string, { variant: 'info' | 'success' | 'warning'; label: string }> = {
      active: { variant: 'warning', label: t('debts.active') },
      paid: { variant: 'success', label: t('debts.paid') },
      settled: { variant: 'info', label: t('debts.settled') },
    };
    const m = map[status] || { variant: 'info' as const, label: status };
    return <Badge variant={m.variant}>{m.label}</Badge>;
  };

  const monthlyColumns: Column<DebtMonthlyStatus>[] = [
    { key: 'month', header: t('debts.month'), render: (s) => new Date(s.month + 'T00:00:00').toLocaleDateString('en-US', { year: 'numeric', month: 'long' }) },
    {
      key: 'amount_due',
      header: t('debts.due'),
      render: (s) => `$${s.amount_due.toFixed(2)}`,
    },
    {
      key: 'amount_paid',
      header: t('debts.paid'),
      render: (s) => `$${s.amount_paid.toFixed(2)}`,
    },
    {
      key: 'paid',
      header: t('debts.status'),
      render: (s) => s.paid
        ? <Badge variant="success">{t('debts.paid')}</Badge>
        : <Badge variant="warning">{t('debts.pending')}</Badge>,
    },
  ];

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    Promise.all([
      debtsApi.getByID(id),
      debtsApi.getMonthlyStatus(id),
    ]).then(([d, m]) => {
      setDebt(d);
      setMonthly(m);
    }).finally(() => setLoading(false));
  }, [id]);

  if (loading) return <Loading text={t('debts.loadingDebts')} />;
  if (!debt) return <div className="p-8 text-center text-[var(--color-text-muted)]">{t('debts.debtNotFound')}</div>;

  return (
    <PageContainer>
      <div className="mb-6">
        <Button variant="ghost" onClick={() => navigate('/debts')}>← {t('debts.backToDebts')}</Button>
      </div>

      <Card title={debt.name} subtitle={`${debt.current_installment}/${debt.total_installments} ${t('debts.installments')}`}>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p className="text-sm text-[var(--color-text-secondary)]">{t('debts.totalAmount')}</p>
            <p className="text-lg font-semibold">${debt.total_amount.toFixed(2)}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--color-text-secondary)]">{t('debts.remaining')}</p>
            <p className="text-lg font-semibold">${debt.remaining_amount.toFixed(2)}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--color-text-secondary)]">{t('debts.interestRate')}</p>
            <p className="text-lg font-semibold">{debt.interest_rate}%</p>
          </div>
          <div>
            <p className="text-sm text-[var(--color-text-secondary)]">{t('debts.status')}</p>
            <div className="mt-1">{statusBadge(debt.status)}</div>
          </div>
        </div>
      </Card>

      <Card title={t('debts.monthlyStatus')} className="mt-8">
        <Table columns={monthlyColumns} data={monthly} variant="striped" />
      </Card>
    </PageContainer>
  );
}
