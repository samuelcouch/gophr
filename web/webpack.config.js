var rucksack = require('rucksack-css')
var webpack = require('webpack')
var path = require('path')
var execSync = require('child_process').execSync

var minikubeIP = '127.0.0.1', gophrWebPort = 30443

// First, grab the minikube ip if not in production.
if (process.env.NODE_ENV !== 'production') {
  try {
    console.log('Attempting to get minikube IP address...')
    minikubeIP = execSync('minikube ip', { encoding: 'utf8' }).trim()
    console.log('Got minikube IP address, now starting webpack...')
  } catch(err) {
    console.error(
      `Failed to read the minikube IP address. ` +
      `Make sure the gophr development environment is running: ${err}.`)
    process.exit(1)
  }
}

module.exports = {
  context: path.join(__dirname, './client'),
  entry: {
    jsx: './index.js',
    html: './index.html',
    vendor: [
      'react',
      'react-dom',
      'react-redux',
      'react-router',
      'react-router-redux',
      'redux'
    ]
  },
  output: {
    path: path.join(__dirname, './static'),
    publicPath: '/static/',
    filename: 'bundle.js',
  },
  module: {
    loaders: [
      {
        test: /\.html$/,
        loader: 'file?name=[name].[ext]'
      },
      {
        test: /\.(svg|eot|ttf|woff|woff2)$/,
        loader: 'url-loader?limit=100000'
      },
      {
        test: /\.css$/,
        include: /client/,
        loaders: [
          'style-loader',
          'css-loader?modules&sourceMap&importLoaders=1&localIdentName=[local]___[hash:base64:5]',
          'postcss-loader'
        ]
      },
      {
        test: /\.css$/,
        exclude: /client/,
        loader: 'style!css'
      },
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        loaders: [
          'react-hot',
          'babel-loader'
        ]
      },
    ],
  },
  resolve: {
    extensions: ['', '.js', '.jsx']
  },
  postcss: [
    rucksack({
      autoprefixer: true
    })
  ],
  plugins: [
    new webpack.optimize.CommonsChunkPlugin('vendor', 'vendor.bundle.js'),
    new webpack.DefinePlugin({
      'process.env': { NODE_ENV: JSON.stringify(process.env.NODE_ENV || 'development') }
    })
  ],
  devServer: {
    contentBase: './client',
    hot: true,
    proxy: {
      '/api/*': {
        target: `https://${minikubeIP}:${gophrWebPort}`,
        secure: false
      }
    }
  }
}
