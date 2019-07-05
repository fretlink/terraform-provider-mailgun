#!/usr/bin/env bash

(sleep 10;
wget --recursive --no-parent --domains localhost --page-requisites --html-extension --convert-links --no-clobber http://localhost:4567/docs/providers/mailgun/;

rm -rf terraform-provider-website;
mv localhost:4567/ terraform-provider-website;

docker stop "$(docker ps -q)") &

make website 
