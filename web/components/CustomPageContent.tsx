interface Props {
  content: string;
}

export default function CustomPageContent(props: Props) {
  const { content } = props;
  return <div dangerouslySetInnerHTML={{ __html: content }} />;
}
