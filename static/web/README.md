TLDR: You can't make web changes here.

This directory contains the compiled, bundled, web application to be distributed.

You should not try to edit any files under `/static/web`, instead look at the `/web` directory to make changes to web source code.

After changes to `/web` get merged into `develop` it will automatically get compiled, and bundled, and will update `/static/web` properly.
