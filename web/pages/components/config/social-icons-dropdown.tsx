import React, { useState } from 'react';
import { PlusOutlined } from "@ant-design/icons";
import { Select, Divider, Input } from "antd";
import classNames from 'classnames';
import { SocialHandleDropdownItem } from "../../../types/config-section";
import { NEXT_PUBLIC_API_HOST } from '../../../utils/apis';


interface DropdownProps {
  iconList: SocialHandleDropdownItem[];
  selectedOption?: string;
}
interface DropdownOptionProps extends SocialHandleDropdownItem {
  isSelected: boolean;
}

// Add "Other" item which creates a text field
// Add fixed custom ones - "url", "donate", "follow", "rss"

function dropdownRender(menu) {
  console.log({menu})
  return 'hi';
}

export default function SocialDropdown({ iconList, selectedOption }: DropdownProps) {
  const [name, setName] = useState('');

  const handleNameChange = event => {
    setName(event.target.value);
  };

  const handleAddItem = () => {
    console.log('addItem');
    // const { items, name } = this.state;
    // this.setState({
    //   items: [...items, name || `New item ${index++}`],
    //   name: '',
    // });
  };


  return (
    <div className="social-dropdown-container">
      <Select
        style={{ width: 240 }}
        className="social-dropdown"
        placeholder="Social platform..."
        // defaultValue
        dropdownRender={menu => (
          <>
            {menu}
            <Divider style={{ margin: '4px 0' }} />
            <div style={{ display: 'flex', flexWrap: 'nowrap', padding: 8 }}>
              <Input style={{ flex: 'auto' }} value="" onChange={handleNameChange} />
              <a
                style={{ flex: 'none', padding: '8px', display: 'block', cursor: 'pointer' }}
                onClick={handleAddItem}
              >
                <PlusOutlined /> Add item
              </a>
            </div>
          </>
        )}
      >
        {iconList.map(item => {
          const { platform, icon, key } =  item;
          return (
            <Select.Option className="social-option" key={`platform-${key}`} value={key}>
              <span className="option-icon">
                <img src={`${NEXT_PUBLIC_API_HOST}${icon}`} alt="" className="option-icon" />
              </span>
              <span className="option-label">{platform}</span>
            </Select.Option>
          );
        })
      }
      </Select>

    </div>
  );
}
