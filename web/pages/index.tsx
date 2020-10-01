import { Card, Alert, Statistic, Row, Col } from "antd";
import { LikeOutlined } from "@ant-design/icons";

const { Meta } = Card;

export default function Home() {
  return (
    <div>
      <Alert
        message="These are some ant design component example I stole from their web site."
        type="success"
      />
      <Row gutter={16}>
        <Col span={12}>
          <Statistic title="Feedback" value={1128} prefix={<LikeOutlined />} />
        </Col>
        <Col span={12}>
          <Statistic title="Unmerged" value={93} suffix="/ 100" />
        </Col>
      </Row>
      <Card
        hoverable
        style={{ width: 240 }}
        cover={
          <img
            alt="example"
            src="https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png"
          />
        }
      >
        <Meta title="Europe Street beat" description="www.instagram.com" />
      </Card>
    </div>
  );
}
