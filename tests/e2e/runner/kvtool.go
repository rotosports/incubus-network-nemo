package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type KvtoolRunnerConfig struct {
	NemoConfigTemplate string

	ImageTag   string
	IncludeIBC bool

	EnableAutomatedUpgrade  bool
	NemoUpgradeName         string
	NemoUpgradeHeight       int64
	NemoUpgradeBaseImageTag string

	SkipShutdown bool
}

// KvtoolRunner implements a NodeRunner that spins up local chains with kvtool.
// It has support for the following:
// - running a Nemo node
// - optionally, running an IBC node with a channel opened to the Nemo node
// - optionally, start the Nemo node on one version and upgrade to another
type KvtoolRunner struct {
	config KvtoolRunnerConfig
}

var _ NodeRunner = &KvtoolRunner{}

// NewKvtoolRunner creates a new KvtoolRunner.
func NewKvtoolRunner(config KvtoolRunnerConfig) *KvtoolRunner {
	return &KvtoolRunner{
		config: config,
	}
}

// StartChains implements NodeRunner.
// For KvtoolRunner, it sets up, runs, and connects to a local chain via kvtool.
func (k *KvtoolRunner) StartChains() Chains {
	// install kvtool if not already installed
	installKvtoolCmd := exec.Command("./scripts/install-kvtool.sh")
	installKvtoolCmd.Stdout = os.Stdout
	installKvtoolCmd.Stderr = os.Stderr
	if err := installKvtoolCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to install kvtool: %s", err.Error()))
	}

	// start local test network with kvtool
	log.Println("starting nemo node")
	kvtoolArgs := []string{"testnet", "bootstrap", "--nemo.configTemplate", k.config.NemoConfigTemplate}
	// include an ibc chain if desired
	if k.config.IncludeIBC {
		kvtoolArgs = append(kvtoolArgs, "--ibc")
	}
	// handle automated upgrade functionality, if defined
	if k.config.EnableAutomatedUpgrade {
		kvtoolArgs = append(kvtoolArgs,
			"--upgrade-name", k.config.NemoUpgradeName,
			"--upgrade-height", fmt.Sprint(k.config.NemoUpgradeHeight),
			"--upgrade-base-image-tag", k.config.NemoUpgradeBaseImageTag,
		)
	}
	// start the chain
	startNemoCmd := exec.Command("kvtool", kvtoolArgs...)
	startNemoCmd.Env = os.Environ()
	startNemoCmd.Env = append(startNemoCmd.Env, fmt.Sprintf("NEMO_TAG=%s", k.config.ImageTag))
	startNemoCmd.Stdout = os.Stdout
	startNemoCmd.Stderr = os.Stderr
	log.Println(startNemoCmd.String())
	if err := startNemoCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to start nemo: %s", err.Error()))
	}

	// wait for chain to be live.
	// if an upgrade is defined, this waits for the upgrade to be completed.
	if err := waitForChainStart(kvtoolNemoChain); err != nil {
		k.Shutdown()
		panic(err)
	}
	log.Println("nemo is started!")

	chains := NewChains()
	chains.Register("nemo", &kvtoolNemoChain)
	if k.config.IncludeIBC {
		chains.Register("ibc", &kvtoolIbcChain)
	}
	return chains
}

// Shutdown implements NodeRunner.
// For KvtoolRunner, it shuts down the local kvtool network.
// To prevent shutting down the chain (eg. to preserve logs or examine post-test state)
// use the `SkipShutdown` option on the config.
func (k *KvtoolRunner) Shutdown() {
	if k.config.SkipShutdown {
		log.Printf("would shut down but SkipShutdown is true")
		return
	}
	log.Println("shutting down nemo node")
	shutdownNemoCmd := exec.Command("kvtool", "testnet", "down")
	shutdownNemoCmd.Stdout = os.Stdout
	shutdownNemoCmd.Stderr = os.Stderr
	if err := shutdownNemoCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to shutdown kvtool: %s", err.Error()))
	}
}
