## A somewhat gauche, brute-force approach alternative to slurp watch

An issue with kqueue was causing [slurp-contrib/watch](https://github.com/slurp-contrib/watch) to only fire for the first change in OS X and Ubuntu linux. This approach is far from elegant (it polls every file once per second to look for changes in the modification time or size) it does the job on all tested platforms.

