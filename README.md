# pcabinet
A tool for capturing and organizing golang profiles

# TODO
 - [X] read in yaml of hardcoded URLs
 - [ ] verify it ends in `/debug/pprof`
 - [X] implement TUI to choose hardcoded URL
 - [X] get user input for a desired profile type
 - [X] get user input for a description
 - [ ] invoke with `?debug=0`
 - [X] display TUI spinner
 - [ ] design naming scheme and write to file

# Extra features
 - [ ] For CPU profiles take a 1 second profile first, open it, and verify CPU usage is over 10%
 - [ ] support variants of an entity (e.g. a dev vs prod instance)
 - [ ] intelligently parse URL to ignore URL parameters specified by the user
 - [ ] display estimated time to completion when capturing a profile
 - [ ] allow multiple types to be captured in sequence with one request