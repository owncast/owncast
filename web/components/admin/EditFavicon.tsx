import { Button, Upload, Popconfirm } from 'antd';
import { RcFile } from 'antd/lib/upload/interface';
import React, { useState, FC } from 'react';
import dynamic from 'next/dynamic';
import { useTranslation } from 'next-export-i18n';
import { FormStatusIndicator } from './FormStatusIndicator';
import {
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  TEXTFIELD_PROPS_FAVICON,
} from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { fetchData, NEXT_PUBLIC_API_HOST, SERVER_CONFIG_UPDATE_URL } from '../../utils/apis';
import { Localization } from '../../types/localization';
import { Translation } from '../ui/Translation/Translation';

import {
  ACCEPTED_FAVICON_TYPES,
  getBase64,
  MAX_FAVICON_FILESIZE,
  readableBytes,
} from '../../utils/images';

const LoadingOutlined = dynamic(() => import('@ant-design/icons/LoadingOutlined'), {
  ssr: false,
});

const UploadOutlined = dynamic(() => import('@ant-design/icons/UploadOutlined'), {
  ssr: false,
});

const UndoOutlined = dynamic(() => import('@ant-design/icons/UndoOutlined'), {
  ssr: false,
});

export const EditFavicon: FC = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [faviconUrl, setFaviconUrl] = useState(null);
  const [faviconCachebuster, setFaviconCacheBuster] = useState(0);
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  let resetTimer = null;

  const { apiPath } = TEXTFIELD_PROPS_FAVICON;

  const resetStates = () => {
    setSubmitStatus(null);
    clearTimeout(resetTimer);
    resetTimer = null;
  };

  // validate file type and size, then create base64 encoded img
  const beforeUpload = (file: RcFile) => {
    setLoading(true);

    // eslint-disable-next-line consistent-return
    return new Promise<void>((res, rej) => {
      if (file.size > MAX_FAVICON_FILESIZE) {
        const msg = t(Localization.Admin.StatusMessages.fileSizeTooBig, {
          size: readableBytes(file.size),
        });
        setSubmitStatus(
          createInputStatus(
            STATUS_ERROR,
            t(Localization.Admin.StatusMessages.thereWasAnError, { message: msg }),
          ),
        );
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
        setLoading(false);
        // eslint-disable-next-line no-promise-executor-return
        return rej();
      }
      if (!ACCEPTED_FAVICON_TYPES.includes(file.type)) {
        const msg = t(Localization.Admin.StatusMessages.fileTypeNotSupported, { type: file.type });
        setSubmitStatus(
          createInputStatus(
            STATUS_ERROR,
            t(Localization.Admin.StatusMessages.thereWasAnError, { message: msg }),
          ),
        );
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
        setLoading(false);
        // eslint-disable-next-line no-promise-executor-return
        return rej();
      }

      getBase64(file, (url: string) => {
        setFaviconUrl(url);
        setTimeout(() => res(), 100);
      });
    });
  };

  const handleFaviconUpdate = async () => {
    if (!faviconUrl) {
      setLoading(false);
      return;
    }

    setSubmitStatus(createInputStatus(STATUS_PROCESSING));
    await postConfigUpdateToAPI({
      apiPath,
      data: { value: faviconUrl },
      onSuccess: () => {
        setSubmitStatus(createInputStatus(STATUS_SUCCESS));
        setLoading(false);
        setFaviconCacheBuster(Math.floor(Math.random() * 100));
      },
      onError: (msg: string) => {
        setSubmitStatus(
          createInputStatus(
            STATUS_ERROR,
            t(Localization.Admin.StatusMessages.thereWasAnError, { message: msg }),
          ),
        );
        setLoading(false);
      },
    });
    setFaviconUrl(null);
    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
  };

  const handleResetFavicon = async () => {
    setLoading(true);
    setSubmitStatus(createInputStatus(STATUS_PROCESSING));

    try {
      const result = await fetchData(`${SERVER_CONFIG_UPDATE_URL}${apiPath}`, {
        method: 'DELETE',
      });

      if (result.success) {
        setSubmitStatus(createInputStatus(STATUS_SUCCESS));
        setFaviconCacheBuster(Math.floor(Math.random() * 100));
      } else {
        setSubmitStatus(
          createInputStatus(
            STATUS_ERROR,
            t(Localization.Admin.StatusMessages.thereWasAnError, { message: result.message }),
          ),
        );
      }
    } catch (error) {
      setSubmitStatus(
        createInputStatus(
          STATUS_ERROR,
          t(Localization.Admin.StatusMessages.thereWasAnError, { message: error.message }),
        ),
      );
    }

    setLoading(false);
    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
  };

  const faviconDisplayUrl = `${NEXT_PUBLIC_API_HOST}favicon.ico?random=${faviconCachebuster}`;

  return (
    <div className="formfield-container logo-upload-container">
      <div className="label-side">
        <span className="formfield-label">
          <Translation
            translationKey={Localization.Admin.EditFavicon.label}
            defaultText="Favicon"
          />
        </span>
      </div>

      <div className="input-side">
        <div className="input-group">
          <img
            src={faviconDisplayUrl}
            alt="favicon"
            style={{
              width: '48px',
              height: '48px',
              imageRendering: 'pixelated',
              marginRight: '10px',
            }}
          />
          <Upload
            name="favicon"
            listType="picture"
            className="avatar-uploader"
            showUploadList={false}
            accept={ACCEPTED_FAVICON_TYPES.join(',')}
            beforeUpload={beforeUpload}
            customRequest={handleFaviconUpdate}
            disabled={loading}
          >
            {loading ? (
              <LoadingOutlined style={{ color: 'white' }} />
            ) : (
              <Button icon={<UploadOutlined />} />
            )}
          </Upload>
          <Popconfirm
            title={t(Localization.Admin.EditFavicon.resetConfirmTitle)}
            onConfirm={handleResetFavicon}
            okText={t(Localization.Admin.EditFavicon.resetConfirmOk)}
            cancelText={t(Localization.Admin.EditFavicon.resetConfirmCancel)}
            disabled={loading}
          >
            <Button icon={<UndoOutlined />} disabled={loading} style={{ marginLeft: '8px' }}>
              <Translation
                translationKey={Localization.Admin.EditFavicon.resetButton}
                defaultText="Reset"
              />
            </Button>
          </Popconfirm>
        </div>
        <FormStatusIndicator status={submitStatus} />
        <p className="field-tip">
          <Translation
            translationKey={Localization.Admin.EditFavicon.tip}
            defaultText="Upload a custom favicon (PNG or ICO format, max 200KB). This icon appears in browser tabs and bookmarks."
          />
        </p>
      </div>
    </div>
  );
};
