import api from './api';
import type { LoginInput, RegisterInput, AuthResponse } from '../types/auth';

export async function login(input: LoginInput): Promise<AuthResponse> {
  const { data } = await api.post('/auth/login', input);
  return data.data;
}

export async function register(input: RegisterInput): Promise<AuthResponse> {
  const { data } = await api.post('/auth/register', input);
  return data.data;
}

export async function refreshToken(token: string): Promise<AuthResponse> {
  const { data } = await api.post('/auth/refresh', { refresh_token: token });
  return data.data;
}
