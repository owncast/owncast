import React, { useContext, useState, useEffect } from 'react';
import { Typography, Table, Modal, Button } from 'antd';
import { CheckOutlined, CloseOutlined } from '@ant-design/icons';
import { ColumnsType } from 'antd/lib/table';
import { ServerStatusContext } from '../../../utils/server-status-context';
import { VideoVariant } from '../../../types/config-section';
import VideoVariantForm from './video-variant-form';
import { DEFAULT_VARIANT_STATE } from './constants';

const { Title } = Typography;

export default function CurrentVariantsTable() {
  const serverStatusData = useContext(ServerStatusContext);
  const [displayModal, setDisplayModal] = useState(false);
  const [editId, setEditId] = useState(0);
  const [dataState, setDataState] = useState(DEFAULT_VARIANT_STATE);
  const { serverConfig } = serverStatusData || {};
  const { videoSettings } = serverConfig || {};
  const { videoQualityVariants } = videoSettings || {};
  if (!videoSettings) {
    return null;
  }

  const handleModalOk = () => {
    setDisplayModal(false);
    setEditId(-1);
  }
  const handleModalCancel = () => {
    setDisplayModal(false);
    setEditId(-1);
  }

  const handleUpdateField = (fieldName: string, value: any) => {
    setDataState({
      ...dataState,
      [fieldName]: value,
    });
  }
  
  const videoQualityColumns: ColumnsType<VideoVariant>  = [
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
    // {
    //   title: "Framerate",
    //   dataIndex: "framerate",
    //   key: "framerate",
    //   render: (framerate: number) =>
    //     !framerate ? "Same as source" : `${framerate} fps`,
    // },
    {
      title: "Encoder preset",
      dataIndex: "encoderPreset",
      key: "encoderPreset",
      render: (preset: string) =>
        !preset ? "n/a" : preset,
    },
    // {
    //   title: "Audio bitrate",
    //   dataIndex: "audioBitrate",
    //   key: "audioBitrate",
    //   render: (bitrate: number) =>
    //     !bitrate ? "Same as source" : `${bitrate} kbps`,
    // },
    // {
    //   title: "Audio passthrough",
    //   dataIndex: "audioPassthrough",
    //   key: "audioPassthrough",
    //   render: (item: boolean) => item ? <CheckOutlined />: <CloseOutlined />,
    // },
    // {
    //   title: "Video passthrough",
    //   dataIndex: "videoPassthrough",
    //   key: "audioPassthrough",
    //   render: (item: boolean) => item ? <CheckOutlined />: <CloseOutlined />,
    // },
    {
      title: '',
      dataIndex: '',
      key: 'edit',
      render: (data: VideoVariant) => (
      <Button type="primary" size="small" onClick={() =>{
        setEditId(data.key - 1);
        setDisplayModal(true);
      }}>
        Edit
      </Button>
      ),
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

      <Modal
        title="Edit Video Variant Details"
        visible={displayModal}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        // confirmLoading={confirmLoading}
      >

        <VideoVariantForm initialValues={{...videoQualityVariants[editId]}} onUpdateField={handleUpdateField} />


      </Modal>

    </>
  );
}