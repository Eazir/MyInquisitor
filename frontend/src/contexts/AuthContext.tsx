import type { ReactNode } from 'react';
import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import api, { setAccessToken } from '../services/api';
import { login as loginApi, register as registerApi } from '../services/auth';
import type { User } from '../types/auth';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string, inviteToken: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  isAdmin: boolean;
  setUser: (user: User | null) => void;
}

const AuthContext = createContext<AuthContextType>({} as AuthContextType);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const stored = localStorage.getItem('refreshToken');
    if (stored) {
      api.post('/auth/refresh', { refresh_token: stored })
        .then(({ data }) => {
          setAccessToken(data.data.access_token);
          localStorage.setItem('refreshToken', data.data.refresh_token);
          setUser(data.data.user);
        })
        .catch(() => {
          localStorage.removeItem('refreshToken');
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    document.documentElement.classList.toggle('super-admin', user?.role === 'super_admin');
  }, [user]);

  const login = useCallback(async (email: string, password: string) => {
    const result = await loginApi({ email, password });
    setAccessToken(result.access_token);
    localStorage.setItem('refreshToken', result.refresh_token);
    setUser(result.user);
  }, []);

  const register = useCallback(async (email: string, password: string, fullName: string, inviteToken: string) => {
    const result = await registerApi({ email, password, full_name: fullName, invite_token: inviteToken });
    setAccessToken(result.access_token);
    localStorage.setItem('refreshToken', result.refresh_token);
    setUser(result.user);
  }, []);

  const logout = useCallback(() => {
    setAccessToken(null);
    localStorage.removeItem('refreshToken');
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{
      user,
      loading,
      login,
      register,
      logout,
      isAuthenticated: !!user,
      isAdmin: user?.role === 'super_admin',
      setUser,
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
