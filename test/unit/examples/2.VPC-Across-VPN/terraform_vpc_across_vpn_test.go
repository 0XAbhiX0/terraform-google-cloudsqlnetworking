package vpc_across_vpn_test

import (
	"testing"
	"golang.org/x/exp/slices"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const terraformDirectoryPath = "../../../../examples/2.VPC-Across-VPN";
var host_project_id            = "pm-singleproject-20";
var service_project_id         = "pm-test-10-e90f";
var database_version 				   = "MYSQL_8_0"
var region                     = "us-central1";
var zone										   = "us-central1-a";
var user_project_id            = "pm-singleproject-30";
var cloudsql_instance_name     = "cn-sqlinstance10-test";
var network_name               = "cloudsql-easy";
var subnetwork_name            = "cloudsql-easy-subnet";
var subnetwork_ip_cidr         = "10.2.0.0/16"
var uservpc_network_name       = "cloudsql-user"
var uservpc_subnetwork_name    = "cloudsql-user-subnet"
var uservpc_subnetwork_ip_cidr = "10.10.30.0/24"
var test_dbname 							 = "test_db"
var user_region                = "us-west1"
var user_zone                  = "us-west1-a"
var deletion_protection 			 = false
var tfVars = map[string]interface{}{
	"host_project_id"            : host_project_id,
	"service_project_id"         : service_project_id,
	"database_version"           : database_version,
	"cloudsql_instance_name"     : cloudsql_instance_name,
	"region"                     : region,
	"zone"                       : zone,
	"create_network"             : true,
	"create_subnetwork"          : true,
	"network_name"               : network_name,
	"subnetwork_name"            : subnetwork_name, // this subnetwork will be created
	"subnetwork_ip_cidr"         : subnetwork_ip_cidr,
	"user_project_id"            : user_project_id,
	"user_region"                : user_region,
	"user_zone"                  : user_zone,
	"create_user_vpc_network"    : true,
	"create_user_vpc_subnetwork" : true,
	"uservpc_network_name"       : uservpc_network_name,
	"uservpc_subnetwork_name"    : uservpc_subnetwork_name,
	"uservpc_subnetwork_ip_cidr" : uservpc_subnetwork_ip_cidr,
	"test_dbname"                : test_dbname,
	"deletion_protection" 	     : deletion_protection,

}
// var backendConfig  						=  map[string]interface{}{
// 	"impersonate_service_account" : "iac-sa-test@pm-singleproject-20.iam.gserviceaccount.com",
// 	"bucket" 											: "pm-cncs-cloudsql-easy-networking",
// 	"prefix" 											: "test/example2",
//  }

func TestInitAndPlanRunWithTfVars(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars : tfVars,
		//BackendConfig : backendConfig,
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
		//VarFiles: [] string {"dev.tfvars" },
	})

	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	assert.Equal(t, 2, planExitCode)
}

func TestInitAndPlanRunWithoutTfVarsExpectFailureScenario(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		//BackendConfig : backendConfig,
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	assert.Equal(t, 1, planExitCode)
}

func TestResourcesCount(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars : tfVars,
		//BackendConfig : backendConfig,
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
		//VarFiles: [] string {"dev.tfvars" },
	})

	//plan *PlanStruct
	planStruct := terraform.InitAndPlan(t, terraformOptions)

	resourceCount := terraform.GetResourceCount(t, planStruct)
	assert.Equal(t,91,resourceCount.Add)
	assert.Equal(t,0,resourceCount.Change)
	assert.Equal(t,0,resourceCount.Destroy)
}

func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModulesAddress := [] string {"module.google_compute_instance.module.compute_instance","module.sql-db.module.mysql[0]","module.gce_sa","module.host_project_vpn","module.user_project_vpn","module.host-vpc","module.terraform_service_accounts","module.user-vpc","module.project_services.module.project_services","module.firewall_rules.module.firewall_rules","module.user_project_services.module.project_services","module.user_google_compute_instance.module.compute_instance","module.user_gce_sa","module.host_project.module.project_services","module.sql-db","module.user_firewall_rules.module.firewall_rules","module.user_google_compute_instance.module.instance_template","module.google_compute_instance.module.instance_template"}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars : tfVars,
		//BackendConfig : backendConfig,
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
		//VarFiles: [] string {"dev.tfvars" },
	})

	//plan *PlanStruct
	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	actualModuleAddress := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		if !slices.Contains(actualModuleAddress, element.ModuleAddress) && len(element.ModuleAddress) > 0 {
			actualModuleAddress = append(actualModuleAddress,element.ModuleAddress)
		}
	}
	if err != nil {
		print(err.Error())
	}
	assert.ElementsMatch(t, expectedModulesAddress, actualModuleAddress);
}
