// Load the modules which are installed through NPM.
var gulp = require('gulp');
var browserify = require('browserify');
var source = require('vinyl-source-stream');
var buffer = require('vinyl-buffer');
var uglify = require('gulp-uglify');
var rename = require('gulp-rename');
var sourcemaps = require('gulp-sourcemaps');
var gutil = require('gulp-util');
var gulpif = require('gulp-if');
var jshint= require('gulp-jshint');

// Define paths.
var paths = {
    js: ['./app/franz.js'],
    dest: './static/js/'
};

var isProd = false;

// If "--env=prod" is passed from the command line then update the flags
if (gutil.env.env === 'prod') {
    isProd = true;
}

var isDev = !isProd;

// The JS task.
gulp.task('js', function () {
    var b = browserify({
        entries: paths.js,
        debug: true
    });
    return b.bundle()
        .pipe(source('bundle.js'))
        .pipe(jshint())
        .pipe(jshint.reporter('default'), { verbose: true })
        .pipe(jshint.reporter('fail'))
        .pipe(gulpif(isDev, gulp.dest(paths.dest)))
        .pipe(buffer())
        .pipe(gulpif(isProd, uglify()))
        .pipe(gulpif(isProd, gulp.dest(paths.dest)));
});

// Rerun tasks whenever a file changes.
gulp.task('watch', function () {
    gulp.watch(paths.js, ['js']);
});

// The default task (called when running `gulp` from cli).
gulp.task('default', ['watch', 'js']);