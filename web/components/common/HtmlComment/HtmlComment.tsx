/* eslint-disable react/no-danger */
export const HtmlComment = ({ text }) => (
  <div dangerouslySetInnerHTML={{ __html: `<!-- ${text} -->` }} />
);
