import { Table } from 'antd';
import format from 'date-fns/format';
import { SortOrder } from 'antd/lib/table/interface';
import { formatDistanceToNow } from 'date-fns';
import { User } from '../types/chat';
import { formatUAstring } from '../utils/format';

export function formatDisplayDate(date: string | Date) {
  return format(new Date(date), 'MMM d H:mma');
}
export default function ViewerTable({ data }: ViewerTableProps) {
  const columns = [
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
      title: 'Watch Time',
      dataIndex: 'firstSeen',
      key: 'firstSeen',
      defaultSortOrder: 'ascend' as SortOrder,
      render: (time: Date) => formatDistanceToNow(new Date(time)),
      sorter: (a: any, b: any) => new Date(a.firstSeen).getTime() - new Date(b.firstSeen).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
  ];

  return (
    <Table
      pagination={{ hideOnSinglePage: true }}
      className="table-container"
      columns={columns}
      dataSource={data}
      size="small"
      rowKey="id"
    />
  );
}

interface ViewerTableProps {
  data: User[];
}
