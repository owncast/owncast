import { Popconfirm, Button, Typography } from 'antd';
import { useContext } from 'react';
import { AlertMessageContext } from '../../utils/alert-message-context';

import { API_YP_RESET, fetchData } from '../../utils/apis';

export default function ResetYP() {
  const { setMessage } = useContext(AlertMessageContext);

  const { Title } = Typography;

  const resetDirectoryRegistration = async () => {
    try {
      await fetchData(API_YP_RESET);
      setMessage('');
    } catch (error) {
      alert(error);
    }
  };
  return (
    <>
      <Title level={3} className="section-title">
        Reset Directory
      </Title>
      <p className="description">
        If you are experiencing issues with your listing on the Owncast Directory and were asked to
        &quot;reset&quot; your connection to the service, you can do that here. The next time you go
        live it will try and re-register your server with the directory from scratch.
      </p>

      <Popconfirm
        placement="topLeft"
        title="Are you sure you want to reset your connection to the Owncast directory?"
        onConfirm={resetDirectoryRegistration}
        okText="Yes"
        cancelText="No"
      >
        <Button type="primary">Reset Directory Connection</Button>
      </Popconfirm>
    </>
  );
}
