module.exports = {
  stripPrefix : 'dist',
  root : 'dist/',
  staticFileGlobs :
      [ 'dist/index.html', 'dist/**.js', 'dist/**.css', 'dist/**.map' ],
  maximumFileSizeToCacheInBytes : 5242880,
  runtimeCaching : [ {
    urlPattern : /^https:\/\/dispotrains\.membrives\.fr\/app\/.*/,
    handler : 'networkFirst'
  } ]
};
