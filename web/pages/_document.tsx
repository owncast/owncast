/* eslint-disable react/no-danger */
import { readFileSync } from 'fs';
import { join } from 'path';
import { Html, Head, Main, NextScript } from 'next/document';

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
      </body>
    </Html>
  );
}
