import { Popconfirm, Button, Typography } from 'antd';
import { useContext, useState } from 'react';
import { AlertMessageContext } from '../../utils/alert-message-context';

import { API_YP_RESET, fetchData } from '../../utils/apis';
import { RESET_TIMEOUT } from '../../utils/config-constants';
import {
  createInputStatus,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import FormStatusIndicator from './form-status-indicator';

export default function ResetYP() {
  const { setMessage } = useContext(AlertMessageContext);

  const [submitStatus, setSubmitStatus] = useState(null);
  let resetTimer = null;
  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  const resetDirectoryRegistration = async () => {
    setSubmitStatus(createInputStatus(STATUS_PROCESSING));
    try {
      await fetchData(API_YP_RESET);
      setMessage('');
      setSubmitStatus(createInputStatus(STATUS_SUCCESS));
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
    } catch (error) {
      setSubmitStatus(createInputStatus(STATUS_ERROR, `There was an error: ${error}`));
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
    }
  };

  return (
    <>
      <Typography.Title level={3} className="section-title">
        Reset Directory
      </Typography.Title>
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
      <p>
        <FormStatusIndicator status={submitStatus} />
      </p>
    </>
  );
}
