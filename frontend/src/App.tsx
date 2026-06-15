import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { LanguageProvider } from './contexts/LanguageContext';
import { AuthProvider } from './contexts/AuthContext';
import { AppLayout } from './components/layout/AppLayout';
import { PrivateRoute } from './components/layout/PrivateRoute';
import { AdminRoute } from './components/layout/AdminRoute';
import { LoginPage } from './pages/auth/LoginPage';
import { RegisterPage } from './pages/auth/RegisterPage';
import { DashboardPage } from './pages/dashboard/DashboardPage';
import { DebtsListPage } from './pages/debts/DebtsListPage';
import { DebtDetailPage } from './pages/debts/DebtDetailPage';
import { ExpensesPage } from './pages/expenses/ExpensesPage';
import { AccountingPage } from './pages/accounting/AccountingPage';
import { PlanningPage } from './pages/planning/PlanningPage';
import { MonthlyPaymentsPage } from './pages/monthly/MonthlyPaymentsPage';
import { AdminPage } from './pages/admin/AdminPage';
import { SettingsPage } from './pages/settings/SettingsPage';
import { ToastContainer } from './components/ui/Toast';

export default function App() {
  return (
    <BrowserRouter>
      <ThemeProvider>
        <LanguageProvider>
          <AuthProvider>
          <ToastContainer />
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register/:token" element={<RegisterPage />} />
            <Route element={<PrivateRoute />}>
              <Route element={<AppLayout />}>
                <Route path="/" element={<Navigate to="/dashboard" replace />} />
                <Route path="/dashboard" element={<DashboardPage />} />
                <Route path="/debts" element={<DebtsListPage />} />
                <Route path="/debts/:id" element={<DebtDetailPage />} />
                <Route path="/expenses" element={<ExpensesPage />} />
                <Route path="/accounting" element={<AccountingPage />} />
                <Route path="/planning" element={<PlanningPage />} />
                <Route path="/monthly-payments" element={<MonthlyPaymentsPage />} />
                <Route path="/settings" element={<SettingsPage />} />
                <Route element={<AdminRoute />}>
                  <Route path="/admin" element={<AdminPage />} />
                </Route>
              </Route>
            </Route>
            <Route path="/register" element={<RegisterPage />} />
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </AuthProvider>
      </LanguageProvider>
    </ThemeProvider>
  </BrowserRouter>
  );
}
