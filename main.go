package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"encoding/json"
)

func main() {

	//http.HandleFunc("/github.png", index)
	http.HandleFunc("/AMT", amt)

	//fs := http.FileServer(http.Dir("/home/arpit/github-flair-server/static"))
	//http.Handle("/.well-known/", http.StripPrefix("/.well-known", fs))

	//http.ListenAndServe(":80",nil)
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

	collection :=session.DB("personal").C("amt")

	result := Entry{}

	var some = make([]interface{},len(ent))
	var thief =make([]interface{},len(ent))

	for i:=0;i<len(ent);i++ {

		some[i] = bson.M{"mac":ent[i].MAC}
		thief[i] = bson.M{"mac":ent[i].MAC,"hostname":ent[i].Hostname,"interfacename":ent[i].InterfaceName}
	}

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
		collection.Insert(thief...)

	}
	buffer := new(bytes.Buffer)
	buffer.WriteString(message)

	defer r.Body.Close()
	w.Write(buffer.Bytes())

}
