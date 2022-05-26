package main

import (
	"e2e-test/cmd"
)

/*
	1. Check if sbx cluster1 is online,
		if it is, check 2, then 3.
        If no cluster is available - fail pipeline
    2. Once a cluster is NOT found, spin up that cluster using context.Background & terraform-exec https://github.com/hashicorp/terraform-exec
    3. In the background using context.Background, pull the code for the e2e testing
    4. once the cluster is up - context is complete, deploy all end-2-end code
    5. Query any DNS records as required and once they return, start tests.
    6. Check all end points to ensure they are up
	7. Attempt to scale cluster up to max nodes
	8. Anything else you can think of.
*/

func main() {
	//viper.SetConfigName("config")         // name of config file (without extension)
	//viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	//viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	//viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	//viper.AddConfigPath(".")              // optionally look for config in the working directory
	//err := viper.ReadInConfig()           // Find and read the config file
	//if err != nil {                       // Handle errors reading the config file
	//	panic(fmt.Errorf("Fatal error config file: %w \n", err))
	//}
	cmd.Execute()
}
