import { Table } from 'antd';
import { SortOrder } from 'antd/lib/table/interface';
import { ColumnsType } from 'antd/es/table';
import { formatDistanceToNow } from 'date-fns';
import { Client } from '../types/chat';
import UserPopover from './user-popover';
import BanUserButton from './ban-user-button';
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
      sorter: (a: any, b: any) => a.user.displayName - b.user.displayName,
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: 'Messages sent',
      dataIndex: 'messageCount',
      key: 'messageCount',
      className: 'number-col',
      sorter: (a: any, b: any) => a.messageCount - b.messageCount,
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: 'Connected Time',
      dataIndex: 'connectedAt',
      key: 'connectedAt',
      defaultSortOrder: 'ascend',
      render: (time: Date) => formatDistanceToNow(new Date(time)),
      sorter: (a: any, b: any) =>
        new Date(a.connectedAt).getTime() - new Date(b.connectedAt).getTime(),
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
