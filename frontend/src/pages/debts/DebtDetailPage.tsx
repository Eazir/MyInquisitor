import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Badge } from '../../components/ui/Badge';
import { Loading } from '../../components/ui/Loading';
import { Table, type Column } from '../../components/ui/Table';
import { Modal } from '../../components/ui/Modal';
import { Input } from '../../components/ui/Input';
import { Select } from '../../components/ui/Select';
import { useLanguage } from '../../contexts/LanguageContext';
import { toast } from '../../components/ui/Toast';
import { debtsApi, type Debt, type DebtMonthlyStatus } from '../../services/debts';
import { formatCurrency } from '../../utils/format';

export function DebtDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [debt, setDebt] = useState<Debt | null>(null);
  const [monthly, setMonthly] = useState<DebtMonthlyStatus[]>([]);
  const [loading, setLoading] = useState(true);

  const [showEditModal, setShowEditModal] = useState(false);
  const [editForm, setEditForm] = useState({
    name: '', total_amount: '', interest_rate: '0', total_installments: '1', start_date: '', status: 'active', due_day: '',
  });

  const [showPayModal, setShowPayModal] = useState(false);
  const [selectedStatus, setSelectedStatus] = useState<DebtMonthlyStatus | null>(null);
  const [payAmount, setPayAmount] = useState('');
  const [paying, setPaying] = useState(false);

  const formatMonth = (monthStr: string) => {
    const d = new Date(monthStr + 'T00:00:00');
    return d.toLocaleDateString('es-ES', { month: 'long', year: 'numeric' });
  };

  const formatDateTime = (dateStr: string | null) => {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    return d.toLocaleDateString('es-ES', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' });
  };

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

  const monthlyColumns: Column<DebtMonthlyStatus>[] = [
    {
      key: 'installment_num',
      header: t('debts.installmentNum'),
      render: (s) => `#${s.installment_num}`,
    },
    {
      key: 'month',
      header: t('debts.month'),
      render: (s) => formatMonth(s.month),
    },
    {
      key: 'principal_amount',
      header: t('debts.principal'),
      render: (s) => `$${formatCurrency(s.principal_amount)}`,
    },
    {
      key: 'interest_amount',
      header: t('debts.interest'),
      render: (s) => `$${formatCurrency(s.interest_amount)}`,
    },
    {
      key: 'amount_due',
      header: t('debts.totalDue'),
      render: (s) => `$${formatCurrency(s.amount_due)}`,
    },
    {
      key: 'paid_at',
      header: t('debts.paidAt'),
      render: (s) => formatDateTime(s.paid_at),
    },
    {
      key: 'paid',
      header: t('debts.status'),
      render: (s) => {
        if (s.paid) return <Badge variant="success">{t('debts.paid')}</Badge>;
        const isOverdue = new Date(s.month + 'T00:00:00') < new Date(new Date().toISOString().slice(0, 7) + '-01');
        return isOverdue
          ? <Badge variant="danger">{t('debts.overdue')}</Badge>
          : <Badge variant="warning">{t('debts.pending')}</Badge>;
      },
    },
    {
      key: 'actions',
      header: t('common.actions'),
      render: (s) => !s.paid ? (
        <Button size="sm" variant="secondary" onClick={() => openPayModal(s)}>
          {t('debts.pay')}
        </Button>
      ) : null,
    },
  ];

  const load = () => {
    if (!id) return;
    setLoading(true);
    Promise.all([
      debtsApi.getByID(id),
      debtsApi.getMonthlyStatus(id),
    ]).then(([d, m]) => {
      setDebt(d);
      setEditForm({
        name: d.name,
        total_amount: d.total_amount.toString(),
        interest_rate: d.interest_rate.toString(),
        total_installments: d.total_installments.toString(),
        start_date: d.start_date,
        status: d.status,
        due_day: d.due_day?.toString() || '',
      });
      setMonthly(m.sort((a: DebtMonthlyStatus, b: DebtMonthlyStatus) => a.month.localeCompare(b.month)));
    }).finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, [id]);

  const openPayModal = (s: DebtMonthlyStatus) => {
    setSelectedStatus(s);
    setPayAmount(s.amount_due.toString());
    setShowPayModal(true);
  };

  const handlePay = async (amount: number) => {
    if (!debt || !selectedStatus) return;
    const [year, month] = selectedStatus.month.split('-');
    setPaying(true);
    try {
      await debtsApi.markAsPaid(debt.id, parseInt(year), parseInt(month), amount);
      toast(t('debts.paySuccess'), 'success');
      setShowPayModal(false);
      setSelectedStatus(null);
      load();
    } catch {
      toast(t('debts.payError'), 'error');
    } finally {
      setPaying(false);
    }
  };

  const handleEdit = async () => {
    if (!debt) return;
    try {
      const updated = await debtsApi.update(debt.id, {
        name: editForm.name,
        total_amount: parseFloat(editForm.total_amount),
        interest_rate: parseFloat(editForm.interest_rate),
        total_installments: parseInt(editForm.total_installments),
        start_date: editForm.start_date,
        status: editForm.status as any,
        due_day: editForm.due_day ? parseInt(editForm.due_day) : undefined,
      });
      setDebt(updated);
      setShowEditModal(false);
      toast(t('debts.debtUpdated'), 'success');
    } catch {
      toast(t('debts.debtUpdateError'), 'error');
    }
  };

  const handleDelete = async () => {
    if (!debt) return;
    if (!window.confirm(t('debts.deleteConfirm'))) return;
    try {
      await debtsApi.delete(debt.id);
      toast(t('debts.debtDeleted'), 'success');
      navigate('/debts');
    } catch {
      toast(t('debts.debtDeleteError'), 'error');
    }
  };

  if (loading) return <Loading text={t('debts.loadingDebts')} />;
  if (!debt) return <div className="p-8 text-center text-[var(--color-text-muted)]">{t('debts.debtNotFound')}</div>;

  return (
    <PageContainer>
      <div className="mb-6 flex items-center justify-between">
        <Button variant="ghost" onClick={() => navigate('/debts')}>← {t('debts.backToDebts')}</Button>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={() => setShowEditModal(true)}>{t('common.edit')}</Button>
          <Button variant="danger" onClick={handleDelete}>{t('common.delete')}</Button>
        </div>
      </div>

      <Card title={debt.name} subtitle={`${monthly.filter(s => s.paid).length}/${debt.total_installments} ${t('debts.installments')}`}>
        <div className="grid grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <p className="text-sm text-[var(--color-text-secondary)]">{t('debts.totalAmount')}</p>
            <p className="text-lg font-semibold">${formatCurrency(debt.total_amount)}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--color-text-secondary)]">{t('debts.remaining')}</p>
            <p className="text-lg font-semibold">${formatCurrency(debt.remaining_amount)}</p>
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

      <Modal isOpen={showPayModal} onClose={() => setShowPayModal(false)} title={t('debts.pay')}>
        <div className="space-y-6">
          {selectedStatus && (
            <div className="text-sm text-[var(--color-text-secondary)]">
              {formatMonth(selectedStatus.month)} — #{selectedStatus.installment_num}
            </div>
          )}

          <Button
            className="w-full"
            loading={paying}
            onClick={() => selectedStatus && handlePay(selectedStatus.amount_due)}
          >
            {t('debts.payFull')} — ${selectedStatus ? formatCurrency(selectedStatus.amount_due) : '0'}
          </Button>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-[var(--color-border)]" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-[var(--color-bg-primary)] px-2 text-[var(--color-text-muted)]">{t('common.or')}</span>
            </div>
          </div>

          <Input
            label={t('debts.enterAmount')}
            type="number"
            value={payAmount}
            onChange={e => setPayAmount(e.target.value)}
            min="0.01"
            step="0.01"
          />

          <Button
            className="w-full"
            variant="secondary"
            loading={paying}
            onClick={() => handlePay(parseFloat(payAmount) || 0)}
            disabled={!payAmount || parseFloat(payAmount) <= 0}
          >
            {t('debts.payCustom')}
          </Button>
        </div>
      </Modal>

      <Modal isOpen={showEditModal} onClose={() => setShowEditModal(false)} title={t('debts.editDebt')}>
        {debt.status === 'paused' ? (
          <div className="space-y-5">
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-info)]/10 text-[var(--color-info)] text-sm">
              {t('debts.pausedEditInfo')}
            </div>
            <Select
              label={t('debts.status')}
              options={[
                { value: 'paused', label: t('debts.paused') },
                { value: 'active', label: t('debts.active') },
              ]}
              value={editForm.status}
              onChange={e => setEditForm(p => ({ ...p, status: e.target.value }))}
            />
            <div className="pt-2">
              <Button className="w-full" onClick={handleEdit}>{t('common.save')}</Button>
            </div>
          </div>
        ) : (
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
                { value: 'paused', label: t('debts.paused') },
                { value: 'paid', label: t('debts.paid') },
              ]}
              value={editForm.status}
              onChange={e => setEditForm(p => ({ ...p, status: e.target.value }))}
            />
            <div className="pt-2">
              <Button className="w-full" onClick={handleEdit}>{t('debts.saveChanges')}</Button>
            </div>
          </div>
        )}
      </Modal>
    </PageContainer>
  );
}
