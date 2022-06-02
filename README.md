# Sasha
A shell for your s3-compliant storage.

A miniproject developed as a micro tool for the Tools team in Seznam's SCIF, intended to be used with its proprietary S3 storage.

## Description
Sasha is supposed to simulate shell-like environment while browsing through an S3-compliant storage. Getting around while checking the contents of an S3 storage can be tricky, and Sasha is here to make it easier.

### Usage
Sasha looks for the following env vars upon start:
- `AWS_ACCESS_KEY`
- `AWS_SECRET_ACCESS_KEY`
- `S3_ENTRYPOINT` - S3 entrypoint URL

Flag options:
- `e` - S3 entrypoint URL. Example value: `https://prop.s3.loc.dom.com`. If not specified, the entrypoint will be sought in the env.

You can find some pre-compiled binaries in this repository.

### Supported commands
- `cd` (one level at the time, including `cd ..`, or naked to go to root)
- `ls`
- `file` - displays the file's basic info + metadata
- `pwd`
- `update`/`reload` - reloads the current level
- `exit`  

...more to be added in the future.