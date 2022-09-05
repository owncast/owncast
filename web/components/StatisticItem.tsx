/* eslint-disable react/no-unused-prop-types */
/* eslint-disable react/no-unstable-nested-components */
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

interface ContentProps {
  prefix: string;
  value: any;
  suffix: string;
  title: string;
}

const Content = ({ prefix, value, suffix, title }: ContentProps) => (
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

const ProgressView = ({ title, value, prefix, suffix, color }: StatisticItemProps) => {
  const endColor = value > 90 ? 'red' : color;
  const content = <Content prefix={prefix} value={value} suffix={suffix} title={title} />;

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
};
ProgressView.defaultProps = defaultProps;

const StatisticView = ({ title, value, prefix, formatter }: StatisticItemProps) => (
  <Statistic title={title} value={value} prefix={prefix} formatter={formatter} />
);
StatisticView.defaultProps = defaultProps;

const StatisticItem = (props: StatisticItemProps) => {
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
};
StatisticItem.defaultProps = defaultProps;

export default StatisticItem;
