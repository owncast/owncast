import React, { ReactElement, useContext, useEffect, useState } from 'react';
import { Table, Input, Select, Tag, Space, Typography, Tabs } from 'antd';
import { ColumnsType, SortOrder } from 'antd/lib/table/interface';
import { useTranslation } from 'next-export-i18n';
import { AdminLayout } from '../../components/layouts/AdminLayout';
import { ServerStatusContext } from '../../utils/server-status-context';
import {
  USERS,
  MODERATORS,
  DISABLED_USERS,
  CONNECTED_CLIENTS,
  BANNED_IPS,
  fetchData,
} from '../../utils/apis';
import { User, Client } from '../../types/chat';
import { formatDateOnly } from '../../utils/format';
import { ClientTable } from '../../components/admin/ClientTable';
import { BannedIPsTable } from '../../components/admin/BannedIPsTable';
import { UserActionsDropdown } from '../../components/admin/UserActionsDropdown';
import { UserDetailsModal } from '../../components/admin/UserDetailsModal';
import { Translation } from '../../components/ui/Translation/Translation';
import { Localization } from '../../types/localization';

const PAGE_SIZE = 25;
const LIVE_FETCH_INTERVAL = 10 * 1000; // 10 sec

const STATUS_OPTIONS = [
  { value: 'all', label: 'All' },
  { value: 'active', label: 'Active' },
  { value: 'bots', label: 'Bots' },
];

