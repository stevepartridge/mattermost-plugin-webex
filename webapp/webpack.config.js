var path = require('path');

module.exports = {
    entry: [
        './src/index.js',
    ],
    resolve: {
        modules: [
            'src',
            'node_modules',
        ],
        extensions: ['*', '.js', '.jsx'],
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules|external/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['env', 'react'],
                        plugins: [
                            'transform-class-properties',
                            'transform-object-rest-spread',
                        ],
                    },
                },
            },
            {
                test: /\.(png|svg|jpg|gif)$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: 'images/[name].[ext]',
                            publicPath: '/static/plugins/webex/'
                        }
                    }
                ]
            },
            {
                test: /\.(js|css)$/,
                include: /external/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: 'external/[name].[ext]',
                            publicPath: '/static/plugins/webex/'
                        }
                    }
                ]
            }
        ],
    },
    externals: {
        react: 'React',
        redux: 'Redux',
        'react-redux': 'ReactRedux',
    },
    output: {
        path: path.join(__dirname, '/dist'),
        publicPath: '/',
        filename: 'main.js'
        // require.context("./images/", true, /\.(png|svg|jpg|gif)$/)
    },
};

// require.context("../images/", true, /\.(png|svg|jpg|gif)$/);
