import { useState, useEffect } from 'react';
import { PageContainer } from '../../components/layout/PageContainer';
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
import { expensesApi, type Expense } from '../../services/expenses';
import { formatCurrency } from '../../utils/format';

export function ExpensesPage() {
  const { t } = useLanguage();
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({
    name: '', amount: '', frequency: 'monthly', due_day: '', billing_month: '', start_date: '',
  });

  const [showEditModal, setShowEditModal] = useState(false);
  const [editingExpense, setEditingExpense] = useState<Expense | null>(null);
  const [editForm, setEditForm] = useState({
    name: '', amount: '', frequency: 'monthly', due_day: '', billing_month: '', status: 'active',
  });

  const frequencyLabels: Record<string, string> = {
    monthly: t('expenses.monthly'),
    yearly: t('expenses.yearly'),
    weekly: t('expenses.weekly'),
    biweekly: t('expenses.biweekly'),
    once: t('expenses.once'),
  };

  const monthOptions = Array.from({ length: 12 }, (_, i) => ({
    value: String(i + 1),
    label: new Date(0, i).toLocaleDateString('default', { month: 'long' }),
  }));

  const statusBadge = (status: string) => {
    const map: Record<string, { variant: 'info' | 'success'; label: string }> = {
      active: { variant: 'info', label: t('expenses.active') },
      cancelled: { variant: 'success', label: t('expenses.cancelled') },
    };
    const m = map[status] || { variant: 'info' as const, label: status };
    return <Badge variant={m.variant}>{m.label}</Badge>;
  };

  const columns: Column<Expense>[] = [
    { key: 'name', header: t('expenses.name') },
    {
      key: 'amount',
      header: t('expenses.amount'),
      render: (e) => <span className="font-medium">${formatCurrency(e.amount)}</span>,
    },
    {
      key: 'frequency',
      header: t('expenses.frequency'),
      render: (e) => frequencyLabels[e.frequency] || e.frequency,
    },
    {
      key: 'due_day',
      header: t('expenses.dueDay'),
      render: (e) => e.due_day ? t('expenses.dayPlaceholder', { day: e.due_day }) : '-',
    },
    {
      key: 'status',
      header: t('expenses.status'),
      render: (e) => statusBadge(e.status),
    },
    {
      key: 'actions',
      header: t('common.actions'),
      render: (e) => (
        <div className="flex gap-2">
          <Button size="sm" variant="secondary" onClick={() => openEdit(e)}>
            {t('common.edit')}
          </Button>
          <Button size="sm" variant="danger" onClick={() => handleDelete(e)}>
            {t('common.delete')}
          </Button>
        </div>
      ),
    },
  ];

  const load = async () => {
    setLoading(true);
    try {
      const data = await expensesApi.list();
      setExpenses(data);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    await expensesApi.create({
      name: form.name,
      amount: parseFloat(form.amount),
      frequency: form.frequency,
      due_day: form.due_day ? parseInt(form.due_day) : undefined,
      billing_month: form.billing_month ? parseInt(form.billing_month) : undefined,
      start_date: form.start_date,
    });
    setShowModal(false);
    setForm({ name: '', amount: '', frequency: 'monthly', due_day: '', billing_month: '', start_date: '' });
    load();
  };

  const openEdit = (expense: Expense) => {
    setEditingExpense(expense);
    setEditForm({
      name: expense.name,
      amount: expense.amount.toString(),
      frequency: expense.frequency,
      due_day: expense.due_day?.toString() || '',
      billing_month: expense.billing_month?.toString() || '',
      status: expense.status,
    });
    setShowEditModal(true);
  };

  const handleEdit = async () => {
    if (!editingExpense) return;
    try {
      await expensesApi.update(editingExpense.id, {
        name: editForm.name,
        amount: parseFloat(editForm.amount),
        frequency: editForm.frequency,
        due_day: editForm.due_day ? parseInt(editForm.due_day) : undefined,
        billing_month: editForm.billing_month ? parseInt(editForm.billing_month) : undefined,
        status: editForm.status,
      });
      setShowEditModal(false);
      setEditingExpense(null);
      toast(t('expenses.expenseUpdated'), 'success');
      load();
    } catch {
      toast(t('expenses.expenseUpdateError'), 'error');
    }
  };

  const handleDelete = async (expense: Expense) => {
    if (!window.confirm(t('expenses.deleteExpenseConfirm'))) return;
    try {
      await expensesApi.delete(expense.id);
      toast(t('expenses.expenseDeleted'), 'success');
      load();
    } catch {
      toast(t('expenses.expenseDeleteError'), 'error');
    }
  };

  return (
    <PageContainer>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('expenses.title')}</h2>
          <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('expenses.description')}</p>
        </div>
        <Button onClick={() => setShowModal(true)}>{t('expenses.addExpense')}</Button>
      </div>

      {loading ? (
        <Loading text={t('expenses.loadingExpenses')} />
      ) : expenses.length === 0 ? (
        <EmptyState
          title={t('expenses.noExpenses')}
          description={t('expenses.noExpensesDescription')}
          action={<Button onClick={() => setShowModal(true)}>{t('expenses.addExpense')}</Button>}
        />
      ) : (
        <Table columns={columns} data={expenses} variant="striped" />
      )}

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={t('expenses.addExpense')}>
        <div className="space-y-5">
          <Input label={t('expenses.name')} value={form.name} onChange={e => setForm(p => ({ ...p, name: e.target.value }))} />
          <Input label={t('expenses.amount')} type="number" value={form.amount} onChange={e => setForm(p => ({ ...p, amount: e.target.value }))} />
          <Select
            label={t('expenses.frequency')}
            options={[
              { value: 'monthly', label: t('expenses.monthly') },
              { value: 'yearly', label: t('expenses.yearly') },
              { value: 'weekly', label: t('expenses.weekly') },
              { value: 'biweekly', label: t('expenses.biweekly') },
              { value: 'once', label: t('expenses.once') },
            ]}
            value={form.frequency}
            onChange={e => setForm(p => ({ ...p, frequency: e.target.value }))}
          />
          {form.frequency === 'yearly' && (
            <Select
              label={t('expenses.billingMonth')}
              options={monthOptions}
              value={form.billing_month}
              onChange={e => setForm(p => ({ ...p, billing_month: e.target.value }))}
            />
          )}
          <Input label={t('expenses.dueDay')} type="number" value={form.due_day} onChange={e => setForm(p => ({ ...p, due_day: e.target.value }))} />
          <Input label={t('expenses.startDate')} type="date" value={form.start_date} onChange={e => setForm(p => ({ ...p, start_date: e.target.value }))} />
          <div className="pt-2">
            <Button className="w-full" onClick={handleCreate}>{t('expenses.createExpense')}</Button>
          </div>
        </div>
      </Modal>

      <Modal isOpen={showEditModal} onClose={() => { setShowEditModal(false); setEditingExpense(null); }} title={t('expenses.editExpense')}>
        <div className="space-y-5">
          <Input label={t('expenses.name')} value={editForm.name} onChange={e => setEditForm(p => ({ ...p, name: e.target.value }))} />
          <Input label={t('expenses.amount')} type="number" value={editForm.amount} onChange={e => setEditForm(p => ({ ...p, amount: e.target.value }))} />
          <Select
            label={t('expenses.frequency')}
            options={[
              { value: 'monthly', label: t('expenses.monthly') },
              { value: 'yearly', label: t('expenses.yearly') },
              { value: 'weekly', label: t('expenses.weekly') },
              { value: 'biweekly', label: t('expenses.biweekly') },
              { value: 'once', label: t('expenses.once') },
            ]}
            value={editForm.frequency}
            onChange={e => setEditForm(p => ({ ...p, frequency: e.target.value }))}
          />
          {editForm.frequency === 'yearly' && (
            <Select
              label={t('expenses.billingMonth')}
              options={monthOptions}
              value={editForm.billing_month}
              onChange={e => setEditForm(p => ({ ...p, billing_month: e.target.value }))}
            />
          )}
          <Input label={t('expenses.dueDay')} type="number" value={editForm.due_day} onChange={e => setEditForm(p => ({ ...p, due_day: e.target.value }))} />
          <Select
            label={t('expenses.status')}
            options={[
              { value: 'active', label: t('expenses.active') },
              { value: 'cancelled', label: t('expenses.cancelled') },
            ]}
            value={editForm.status}
            onChange={e => setEditForm(p => ({ ...p, status: e.target.value }))}
          />
          <div className="pt-2">
            <Button className="w-full" onClick={handleEdit}>{t('common.save')}</Button>
          </div>
        </div>
      </Modal>
    </PageContainer>
  );
}
