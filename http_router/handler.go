package http_router

import (
	"context"
	"encoding/json"
	"fmt"
	"k8stool/http_router/message"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/zerolog/log"

	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var cli *wfclientset.Clientset

// var bearToken string

func InitK8SClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf(err.Error())
		panic(err.Error())
	}
	// bearToken = config.BearerToken
	// log.Printf("beakToken: %s", bearToken)
	cli = wfclientset.NewForConfigOrDie(config)
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Printf(err.Error())
	// 	return err
	// }
	// cli = clientset
	return nil
}

func GetWorkFlowsStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	ns := mux.Vars(r)["namespace"]
	// if bearToken != "" {
	// 	token := r.Header.Get("Authorization")
	// 	log.Printf("http request token: %s", token)
	// 	if bearToken != token {
	// 		err := errors.New("token not valid for running mode")
	// 		log.Printf(err.Error())
	// 		WrapResponse(w, message.CommonResponse{
	// 			Status:   message.FAIL,
	// 			Messages: err.Error(),
	// 		})
	// 		return
	// 	}
	// }

	var req []message.GetWorkFlowsStatusRequest
	err := ResolveRequest(r, &req)
	if err != nil {
		log.Printf(err.Error())
		WrapResponse(w, message.CommonResponse{
			Status:   message.FAIL,
			Messages: err.Error(),
		})
		return
	}
	res := make([]message.WorkFlowInfo, 0)
	for _, r := range req {
		wf, err := cli.ArgoprojV1alpha1().Workflows(ns).Get(context.Background(), r.WorkflowName, v1.GetOptions{})
		if err != nil {
			log.Printf(err.Error())
			WrapResponse(w, message.CommonResponse{
				Status:   message.FAIL,
				Messages: err.Error(),
			})
			return
		}
		item := message.WorkFlowInfo{}
		item.StatusPhase = string(wf.Status.Phase)
		item.WorkflowName = wf.Name
		res = append(res, item)
	}
	messageResp := message.CommonResponse{}
	messageResp.Status = message.SUCCESS
	messageResp.Result = res
	WrapResponse(w, messageResp)
}

func ResolveRequest(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(req)
	return err
}

func WrapResponse(w http.ResponseWriter, msg interface{}) {
	jsonBytes, _ := json.Marshal(msg)
	fmt.Fprintln(w, string(jsonBytes))
}
