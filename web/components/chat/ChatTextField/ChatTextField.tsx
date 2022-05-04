import { SmileOutlined } from '@ant-design/icons';
import { Button, Input } from 'antd';
import { useRef, useState } from 'react';
import { useRecoilValue } from 'recoil';
import ContentEditable from 'react-contenteditable';
import WebsocketService from '../../../services/websocket-service';
import { websocketServiceAtom } from '../../stores/ClientConfigStore';
import s from './ChatTextField.module.scss';

interface Props {
  value?: string;
}

export default function ChatTextField(props: Props) {
  const { value: originalValue } = props;
  const [value, setValue] = useState(originalValue);
  const [showEmojis, setShowEmojis] = useState(false);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);

  const text = useRef(value);

  // large is 40px
  const size = 'small';

  const handleChange = evt => {
    text.current = evt.target.value;
  };

  return (
    <div>
      <Input.Group compact style={{ display: 'flex', width: '100%', position: 'absolute' }}>
        <ContentEditable
          style={{ width: '60%', maxHeight: '50px', padding: '5px' }}
          html={text.current}
          onChange={handleChange}
        />
        <Button type="default" ghost title="Emoji" onClick={() => setShowEmojis(!showEmojis)}>
          <SmileOutlined style={{ color: 'rgba(0,0,0,.45)' }} />
        </Button>
        <Button size={size} type="primary">
          Submit
        </Button>
      </Input.Group>
    </div>
  );
}

ChatTextField.defaultProps = {
  value: '',
};
