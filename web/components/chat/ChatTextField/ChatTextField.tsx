import { MoreOutlined } from '@ant-design/icons';
import { Input, Button } from 'antd';
import { useEffect, useState } from 'react';
import s from './ChatTextField.module.scss';

interface Props {}

export default function ChatTextField(props: Props) {
  const [value, setValue] = useState('');
  const [showEmojis, setShowEmojis] = useState(false);
  const size = 'large';

  useEffect(() => {
    console.log({ value });
  }, [value]);

  return (
    <div>
      <Input.Group compact style={{ display: 'flex' }}>
        <Input
          onChange={e => setValue(e.target.value)}
          size={size}
          placeholder="Enter text and hit enter!"
        />
        <Button
          size={size}
          icon={<MoreOutlined />}
          type="default"
          onClick={() => setShowEmojis(!showEmojis)}
        />
        <Button size={size} type="primary">
          Submit
        </Button>
      </Input.Group>
    </div>
  );
}
