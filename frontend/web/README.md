This is a [Next.js](https://nextjs.org/) project bootstrapped with [`create-next-app`](https://github.com/vercel/next.js/tree/canary/packages/create-next-app).

## Getting Started Now

Get your PAT (classic) from [here](https://github.com/settings/tokens) you only need to allow:

- `read:packages`

Add the following to your `.zshrc` or `.bashrc` file:

```
export GH_NPM_TOKEN={Your GitHub Personal Access Token (classic)}
```

Then, install your deps:

```bash
bun install
```

Then start the dev server:

```bash
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Graphql

We're using [GraphQL Code Generator](https://the-guild.dev/graphql/codegen/docs/getting-started) to generate type from our graph.

Queries can be kept in any .tsx file. For example: [src/app/(app)/documents/page.tsx](</frontend/web/src/app/(app)/documents/page.tsx>). Then run:

```bash
bun run gen
```

Which will inspect the graph running on http://localhost:3000/graphql and generate the types in the gql directory.

## Rogue

### Working with a local copy of Rogue

Update your package.json with:

```json
"dependencies": {
    "@teamreviso/rogue": "file:../../../rogue/javascript/"
```

### Build Rogue

In the Rogue directory:

```bash
npm run build
```

### Install Rogue:

In the web directory:

```bash
bun install
```

### Clear the build cache and run dev:

In the web directory:

```bash
rm -rf .next && bun run dev
```

> Each time you change Rogue, you need to build rogue, clear the build cache and run dev.
