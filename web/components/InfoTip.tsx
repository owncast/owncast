import { InfoCircleOutlined } from '@ant-design/icons';
import { Tooltip } from 'antd';
import { FC } from 'react';

export type InfoTipProps = {
  tip: string | null;
};

export const InfoTip: FC<InfoTipProps> = ({ tip }) => {
  if (tip === '' || tip === null) {
    return null;
  }

  return (
    <span className="info-tip">
      <Tooltip title={tip}>
        <InfoCircleOutlined />
      </Tooltip>
    </span>
  );
};
