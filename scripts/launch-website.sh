#!/usr/bin/env bash

#the sleep and the subquery & are here to ensure that wget is called after make website have done its work so that the website is gotten correctly
(sleep 10;
wget --recursive --no-parent --domains localhost --page-requisites --html-extension --convert-links --no-clobber http://localhost:4567/docs/providers/mailgun/;

rm -rf terraform-provider-website;
mv localhost:4567/ terraform-provider-website;
mv index.html terraform-provider-website/index.html

docker stop "$(docker ps -q)") &

make website 
