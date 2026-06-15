import { useState, useEffect, useCallback } from 'react';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Badge } from '../../components/ui/Badge';
import { Loading } from '../../components/ui/Loading';
import { EmptyState } from '../../components/ui/EmptyState';
import { useLanguage } from '../../contexts/LanguageContext';
import { toast } from '../../components/ui/Toast';
import { debtsApi, type DebtMonthlyStatus } from '../../services/debts';
import { expensesApi, type Expense } from '../../services/expenses';
import { accountingApi } from '../../services/accounting';
import { formatCurrency } from '../../utils/format';

interface ObligationItem {
  id: string;
  name: string;
  type: 'debt' | 'expense';
  amount: number;
  paid: boolean;
  refId: string;
  month: string;
  installmentNum?: number;
  dueDay?: number | null;
}

export function MonthlyPaymentsPage() {
  const { t } = useLanguage();
  const [items, setItems] = useState<ObligationItem[]>([]);
  const [balance, setBalance] = useState<{ total_income: number; total_expenses: number; net_balance: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [payingId, setPayingId] = useState<string | null>(null);

  const now = new Date();
  const currentYear = now.getFullYear();
  const currentMonth = now.getMonth() + 1;
  const currentMonthStr = `${currentYear}-${String(currentMonth).padStart(2, '0')}`;

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [debts, expenses, bal] = await Promise.all([
        debtsApi.list(),
        expensesApi.list(),
        accountingApi.getMonthlyBalance(currentYear, currentMonth),
      ]);
      setBalance(bal);

      const result: ObligationItem[] = [];

      for (const debt of debts) {
        let monthlyStatuses: DebtMonthlyStatus[];
        try {
          monthlyStatuses = await debtsApi.getMonthlyStatus(debt.id);
        } catch {
          continue;
        }

        const current = monthlyStatuses.find(s => s.month.startsWith(currentMonthStr));
        if (current) {
          result.push({
            id: current.id,
            name: debt.name,
            type: 'debt',
            amount: current.amount_due,
            paid: current.paid,
            refId: debt.id,
            month: current.month,
            installmentNum: current.installment_num,
            dueDay: debt.due_day,
          });
        }
      }

      const filteredExpenses = expenses.filter((exp: Expense) => {
        if (exp.frequency === 'yearly') {
          const billingMonth = exp.billing_month || new Date(exp.start_date + 'T00:00:00').getMonth() + 1;
          return billingMonth === currentMonth;
        }
        if (exp.frequency === 'once') {
          const expMonth = new Date(exp.start_date + 'T00:00:00').getMonth() + 1;
          const expYear = new Date(exp.start_date + 'T00:00:00').getFullYear();
          return expMonth === currentMonth && expYear === currentYear;
        }
        return true;
      });

      const expenseStatusResults = await Promise.all(
        filteredExpenses.map((exp: Expense) =>
          expensesApi.getMonthlyStatus(exp.id, currentYear, currentMonth).catch(() => null)
        )
      );

      for (let i = 0; i < filteredExpenses.length; i++) {
        const expense = filteredExpenses[i];
        const status = expenseStatusResults[i];
        result.push({
          id: expense.id,
          name: expense.name,
          type: 'expense',
          amount: expense.amount,
          paid: status?.paid ?? false,
          refId: expense.id,
          month: currentMonthStr,
          dueDay: expense.due_day,
        });
      }

      setItems(result);
    } finally {
      setLoading(false);
    }
  }, [currentMonthStr]);

  useEffect(() => { load(); }, [load]);

  const handleMarkAsPaid = async (item: ObligationItem) => {
    setPayingId(item.id);
    try {
      if (item.type === 'debt') {
        await debtsApi.markAsPaid(item.refId, currentYear, currentMonth, item.amount);
        setItems(prev => prev.map(i => i.id === item.id ? { ...i, paid: true } : i));
      } else {
        const updated = await expensesApi.togglePaid(item.refId, currentYear, currentMonth);
        setItems(prev => prev.map(i => i.id === item.id ? { ...i, paid: updated.paid } : i));
      }
      const bal = await accountingApi.getMonthlyBalance(currentYear, currentMonth);
      setBalance(bal);
      toast(t('monthlyPayments.paySuccess'), 'success');
    } catch {
      toast(t('monthlyPayments.payError'), 'error');
    } finally {
      setPayingId(null);
    }
  };

  const totalDue = items.reduce((sum, i) => sum + i.amount, 0);
  const totalPaid = items.filter(i => i.paid).reduce((sum, i) => sum + i.amount, 0);
  const totalPending = totalDue - totalPaid;

  if (loading) return <Loading text={t('common.loading')} />;

  return (
    <PageContainer>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('monthlyPayments.title')}</h2>
        <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('monthlyPayments.description')}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6 mb-8">
        <Card variant="stats">
          <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('dashboard.available')}</p>
          <p className={`text-3xl font-bold ${(balance?.net_balance ?? 0) >= 0 ? 'text-[var(--color-success)]' : 'text-[var(--color-danger)]'}`}>
            ${formatCurrency(balance?.net_balance ?? 0)}
          </p>
        </Card>
        <Card variant="stats">
          <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('monthlyPayments.totalDue')}</p>
          <p className="text-3xl font-bold text-[var(--color-text-primary)]">${formatCurrency(totalDue)}</p>
        </Card>
        <Card variant="stats">
          <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('monthlyPayments.totalPending')}</p>
          <p className="text-3xl font-bold text-[var(--color-danger)]">${formatCurrency(totalPending)}</p>
        </Card>
        <Card variant="stats">
          <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('monthlyPayments.totalPaid')}</p>
          <p className="text-3xl font-bold text-[var(--color-success)]">${formatCurrency(totalPaid)}</p>
        </Card>
      </div>

      {items.length === 0 ? (
        <EmptyState
          title={t('monthlyPayments.noObligations')}
          description=""
        />
      ) : (
        <Card>
          <div className="divide-y divide-[var(--color-border)]">
            <div className="grid grid-cols-12 gap-4 px-4 py-3 text-xs font-medium uppercase text-[var(--color-text-muted)]">
              <div className="col-span-1">{t('monthlyPayments.type')}</div>
              <div className="col-span-3">{t('monthlyPayments.name')}</div>
              <div className="col-span-2">{t('monthlyPayments.dueDate')}</div>
              <div className="col-span-2 text-right">{t('monthlyPayments.amount')}</div>
              <div className="col-span-2 text-center">{t('monthlyPayments.status')}</div>
              <div className="col-span-2 text-right">{t('common.actions')}</div>
            </div>
            {items.map(item => (
              <div key={item.id} className="grid grid-cols-12 gap-4 px-4 py-4 items-center">
                <div className="col-span-1">
                  <Badge variant={item.type === 'debt' ? 'warning' : 'info'}>
                    {item.type === 'debt' ? t('monthlyPayments.debt') : t('monthlyPayments.expense')}
                  </Badge>
                </div>
                <div className="col-span-3 font-medium text-[var(--color-text-primary)]">
                  {item.name}
                  {item.installmentNum && (
                    <span className="text-xs text-[var(--color-text-muted)] ml-2">#{item.installmentNum}</span>
                  )}
                </div>
                <div className="col-span-2 text-sm text-[var(--color-text-secondary)]">
                  {item.dueDay ? (
                    item.type === 'debt'
                      ? `${item.dueDay}-${item.month.slice(5, 7)}`
                      : `${item.dueDay}-${currentMonthStr.slice(5, 7)}`
                  ) : (
                    item.type === 'debt' ? item.month.slice(5, 7) : '—'
                  )}
                </div>
                <div className="col-span-2 text-right font-semibold">${formatCurrency(item.amount)}</div>
                <div className="col-span-2 text-center">
                  {item.paid
                    ? <Badge variant="success">{t('monthlyPayments.paidLabel')}</Badge>
                    : <Badge variant="warning">{t('monthlyPayments.pendingLabel')}</Badge>
                  }
                </div>
                <div className="col-span-2 text-right">
                  {!item.paid && (
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => handleMarkAsPaid(item)}
                      disabled={payingId === item.id}
                    >
                      {payingId === item.id ? '...' : t('monthlyPayments.markAsPaid')}
                    </Button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}
    </PageContainer>
  );
}
