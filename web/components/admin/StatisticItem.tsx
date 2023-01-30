/* eslint-disable react/no-unused-prop-types */
/* eslint-disable react/no-unstable-nested-components */
// TODO: This component should be cleaned up and usage should be re-examined. The types should be reconsidered as well.

import { Typography, Statistic, Card, Progress } from 'antd';
import { FC } from 'react';

const { Text } = Typography;

export type StatisticItemProps = {
  title?: string;
  value?: any;
  prefix?: any;
  suffix?: string;
  color?: string;
  progress?: boolean;
  centered?: boolean;
  formatter?: any;
};

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

export type ContentProps = {
  prefix: string;
  value: any;
  suffix: string;
  title: string;
};

const Content: FC<ContentProps> = ({ prefix, value, suffix, title }) => (
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

const ProgressView: FC<StatisticItemProps> = ({ title, value, prefix, suffix, color }) => {
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

const StatisticView: FC<StatisticItemProps> = ({ title, value, prefix, formatter }) => (
  <Statistic title={title} value={value} prefix={prefix} formatter={formatter} />
);
StatisticView.defaultProps = defaultProps;

export const StatisticItem: FC<StatisticItemProps> = props => {
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
