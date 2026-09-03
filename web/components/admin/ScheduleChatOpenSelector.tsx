import { Select } from 'antd';
import type { FC } from 'react';

export interface ScheduleChatOpenSelectorProps {
  value: number;
  options: Array<{ value: number; label: string }>;
  label: string;
  tip: string;
  onChange: (value: number) => void;
}

export const ScheduleChatOpenSelector: FC<ScheduleChatOpenSelectorProps> = ({
  value,
  options,
  label,
  tip,
  onChange,
}) => (
  <div className="formfield-container toggleswitch-container">
    <div className="label-side">
      <span className="formfield-label">{label}</span>
    </div>
    <div className="input-side">
      <div className="input-group">
        <Select<number>
          aria-label={label}
          value={value}
          options={options}
          onChange={onChange}
          style={{ width: '100%' }}
        />
      </div>
      <p className="field-tip">{tip}</p>
    </div>
  </div>
);
