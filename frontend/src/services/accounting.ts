import api from './api';

export interface Transaction {
  id: string;
  type: 'income' | 'expense' | 'transfer';
  amount: number;
  description: string | null;
  source: string | null;
  reference_date: string;
  category_id?: string;
}

export interface CreateTransactionInput {
  type: string;
  amount: number;
  description?: string;
  source?: string;
  reference_date: string;
  category_id?: string;
}

export interface MonthlyBalance {
  month: string;
  total_income: number;
  total_expenses: number;
  total_debt_payments: number;
  total_obligations: number;
  net_balance: number;
  projected_income: number | null;
  income_breakdown: Record<string, number>;
  expense_breakdown: Record<string, number>;
}

export interface CashFlowItem {
  period: string;
  income: number;
  expenses: number;
  balance: number;
}

export interface ExtraExpenseItem {
  description: string;
  amount: number;
}

export interface FixedExpenseItem {
  id: string;
  name: string;
  amount: number;
  frequency: string;
  due_day: number | null;
}

export interface Projection {
  month: string;
  month_label: string;
  base_income: number;
  income_modifier: number;
  projected_income: number;
  fixed_expenses: number;
  fixed_expenses_list: FixedExpenseItem[];
  extra_budgetary_avg: number;
  extra_expenses_total: number;
  extra_expenses_list: ExtraExpenseItem[];
  debt_payments: number;
  total_expenses: number;
  projected_balance: number;
}

export interface Category {
  id: string;
  name: string;
  type: string;
  icon: string | null;
  color: string | null;
}

export const accountingApi = {
  recordTransaction: (input: CreateTransactionInput) =>
    api.post('/accounting/transactions', input).then(r => r.data.data),
  listTransactions: (params?: { year?: number; month?: number; page?: number; limit?: number }) =>
    api.get('/accounting/transactions', { params }).then(r => r.data),
  getMonthlyBalance: (year: number, month: number) =>
    api.get(`/accounting/balance/${year}/${month}`).then(r => r.data.data),
  getCashFlow: () => api.get('/accounting/cash-flow').then(r => r.data.data),
  getProjections: (months?: number, historyMonths?: number) =>
    api.get('/accounting/projections', { params: { months, history_months: historyMonths } }).then(r => r.data.data),
};

export const categoriesApi = {
  list: () => api.get('/categories').then(r => r.data.data),
  create: (input: { name: string; type: string; icon?: string; color?: string }) =>
    api.post('/categories', input).then(r => r.data.data),
  delete: (id: string) => api.delete(`/categories/${id}`),
};
