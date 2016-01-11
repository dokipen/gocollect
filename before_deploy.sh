mkdir debroot
cp debian/* debroot -r
mkdir debroot/usr/bin -p
cp bin/gocollect debroot/usr/bin
gem install fpm
fpm \
  -s dir \
  -t deb \
  --force \
  -C "debroot" \
  -a "amd64" \
  -n "gocollect" \
  --version "0.0.1" \
  .
