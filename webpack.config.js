// const path = require('path');
// const webpack = require('webpack');
// const CopyWebpackPlugin = require('copy-webpack-plugin');
// const CleanWebpackPlugin = require('clean-webpack-plugin');
// const CustomPlugin = require('custom-plugin');

module.exports.getWebpackConfig = (config, options) => {

  let newConf = {
    ...config,
  };

  // newConf.externals.reverse();
  // newConf.externals.push('aws-sdk');
  // newConf.externals.reverse();
  // newConf.externals.push((context, request, callback) => {
  //   console.log('context', context);
  //   console.log('request', request);
  //   callback();
  // });
  newConf.module.rules.push({
    type: 'javascript/auto',
    test: /\.json$/,
    use: 'json-loader'
  });

  newConf.resolve.extensions.push('.json')

  console.log('webpack config', newConf);

  return newConf;
};

// module.exports = {
//   node: {
//     fs: 'empty',
//   },
//   context: path.join(__dirname, 'src'),
//   entry: {
//     module: './module.ts',
//   },
//   devtool: 'source-map',
//   output: {
//     filename: '[name].js',
//     path: path.join(__dirname, 'dist'),
//     libraryTarget: 'amd',
//   },
//   externals: [
//     'aws-sdk',
//     function(context, request, callback) {
//       var prefix = 'grafana/';
//       if (request.indexOf(prefix) === 0) {
//         return callback(null, request.substr(prefix.length));
//       }
//       callback();
//     },
//   ],
//   plugins: [
//     new CleanWebpackPlugin.CleanWebpackPlugin({}),
//     new webpack.optimize.OccurrenceOrderPlugin(true),
//     new CopyWebpackPlugin({
//       patterns:[
//         {from: 'plugin.json', to: '.'},
//         {from: '../README.md', to: '.'},
//         // {from: '../CHANGELOG.md', to: '.'},
//         {from: '../LICENSE', to: '.'},
//         // {from: 'partials/*', to: '.'},
//         // {from: 'images/*', to: '.'},
//         // {from: 'css/*', to: '.'},
//         // {from: 'data/*', to: '.'},
//       ]
//     }),
//   ],
//   resolve: {
//     extensions: ['.ts', '.js'],
//   },
//   module: {
//     rules: [
//       {
//         test: /\.(png|jpg|gif|svg|ico)$/,
//         loader: 'file-loader'
//         // query: {
//         //   outputPath: './images/',
//         //   name: '[name].[ext]',
//         // },
//       },
//       {
//         test: /\.tsx?$/,
//         use: 'ts-loader',
//         exclude: /node_modules/
//       },
//       // {
//       //   test: /\.tsx?$/,
//       //   loaders: [
//       //     {
//       //       loader: 'babel-loader',
//       //       options: { presets: ['env'] },
//       //     },
//       //     'ts-loader',
//       //   ],
//       //   exclude: /(node_modules)/,
//       // },
//       {
//         test: /\.css$/,
//         use: [
//           {
//             loader: 'style-loader',
//           },
//           {
//             loader: 'css-loader',
//             options: {
//               importLoaders: 1,
//               sourceMap: true,
//             },
//           }
//         ],
//       },
//     ],
//   },
// };
