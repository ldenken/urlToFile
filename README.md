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
The *.info file contains file, http headers with internal and external links if the "Content-Type" = "text/html". JSON structure:

    File 			map[string]string 	`json:"File"` 
    Request 		map[string]string 	`json:"Request"` 
    Header 			map[string]string 	`json:"Header"`
    Response 		map[string]string 	`json:"Response"` 
    LinksInternal	[][]string 			`json:"LinksInternal"`
    LinksExternal	[][]string 			`json:"LinksExternal"`



$ cat /data/download/www.bbc.co.uk/11f2e26b746b0b07607feb09f10c1431.info | jq '.| {Header} '
{
  "Header": {
    "X-Pal-Host": "pal014.back.live.cwwtf.local:80",
    "X-News-Data-Centre": "cwwtf",
    "X-News-Cache-Id": "85512",
    "X-Lb-Nocache": "true",
    "X-Cache-Hits": "30",
    "X-Cache-Age": "19",
    "X-Cache-Action": "HIT",
    "Cache-Control": "private, max-age=30, stale-while-revalidate",
    "Connection": "keep-alive",
    "Content-Language": "en-GB",
    "Content-Type": "text/html; charset=utf-8",
    "Date": "Thu, 07 Jan 2016 16:50:13 GMT",
    "Server": "Apache",
    "Set-Cookie": "BBC-UID=b5c6983e5957a4556e6fea71f1c7a8beebf9e488a474a1deaa7174245e85e8e00Mozilla/5.0%20(X11%3b%20Ubuntu%3b%20Linux%20x86_64%3b%20rv:43.0)%20Gecko/20100101%20Firefox/43.0; expires=Mon, 06-Jan-20 16:50:13 GMT; path=/; domain=.bbc.co.uk",
    "Vary": "X-CDN,X-BBC-Edge-Cache,Accept-Encoding"
  }
}




### ./jq
[jq](http://stedolan.github.com/jq) is a lightweight and flexible command-line JSON processor.

We can use jq to extract all of the internal links.

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
