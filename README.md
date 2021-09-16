# pcabinet
A tool for capturing and organizing golang profiles

## Usage
1. Edit `pcabinet_config.yml` to contain the base `debug/pprof` URLs for your desired services (samples included
2. (optional) copy this value into `$HOME/.config/pcabinet/pcabinet_config.yml` for the ability to run anywhere with `go install`
2. `go run .` to start up the interface
3. select your service from the presented list
4. select the profile(s) you'd like to capture from the presented list
5. (optional) give a description of what's different with this profile
6. your profile is captured and stored into the `$NAME/$NAME.$DATE.$DESC.$TYPE` file

![pcabinet](recording.gif?raw=true)