const CleanWebpackPlugin = require('clean-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const path = require('path')
const webpack = require('webpack')

// Plugin to convert SASS into CSS
const extractSass = new ExtractTextPlugin({
  filename: 'css/[name]-[contenthash].css',
  disable: process.env.NODE_ENV !== 'production'
})

// List of plugins, which can be modified depending on build environment
var plugins = [
  new CleanWebpackPlugin(['public'], { root: __dirname }),
  new webpack.ProvidePlugin({

  }),
  new webpack.DefinePlugin({
    VERSION: JSON.stringify(require('./package.json').version)
  }),
  new webpack.optimize.CommonsChunkPlugin({
    names: ['vendor'],
    filename: 'js/[name]-[id]-[hash].js'
  }),
  new HtmlWebpackPlugin({
    title: 'Gringotts',
    filename: path.join(__dirname, 'public', 'index.html'),
    template: path.join(__dirname, 'index.html'),
    chunks: ['vendor', 'bundle']
  }),
  extractSass
]

// Do additional stuff for production environment
if (process.env.NODE_ENV === 'production') {
  // Uglify Javascript code
  plugins.push(new webpack.optimize.UglifyJsPlugin())
}

module.exports = {
  entry: {
    bundle: path.join(__dirname, 'app.jsx'),
    vendor: path.join(__dirname, 'vendor.js')
  },
  target: 'electron-renderer',
  output: {
    filename: 'js/[name]-[hash].js',
    path: path.resolve(__dirname, 'public'),
    publicPath: ''
  },
  module: {
    rules: [{
      test: /\.jsx?$/,
      exclude: /node_modules/,
      use: {
        loader: 'babel-loader',
        options: {
          presets: ['env', 'react']
        }
      }
    }, {
      test: /\.scss$/,
      use: extractSass.extract({
        use: [{
          loader: 'css-loader',
          options: {
            minimize: process.env.NODE_ENV === 'production'
          }
        }, {
          loader: 'sass-loader'
        }],
        fallback: 'style-loader'
      })
    }]
  },
  plugins
}

