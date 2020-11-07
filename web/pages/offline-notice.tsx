import { Result, List, Card } from "antd";
import OwncastLogo from "./components/logo"

export default function Offline() {
  const data = [
    {
      title: "Send some test content",
      content: (
        <div>
          With any video you have around you can pass it to the test script and start streaming it.
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
      title: "Something else",
    },
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
          md: 4,
          lg: 4,
          xl: 6,
          xxl: 3,
        }}
        dataSource={data}
        renderItem={(item) => (
          <List.Item>
            <Card title={item.title}>{item.content}</Card>
          </List.Item>
        )}
      />
    </div>
  );
}
