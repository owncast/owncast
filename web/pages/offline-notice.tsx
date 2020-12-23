import { Result, Card } from "antd";
import { MessageTwoTone, QuestionCircleTwoTone, BookTwoTone, PlaySquareTwoTone } from '@ant-design/icons';
import OwncastLogo from "./components/logo"
import LogTable from "./components/log-table";

const { Meta } = Card;

export default function Offline({ logs = [] }) {
  const data = [
    {
      icon: <BookTwoTone twoToneColor="#6f42c1" />,
      title: "Use your broadcasting software",
      content: (
        <div>
          <a href="https://owncast.online/docs/broadcasting/">Learn how to point your existing software to your new server and start streaming your content.</a>
        </div>
      )
    },
    {
      icon: <MessageTwoTone twoToneColor="#0366d6" />,
      title: "Chat is disabled",
      content: "Chat will continue to be disabled until you begin a live stream."
    },
    {
      icon: <PlaySquareTwoTone twoToneColor="#f9826c" />,
      title: "Embed your video onto other sites",
      content: (
        <div>
          <a href="https://owncast.online/docs/embed">Learn how you can add your Owncast stream to other sites you control.</a>
        </div>
      )
    },
    {
      icon: <QuestionCircleTwoTone twoToneColor="#ffd33d" />,
      title: "Need some help?",
      content: (
        <div>
          Take a look at our <a target="_blank" href="https://owncast.online/docs/troubleshooting">troubleshooting steps </a>
          to see if there's any common problems you're running into.  If you still have questions <a target="_blank" href="https://github.com/owncast/owncast/issues/new/choose">
            let us know how we can help</a>.
        </div>
      ),
    }
  ];

  return (
    <>
      <div className="offline-content">
        <div className="logo-section">
          <Result
            icon={<OwncastLogo />}
            title="No stream is active."
            subTitle="You should start one."
          />
        </div>
        <div className="list-section">
          {
            data.map(item => (
              <Card key={item.title}>
                <Meta
                  avatar={item.icon}
                  title={item.title}
                  description={item.content}
                />
              </Card>
            ))
          }
        </div>

      </div>
      <LogTable logs={logs} pageSize={5} />
    </>
  );
}
