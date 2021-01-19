import React, { useState, useContext, useEffect } from 'react';
import { Typography } from 'antd';
import SocialDropdown from './components/config/social-icons-dropdown';
import { fetchData, SOCIAL_PLATFORMS_LIST } from '../utils/apis';
import { ServerStatusContext } from '../utils/server-status-context';

const { Title } = Typography;


// get icons

export default function ConfigSocialLinks() {
  const [availableIconsList, setAvailableIconsList] = useState([]);
  const [currentSocialHandles, setCurrentSocialHandles] = useState([]);

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { socialHandles: initialSocialHandles } = instanceDetails;

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



  return (
    <div className="config-social-links">
      <Title level={2}>Social Links</Title>
      <p>Add all your social media handles and links to your other profiles here.</p>
        
      <SocialDropdown iconList={availableIconsList} />
    </div>
  ); 
}

