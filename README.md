# urlToFile

A basic command line utility that fetches a url and saves the contents to a file.



		text = regexp.MustCompile(" {2,}").ReplaceAllString(text, " ")


$ go install problemchild.local/gogs/urltofile; urltofile -d /data/download/ -o -v -url http://www.bbc.co.uk/news/world



$ cat /data/download/www.eff.org/7f0f99e5c048890d2f6bd22adb63d155.info | jq '.LinksInternal[][0]'
$ cat /data/download/www.eff.org/7f0f99e5c048890d2f6bd22adb63d155.info | jq '.LinksExternal[][0]'



$ clear; go install problemchild.local/gogs/urltofile; 

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


