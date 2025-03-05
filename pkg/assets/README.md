# Assets

This package contains the assets used in the frontend.


## How it's built

The assets are built using [esbuild](https://esbuild.github.io/). In development, each time an asset it requested it is rebuilt and sent.
In production, before building the reviso binary a build step is ran by calling --build on the reviso main package: `go run cmd/reviso/main.go --build`.
This will build the assets and put them into the `asssets/src/dist` directory. Then when the binary is built it that directory is embedded into the
binary. So when the binary is ran, the assets do not need to be built on request and can be loaded from the `assets/src/dist` directory.

```
developemnt -> esbuild -> js/css file

build -> esbuild -> js/css file stored in assets/src/dist
production -> embedded fs (assets/src/dist) -> js/css file
```

To see how the assets are build see the [assets.go](./assets.go) file.


## Adding new shadcn-ui components

You'll need to cd into the `pkg/assets/src` directory and run the add command with the path:

```sh
(cd pkg/assets/src && npx shadcn-ui@latest add button --path app/components/ui/)
```

## Graphql queries

When you make/change a graphql query, you'll need to run the `codegen` command:

```sh
make gen_web
```

> Note: this will run the codegen for the frontend/web and the pkg/assets/src directories.


## Todo

- [ ] Right now when you request an asset, it rebuilds ALL the assets. This is not ideal. We should only be rebuilding the assets that was requested.
- [ ] Getting the assets to build is a little funky, would be better if there was an actual assets.Build() function rather than building the HttpHandler.

