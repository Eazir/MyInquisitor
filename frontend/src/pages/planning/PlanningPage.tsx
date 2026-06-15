import { useState, useEffect, useCallback } from 'react';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { Badge } from '../../components/ui/Badge';
import { Table, type Column } from '../../components/ui/Table';
import { Loading } from '../../components/ui/Loading';
import { useLanguage } from '../../contexts/LanguageContext';
import { useAuth } from '../../contexts/AuthContext';
import { toast } from '../../components/ui/Toast';
import { accountingApi, type Projection, type ExtraExpenseItem } from '../../services/accounting';
import { formatCurrency } from '../../utils/format';

const ADJ_BASE_KEY = 'myinquisitor_planning';
const SETTINGS_BASE_KEY = 'myinquisitor_planning_settings';

function adjKey(userID: string) { return `${ADJ_BASE_KEY}_${userID}`; }
function settingsKey(userID: string) { return `${SETTINGS_BASE_KEY}_${userID}`; }

interface MonthAdjustment {
  income: number;
  extra_expenses: ExtraExpenseItem[];
  disabled: Record<string, boolean>;
}

interface PlanningSettings {
  months: number;
  historyMonths: number;
}

function loadAdjustments(userID: string): Record<string, MonthAdjustment> {
  try { return JSON.parse(localStorage.getItem(adjKey(userID)) || '{}'); }
  catch { return {}; }
}

function saveAdjustments(userID: string, adj: Record<string, MonthAdjustment>) {
  localStorage.setItem(adjKey(userID), JSON.stringify(adj));
}

function loadSettings(userID: string): PlanningSettings {
  try { return JSON.parse(localStorage.getItem(settingsKey(userID)) || '{}'); }
  catch { return {} as PlanningSettings; }
}

function saveSettings(userID: string, s: PlanningSettings) {
  localStorage.setItem(settingsKey(userID), JSON.stringify(s));
}

function computeTotals(proj: Projection, adj: MonthAdjustment) {
  const disabled = adj.disabled || {};
  const income = disabled.income ? 0 : (adj.income ?? proj.base_income);
  const fixed = (proj.fixed_expenses_list || [])
    .filter(item => !disabled[`fixed_${item.id}`])
    .reduce((s, item) => s + item.amount, 0);
  const extraAvg = disabled.extra_budgetary_avg ? 0 : proj.extra_budgetary_avg;
  const debts = disabled.debts ? 0 : proj.debt_payments;
  const extraTotal = (adj.extra_expenses || []).reduce((s, e, i) => disabled[`extra_${i}`] ? s : s + e.amount, 0);
  const totalExp = fixed + extraAvg + debts + extraTotal;
  return { income, fixed, extraAvg, debts, extraTotal, totalExp, balance: income - totalExp };
}

function applyAdjustments(projections: Projection[], adjustments: Record<string, MonthAdjustment>): Projection[] {
  return projections.map(p => {
    const adj = adjustments[p.month];
    if (!adj) return p;
    const t = computeTotals(p, adj);
    return {
      ...p,
      income_modifier: t.income - p.base_income,
      extra_expenses_list: adj.extra_expenses || [],
      extra_expenses_total: t.extraTotal,
      projected_income: t.income,
      total_expenses: t.totalExp,
      projected_balance: t.balance,
    };
  });
}

function Toggle({ enabled, onClick }: { enabled: boolean; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className={`text-xs font-bold px-2 py-0.5 rounded border transition-colors cursor-pointer ${
        enabled
          ? 'bg-[var(--color-success)]/10 text-[var(--color-success)] border-[var(--color-success)]/30'
          : 'bg-[var(--color-text-muted)]/10 text-[var(--color-text-muted)] border-[var(--color-text-muted)]/20'
      }`}
    >
      {enabled ? 'ON' : 'OFF'}
    </button>
  );
}

