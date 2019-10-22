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

type ViewGenerator struct {
	SnowflakeService
}

func (g ViewGenerator) createResources(viewList []view) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, view := range viewList {
		resources = append(resources, terraform_utils.NewResource(
			fmt.Sprintf("%s|%s|%s", view.DatabaseName.String, view.SchemaName.String, view.Name.String),
			fmt.Sprintf("%s__%s__%s", view.DatabaseName.String, view.SchemaName.String, view.Name.String),
			"snowflake_view",
			"snowflake",
			map[string]string{},
			[]string{},
			map[string]interface{}{
				"database":  view.DatabaseName.String,
				"statement": view.Text.String,
			},
		))
	}
	return resources
}

func (g *ViewGenerator) InitResources() error {
	db, err := g.generateService()
	if err != nil {
		return err
	}
	output, err := db.ListViews()
	if err != nil {
		return err
	}
	g.Resources = g.createResources(output)
	return nil
}