import { Button, Card, Col, Divider, Row } from 'antd'
import Meta from 'antd/lib/card/Meta'
import Title from 'antd/lib/typography/Title'
import {
    AlertOutlined,
    AlertTwoTone,
    ApiTwoTone,
    BookOutlined,
    BugTwoTone,
    CameraTwoTone,
    DatabaseTwoTone,
    DislikeTwoTone,
    EditTwoTone,
    FireFilled,
    FireOutlined,
    Html5TwoTone,
    LinkOutlined,
    QuestionCircleFilled,
    QuestionCircleTwoTone,
    SettingTwoTone,
    SlidersTwoTone,
    VideoCameraTwoTone
} from '@ant-design/icons';
import React from 'react'
import Text from 'antd/lib/typography/Text';

interface Props { }

export default function Help(props: Props) {
    const questions = [
        {
            icon: <SettingTwoTone style={{ fontSize: '24px' }} />,
            title: "I want to configure my owncast instance",
            content: (
                <div>
                    <a href="https://owncast.online/docs/configuration/"  target="_blank" rel="noopener noreferrer"><LinkOutlined/> Learn more</a>
                </div>
            )
        },
        {
            icon: <CameraTwoTone style={{ fontSize: '24px' }} />,
            title: "I need help configuring my broadcasting software",
            content: (
                <div>
                    <a href="https://owncast.online/docs/broadcasting/"  target="_blank" rel="noopener noreferrer"><LinkOutlined/> Learn more</a>
                </div>
            )
        },
        {
            icon: <Html5TwoTone style={{ fontSize: '24px' }} />,
            title: "I want to embed my stream into another site",
            content: (
                <div>
                    <a href="https://owncast.online/docs/embed/"  target="_blank" rel="noopener noreferrer"><LinkOutlined/> Learn more</a>
                </div>
            )
        },
        {
            icon: <EditTwoTone style={{ fontSize: '24px' }} />,
            title: "I want to customize my website",
            content: (
                <div>
                    <a href="https://owncast.online/docs/website/"  target="_blank" rel="noopener noreferrer"><LinkOutlined/> Learn more</a>
                </div>
            )
        },
        {
            icon: <SlidersTwoTone style={{ fontSize: '24px' }} />,
            title: "I want to tweak my encoding",
            content: (
                <div>
                    <a href="https://owncast.online/docs/encoding/"  target="_blank" rel="noopener noreferrer"><LinkOutlined/> Learn more</a>
                </div>
            )
        },
        {
            icon: <DatabaseTwoTone style={{ fontSize: '24px' }} />,
            title: "I want to offload my video to an external storage provider",
            content: (
                <div>
                    <a href="https://owncast.online/docs/encoding/"  target="_blank" rel="noopener noreferrer"><LinkOutlined/> Learn more</a>
                </div>
            )
        },
    ]

    const otherResources = [
        {
            icon: <BugTwoTone style={{ fontSize: '24px' }}  />,
            title: "I found a bug",
            content: (
                <div>
                    If you found a bug, then report it in our 
                    <a href="https://owncast.online/docs/encoding/" target="_blank" rel="noopener noreferrer"> Github Issues</a>
                </div>
            )
        },
        {
            icon: <QuestionCircleTwoTone style={{ fontSize: '24px' }}  />,
            title: "I have a general question",
            content: (
                <div>
                    Most general questions are answered in our 
                    <a href="https://owncast.online/docs/encoding/" target="_blank" rel="noopener noreferrer"> FAQ</a>
                </div>
            )
        },
        {
            icon: <ApiTwoTone style={{ fontSize: '24px' }}  />,
            title: "I want to use the API",
            content: (
                <div>
                    You can view the API documentation for either the 
                    <a href="https://owncast.online/api/latest" target="_blank" rel="noopener noreferrer"> latest</a>
                    or 
                    <a href="https://owncast.online/api/development" target="_blank" rel="noopener noreferrer"> development</a>
                    release.
                </div>
            )
        }
    ]

    return (
        <div>
            <Title style={{textAlign: 'center'}}>How can we help you?</Title>
            <Row gutter={[16, 16]} justify="space-around" align="middle">
                <Col span={12} style={{textAlign: 'center'}}>
                    <AlertOutlined style={{ fontSize: '64px' }}/>
                    <Title level={2}>Troubleshooting</Title>
                    <Button href="https://owncast.online/docs/troubleshooting/" icon={<LinkOutlined/>} type="primary">Read Troubleshoting</Button>
                </Col>
                <Col span={12} style={{textAlign: 'center'}}>
                    <BookOutlined style={{ fontSize: '64px' }}/>
                    <Title level={2}>Documentation</Title>
                    <Button href="https://owncast.online/docs/faq/" icon={<LinkOutlined/>} type="primary">Read the Docs</Button>
                </Col>
            </Row>
            <Divider />
            <Title level={2}>Common tasks</Title>
            <Row gutter={[16, 16]}>
                {
                    questions.map(question => (
                        <Col span={12}>
                            <Card key={question.title}>
                                <Meta
                                    avatar={question.icon}
                                    title={question.title}
                                    description={question.content}
                                />
                            </Card>
                        </Col>
                    ))
                }
            </Row>
            <Divider />
            <Title level={2}>Other</Title>
            <Row gutter={[16, 16]}>
                {
                    otherResources.map(question => (
                        <Col span={12}>
                            <Card key={question.title}>
                                <Meta
                                    avatar={question.icon}
                                    title={question.title}
                                    description={question.content}
                                />
                            </Card>
                        </Col>
                    ))
                }
            </Row>
        </div>
    )
}
