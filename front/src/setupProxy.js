const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  const target = process.env.REACT_APP_API_URL || 'http://localhost:8080';
  
  // Proxy all API paths to backend
  app.use(
    ['/auth', '/posts', '/user', '/users', '/admin', '/subscription', '/categories', '/groups'],
    createProxyMiddleware({
      target: target,
      changeOrigin: true,
      logLevel: 'debug',
      onError: (err) => {
        console.log('Proxy error:', err.message);
      },
      onProxyReq: (proxyReq, req) => {
        console.log('Proxying request:', req.method, req.url, '-> ', target + req.url);
      }
    })
  );

  // Proxy uploads to backend
  app.use(
    '/uploads',
    createProxyMiddleware({
      target: target,
      changeOrigin: true,
      onError: (err) => {
        console.log('Upload proxy error:', err.message);
      }
    })
  );
};