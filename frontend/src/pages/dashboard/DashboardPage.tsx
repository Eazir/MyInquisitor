import { useState, useEffect } from 'react';
import { Card } from '../../components/ui/Card';
import { Badge } from '../../components/ui/Badge';
import { Loading } from '../../components/ui/Loading';
import { PageContainer } from '../../components/layout/PageContainer';
import { useLanguage } from '../../contexts/LanguageContext';
import { accountingApi } from '../../services/accounting';
import { debtsApi, type Debt, type DebtMonthlyStatus } from '../../services/debts';
import { expensesApi, type Expense } from '../../services/expenses';
import { formatCurrency } from '../../utils/format';

interface PaymentItem {
  id: string;
  name: string;
  type: 'debt' | 'expense';
  amount: number;
  paid: boolean;
  detail: string;
}

export function DashboardPage() {
  const { t } = useLanguage();
  const [balance, setBalance] = useState<{ total_income: number; total_expenses: number; net_balance: number } | null>(null);
  const [debts, setDebts] = useState<Debt[]>([]);
  const [paymentItems, setPaymentItems] = useState<PaymentItem[]>([]);
  const [loading, setLoading] = useState(true);

  const now = new Date();
  const currentYear = now.getFullYear();
  const currentMonth = now.getMonth() + 1;
  const currentMonthStr = `${currentYear}-${String(currentMonth).padStart(2, '0')}`;

  const load = async () => {
    setLoading(true);
    try {
      const [bal, debtData, expData] = await Promise.all([
        accountingApi.getMonthlyBalance(currentYear, currentMonth),
        debtsApi.list(),
        expensesApi.list(),
      ]);
      setBalance(bal);
      setDebts(debtData);

      const activeDebts = debtData.filter((d: Debt) => d.status === 'active');
      const activeExpenses = expData.filter((e: Expense) => {
        if (e.status !== 'active') return false;
        if (e.frequency === 'yearly') {
          const billingMonth = e.billing_month || new Date(e.start_date + 'T00:00:00').getMonth() + 1;
          return billingMonth === currentMonth;
        }
        if (e.frequency === 'once') {
          const expMonth = new Date(e.start_date + 'T00:00:00').getMonth() + 1;
          const expYear = new Date(e.start_date + 'T00:00:00').getFullYear();
          return expMonth === currentMonth && expYear === currentYear;
        }
        return true;
      });

      const debtMonthMap = new Map<string, { paid: boolean; amount_due: number }>();
      await Promise.all(
        activeDebts.map(async (d: Debt) => {
          try {
            const statuses: DebtMonthlyStatus[] = await debtsApi.getMonthlyStatus(d.id);
            const current = statuses.find(s => s.month.startsWith(currentMonthStr));
            debtMonthMap.set(d.id, { paid: current?.paid ?? false, amount_due: current?.amount_due ?? d.remaining_amount });
          } catch {
            debtMonthMap.set(d.id, { paid: false, amount_due: d.remaining_amount });
          }
        })
      );

      const expenseStatusResults = await Promise.all(
        activeExpenses.map((e: Expense) =>
          expensesApi.getMonthlyStatus(e.id, currentYear, currentMonth).catch(() => null)
        )
      );

      const items: PaymentItem[] = [
        ...activeDebts.map((d: Debt) => {
          const m = debtMonthMap.get(d.id);
          return {
            id: d.id,
            name: d.name,
            type: 'debt' as const,
            amount: m?.amount_due ?? d.remaining_amount,
            paid: m?.paid ?? false,
            detail: `${d.current_installment}/${d.total_installments} ${t('debts.installments')}`,
          };
        }),
        ...activeExpenses.map((e: Expense, i: number) => ({
          id: e.id,
          name: e.name,
          type: 'expense' as const,
          amount: e.amount,
          paid: expenseStatusResults[i]?.paid ?? false,
          detail: e.due_day ? t('expenses.dayPlaceholder', { day: e.due_day }) : '',
        })),
      ];

      setPaymentItems(items.filter(item => !item.paid));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const reservedAmount = debts.filter(d => d.status === 'active').reduce((sum, d) => sum + d.remaining_amount, 0);
  const totalManaged = balance?.total_income ?? 0;

  if (loading) return <Loading text={t('common.loading')} />;

  return (
    <PageContainer>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('dashboard.title')}</h2>
        <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('dashboard.description')}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6 lg:gap-8 mb-8">
        <Card variant="stats">
          <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('dashboard.totalManaged')}</p>
          <p className="text-3xl font-bold text-[var(--color-text-primary)]">${formatCurrency(totalManaged)}</p>
        </Card>
        <Card variant="stats">
          <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('dashboard.available')}</p>
          <p className="text-3xl font-bold text-[var(--color-success)]">${formatCurrency(balance?.net_balance ?? 0)}</p>
        </Card>
        <Card variant="stats">
          <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('dashboard.reserved')}</p>
          <p className="text-3xl font-bold text-[var(--color-danger)]">${formatCurrency(reservedAmount)}</p>
        </Card>
      </div>

      <Card title={t('dashboard.upcomingPayments')} className="mb-8">
        {paymentItems.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <span className="text-4xl mb-4">{'📅'}</span>
            <p className="text-[var(--color-text-muted)] text-sm">{t('dashboard.noUpcomingPayments')}</p>
          </div>
        ) : (
          <div className="divide-y divide-[var(--color-border)]">
            {paymentItems.map(item => (
              <div key={`${item.type}-${item.id}`} className="flex items-center justify-between py-4 px-2 hover:bg-[var(--color-bg-hover)] transition-colors">
                <div className="flex items-center gap-3 min-w-0">
                  <span className="text-lg flex-shrink-0">{item.type === 'debt' ? '💳' : '📄'}</span>
                  <div className="min-w-0">
                    <p className="font-medium text-[var(--color-text-primary)] truncate">{item.name}</p>
                    <p className="text-xs text-[var(--color-text-secondary)]">
                      <Badge variant={item.type === 'debt' ? 'warning' : 'info'}>{item.type === 'debt' ? t('monthlyPayments.debt') : t('monthlyPayments.expense')}</Badge>
                      {item.detail && <span className="ml-2">{item.detail}</span>}
                    </p>
                  </div>
                </div>
                <div className="text-right flex-shrink-0">
                  <p className={`font-semibold ${item.type === 'debt' ? 'text-[var(--color-danger)]' : 'text-[var(--color-accent)]'}`}>
                    ${formatCurrency(item.amount)}
                  </p>
                  {item.paid
                    ? <Badge variant="success">{t('monthlyPayments.paidLabel')}</Badge>
                    : <Badge variant="warning">{t('debts.pending')}</Badge>
                  }
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

      <Card title={t('dashboard.balanceTrend')}>
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <span className="text-4xl mb-4">{'📈'}</span>
          <p className="text-[var(--color-text-muted)] text-sm">{t('dashboard.chartComingSoon')}</p>
        </div>
      </Card>
    </PageContainer>
  );
}
