//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput

package file

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/zclconf/go-cty/cty"
	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/decrypt"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	File   string `mapstructure:"file" required:"true"`
	Format string `mapstructure:"format" required:"true"`

	ctx interpolate.Context
}

type DataSource struct {
	config Config
}

type DatasourceOutput struct {
	DecryptedRaw []byte `mapstructure:"decrypted_raw"`
	Decrypted    string `mapstructure:"decrypted"`
}

func (d *DataSource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *DataSource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	var errs *packersdk.MultiError

	if len(d.config.File) == 0 {
		errs = packersdk.MultiErrorAppend(errs, fmt.Errorf("the `file` must be specified"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (d *DataSource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *DataSource) Execute() (cty.Value, error) {
	decrypted, err := decrypt.File(d.config.File, d.config.Format)

	if userErr, ok := err.(sops.UserError); ok {
		err = fmt.Errorf(userErr.UserError())
	}

	if err != nil {
		return cty.Value{}, err
	}

	output := DatasourceOutput{
		DecryptedRaw: decrypted,
		Decrypted:    string(decrypted),
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
