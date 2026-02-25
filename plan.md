# catalog/catalogs support

- `pnpm`: supports `catalog` + `catalogs`.
- `bun`: supports `catalog` + `catalogs`.
- `yarn`: supports catalogs.
- `npm`: no catalog support.

Decision: implement `catalog`/`catalogs` for pnpm, bun, and yarn; accept npm as input but return a clear unsupported error.
