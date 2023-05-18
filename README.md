# foxyshot

![build](https://github.com/elnoro/foxyshot/workflows/go/badge.svg) [![go report](https://goreportcard.com/badge/github.com/elnoro/foxyshot)](https://goreportcard.com/report/github.com/elnoro/foxyshot)

A lightweight tool to upload MacOS screenshots to an S3-compatible provider. 

## Install 
Install with brew:
```
$ brew tap elnoro/tap
$ brew install foxyshot
```

Install from source:

```
$ git clone https://github.com/elnoro/foxyshot.git
$ cd foxyshot && make install
```

## Configure

1. Change the default MacOS screenshot location to a designated folder, e. g. `~/Desktop/Screenshots`
2. Run `foxyshot configure` (it creates a config file in ~/.config/foxyshot/config.json; see the format [here](https://github.com/elnoro/foxyshot/blob/master/config/testdata/full.json)). For S3 credentials, refer to your S3 provider.

## Run

### brew services (program starts via launchctl)
```
$ brew services start foxyshot
```
```
$ brew services stop foxyshot
```


### Manually 
```
$ foxyshot start
```
The program starts in the background. To stop it, run:
```
$ foxyshot stop
```

## Known issues

If you decide to keep the original screenshot files (setting "removeOriginals" to false), on MacOS you will eventually run into a "too many open files" error.

At this point, either set a higher ulimit or remove the old files manually.

This is because of kqueue, see more technical details [here](https://github.com/fsnotify/fsnotify/issues/11#issuecomment-1279133120).