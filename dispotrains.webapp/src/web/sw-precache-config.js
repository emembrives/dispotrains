module.exports = {
  navigateFallback : '/index.html',
  stripPrefix : 'dist',
  root : 'dist/',
  staticFileGlobs :
      [ 'dist/index.html', 'dist/**.js', 'dist/**.css', 'dist/**.map' ],
  maximumFileSizeToCacheInBytes : 5242880,
  runtimeCaching : [ {
    urlPattern : /^https:\/\/dispotrains\.membrives\.fr\/app\/.*/,
    handler : 'networkFirst'
  },
  {
    urlPattern : /^https:\/\/dispotrains\.membrives\.fr\/static\/data\/.*/,
    handler : 'networkOnly'
  },
  {
    urlPattern : /^https:\/\/dispotrains\.membrives\.fr\/assets\/.*/,
    handler : 'networkFirst'
  }]
};
