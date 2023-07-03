package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type FutoolRunnerConfig struct {
	NemoConfigTemplate string

	ImageTag   string
	IncludeIBC bool

	EnableAutomatedUpgrade  bool
	NemoUpgradeName         string
	NemoUpgradeHeight       int64
	NemoUpgradeBaseImageTag string

	SkipShutdown bool
}

// FutoolRunner implements a NodeRunner that spins up local chains with futool.
// It has support for the following:
// - running a Nemo node
// - optionally, running an IBC node with a channel opened to the Nemo node
// - optionally, start the Nemo node on one version and upgrade to another
type FutoolRunner struct {
	config FutoolRunnerConfig
}

var _ NodeRunner = &FutoolRunner{}

// NewFutoolRunner creates a new FutoolRunner.
func NewFutoolRunner(config FutoolRunnerConfig) *FutoolRunner {
	return &FutoolRunner{
		config: config,
	}
}

// StartChains implements NodeRunner.
// For FutoolRunner, it sets up, runs, and connects to a local chain via futool.
func (k *FutoolRunner) StartChains() Chains {
	// install futool if not already installed
	installFutoolCmd := exec.Command("./scripts/install-futool.sh")
	installFutoolCmd.Stdout = os.Stdout
	installFutoolCmd.Stderr = os.Stderr
	if err := installFutoolCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to install futool: %s", err.Error()))
	}

	// start local test network with futool
	log.Println("starting nemo node")
	futoolArgs := []string{"testnet", "bootstrap", "--nemo.configTemplate", k.config.NemoConfigTemplate}
	// include an ibc chain if desired
	if k.config.IncludeIBC {
		futoolArgs = append(futoolArgs, "--ibc")
	}
	// handle automated upgrade functionality, if defined
	if k.config.EnableAutomatedUpgrade {
		futoolArgs = append(futoolArgs,
			"--upgrade-name", k.config.NemoUpgradeName,
			"--upgrade-height", fmt.Sprint(k.config.NemoUpgradeHeight),
			"--upgrade-base-image-tag", k.config.NemoUpgradeBaseImageTag,
		)
	}
	// start the chain
	startNemoCmd := exec.Command("futool", futoolArgs...)
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
	if err := waitForChainStart(futoolNemoChain); err != nil {
		k.Shutdown()
		panic(err)
	}
	log.Println("nemo is started!")

	chains := NewChains()
	chains.Register("nemo", &futoolNemoChain)
	if k.config.IncludeIBC {
		chains.Register("ibc", &futoolIbcChain)
	}
	return chains
}

// Shutdown implements NodeRunner.
// For FutoolRunner, it shuts down the local futool network.
// To prevent shutting down the chain (eg. to preserve logs or examine post-test state)
// use the `SkipShutdown` option on the config.
func (k *FutoolRunner) Shutdown() {
	if k.config.SkipShutdown {
		log.Printf("would shut down but SkipShutdown is true")
		return
	}
	log.Println("shutting down nemo node")
	shutdownNemoCmd := exec.Command("futool", "testnet", "down")
	shutdownNemoCmd.Stdout = os.Stdout
	shutdownNemoCmd.Stderr = os.Stderr
	if err := shutdownNemoCmd.Run(); err != nil {
		panic(fmt.Sprintf("failed to shutdown futool: %s", err.Error()))
	}
}
