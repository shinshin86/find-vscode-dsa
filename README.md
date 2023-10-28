# Find VSCode DSA

Find the VSCode project you `don't show again`.


## Development

### Windows

Windows environment requires setup before using go-sqlite3

Reference: https://github.com/mattn/go-sqlite3#windows

Before starting the application, please download and install tdm-gcc-webdl.exe from [this page](https://jmeubank.github.io/tdm-gcc/download/).

------------------

# README

## About

This is the official Wails React-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
