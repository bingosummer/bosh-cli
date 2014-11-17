package cpi

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshcmd "github.com/cloudfoundry/bosh-agent/platform/commands"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bmcloud "github.com/cloudfoundry/bosh-micro-cli/cloud"
	bmcomp "github.com/cloudfoundry/bosh-micro-cli/cpi/compile"
	bmcpiinstall "github.com/cloudfoundry/bosh-micro-cli/cpi/install"
	bmdepl "github.com/cloudfoundry/bosh-micro-cli/deployment"
	bmrel "github.com/cloudfoundry/bosh-micro-cli/release"
	bmrelvalidation "github.com/cloudfoundry/bosh-micro-cli/release/validation"
	bmui "github.com/cloudfoundry/bosh-micro-cli/ui"
)

type Installer interface {
	Extract(releaseTarballPath string) (bmrel.Release, error)
	Install(deployment bmdepl.Deployment, release bmrel.Release) (bmcloud.Cloud, error)
}

type cpiInstaller struct {
	ui              bmui.UI
	fs              boshsys.FileSystem
	extractor       boshcmd.Compressor
	validator       bmrelvalidation.ReleaseValidator
	releaseCompiler bmcomp.ReleaseCompiler
	jobInstaller    bmcpiinstall.JobInstaller
	cloudFactory    bmcloud.Factory
	logger          boshlog.Logger
	logTag          string
}

func NewInstaller(
	ui bmui.UI,
	fs boshsys.FileSystem,
	extractor boshcmd.Compressor,
	validator bmrelvalidation.ReleaseValidator,
	releaseCompiler bmcomp.ReleaseCompiler,
	jobInstaller bmcpiinstall.JobInstaller,
	cloudFactory bmcloud.Factory,
	logger boshlog.Logger,
) Installer {
	return &cpiInstaller{
		ui:              ui,
		fs:              fs,
		extractor:       extractor,
		validator:       validator,
		releaseCompiler: releaseCompiler,
		jobInstaller:    jobInstaller,
		cloudFactory:    cloudFactory,
		logger:          logger,
		logTag:          "cpiInstaller",
	}
}

// Extract decompresses a release tarball into a temp directory (release.extractedPath),
// parses and validates the release manifest, and decompresses the packages and jobs.
// Use release.Delete() to clean up the temp directory.
func (c *cpiInstaller) Extract(releaseTarballPath string) (bmrel.Release, error) {
	c.logger.Info(c.logTag, "Extracting CPI release")
	extractedReleasePath, err := c.fs.TempDir("cmd-deployCmd")
	if err != nil {
		c.ui.Error("Could not create a temporary directory")
		return nil, bosherr.WrapError(err, "Creating temp directory")
	}

	releaseReader := bmrel.NewReader(releaseTarballPath, extractedReleasePath, c.fs, c.extractor)
	release, err := releaseReader.Read()
	if err != nil {
		c.ui.Error(fmt.Sprintf("CPI release at `%s' is not a BOSH release", releaseTarballPath))
		return nil, bosherr.WrapError(err, fmt.Sprintf("Reading CPI release from `%s'", releaseTarballPath))
	}

	c.logger.Info(c.logTag, "Extracted CPI release `%s' to `%s'", release.Name, extractedReleasePath)

	c.logger.Info(c.logTag, "Validating CPI release `%s'", release.Name)
	err = c.validator.Validate(release)
	if err != nil {
		return nil, bosherr.WrapError(err, "Validating CPI release `%s'", release.Name)
	}

	return release, nil
}

func (c *cpiInstaller) Install(deployment bmdepl.Deployment, release bmrel.Release) (bmcloud.Cloud, error) {
	c.logger.Info(c.logTag, fmt.Sprintf("Compiling CPI release `%s'", release.Name))
	c.logger.Debug(c.logTag, fmt.Sprintf("Compiling CPI release `%s': %#v", release.Name, release))
	err := c.releaseCompiler.Compile(release, deployment)
	if err != nil {
		c.ui.Error("Could not compile CPI release")
		return nil, bosherr.WrapError(err, "Compiling CPI release")
	}

	jobs := deployment.Jobs
	if len(jobs) != 1 {
		c.ui.Error("Invalid CPI deployment: exactly one job required")
		return nil, bosherr.New("Invalid CPI deployment: exactly one job required, %d jobs found", len(jobs))
	}
	cpiJob := jobs[0]

	instances := cpiJob.Instances
	if instances != 1 {
		c.ui.Error("Invalid CPI deployment: exactly one job instance required")
		return nil, bosherr.New(
			"Invalid CPI deployment: exactly one instance required, found %d instances in job `%s'",
			instances,
			cpiJob.Name,
		)
	}

	installedJobs, err := c.installJob(cpiJob, release)
	if err != nil {
		c.ui.Error("Could not install CPI deployment job")
		return nil, bosherr.WrapError(err, "Installing CPI deployment job")
	}

	cloud, err := c.cloudFactory.NewCloud(installedJobs)
	if err != nil {
		c.ui.Error("Invalid CPI deployment")
		return nil, bosherr.WrapError(err, "Validating CPI deployment job installation")
	}

	return cloud, nil
}

func (c *cpiInstaller) installJob(deploymentJob bmdepl.Job, release bmrel.Release) ([]bmcpiinstall.InstalledJob, error) {
	installedJobs := make([]bmcpiinstall.InstalledJob, 0, len(deploymentJob.Templates))
	for _, releaseJobRef := range deploymentJob.Templates {
		releaseJobName := releaseJobRef.Name
		releaseJob, found := release.FindJobByName(releaseJobName)

		if !found {
			c.ui.Error(fmt.Sprintf("Could not find CPI job `%s' in release `%s'", releaseJobName, release.Name))
			return installedJobs, bosherr.New("Invalid CPI deployment manifest: job `%s' not found in release `%s'", releaseJobName, release.Name)
		}

		installedJob, err := c.jobInstaller.Install(releaseJob)
		if err != nil {
			c.ui.Error(fmt.Sprintf("Could not install `%s' job", releaseJobName))
			return installedJobs, bosherr.WrapError(err, "Installing `%s' job for CPI release", releaseJobName)
		}
		installedJobs = append(installedJobs, installedJob)
	}
	return installedJobs, nil
}