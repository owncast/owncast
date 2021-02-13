import React from 'react';
import { Select } from 'antd';
import { SocialHandleDropdownItem } from '../../types/config-section';
import { NEXT_PUBLIC_API_HOST } from '../../utils/apis';
import { OTHER_SOCIAL_HANDLE_OPTION } from '../../utils/config-constants';

interface DropdownProps {
  iconList: SocialHandleDropdownItem[];
  selectedOption: string;
  onSelected: any;
}

export default function SocialDropdown({ iconList, selectedOption, onSelected }: DropdownProps) {
  const handleSelected = (value: string) => {
    if (onSelected) {
      onSelected(value);
    }
  };
  const inititalSelected = selectedOption === '' ? null : selectedOption;
  return (
    <div className="social-dropdown-container">
      <p className="description">
        If you are looking for a platform name not on this list, please select Other and type in
        your own name. A logo will not be provided.
      </p>
      <p className="description">
        If you DO have a logo, drop it in to the <code>/webroot/img/platformicons</code> directory
        and update the <code>/socialHandle.go</code> list. Then restart the server and it will show
        up in the list.
      </p>

      <Select
        style={{ width: 240 }}
        className="social-dropdown"
        placeholder="Social platform..."
        defaultValue={inititalSelected}
        value={inititalSelected}
        onSelect={handleSelected}
      >
        {iconList.map(item => {
          const { platform, icon, key } = item;
          return (
            <Select.Option className="social-option" key={`platform-${key}`} value={key}>
              <span className="option-icon">
                <img src={`${NEXT_PUBLIC_API_HOST}${icon}`} alt="" className="option-icon" />
              </span>
              <span className="option-label">{platform}</span>
            </Select.Option>
          );
        })}
        <Select.Option
          className="social-option"
          key={`platform-${OTHER_SOCIAL_HANDLE_OPTION}`}
          value={OTHER_SOCIAL_HANDLE_OPTION}
        >
          Other...
        </Select.Option>
      </Select>
    </div>
  );
}
