import { defineConfig } from '@umijs/max';

export default defineConfig({
  publicPath: '/',
  title: 'agentic-chain',
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: {},
  layout: {
    // title: '@umijs/max',
  },
  define:{
    'process.env.BASE_URL' : process.env.BASE_URL
  },
  routes: [
    {
      path: '/',
      redirect: '/home',
    },
    {
      name: 'home',
      path: '/home',
      component: './Home',
    },
    {
      name: 'proposalDetail',
      path: '/proposalDetail/:id',
      component: './Proposal/proposalDetail.tsx',
    },
    {
      name: 'agentDetail',
      path: '/agentDetail/:id',
      component: './Agent/agentDetail.tsx',
    }
  ],
  npmClient: 'yarn',
  proxy: {
    '/api': {
      target: 'http://52.221.224.189:8631', // 你的后端域名
      changeOrigin: true,
      // pathRewrite: { '^/api': '' }, // 根据需求决定是否重写路径
    },
  },
  // chainWebpack(config) {
  //   // 清空所有默认图片处理规则
  //   ['png','jpg','jpeg','gif','webp','svg'].forEach(ext => {
  //     config.module.rules.delete(ext);
  //   });
  //
  //   // 统一配置图片处理
  //   config.module
  //     .rule('image')
  //     .test(/\.(png|jpe?g|gif|webp|svg)$/i)
  //     .use('url-loader')
  //     .loader(require.resolve('url-loader'))
  //     .options({
  //       limit: 4096, // 根据需求调整
  //       esModule: false, // 关键参数
  //       fallback: {
  //         loader: require.resolve('file-loader'),
  //         options: {
  //           name: 'static/[name].[hash:8].[ext]',
  //           esModule: false
  //         }
  //       }
  //     });
  // }
});

