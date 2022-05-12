interface Props {
  content: string;
}

export default function CustomPageContent(props: Props) {
  const { content } = props;
  // eslint-disable-next-line react/no-danger
  return <div dangerouslySetInnerHTML={{ __html: content }} />;
}
