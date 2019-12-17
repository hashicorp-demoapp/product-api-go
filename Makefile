# Run tests continually with  a watcher
autotest:
	filewatcher --idle-timeout 24h -x **/functional_tests gotestsum --format standard-verbose