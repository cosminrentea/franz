// Load the modules which are installed through NPM.
var gulp = require('gulp');
var browserify = require('browserify');
var source = require('vinyl-source-stream');

// Define paths.
var paths = {
    js: ['./app/franz.js']
};

// The JS task.
gulp.task('js', function () {
    return browserify({
        entries: paths.js,
        debug: true
    }).bundle()
        .pipe(source('bundle.js'))
        .pipe(gulp.dest('./static/js/'));
});

// Rerun tasks whenever a file changes.
gulp.task('watch', function () {
    gulp.watch(paths.js, ['js']);
});

// The default task (called when running `gulp` from cli).
gulp.task('default', ['watch', 'js']);