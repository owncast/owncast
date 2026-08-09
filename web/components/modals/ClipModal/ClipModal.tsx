import { FC, useEffect, useState } from 'react';
import { Alert, Button, Input, Modal, Space, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../../types/localization';
import styles from './ClipModal.module.scss';

export type ClipModalProps = {
  open: boolean;
  handleClose: () => void;
  completedClipId: string;
  onSaveTitle: (title: string) => void;
  onDiscard: () => void;
  saving: boolean;
  error: string;
};

// ClipModal asks for a title before the explicit clip window is persisted. Once
// creation succeeds it shows the share link for the saved clip.
export const ClipModal: FC<ClipModalProps> = ({
  open,
  handleClose,
  completedClipId,
  onSaveTitle,
  onDiscard,
  saving,
  error,
}) => {
  const { t } = useTranslation();
  const [title, setTitle] = useState('');
  const [copied, setCopied] = useState(false);
  const clipUrl = completedClipId ? `${window.location.origin}/clips/${completedClipId}` : '';

  useEffect(() => {
    if (open) {
      setTitle('');
      setCopied(false);
    }
  }, [open, completedClipId]);

  const saveTitle = () => onSaveTitle(title.trim());
  const saveWithoutTitle = () => onSaveTitle('');
  const copyLink = async () => {
    try {
      await navigator.clipboard.writeText(clipUrl);
      setCopied(true);
    } catch {
      setCopied(false);
    }
  };

  return (
    <Modal
      title={
        completedClipId
          ? t(Localization.Frontend.Clips.created)
          : t(Localization.Frontend.Clips.titlePrompt)
      }
      open={open}
      onCancel={handleClose}
      width="450px"
      zIndex={999}
      footer={null}
      centered
    >
      <Space orientation="vertical" className={styles.body} id="clip-modal">
        {!completedClipId ? (
          <>
            <Typography.Text>{t(Localization.Frontend.Clips.titlePrompt)}</Typography.Text>
            <Input
              placeholder={t(Localization.Frontend.Clips.titlePlaceholder)}
              value={title}
              maxLength={100}
              onChange={e => setTitle(e.target.value)}
              onPressEnter={saveTitle}
              autoFocus
            />
            {error && <Alert type="error" message={error} showIcon />}
            <Space>
              <Button type="primary" loading={saving} onClick={saveTitle} disabled={!title.trim()}>
                {t(Localization.Frontend.Clips.saveTitle)}
              </Button>
              <Button loading={saving} onClick={saveWithoutTitle}>
                {t(Localization.Frontend.Clips.skipTitle)}
              </Button>
              <Button danger onClick={onDiscard} disabled={saving} id="discard-clip-button">
                {t(Localization.Frontend.Clips.discard)}
              </Button>
            </Space>
          </>
        ) : (
          <>
            <Alert type="success" message={t(Localization.Frontend.Clips.created)} showIcon />
            <Input readOnly value={clipUrl} onFocus={e => e.target.select()} />
            <Button onClick={copyLink}>
              {copied
                ? t(Localization.Frontend.Clips.linkCopied)
                : t(Localization.Frontend.Clips.copyLink)}
            </Button>
          </>
        )}
      </Space>
    </Modal>
  );
};
