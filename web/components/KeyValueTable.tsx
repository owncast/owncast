import { Table, Typography } from 'antd';

const { Title } = Typography;

export const KeyValueTable = ({ title, data }: KeyValueTableProps) => {
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
export default KeyValueTable;

interface KeyValueTableProps {
  title: string;
  data: any;
}
