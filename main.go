package main

import (
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"encoding/json"
	"log"
	"net/smtp"
	"github.com/arpitjindal97/github-flair-server"
)

func main() {

	github_flair.SetTorProxy()
	github_flair.DownloadImages()

	http.HandleFunc("/AMT", amt)
	http.HandleFunc("/github.png",github_flair.Flair)

	http.ListenAndServeTLS(":8080", "certificate.pem",
		"private.key", nil)

}

type Entry struct {
	Hostname      string
	InterfaceName string
	MAC           string
}

func amt(w http.ResponseWriter, r *http.Request) {

	var ent []Entry

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&ent)

	session,_ := mgo.Dial("127.0.0.1")
	defer session.Close()

	collection :=session.DB("personal").C("amt_thief")

	result := Entry{}

	var some = make([]interface{},len(ent))
	var thief =make([]interface{},len(ent))

	for i:=0;i<len(ent);i++ {

		some[i] = bson.M{"mac":ent[i].MAC}
		thief[i] = bson.M{"mac":ent[i].MAC,"hostname":ent[i].Hostname,"interfacename":ent[i].InterfaceName}
		collection.RemoveAll(thief[i])
	}

	collection =session.DB("personal").C("amt")
	collection.Find(bson.M{"$or": some}).One(&result)

	message := "You are not authorized person.\nTo gain access to this software, Please contact"+
		" Arpit at this number +918285283150"
	if result.MAC != "" {
		message = "OK"
		for i:=0;i<len(ent);i++{
			if ent[i].MAC == result.MAC{
				collection.Update(
					bson.M{"mac":result.MAC},
					bson.M{"hostname":ent[i].Hostname,
						"mac":ent[i].MAC,
						"interfacename":ent[i].InterfaceName})
				break;
			}
		}
	} else {
		collection =session.DB("personal").C("amt_thief")
		//log.Println("Thief caught!!!")
		send_mail(ent)
		collection.Insert(thief...)
	}
	buffer := new(bytes.Buffer)
	buffer.WriteString(message)

	defer r.Body.Close()
	w.Write(buffer.Bytes())

}

var mail_user,mail_pass,mail_smtp_host,mail_smtp_port string
func send_mail(ent []Entry) {

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		mail_user,
		mail_pass,
		mail_smtp_host,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	email_string:="From: arpitjindal97@gmail.com\nSubject: Aircel-AMT Software\nTo: arpitjindal97@gmail.com"+
		"\nYour mechanism just caught a thief. Here are the details.\n\n"

		for i:=0;i<len(ent);i++{
			email_string += "Hostname : "+ent[i].Hostname+
				"\tInterfaceName : "+ent[i].InterfaceName+"\tMAC : "+ent[i].MAC+"\n";
		}

	err := smtp.SendMail(
		mail_smtp_host+":"+mail_smtp_port,
		auth,
		mail_user,
		[]string{"arpitjindal97@gmail.com"},
		[]byte(email_string),
	)
	if err != nil {
		log.Fatal(err)
	}
}
