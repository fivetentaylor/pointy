# Rogue WASM

The build version of this project is stored in the public dir of the web package.

To work on this locally run:

```
make watch_wasm
```

That will rebuild the wasm package and watch for changes.

Also if you want your editor's LSP to work correctly with the wasm package you might need to set some environment variables:

```
export GOOS=js
export GOARCH=wasm
```
