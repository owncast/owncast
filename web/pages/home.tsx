import React, { useState, useEffect } from 'react';
import { Layout } from 'antd';


import BroadcastInfo from './components/broadcast-info';
import HardwareInfo from './components/hardware-info';
import ViewerInfo from './components/viewer-info';
import ServerConfig from './components/server-config';
import ConnectedClients from './components/connected-clients';

export default function HomeView(props) {
  const { broadcastActive, broadcaster, message } = props;

  const { Header, Footer, Content } = Layout;

  const broadcastDetails = broadcastActive ? (
    <>
      {/* <BroadcastInfo {...broadcaster} />
      <HardwareInfo /> */}
      <ViewerInfo />
      <ConnectedClients />
      {/* <ServerConfig /> */}
    </>
  ) : null;

  const disconnectButton = broadcastActive ? <button type="button">Boot (Disconnect)</button> : null;

  return (
    <Layout className="layout">
      <Header>
        <div className="logo">logo</div>
        {/* <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['2']}>
          <Menu.Item key="1">nav 1</Menu.Item>
          <Menu.Item key="2">nav 2</Menu.Item>
          <Menu.Item key="3">nav 3</Menu.Item>
        </Menu> */}
      </Header>
      <Content style={{ padding: '0 50px' }}>
        <div className="site-layout-content">
          <p>
            <b>Status: {broadcastActive ? 'on' : 'off'}</b>
          </p>

          <h2>Utilities</h2>
          <p>(these dont do anything yet)</p>

          {disconnectButton}
          <button type="button">Change Stream Key</button>
          <button type="button">Server Config</button>

          <br />
          <br />

          {broadcastDetails}


        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}><a href="https://owncast.online/">About Owncast</a></Footer>
    </Layout>

    

  );
}
