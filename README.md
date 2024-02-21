# monoize

`monoize` makes your Git repositories monorepo.

## Installation

### Build from source

1. Clone the repository

   ```sh
   $ git clone https://github.com/0x1306e6d/monoize.git monoize
   $ cd monoize
   ```

2. Build and install

   ```sh
   $ go install github.com/0x1306e6d/monoize
   ```

   See the [Go Modules Reference](https://go.dev/ref/mod#go-install) for more
   information.

## Usage

```sh
$ monoize [<options>] <source>... <target>
```

Merges the specified `<source>` repositories into a single monorepo in the
`<target>` directory. If the target directory is already a Git repository,
commits will be appended.

All commits from the source repositories are applied to the target repository
in the order of the author date.

Files in each source are placed in a subdirectory named the same as the
repository name. You can specify the name by appending `>><name>`.
For example, `monoize https://github.com/0x1306e6d/monoize>>mono tools` will
place the files of the `monoize` repository in the `tools/mono` directory.

### Options

#### -f, --force

Deletes the target directory before merging.

## How it works

`monoize` uses `git format-patch` to export the commits of the source
repositories. And it sorts all patches by the author date and applies them to
the target repository by `git am`.

## License

```
MIT License

Copyright (c) 2024 Gihwan Kim

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
