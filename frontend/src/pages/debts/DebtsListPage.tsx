import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { PageContainer } from '../../components/layout/PageContainer';
import { Button } from '../../components/ui/Button';
import { Table, type Column } from '../../components/ui/Table';
import { Badge } from '../../components/ui/Badge';
import { Modal } from '../../components/ui/Modal';
import { Input } from '../../components/ui/Input';
import { Loading } from '../../components/ui/Loading';
import { EmptyState } from '../../components/ui/EmptyState';
import { useLanguage } from '../../contexts/LanguageContext';
import { debtsApi, type Debt } from '../../services/debts';

export function DebtsListPage() {
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [debts, setDebts] = useState<Debt[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({
    name: '', total_amount: '', interest_rate: '0', total_installments: '1', start_date: '',
  });

  const statusBadge = (status: string) => {
    const map: Record<string, { variant: 'info' | 'success' | 'warning'; label: string }> = {
      active: { variant: 'warning', label: t('debts.active') },
      paid: { variant: 'success', label: t('debts.paid') },
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
      render: (d) => <span className="font-medium">${d.total_amount.toFixed(2)}</span>,
    },
    {
      key: 'remaining_amount',
      header: t('debts.remaining'),
      render: (d) => <span className="font-medium text-[var(--color-danger)]">${d.remaining_amount.toFixed(2)}</span>,
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
    });
    setShowModal(false);
    setForm({ name: '', total_amount: '', interest_rate: '0', total_installments: '1', start_date: '' });
    load();
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
          <div className="pt-2">
            <Button className="w-full" onClick={handleCreate}>{t('debts.createDebt')}</Button>
          </div>
        </div>
      </Modal>
    </PageContainer>
  );
}
