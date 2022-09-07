import { Table, Typography } from 'antd';
import { FC } from 'react';

const { Title } = Typography;

export type KeyValueTableProps = {
  title: string;
  data: any;
};

export const KeyValueTable: FC<KeyValueTableProps> = ({ title, data }) => {
  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Value',
      dataIndex: 'value',
      key: 'value',
    },
  ];

  return (
    <>
      <Title level={2}>{title}</Title>
      <Table pagination={false} columns={columns} dataSource={data} rowKey="name" />
    </>
  );
};
