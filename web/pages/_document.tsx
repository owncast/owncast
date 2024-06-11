/* eslint-disable react/no-danger */
import { readFileSync } from 'fs';
import { join } from 'path';
import Document, { Html, Head, Main, NextScript } from 'next/document';
import type { DocumentContext } from 'next/document';
import { createCache, extractStyle, StyleProvider } from '@ant-design/cssinjs';

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

const Doc = () => (
  <Html lang="en">
    <InlineStylesHead />
    <body>
      <Main />
      <NextScript />
    </body>
  </Html>
);

/*

Yoinked from https://ant.design/docs/react/use-with-next#using-pages-router. Change if someone figures out a way to
mesh this with the InlineStyles above.

*/
Doc.getInitialProps = async (ctx: DocumentContext) => {
  const cache = createCache();
  const originalRenderPage = ctx.renderPage;
  ctx.renderPage = () =>
    originalRenderPage({
      enhanceApp: App => props => (
        // layer forces antd to generate new css so that the cache can be safely removed
        <StyleProvider cache={cache} layer autoClear>
          <App {...props} />
        </StyleProvider>
      ),
    });

  const initialProps = await Document.getInitialProps(ctx);
  const style = extractStyle(cache, true);
  return {
    ...initialProps,
    styles: (
      <>
        {initialProps.styles}
        <style id="antd-init-cache" dangerouslySetInnerHTML={{ __html: style }} />
      </>
    ),
  };
};

export default Doc;
