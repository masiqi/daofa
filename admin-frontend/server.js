const express = require('express');
const http = require('http');
const path = require('path');
const reload = require('reload');
const bodyParser = require('body-parser');
const logger = require('morgan');
const { createProxyMiddleware } = require('http-proxy-middleware');
const jwt = require('jsonwebtoken');
const cookieParser = require('cookie-parser');

const app = express();

app.set('port', process.env.PORT || 3000);
app.set('host', process.env.HOST || '0.0.0.0');
app.use(logger('dev'));
app.use(bodyParser.json());
app.use(cookieParser());

// 中间件：检查 JWT token
const checkToken = (req, res, next) => {
    const token = req.cookies.jwt;
    if (token) {
        jwt.verify(token, process.env.JWT_SECRET, (err, decoded) => {
            if (err) {
                return res.redirect('/login');
            }
            req.decoded = decoded;
            next();
        });
    } else {
        res.redirect('/login');
    }
};

// 静态文件服务
app.use(express.static(path.join(__dirname, 'public')));

app.get('/login', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'login.html'));
});

app.get('/', checkToken, (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// API 代理
app.use('/api', createProxyMiddleware({
    target: 'http://localhost:8080',
    changeOrigin: true,
    pathRewrite: {
        '^/api': '/admin'
    }
}));

// 捕获所有其他路由
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

const server = http.createServer(app);

// Reload code here
reload(app)
    .then(function (reloadReturned) {
        server.listen(app.get('port'), app.get('host'), function () {
            console.log(
                'Admin frontend server running at http://' + app.get('host') + ':' + app.get('port')
            );
        });
    })
    .catch(function (err) {
        console.error(
            'Reload could not start, could not start server/sample app',
            err
        );
    });