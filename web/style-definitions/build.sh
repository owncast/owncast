#!/bin/sh

mv build/variables.css ../styles/variables.css
mv build/variables.less ../styles/theme.less

# Served plugin stylesheet: the generated :root token block (build/plugin.css)
# followed by the hand-authored element baseline. One file authors <link>,
# placed under public/ so the Owncast server serves it at /styles/plugin.css.
cat templates/plugin-elements.css >> build/plugin.css
mv build/plugin.css ../public/styles/plugin.css
