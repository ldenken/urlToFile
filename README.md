# urlToFile

A basic command line utility that fetches a url and saves the contents to a file.

To download and install this package run:

go get github.com/dyatlov/go-opengraph/opengraph


### Usage:

	$ urltofile -u http://www.bbc.co.uk/news/world
	
	urlToFile v0.2 https://github.com/ldenken

	    URL : http://www.bbc.co.uk/news/world
	created : download
	created : download/www.bbc.co.uk
	   Info : download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info
	   file : download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.html


## ./jq
[jq](http://stedolan.github.com/jq) is a lightweight and flexible command-line JSON processor.




	$ cat download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info | jq '.LinksInternal[]'

	...
	[
	  "http://www.bbc.co.uk/news/",
	  "News"
	]
	[
	  "http://www.bbc.co.uk/sport/",
	  "Sport"
	]
	...





$ cat /data/download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info | jq '.LinksExternal[][0]'



$ for url in $(cat url.lst); do echo $url; urltofile -d /data/download -u $url; done

$ find /data/download/www.bbc.co.uk/ -type f -name '*.info' -print0 |xargs --nul cat |jq '.LinksInternal[]' |sed 's/"//g' >tmp.txt 

|sort -u >> go_scrape/urls.txt


$ cat tmp.txt |sed 's/[ |,]//g' |grep -E '\-[0-9]{8}' |sort -u >url.lst


## License
License: [GNU Lesser General Public License Version 3, 29 June 2007](http://fsf.org/)
