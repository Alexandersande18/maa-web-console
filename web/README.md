# MAA Web frontend

Vue 3 + Vant frontend for the maintained GoFrame backend.
See [the root README](../README.MD) for installation, screenshots, authentication and deployment.

```sh
npm ci
npm run serve
```

Development URL: http://127.0.0.1:8083. Requests to `/maasst` are proxied to the local backend on port 5244.

```sh
npm run test:task-config
npm run build
```

`prebuild` synchronizes shared task configuration into the backend fallback UI.
`dist/` is generated and excluded from Git. Never place secrets in `VUE_APP_*` variables; they are embedded in the public JavaScript bundle.
