package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"time"
	"github.com/jcatala/gqm/utility"
)


func notify2Telegram(builder strings.Builder){

	m := utility.ParseConfig(false)
	bot := utility.GenBot(m["apikey"])
	chatId, err := utility.GetNewChatId(bot,m["savedChatId"] , false)
	if err != nil{
		log.Fatalln(err)
		return
	}
	chatIdInt, err := strconv.ParseInt(chatId, 10 , 64)
	if err != nil{
		log.Fatalln(err)
	}
	// Here we craft a new msg
	messageBytes := []byte(builder.String())
	utility.SendMsgPredefined(bot, chatIdInt, true, false, messageBytes)
	return


}

func generateResponse(v bool, rBody string, redirect string) (string, error){

	// If the rBody is empty, the default is the unix time
	if rBody == "" {
		rBody = strconv.Itoa(int(time.Now().Unix()))
	}
	if redirect == "" {
		resp := `HTTP/1.1 200 OK
Content-Type: text/plain
Server: Go blind Receptor

`
		resp = resp + rBody
		return resp, nil
	}
	resp := `HTTP/1.1 302
Content-Type: text/plain
Server: Go Blind Receptor
Location: `
	resp = resp + redirect
	resp = resp + "\n\n"
	resp = resp + rBody

	return resp, nil
}

func handleConnection(c net.Conn, v bool, rBody string, redirect string, notify bool, prefix string){
	//buf := make([]byte, 256)
	//netData, err := bufio.NewReader(c).Read(buf)
	//netData, err := bufio.NewReader(c).Read
	//reader := bufio.NewReader(c)

	if v {
		fmt.Printf("Serving: %s\n", c.RemoteAddr().String())
		fmt.Printf("Using prefix: %s\n", prefix)
	}
	isValid := false
	requestText := strings.Builder{}
	go func(c net.Conn, isValid *bool, prefix string, request *strings.Builder) {
		reader := bufio.NewReader(c)
		for {
			line,err := reader.ReadString('\n')
			if strings.Contains(line, prefix){
				*isValid = true
			}
			if err != nil {
				break
			}
			request.WriteString(line)
			//request.WriteString("\n")
		}
	}(c, &isValid, prefix, &requestText)

	// If the prefix is "", means that the flag is empty, and every request is valid
	if prefix == ""{
		isValid = true
	}

	time.Sleep(2 * time.Second)

	// If the request is not valid, write out invalid, and do not proceed
	if !isValid{
		invalid,_ := generateResponse(v, "Invalid!", "")
		c.Write([]byte(string(invalid)))
		c.Close()
		return
	}

	// Generate the response
	resp, err := generateResponse(v, rBody, redirect)
	if err != nil{
		log.Println(err)
	}
	if v{
		fmt.Printf("Received:\n%s\n", requestText.String())
		fmt.Printf("Sending: %s\n", resp)
	}
	c.Write([]byte(string(resp)))
	c.Close()
	if notify {
		notify2Telegram(requestText)
	}

}

func listenServer(port int, v bool, rBody string, redirect string, notify bool, prefix string){
	PORT := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp4", PORT)
	if err != nil{
		log.Fatalln(err)
	}
	// We start to listen for new connections
	for {
		c, err := l.Accept()
		if err != nil{
			fmt.Println(err)
			continue
		}
		go handleConnection(c, v, rBody, redirect, notify, prefix)
	}


}

func main() {

	// They're just pointers
	verbose := flag.Bool("verbose",false, "To be verbose")
	redirect := flag.String("redirect","", "To make the server redirect somewhere")
	rBody := flag.String("rbody","", "Custom response body, default: UNIX TIMESTAMP")
	port := flag.Int("port", 9080, "Specify another port, default: 9080")
	//debugInfo := flag.Bool("debugInfo", false, "To get debug information")
	notify := flag.Bool("notify", false, "Notify the incoming request via telegram bot (Not recommended to listen to root directory)")
	prefix := flag.String("prefix", "", "Receive and notify just the requests with some certain prefix.")
	flag.Parse()

	listenServer(*port, *verbose, *rBody, *redirect, *notify, *prefix)





}