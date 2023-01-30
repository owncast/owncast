import { Tooltip } from 'antd';
import dynamic from 'next/dynamic';
import { FC } from 'react';

// Lazy loaded components

const InfoCircleOutlined = dynamic(() => import('@ant-design/icons/InfoCircleOutlined'), {
  ssr: false,
});

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
