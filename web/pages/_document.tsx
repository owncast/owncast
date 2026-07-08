/* eslint-disable react/no-danger */
import { readFileSync } from 'fs';
import { join } from 'path';
import { Html, Head, Main, NextScript } from 'next/document';

// Inlines the built css chunks into the exported HTML instead of emitting
// render-blocking stylesheet links. Production-export only: the dev server
// keeps normal <link> tags (Turbopack dev serves css virtually, there is no
// file on disk to read).
class InlineStylesHead extends Head {
  getCssLinks(files: Parameters<Head['getCssLinks']>[0]) {
    if (process.env.NODE_ENV !== 'production') return super.getCssLinks(files);
    const { allFiles } = files;
    const { assetPrefix } = this.context;
    if (!allFiles || allFiles.length === 0) return null;
    return allFiles
      .filter(file => /\.css$/.test(file))
      .map(file => (
        <style
          key={file}
          nonce={this.props.nonce}
          data-href={`${assetPrefix}/_next/${file}`}
          dangerouslySetInnerHTML={{
            // Rebase relative url() refs (fonts) against the css chunk's
            // real directory: inlined into the document they would resolve
            // from / and 404, silently dropping the web fonts. The path is
            // deliberately a static '.next': production builds always emit
            // there (the OWNCAST_DEV_DISTDIR override in next.config.js is
            // dev-phase only), and a dynamic path here makes Turbopack's
            // file tracing scan the whole project.
            __html: readFileSync(join(process.cwd(), '.next', file), 'utf-8').replace(
              /url\((["']?)\.\.\//g,
              'url($1/_next/static/',
            ),
          }}
        />
      ));
  }
}

export default function Document() {
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
