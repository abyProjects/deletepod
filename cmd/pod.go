package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var ErrMsg string =`
Error: pod needs all the below flags for execution:
Usage:
	deletepod pod [flags]

Flags:
	--name string        pod name to be deleted
	--namespace string   namespace of the pod
	--token string       token for authentication
Example:
	deletepod pod --name <pod_name> --namespace <default> --token <token/path_to_token_file>`

// variable that holds the data
var PodName string
var PodNamespace string
var Token string

// podCmd represents the pod command
var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "command to provide pod details",
	Long: `command to provide pod details.`,
	Run: func(cmd *cobra.Command, args []string) {
		isNameSet:=cmd.Flags().Lookup("name").Changed
		isNamespaceSet:=cmd.Flags().Lookup("namespace").Changed
		isTokenSet:=cmd.Flags().Lookup("token").Changed

		if isNameSet && isNamespaceSet && isTokenSet{
			log.Println("got value for all flags")
		}else{
			log.Fatalln(ErrMsg)
		}

		PodName, _ = cmd.Flags().GetString("name")
		PodNamespace, _ = cmd.Flags().GetString("namespace")
		Token, _ = cmd.Flags().GetString("token")

		if PodName != "" || PodNamespace != "" || Token != "" {
			log.Printf("podname: %s\tpodnamespace: %s\ttoken: %s", PodName, PodNamespace, Token)
		}else{
			log.Fatalln("invalid arguement passed")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(podCmd)

	podCmd.PersistentFlags().String("name", "", "pod name to be deleted")
	podCmd.PersistentFlags().String("namespace", "", "namespace of the pod")
	podCmd.PersistentFlags().String("token", "", "token for authentication")
	
}
