import { useState, useEffect } from 'react';
import { PageContainer } from '../../components/layout/PageContainer';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Table, type Column } from '../../components/ui/Table';
import { Badge } from '../../components/ui/Badge';
import { Modal } from '../../components/ui/Modal';
import { Input } from '../../components/ui/Input';
import { Select } from '../../components/ui/Select';
import { Loading } from '../../components/ui/Loading';
import { useLanguage } from '../../contexts/LanguageContext';
import { toast } from '../../components/ui/Toast';
import { adminApi, type AdminUser, type UpdateUserInput } from '../../services/admin';

export function AdminPage() {
  const { t } = useLanguage();
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [showCreateUser, setShowCreateUser] = useState(false);
  const [generatingInvite, setGeneratingInvite] = useState(false);
  const [form, setForm] = useState({ email: '', password: '', full_name: '' });

  const [showEditUser, setShowEditUser] = useState(false);
  const [editingUser, setEditingUser] = useState<AdminUser | null>(null);
  const [editForm, setEditForm] = useState({ full_name: '', email: '', role: 'user', admin_password: '' });
  const [editError, setEditError] = useState('');

  const columns: Column<AdminUser>[] = [
    { key: 'full_name', header: t('admin.name') },
    { key: 'email', header: t('admin.email') },
    {
      key: 'role',
      header: t('admin.role'),
      render: (u) => u.role === 'super_admin'
        ? <Badge variant="info">{t('admin.superAdmin')}</Badge>
        : <Badge variant="success">{t('admin.user')}</Badge>,
    },
    {
      key: 'active',
      header: t('admin.active'),
      render: (u) => u.active
        ? <Badge variant="success">{t('admin.activeLabel')}</Badge>
        : <Badge variant="danger">{t('admin.inactive')}</Badge>,
    },
    {
      key: 'created_at',
      header: t('admin.created'),
      render: (u) => new Date(u.created_at).toLocaleDateString(),
    },
  ];

  const load = async () => {
    setLoading(true);
    try {
      const { data, meta } = await adminApi.listUsers();
      setUsers(data);
      setTotal(meta?.total || 0);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    await adminApi.createUser(form);
    setShowCreateUser(false);
    setForm({ email: '', password: '', full_name: '' });
    load();
  };

  const handleToggleActive = async (user: AdminUser) => {
    await adminApi.setActive(user.id, !user.active);
    load();
  };

  const handleGenerateInvite = async () => {
    setGeneratingInvite(true);
    try {
      const result = await adminApi.generateInvite();
      const url = `${window.location.origin}/register/${result.token}`;
      await navigator.clipboard.writeText(url);
      toast(t('admin.inviteCopied'), 'success');
    } catch {
      toast(t('admin.inviteError'), 'error');
    } finally {
      setGeneratingInvite(false);
    }
  };

  const openEdit = (user: AdminUser) => {
    setEditingUser(user);
    setEditForm({
      full_name: user.full_name,
      email: user.email,
      role: user.role,
      admin_password: '',
    });
    setEditError('');
    setShowEditUser(true);
  };

  const handleEdit = async () => {
    if (!editingUser) return;
    if (!editForm.admin_password) {
      setEditError(t('admin.youMustEnterAdminPassword'));
      return;
    }
    try {
      const input: UpdateUserInput = {
        full_name: editForm.full_name,
        email: editForm.email,
        role: editForm.role,
        admin_password: editForm.admin_password,
      };
      await adminApi.updateUser(editingUser.id, input);
      setShowEditUser(false);
      setEditingUser(null);
      load();
    } catch {
      setEditError(t('admin.invalidAdminPassword'));
    }
  };

  return (
    <PageContainer>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h2 className="text-2xl font-bold text-[var(--color-text-primary)]">{t('admin.title')}</h2>
          <p className="text-sm text-[var(--color-text-secondary)] mt-1">{t('admin.description')}</p>
        </div>
        <div className="flex gap-3">
          <Button variant="secondary" loading={generatingInvite} onClick={handleGenerateInvite}>{t('admin.generateInvite')}</Button>
          <Button onClick={() => setShowCreateUser(true)}>{t('admin.createUser')}</Button>
        </div>
      </div>

      <Card title={`${t('admin.users')} (${total})`}>
        {loading ? (
          <Loading text={t('admin.loadingUsers')} />
        ) : (
          <Table
            columns={[
              ...columns,
              {
                key: 'actions',
                header: t('admin.actions'),
                render: (u: AdminUser) => (
                  <div className="flex gap-2">
                    <Button size="sm" variant="secondary" onClick={() => openEdit(u)}>
                      {t('admin.edit')}
                    </Button>
                    <Button
                      size="sm"
                      variant={u.active ? 'danger' : 'secondary'}
                      onClick={() => handleToggleActive(u)}
                    >
                      {u.active ? t('admin.deactivate') : t('admin.activate')}
                    </Button>
                  </div>
                ),
              },
            ]}
            data={users}
            variant="striped"
          />
        )}
      </Card>

      <Modal isOpen={showCreateUser} onClose={() => setShowCreateUser(false)} title={t('admin.createUserForm')}>
        <div className="space-y-5">
          <Input label={t('admin.fullName')} value={form.full_name} onChange={e => setForm(p => ({ ...p, full_name: e.target.value }))} />
          <Input label={t('admin.email')} type="email" value={form.email} onChange={e => setForm(p => ({ ...p, email: e.target.value }))} />
          <Input label={t('admin.password')} type="password" value={form.password} onChange={e => setForm(p => ({ ...p, password: e.target.value }))} />
          <div className="pt-2">
            <Button className="w-full" onClick={handleCreate}>{t('admin.createUser')}</Button>
          </div>
        </div>
      </Modal>

      <Modal isOpen={showEditUser} onClose={() => { setShowEditUser(false); setEditingUser(null); }} title={t('admin.editUser', { name: editingUser?.full_name || '' })}>
        <div className="space-y-5">
          {editError && (
            <div className="p-4 rounded-[var(--radius-md)] bg-[var(--color-danger)]/10 text-[var(--color-danger)] text-sm">
              {editError}
            </div>
          )}
          <Input
            label={t('admin.fullName')}
            value={editForm.full_name}
            onChange={e => setEditForm(p => ({ ...p, full_name: e.target.value }))}
          />
          <Input
            label={t('admin.email')}
            type="email"
            value={editForm.email}
            onChange={e => setEditForm(p => ({ ...p, email: e.target.value }))}
          />
          <Select
            label={t('admin.role')}
            options={[
              { value: 'user', label: t('admin.user') },
              { value: 'super_admin', label: t('admin.superAdmin') },
            ]}
            value={editForm.role}
            onChange={e => setEditForm(p => ({ ...p, role: e.target.value }))}
          />
          <Input
            label={t('admin.yourAdminPassword')}
            type="password"
            value={editForm.admin_password}
            onChange={e => setEditForm(p => ({ ...p, admin_password: e.target.value }))}
            placeholder={t('admin.requiredToSave')}
          />
          <div className="pt-2">
            <Button className="w-full" onClick={handleEdit}>{t('admin.saveChanges')}</Button>
          </div>
        </div>
      </Modal>
    </PageContainer>
  );
}
