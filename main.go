/*
Copyright 2022 CMU-SV.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/utils/clock"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	selfsignedissuerv1alpha1 "github.com/kxk-4498/Venafi-test-wizard/api/v1alpha1"
	"github.com/kxk-4498/Venafi-test-wizard/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(selfsignedissuerv1alpha1.AddToScheme(scheme))

	utilruntime.Must(cmapi.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var disableApprovedCheck bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. "+"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&disableApprovedCheck, "disable-approved-check", false, "Disables waiting for CertificateRequests to have an approved condition before signing.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.ChaosIssuerReconciler{
		Client:   mgr.GetClient(),
		Clock:    clock.RealClock{},
		Recorder: mgr.GetEventRecorderFor("ChaosIssuer-controller"),
		Log:      ctrl.Log.WithName("controllers").WithName("ChaosIssuer"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ChaosIssuer")
		os.Exit(1)
	}
	if err = (&controllers.ChaosIssuerReconciler{
		Client:   mgr.GetClient(),
		Clock:    clock.RealClock{},
		Recorder: mgr.GetEventRecorderFor("ChaosIssuer-controller"),
		Log:      ctrl.Log.WithName("controllers").WithName("ChaosIssuer"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ChaosClusterIssuer")
		os.Exit(1)
	}
	if err = (&controllers.CertificateRequestReconciler{
		Client:                 mgr.GetClient(),
		Clock:                  clock.RealClock{},
		Log:                    ctrl.Log.WithName("controllers").WithName("CertificateRequest"),
		CheckApprovedCondition: !disableApprovedCheck,
		Recorder:               mgr.GetEventRecorderFor("certificaterequests-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CertificateRequest")
		os.Exit(1)
	}

	//+kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
