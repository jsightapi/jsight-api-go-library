Snapshot checker
================

Run test tests from `/test/` package, checks that JSON output from test is correct, show diff if not, and give an ability to update expected JSON.

Usage
-----

1. Run `go run ./internal/cmd/snapshot_checker`
2. Command will run all test in `/testdata` directory
3. It will search for a first failed snapshot test and print out difference between current and expected JSON.
4. Press 'y + return' if expected JSON should be patched, or 'n + return' if no.
