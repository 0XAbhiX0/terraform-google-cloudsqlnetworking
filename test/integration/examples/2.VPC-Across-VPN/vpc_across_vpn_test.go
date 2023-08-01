package vpc_across_vpn_test

import (
	"fmt"
	"time"
	"testing"
	"github.com/tidwall/gjson"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/shell"
)

const terraformDirectoryPath   = "../../../../../cloudsql-easy-networking/examples/2.VPC-Across-VPN";
var host_project_id          = "pm-singleproject-20";
var service_project_id       = "pm-test-10-e90f";
var user_project_id          = "pm-singleproject-30";
var cloudsql_instance_name   = "cn-sqlinstance10-test";
var subnetwork_name          = "cloudsql-easy-subnet";
var region                   = "us-central1";
var test_dbname              = "test_db"
var database_version         = "MYSQL_8_0"

// name the function as Test*
func TestMySqlPrivateAndVPNModule(t *testing.T) {
	host_project_id          = "pm-singleproject-20";
	service_project_id       = "pm-test-10-e90f";
	region                   = "us-central1";
	user_project_id          = "pm-singleproject-30";
	cloudsql_instance_name   = "cn-sqlinstance10-test";
	network_name             := "cloudsql-easy";
	subnetwork_name          = "cloudsql-easy-subnet";

	tfVars := map[string]interface{}{

		"cloudsql_instance_name" : cloudsql_instance_name,
		"region"                 : region,
		"network_name"           : network_name,
		"subnetwork_name"        : subnetwork_name, // this subnetwork will be created
		"test_dbname"            : test_dbname,
		"host_project_id"        : host_project_id,
    "service_project_id"     : service_project_id,
    "user_project_id"        : user_project_id,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		//PlanFilePath: "./plan",
		NoColor: true,
		SetVarsAfterVarFiles: true,
		VarFiles: [] string {"dev.tfvars" },
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	//wait for 60 seconds to let resource acheive stable state
	time.Sleep(60 * time.Second)


	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, network_name, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSqlInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkId := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",host_project_id,region,subnetwork_name)
	assert.Equal(t,subnetworkId , output)

	// Validate if SQL instance wih private IP is up and running
	text := "sql"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err := shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result := gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("===Error %s Encountered while executing %s", err ,text)
	}
	fmt.Println(" ========= Verify if public IP is disabled ========= ")
	assert.Equal(t, false, gjson.Get(result.String(),"settings.ipConfiguration.ipv4Enabled").Bool())
	fmt.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(),"state").String())

	// Validate if VPN tunnels are up & running with Established Connection
	fmt.Println(" ====================================================== ")
	fmt.Println(" ========= Verify VPN Tunnel ========= ")

	var vpnTunnelName = []string { "ha-vpn-tunnel1","ha-vpn-tunnel2","ha-vpn-tunnel3","ha-vpn-tunnel4"}
	var projectId = ""
	for _, v := range vpnTunnelName {
		if v == "ha-vpn-tunnel1" || v=="ha-vpn-tunnel2" {
			projectId = host_project_id;
		} else {
			projectId = user_project_id;
		}
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"compute","vpn-tunnels","describe",v,"--project",projectId,"--region",region,"--format=json"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result = gjson.Parse(op)
		if err != nil {
			fmt.Sprintf("===Error %s Encountered while executing %s", err ,text)
		}
		fmt.Printf(" \n========= validating tunnel %s ============\n",v);
		fmt.Println(" ========= check if tunnel is up & running ========= ",)
		assert.Equal(t, "Tunnel is up and running.", gjson.Get(result.String(),"detailedStatus").String())
		fmt.Println(" ========= check if connection is established ========= ")
		assert.Equal(t, "ESTABLISHED", gjson.Get(result.String(),"status").String())
	}

	//Iterate through list of database to ensure a new db was created
	fmt.Println(" ====================================================== ")
	fmt.Println(" =========== Verify DB Creation =========== ")
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err = shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("======= Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}

func TestUsingExistingNetworkMySqlPrivateAndVPNModule(t *testing.T) {
	host_project_id          = "pm-host-networking";
  service_project_id       = "pm-service1-networking";
  user_project_id          = "pm-userproject-networking";
	network_name             = "host-cloudsql-easy";
	subnetwork_name          = "host-cloudsql-easy-subnet";
	region                   = "us-central1";
	user_region              = "us-west1";
	user_zone                = "us-west1-a";
	uservpc_network_name     = "user-cloudsql-easy";
	uservpc_subnetwork_name  = "user-cloudsql-easy-subnet";

	tfVars := map[string]interface{}{
		"host_project_id"            : host_project_id,
    "service_project_id"         : service_project_id,
		"region"                     : region,
		"create_subnetwork"          : false,
		"create_network"             : false,
		"network_name"               : network_name,
		"subnetwork_name"            : subnetwork_name,
		"cloudsql_instance_name"     : cloudsql_instance_name,
		"database_version"           : database_version,
		"test_dbname"                : test_dbname,
    "user_project_id"            : user_project_id,
		"user_region"                : user_region,
    "user_zone"                  : user_zone,
		"create_user_vpc_network"    : false,
    "create_user_vpc_subnetwork" : false,
    "uservpc_network_name"       : uservpc_network_name,
    "uservpc_subnetwork_name"    : uservpc_subnetwork_name,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		//PlanFilePath: "./plan",
		NoColor: true,
		SetVarsAfterVarFiles: true,
		VarFiles: [] string {"dev.tfvars" },
	})


	//validate if the VPC already exists, else create one outside of the terraform


	//validate if the subnet already exists, else create one outside of the terraform



	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	//wait for 60 seconds to let resource acheive stable state
	time.Sleep(60 * time.Second)


	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, network_name, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSqlInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkId := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",host_project_id,region,subnetwork_name)
	assert.Equal(t,subnetworkId , output)

	// Validate if SQL instance wih private IP is up and running
	text := "sql"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err := shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result := gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("===Error %s Encountered while executing %s", err ,text)
	}
	fmt.Println(" ========= Verify if public IP is disabled ========= ")
	assert.Equal(t, false, gjson.Get(result.String(),"settings.ipConfiguration.ipv4Enabled").Bool())
	fmt.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(),"state").String())

	// Validate if VPN tunnels are up & running with Established Connection
	fmt.Println(" ====================================================== ")
	fmt.Println(" ========= Verify VPN Tunnel ========= ")

	var vpnTunnelName = []string { "ha-vpn-tunnel1","ha-vpn-tunnel2","ha-vpn-tunnel3","ha-vpn-tunnel4"}
	var projectId = ""
	for _, v := range vpnTunnelName {
		if v == "ha-vpn-tunnel1" || v=="ha-vpn-tunnel2" {
			projectId = host_project_id;
		} else {
			projectId = user_project_id;
		}
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"compute","vpn-tunnels","describe",v,"--project",projectId,"--region",region,"--format=json"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result = gjson.Parse(op)
		if err != nil {
			fmt.Sprintf("===Error %s Encountered while executing %s", err ,text)
		}
		fmt.Printf(" \n========= validating tunnel %s ============\n",v);
		fmt.Println(" ========= check if tunnel is up & running ========= ",)
		assert.Equal(t, "Tunnel is up and running.", gjson.Get(result.String(),"detailedStatus").String())
		fmt.Println(" ========= check if connection is established ========= ")
		assert.Equal(t, "ESTABLISHED", gjson.Get(result.String(),"status").String())
	}

	//Iterate through list of database to ensure a new db was created
	fmt.Println(" ====================================================== ")
	fmt.Println(" =========== Verify DB Creation =========== ")
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err = shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("======= Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}
