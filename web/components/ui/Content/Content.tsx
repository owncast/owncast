import { useRecoilValue } from 'recoil';
import { clientConfigState } from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { Layout, Row, Col } from 'antd';

const { Content } = Layout;

export default function FooterComponent() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigState);
  const { extraPageContent } = clientConfig;

  return (
    <Content style={{ margin: '80px 16px 0', overflow: 'initial' }}>
      <div>
        <Row>
          <Col span={24}>Video player goes here</Col>
        </Row>
        <Row>
          <Col span={24}>
            <Content dangerouslySetInnerHTML={{ __html: extraPageContent }} />
          </Col>
        </Row>
      </div>
    </Content>
  );
}
