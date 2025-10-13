import React, { FC, useState } from 'react';
import { Table, Tag, Typography } from 'antd';
import Linkify from 'react-linkify';
import { SortOrder, TablePaginationConfig } from 'antd/lib/table/interface';
import { format } from 'date-fns';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../types/localization';

const { Title } = Typography;

function renderColumnLevel(text, entry) {
  let color = 'black';

  if (entry.level === 'warning') {
    color = 'orange';
  } else if (entry.level === 'error') {
    color = 'red';
  }

  return <Tag color={color}>{text}</Tag>;
}

export type LogTableProps = {
  logs: object[];
  initialPageSize: number;
};

export const LogTable: FC<LogTableProps> = ({ logs, initialPageSize }) => {
  const { t } = useTranslation();
  const [pageSize, setPageSize] = useState(initialPageSize);

  const handleTableChange = (pagination: TablePaginationConfig) => {
    setPageSize(pagination.pageSize);
  };

  if (!logs?.length) {
    return null;
  }

  const columns = [
    {
      title: t(Localization.Admin.LogTable.level),
      dataIndex: 'level',
      key: 'level',
      filters: [
        {
          text: t(Localization.Admin.LogTable.info),
          value: 'info',
        },
        {
          text: t(Localization.Admin.LogTable.warning),
          value: 'warning',
        },
        {
          text: t(Localization.Admin.LogTable.error),
          value: 'error',
        },
      ],
      onFilter: (level, row) => row.level === level,
      render: renderColumnLevel,
    },
    {
      title: t(Localization.Admin.LogTable.timestamp),
      dataIndex: 'time',
      key: 'time',
      render: (timestamp: Date) => {
        const dateObject = new Date(timestamp);
        return format(dateObject, 'p P');
      },
      sorter: (a: any, b: any) => new Date(a.time).getTime() - new Date(b.time).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
      defaultSortOrder: 'descend' as SortOrder,
    },

    {
      title: t(Localization.Admin.LogTable.message),
      dataIndex: 'message',
      key: 'message',
      render: (message: string) => <Linkify>{message}</Linkify>,
    },
  ];

  return (
    <div className="logs-section">
      <Title>{t(Localization.Admin.LogTable.logs)}</Title>
      <Table
        size="middle"
        dataSource={logs}
        columns={columns}
        rowKey={row => row.time}
        pagination={{
          pageSize,
        }}
        onChange={handleTableChange}
      />
    </div>
  );
};
