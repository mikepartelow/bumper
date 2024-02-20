# bumper

`bumper` bumps image versions in a text file to their latest version in a given registry.

## Usage

```bash
% make local 
% echo "here is an image url: ghcr.io/mikepartelow/bumper:latest" > file.txt
% ./bumper ghcr.io/mikepartelow/ file.txt
here is an image url: ghcr.io/mikepartelow/bumper:main.v1.0.0@sha256:929ed1ece9808d842389cb4afea961a8e2422514
```

## How is this useful?

You should [pin](https://docs.docker.com/engine/reference/commandline/image_pull/#pull-an-image-by-digest-immutable-identifier) your container images, and you should frequently update your container dependencies so you benefit from security updates.

`bumper` helps with both, by always converting unpinned image URLs to pinned ones, and by update image URLs to use the latest available image.

`bumper` makes some assumptions, namely, it only considers tags that begin with `main.` when choosing the "latest" image.
