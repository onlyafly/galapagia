var gulp = require('gulp');
var plumber = require('gulp-plumber'); // Prevents build errors from killing gulp
var browserify = require('browserify');
var babelify = require('babelify');
var vinylSourceStream = require('vinyl-source-stream');
var gulpUtil = require('gulp-util');

// ---------- scripts

var handleBrowserifyErrors = function () {
  var args = Array.prototype.slice.call(arguments);
  gulpUtil.log(args);
  this.emit('end'); // Keep gulp from hanging on this task
};

gulp.task('scripts:compile', function () {
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

gulp.task('scripts', ['scripts:compile'], function () {
  gulp.watch('web/scripts/**/*.js*', ['scripts:compile']);
});

// ---------- default

gulp.task('default', ['scripts']);
