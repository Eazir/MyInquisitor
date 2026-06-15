import api from './api';

export interface Expense {
  id: string;
  name: string;
  description: string | null;
  amount: number;
  frequency: string;
  due_day: number | null;
  billing_month: number | null;
  status: string;
  start_date: string;
  end_date: string | null;
  category_id?: string;
}

export interface CreateExpenseInput {
  name: string;
  description?: string;
  amount: number;
  frequency: string;
  due_day?: number;
  billing_month?: number;
  start_date: string;
  end_date?: string;
  status?: string;
  category_id?: string;
}

export const expensesApi = {
  list: () => api.get('/expenses', { params: { status: 'active' } }).then(r => r.data.data),
  getByID: (id: string) => api.get(`/expenses/${id}`).then(r => r.data.data),
  create: (input: CreateExpenseInput) => api.post('/expenses', input).then(r => r.data.data),
  update: (id: string, input: Partial<CreateExpenseInput>) => api.put(`/expenses/${id}`, input).then(r => r.data.data),
  delete: (id: string) => api.delete(`/expenses/${id}`),
  togglePaid: (id: string, year: number, month: number) =>
    api.put(`/expenses/${id}/monthly/${year}/${month}/toggle`).then(r => r.data.data),
  getMonthlyStatus: (id: string, year: number, month: number) =>
    api.get(`/expenses/${id}/monthly/${year}/${month}`).then(r => r.data.data),
};
