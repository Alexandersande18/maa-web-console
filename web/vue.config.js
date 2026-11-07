const { defineConfig } = require('@vue/cli-service')

// 使用 Webpack 版本的 unplugin-vue-components
const Components = require('unplugin-vue-components/webpack').default
const { VantResolver } = require('unplugin-vue-components/resolvers')

module.exports = defineConfig({
  transpileDependencies: true,
  configureWebpack: {
    plugins: [
      Components({
        resolvers: [VantResolver()],
      }),
    ],
  },
  devServer: {
    host: process.env.DEV_SERVER_HOST || '127.0.0.1',
    port: Number(process.env.DEV_SERVER_PORT || 8083),
    allowedHosts: process.env.DEV_SERVER_ALLOWED_HOSTS === 'all'
      ? 'all'
      : (process.env.DEV_SERVER_ALLOWED_HOSTS
        ? process.env.DEV_SERVER_ALLOWED_HOSTS.split(',').map(host => host.trim())
        : 'auto'),
    proxy: {
      '/maasst': {
        target: 'http://localhost:5244',
        changeOrigin: true,
        ws: true
      }
    }
  }
})
