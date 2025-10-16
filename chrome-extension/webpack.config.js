const path = require('path');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
  entry: {
    popup: './popup/popup.js',
    background: './background/background.js',
    content: './content/content.js',
    options: './options/options.js'
  },
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: '[name].js',
    clean: true
  },
  mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
  devtool: process.env.NODE_ENV === 'production' ? false : 'inline-source-map',
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env']
          }
        }
      }
    ]
  },
  plugins: [
    new CopyWebpackPlugin({
      patterns: [
        { from: 'manifest.json', to: 'manifest.json' },
        { from: 'popup/popup.html', to: 'popup.html' },
        { from: 'popup/popup.css', to: 'popup.css' },
        { from: 'options/options.html', to: 'options.html' },
        { from: 'icons', to: 'icons' },
        { from: 'jest.setup.js', to: 'jest.setup.js' }
      ]
    })
  ],
  resolve: {
    extensions: ['.js']
  }
};