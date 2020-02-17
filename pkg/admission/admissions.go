package main

import (
	"fmt"
	"github.com/jfrog/kubenab/log"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	json "github.com/json-iterator/go"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"api_endpoint"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_milliseconds",
		Help: "The HTTP request Duration in Milliseconds",

		// Latency Distribution
		// 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	}, []string{"api_method"})
)

var (
	dockerRegistryUrl     = os.Getenv("DOCKER_REGISTRY_URL")
	replaceRegistryUrl    = os.Getenv("REPLACE_REGISTRY_URL")
	registrySecretName    = os.Getenv("REGISTRY_SECRET_NAME")
	whitelistRegistries   = os.Getenv("WHITELIST_REGISTRIES")
	whitelistNamespaces   = os.Getenv("WHITELIST_NAMESPACES")
	whitelistedNamespaces = strings.Split(whitelistNamespaces, ",")
	whitelistedRegistries = strings.Split(whitelistRegistries, ",")
)

type patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func mutateAdmissionReviewHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.With(prometheus.Labels{"api_endpoint": "mutate"}).Inc()

	// log Request Duration
	promTimer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		httpRequestDuration.WithLabelValues("mutate").Observe(v * 1000) // add Milliseconds
	}))
	defer promTimer.ObserveDuration()

	log.Printf("Serving request: %s", r.URL.Path)
	//set header
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debugln(data)

	ar := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(data, &ar); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namespace := ar.Request.Namespace
	log.Printf("AdmissionReview Namespace is: %s", namespace)

	admissionResponse := v1beta1.AdmissionResponse{Allowed: false}
	patches := []patch{}
	if !contains(whitelistedNamespaces, namespace) {
		pod := v1.Pod{}
		if err := json.Unmarshal(ar.Request.Object.Raw, &pod); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var patchPath strings.Builder

		// Handle Containers
		for i, container := range pod.Spec.Containers {
			createPatch := handleContainer(&container, dockerRegistryUrl)
			if createPatch {
				patchPath.Reset()
				// TODO: Check the returned Error
				_, _ = patchPath.WriteString("/spec/containers/")
				_, _ = patchPath.WriteString(strconv.Itoa(i))
				_, _ = patchPath.WriteString("/image")

				patches = append(patches, patch{
					Op:    "replace",
					Path:  patchPath.String(),
					Value: container.Image,
				})
			}
		}

		// Handle init containers
		for i, container := range pod.Spec.InitContainers {
			createPatch := handleContainer(&container, dockerRegistryUrl)
			if createPatch {
				patchPath.Reset()

				// TODO: Check the returned Error
				_, _ = patchPath.WriteString("/spec/initContainers/")
				_, _ = patchPath.WriteString(strconv.Itoa(i))
				_, _ = patchPath.WriteString("/image")

				patches = append(patches, patch{
					Op:    "replace",
					Path:  patchPath.String(),
					Value: container.Image,
				})
			}
		}
	} else {
		log.Printf("Namespace is %s Whitelisted", namespace)
	}

	admissionResponse.Allowed = true
	if len(patches) > 0 {
		// add image pull secret patch if User has added it
		if len(registrySecretName) > 0 {
			patches = append(patches, patch{
				Op:   "add",
				Path: "/spec/imagePullSecrets",
				Value: []v1.LocalObjectReference{
					v1.LocalObjectReference{
						Name: registrySecretName,
					},
				},
			})
		}

		patchContent, err := json.Marshal(patches)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		admissionResponse.Patch = patchContent
		pt := v1beta1.PatchTypeJSONPatch
		admissionResponse.PatchType = &pt
	}

	ar = v1beta1.AdmissionReview{
		Response: &admissionResponse,
	}

	data, err = json.Marshal(ar)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleContainer(container *v1.Container, dockerRegistryUrl string) bool {
	log.Println("Container Image is", container.Image)

	if !containsRegisty(whitelistedRegistries, container.Image) {
		message := fmt.Sprintf("Image is not being pulled from Private Registry: %s", container.Image)
		log.Printf(message)

		imageParts := strings.Split(container.Image, "/")
		newImage := ""

		// pre-pend new Docker Registry Domain
		repRegUrl, _ := strconv.ParseBool(replaceRegistryUrl) // we do not need to check for errors here, since we have done this already in checkArguments()
		if (len(imageParts) < 3) || !repRegUrl {
			newImage = dockerRegistryUrl + "/" + container.Image
		} else {
			imageParts[0] = dockerRegistryUrl
			newImage = strings.Join(imageParts, "/")
		}
		log.Printf("Changing image registry to: %s", newImage)

		container.Image = newImage
		return true
	} else {
		log.Printf("Image is being pulled from Private Registry: %s", container.Image)
	}
	return false
}

func validateAdmissionReviewHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.With(prometheus.Labels{"api_endpoint": "validate"}).Inc()

	// log Request Duration
	promTimer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		httpRequestDuration.WithLabelValues("validate").Observe(v * 1000) // add Milliseconds
	}))
	defer promTimer.ObserveDuration()

	log.Printf("Serving request: %s", r.URL.Path)
	//set header
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug(data)

	ar := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(data, &ar); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namespace := ar.Request.Namespace
	log.Printf("AdmissionReview Namespace is: %s", namespace)

	admissionResponse := v1beta1.AdmissionResponse{Allowed: false}
	if !contains(whitelistedNamespaces, namespace) {
		pod := v1.Pod{}
		if err := json.Unmarshal(ar.Request.Object.Raw, &pod); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Handle containers
		for _, container := range pod.Spec.Containers {
			log.Println("Container Image is", container.Image)

			if !containsRegisty(whitelistedRegistries, container.Image) {
				message := fmt.Sprintf("Image is not being pulled from Private Registry: %s", container.Image)
				log.Printf(message)
				admissionResponse.Result = getInvalidContainerResponse(message)
				goto done
			} else {
				log.Printf("Image is being pulled from Private Registry: %s", container.Image)
				admissionResponse.Allowed = true
			}
		}

		// Handle init containers
		for _, container := range pod.Spec.InitContainers {
			log.Println("Init Container Image is", container.Image)

			if !containsRegisty(whitelistedRegistries, container.Image) {
				message := fmt.Sprintf("Image is not being pulled from Private Registry: %s", container.Image)
				log.Printf(message)
				admissionResponse.Result = getInvalidContainerResponse(message)
				goto done
			} else {
				log.Printf("Image is being pulled from Private Registry: %s", container.Image)
				admissionResponse.Allowed = true
			}
		}
	} else {
		log.Printf("Namespace is %s Whitelisted", namespace)
		admissionResponse.Allowed = true
	}

done:
	ar = v1beta1.AdmissionReview{
		Response: &admissionResponse,
	}

	data, err = json.Marshal(ar)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getInvalidContainerResponse(message string) *metav1.Status {
	return &metav1.Status{
		Reason: metav1.StatusReasonInvalid,
		Details: &metav1.StatusDetails{
			Causes: []metav1.StatusCause{
				{Message: message},
			},
		},
	}
}

// if current namespace is part of whitelisted namespaces
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str || strings.Contains(a, str) {
			return true
		}
	}
	return false
}

// if current registry is part of whitelisted registries
func containsRegisty(arr []string, str string) bool {
	for _, a := range arr {
		if a == str || strings.Contains(str, a) {
			return true
		}
	}
	return false
}

// ping responds to the request with a plain-text "Ok" message.
func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	fmt.Fprintf(w, "Ok")
}

// check if all (required) Arguments are set and valid
func checkArguments() {
	if len(dockerRegistryUrl) == 0 {
		log.Fatalln("Environment Variable 'DOCKER_REGISTRY_URL' not set")
	}

	if len(replaceRegistryUrl) == 0 {
		log.Fatalln("Environment Variable 'REPLACE_REGISTRY_URL' not set")
	}

	_, err := strconv.ParseBool(replaceRegistryUrl)
	if err != nil {
		log.Fatalln("Invalid Value in Environment Variable 'REPLACE_REGISTRY_URL'")
	}
}
