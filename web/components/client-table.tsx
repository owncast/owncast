import { Input, Table } from 'antd';
import { FilterDropdownProps, SortOrder } from 'antd/lib/table/interface';
import { ColumnsType } from 'antd/es/table';
import { SearchOutlined } from '@ant-design/icons';
import { formatDistanceToNow } from 'date-fns';
import { Client } from '../types/chat';
import UserPopover from './user-popover';
import BanUserButton from './other/ban-user-button/ban-user-button';
import { formatUAstring } from '../utils/format';

export default function ClientTable({ data }: ClientTableProps) {
  const columns: ColumnsType<Client> = [
    {
      title: 'Display Name',
      key: 'username',
      // eslint-disable-next-line react/destructuring-assignment
      render: (client: Client) => {
        const { user, connectedAt, messageCount, userAgent } = client;
        const connectionInfo = { connectedAt, messageCount, userAgent };
        return (
          <UserPopover user={user} connectionInfo={connectionInfo}>
            <span className="display-name">{user.displayName}</span>
          </UserPopover>
        );
      },
      sorter: (a: any, b: any) => b.user.displayName.localeCompare(a.user.displayName),
      filterIcon: <SearchOutlined />,
      // eslint-disable-next-line react/no-unstable-nested-components
      filterDropdown: ({ setSelectedKeys, selectedKeys, confirm }: FilterDropdownProps) => (
        <div style={{ padding: 8 }}>
          <Input
            placeholder="Search display names..."
            value={selectedKeys[0]}
            onChange={e => {
              setSelectedKeys(e.target.value ? [e.target.value] : []);
              confirm({ closeDropdown: false });
            }}
          />
        </div>
      ),
      onFilter: (value: string, record: Client) => record.user.displayName.includes(value),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: 'Messages sent',
      dataIndex: 'messageCount',
      key: 'messageCount',
      className: 'number-col',
      width: '12%',
      sorter: (a: any, b: any) => a.messageCount - b.messageCount,
      sortDirections: ['descend', 'ascend'] as SortOrder[],
      render: (count: number) => <div style={{ textAlign: 'center' }}>{count}</div>,
    },
    {
      title: 'Connected Time',
      dataIndex: 'connectedAt',
      key: 'connectedAt',
      defaultSortOrder: 'ascend',
      render: (time: Date) => formatDistanceToNow(new Date(time)),
      sorter: (a: any, b: any) =>
        new Date(b.connectedAt).getTime() - new Date(a.connectedAt).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: 'User Agent',
      dataIndex: 'userAgent',
      key: 'userAgent',
      render: (ua: string) => formatUAstring(ua),
    },
    {
      title: 'Location',
      dataIndex: 'geo',
      key: 'geo',
      render: geo => (geo ? `${geo.regionName}, ${geo.countryCode}` : '-'),
    },
    {
      title: '',
      key: 'block',
      className: 'actions-col',
      render: (_, row) => <BanUserButton user={row.user} isEnabled={!row.user.disabledAt} />,
    },
  ];

  return (
    <Table
      className="table-container"
      pagination={{ hideOnSinglePage: true }}
      columns={columns}
      dataSource={data}
      size="small"
      rowKey="id"
    />
  );
}

interface ClientTableProps {
  data: Client[];
}
