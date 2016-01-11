mkdir debroot
cp debian/* debroot -r
mkdir debroot/usr/bin -p
cp gocollect debroot/usr/bin
fpm \
  -s dir \
  -t deb \
  --force \
  -C "debroot" \
  -a "amd64" \
  --version "0.0.1" \
  .
