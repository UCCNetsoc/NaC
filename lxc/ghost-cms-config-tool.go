package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

const port = 80

const form = "<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><meta name='viewport' content='width=device-width, initial-scale=1.0'><link rel='icon' type='image/png' href='https://raw.githubusercontent.com/UCCNetsoc/cloud/master/ui/public/favicon.ico'><title>Configuration</title></head><style>body{padding:0;margin:0;background-color:#212121;font-size:14px;font-weight:400;font-family:'Roboto',sans-serif !important;color:white}#header{background-color:#2196F3;text-align:center}#header img{height:48px;margin:0 auto;padding:16px}#content{text-align:center;max-width:700px;margin:0 auto}#content-body{padding-top:16px;border-top:1px solid #333;color:gray;border-bottom:1px solid #333}h1{font-size:24px;text-transform:lowercase;font-weight:200;padding:8px 0}label{position:relative;opacity:0.5;margin:0px auto;font-size:1.2em;padding:8px}input{margin:8px auto;height:16px;background-color:#00000040;border:none;padding:8px;color:white;width:calc(100% - 16px);font-size:16px}form{max-width:500px;width:100%;margin:0 auto;text-align:left}a{color:#fff}#submit{background-color:#2196F3;font-weight:500;margin-top:20px;max-width:400px;padding:20px 0;text-transform:uppercase;font-family:'Roboto',sans-serif !important;border-radius:5px;line-height:0;letter-spacing:0.09em}</style><body><div id='header'> <img src='https://raw.githubusercontent.com/UCCNetsoc/wiki/master/assets/logo-horizontal.svg'></div><div id='content'><div><h1>Ghost CMS Configuration</h1></div><div id='content-body'><form id='form' action='/config'> <label for='host'>host</label> <input type='text' id='host' name='host'></form><div id='info'><p>Ensure that the entered host is correctly configured as a vHost</p><p>Contact Netsoc SysAdmins on our <a href='discord.netsoc.co'>Discord</a> or by <a href='mailto:netsoc@uccsocieties.ie'>Email</a> if you are unsure or require assistance!</p></div></div> <input type='submit' value='Set host' id='submit' form='form'></div></div><script>document.getElementById('host').placeholder=window.location.host;document.getElementById('host').value = window.location.host;</script></body></html>"
const response = "<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><meta name='viewport' content='width=device-width, initial-scale=1.0'><link rel='icon' type='image/png' href='https://raw.githubusercontent.com/UCCNetsoc/cloud/master/ui/public/favicon.ico'><title>Configuration</title></head><style>body{padding:0;margin:0;background-color:#212121;font-size:14px;font-weight:400;font-family:'Roboto',sans-serif !important;color:white}#header{background-color:#2196F3;text-align:center}#header img{height:48px;margin:0 auto;padding:16px}#content{text-align:center;max-width:700px;margin:0 auto}#content-body{padding-top:16px;border-top:1px solid #333;border-bottom:1px solid #333;color:gray}h1{font-size:24px;text-transform:lowercase;font-weight:200;padding:8px 0}a{color:#fff}</style><body><div id='header'> <img src='https://raw.githubusercontent.com/UCCNetsoc/wiki/master/assets/logo-horizontal.svg'></div><div id='content'><div><h1>Ghost CMS Configuration</h1></div><div id='content-body'> {{if .Success}}<h2>Configured Ghost CMS to use {{.Host}}</h2><p>Find the admin page here <a href='https://{{.Host}}/ghost' target='_blank'>{{.Host}}</a></p><p>It may take a minute to start. If you have any issues don't hesitate to contact a Netsoc SysAdmin :D</p> {{else}}<h2>Internal Server Error</h2><p>Unable to complete the request</p><p>Manual configuration can be completed by following this guide (TODO). Contact a SysAdmin if assistance is required :D</p> {{end}}<p>Contact Netsoc SysAdmins on our <a href='discord.netsoc.co'>Discord</a> or by <a href='mailto:sysadmins@netsoc.co'>Email</a></p></div></div></body></html>"

func main() {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := startHTTPServer(httpServerExitDone)

	done := make(chan bool)
	configureFunc := configure(done)

	http.HandleFunc("/", showForm)
	http.HandleFunc("/config", configureFunc)

	<-done

	// Wait for server to shutdown and wait additional 10s for good measure
	_ = srv.Shutdown(context.TODO())
	httpServerExitDone.Wait()
	time.Sleep(10 * time.Second)

	// Commands to execute after web server shutdown
	exec.Command("systemctl", "disable", "config_tool.service").Run()

	exec.Command("systemctl", "enable", "nginx.service").Run()
	exec.Command("systemctl", "start", "nginx.service").Run()

	exec.Command("systemctl", "enable", "ghost_ghost.service").Run()
	exec.Command("systemctl", "start", "ghost_ghost.service").Run()


}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port)}

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

/// Web Handler Funcs

func showForm(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("form").Parse(form)
	t.Execute(w, nil)
}

func configure(done chan bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		host := r.FormValue("host")
		re := regexp.MustCompile(`https?:\/\/`)
		host = re.ReplaceAllString(host, "")
		host = strings.TrimSuffix(host, "/")
		// Update ghost config with new host
		err := exec.Command("sed", "-i", fmt.Sprintf("s/https:\\/\\/PLACEHOLDER/https:\\/\\/%s/", host), "/ghost/config.production.json").Run()

		// Configure nginx
		nginx := exec.Command("su", "ghost_manager", "-c", "/usr/local/bin/ghost setup nginx")
		nginx.Dir = "/ghost"
		nginx.Run()
		exec.Command("sed", "-i", "s/\\$scheme/https/", fmt.Sprintf("/etc/nginx/sites-enabled/%s.conf", host)).Run()

		t, _ := template.New("response").Parse(response)
		t.Execute(w, map[string]interface{}{"Host": host, "Success": err == nil})
		done <- true
	}
}
