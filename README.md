# urlToFile

A basic command line utility that fetches a url and saves the contents to a file.

To download and install this package run:

go get github.com/dyatlov/go-opengraph/opengraph


### Usage:

	$ urltofile -d /data/download -u http://www.bbc.co.uk/news/world

	urlToFile v0.2 https://github.com/ldenken

	    URL : http://www.bbc.co.uk/news/world
	   Info : /data/download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info
	   file : /data/download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.html



http://stedolan.github.com/jq

$ cat /data/download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info | jq '.LinksInternal[][0]'
$ cat /data/download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info | jq '.LinksExternal[][0]'



$ for url in $(cat url.lst); do echo $url; urltofile -d /data/download -u $url; done

$ find /data/download/www.bbc.co.uk/ -type f -name '*.info' -print0 |xargs --nul cat |jq '.LinksInternal[]' |sed 's/"//g' >tmp.txt 

|sort -u >> go_scrape/urls.txt


$ cat tmp.txt |sed 's/[ |,]//g' |grep -E '\-[0-9]{8}' |sort -u >url.lst


## License
License: [GNU Lesser General Public License Version 3, 29 June 2007](http://fsf.org/)
