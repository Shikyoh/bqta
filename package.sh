# Variables:
# NAME - name of the service/package without f8- prefix set in Makefile
# DEB_DIR - deb_$NAME
#
# Functions:
# mkdebdir - mkdir -p $DEB_DIR$1 if it does not exist

mkdebdir /usr/local/bin
mkdebdir /etc/$NAME
mkdebdir /etc/init

cp ./bin/$NAME $DEB_DIR/usr/local/bin/$NAME
