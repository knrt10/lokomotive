// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage/driver"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kinvolk/lokomotive/pkg/components"
	"github.com/kinvolk/lokomotive/pkg/k8sutil"
)

func ensureNamespaceExists(name string, kubeconfigPath string) error {
	kubeconfig, err := ioutil.ReadFile(kubeconfigPath) // #nosec G304
	if err != nil {
		return fmt.Errorf("reading kubeconfig file: %w", err)
	}

	cs, err := k8sutil.NewClientset(kubeconfig)
	if err != nil {
		return fmt.Errorf("creating clientset: %w", err)
	}

	if name == "" {
		return fmt.Errorf("namespace name can't be empty")
	}

	// Ensure the namespace in which we create release and resources exists.
	_, err = cs.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// RunAdmissionWebhook is used to create all required resources to start mutating admission controller.
func RunAdmissionWebhook(kubeconfigAbsoutePath string) error {
	dir, err := ioutil.TempDir("", "mutating")
	if err != nil {
		return fmt.Errorf("error in creating directory: %v", err)
	}

	defer func() {
		err = os.RemoveAll(dir)
		if err != nil {
			fmt.Printf("error when removing files: %v", err)
		}
	}()

	//nolint:lll
	fileURL := "https://gist.github.com/knrt10/274228a6626894f97466751a30f82b87/archive/809b3825c0411f72330031a5f3f12e4e947f647d.zip"

	err = downloadFile(dir+"/install.zip", fileURL)
	if err != nil {
		return fmt.Errorf("error in dowloading zip file: %v", err)
	}

	_, err = unzip(dir+"/install.zip", dir)
	if err != nil {
		return fmt.Errorf("error in unzipping file: %v", err)
	}

	_, err = exec.Command("chmod", "777", dir+"/webhook-run.sh").Output() //nolint:gosec
	if err != nil {
		return fmt.Errorf("changing file permission failed failed: %v", err)
	}

	output, err := exec.Command("/bin/sh", dir+"/webhook-run.sh", "--kubeconfigpath", kubeconfigAbsoutePath).Output() //nolint
	if err != nil {
		fmt.Println("testing ci output", string(output))
		return fmt.Errorf("applying mutating admission webhook failed: %v", err)
	}

	return nil
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("error when closing response body: %v", err)
		}
	}()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer func() {
		err = out.Close()
		if err != nil {
			fmt.Printf("error when closing file: %v", err)
		}
	}()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//nolint
// unzip is used to unzip the desired zip file.
func unzip(src string, dest string) ([]string, error) {
	var filenames []string //nolint:gosec

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}

	defer func() {
		err = r.Close()
		if err != nil {
			fmt.Printf("error when closing file: %v", err)
		}
	}()

	for _, f := range r.File {
		fileSlicePath := strings.Split(f.Name, "/")
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, fileSlicePath[1], "/")

		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return filenames, err
			}

			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)
		// Close the file without defer to close before next iteration of loop
		err = outFile.Close()
		if err != nil {
			fmt.Printf("error in closing file %v:", err)
		}

		err = rc.Close()
		if err != nil {
			fmt.Printf("error in closing file %v:", err)
		}

		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}

// InstallComponent installs given component using given kubeconfig as a Helm release using a Helm client.
func InstallComponent(c components.Component, kubeconfig string) error {
	name := c.Metadata().Name
	ns := c.Metadata().Namespace

	if err := ensureNamespaceExists(ns, kubeconfig); err != nil {
		return fmt.Errorf("failed ensuring that namespace %q for component %q exists: %w", ns, name, err)
	}

	actionConfig, err := HelmActionConfig(ns, kubeconfig)
	if err != nil {
		return fmt.Errorf("failed preparing helm client: %w", err)
	}

	chart, err := chartFromComponent(c)
	if err != nil {
		return err
	}

	if err := chart.Validate(); err != nil {
		return fmt.Errorf("chart is invalid: %w", err)
	}

	exists, err := ReleaseExists(*actionConfig, name)
	if err != nil {
		return fmt.Errorf("failed checking if component is installed: %w", err)
	}

	wait := c.Metadata().Helm.Wait

	helmAction := &helmAction{
		releaseName:  name,
		chart:        chart,
		actionConfig: actionConfig,
		wait:         wait,
	}

	if !exists {
		return install(helmAction, ns)
	}

	return upgrade(helmAction)
}

type helmAction struct {
	releaseName  string
	chart        *chart.Chart
	actionConfig *action.Configuration
	wait         bool
}

func install(helmAction *helmAction, namespace string) error {
	install := action.NewInstall(helmAction.actionConfig)
	install.ReleaseName = helmAction.releaseName
	install.Namespace = namespace

	// Currently, we install components one-by-one, in the order how they are
	// defined in the configuration and we do not support any dependencies between
	// the components.
	//
	// If it is critical for component to have it's dependencies ready before it is
	// installed, all dependencies should set Wait field to 'true' in components.HelmMetadata
	// struct.
	//
	// The example of such dependency is between prometheus-operator and openebs-storage-class, where
	// both openebs-operator and openebs-storage-class components must be fully functional, before
	// prometheus-operator is deployed, otherwise it won't pick the default storage class.
	install.Wait = helmAction.wait

	if _, err := install.Run(helmAction.chart, map[string]interface{}{}); err != nil {
		return fmt.Errorf("installing release failed: %w", err)
	}

	return nil
}

func upgrade(helmAction *helmAction) error {
	upgrade := action.NewUpgrade(helmAction.actionConfig)
	upgrade.Wait = helmAction.wait
	upgrade.RecreateResources = true

	if _, err := upgrade.Run(helmAction.releaseName, helmAction.chart, map[string]interface{}{}); err != nil {
		return fmt.Errorf("upgrading release failed: %w", err)
	}

	return nil
}

// HelmActionConfig creates initialized Helm action configuration.
func HelmActionConfig(ns string, kubeconfig string) (*action.Configuration, error) {
	actionConfig := &action.Configuration{}

	// TODO: Add some logging implementation? We currently just pass an empty function for logging.
	kubeConfig := kube.GetConfig(kubeconfig, "", ns)
	logF := func(format string, v ...interface{}) {}

	if err := actionConfig.Init(kubeConfig, ns, "secret", logF); err != nil {
		return nil, fmt.Errorf("failed initializing helm: %w", err)
	}

	return actionConfig, nil
}

// ReleaseExists checks if given Helm release exists.
func ReleaseExists(actionConfig action.Configuration, name string) (bool, error) {
	histClient := action.NewHistory(&actionConfig)
	histClient.Max = 1

	_, err := histClient.Run(name)
	if err != nil && err != driver.ErrReleaseNotFound {
		return false, fmt.Errorf("failed checking for chart history: %w", err)
	}

	return err != driver.ErrReleaseNotFound, nil
}
