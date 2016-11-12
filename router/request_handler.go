package main

import (
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
)

const (
	healthCheckRoute       = "/status"
	wildcardHandlerPattern = "/"
)

var (
	statusCheckResponse = []byte("ok")
)

// RequestHandler creates an HTTP request handler that responds to all incoming
// router requests.
func RequestHandler(
	conf *config.Config,
	client db.Client,
	creds *config.Credentials,
	datadogClient *statsd.Client,
) func(http.ResponseWriter, *http.Request) {
	// Instantiate the IO module for use in package downloading and versioning.
	io := io.NewIO()

	// Instantiate the the github request service to pass into new
	// package requests.
	ghSvc := github.NewRequestService(github.RequestServiceArgs{
		Conf:       conf,
		Queryable:  client,
		ForIndexer: false,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		trackingArgs := datadog.TrackTranscationArgs{
			Tags: []string{
				"package-download",
				"external",
			},
			Client:    datadogClient,
			StartTime: time.Now(),
			EventInfo: []string{
				r.URL.Path, r.UserAgent(),
			},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: "package.download",
		}

		// Make sure that this isn't a simple health check before getting more
		// complicated.
		if r.URL.Path == healthCheckRoute {
			w.Write(statusCheckResponse)
			return
		}

		// First, create the necessary variables.
		var (
			pr  *packageRequest
			err error
		)

		// Create a new package request.
		if pr, err = newPackageRequest(newPackageRequestArgs{
			req:          r,
			downloadRefs: lib.FetchRefs,
			fetchFullSHA: github.FetchFullSHAFromPartialSHA,
			doHTTPHead:   github.DoHTTPHeadReq,
		}); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			defer datadog.TrackTranscation(trackingArgs)
			errors.RespondWithError(w, err)
			return
		}

		// Use the package request to respond.
		if err = pr.respond(respondToPackageRequestArgs{
			io:                    io,
			db:                    client,
			res:                   w,
			conf:                  conf,
			creds:                 creds,
			ghSvc:                 ghSvc,
			versionPackage:        versionAndArchivePackage,
			isPackageArchived:     isPackageArchived,
			recordPackageDownload: recordPackageDownload,
			recordPackageArchival: recordPackageArchival,
		}); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			defer datadog.TrackTranscation(trackingArgs)
			errors.RespondWithError(w, err)
			return
		}

		trackingArgs.AlertType = datadog.Success
		defer datadog.TrackTranscation(trackingArgs)
	}
}
