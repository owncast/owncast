/* eslint-disable react/no-danger */
import { Html, Head, Main, NextScript } from 'next/document';
import { readFileSync } from 'fs';
import { join } from 'path';
import { HtmlComment } from '../components/common/HtmlComment/HtmlComment';

class InlineStylesHead extends Head {
  getCssLinks: Head['getCssLinks'] = ({ allFiles }) => {
    const { assetPrefix } = this.context;
    if (!allFiles || allFiles.length === 0) return null;
    return allFiles
      .filter((file: any) => /\.css$/.test(file))
      .map((file: any) => (
        <style
          key={file}
          nonce={this.props.nonce}
          data-href={`${assetPrefix}/_next/${file}`}
          dangerouslySetInnerHTML={{
            __html: readFileSync(join(process.cwd(), '.next', file), 'utf-8'),
          }}
        />
      ));
  };
}

export default function Document() {
  return (
    <Html lang="en">
      <InlineStylesHead />
      <body>
        <Main />
        <NextScript />
        {'\n\n'}
        <HtmlComment text="If you're reading this we'd love your help with Owncast! Visit https://owncast.online/help to see how you can contribute to make independent, open source live streaming better with your help." />
      </body>
    </Html>
  );
}
