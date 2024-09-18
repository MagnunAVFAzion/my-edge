"use strict";

const next = require('next');

const isDev = false;

const nextConfig = require('./next.config');
const nextApp = next({
  dev: isDev,
  dir: __dirname,
  conf: nextConfig,
});

const handle = nextApp.getRequestHandler();

module.exports = async function (context, callback) {
  nextApp.prepare().then(() => handle(context.request, context.response));
}
