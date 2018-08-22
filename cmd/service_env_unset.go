package cmd

import (
	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	ECS "github.com/turnerlabs/fargate/ecs"
)

type ServiceEnvUnsetOperation struct {
	ServiceName string
	Keys        []string
}

func (o *ServiceEnvUnsetOperation) Validate() {
	if len(o.Keys) == 0 {
		console.IssueExit("No keys specified")
	}
}

var serviceEnvUnsetCmd = &cobra.Command{
	Use:   "unset --key <key-name> [--key <key-name>] ...",
	Short: "Unset environment variables",
	Long: `Unset environment variables

Unsets the environment variable specified via the --key flag. Specify --key with
a key name multiple times to unset multiple variables.`,
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceEnvUnsetOperation{
			ServiceName: getServiceName(),
		}

		operation.Keys = flagServiceEnvUnsetKeys
		operation.Validate()
		serviceEnvUnset(operation)
	},
}

var flagServiceEnvUnsetKeys []string

func init() {
	serviceEnvUnsetCmd.Flags().StringSliceVarP(&flagServiceEnvUnsetKeys, "key", "k", []string{}, "Environment variable keys to unset [e.g. KEY, NGINX_PORT]")

	serviceEnvCmd.AddCommand(serviceEnvUnsetCmd)
}

func serviceEnvUnset(operation *ServiceEnvUnsetOperation) {
	ecs := ECS.New(sess, getClusterName())
	service := ecs.DescribeService(operation.ServiceName)
	taskDefinitionArn := ecs.RemoveEnvVarsFromTaskDefinition(service.TaskDefinitionArn, operation.Keys)

	ecs.UpdateServiceTaskDefinition(operation.ServiceName, taskDefinitionArn)

	console.Info("Unset %s environment variables:", operation.ServiceName)

	for _, key := range operation.Keys {
		console.Info("- %s", key)
	}
}
