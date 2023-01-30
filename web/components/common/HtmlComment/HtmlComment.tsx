/* eslint-disable react/no-danger */
export const HtmlComment = ({ text }) => (
  <span style={{ display: 'none' }} dangerouslySetInnerHTML={{ __html: `\n\n<!-- ${text} -->` }} />
);
