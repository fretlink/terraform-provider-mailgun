#!/usr/bin/env bash

# sleep and the sub-command are there to make sure wget is called *after* `make website` has started. This ensures the website to be fetched correctly.
(sleep 10;
wget --recursive --no-parent --domains localhost --page-requisites --html-extension --convert-links --no-clobber http://localhost:4567/docs/providers/mailgun/;

rm -rf terraform-provider-website;
mv localhost:4567/ terraform-provider-website;
mv index.html terraform-provider-website/index.html

docker stop "$(docker ps -q)") &

make website 
