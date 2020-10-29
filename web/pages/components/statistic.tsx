import { Statistic, Card, Col} from "antd";

interface ItemProps {
  title: string, 
  value: string,
  prefix: JSX.Element,
};

export default function StatisticItem(props: ItemProps) {
  const { title, value, prefix } = props;
  const valueStyle = { color: "#334", fontSize: "1.8rem" };

  return (
    <Col span={8}>
      <Card>
        <Statistic
          title={title}
          value={value}
          valueStyle={valueStyle}
          prefix={prefix}
        />
      </Card>
    </Col>
  );
}