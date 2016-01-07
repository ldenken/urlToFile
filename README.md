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


### Information File:
The *.info file contains a JSON structure containing information about the downloaded file, http headers and internal/external links if the "Content-Type" = "text/html". 

    File 			map[string]string 	`json:"File"` 
    Request 		map[string]string 	`json:"Request"` 
    Header 			map[string]string 	`json:"Header"`
    Response 		map[string]string 	`json:"Response"` 
    LinksInternal	[][]string 			`json:"LinksInternal"`
    LinksExternal	[][]string 			`json:"LinksExternal"`


### [./jq](http://stedolan.github.com/jq)
jq is a lightweight and flexible command-line JSON processor and can be used to extract information from the *.info file(s).

	$ jq '.| {File}' download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info
	{
	  "File": {
	    "url": "http://www.bbc.co.uk/news/world",
	    "timestamp": "2016-01-07T16:54:00Z",
	    "filename": "download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.html"
	  }
	}

	$ jq '.| {Request}' download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info 
	{
	  "Request": {
	  	...
	    "User-Agent": "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0",
	    "URL": "http://www.bbc.co.uk/news/world",
	    "RequestURI": "",
	    ...
	  }
	}

	$ jq '.| {Header}' download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info{
	  "Header": {
	  	...
	    "Connection": "keep-alive",
	    "Content-Language": "en-GB",
	    "Content-Type": "text/html; charset=utf-8",
	    "Date": "Thu, 07 Jan 2016 16:54:00 GMT",
	    ...
	  }
	}

	$ jq '.| {Response}' download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info
	{
	  "Response": {
	  	...
	    "StatusCode": "200",
	    "Status": "200 OK",
	    ...
	  }
	}

	$ jq '.LinksInternal[]' download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info
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
