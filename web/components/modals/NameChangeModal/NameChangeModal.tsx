import React, { CSSProperties, FC, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { Input, Button, Select } from 'antd';
import { MessageType } from '../../../interfaces/socket-events';
import WebsocketService from '../../../services/websocket-service';
import { websocketServiceAtom, currentUserAtom } from '../../stores/ClientConfigStore';

const { Option } = Select;

export type UserColorProps = {
  color: number;
};

const UserColor: FC<UserColorProps> = ({ color }) => {
  const style: CSSProperties = {
    textAlign: 'center',
    backgroundColor: `var(--theme-color-users-${color})`,
    width: '100%',
    height: '100%',
  };
  return <div style={style} />;
};

export const NameChangeModal: FC = () => {
  const currentUser = useRecoilValue(currentUserAtom);
  if (!currentUser) {
    return null;
  }

  const { displayName, displayColor } = currentUser;
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const [newName, setNewName] = useState<any>(displayName);

  const handleNameChange = () => {
    const nameChange = {
      type: MessageType.NAME_CHANGE,
      newName,
    };
    websocketService.send(nameChange);
  };

  const saveEnabled = newName !== displayName && newName !== '' && websocketService?.isConnected();

  const handleColorChange = (color: string) => {
    const colorChange = {
      type: MessageType.COLOR_CHANGE,
      newColor: Number(color),
    };
    websocketService.send(colorChange);
  };

  const maxColor = 8; // 0...n
  const colorOptions = [...Array(maxColor)].map((e, i) => i);

  return (
    <div>
      Your chat display name is what people see when you send chat messages. Other information can
      go here to mention auth, and stuff.
      <Input
        id="name-change-field"
        value={newName}
        onChange={e => setNewName(e.target.value)}
        placeholder="Your chat display name"
        maxLength={30}
        showCount
        defaultValue={displayName}
      />
      <Button id="name-change-submit" disabled={!saveEnabled} onClick={handleNameChange}>
        Change name
      </Button>
      <div>
        Your Color
        <Select
          style={{ width: 120 }}
          onChange={handleColorChange}
          defaultValue={displayColor.toString()}
          getPopupContainer={triggerNode => triggerNode.parentElement}
        >
          {colorOptions.map(e => (
            <Option key={e.toString()} title={e}>
              <UserColor color={e} />
            </Option>
          ))}
        </Select>
      </div>
    </div>
  );
};
