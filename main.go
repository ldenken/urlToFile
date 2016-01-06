/*

TODO

1.	
2. 	

$ clear; go install problemchild.local/gogs/urltofile; urltofile


*/

package main

import (

	"crypto/md5"
	"flag"
	"fmt"
	"encoding/hex"
	"reflect"
	"log"
	"os"
	"regexp"
	"strings"
	"bytes"

)

const (
	// VERSION is the binary version.
	VERSION = "v0.0.1"
	// BANNER is what is printed for help/info output.
	BANNER = `            _ _____    ______ _ _      
           | |_   _|   |  ___(_) |     
 _   _ _ __| | | | ___ | |_   _| | ___ 
| | | | '__| | | |/ _ \|  _| | | |/ _ \
| |_| | |  | | | | (_) | |   | | |  __/
 \__,_|_|  |_| \_/\___/\_|   |_|_|\___|
                                       
Download a URL and save the contents to a file
urlToFile ` + VERSION + ` https://github.com/ldenken`

	HELP = BANNER + `

Usage: urlToFile [-u|-url] {url} [-d] {directory} [-o] {overwrite}
URL:
  -u | -url      url to download
DIRECTORY:
  -d             root path for the download directory 
OVERWRITE:
  -o             overwrite any existing downloaded file
HELP:
  -h | -help     print help information and exit
VERSION:
  -v | -version  print version information and exit
EXAMPLES:

`

)

var (

	COLUMN int = 14

    url string = ""
    directoryBase string = ""
    directoryHost string = ""
	overwrite bool = false
	help bool = false
	version bool = false
	filenameInfo string = ""

)

func init() {
	// parse flags

	//url_string = flag.BoolVar(&url, "url", false, "url to download")
	//flag.BoolVar(&url, "u", false, "url to download (shorthand)")
	//u = flag.String("u", "", "a url string")


    flag.StringVar(&url, "url", "", "a url to download")
    flag.StringVar(&url, "u", "", "a url to download (shorthand)")

    flag.StringVar(&directoryBase, "d", "", "root path for the download directory")
	flag.BoolVar(&overwrite, "o", false, "overwrite any existing downloaded file")

	flag.BoolVar(&help, "help", false, "print help and exit")
	flag.BoolVar(&help, "h", false, "print help and exit (shorthand)")

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")

	flag.Parse()
}

func main() {
	// ----- parse the args ----------------------------------------------------

	if len(os.Args) <= 1 {
		fmt.Println(HELP)
		return
	}

	if help {
		fmt.Println(HELP)
		return
	}

	if version {
		fmt.Println(VERSION)
		return
	}

	if url != "" {
		reg, err := regexp.Compile("^(ftp|http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&@!-/]))?")
	    if err != nil {
	        log.Fatal(err)
	    }
		if reg.MatchString(url) == true {
			printKeyValue("url", url, COLUMN)
		} else {
	        log.Fatal("url failed regexp! ", url)
		}
	}

	if directoryBase != "" {
		directoryBase = strings.TrimRight(directoryBase, "/")
		if existsTF(directoryBase) == false {
			createDirectory(directoryBase)
			printKeyValue("created", directoryBase, COLUMN)
		}
		printKeyValue("directoryBase", directoryBase, COLUMN)
	}

	if overwrite {
		printKeyValue("overwrite", "true", COLUMN)
	}



	// ----- parse valid url into its component parts --------------------------

	url_slice := []string(strings.Split(url, "://"))

	protocol := url_slice[0]
	printKeyValue("protocol", protocol, COLUMN)

	urlMD5 := getMD5(url_slice[1])
	printKeyValue("urlMD5", urlMD5, COLUMN)

	host := []string(strings.Split(url_slice[1], "/"))[0]
	printKeyValue("host", host, COLUMN)

	directoryHost := directoryBase + "/" + host
	if existsTF(directoryHost) == false {
		createDirectory(directoryHost)
		fmt.Println("created:", directoryHost)
	}
	printKeyValue("directoryHost", directoryHost, COLUMN)



	// ----- check if the information file exists ------------------------------
	filenameInfo = directoryBase + "/" + host + "/" + urlMD5 + ".info"
    test, err := existsTFE(filenameInfo)
    if err != nil {
        log.Fatal(err)
    }
    if test == true && overwrite == false {
	    fmt.Println(filenameInfo, "Exists!\n")
        os.Exit(1)
    }
	printKeyValue("filenameInfo", filenameInfo, COLUMN)


	fmt.Println("\n\nurl:", reflect.TypeOf(url).Kind(), "len:", len(url), "->\n", url)
}



func printKeyValue(key string, value string, column int) {
	//fmt.Println("printKeyValue", key:", key, "value:", value, "column:", column)
    var buffer bytes.Buffer
	if len(key) <= column {
		for i := 1; i <= (column - (len(key)+1)); i++ {
	        buffer.WriteString(" ")
		}
	}
	var txt string = buffer.String() + key
	fmt.Println(txt, ":", value)	
}

func existsTFE(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { 
    	return true, nil 
    }
    if os.IsNotExist(err) { 
    	return false, nil 
    }
    return true, err
}

func existsTF(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { 
    	return true
    }
    if os.IsNotExist(err) { 
    	return false
    }
    return true	
}

func createDirectory(directory string) {
	err := os.MkdirAll(directory, 0711)
	if err != nil {
	  log.Fatal(err)
	}
}

func getMD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

/*


*/