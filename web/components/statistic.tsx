/* eslint-disable react/no-unused-prop-types */
// TODO: This component should be cleaned up and usage should be re-examined. The types should be reconsidered as well.

import { Typography, Statistic, Card, Progress } from 'antd';

const { Text } = Typography;

interface StatisticItemProps {
  title?: string;
  value?: any;
  prefix?: any;
  suffix?: string;
  color?: string;
  progress?: boolean;
  centered?: boolean;
  formatter?: any;
}
const defaultProps = {
  title: '',
  value: 0,
  prefix: null,
  suffix: null,
  color: '',
  progress: false,
  centered: false,
  formatter: null,
};

function ProgressView({ title, value, prefix, suffix, color }: StatisticItemProps) {
  const endColor = value > 90 ? 'red' : color;
  const content = (
    <div>
      {prefix}
      <div>
        <Text type="secondary">{title}</Text>
      </div>
      <div>
        <Text type="secondary">
          {value}
          {suffix || '%'}
        </Text>
      </div>
    </div>
  );
  return (
    <Progress
      type="dashboard"
      percent={value}
      width={120}
      strokeColor={{
        '0%': color,
        '90%': endColor,
      }}
      format={() => content}
    />
  );
}
ProgressView.defaultProps = defaultProps;

function StatisticView({ title, value, prefix, formatter }: StatisticItemProps) {
  return <Statistic title={title} value={value} prefix={prefix} formatter={formatter} />;
}
StatisticView.defaultProps = defaultProps;

export default function StatisticItem(props: StatisticItemProps) {
  const { progress, centered } = props;
  const View = progress ? ProgressView : StatisticView;

  const style = centered ? { display: 'flex', alignItems: 'center', justifyContent: 'center' } : {};

  return (
    <Card type="inner">
      <div style={style}>
        <View {...props} />
      </div>
    </Card>
  );
}
StatisticItem.defaultProps = defaultProps;
