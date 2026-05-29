import React from 'react';
import { List, Modal, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Plugin } from '../../../interfaces/plugin';
import { Localization } from '../../../types/localization';
import { permissionDescriptionKey, permissionNameKey } from './permissionDescriptions';

const { Paragraph } = Typography;

export type InstallConfirmModalProps = {
  // The just-installed plugin, or null when the modal should be
  // hidden. Open state is derived from this rather than a separate
  // boolean so the parent only owns one piece of state.
  plugin: Plugin | null;
  onCancel: () => void;
  onEnable: () => void;
};

// InstallConfirmModal pops after a plugin is downloaded + installed
// from the registry. Lists the permissions the plugin's manifest
// declares (in plain language, with descriptions on hover) and asks
// the admin whether to enable it right now. Cancel leaves the plugin
// installed but disabled; the admin can flip the switch later from
// the Installed tab.
export const InstallConfirmModal = ({ plugin, onCancel, onEnable }: InstallConfirmModalProps) => {
  const { t } = useTranslation();
  const permissions = plugin?.permissions || [];
  const name = plugin?.name ?? '';

  return (
    <Modal
      open={!!plugin}
      title={t(Localization.Admin.Plugins.installConfirmTitle, { name })}
      onCancel={onCancel}
      onOk={onEnable}
      okText={t(Localization.Admin.Plugins.installConfirmEnable)}
      cancelText={t(Localization.Admin.Plugins.installConfirmCancel)}
      okButtonProps={{ type: 'primary' }}
    >
      {permissions.length === 0 ? (
        <Paragraph>{t(Localization.Admin.Plugins.installConfirmNoPermissions)}</Paragraph>
      ) : (
        <>
          <Paragraph>{t(Localization.Admin.Plugins.installConfirmPrompt)}</Paragraph>
          <List
            size="small"
            bordered
            dataSource={permissions}
            renderItem={p => {
              const nameKey = permissionNameKey[p];
              const descKey = permissionDescriptionKey[p];
              return (
                <List.Item>
                  <List.Item.Meta
                    title={nameKey ? t(nameKey) : p}
                    description={descKey ? t(descKey) : null}
                  />
                </List.Item>
              );
            }}
          />
        </>
      )}
    </Modal>
  );
};
