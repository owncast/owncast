import { Col, Row } from 'antd';
import s from './ChatModerationDetailsModal.module.scss';

/* eslint-disable @typescript-eslint/no-unused-vars */
interface Props {
  // userID: string;
}

export default function ChatModerationDetailsModal(props: Props) {
  return (
    <div className={s.modalContainer}>
      <Row justify="space-around" align="middle">
        <Col span={12}>User created</Col>
        <Col span={12}>xxxx</Col>
      </Row>

      <Row justify="space-around" align="middle">
        <Col span={12}>Previous names</Col>
        <Col span={12}>xxxx</Col>
      </Row>

      <h1>Recent Chat Messages</h1>

      <div className={s.chatHistory} />
    </div>
  );
}
