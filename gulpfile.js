var gulp = require('gulp');
var plumber = require('gulp-plumber'); // Prevents build errors from killing gulp
var browserify = require('browserify');
var babelify = require('babelify');
var vinylSourceStream = require('vinyl-source-stream');
var gulpUtil = require('gulp-util');
var child      = require('child_process'); // Node.js module
var reload     = require('gulp-livereload');
var sync       = require('gulp-sync')(gulp).sync;

// ---------- local variables

var server = null;

// ---------- server
// Inspired by http://struct.cc/blog/2015/05/08/building-web-applications-in-golang-with-gulp-and-livereload/

// Build application server
gulp.task('server:build', function() {
  // Build is performed synchronously via spawnSync, as we have to wait for
  // the build to complete before restarting the application server.
  result = child.spawnSync('go', ['install', './...']);

  if (result.status !== 0) {
    gulpUtil.log(gulpUtil.colors.red(result.stderr));
    gulpUtil.beep();
  }

  return result;
});

// Restart the application server - The spawn task checks if the server is
// currently running, terminates it by sending SIGTERM to the child process and
// restarts it. When the server sends the first log output (ideally something like
// Listening on 127.0.0.1:3000), a reload is triggered and the browser will
// refresh instantly.
gulp.task('server:spawn', function() {
  if (server) {
    server.kill();
  }

  // Spawn application server
  server = child.spawn('galapagia');

  // Trigger LiveReload upon server start
  server.stdout.once('data', function() {
    reload.reload('/'); // The url to reload to
  });

  // Pretty print server log output
  server.stdout.on('data', function(data) {
    var lines = data.toString().split('\n');
    for (var l in lines) {
      if (lines[l].length) {
        gulpUtil.log(lines[l]);
      }
    }
  });

  // Print errors to stdout
  server.stderr.on('data', function(data) {
    process.stdout.write(data.toString());
  });
});

// Watch source for changes and restart application server.
gulp.task('server:watch', function() {

  // Restart application server
  gulp.watch([
    'web/templates/**/*.tmpl'
  ], ['server:spawn']);

  // Rebuild and restart application server
  gulp.watch([
    '*/**/*.go',
  ], sync([
    'server:build',
    'server:spawn'
  ], 'server'));
});

// ---------- scripts

var handleBrowserifyErrors = function () {
  var args = Array.prototype.slice.call(arguments);
  gulpUtil.log(gulpUtil.colors.red("Error during browserify"), args);
  gulpUtil.beep();
  this.emit('end'); // Keep gulp from hanging on this task
};

gulp.task('scripts:build', function () {
  var b = browserify(
    {extensions: ['.jsx']}
  );

  // babelify transform turns ES6 into ES5
  b.transform("babelify", {presets: ["es2015", "react"]});

  b.add('web/scripts/app.js');

  return b.bundle()
    .on('error', handleBrowserifyErrors)
    .pipe(plumber())
    .pipe(vinylSourceStream('bundle.js')) // desired output filename
    .pipe(gulp.dest('web/public/js')); // desired output location; file is not minified
});

gulp.task('scripts:watch', function () {
  gulp.watch('web/scripts/**/*.js*', ['scripts:build']);
});

// ---------- default

// Start asset and server watchdogs and initialize livereload.
gulp.task('watch', ['scripts:build','server:build'], function() {
  reload.listen();
  return gulp.start([
    'scripts:watch',
    'server:watch',
    'server:spawn'
  ]);
});

gulp.task('default', ['watch']);
