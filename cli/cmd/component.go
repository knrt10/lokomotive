package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kinvolk/lokoctl/pkg/components"
	"github.com/kinvolk/lokoctl/pkg/k8sutil"
)

var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "Interact with components of a Lokomotive cluster",
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a component",
	Run:   runInstall,
}

var (
	namespace string
)

func init() {
	rootCmd.AddCommand(componentCmd)
	componentCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace where the component will be installed")
}

func runInstall(cmd *cobra.Command, args []string) {
	contextLogger := log.WithFields(log.Fields{
		"command": "lokoctl components install",
		"namespace": namespace,
		"args":    args,
	})

	if len(args) == 0 {
		contextLogger.Fatalf("Component name missing from command. Must be one of: %q", components.List())
	}

	c, err := components.Get(args[0])
	if err != nil {
		contextLogger.Fatalf("%q", err)
	}

	client, err := k8sutil.NewClientset(kubeconfig)
	if err != nil {
		contextLogger.Fatalf("Error in setting up Kubernetes client: %q", err)
	}

	if err = c.Install(client, namespace); err != nil {
		contextLogger.Fatalf("Installation of component %q failed: %q", c.Name(), err)
	}

}
