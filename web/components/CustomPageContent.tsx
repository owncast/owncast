interface Props {
  content: string;
}

export default function CustomPageContent(props: Props) {
  const { content } = props;
  return <div>{content}</div>;
}
