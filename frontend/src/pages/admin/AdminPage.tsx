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
import { adminApi, type AdminUser, type UpdateUserInput, type InviteToken } from '../../services/admin';

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

  const [invites, setInvites] = useState<InviteToken[]>([]);
  const [loadingInvites, setLoadingInvites] = useState(true);

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

  const inviteColumns: Column<InviteToken>[] = [
    {
      key: 'token',
      header: t('admin.inviteToken'),
      render: (inv) => <code className="text-xs font-mono">{inv.token.slice(0, 16)}...</code>,
    },
    {
      key: 'status',
      header: t('admin.inviteStatus'),
      render: (inv) => {
        if (inv.used) return <Badge variant="info">{t('admin.inviteUsed')}</Badge>;
        if (inv.expired) return <Badge variant="danger">{t('admin.inviteExpired')}</Badge>;
        return <Badge variant="success">{t('admin.inviteActive')}</Badge>;
      },
    },
    {
      key: 'creator_name',
      header: t('admin.inviteCreatedBy'),
      render: (inv) => inv.creator_name || '-',
    },
    {
      key: 'expires_at',
      header: t('admin.inviteExpires'),
      render: (inv) => new Date(inv.expires_at).toLocaleDateString(),
    },
    {
      key: 'created_at',
      header: t('admin.inviteCreated'),
      render: (inv) => new Date(inv.created_at).toLocaleDateString(),
    },
    {
      key: 'actions',
      header: t('admin.inviteActions'),
      render: (inv) => (
        <div className="flex gap-2">
          <Button size="sm" variant="secondary" onClick={() => handleCopyInvite(inv)}>
            {t('admin.inviteCopy')}
          </Button>
          <Button size="sm" variant="danger" onClick={() => handleDeleteInvite(inv)}>
            {t('admin.inviteDelete')}
          </Button>
        </div>
      ),
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

  const loadInvites = async () => {
    setLoadingInvites(true);
    try {
      const data = await adminApi.listInvites();
      setInvites(data || []);
    } finally {
      setLoadingInvites(false);
    }
  };

  useEffect(() => { load(); loadInvites(); }, []);

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
      loadInvites();
    } catch {
      toast(t('admin.inviteError'), 'error');
    } finally {
      setGeneratingInvite(false);
    }
  };

  const handleCopyInvite = async (inv: InviteToken) => {
    const url = `${window.location.origin}${inv.url}`;
    await navigator.clipboard.writeText(url);
    toast(t('admin.inviteCopied'), 'success');
  };

  const handleDeleteInvite = async (inv: InviteToken) => {
    if (!window.confirm(t('admin.inviteDeleteConfirm'))) return;
    try {
      await adminApi.deleteInvite(inv.id);
      toast(t('admin.inviteTokenDeleted'), 'success');
      loadInvites();
    } catch {
      toast(t('admin.inviteDeleteError'), 'error');
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

      <Card title={t('admin.pendingInvites')} className="mt-8">
        {loadingInvites ? (
          <Loading text={t('common.loading')} />
        ) : invites.length === 0 ? (
          <p className="text-sm text-[var(--color-text-secondary)] py-4 text-center">{t('admin.noInvites')}</p>
        ) : (
          <Table columns={inviteColumns} data={invites} variant="striped" />
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
