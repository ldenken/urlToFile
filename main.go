/*

TODO

1.
2.

$ clear; go install problemchild.local/gogs/urltofile; urltofile


*/

package main

import (

	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

)

const (
	// VERSION is the binary version.
	VERSION = "v0.0.2"
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
	USERAGENT = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0"
	//USERAGENT = "Golang Spider Bot " + VERSION

)

var (

	COLUMN int = 21

    url string = ""
    directoryBase string = ""
    directoryHost string = ""
	overwrite bool = false
	help bool = false
	version bool = false
	filenameInfo string = ""

)

func init() {
	// ----- parse flags -------------------------------------------------------

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



	// ----- getUrl ------------------------------------------------------------
	//fmt.Println(title_string("get url"))
	Response, Header, Request, body := getUrl(url)

	fmt.Println("\nResponse type:", reflect.TypeOf(Response).Kind(), "len:", len(Response))
    for key, value := range Response {
		printKeyValue(key, value, COLUMN)
    }
	fmt.Println("\nHeader type:", reflect.TypeOf(Header).Kind(), "len:", len(Header))
    for key, value := range Header {
		printKeyValue(key, value, COLUMN)
    }
	fmt.Println("\nRequest type:", reflect.TypeOf(Request).Kind(), "len:", len(Request))
    for key, value := range Header {
		printKeyValue(key, value, COLUMN)
    }
	fmt.Println("\nbody type:", reflect.TypeOf(body).Kind(), "len:", len(body))
	fmt.Println("\n")



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

func getUrl(url string) (map[string]string, map[string]string, map[string]string, []byte) {
	//fmt.Println("\ngetUrl ->", "url:", url)

	client := &http.Client {}
	fmt.Println("client type:", reflect.TypeOf(client).Kind(), "->\n", client)
	fmt.Println("client.Transport type:", reflect.TypeOf(client.transport).Kind())






	req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
	//fmt.Println("*req type:", reflect.TypeOf(*req).Kind(), "->\n", *req)

	req.Header.Add("User-Agent", USERAGENT)
	fmt.Println("\nreq type:", reflect.TypeOf(req).Kind(), "->\n", req)

	resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
	defer resp.Body.Close()


	// Response = resp *
	Response := make(map[string]string)
	//fmt.Println("\nresp type:", reflect.TypeOf(resp).Kind(), "->\n", resp)
	//fmt.Println("\n*resp type:", reflect.TypeOf(*resp).Kind(), "->\n", *resp)
	
	// fmt.Println("resp.Status type:", reflect.TypeOf(resp.Status).Kind(), resp.Status) 
	// resp.Status type: string 200 OK
	Response["status"] = resp.Status

	// fmt.Println("resp.StatusCode type:", reflect.TypeOf(resp.StatusCode).Kind(), resp.StatusCode) 
	// resp.StatusCode type: int 200
	Response["statuscode"] = strconv.Itoa(resp.StatusCode)

	// fmt.Println("resp.Proto type:", reflect.TypeOf(resp.Proto).Kind(), resp.Proto) 
	// resp.Proto type: string HTTP/1.1
	Response["proto"] = resp.Proto

	// fmt.Println("resp.ProtoMajor type:", reflect.TypeOf(resp.ProtoMajor).Kind(), resp.ProtoMajor) 
	// resp.ProtoMajor type: int 1
	Response["protomajor"] = strconv.Itoa(resp.ProtoMajor)

	// fmt.Println("resp.ProtoMinor type:", reflect.TypeOf(resp.ProtoMinor).Kind(), resp.ProtoMinor) 
	// resp.ProtoMinor type: int 1 
	Response["protominor"] = strconv.Itoa(resp.ProtoMinor)

	// fmt.Println("resp.ContentLength type:", reflect.TypeOf(resp.ContentLength).Kind(), resp.ContentLength) 
	// resp.ContentLength type: int64 12150
	Response["contentlength"] = strconv.FormatInt(resp.ContentLength, 10)

	// fmt.Println("resp.TransferEncoding type:", reflect.TypeOf(resp.TransferEncoding).Kind(), resp.TransferEncoding) 
	// resp.TransferEncoding type: slice []
	//Response["transferencoding"] = resp.TransferEncoding

	// fmt.Println("resp.Close type:", reflect.TypeOf(resp.Close).Kind(), resp.Close) 
	// resp.Close type: bool false
	Response["close"] = strconv.FormatBool(resp.Close)

	// fmt.Println("resp.Trailer type:", reflect.TypeOf(resp.Trailer).Kind(), resp.Trailer) 
	// resp.Trailer type: map map[]
	//Response["trailer"] = resp.Trailer

	// fmt.Println("resp.TLS type:", reflect.TypeOf(resp.TLS).Kind(), resp.TLS) 
	// resp.TLS type: ptr <nil>
	//Response["tls"] = resp.TLS

	// fmt.Println("\nresponse type:", reflect.TypeOf(Response).Kind(), "->\n", Response, "\n")
	//fmt.Println("Response")
	//for k, v := range Response {
	//	fmt.Println(" ", k, ":", v)
	//}



	// Header = resp.Header
	Header := make(map[string]string)
	// fmt.Println("\nresp.Header type:", reflect.TypeOf(resp.Header).Kind(), "->\n", resp.Header, "\n")
	for k, v := range resp.Header {
		//k := strings.ToLower(k)
		v := sliceToString(v)
		v = stripString(v)
		Header[k] = v
	}
	// fmt.Println("\nheader type:", reflect.TypeOf(Header).Kind(), "->\n", Header, "\n")
	//fmt.Println("Header")
	//for k, v := range Header {
	//	fmt.Println(" ", k, ":", v)
	//}



	// Request = resp.Request.*
	Request := make(map[string]string)
	//fmt.Println("\nresp.Request type:", reflect.TypeOf(*resp.Request).Kind(), "->\n", resp.Request, "\n")

	// fmt.Println("resp.Request.Method type:", reflect.TypeOf(resp.Request.Method).Kind(), resp.Request.Method)
	// resp.Request.Method type: string GET
	Request["method"] = resp.Request.Method

	// fmt.Println("*resp.Request.URL type:", reflect.TypeOf(*resp.Request.URL).Kind(), *resp.Request.URL)
	// *resp.Request.URL type: struct {http  <nil> httpbin.org /   }
	Request["url"] = resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + resp.Request.URL.Path

	// fmt.Println("resp.Request.Proto type:", reflect.TypeOf(resp.Request.Proto).Kind(), resp.Request.Proto)
	// resp.Request.Proto type: string HTTP/1.1
	Request["proto"] = resp.Request.Proto

	// fmt.Println("resp.Request.ProtoMajor type:", reflect.TypeOf(resp.Request.ProtoMajor).Kind(), resp.Request.ProtoMajor)
	// resp.Request.ProtoMajor type: int 1
	Request["protomajor"] = strconv.Itoa(resp.Request.ProtoMajor)

	// fmt.Println("resp.Request.ProtoMinor type:", reflect.TypeOf(resp.Request.ProtoMinor).Kind(), resp.Request.ProtoMinor)
	// resp.Request.ProtoMinor type: int 1
	Request["protominor"] = strconv.Itoa(resp.Request.ProtoMinor)

	// fmt.Println("resp.Request.Header type:", reflect.TypeOf(resp.Request.Header).Kind(), resp.Request.Header)
	for k, v := range resp.Request.Header {
		//k := strings.ToLower(k)
		v := sliceToString(v)
		v = stripString(v)
		Request[k] = v
	}

	// fmt.Println("resp.Request.Body type:", reflect.TypeOf(resp.Request.Body).Kind(), resp.Request.Body)

	// fmt.Println("resp.Request.ContentLength type:", reflect.TypeOf(resp.Request.ContentLength).Kind(), resp.Request.ContentLength)
	// resp.Request.ContentLength type: int64 0
	Request["contentlength"] = strconv.FormatInt(resp.Request.ContentLength, 10)

	// fmt.Println("resp.Request.TransferEncoding type:", reflect.TypeOf(resp.Request.TransferEncoding).Kind(), resp.Request.TransferEncoding)
	// resp.Request.TransferEncoding type: slice []

	// fmt.Println("resp.Request.Close type:", reflect.TypeOf(resp.Request.Close).Kind(), resp.Request.Close)
	// resp.Request.Close type: bool false
	Request["close"] = strconv.FormatBool(resp.Request.Close)

	// fmt.Println("resp.Request.Host type:", reflect.TypeOf(resp.Request.Host).Kind(), resp.Request.Host)
	// resp.Request.Host type: string httpbin.org
	Request["host"] = resp.Request.Host

	// fmt.Println("resp.Request.Form type:", 			reflect.TypeOf(resp.Request.Form).Kind(), 			resp.Request.Form)
	// resp.Request.Form type: map map[]
	// fmt.Println("resp.Request.PostForm type:", 		reflect.TypeOf(resp.Request.PostForm).Kind(), 		resp.Request.PostForm)
	// resp.Request.PostForm type: map map[]

	//fmt.Println("resp.Request.MultipartForm type:", reflect.TypeOf(resp.Request.MultipartForm).Kind(), resp.Request.MultipartForm)
	//fmt.Println("resp.Request.Trailer type:", 		reflect.TypeOf(resp.Request.Trailer).Kind(), 		resp.Request.Trailer)

	// fmt.Println("resp.Request.RemoteAddr type:", reflect.TypeOf(resp.Request.RemoteAddr).Kind(), resp.Request.RemoteAddr)
	// resp.Request.RemoteAddr type: string
	Request["remoteaddr"] = resp.Request.RemoteAddr

	// fmt.Println("resp.Request.RequestURI type:", reflect.TypeOf(resp.Request.RequestURI).Kind(), resp.Request.RequestURI)
	// resp.Request.RequestURI type: string
	Request["requesturi"] = resp.Request.RequestURI

	//fmt.Println("*resp.Request.TLS type:", reflect.TypeOf(*resp.Request.TLS).Kind(), *resp.Request.TLS)
	// fmt.Println("resp.Request.Cancel type:", 		reflect.TypeOf(resp.Request.Cancel).Kind(), 		resp.Request.Cancel)

	// fmt.Println("\nrequest type:", reflect.TypeOf(Request).Kind(), "->\n", Request, "\n")
	//fmt.Println("Request")
	//for k, v := range Request {
	//	fmt.Println(" ", k, ":", v)
	//}
	//fmt.Println("")


	// Body io.ReadCloser
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("\nbody type:", reflect.TypeOf(body).Kind(), "len:", len(body))

	return Response, Header, Request, body
}

func sliceToString(s []string) string {
	var buffer bytes.Buffer
    for _, value := range s {
		buffer.WriteString(value)
    }
	return buffer.String()
}

func stripString(s string) string {
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\t")
	s = strings.Trim(s, " ")
	return s
}


/*


*/