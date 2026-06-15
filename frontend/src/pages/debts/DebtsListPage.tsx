import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
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
import { debtsApi, type Debt } from '../../services/debts';
import { formatCurrency } from '../../utils/format';

export function DebtsListPage() {
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [debts, setDebts] = useState<Debt[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({
    name: '', total_amount: '', interest_rate: '0', total_installments: '1', start_date: '', due_day: '',
  });

  const [showEditModal, setShowEditModal] = useState(false);
  const [editingDebt, setEditingDebt] = useState<Debt | null>(null);
  const [editForm, setEditForm] = useState({
    name: '', total_amount: '', interest_rate: '0', total_installments: '1', start_date: '', status: 'active', due_day: '',
  });

  const statusBadge = (status: string) => {
    const map: Record<string, { variant: 'info' | 'success' | 'warning' | 'danger'; label: string }> = {
      active: { variant: 'warning', label: t('debts.active') },
      paid: { variant: 'success', label: t('debts.paid') },
      paused: { variant: 'info', label: t('debts.paused') },
      settled: { variant: 'info', label: t('debts.settled') },
    };
    const m = map[status] || { variant: 'info' as const, label: status };
    return <Badge variant={m.variant}>{m.label}</Badge>;
  };

  const columns: Column<Debt>[] = [
    { key: 'name', header: t('debts.name') },
    {
      key: 'total_amount',
      header: t('debts.total'),
      render: (d) => <span className="font-medium">${formatCurrency(d.total_amount)}</span>,
    },
    {
      key: 'remaining_amount',
      header: t('debts.remaining'),
      render: (d) => <span className="font-medium text-[var(--color-danger)]">${formatCurrency(d.remaining_amount)}</span>,
    },
    {
      key: 'current_installment',
      header: t('debts.progress'),
      render: (d) => `${d.current_installment}/${d.total_installments}`,
    },
    {
      key: 'status',
      header: t('debts.status'),
      render: (d) => statusBadge(d.status),
    },
    {
      key: 'actions',
      header: t('common.actions'),
      render: (d) => (
        <div className="flex gap-2">
          <Button size="sm" variant="secondary" onClick={(e) => { e.stopPropagation(); openEdit(d); }}>
            {t('common.edit')}
          </Button>
          <Button size="sm" variant="danger" onClick={(e) => { e.stopPropagation(); handleDelete(d); }}>
            {t('common.delete')}
          </Button>
        </div>
      ),
    },
  ];

  const load = async () => {
    setLoading(true);
    try {
      const data = await debtsApi.list();
      setDebts(data);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    await debtsApi.create({
      name: form.name,
      total_amount: parseFloat(form.total_amount),
      interest_rate: parseFloat(form.interest_rate),
      total_installments: parseInt(form.total_installments),
      start_date: form.start_date,
      due_day: form.due_day ? parseInt(form.due_day) : undefined,
    });
    setShowModal(false);
    setForm({ name: '', total_amount: '', interest_rate: '0', total_installments: '1', start_date: '', due_day: '' });
    load();
  };

  const openEdit = (debt: Debt) => {
    setEditingDebt(debt);
    setEditForm({
      name: debt.name,
      total_amount: debt.total_amount.toString(),
      interest_rate: debt.interest_rate.toString(),
      total_installments: debt.total_installments.toString(),
      start_date: debt.start_date,
      status: debt.status,
      due_day: debt.due_day?.toString() || '',
    });
    setShowEditModal(true);
  };

  const handleEdit = async () => {
    if (!editingDebt) return;
    try {
      await debtsApi.update(editingDebt.id, {
        name: editForm.name,
        total_amount: parseFloat(editForm.total_amount),
        interest_rate: parseFloat(editForm.interest_rate),
        total_installments: parseInt(editForm.total_installments),
        start_date: editForm.start_date,
        status: editForm.status as any,
        due_day: editForm.due_day ? parseInt(editForm.due_day) : undefined,
      });
      setShowEditModal(false);
      setEditingDebt(null);
      toast(t('debts.debtUpdated'), 'success');
      load();
    } catch {
      toast(t('debts.debtUpdateError'), 'error');
    }
  };

  const handleDelete = async (debt: Debt) => {
    if (!window.confirm(t('debts.deleteConfirm'))) return;
    try {
      await debtsApi.delete(debt.id);
      toast(t('debts.debtDeleted'), 'success');
      load();
    } catch {
      toast(t('debts.debtDeleteError'), 'error');
    }
  };

  return (
    <PageContainer>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('debts.title')}</h2>
          <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('debts.description')}</p>
        </div>
        <Button onClick={() => setShowModal(true)}>{t('debts.addDebt')}</Button>
      </div>

      {loading ? (
        <Loading text={t('debts.loadingDebts')} />
      ) : debts.length === 0 ? (
        <EmptyState
          title={t('debts.noDebts')}
          description={t('debts.noDebtsDescription')}
          action={<Button onClick={() => setShowModal(true)}>{t('debts.addDebt')}</Button>}
        />
      ) : (
        <Table columns={columns} data={debts} variant="striped" onRowClick={(d) => navigate(`/debts/${d.id}`)} />
      )}

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={t('debts.addDebt')}>
        <div className="space-y-5">
          <Input label={t('debts.name')} value={form.name} onChange={e => setForm(p => ({ ...p, name: e.target.value }))} />
          <Input label={t('debts.totalAmount')} type="number" value={form.total_amount} onChange={e => setForm(p => ({ ...p, total_amount: e.target.value }))} />
          <Input label={t('debts.interestRate')} type="number" value={form.interest_rate} onChange={e => setForm(p => ({ ...p, interest_rate: e.target.value }))} />
          <Input label={t('debts.totalInstallments')} type="number" value={form.total_installments} onChange={e => setForm(p => ({ ...p, total_installments: e.target.value }))} />
          <Input label={t('debts.startDate')} type="date" value={form.start_date} onChange={e => setForm(p => ({ ...p, start_date: e.target.value }))} />
          <Input label={t('debts.dueDay')} type="number" value={form.due_day} onChange={e => setForm(p => ({ ...p, due_day: e.target.value }))} />
          <div className="pt-2">
            <Button className="w-full" onClick={handleCreate}>{t('debts.createDebt')}</Button>
          </div>
        </div>
      </Modal>

      <Modal isOpen={showEditModal} onClose={() => { setShowEditModal(false); setEditingDebt(null); }} title={t('debts.editDebt')}>
        <div className="space-y-5">
          <Input label={t('debts.name')} value={editForm.name} onChange={e => setEditForm(p => ({ ...p, name: e.target.value }))} />
          <Input label={t('debts.totalAmount')} type="number" value={editForm.total_amount} onChange={e => setEditForm(p => ({ ...p, total_amount: e.target.value }))} />
          <Input label={t('debts.interestRate')} type="number" value={editForm.interest_rate} onChange={e => setEditForm(p => ({ ...p, interest_rate: e.target.value }))} />
          <Input label={t('debts.totalInstallments')} type="number" value={editForm.total_installments} onChange={e => setEditForm(p => ({ ...p, total_installments: e.target.value }))} />
          <Input label={t('debts.startDate')} type="date" value={editForm.start_date} onChange={e => setEditForm(p => ({ ...p, start_date: e.target.value }))} />
          <Input label={t('debts.dueDay')} type="number" value={editForm.due_day} onChange={e => setEditForm(p => ({ ...p, due_day: e.target.value }))} />
          <Select
            label={t('debts.status')}
            options={[
              { value: 'active', label: t('debts.active') },
              { value: 'paid', label: t('debts.paid') },
              { value: 'settled', label: t('debts.settled') },
            ]}
            value={editForm.status}
            onChange={e => setEditForm(p => ({ ...p, status: e.target.value }))}
          />
          <div className="pt-2">
            <Button className="w-full" onClick={handleEdit}>{t('debts.saveChanges')}</Button>
          </div>
        </div>
      </Modal>
    </PageContainer>
  );
}
