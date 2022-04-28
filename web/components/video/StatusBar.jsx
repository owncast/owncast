interface Props {
  online: boolean;
  viewers: number;
  timer: string;
}

export default function StatusBar(props: Props) {
  return <div>Stream status bar goes here</div>;
}
