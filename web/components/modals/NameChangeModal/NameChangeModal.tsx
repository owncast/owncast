import React, { CSSProperties, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { Input, Button, Select } from 'antd';
import { MessageType } from '../../../interfaces/socket-events';
import WebsocketService from '../../../services/websocket-service';
import {
  websocketServiceAtom,
  chatDisplayNameAtom,
  chatDisplayColorAtom,
} from '../../stores/ClientConfigStore';

const { Option } = Select;

/* eslint-disable @typescript-eslint/no-unused-vars */
interface Props {}

const UserColor = (props: { color: number }): React.ReactElement => {
  const { color } = props;
  const style: CSSProperties = {
    textAlign: 'center',
    backgroundColor: `var(--theme-color-users-${color})`,
    width: '100%',
    height: '100%',
  };
  return <div style={style} />;
};

const NameChangeModal = (props: Props) => {
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const chatDisplayColor = useRecoilValue<number>(chatDisplayColorAtom) || 0;
  const [newName, setNewName] = useState<any>(chatDisplayName);

  const handleNameChange = () => {
    const nameChange = {
      type: MessageType.NAME_CHANGE,
      newName,
    };
    websocketService.send(nameChange);
  };

  const saveEnabled =
    newName !== chatDisplayName && newName !== '' && websocketService?.isConnected();

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
        value={newName}
        onChange={e => setNewName(e.target.value)}
        placeholder="Your chat display name"
        maxLength={30}
        showCount
        defaultValue={chatDisplayName}
      />
      <Button disabled={!saveEnabled} onClick={handleNameChange}>
        Change name
      </Button>
      <div>
        Your Color
        <Select
          style={{ width: 120 }}
          onChange={handleColorChange}
          defaultValue={chatDisplayColor.toString()}
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
export default NameChangeModal;
