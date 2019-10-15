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
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

type DatabaseGrantGenerator struct {
	SnowflakeService
}

func (g DatabaseGrantGenerator) createResources(databaseGrantList []databaseGrant) []terraform_utils.Resource {
	type tfGrant struct {
		DB        string
		Privilege string
		Roles     []string
		Shares    []string
	}
	groupedResources := map[string]*tfGrant{}
	for _, grant := range databaseGrantList {
                // TODO(ad): Fix this csv delimited when fixed in the provider. We should use the same functionality.
		id := fmt.Sprintf("%v|||%v", grant.Name.String, grant.Privilege.String)
		_, ok := groupedResources[id]
		if !ok {
			groupedResources[id] = &tfGrant{
				DB:        grant.Name.String,
				Privilege: grant.Privilege.String,
				Roles:     []string{},
				Shares:    []string{},
			}
		}
		tfGrant := groupedResources[id]
		if grant.GrantedTo.String == "ROLE" {
			tfGrant.Roles = append(tfGrant.Roles, grant.GranteeName.String)
		}
		if grant.GrantedTo.String == "SHARE" {
			tfGrant.Shares = append(tfGrant.Shares, grant.GranteeName.String)
		}
	}
	var resources []terraform_utils.Resource
	for id, grant := range groupedResources {
		resources = append(resources, terraform_utils.NewResource(
			id,
			fmt.Sprintf("%v_%v", grant.DB, grant.Privilege),
			"snowflake_database_grant",
			"snowflake",
			map[string]string{
				"privilege": grant.Privilege,
			},
			[]string{},
			map[string]interface{}{
				"roles":  grant.Roles,
				"shares": grant.Shares,
			},
		))
	}
	return resources
}

func (g *DatabaseGrantGenerator) InitResources() error {
	db, err := g.generateService()
	if err != nil {
		return err
	}
	databases, err := db.ListDatabases()
	if err != nil {
		return err
	}
	allGrants := []databaseGrant{}
	for _, database := range databases {
		if database.Origin.String == "" {
			// Provider does not support grants on imported databases yet
			continue
		}
		grants, err := db.ListDatabaseGrants(database)
		if err != nil {
			return err
		}
		allGrants = append(allGrants, grants...)
	}
	g.Resources = g.createResources(allGrants)
	return nil
}