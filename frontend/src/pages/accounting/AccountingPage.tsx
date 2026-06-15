import { useState, useEffect } from 'react';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Table, type Column } from '../../components/ui/Table';
import { Badge } from '../../components/ui/Badge';
import { Modal } from '../../components/ui/Modal';
import { Input } from '../../components/ui/Input';
import { Select } from '../../components/ui/Select';
import { Loading } from '../../components/ui/Loading';
import { EmptyState } from '../../components/ui/EmptyState';
import { useLanguage } from '../../contexts/LanguageContext';
import { toast } from '../../components/ui/Toast';
import { accountingApi, type Transaction, type MonthlyBalance } from '../../services/accounting';
import { formatCurrency } from '../../utils/format';

export function AccountingPage() {
  const { t } = useLanguage();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [balance, setBalance] = useState<MonthlyBalance | null>(null);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const today = new Date().toISOString().slice(0, 10);
  const [form, setForm] = useState({
    type: 'expense', amount: '', description: '', reference_date: today,
  });

  const now = new Date();
  const currentMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;

  const typeBadge = (type: string) => {
    const map: Record<string, { variant: 'success' | 'warning' | 'info'; label: string }> = {
      income: { variant: 'success', label: t('accounting.incomeType') },
      expense: { variant: 'warning', label: t('accounting.expenseType') },
      transfer: { variant: 'info', label: t('accounting.transferType') },
    };
    const m = map[type] || { variant: 'info' as const, label: type };
    return <Badge variant={m.variant}>{m.label}</Badge>;
  };

  const columns: Column<Transaction>[] = [
    {
      key: 'type',
      header: t('accounting.type'),
      render: (t) => typeBadge(t.type),
    },
    {
      key: 'amount',
      header: t('accounting.amount'),
      render: (t) => <span className="font-medium">${formatCurrency(t.amount)}</span>,
    },
    { key: 'description', header: t('accounting.description'), render: (t) => t.description || '-' },
    {
      key: 'reference_date',
      header: t('accounting.date'),
      render: (t) => new Date(t.reference_date + 'T00:00:00').toLocaleDateString(),
    },
  ];

  const load = async () => {
    setLoading(true);
    try {
      const [year, month] = currentMonth.split('-').map(Number);
      const [txData, bal] = await Promise.all([
        accountingApi.listTransactions({ year, month }),
        accountingApi.getMonthlyBalance(year, month),
      ]);
      setTransactions(txData.data || []);
      setBalance(bal);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    const amount = parseFloat(form.amount);
    if (!amount || amount <= 0) {
      toast(t('accounting.invalidAmount'), 'error');
      return;
    }
    if (!form.reference_date) {
      toast(t('accounting.requiredDate'), 'error');
      return;
    }
    try {
      await accountingApi.recordTransaction({
        type: form.type,
        amount,
        description: form.description || undefined,
        reference_date: form.reference_date,
      });
      setShowModal(false);
      setForm({ type: 'expense', amount: '', description: '', reference_date: today });
      toast(t('accounting.transactionCreated'), 'success');
      load();
    } catch {
      toast(t('accounting.transactionError'), 'error');
    }
  };

  return (
    <PageContainer>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('accounting.title')}</h2>
          <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('accounting.subtitle')}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6 lg:gap-8 mb-8">
        <Card variant="stats" title={t('accounting.income')} subtitle={currentMonth}>
          <p className="text-2xl font-bold text-[var(--color-success)]">
            ${formatCurrency(balance?.total_income ?? 0)}
          </p>
        </Card>
        <Card variant="stats" title={t('accounting.expenses')} subtitle={currentMonth}>
          <p className="text-2xl font-bold text-[var(--color-danger)]">
            ${formatCurrency(balance?.total_expenses ?? 0)}
          </p>
        </Card>
        <Card variant="stats" title={t('accounting.netBalance')} subtitle={currentMonth}>
          <p className={`text-2xl font-bold ${(balance?.net_balance || 0) >= 0 ? 'text-[var(--color-success)]' : 'text-[var(--color-danger)]'}`}>
            ${formatCurrency(balance?.net_balance ?? 0)}
          </p>
        </Card>
      </div>

      <div className="flex items-center justify-between mb-8">
        <h3 className="text-lg font-semibold text-[var(--color-text-primary)]">{t('accounting.transactions')}</h3>
        <Button onClick={() => setShowModal(true)}>{t('accounting.addTransaction')}</Button>
      </div>

      {loading ? (
        <Loading text={t('accounting.loadingTransactions')} />
      ) : transactions.length === 0 ? (
        <EmptyState title={t('accounting.noTransactions')} description={t('accounting.noTransactionsDescription')} />
      ) : (
        <Table columns={columns} data={transactions} variant="striped" />
      )}

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={t('accounting.addTransaction')}>
        <div className="space-y-5">
          <Select
            label={t('accounting.type')}
            options={[
              { value: 'income', label: t('accounting.incomeType') },
              { value: 'expense', label: t('accounting.expenseType') },
              { value: 'transfer', label: t('accounting.transferType') },
            ]}
            value={form.type}
            onChange={e => setForm(p => ({ ...p, type: e.target.value }))}
          />
          <Input label={t('accounting.amount')} type="number" value={form.amount} onChange={e => setForm(p => ({ ...p, amount: e.target.value }))} />
          <Input label={t('accounting.description')} value={form.description} onChange={e => setForm(p => ({ ...p, description: e.target.value }))} />
          <Input label={t('accounting.date')} type="date" value={form.reference_date} onChange={e => setForm(p => ({ ...p, reference_date: e.target.value }))} />
          <div className="pt-2">
            <Button className="w-full" onClick={handleCreate}>{t('accounting.saveTransaction')}</Button>
          </div>
        </div>
      </Modal>
    </PageContainer>
  );
}
