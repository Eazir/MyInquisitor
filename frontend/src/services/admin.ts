import api from './api';

export interface AdminUser {
  id: string;
  email: string;
  full_name: string;
  role: string;
  active: boolean;
  created_at: string;
}

export interface CreateUserInput {
  email: string;
  password: string;
  full_name: string;
  role?: string;
}

export interface UpdateUserInput {
  email?: string;
  full_name?: string;
  role?: string;
  active?: boolean;
  admin_password: string;
}

export interface InviteToken {
  id: string;
  token: string;
  created_by: string;
  creator_name: string | null;
  expires_at: string;
  url: string;
  used: boolean;
  created_at: string;
  expired: boolean;
}

export const adminApi = {
  listUsers: (page = 1, limit = 20) =>
    api.get('/admin/users', { params: { page, limit } }).then(r => r.data),
  createUser: (input: CreateUserInput) =>
    api.post('/admin/users', input).then(r => r.data.data),
  updateUser: (id: string, input: UpdateUserInput) =>
    api.put(`/admin/users/${id}`, input).then(r => r.data.data),
  setActive: (id: string, active: boolean) =>
    api.put(`/admin/users/${id}/activate/${active}`).then(r => r.data.data),
  generateInvite: () =>
    api.post('/admin/invite').then(r => r.data.data),
  listInvites: () =>
    api.get('/admin/invites').then(r => r.data.data),
  deleteInvite: (id: string) =>
    api.delete(`/admin/invites/${id}`),
};
