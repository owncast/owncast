import { Table, Typography } from "antd";

const { Title } = Typography;

export default function KeyValueTable({ title, data }) {
  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Value",
      dataIndex: "value",
      key: "value",
    },
  ];

  return (
    <div>
      <Title>{title}</Title>
      <Table pagination={false} columns={columns} dataSource={data} />
    </div>
  );
}
