package main

import (

	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

)

const (
	// VERSION is the binary version.
	VERSION = "v0.2"
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

Usage: urlToFile [-u|-url] {url} [-d] {directory} [-o] [-v]
URL:
  -u|-url      url to download
DIRECTORY:
  -d           root path for the download directory
OVERWRITE:
  -o           overwrite existing downloaded file
HELP:
  -h|-help     print help information and exit
verbose:
  -v           print verbose output
EXAMPLES:
  $ urltofile -u http://www.bbc.co.uk/news/world
  $ urltofile -o -d /code/download -u http://www.bbc.co.uk/news/world

`
	USERAGENT = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0"
	//USERAGENT = "Golang Bot " + VERSION

)

var (

	column int = 8

	help bool = false
	overwrite bool = false
	verbose bool = false
    directory string = ""

    url string = ""
    text string = ""
	filename string = ""

)

func init() {
	// ----- parse flags -------------------------------------------------------

    flag.StringVar(&url, "url", "", "a url to download")
    flag.StringVar(&url, "u", "", "a url to download (shorthand)")

    flag.StringVar(&directory, "d", "download", "root path for the download directory")
	flag.BoolVar(&overwrite, "o", false, "overwrite any existing downloaded file")

	flag.BoolVar(&help, "help", false, "print help and exit")
	flag.BoolVar(&help, "h", false, "print help and exit (shorthand)")

	flag.BoolVar(&verbose, "v", false, "print verbose output (shorthand)")

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

	fmt.Println("\nurlToFile", VERSION, "https://github.com/ldenken\n")

	if url != "" {
		reg, err := regexp.Compile("^(ftp|http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&@!-/]))?")
	    if err != nil {
	        log.Fatal(err)
	    }
		if reg.MatchString(url) == true {
			printKeyValue("URL", url, column)
		} else {
	        log.Fatal("url failed regexp! ", url)
		}
	}

	if directory != "" {
		directory = strings.TrimRight(directory, "/")
		if existsTF(directory) == false {
			createDirectory(directory)
			printKeyValue("created", directory, column)
		}
		//printKeyValue("Base", directory, column)
	}

	//if overwrite {
	//	printKeyValue("Write", "true", column)
	//}



	// ----- parse valid url into its component parts --------------------------
	url_slice := []string(strings.Split(url, "://"))

	protocol := url_slice[0]
	//printKeyValue("protocol", protocol, column)

	urlMD5 := getMD5(url_slice[1])
	//printKeyValue("urlMD5", urlMD5, column)

	host := []string(strings.Split(url_slice[1], "/"))[0]
	//printKeyValue("host", host, column)

	//directoryHost := directory + "/" + host
	if existsTF(directory + "/" + host) == false {
		createDirectory(directory + "/" + host)
		printKeyValue("created", directory + "/" + host, column)
	}
	//printKeyValue("Host", directoryHost, column)



	// ----- check if the information file exists ------------------------------
	filename = directory + "/" + host + "/" + urlMD5 + ".info"
    test, err := existsTFE(filename)
    if err != nil {
        log.Fatal(err)
    }
    if test == true && overwrite == false {
		printKeyValue("Exists", filename, column)
		fmt.Println("")
        os.Exit(1)
    }
	printKeyValue("Info", filename, column)



	// ----- getUrl ------------------------------------------------------------
	Request, Header, Response, body := getUrl(url)

	if verbose {
		column = 21

		fmt.Println("\nRequest type:", reflect.TypeOf(Request).Kind(), "len:", len(Request))
	    for key, value := range Request {
			printKeyValue(key, value, column)
	    }

		fmt.Println("\nHeader type:", reflect.TypeOf(Header).Kind(), "len:", len(Header))
	    for key, value := range Header {
			printKeyValue(key, value, column)
	    }

		fmt.Println("\nResponse type:", reflect.TypeOf(Response).Kind(), "len:", len(Response))
	    for key, value := range Response {
			printKeyValue(key, value, column)
	    }
		fmt.Println("\nbody type:", reflect.TypeOf(body).Kind(), "len:", len(body))

		fmt.Println("")
		column = 14
	}



	// ----- write body to file ------------------------------------------------
	x := strings.Split(Header["Content-Type"], ";")
	y := strings.SplitAfterN(x[0], "/", -1)
	var contentType string = y[len(y)-1]
	//printKeyValue("contentType", contentType, column)

	//var filename string = directoryHost + "/" + urlMD5 + "." + contentType
	filename = directory + "/" + host + "/" + urlMD5 + "." + contentType

	printKeyValue("file", filename, column)

	var fileType string = ""
	if contentType == "html" {
		fileType = "string"
	}
	if contentType == "pdf" {
		fileType = "byte"
	}

    switch fileType {
    	case "string":
			wirteFile(filename, []byte(string(body)))
    	case "byte":
			wirteFile(filename, body)
		default:
			fmt.Println("unknown fileType ->", fileType)
    }



	// ----- extract links if contentType == "html" -----------------------
	LinksInternal := [][]string{}
	LinksExternal := [][]string{}

	if len(body) > 0 && contentType == "html" {

		text = string(body)
		text = html.UnescapeString(text)

		text = regexp.MustCompile("\r").ReplaceAllString(text, " ")
		text = regexp.MustCompile("\n").ReplaceAllString(text, " ")
		text = regexp.MustCompile("\t").ReplaceAllString(text, " ")
		text = regexp.MustCompile(" {2,}").ReplaceAllString(text, " ")

		text = regexp.MustCompile("(?i)<div").ReplaceAllString(text, "\n<div")

		// remove <span.*, .*>
		text = regexp.MustCompile("(?i)<span").ReplaceAllString(text, "\n<span")
		text = regexp.MustCompile("(?i)>").ReplaceAllString(text, ">\n")
		text = regexp.MustCompile("(?i)<span.*").ReplaceAllString(text, "")
		text = regexp.MustCompile("\n").ReplaceAllString(text, "")

		// remove <h[1-6].*, .*>
		text = regexp.MustCompile("(?i)<h[1-6]").ReplaceAllString(text, "\n<h1")
		text = regexp.MustCompile("(?i)>").ReplaceAllString(text, ">\n")
		text = regexp.MustCompile("(?i)<h1.*").ReplaceAllString(text, "")
		text = regexp.MustCompile("\n").ReplaceAllString(text, "")

		// remove image links from <a href regexp Response
		text = regexp.MustCompile("(?i)<img").ReplaceAllString(text, "\n<img")
		text = regexp.MustCompile("(?i)<svg").ReplaceAllString(text, "\n<svg")

		// remove protocol links
		text = regexp.MustCompile("(?i)<a href=\"mailto:").ReplaceAllString(text, "\n*****")
		text = regexp.MustCompile("(?i)<a href=\"whatsapp:").ReplaceAllString(text, "\n*****")

		text = regexp.MustCompile(" {2,}").ReplaceAllString(text, " ")

		text = regexp.MustCompile("(?i)<a href=").ReplaceAllString(text, "\n<a href=")
		text = regexp.MustCompile("(?i)</a>").ReplaceAllString(text, "</a>\n")

		//fmt.Println("\ntext type:", reflect.TypeOf(text).Kind(), "len:", len(text), "\n", text)

		href := regexp.MustCompile("(?i)<a href=.*</a>").FindAllString(text, -1)
		//fmt.Println("\nhref type:", reflect.TypeOf(href).Kind(), "len:", len(href), "\n", href)
		// href type: slice len: 62

		for _, v := range href {
			v = strings.Trim(v, " \t\n\r")
			tmpslice := strings.SplitAfterN(v, ">", 2)
			tmpslice[0] = regexp.MustCompile("(?i)<a href=(\"|')").ReplaceAllString(tmpslice[0], "")
			tmpslice[0] = regexp.MustCompile("(\"|').*").ReplaceAllString(tmpslice[0], "")
			tmpslice[0] = strings.Trim(tmpslice[0], " \t\n\r")
			tmpslice[1] = regexp.MustCompile("(?i)</a>").ReplaceAllString(tmpslice[1], "")
			tmpslice[1] = regexp.MustCompile("(?i)<[/a-z0-9]{1,}>").ReplaceAllString(tmpslice[1], "")
			tmpslice[1] = strings.Trim(tmpslice[1], " \t\n\r")
			if regexp.MustCompile("^/").MatchString(tmpslice[0]) == true {
				tmpslice[0] = protocol + "://" + host + tmpslice[0]
			}
			if regexp.MustCompile(host).MatchString(tmpslice[0]) == true {
				link := []string{}
				link = append(link, tmpslice[0])
				link = append(link, tmpslice[1])
				LinksInternal = append(LinksInternal, link)
			} else {
				if regexp.MustCompile("^http").MatchString(tmpslice[0]) == true {
					link := []string{}
					link = append(link, tmpslice[0])
					link = append(link, tmpslice[1])
					LinksExternal = append(LinksExternal, link)
				}
			}
		}
		if verbose {
			fmt.Println("\nLinksInternal type:", reflect.TypeOf(LinksInternal).Kind(), "len:", len(LinksInternal))
			for i, v := range LinksInternal {
				fmt.Println(i, v[0], v[1])
			}

			fmt.Println("\nLinksExternal type:", reflect.TypeOf(LinksExternal).Kind(), "len:", len(LinksExternal))
			for i, v := range LinksExternal {
				fmt.Println(i, v[0], v[1])
			}
		}
	}



	// ----- build File information map ------------------------------------
	File := make(map[string]string)
	File["filename"] = filename
	File["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	File["url"] = url


	// ----- build json Information struct ---------------------------------
	type Information struct {
	    File 			map[string]string 	`json:"File"`
	    Request 		map[string]string 	`json:"Request"`
	    Header 			map[string]string 	`json:"Header"`
	    Response 		map[string]string 	`json:"Response"`
	    LinksInternal	[][]string 			`json:"LinksInternal"`
	    LinksExternal	[][]string 			`json:"LinksExternal"`
	}

	i := Information{}

	i.File 				= File
	i.Request 			= Request
	i.Response 			= Response
	i.Header 			= Header
	i.LinksInternal 	= LinksInternal
	i.LinksExternal 	= LinksExternal

	infoJson, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("\ninfoJson type:", reflect.TypeOf(infoJson).Kind(), "len:", len(infoJson))


	// ----- write infoJson to file ---------------------------------------
	filename = directory + "/" + host + "/" + urlMD5 + ".info"
	wirteFile(filename, infoJson)


	fmt.Println("")
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
	//fmt.Println("client type:", reflect.TypeOf(client).Kind(), "->\n", client)
	//fmt.Println("client.Transport type:", reflect.TypeOf(client.transport).Kind())

	req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
	//fmt.Println("*req type:", reflect.TypeOf(*req).Kind(), "->\n", *req)

	req.Header.Add("User-Agent", USERAGENT)
	//fmt.Println("\nreq type:", reflect.TypeOf(req).Kind(), "->\n", req)

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
	Response["Status"] = resp.Status

	// fmt.Println("resp.StatusCode type:", reflect.TypeOf(resp.StatusCode).Kind(), resp.StatusCode)
	// resp.StatusCode type: int 200
	Response["StatusCode"] = strconv.Itoa(resp.StatusCode)

	// fmt.Println("resp.Proto type:", reflect.TypeOf(resp.Proto).Kind(), resp.Proto)
	// resp.Proto type: string HTTP/1.1
	Response["Proto"] = resp.Proto

	// fmt.Println("resp.ProtoMajor type:", reflect.TypeOf(resp.ProtoMajor).Kind(), resp.ProtoMajor)
	// resp.ProtoMajor type: int 1
	Response["ProtoMajor"] = strconv.Itoa(resp.ProtoMajor)

	// fmt.Println("resp.ProtoMinor type:", reflect.TypeOf(resp.ProtoMinor).Kind(), resp.ProtoMinor)
	// resp.ProtoMinor type: int 1
	Response["ProtoMinor"] = strconv.Itoa(resp.ProtoMinor)

	//fmt.Println("resp.ContentLength type:", reflect.TypeOf(resp.ContentLength).Kind(), resp.ContentLength)
	// resp.ContentLength type: int64 12150
	Response["ContentLength"] = strconv.FormatInt(resp.ContentLength, 10)

	// fmt.Println("resp.TransferEncoding type:", reflect.TypeOf(resp.TransferEncoding).Kind(), resp.TransferEncoding)
	// resp.TransferEncoding type: slice []
	//Response["transferencoding"] = resp.TransferEncoding

	// fmt.Println("resp.Close type:", reflect.TypeOf(resp.Close).Kind(), resp.Close)
	// resp.Close type: bool false
	//Response["close"] = strconv.FormatBool(resp.Close)

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
	Request["Method"] = resp.Request.Method

	// fmt.Println("*resp.Request.URL type:", reflect.TypeOf(*resp.Request.URL).Kind(), *resp.Request.URL)
	// *resp.Request.URL type: struct {http  <nil> httpbin.org /   }
	Request["URL"] = resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + resp.Request.URL.Path

	// fmt.Println("resp.Request.Proto type:", reflect.TypeOf(resp.Request.Proto).Kind(), resp.Request.Proto)
	// resp.Request.Proto type: string HTTP/1.1
	Request["Proto"] = resp.Request.Proto

	// fmt.Println("resp.Request.ProtoMajor type:", reflect.TypeOf(resp.Request.ProtoMajor).Kind(), resp.Request.ProtoMajor)
	// resp.Request.ProtoMajor type: int 1
	Request["ProtoMajor"] = strconv.Itoa(resp.Request.ProtoMajor)

	// fmt.Println("resp.Request.ProtoMinor type:", reflect.TypeOf(resp.Request.ProtoMinor).Kind(), resp.Request.ProtoMinor)
	// resp.Request.ProtoMinor type: int 1
	Request["ProtoMinor"] = strconv.Itoa(resp.Request.ProtoMinor)

	// fmt.Println("resp.Request.Header type:", reflect.TypeOf(resp.Request.Header).Kind(), resp.Request.Header)
	for k, v := range resp.Request.Header {
		v := sliceToString(v)
		v = stripString(v)
		Request[k] = v
	}

	// fmt.Println("resp.Request.Body type:", reflect.TypeOf(resp.Request.Body).Kind(), resp.Request.Body)

	// fmt.Println("resp.Request.ContentLength type:", reflect.TypeOf(resp.Request.ContentLength).Kind(), resp.Request.ContentLength)
	// resp.Request.ContentLength type: int64 0
	Request["ContentLength"] = strconv.FormatInt(resp.Request.ContentLength, 10)

	// fmt.Println("resp.Request.TransferEncoding type:", reflect.TypeOf(resp.Request.TransferEncoding).Kind(), resp.Request.TransferEncoding)
	// resp.Request.TransferEncoding type: slice []

	// fmt.Println("resp.Request.Close type:", reflect.TypeOf(resp.Request.Close).Kind(), resp.Request.Close)
	// resp.Request.Close type: bool false
	Request["Close"] = strconv.FormatBool(resp.Request.Close)

	// fmt.Println("resp.Request.Host type:", reflect.TypeOf(resp.Request.Host).Kind(), resp.Request.Host)
	// resp.Request.Host type: string httpbin.org
	Request["Host"] = resp.Request.Host

	// fmt.Println("resp.Request.Form type:", 			reflect.TypeOf(resp.Request.Form).Kind(), 			resp.Request.Form)
	// resp.Request.Form type: map map[]
	// fmt.Println("resp.Request.PostForm type:", 		reflect.TypeOf(resp.Request.PostForm).Kind(), 		resp.Request.PostForm)
	// resp.Request.PostForm type: map map[]

	//fmt.Println("resp.Request.MultipartForm type:", reflect.TypeOf(resp.Request.MultipartForm).Kind(), resp.Request.MultipartForm)
	//fmt.Println("resp.Request.Trailer type:", 		reflect.TypeOf(resp.Request.Trailer).Kind(), 		resp.Request.Trailer)

	// fmt.Println("resp.Request.RemoteAddr type:", reflect.TypeOf(resp.Request.RemoteAddr).Kind(), resp.Request.RemoteAddr)
	// resp.Request.RemoteAddr type: string
	Request["RemoteAddr"] = resp.Request.RemoteAddr

	// fmt.Println("resp.Request.RequestURI type:", reflect.TypeOf(resp.Request.RequestURI).Kind(), resp.Request.RequestURI)
	// resp.Request.RequestURI type: string
	Request["RequestURI"] = resp.Request.RequestURI

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

	return Request, Header, Response, body
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


func wirteFile(filename string, content []byte) {
	//fmt.Println("filename:", filename)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		log.Fatal(err)
	}
	f.Sync()
}


/*


*/