// Copyright 2019 The Terraformer Authors.
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
package snowflake

import (
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

type RoleGenerator struct {
	SnowflakeService
}

func (g RoleGenerator) createResources(roleList []role) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, role := range roleList {
		resources = append(resources, terraform_utils.NewSimpleResource(
			role.Name.String,
			role.Name.String,
			"snowflake_role",
			"snowflake",
			[]string{}))
	}
	return resources
}

func (g *RoleGenerator) InitResources() error {
	db, err := g.generateService()
	if err != nil {
		return err
	}
	output, err := db.ListRoles()
	if err != nil {
		return err
	}
	g.Resources = g.createResources(output)
	return nil
}