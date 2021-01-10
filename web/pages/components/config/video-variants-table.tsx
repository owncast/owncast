import React, { useContext } from 'react';
import { Typography, Table, Modal } from 'antd';
import { ServerStatusContext } from '../../../utils/server-status-context';

const { Title } = Typography;


export default function CurrentVariantsTable() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};
  const { videoSettings } = serverConfig || {};
  const { videoQualityVariants } = videoSettings || {};
  if (!videoSettings) {
    return null;
  }

  
  const videoQualityColumns = [
    {
      title: "#",
      dataIndex: "key",
      key: "key"
    },
    {
      title: "Video bitrate",
      dataIndex: "videoBitrate",
      key: "videoBitrate",
      render: (bitrate: number) =>
        !bitrate ? "Same as source" : `${bitrate} kbps`,
    },
    {
      title: "Framerate",
      dataIndex: "framerate",
      key: "framerate",
      render: (framerate: number) =>
        !framerate ? "Same as source" : `${framerate} fps`,
    },
    {
      title: "Encoder preset",
      dataIndex: "encoderPreset",
      key: "encoderPreset",
      render: (preset: string) =>
        !preset ? "n/a" : preset,
    },
    {
      title: "Audio bitrate",
      dataIndex: "audioBitrate",
      key: "audioBitrate",
      render: (bitrate: number) =>
        !bitrate ? "Same as source" : `${bitrate} kbps`,
    },
    {
      title: '',
      dataIndex: '',
      key: 'edit',
      render: () => "edit.. populate modal",
    },
  ];


  const videoQualityVariantData = videoQualityVariants.map((variant, index) => ({ key: index + 1, ...variant }));

  return (
    <>
      <Title level={3}>Current Variants</Title>
      
      <Table
        pagination={false}
        size="small"
        columns={videoQualityColumns}
        dataSource={videoQualityVariantData}
      />

    </>
  );
}