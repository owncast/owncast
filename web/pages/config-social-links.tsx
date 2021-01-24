import React, { useState, useContext, useEffect } from 'react';
import { Typography, Table, Button, Modal } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { DeleteOutlined } from '@ant-design/icons';
import SocialDropdown from './components/config/social-icons-dropdown';
import { fetchData, NEXT_PUBLIC_API_HOST, SOCIAL_PLATFORMS_LIST } from '../utils/apis';
import { ServerStatusContext } from '../utils/server-status-context';
import { API_SOCIAL_HANDLES, postConfigUpdateToAPI, RESET_TIMEOUT, SUCCESS_STATES, DEFAULT_SOCIAL_HANDLE } from './components/config/constants';
import { SocialHandle } from '../types/config-section';

const { Title } = Typography;


// get icons

export default function ConfigSocialLinks() {
  const [availableIconsList, setAvailableIconsList] = useState([]);
  const [currentSocialHandles, setCurrentSocialHandles] = useState([]);

  const [displayModal, setDisplayModal] = useState(false);
  const [modalProcessing, setModalProcessing] = useState(false);
  const [editId, setEditId] = useState(0);
  
  // current data inside modal
  const [modalDataState, setModalDataState] = useState(DEFAULT_SOCIAL_HANDLE);

  const [submitStatus, setSubmitStatus] = useState(null);
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');


  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { socialHandles: initialSocialHandles } = instanceDetails;

  let resetTimer = null;

  const getAvailableIcons = async () => {
    try {
      const result = await fetchData(SOCIAL_PLATFORMS_LIST, { auth: false });
      const list = Object.keys(result).map(item => ({
        key: item,
        ...result[item],
      }));
      console.log({result})
      setAvailableIconsList(list);

    } catch (error) {
      console.log(error)
      //  do nothing
    }
  };

  useEffect(() => {
    getAvailableIcons();
  }, []);

  useEffect(() => {
    setCurrentSocialHandles(initialSocialHandles);
  }, [instanceDetails]);


  const resetStates = () => {
    setSubmitStatus(null);
    setSubmitStatusMessage('');
    resetTimer = null;
    clearTimeout(resetTimer);
  }

  const handleModalCancel = () => {
    setDisplayModal(false);
    setEditId(-1);
  }
  

  // posts all the variants at once as an array obj
  const postUpdateToAPI = async (postValue: any) => {
    await postConfigUpdateToAPI({
      apiPath: API_SOCIAL_HANDLES,
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({ fieldName: 'socialHandles', value: postValue, path: 'instancesDetails' });

        // close modal
        setModalProcessing(false);
        handleModalCancel();

        setSubmitStatus('success');
        setSubmitStatusMessage('Variants updated.');
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus('error');
        setSubmitStatusMessage(message);
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };


  // on Ok, send all of dataState to api
  // show loading
  // close modal when api is done
  const handleModalOk = () => {
    setModalProcessing(true);
    
    const postData = [
      ...currentSocialHandles,
    ];
    if (editId === -1) {
      postData.push(modalDataState);
    } else {
      postData.splice(editId, 1, modalDataState);
    }
    postUpdateToAPI(postData);
  };

  const handleDeleteVariant = index => {
    const postData = [
      ...currentSocialHandles,
    ];
    postData.splice(index, 1);
    postUpdateToAPI(postData)
  };


  const socialHandlesColumns: ColumnsType<SocialHandle>  = [
    {
      title: "#",
      dataIndex: "key",
      key: "key"
    },
    {
      title: "Platform",
      dataIndex: "platform",
      key: "platform",
      render: (platform: string) => {
        const platformInfo = availableIconsList[platform];
        if (!platformInfo) {
          return platform;
        }
        const { icon, platform: platformName } = platformInfo;
        return (
          <>
            <span className="option-icon">
              <img src={`${NEXT_PUBLIC_API_HOST}${icon}`} alt="" className="option-icon" />
            </span>
            <span className="option-label">{platformName}</span>
          </>
        );
      },
    },

    {
      title: "Url to profile",
      dataIndex: "url",
      key: "url",
    },
    {
      title: '',
      dataIndex: '',
      key: 'edit',
      render: (data, record, index) => {
        return (
          <span className="actions">
            <Button type="primary" size="small" onClick={() => {
              setEditId(index);
              setModalDataState(currentSocialHandles[index]);
              setDisplayModal(true);
            }}>
              Edit
            </Button>
            <Button
              className="delete-button"
              icon={<DeleteOutlined />}
              size="small"
              onClick={() => {
                handleDeleteVariant(index);
              }}
            />              
          </span>
        )},
      },
    ];
  
  const {
    icon: newStatusIcon = null,
    message: newStatusMessage = '',
  } = SUCCESS_STATES[submitStatus] || {};
  const statusMessage = (
    <div className={`status-message ${submitStatus || ''}`}>
      {newStatusIcon} {newStatusMessage} {submitStatusMessage}
    </div>
  );

  return (
    <div className="config-social-links">
      <Title level={2}>Social Links</Title>
      <p>Add all your social media handles and links to your other profiles here.</p>
        

      {statusMessage}

      <Table
        className="variants-table"
        pagination={false}
        size="small"
        columns={socialHandlesColumns}
        dataSource={currentSocialHandles}
      />

      <Modal
        title="Edit Social Handle"
        visible={displayModal}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={modalProcessing}
      >
        <SocialDropdown iconList={availableIconsList} />

        
        {statusMessage}
      </Modal>
      <br />
      <Button type="primary" onClick={() => {
          setEditId(-1);
          setModalDataState(DEFAULT_SOCIAL_HANDLE);
          setDisplayModal(true);
        }}>
        Add a new variant
      </Button>

    </div>
  ); 
}

