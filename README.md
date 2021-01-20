# foxyshot

![build](https://github.com/elnoro/foxyshot/workflows/build/badge.svg) [![go report](https://goreportcard.com/badge/github.com/elnoro/foxyshot)](https://goreportcard.com/report/github.com/elnoro/foxyshot)

A lightweight tool to upload MacOS screenshots to an S3-compatible provider. 

## Install 

```
$ git clone https://github.com/elnoro/foxyshot.git
$ cd foxyshot && make install
```

## Configure

1. Change the default MacOS screenshot location to a designated folder, e. g. `~/Desktop/Screenshots`
2. Create a config file and put it in ~/.config/foxyshot/config.json (see the format [here](https://github.com/elnoro/foxyshot/blob/master/config/testdata/full.json). For S3 credentials, refer to your S3 provider.
3. Launch the program: 
```
$ bash -c "nohup sh -c 'foxyshot' &"
```
