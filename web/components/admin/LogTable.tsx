import React, { FC, useState } from 'react';
import { Table, Tag, Typography } from 'antd';
import Linkify from 'react-linkify';
import { SortOrder, TablePaginationConfig } from 'antd/lib/table/interface';
import format from 'date-fns/format';

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

function renderMessage(text) {
  return <Linkify>{text}</Linkify>;
}

export type LogTableProps = {
  logs: object[];
  initialPageSize: number;
};

export const LogTable: FC<LogTableProps> = ({ logs, initialPageSize }) => {
  if (!logs?.length) {
    return null;
  }

  const [pageSize, setPageSize] = useState(initialPageSize);
  const handleTableChange = (pagination: TablePaginationConfig) => {
    setPageSize(pagination.pageSize);
  };

  const columns = [
    {
      title: 'Level',
      dataIndex: 'level',
      key: 'level',
      filters: [
        {
          text: 'Info',
          value: 'info',
        },
        {
          text: 'Warning',
          value: 'warning',
        },
        {
          text: 'Error',
          value: 'Error',
        },
      ],
      onFilter: (level, row) => row.level.indexOf(level) === 0,
      render: renderColumnLevel,
    },
    {
      title: 'Timestamp',
      dataIndex: 'time',
      key: 'time',
      render: timestamp => {
        const dateObject = new Date(timestamp);
        return format(dateObject, 'pp P');
      },
      sorter: (a, b) => new Date(a.time).getTime() - new Date(b.time).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
      defaultSortOrder: 'descend' as SortOrder,
    },
    {
      title: 'Message',
      dataIndex: 'message',
      key: 'message',
      render: renderMessage,
    },
  ];

  return (
    <div className="logs-section">
      <Title>Logs</Title>
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
