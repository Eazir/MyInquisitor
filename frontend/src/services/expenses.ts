import api from './api';

export interface Expense {
  id: string;
  name: string;
  description: string | null;
  amount: number;
  frequency: string;
  due_day: number | null;
  status: string;
  category_id?: string;
}

export interface CreateExpenseInput {
  name: string;
  description?: string;
  amount: number;
  frequency: string;
  due_day?: number;
  start_date: string;
  end_date?: string;
  category_id?: string;
}

export const expensesApi = {
  list: () => api.get('/expenses').then(r => r.data.data),
  getByID: (id: string) => api.get(`/expenses/${id}`).then(r => r.data.data),
  create: (input: CreateExpenseInput) => api.post('/expenses', input).then(r => r.data.data),
  update: (id: string, input: Partial<CreateExpenseInput>) => api.put(`/expenses/${id}`, input).then(r => r.data.data),
  delete: (id: string) => api.delete(`/expenses/${id}`),
  togglePaid: (id: string, year: number, month: number) =>
    api.put(`/expenses/${id}/monthly/${year}/${month}/toggle`).then(r => r.data.data),
};
