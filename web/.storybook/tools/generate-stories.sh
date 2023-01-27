#!/bin/sh

# Generate the custom Emoji story
node generate-emoji-story.mjs >../stories-category-doc-pages/Emoji.stories.mdx

# Generate stories out of documentation

# Pull down the doc about development
curl -s https://raw.githubusercontent.com/owncast/owncast.github.io/master/content/development.md >/tmp/development.md
node generate-document-stories.mjs

# Project image assets

node generate-image-story.mjs ../../public/img/ Images >../stories-category-doc-pages/Images.stories.mdx
node generate-image-story.mjs ../../public/img/platformlogos/ "Social Platform Images" >../stories-category-doc-pages/SocialPlatformImages.stories.mdx
