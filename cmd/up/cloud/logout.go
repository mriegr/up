// Copyright 2021 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloud

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/upbound/up-sdk-go"

	"github.com/upbound/up/internal/config"
	"github.com/upbound/up/internal/upbound"
)

const (
	logoutPath = "/v1/logout"

	errLogoutFailed = "unable to logout"
)

// AfterApply sets default values in login after assignment and validation.
func (c *logoutCmd) AfterApply(upCtx *upbound.Context) error {
	var profile config.Profile
	var err error
	if upCtx.Profile == "" {
		_, profile, err = upCtx.Cfg.GetDefaultUpboundProfile()
		if err != nil {
			return err
		}
	} else {
		profile, err = upCtx.Cfg.GetUpboundProfile(upCtx.Profile)
		if err != nil {
			return err
		}
	}
	cfg, err := upbound.BuildSDKConfig(profile.Session, upCtx.Endpoint)
	if err != nil {
		return err
	}
	c.client = cfg.Client
	return nil
}

// logoutCmd invalidates a stored session token for a given profile.
type logoutCmd struct {
	client up.Client
}

// Run executes the logout command.
func (c *logoutCmd) Run(upCtx *upbound.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	req, err := c.client.NewRequest(ctx, http.MethodPost, logoutPath, "", nil)
	if err != nil {
		return errors.Wrap(err, errLogoutFailed)
	}
	// TODO(hasheddan): consider removing session token from config.
	return errors.Wrap(c.client.Do(req, nil), errLogoutFailed)
}
