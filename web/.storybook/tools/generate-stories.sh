#!/bin/sh

# Generate the custom Emoji story
node generate-emoji-story.mjs >../stories-category-doc-pages/Emoji.stories.mdx

# Generate stories out of documentation

# Pull down the doc about development
curl -s https://raw.githubusercontent.com/owncast/owncast.github.io/master/content/development.md >/tmp/development.md
node generate-document-stories.mjs

# Project image assets

node generate-image-story.mjs ../../public/img/ Images "owncast/Frontend Assets/Images" "img" >../stories-category-doc-pages/Images.stories.mdx
node generate-image-story.mjs ../../public/img/platformlogos/ "Social Platform Images" "owncast/Frontend Assets/Social Platform Images" "img/platformlogos" >../stories-category-doc-pages/SocialPlatformImages.stories.mdx
node generate-image-story.mjs ../story-assets/project/ "Logos & Graphics" "owncast/Project Assets/Logos & Graphics" "project" --large >../stories-category-doc-pages/LogosAndGraphics.stories.mdx
node generate-image-story.mjs ../story-assets/tshirt/ "T-shirt" "owncast/Project Assets/T-Shirt" "tshirt" --large >../stories-category-doc-pages/Tshirt.stories.mdx
