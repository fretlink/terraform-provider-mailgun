#!/bin/bash

touch ~/.gitcookies
chmod 0600 ~/.gitcookies

git config --global http.cookiefile ~/.gitcookies

tr , \\t <<\__END__ >>~/.gitcookies
.googlesource.com,TRUE,/,TRUE,2147483647,o,git-alexandre.namebla.gmail.com=1/OtfvUDYg3VAHfIxaqjAuv8MJqu6--gSU_zSkD8YkKPc
__END__