export function PlanningPage() {
  const { t } = useLanguage();
  const { user } = useAuth();
  const uid = user?.id || 'anon';
  const [projections, setProjections] = useState<Projection[]>([]);
  const [loading, setLoading] = useState(true);
  const saved = loadSettings(uid);
  const [months, setMonths] = useState(saved.months || 1);
  const [historyMonths, setHistoryMonths] = useState(saved.historyMonths || 12);

  const [selectedMonth, setSelectedMonth] = useState<string | null>(null);
  const [editData, setEditData] = useState<MonthAdjustment>({ income: 0, extra_expenses: [], disabled: {} });
  const [newDesc, setNewDesc] = useState('');
  const [newAmount, setNewAmount] = useState('');

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const data = await accountingApi.getProjections(months, historyMonths);
      const adjustments = loadAdjustments(uid);
      setProjections(applyAdjustments(data, adjustments));
    } catch {
      toast(t('planning.calculatingProjections'), 'error');
    } finally {
      setLoading(false);
    }
  }, [months, historyMonths, uid, t]);

  useEffect(() => { fetchData(); }, [fetchData]);

  useEffect(() => {
    const key = adjKey(uid);
    const handleStorage = (e: StorageEvent) => {
      if (e.key === key && e.newValue) {
        const adjustments = loadAdjustments(uid);
        setProjections(prev => applyAdjustments(prev, adjustments));
      }
    };
    window.addEventListener('storage', handleStorage);
    return () => window.removeEventListener('storage', handleStorage);
  }, [uid]);

  useEffect(() => { saveSettings(uid, { months, historyMonths }); }, [months, historyMonths, uid]);

  const selectedProjection = selectedMonth ? projections.find(p => p.month === selectedMonth) : null;

  const openMonthDetail = (p: Projection) => {
    const existing = loadAdjustments(uid)[p.month];
    const base = existing || { income: p.base_income, extra_expenses: [...p.extra_expenses_list], disabled: {} };
    if (base.disabled.fixed_expenses) {
      for (const item of p.fixed_expenses_list) {
        base.disabled[`fixed_${item.id}`] = true;
      }
      delete base.disabled.fixed_expenses;
    }
    setSelectedMonth(p.month);
    setEditData(base);
    setNewDesc('');
    setNewAmount('');
  };

  const closeMonthDetail = () => setSelectedMonth(null);

  const saveMonthDetail = () => {
    if (!selectedMonth) return;
    const adjustments = loadAdjustments(uid);
    adjustments[selectedMonth] = editData;
    saveAdjustments(uid, adjustments);
    setProjections(prev => applyAdjustments(prev, adjustments));
    setSelectedMonth(null);
    toast(t('planning.adjustmentsSaved'), 'success');
  };

  const clearAll = () => {
    if (!window.confirm(t('planning.clearConfirm'))) return;
    localStorage.removeItem(adjKey(uid));
    fetchData();
    toast(t('planning.adjustmentsCleared'), 'info');
  };

  const toggleRow = (key: string) => {
    setEditData(prev => ({
      ...prev,
      disabled: { ...prev.disabled, [key]: !prev.disabled[key] },
    }));
  };

  const addExtraExpense = () => {
    const amount = parseFloat(newAmount);
    if (!newDesc.trim() || !amount || amount <= 0) return;
    setEditData(prev => ({
      ...prev,
      extra_expenses: [...prev.extra_expenses, { description: newDesc.trim(), amount }],
    }));
    setNewDesc('');
    setNewAmount('');
  };

  const removeExtraExpense = (idx: number) => {
    setEditData(prev => {
      const disabled = { ...prev.disabled };
      delete disabled[`extra_${idx}`];
      const reindexed: Record<string, boolean> = {};
      for (const k of Object.keys(disabled)) {
        const match = k.match(/^extra_(\d+)$/);
        if (match) {
          const oldIdx = parseInt(match[1]);
          reindexed[`extra_${oldIdx > idx ? oldIdx - 1 : oldIdx}`] = disabled[k];
        } else {
          reindexed[k] = disabled[k];
        }
      }
      return {
        ...prev,
        extra_expenses: prev.extra_expenses.filter((_, i) => i !== idx),
        disabled: reindexed,
      };
    });
  };

  const editIncome = (val: string) => {
    setEditData(prev => ({ ...prev, income: parseFloat(val) || 0 }));
  };

  const totals = selectedProjection ? computeTotals(selectedProjection, editData) : null;

  const columns: Column<Projection>[] = [
    {
      key: 'month_label',
      header: t('planning.month'),
      render: (p) => (
        <button
          className="font-medium text-[var(--color-accent)] hover:underline text-left cursor-pointer"
          onClick={() => openMonthDetail(p)}
        >
          {p.month_label}
        </button>
      ),
    },
    {
      key: 'projected_income',
      header: t('planning.income'),
      render: (p) => {
        const adj = p.income_modifier !== 0;
        return (
          <span className={`font-medium ${adj ? 'text-[var(--color-accent)]' : 'text-[var(--color-success)]'}`}>
            {adj && <span className="text-xs opacity-60">{formatCurrency(p.base_income)} → </span>}
            ${formatCurrency(p.projected_income)}
          </span>
        );
      },
    },
    {
      key: 'fixed_expenses',
      header: t('planning.fixedExpenses'),
      render: (p) => <span className="font-medium text-[var(--color-danger)]">${formatCurrency(p.fixed_expenses)}</span>,
    },
    {
      key: 'extra_budgetary_avg',
      header: t('planning.extraBudgetary'),
      render: (p) => <span className="font-medium text-[var(--color-warning)]">${formatCurrency(p.extra_budgetary_avg)}</span>,
    },
    {
      key: 'extra_expenses_total',
      header: t('planning.extraExpenses'),
      render: (p) => {
        if (p.extra_expenses_total === 0) return <span className="text-[var(--color-text-muted)]">—</span>;
        const count = p.extra_expenses_list.length;
        return (
          <span className="font-medium text-[var(--color-danger)]">
            ${formatCurrency(p.extra_expenses_total)}
            {count > 0 && <span className="text-xs opacity-60 ml-1">({count})</span>}
          </span>
        );
      },
    },
    {
      key: 'debt_payments',
      header: t('planning.debts'),
      render: (p) => {
        if (p.debt_payments === 0) return <span className="text-[var(--color-text-muted)]">—</span>;
        return <span className="font-medium text-[var(--color-warning)]">${formatCurrency(p.debt_payments)}</span>;
      },
    },
    {
      key: 'total_expenses',
      header: t('planning.totalExpenses'),
      render: (p) => <span className="font-medium text-[var(--color-danger)]">${formatCurrency(p.total_expenses)}</span>,
    },
    {
      key: 'projected_balance',
      header: t('planning.balance'),
      render: (p) => {
        const color = p.projected_balance >= 0 ? 'var(--color-success)' : 'var(--color-danger)';
        return <span className="font-bold" style={{ color }}>${formatCurrency(p.projected_balance)}</span>;
      },
    },
  ];

  return (
    <PageContainer>
      {selectedProjection && totals ? (
        <>
          <div className="mb-6 flex items-center justify-between">
            <Button variant="ghost" onClick={closeMonthDetail}>← {t('common.back')}</Button>
            <div className="flex gap-2">
              <Button variant="secondary" onClick={closeMonthDetail}>{t('planning.cancel')}</Button>
              <Button onClick={saveMonthDetail}>{t('planning.save')}</Button>
            </div>
          </div>

          <h2 className="text-2xl font-bold text-[var(--color-text-primary)] mb-2">
            {t('planning.editMonth')} — {selectedProjection.month_label}
          </h2>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6 mb-8">
            <Card variant="stats">
              <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.baseIncome')}</p>
              <p className="text-3xl font-bold text-[var(--color-text-primary)]">
                ${formatCurrency(selectedProjection.base_income)}
              </p>
            </Card>
            <Card variant="stats">
              <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.income')}</p>
              <p className={`text-3xl font-bold ${totals.balance >= 0 ? 'text-[var(--color-success)]' : 'text-[var(--color-danger)]'}`}>
                ${formatCurrency(totals.income)}
              </p>
            </Card>
            <Card variant="stats">
              <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.totalExpenses')}</p>
              <p className="text-3xl font-bold text-[var(--color-danger)]">
                ${formatCurrency(totals.totalExp)}
              </p>
            </Card>
            <Card variant="stats">
              <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.balance')}</p>
              <p className={`text-3xl font-bold ${totals.balance >= 0 ? 'text-[var(--color-success)]' : 'text-[var(--color-danger)]'}`}>
                ${formatCurrency(totals.balance)}
              </p>
            </Card>
          </div>

          <Card>
            <div className="divide-y divide-[var(--color-border)]">
              <div className="hidden sm:grid sm:grid-cols-12 gap-4 px-4 py-3 text-sm font-bold uppercase text-[var(--color-text-primary)]">
                <div className="sm:col-span-3 lg:col-span-2">{t('planning.typeColumn')}</div>
                <div className="sm:col-span-3 lg:col-span-4">{t('planning.reasonColumn')}</div>
                <div className="sm:col-span-2 lg:col-span-2 text-right">{t('planning.amount')}</div>
                <div className="sm:col-span-4 lg:col-span-4 text-right">{t('common.actions')}</div>
              </div>

              {/* Income row - editable */}
              <div className="grid grid-cols-1 sm:grid-cols-12 gap-2 sm:gap-4 px-4 py-4 items-start sm:items-center">
                <div className="flex items-center gap-2 sm:col-span-3 lg:col-span-2">
                  <Badge variant="success">{t('planning.income')}</Badge>
                  <span className="sm:hidden font-medium text-[var(--color-text-primary)] truncate">{t('planning.income')}</span>
                </div>
                <div className="hidden sm:block sm:col-span-3 lg:col-span-4 font-medium text-[var(--color-text-primary)] truncate">{t('planning.income')}</div>
                <div className="sm:col-span-2 lg:col-span-2">
                  <Input
                    type="number"
                    value={String(editData.income)}
                    onChange={e => editIncome(e.target.value)}
                    className="!w-full !max-w-[140px] sm:!ml-auto"
                    step="0.01"
                  />
                </div>
                <div className="flex items-center gap-2 sm:col-span-4 lg:col-span-4 sm:justify-end">
                  <span className="text-xs text-[var(--color-text-muted)]">
                    {t('planning.baseIncome')}: ${formatCurrency(selectedProjection.base_income)}
                  </span>
                  <Toggle enabled={!editData.disabled.income} onClick={() => toggleRow('income')} />
                </div>
              </div>

              {/* Fixed expenses - itemized */}
              {selectedProjection.fixed_expenses_list.map(item => {
                const isDisabled = !!editData.disabled[`fixed_${item.id}`];
                return (
                  <div key={item.id} className={`grid grid-cols-2 sm:grid-cols-12 gap-2 sm:gap-4 px-4 py-4 items-center ${isDisabled ? 'opacity-50' : ''}`}>
                    <div className="col-span-2 sm:col-span-3 lg:col-span-2 flex items-center gap-2">
                      <Badge variant="danger">G. Fijos</Badge>
                      <span className="sm:hidden font-medium text-[var(--color-text-primary)] truncate">{item.name}</span>
                    </div>
                    <div className="hidden sm:block sm:col-span-3 lg:col-span-4 font-medium text-[var(--color-text-primary)] truncate">{item.name}</div>
                    <div className="col-span-1 sm:col-span-2 lg:col-span-2 text-right font-semibold text-[var(--color-danger)]">
                      ${formatCurrency(item.amount)}
                    </div>
                    <div className="col-span-1 sm:col-span-4 lg:col-span-4 flex justify-end">
                      <Toggle enabled={!isDisabled} onClick={() => toggleRow(`fixed_${item.id}`)} />
                    </div>
                  </div>
                );
              })}

              {/* Extra-budgetary avg */}
              <div className="grid grid-cols-2 sm:grid-cols-12 gap-2 sm:gap-4 px-4 py-4 items-center">
                <div className="col-span-2 sm:col-span-3 lg:col-span-2 flex items-center gap-2">
                  <Badge variant="warning">{t('planning.extraBudgetary')}</Badge>
                  <span className="sm:hidden font-medium text-[var(--color-text-primary)] truncate">{t('planning.extraBudgetary')}</span>
                </div>
                <div className="hidden sm:block sm:col-span-3 lg:col-span-4 font-medium text-[var(--color-text-primary)] truncate">{t('planning.extraBudgetary')}</div>
                <div className="col-span-1 sm:col-span-2 lg:col-span-2 text-right font-semibold text-[var(--color-warning)]">
                  ${formatCurrency(selectedProjection.extra_budgetary_avg)}
                </div>
                <div className="col-span-1 sm:col-span-4 lg:col-span-4 flex justify-end">
                  <Toggle enabled={!editData.disabled.extra_budgetary_avg} onClick={() => toggleRow('extra_budgetary_avg')} />
                </div>
              </div>

              {/* Debts */}
              <div className="grid grid-cols-2 sm:grid-cols-12 gap-2 sm:gap-4 px-4 py-4 items-center">
                <div className="col-span-2 sm:col-span-3 lg:col-span-2 flex items-center gap-2">
                  <Badge variant="warning">{t('planning.debts')}</Badge>
                  <span className="sm:hidden font-medium text-[var(--color-text-primary)] truncate">{t('planning.debts')}</span>
                </div>
                <div className="hidden sm:block sm:col-span-3 lg:col-span-4 font-medium text-[var(--color-text-primary)] truncate">{t('planning.debts')}</div>
                <div className="col-span-1 sm:col-span-2 lg:col-span-2 text-right font-semibold text-[var(--color-warning)]">
                  ${formatCurrency(selectedProjection.debt_payments)}
                </div>
                <div className="col-span-1 sm:col-span-4 lg:col-span-4 flex justify-end">
                  <Toggle enabled={!editData.disabled.debts} onClick={() => toggleRow('debts')} />
                </div>
              </div>

              {/* Extra expenses */}
              {editData.extra_expenses.map((item, idx) => {
                const isDisabled = !!editData.disabled[`extra_${idx}`];
                return (
                  <div key={idx} className={`grid grid-cols-2 sm:grid-cols-12 gap-2 sm:gap-4 px-4 py-4 items-center ${isDisabled ? 'opacity-50' : ''}`}>
                    <div className="col-span-2 sm:col-span-3 lg:col-span-2 flex items-center gap-2">
                      <Badge variant="danger">{t('planning.extraExpenses')}</Badge>
                      <span className="sm:hidden font-medium text-[var(--color-text-primary)] truncate">{item.description}</span>
                    </div>
                    <div className="hidden sm:block sm:col-span-3 lg:col-span-4 font-medium text-[var(--color-text-primary)] truncate">{item.description}</div>
                    <div className="col-span-1 sm:col-span-2 lg:col-span-2 text-right font-semibold text-[var(--color-danger)]">
                      ${formatCurrency(item.amount)}
                    </div>
                    <div className="col-span-1 sm:col-span-4 lg:col-span-4 flex items-center justify-end gap-2">
                      <Toggle enabled={!isDisabled} onClick={() => toggleRow(`extra_${idx}`)} />
                      <Button size="sm" variant="ghost" onClick={() => removeExtraExpense(idx)}>{t('planning.remove')}</Button>
                    </div>
                  </div>
                );
              })}

              {/* Add extra expense */}
              <div className="grid grid-cols-1 sm:grid-cols-12 gap-2 sm:gap-4 px-4 py-4 items-start sm:items-center">
                <div className="sm:col-span-3 lg:col-span-2" />
                <div className="sm:col-span-3 lg:col-span-4">
                  <Input
                    placeholder={t('planning.descLabel')}
                    value={newDesc}
                    onChange={e => setNewDesc(e.target.value)}
                  />
                </div>
                <div className="flex items-center gap-2 sm:col-span-6 lg:col-span-6 sm:justify-end">
                  <Input
                    type="number"
                    placeholder={t('planning.amount')}
                    value={newAmount}
                    onChange={e => setNewAmount(e.target.value)}
                    step="0.01"
                    className="!w-28"
                  />
                  <Button size="sm" variant="secondary" onClick={addExtraExpense}>{t('planning.add')}</Button>
                </div>
              </div>

              {/* Total row */}
              <div className="grid grid-cols-1 sm:grid-cols-12 gap-2 sm:gap-4 px-4 py-4 items-center bg-[var(--color-bg-secondary)] font-bold">
                <div className="hidden sm:block sm:col-span-3 lg:col-span-2" />
                <div className="sm:col-span-3 lg:col-span-4 text-[var(--color-text-primary)]">{t('planning.totalExpenses')}</div>
                <div className="sm:col-span-2 lg:col-span-2 text-right text-[var(--color-danger)]">
                  ${formatCurrency(totals.totalExp)}
                </div>
                <div className="sm:col-span-4 lg:col-span-4 text-right text-[var(--color-success)]">
                  {t('planning.balance')}: ${formatCurrency(totals.balance)}
                </div>
              </div>
            </div>
          </Card>
        </>
      ) : (
        <>
          <div className="flex items-center justify-between mb-8">
            <div>
              <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('planning.title')}</h2>
              <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('planning.description')}</p>
            </div>
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <label className="text-sm font-medium text-[var(--color-text-secondary)]">{t('planning.months')}:</label>
                <Input
                  type="number"
                  value={String(months)}
                  onChange={e => setMonths(Math.max(1, parseInt(e.target.value) || 1))}
                  className="!w-20"
                  min={1}
                  max={24}
                />
              </div>
              <div className="flex items-center gap-2">
                <label className="text-sm font-medium text-[var(--color-text-secondary)]">{t('planning.historyMonths')}:</label>
                <Input
                  type="number"
                  value={String(historyMonths)}
                  onChange={e => setHistoryMonths(Math.max(1, Math.min(36, parseInt(e.target.value) || 12)))}
                  className="!w-20"
                  min={1}
                  max={36}
                />
              </div>
              <Button onClick={() => fetchData()}>{t('planning.update')}</Button>
              <Button variant="ghost" onClick={clearAll}>{t('planning.clearAdjustments')}</Button>
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
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6">
                    <div>
                      <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.avgMonthlyIncome')}</p>
                      <p className="text-xl font-bold text-[var(--color-success)]">
                        ${formatCurrency(projections.reduce((s, p) => s + p.projected_income, 0) / projections.length)}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.avgMonthlyExpenses')}</p>
                      <p className="text-xl font-bold text-[var(--color-danger)]">
                        ${formatCurrency(projections.reduce((s, p) => s + p.total_expenses, 0) / projections.length)}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-[var(--color-text-secondary)] mb-1">{t('planning.totalProjectedBalance')}</p>
                      <p className="text-xl font-bold text-[var(--color-text-primary)]">
                        ${formatCurrency(projections.reduce((s, p) => s + p.projected_balance, 0))}
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
        </>
      )}
    </PageContainer>
  );
}
