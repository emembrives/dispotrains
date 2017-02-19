module.exports = {
  navigateFallback : '/index.html',
  stripPrefix : 'dist',
  root : 'dist/',
  staticFileGlobs :
      [ 'dist/index.html', 'dist/**.js', 'dist/**.css', 'dist/**.map' ],
  maximumFileSizeToCacheInBytes : 5242880,
  runtimeCaching : [ {
    urlPattern : /^http:\/\/dispotrains\.membrives\.fr\/app\/Get/,
    handler : 'networkFirst'
  } ]
};
