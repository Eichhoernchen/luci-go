// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package buildbucket

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/luci/luci-go/common/api/buildbucket/buildbucket/v1"
	"github.com/luci/luci-go/common/api/buildbucket/swarmbucket/v1"
	"github.com/luci/luci-go/milo/api/resp"
	"github.com/luci/luci-go/server/auth"
)

const (
	// StatusScheduled means a build is pending.
	StatusScheduled = "SCHEDULED"
	// StatusStarted means a build is executing.
	StatusStarted = "STARTED"
	// StatusCompleted means a build is completed (successfully or not).
	StatusCompleted = "COMPLETED"
)

// ParseTags parses buildbucket build tags to a map.
// Ignores tags that doesn't have colon (we don't have them in practice because
// buildbucket validates tags).
func ParseTags(tags []string) map[string]string {
	result := make(map[string]string, len(tags))
	for _, t := range tags {
		parts := strings.SplitN(t, ":", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func newSwarmbucketClient(c context.Context, server string) (*swarmbucket.Service, error) {
	c, _ = context.WithTimeout(c, time.Minute)
	t, err := auth.GetRPCTransport(c, auth.AsUser)
	if err != nil {
		return nil, err
	}
	client, err := swarmbucket.New(&http.Client{Transport: t})
	if err != nil {
		return nil, err
	}
	client.BasePath = fmt.Sprintf("https://%s/api/swarmbucket/v1/", server)
	return client, nil
}

func newBuildbucketClient(c context.Context, server string) (*buildbucket.Service, error) {
	c, _ = context.WithTimeout(c, time.Minute)
	t, err := auth.GetRPCTransport(c, auth.AsUser)
	if err != nil {
		return nil, err
	}
	client, err := buildbucket.New(&http.Client{Transport: t})
	if err != nil {
		return nil, err
	}
	client.BasePath = fmt.Sprintf("https://%s/api/buildbucket/v1/", server)
	return client, nil
}

// parseStatus converts a buildbucket build status to resp.Status.
func parseStatus(build *buildbucket.ApiBuildMessage) (resp.Status, error) {
	switch build.Status {
	case StatusScheduled:
		return resp.NotRun, nil

	case StatusStarted:
		return resp.Running, nil

	case StatusCompleted:
		switch build.Result {
		case "SUCCESS":
			return resp.Success, nil

		case "FAILURE":
			switch build.FailureReason {
			case "BUILD_FAILURE":
				return resp.Failure, nil
			default:
				return resp.InfraFailure, nil
			}

		case "CANCELED":
			return resp.InfraFailure, nil

		default:
			return 0, fmt.Errorf("unexpected buildbucket build result %q", build.Result)
		}

	default:
		return 0, fmt.Errorf("unexpected buildbucket build status %q", build.Status)
	}
}

// getChangeList tries to extract CL information from a buildbucket build.
func getChangeList(
	build *buildbucket.ApiBuildMessage, params *buildParameters,
	resultDetails *resultDetails) (result *resp.Commit) {

	prop := &params.Properties
	switch prop.PatchStorage {
	case "rietveld":
		if prop.RietveldURL != "" && prop.Issue != 0 {
			result = &resp.Commit{
				RequestRevision: &resp.Link{Label: prop.Revision},
				Changelist: &resp.Link{
					Label: fmt.Sprintf("Rietveld CL %d", prop.Issue),
					URL:   fmt.Sprintf("%s/%d/#ps%d", prop.RietveldURL, prop.Issue, prop.PatchSet),
				},
			}
			if resultDetails.Properties.GotRevision != "" {
				// TODO(hinoka): Figure out the full URL for these revisions, add it
				// to the URL field.
				result.Revision = &resp.Link{Label: resultDetails.Properties.GotRevision}
			}
		}

	case "gerrit":
		if prop.GerritPatchURL != "" && prop.GerritPatchIssue != 0 {
			path := fmt.Sprintf("%d", prop.GerritPatchIssue)
			if prop.GerritPatchSet != 0 {
				path = fmt.Sprintf("%d/%d", prop.GerritPatchIssue, prop.GerritPatchSet)
			}
			result = &resp.Commit{
				Changelist: &resp.Link{
					Label: fmt.Sprintf("Gerrit CL %d", prop.GerritPatchIssue),
					URL: fmt.Sprintf(
						"%s/c/%s", prop.GerritPatchURL, path),
				},
			}
		}
	}

	if result != nil && len(params.Changes) != 0 {
		// tryjobs have one change and it is the CL author
		result.AuthorEmail = params.Changes[0].Author.Email
	}

	return
}
