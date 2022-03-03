// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
)

// CreateOrUpdate indicates an expected call of Network interface CreateOrUpdate.
func (c NetworkInterfaceClient) CreateOrUpdate(ctx context.Context, resourceGroupName, name string, parameters network.Interface) (*network.Interface, error) {
	future, err := c.client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return nil, err
	}
	if err := future.WaitForCompletionRef(ctx, c.client.Client); err != nil {
		return nil, err
	}
	nic, err := future.Result(c.client)
	if err != nil {
		return nil, err
	}
	return &nic, nil
}

// Get will get a Network interface.
func (c NetworkInterfaceClient) Get(ctx context.Context, resourceGroupName string, name string, expander string) (*network.Interface, error) {
	nic, err := c.client.Get(ctx, resourceGroupName, name, expander)
	if err != nil {
		return nil, err
	}
	return &nic, nil
}

// Delete will delete a Network interface.
func (c NetworkInterfaceClient) Delete(ctx context.Context, resourceGroupName, name string) error {
	future, err := c.client.Delete(ctx, resourceGroupName, name)
	if err != nil {
		return err
	}
	if err := future.WaitForCompletionRef(ctx, c.client.Client); err != nil {
		return err
	}
	result, err := future.Result(c.client)
	if err != nil {
		return err
	}
	if result.StatusCode == http.StatusOK || result.StatusCode == http.StatusAccepted || result.StatusCode == http.StatusNoContent || result.StatusCode == http.StatusNotFound {
		return nil
	}
	return fmt.Errorf("deletion of network interfaces %s failed. statuscode=%d", name, result.StatusCode)
}