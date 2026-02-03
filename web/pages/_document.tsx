/* eslint-disable react/no-danger, max-classes-per-file */
import { readFileSync, existsSync } from 'fs';
import { join } from 'path';
import Document, { Html, Head, Main, NextScript, DocumentContext } from 'next/document';
import { createCache, extractStyle, StyleProvider } from '@ant-design/cssinjs';

class InlineStylesHead extends Head {
  getCssLinks: Head['getCssLinks'] = ({ allFiles }) => {
    const { assetPrefix } = this.context;
    if (!allFiles || allFiles.length === 0) return null;
    return allFiles
      .filter((file: string) => /\.css$/.test(file))
      .map((file: string) => {
        const filePath = join(process.cwd(), '.next', file);
        if (!existsSync(filePath)) {
          return null;
        }
        return (
          <style
            key={file}
            nonce={this.props.nonce}
            data-href={`${assetPrefix}/_next/${file}`}
            dangerouslySetInnerHTML={{
              __html: readFileSync(filePath, 'utf-8'),
            }}
          />
        );
      })
      .filter(Boolean);
  };
}

export default class MyDocument extends Document {
  static async getInitialProps(ctx: DocumentContext) {
    const cache = createCache();
    const originalRenderPage = ctx.renderPage;

    ctx.renderPage = () =>
      originalRenderPage({
        enhanceApp: App =>
          function EnhancedApp(props) {
            return (
              <StyleProvider cache={cache}>
                <App {...props} />
              </StyleProvider>
            );
          },
      });

    const initialProps = await Document.getInitialProps(ctx);
    const antdStyle = extractStyle(cache, true);

    return {
      ...initialProps,
      styles: (
        <>
          {initialProps.styles}
          <style dangerouslySetInnerHTML={{ __html: antdStyle }} />
        </>
      ),
    };
  }

  render() {
    return (
      <Html>
        <InlineStylesHead />
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}
