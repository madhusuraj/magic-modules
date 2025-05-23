package discoveryengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type DiscoveryEngineOperationWaiter struct {
	Config    *transport_tpg.Config
	UserAgent string
	Project   string
	tpgresource.CommonOperationWaiter
}

func (w *DiscoveryEngineOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("cannot query operation, it's unset or nil")
	}

	opName := w.CommonOperationWaiter.Op.Name

	// Extract location from opName.  Example:
	// "projects/<project>/locations/<location>/operations/<operation-id>"
	// Requires error handling if not properly formatted.
	parts := strings.Split(opName, "/")
	if len(parts) < 5 || parts[2] != "locations" {
		return nil, fmt.Errorf("unexpected operation name format: %s", opName) // Handle invalid format appropriately
	}
	location := parts[3]

	basePath := strings.Replace(w.Config.DiscoveryEngineBasePath, "{{location}}", location, 1)
	url := fmt.Sprintf("%s%s", basePath, opName)

	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    w.Config,
		Method:    "GET",
		Project:   w.Project,
		RawURL:    url,
		UserAgent: w.UserAgent,
	})
}

func createDiscoveryEngineWaiter(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string) (*DiscoveryEngineOperationWaiter, error) {
	w := &DiscoveryEngineOperationWaiter{
		Config:    config,
		UserAgent: userAgent,
		Project:   project,
	}
	if err := w.CommonOperationWaiter.SetOp(op); err != nil {
		return nil, err
	}
	return w, nil
}

// nolint: deadcode,unused
func DiscoveryEngineOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createDiscoveryEngineWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	rawResponse := []byte(w.CommonOperationWaiter.Op.Response)
	if len(rawResponse) == 0 {
		return errors.New("`resource` not set in operation response")
	}
	return json.Unmarshal(rawResponse, response)
}

func DiscoveryEngineOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createDiscoveryEngineWaiter(config, op, project, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}
