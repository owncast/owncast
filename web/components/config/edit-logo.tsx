import { Button, Upload } from 'antd';
import { RcFile } from 'antd/lib/upload/interface';
import { LoadingOutlined, UploadOutlined } from '@ant-design/icons';
import React, { useState, useContext } from 'react';
import FormStatusIndicator from './form-status-indicator';
import { ServerStatusContext } from '../../utils/server-status-context';
import {
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  TEXTFIELD_PROPS_LOGO,
} from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { NEXT_PUBLIC_API_HOST } from '../../utils/apis';

const ACCEPTED_FILE_TYPES = ['image/svg+xml', 'image/png', 'image/jpeg', 'image/gif'];

function getBase64(img: File | Blob, callback: (imageUrl: string | ArrayBuffer) => void) {
  const reader = new FileReader();
  reader.addEventListener('load', () => callback(reader.result));
  reader.readAsDataURL(img);
}

export default function EditLogo() {
  const [logoUrl, setlogoUrl] = useState(null);
  const [loading, setLoading] = useState(false);
  const [logoCachedbuster, setLogoCacheBuster] = useState(0);

  const serverStatusData = useContext(ServerStatusContext);
  const { setFieldInConfigState, serverConfig } = serverStatusData || {};
  const currentLogo = serverConfig?.instanceDetails?.logo;

  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  let resetTimer = null;

  const { apiPath, tip } = TEXTFIELD_PROPS_LOGO;

  // Clear out any validation states and messaging
  const resetStates = () => {
    setSubmitStatus(null);
    clearTimeout(resetTimer);
    resetTimer = null;
  };

  // validate file type and create base64 encoded img
  const beforeUpload = (file: RcFile) => {
    setLoading(true);

    // eslint-disable-next-line consistent-return
    return new Promise<void>((res, rej) => {
      if (!ACCEPTED_FILE_TYPES.includes(file.type)) {
        const msg = `File type is not supported: ${file.type}`;
        setSubmitStatus(createInputStatus(STATUS_ERROR, `There was an error: ${msg}`));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
        setLoading(false);
        return rej();
      }

      getBase64(file, (url: string) => {
        setlogoUrl(url);
        return res();
      });
    });
  };

  // Post new logo to api
  const handleLogoUpdate = async () => {
    if (logoUrl !== currentLogo) {
      setSubmitStatus(createInputStatus(STATUS_PROCESSING));

      await postConfigUpdateToAPI({
        apiPath,
        data: { value: logoUrl },
        onSuccess: () => {
          setFieldInConfigState({ fieldName: 'logo', value: logoUrl, path: '' });
          setSubmitStatus(createInputStatus(STATUS_SUCCESS));
          setLoading(false);
          setLogoCacheBuster(Math.floor(Math.random() * 100)); // Force logo to re-load
        },
        onError: (msg: string) => {
          setSubmitStatus(createInputStatus(STATUS_ERROR, `There was an error: ${msg}`));
          setLoading(false);
        },
      });
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
    }
  };

  const logoDisplayUrl = `${NEXT_PUBLIC_API_HOST}logo?random=${logoCachedbuster}`;

  return (
    <div className="formfield-container logo-upload-container">
      <div className="label-side">
        <span className="formfield-label">Logo</span>
      </div>

      <div className="input-side">
        <div className="input-group">
          <img src={logoDisplayUrl} alt="avatar" className="logo-preview" />
          <Upload
            name="logo"
            listType="picture"
            className="avatar-uploader"
            showUploadList={false}
            accept={ACCEPTED_FILE_TYPES.join(',')}
            beforeUpload={beforeUpload}
            customRequest={handleLogoUpdate}
            disabled={loading}
          >
            {loading ? (
              <LoadingOutlined style={{ color: 'white' }} />
            ) : (
              <Button icon={<UploadOutlined />} />
            )}
          </Upload>
        </div>
        <FormStatusIndicator status={submitStatus} />
        <p className="field-tip">{tip}</p>
      </div>
    </div>
  );
}
