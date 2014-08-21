package api


import (
	"fmt"
	"net/http"
	"atlantis/routerapi/zk"
	"atlantis/routerapi/auth"
)

func Healthz(w http.ResponseWriter, r *http.Request){

	//check zk health
	zkOk := zk.IsZkConnOk()
	//check auth health
	authOk := auth.GetAuthorizer().AuthOk()

	if !zkOk || !authOk {
		w.Header().Add("Server-Status", "CRITICAL")
		fmt.Fprintf(w, "%s", "CRITICAL")
	} else {

		w.Header().Add("Server-Status", "OK")	
		fmt.Fprintf(w, "%s", "OK")
	}

}
