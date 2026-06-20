#!/bin/sh

mv build/variables.css ../styles/variables.css
mv build/variables.less ../styles/theme.less

# Served plugin stylesheet: the generated :root token block (build/plugin.css)
# followed by the hand-authored element baseline. One file authors <link>,
# placed under public/ so the Owncast server serves it at /styles/plugin.css.
cat templates/plugin-elements.css >> build/plugin.css
mv build/plugin.css ../public/styles/plugin.css

# Normalize the generated outputs with Prettier so `npm run build-styles`
# reproduces the committed files exactly. Style Dictionary emits one long line
# per token; Prettier wraps them deterministically. Run from the web/ root so
# Prettier picks up the project config (and so styles/variables.css must not be
# listed in .prettierignore).
(cd .. && npx --yes prettier --write \
  styles/variables.css \
  styles/theme.less \
  public/styles/plugin.css)
