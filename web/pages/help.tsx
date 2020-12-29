import { Button, Card, Col, Divider, Row } from 'antd'
import Meta from 'antd/lib/card/Meta'
import Title from 'antd/lib/typography/Title'
import {
    BugTwoTone,
    CameraTwoTone,
    DatabaseTwoTone,
    DislikeTwoTone,
    EditTwoTone,
    FireFilled,
    FireOutlined,
    Html5TwoTone,
    QuestionCircleFilled,
    SettingTwoTone,
    SlidersTwoTone,
    VideoCameraTwoTone
} from '@ant-design/icons';
import React from 'react'

interface Props { }

export default function Help(props: Props) {
    const questions = [
        {
            icon: <SettingTwoTone />,
            title: "I want to configure my owncast instance",
            content: (
                <div>
                    <a href="https://owncast.online/docs/configuration/"  target="_blank" rel="noopener noreferrer">Learn more</a>
                </div>
            )
        },
        {
            icon: <CameraTwoTone />,
            title: "I need help configuring my broadcasting software",
            content: (
                <div>
                    <a href="https://owncast.online/docs/broadcasting/"  target="_blank" rel="noopener noreferrer">Learn more</a>
                </div>
            )
        },
        {
            icon: <Html5TwoTone />,
            title: "I want to embed my stream into another site",
            content: (
                <div>
                    <a href="https://owncast.online/docs/embed/"  target="_blank" rel="noopener noreferrer">Learn more</a>
                </div>
            )
        },
        {
            icon: <EditTwoTone />,
            title: "I want to customize my website",
            content: (
                <div>
                    <a href="https://owncast.online/docs/website/"  target="_blank" rel="noopener noreferrer">Learn more</a>
                </div>
            )
        },
        {
            icon: <SlidersTwoTone />,
            title: "I want to tweak my encoding",
            content: (
                <div>
                    <a href="https://owncast.online/docs/encoding/"  target="_blank" rel="noopener noreferrer">Learn more</a>
                </div>
            )
        },
        {
            icon: <DatabaseTwoTone />,
            title: "I want to offload my video to an external storage provider",
            content: (
                <div>
                    <a href="https://owncast.online/docs/encoding/"  target="_blank" rel="noopener noreferrer">Learn more</a>
                </div>
            )
        },
    ]

    const otherResources = [
        {
            icon: <BugTwoTone />,
            title: "I found a bug",
            content: (
                <div>
                    If you found a bug, then report it in our 
                    <a href="https://owncast.online/docs/encoding/" target="_blank" rel="noopener noreferrer"> Github Issues</a>
                </div>
            )
        }
    ]

    return (
        <div>
            <Title>How can we help you?</Title>
            <Row gutter={[16, 16]} justify="space-around" align="middle">
                <Col span={12}>
                    <Title level={2}>Having issues with owncast?</Title>
                    <Button href="https://owncast.online/docs/troubleshooting/" icon={<FireFilled/>} type="primary">Try Troubleshooting</Button>
                </Col>
                <Col span={12}>
                    <Title level={2}>Having any questions about owncast?</Title>
                    <Button href="https://owncast.online/docs/faq/" icon={<QuestionCircleFilled/>} type="primary">Read our FAQ</Button>
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
