import { Button, Upload, Popconfirm } from 'antd';
import { RcFile } from 'antd/lib/upload/interface';
import React, { useState, useRef, FC } from 'react';
import dynamic from 'next/dynamic';
import { useTranslation } from 'next-export-i18n';
import { FormStatusIndicator } from './FormStatusIndicator';
import { RESET_TIMEOUT } from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { NEXT_PUBLIC_API_HOST } from '../../utils/apis';
import { Localization } from '../../types/localization';
import { Translation } from '../ui/Translation/Translation';

import { ACCEPTED_FAVICON_TYPES, MAX_FAVICON_FILESIZE, readableBytes } from '../../utils/images';

const LoadingOutlined = dynamic(() => import('@ant-design/icons/LoadingOutlined'), {
  ssr: false,
});

const UploadOutlined = dynamic(() => import('@ant-design/icons/UploadOutlined'), {
  ssr: false,
});

const UndoOutlined = dynamic(() => import('@ant-design/icons/UndoOutlined'), {
  ssr: false,
});

const ADMIN_USERNAME = process.env.NEXT_PUBLIC_ADMIN_USERNAME;
const ADMIN_STREAMKEY = process.env.NEXT_PUBLIC_ADMIN_STREAMKEY;

export const EditFavicon: FC = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [faviconCachebuster, setFaviconCacheBuster] = useState(0);
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const pendingFile = useRef<RcFile | null>(null);
  let resetTimer = null;

  const resetStates = () => {
    setSubmitStatus(null);
    clearTimeout(resetTimer);
    resetTimer = null;
  };

  // validate file type and size
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

      pendingFile.current = file;
      setTimeout(() => res(), 100);
    });
  };

  const handleFaviconUpdate = async () => {
    if (!pendingFile.current) {
      setLoading(false);
      return;
    }

    setSubmitStatus(createInputStatus(STATUS_PROCESSING));

    const formData = new FormData();
    formData.append('favicon', pendingFile.current);

    try {
      const encoded = btoa(`${ADMIN_USERNAME}:${ADMIN_STREAMKEY}`);
      const response = await fetch(`${NEXT_PUBLIC_API_HOST}api/admin/config/favicon`, {
        method: 'POST',
        headers: {
          Authorization: `Basic ${encoded}`,
        },
        body: formData,
      });

      const result = await response.json();

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

    pendingFile.current = null;
    setLoading(false);
    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
  };

  const handleResetFavicon = async () => {
    setLoading(true);
    setSubmitStatus(createInputStatus(STATUS_PROCESSING));

    try {
      const encoded = btoa(`${ADMIN_USERNAME}:${ADMIN_STREAMKEY}`);
      const response = await fetch(`${NEXT_PUBLIC_API_HOST}api/admin/config/favicon`, {
        method: 'DELETE',
        headers: {
          Authorization: `Basic ${encoded}`,
        },
      });

      const result = await response.json();

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
