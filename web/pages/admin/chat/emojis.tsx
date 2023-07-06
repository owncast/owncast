import { Avatar, Button, Card, Col, Row, Tooltip, Typography } from 'antd';
import Upload, { RcFile } from 'antd/lib/upload';
import React, { ReactElement, useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import FormStatusIndicator from '../../../components/admin/FormStatusIndicator';
import { DELETE_EMOJI, fetchData, UPLOAD_EMOJI } from '../../../utils/apis';
import { ACCEPTED_IMAGE_TYPES, getBase64 } from '../../../utils/images';
import {
  createInputStatus,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../../utils/input-statuses';
import { RESET_TIMEOUT } from '../../../utils/config-constants';
import { AdminLayout } from '../../../components/layouts/AdminLayout';

const URL_CUSTOM_EMOJIS = `/api/emoji`;

const { Meta } = Card;
// Lazy loaded components

const CloseOutlined = dynamic(() => import('@ant-design/icons/CloseOutlined'), {
  ssr: false,
});

type CustomEmoji = {
  name: string;
  url: string;
};

const { Title, Paragraph } = Typography;

const Emoji = () => {
  const [emojis, setEmojis] = useState<CustomEmoji[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitStatus, setSubmitStatus] = useState(null);
  const [uploadFile, setUploadFile] = useState<RcFile>(null);

  let resetTimer = null;
  const resetStates = () => {
    setSubmitStatus(null);
    clearTimeout(resetTimer);
    resetTimer = null;
  };

  async function getEmojis() {
    setLoading(true);
    try {
      const response = await fetchData(URL_CUSTOM_EMOJIS);
      setEmojis(response);
    } catch (error) {
      console.error('error fetching emojis', error);
    }
    setLoading(false);
  }
  useEffect(() => {
    getEmojis();
  }, []);

  async function handleDelete(fullPath: string) {
    const name = `/${fullPath.split('/').slice(3).join('/')}`;
    console.log(name);

    setLoading(true);

    setSubmitStatus(createInputStatus(STATUS_PROCESSING, 'Deleting emoji...'));

    try {
      const response = await fetchData(DELETE_EMOJI, {
        method: 'POST',
        data: { name },
      });

      if (response instanceof Error) {
        throw response;
      }

      setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Emoji deleted'));
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
    } catch (error) {
      setSubmitStatus(createInputStatus(STATUS_ERROR, `${error}`));
      setLoading(false);
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
    }

    getEmojis();
  }

  async function handleEmojiUpload() {
    setLoading(true);
    try {
      setSubmitStatus(createInputStatus(STATUS_PROCESSING, 'Converting emoji...'));

      // eslint-disable-next-line consistent-return
      const emojiData = await new Promise<CustomEmoji>((res, rej) => {
        if (!ACCEPTED_IMAGE_TYPES.includes(uploadFile.type)) {
          const msg = `File type is not supported: ${uploadFile.type}`;
          // eslint-disable-next-line no-promise-executor-return
          return rej(msg);
        }

        getBase64(uploadFile, (url: string) =>
          res({
            name: uploadFile.name,
            url,
          }),
        );
      });

      setSubmitStatus(createInputStatus(STATUS_PROCESSING, 'Uploading emoji...'));

      const response = await fetchData(UPLOAD_EMOJI, {
        method: 'POST',
        data: {
          name: emojiData.name,
          data: emojiData.url,
        },
      });

      if (response instanceof Error) {
        throw response;
      }

      setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Emoji uploaded successfully!'));
      getEmojis();
    } catch (error) {
      setSubmitStatus(createInputStatus(STATUS_ERROR, `${error}`));
    }

    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
    setLoading(false);
  }

  return (
    <div>
      <Title>Emojis</Title>
      <Paragraph>
        Here you can upload new custom emojis for usage in the chat. When uploading a new emoji, the
        filename without extension will be used as emoji name. Additionally, emoji names are
        case-insensitive. For best results, ensure all emoji have unique names.
      </Paragraph>
      <br />
      <Upload
        name="emoji"
        listType="picture"
        className="emoji-uploader"
        showUploadList={false}
        accept={ACCEPTED_IMAGE_TYPES.join(',')}
        beforeUpload={setUploadFile}
        customRequest={handleEmojiUpload}
        disabled={loading}
      >
        <Button type="primary" disabled={loading}>
          Upload new emoji
        </Button>
      </Upload>
      <FormStatusIndicator status={submitStatus} />
      <br />
      <Row>
        {emojis.map(record => (
          <Col style={{ padding: '10px' }} key={record.name}>
            <Card style={{ width: 120, marginTop: 16 }} actions={[]}>
              <Meta
                description={[
                  <div
                    style={{
                      display: 'flex',
                      justifyItems: 'center',
                      alignItems: 'center',
                      flexDirection: 'column',
                      gap: '20px',
                    }}
                  >
                    <Tooltip title={record.name}>
                      <Avatar style={{ height: 50, width: 50 }} src={record.url} />
                    </Tooltip>
                    <Button
                      size="small"
                      type="ghost"
                      title="Delete emoji"
                      style={{
                        position: 'absolute',
                        right: 0,
                        top: 0,
                        height: 24,
                        width: 24,
                        border: 'none',
                        color: 'gray',
                      }}
                      onClick={() => handleDelete(record.url)}
                      icon={<CloseOutlined />}
                    />
                  </div>,
                ]}
              />
            </Card>
          </Col>
        ))}
      </Row>
      <br />
    </div>
  );
};

Emoji.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};

export default Emoji;
