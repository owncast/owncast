import { InfoCircleOutlined } from '@ant-design/icons';
import { Tooltip } from 'antd';

interface InfoTipProps {
  tip: string | null;
}

export default function InfoTip({ tip }: InfoTipProps) {
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
}
