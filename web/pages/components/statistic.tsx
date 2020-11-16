import { Typography, Statistic, Card, Col, Progress} from "antd";
const { Text } = Typography;

interface StatisticItemProps {
  title?: string, 
  value?: any,
  prefix?: JSX.Element,
  color?: string,
  progress?: boolean,
  centered?: boolean,
};
const defaultProps = {
  title: '', 
  value: 0,
  prefix: null,
  color: '',
  progress: false,
  centered: false,
};


function ProgressView({ title, value, prefix, color }: StatisticItemProps) {
  const endColor = value > 90 ? 'red' : color;
  const content = (
    <div>
    {prefix}
    <div><Text type="secondary">{title}</Text></div>
    <div><Text type="secondary">{value}%</Text></div>
    </div>
  )
  return (
    <Progress
      type="dashboard"
      percent={value}
      width={120}
      strokeColor={{
        '0%': color,
        '90%': endColor,
      }}
      format={percent => content}
    />
  )
}
ProgressView.defaultProps = defaultProps;

function StatisticView({ title, value, prefix }: StatisticItemProps) {
  return (
    <Statistic
      title={title}
      value={value}
      prefix={prefix}
  />
  )
}
StatisticView.defaultProps = defaultProps;

export default function StatisticItem(props: StatisticItemProps) {
  const { progress, centered } = props;
  const View = progress ? ProgressView : StatisticView;

  const style = centered ? {display: 'flex', alignItems: 'center', justifyContent: 'center'} : {};

  return (
      <Card type="inner">
        <div style={style}>
          <View {...props} />
        </div>
      </Card>
  );
}
StatisticItem.defaultProps = defaultProps;
