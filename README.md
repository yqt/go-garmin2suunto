# go-garmin2suunto

**UNSTABLE**

## Usage notice

Requesting unofficial APIs and simulated website login behaviors can be suspicious to Garmin.

Use it AT YOUR OWN RISK.

## Usage

```
cp config/config.sample config/config.go
# Filling config.go with garmin and suunto(movescount) account info
cd main
go build -ldflags '-s -w' -o garmin2suunto

# Specific port to listen. default is 38080.
PORT=38080 ./garmin2suunto
# Sync latest garmin activities of TODAY(up to 3 activities) to movescount
# It will try to log in to Garmin website in ervery sync process since login session persistence is not implemented.
curl 'http://localhost:38080/api/sync'
```

## Thanks

[garmin2suunto](https://github.com/goodtools/garmin2suunto) (more details inside)

[tapiriik](https://github.com/cpfair/tapiriik)

[garmin-connect](https://github.com/abrander/garmin-connect)

[garminexport](https://github.com/petergardfjall/garminexport)
