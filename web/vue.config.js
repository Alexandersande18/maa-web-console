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
    host: '127.0.0.1',
    port: 8083,
    proxy: {
      '/maasst': {
        target: 'http://localhost:5244',
        changeOrigin: true,
        ws: true
      }
    }
  }
})
