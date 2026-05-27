import api from './api';

export interface Debt {
  id: string;
  name: string;
  description: string | null;
  total_amount: number;
  remaining_amount: number;
  interest_rate: number;
  total_installments: number;
  current_installment: number;
  status: string;
  start_date: string;
  end_date: string | null;
  category_id?: string;
}

export interface CreateDebtInput {
  name: string;
  description?: string;
  total_amount: number;
  interest_rate: number;
  total_installments: number;
  start_date: string;
  end_date?: string;
  category_id?: string;
}

export interface DebtMonthlyStatus {
  id: string;
  debt_id: string;
  month: string;
  installment_num: number;
  amount_due: number;
  amount_paid: number;
  paid: boolean;
  paid_at: string | null;
  notes: string | null;
}

export const debtsApi = {
  list: () => api.get('/debts').then(r => r.data.data),
  getByID: (id: string) => api.get(`/debts/${id}`).then(r => r.data.data),
  create: (input: CreateDebtInput) => api.post('/debts', input).then(r => r.data.data),
  update: (id: string, input: Partial<CreateDebtInput>) => api.put(`/debts/${id}`, input).then(r => r.data.data),
  delete: (id: string) => api.delete(`/debts/${id}`),
  getMonthlyStatus: (id: string) => api.get(`/debts/${id}/monthly`).then(r => r.data.data),
  markAsPaid: (id: string, year: number, month: number, amount_paid: number, notes?: string) =>
    api.put(`/debts/${id}/monthly/${year}/${month}/pay`, { amount_paid, notes }).then(r => r.data.data),
};