export default function UsersAdmin() {
  const context = useContext(ServerStatusContext);
  const { online } = context || {};
  const { t } = useTranslation();

  // "All Users" tab state (server-side paginated + filtered).
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [search, setSearch] = useState<string>('');
  const [status, setStatus] = useState<string>('all');
  // Server-side sort of the paginated list by creation date ('desc' = newest
  // first, 'asc' = oldest first). The API sorts across every page, not just the
  // rows already fetched.
  const [sort, setSort] = useState<'asc' | 'desc'>('desc');
  const [loading, setLoading] = useState<boolean>(false);

  // Exhaustive lists (the paginated endpoint can't return a complete set).
  const [moderators, setModerators] = useState<User[]>([]);
  const [bannedUsers, setBannedUsers] = useState<User[]>([]);

  // Live tabs state (connected clients + IP bans).
  const [clients, setClients] = useState<Client[]>([]);
  const [ipBans, setIPBans] = useState([]);

  // The user whose detail modal is open (null when closed). Set by clicking a
  // table row.
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const getUsers = async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams({
        offset: String((currentPage - 1) * PAGE_SIZE),
        limit: String(PAGE_SIZE),
        sort,
      });
      if (search) {
        params.set('search', search);
      }
      if (status && status !== 'all') {
        params.set('status', status);
      }
      const result = await fetchData(`${USERS}?${params.toString()}`, { auth: true });
      setUsers(result?.results || []);
      setTotal(result?.total || 0);
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error(e);
      setUsers([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  // Moderators and banned users are complete sets, not pages, so they come from
  // their dedicated endpoints rather than the paginated users API.
  const getLists = async () => {
    try {
      setModerators((await fetchData(MODERATORS, { auth: true })) || []);
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error('error fetching moderators', e);
    }
    try {
      setBannedUsers((await fetchData(DISABLED_USERS, { auth: true })) || []);
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error('error fetching banned users', e);
    }
  };

  const getLiveInfo = async () => {
    try {
      setClients(await fetchData(CONNECTED_CLIENTS));
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error('error fetching connected clients', e);
    }
    try {
      setIPBans(await fetchData(BANNED_IPS));
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error('error fetching banned ips', e);
    }
  };

  // A user action (ban, moderator toggle, delete) can move a user between the
  // paginated view and either list, so refresh all three.
  const refresh = () => {
    getUsers();
    getLists();
  };

  useEffect(() => {
    getUsers();
  }, [currentPage, search, status, sort]);

  useEffect(() => {
    getLists();
  }, []);

  useEffect(() => {
    getLiveInfo();
    const id = setInterval(getLiveInfo, LIVE_FETCH_INTERVAL);
    return () => clearInterval(id);
  }, [online]);

  const columns: ColumnsType<User> = [
    {
      title: 'Display Name',
      dataIndex: 'displayName',
      key: 'displayName',
      render: (displayName: string) => <span className="display-name">{displayName}</span>,
    },
    {
      title: 'Status',
      key: 'status',
      render: (_, user) => (
        <Space size={[0, 4]} wrap>
          {user.disabledAt ? <Tag color="red">Banned</Tag> : <Tag color="green">Active</Tag>}
          {user.scopes?.includes('MODERATOR') && <Tag color="blue">Moderator</Tag>}
          {user.isBot && <Tag>Bot</Tag>}
        </Space>
      ),
    },
    {
      title: 'Authentication',
      key: 'authentication',
      render: (_, user) => {
        if (user.authProviders?.length) {
          return (
            <Space size={[0, 4]} wrap>
              {user.authProviders.map(provider => (
                <Tag color="gold" key={provider} title={`Authenticated with ${provider}`}>
                  {provider}
                </Tag>
              ))}
            </Space>
          );
        }
        if (user.authenticated) {
          return <Tag color="gold">Authenticated</Tag>;
        }
        return <Typography.Text type="secondary">Anonymous</Typography.Text>;
      },
    },
    {
      title: 'Created',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: Date) => formatDateOnly(date),
      sorter: (a: User, b: User) =>
        new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: '',
      key: 'actions',
      className: 'actions-col',
      // Stop clicks in the actions cell from also opening the row's detail modal.
      render: (_, user) => (
        <span
          onClick={e => e.stopPropagation()}
          onKeyDown={e => e.stopPropagation()}
          role="presentation"
        >
          <UserActionsDropdown user={user} onChanged={refresh} />
        </span>
      ),
    },
  ];

  // Clicking a row opens that user's detail modal. Shared by the paginated "All
  // Users" table and the Moderators/Banned list tables.
  const rowInteraction = (record: User) => ({
    onClick: () => setSelectedUser(record),
    style: { cursor: 'pointer' },
  });

  // The "All Users" table is server-paginated, so its Created column sorts via
  // the API across every page rather than only the rows in hand. The dedicated
  // Moderators/Banned tables hold complete sets, so they keep the shared
  // client-side sorter from `columns`.
  const usersColumns: ColumnsType<User> = columns.map(col =>
    col.key === 'createdAt'
      ? { ...col, sorter: true, sortOrder: (sort === 'asc' ? 'ascend' : 'descend') as SortOrder }
      : col,
  );

  const allUsersTab = (
    <>
      <Typography.Paragraph>
        <Translation
          translationKey={Localization.Admin.Users.pageDescription}
          defaultText="Everyone known to your server: chat viewers, authenticated and plugin-created users, and API integrations. Change moderator access, ban, or permanently remove users."
        />
      </Typography.Paragraph>

      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search
          placeholder="Search by display name"
          allowClear
          enterButton
          onSearch={value => {
            setCurrentPage(1);
            setSearch(value);
          }}
          style={{ width: 320 }}
        />
        <Select
          value={status}
          onChange={value => {
            setCurrentPage(1);
            setStatus(value);
          }}
          options={STATUS_OPTIONS}
          style={{ width: 160 }}
        />
      </Space>

      <Table
        loading={loading}
        className="table-container"
        columns={usersColumns}
        dataSource={users}
        size="small"
        rowKey="id"
        onRow={rowInteraction}
        pagination={{
          current: currentPage,
          pageSize: PAGE_SIZE,
          total,
          showSizeChanger: false,
          hideOnSinglePage: true,
        }}
        onChange={(pagination, _filters, sorter) => {
          // A sort change resets to the first page; a page change keeps the
          // current sort.
          const activeSorter = Array.isArray(sorter) ? sorter[0] : sorter;
          if (activeSorter?.columnKey === 'createdAt') {
            const nextSort = activeSorter.order === 'ascend' ? 'asc' : 'desc';
            if (nextSort !== sort) {
              setSort(nextSort);
              setCurrentPage(1);
              return;
            }
          }
          setCurrentPage(pagination.current || 1);
        }}
      />
    </>
  );

  const connectedTab = online ? (
    <>
      <ClientTable data={clients} />
      <p className="description">
        {t('Visit the')}{' '}
        <a
          href="https://owncast.online/docs/viewers/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          {t('documentation')}
        </a>{' '}
        {t('to configure additional details about your viewers.')}
      </p>
    </>
  ) : (
    <p className="description">
      {t(
        'When a stream is active and chat is enabled, connected chat clients will be displayed here.',
      )}
    </p>
  );

  // Full, unpaginated tables for the dedicated lists.
  const listTable = (data: User[]) => (
    <Table
      className="table-container"
      columns={columns}
      dataSource={data}
      size="small"
      rowKey="id"
      onRow={rowInteraction}
      pagination={{ pageSize: PAGE_SIZE, showSizeChanger: false, hideOnSinglePage: true }}
    />
  );

  const items = [
    { label: t('All Users'), key: 'all', children: allUsersTab },
    {
      label: <span>{`${t('Moderators')} (${moderators.length})`}</span>,
      key: 'moderators',
      children: listTable(moderators),
    },
    {
      label: <span>{`${t('Banned')} (${bannedUsers.length})`}</span>,
      key: 'banned',
      children: listTable(bannedUsers),
    },
    {
      label: <span>{`${t('Connected')} (${online ? clients.length : t('offline')})`}</span>,
      key: 'connected',
      children: connectedTab,
    },
    {
      label: <span>{`${t('IP Bans')} (${ipBans.length})`}</span>,
      key: 'ipbans',
      children: <BannedIPsTable data={ipBans} />,
    },
  ];

  return (
    <div>
      <Typography.Title level={3}>
        <Translation translationKey={Localization.Admin.Users.pageTitle} defaultText="Users" />
      </Typography.Title>
      <Tabs defaultActiveKey="all" items={items} />
      {selectedUser && (
        <UserDetailsModal
          user={selectedUser}
          open
          onClose={() => {
            setSelectedUser(null);
            // A ban or moderator change made from the modal should be reflected
            // in the lists.
            refresh();
          }}
        />
      )}
    </div>
  );
}

UsersAdmin.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
