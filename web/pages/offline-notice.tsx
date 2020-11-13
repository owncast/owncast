import { Result, List, Card } from "antd";
import OwncastLogo from "./components/logo"

export default function Offline() {
  const data = [
    {
      title: "Send some test content",
      content: (
        <div>
          Test your server with any video you have around. Pass it to the test script and start streaming it.
          <blockquote>
            <em>./test/ocTestStream.sh yourVideo.mp4</em>
          </blockquote>
        </div>
      ),
    },
    {
      title: "Use your broadcasting software",
      content: (
        <div>
          <a href="https://owncast.online/docs/broadcasting/">Learn how to point your existing software to your new server and start streaming your content.</a>
        </div>
      )
    },
    {
      title: "Chat is disabled",
      content: "Chat will continue to be disabled until you begin a live stream."
    },
    {
      title: "Embed your video onto other sites",
      content: (
        <div>
          <a href="https://owncast.online/docs/embed">Learn how you can add your Owncast stream to other sites you control.</a>
        </div>
      )
    }
  ];
  return (
    <div>
      <Result
        icon={<OwncastLogo />}
        title="No stream is active."
        subTitle="You should start one."
      />

      <List
        grid={{
          gutter: 16,
          xs: 1,
          sm: 2,
          md: 2,
          lg: 6,
          xl: 3,
          xxl: 3,
        }}
        dataSource={data}
        renderItem={(item) => (
          <List.Item>
            <Card title={item.title}>{item.content}</Card>
          </List.Item>
        )}
      />
      {logTable}
    </div>
  );
}
