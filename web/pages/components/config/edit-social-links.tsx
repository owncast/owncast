import React, { useState, useContext, useEffect } from 'react';
import { Typography, Table, Button, Modal, Input } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { DeleteOutlined } from '@ant-design/icons';
import SocialDropdown from './social-icons-dropdown';
import { fetchData, NEXT_PUBLIC_API_HOST, SOCIAL_PLATFORMS_LIST } from '../../../utils/apis';
import { ServerStatusContext } from '../../../utils/server-status-context';
import { API_SOCIAL_HANDLES, postConfigUpdateToAPI, RESET_TIMEOUT, SUCCESS_STATES, DEFAULT_SOCIAL_HANDLE, OTHER_SOCIAL_HANDLE_OPTION } from './constants';
import { SocialHandle } from '../../../types/config-section';
import {isValidUrl} from '../../../utils/urls';

import configStyles from '../../../styles/config-pages.module.scss';

const { Title } = Typography;

export default function EditSocialLinks() {
  const [availableIconsList, setAvailableIconsList] = useState([]);
  const [currentSocialHandles, setCurrentSocialHandles] = useState([]);

  const [displayModal, setDisplayModal] = useState(false);
  const [displayOther, setDisplayOther] = useState(false);
  const [modalProcessing, setModalProcessing] = useState(false);
  const [editId, setEditId] = useState(-1);
  
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
      setAvailableIconsList(list);

    } catch (error) {
      console.log(error)
      //  do nothing
    }
  };

  const selectedOther = modalDataState.platform !== '' && !availableIconsList.find(item => item.key === modalDataState.platform);


  useEffect(() => {
    getAvailableIcons();
  }, []);

  useEffect(() => {
    if (instanceDetails.socialHandles) {
      setCurrentSocialHandles(initialSocialHandles);
    }
  }, [instanceDetails]);


  const resetStates = () => {
    setSubmitStatus(null);
    setSubmitStatusMessage('');
    resetTimer = null;
    clearTimeout(resetTimer);
  };
  const resetModal = () => {
    setDisplayModal(false);
    setEditId(-1);
    setDisplayOther(false);
    setModalProcessing(false);
    setModalDataState({...DEFAULT_SOCIAL_HANDLE});
  };

  const handleModalCancel = () => {
    resetModal();
  };

  const updateModalState = (fieldName: string, value: string) => {
    setModalDataState({
      ...modalDataState,
      [fieldName]: value,
    });
  };
  const handleDropdownSelect = (value: string) => {
    if (value === OTHER_SOCIAL_HANDLE_OPTION) {
      setDisplayOther(true);
      updateModalState('platform', '');
    } else {
      setDisplayOther(false);
      updateModalState('platform', value);
    }
  };
  const handleOtherNameChange = event => {
    const { value } = event.target;
    updateModalState('platform', value);
  };
  const handleUrlChange = event => {
    const { value } = event.target;
    updateModalState('url', value);
  };
  

  // posts all the variants at once as an array obj
  const postUpdateToAPI = async (postValue: any) => {
    await postConfigUpdateToAPI({
      apiPath: API_SOCIAL_HANDLES,
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({ fieldName: 'socialHandles', value: postValue, path: 'instanceDetails' });

        // close modal
        setModalProcessing(false);
        handleModalCancel();

        setSubmitStatus('success');
        setSubmitStatusMessage('Social Handles updated.');
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
    const postData = currentSocialHandles.length ? [
      ...currentSocialHandles,
    ]: [];
    if (editId === -1) {
      postData.push(modalDataState);
    } else {
      postData.splice(editId, 1, modalDataState);
    }
    postUpdateToAPI(postData);
  };

  const handleDeleteItem = index => {
    const postData = [
      ...currentSocialHandles,
    ];
    postData.splice(index, 1);
    postUpdateToAPI(postData);
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
        const platformInfo = availableIconsList.find(item => item.key === platform);
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
      title: "Url Link",
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
              setModalDataState({...currentSocialHandles[index]});
              setDisplayModal(true);
            }}>
              Edit
            </Button>
            <Button
              className="delete-button"
              icon={<DeleteOutlined />}
              size="small"
              onClick={() =>  handleDeleteItem(index)}
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

  const okButtonProps = {
    disabled: !isValidUrl(modalDataState.url)
};


  return (
    <div className={configStyles.socialLinksEditor}>
      <Title level={2}>Social Links</Title>
      <p>Add all your social media handles and links to your other profiles here.</p>

      {statusMessage}

      <Table
        className="dataTable"
        pagination={false}
        size="small"
        rowKey={record => record.url}
        columns={socialHandlesColumns}
        dataSource={currentSocialHandles}
      />

      <Modal
        title="Edit Social Handle"
        visible={displayModal}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={modalProcessing}
        okButtonProps={okButtonProps}
      >
        <SocialDropdown
          iconList={availableIconsList}
          selectedOption={selectedOther ? OTHER_SOCIAL_HANDLE_OPTION : modalDataState.platform}
          onSelected={handleDropdownSelect}
        />
        {
          displayOther
            ? (
            <>
              <Input
                placeholder="Other"
                defaultValue={modalDataState.platform}
                onChange={handleOtherNameChange}
              />
              <br/>
            </>
            ) : null
        }
        <br/>

        URL
        <Input
          placeholder="Url to page"
          defaultValue={modalDataState.url}
          value={modalDataState.url}
          onChange={handleUrlChange}
        />

        {statusMessage}
      </Modal>
      <br />
      <Button type="primary" onClick={() => {
          resetModal();
          setDisplayModal(true);
        }}>
        Add a new social link
      </Button>
    </div>
  ); 
}

