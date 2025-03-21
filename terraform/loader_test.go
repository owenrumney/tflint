package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"github.com/zclconf/go-cty/cty"
)

func TestLoadConfig(t *testing.T) {
	withinFixtureDir(t, "v0.15.0_module", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		config, diags := loader.LoadConfig(".", CallAllModule)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		// root
		if config.Module.SourceDir != "." {
			t.Fatalf("root module path: want=%s, got=%s", ".", config.Module.SourceDir)
		}
		// module.instance
		testChildModule(t, config, "instance", "ec2")
		// module.consul
		testChildModule(t, config, "consul", ".terraform/modules/consul")
		// module.consul.module.consul_clients
		testChildModule(
			t,
			config.Children["consul"],
			"consul_clients",
			".terraform/modules/consul/modules/consul-cluster",
		)
		// module.consul.module.consul_clients.module.iam_policies
		testChildModule(
			t,
			config.Children["consul"].Children["consul_clients"],
			"iam_policies",
			".terraform/modules/consul/modules/consul-iam-policies",
		)
		// module.consul.module.consul_clients.module.security_group_rules
		testChildModule(
			t,
			config.Children["consul"].Children["consul_clients"],
			"security_group_rules",
			".terraform/modules/consul/modules/consul-security-group-rules",
		)
		// module.consul.module.consul_servers
		testChildModule(
			t,
			config.Children["consul"],
			"consul_servers",
			".terraform/modules/consul/modules/consul-cluster",
		)
		// module.consul.module.consul_servers.module.iam_policies
		testChildModule(
			t,
			config.Children["consul"].Children["consul_servers"],
			"iam_policies",
			".terraform/modules/consul/modules/consul-iam-policies",
		)
		// module.consul.module.consul_servers.module.security_group_rules
		testChildModule(
			t,
			config.Children["consul"].Children["consul_servers"],
			"security_group_rules",
			".terraform/modules/consul/modules/consul-security-group-rules",
		)
	})
}

func TestLoadConfig_withBaseDir(t *testing.T) {
	withinFixtureDir(t, "v0.15.0_module", func(dir string) {
		// The current dir is test-fixtures/v0.15.0_module, but the base dir is test-fixtures
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, filepath.Dir(dir))
		if err != nil {
			t.Fatal(err)
		}
		config, diags := loader.LoadConfig(".", CallAllModule)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		// root
		if config.Module.SourceDir != "." {
			t.Fatalf("root module path: want=%s, got=%s", ".", config.Module.SourceDir)
		}
		// module.instance
		testChildModule(t, config, "instance", "ec2")
		// module.consul
		testChildModule(t, config, "consul", ".terraform/modules/consul")
		// module.consul.module.consul_clients
		testChildModule(
			t,
			config.Children["consul"],
			"consul_clients",
			".terraform/modules/consul/modules/consul-cluster",
		)
		// module.consul.module.consul_clients.module.iam_policies
		testChildModule(
			t,
			config.Children["consul"].Children["consul_clients"],
			"iam_policies",
			".terraform/modules/consul/modules/consul-iam-policies",
		)
		// module.consul.module.consul_clients.module.security_group_rules
		testChildModule(
			t,
			config.Children["consul"].Children["consul_clients"],
			"security_group_rules",
			".terraform/modules/consul/modules/consul-security-group-rules",
		)
		// module.consul.module.consul_servers
		testChildModule(
			t,
			config.Children["consul"],
			"consul_servers",
			".terraform/modules/consul/modules/consul-cluster",
		)
		// module.consul.module.consul_servers.module.iam_policies
		testChildModule(
			t,
			config.Children["consul"].Children["consul_servers"],
			"iam_policies",
			".terraform/modules/consul/modules/consul-iam-policies",
		)
		// module.consul.module.consul_servers.module.security_group_rules
		testChildModule(
			t,
			config.Children["consul"].Children["consul_servers"],
			"security_group_rules",
			".terraform/modules/consul/modules/consul-security-group-rules",
		)
	})
}

func TestLoadConfig_callLocalModules(t *testing.T) {
	withinFixtureDir(t, "v0.15.0_module", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		config, diags := loader.LoadConfig(".", CallLocalModule)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		// root
		if config.Module.SourceDir != "." {
			t.Fatalf("root module path: want=%s, got=%s", ".", config.Module.SourceDir)
		}
		// module.instance
		testChildModule(t, config, "instance", "ec2")

		if len(config.Children) != 1 {
			t.Fatalf("Root module has children unexpectedly: %#v", config.Children)
		}
	})
}

func TestLoadConfig_withoutModuleManifest(t *testing.T) {
	withinFixtureDir(t, "without_module_manifest", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		_, diags := loader.LoadConfig(".", CallAllModule)
		if !diags.HasErrors() {
			t.Fatal("Expected error is not occurred")
		}

		expected := `module.tf:6,1-16: "consul" module is not found. Did you run "terraform init"?; `
		if diags.Error() != expected {
			t.Fatalf(`Expected error is "%s", but got "%s"`, expected, diags)
		}
	})
}

func TestLoadConfig_withoutModuleManifest_callLocalModules(t *testing.T) {
	withinFixtureDir(t, "without_module_manifest", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		config, diags := loader.LoadConfig(".", CallLocalModule)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		// root
		if config.Module.SourceDir != "." {
			t.Fatalf("root module path: want=%s, got=%s", ".", config.Module.SourceDir)
		}
		// module.instance
		testChildModule(t, config, "instance", "ec2")

		if len(config.Children) != 1 {
			t.Fatalf("Root module has children unexpectedly: %#v", config.Children)
		}
	})
}

func TestLoadConfig_moduleNotFound(t *testing.T) {
	withinFixtureDir(t, "module_not_found", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		_, diags := loader.LoadConfig(".", CallLocalModule)
		if !diags.HasErrors() {
			t.Fatal("Expected error is not occurred")
		}

		expected := `module.tf:1,1-22: "ec2_instance" module is not found; The module directory "tf_aws_ec2_instance" does not exist or cannot be read.`
		if diags.Error() != expected {
			t.Fatalf(`Expected error is "%s", but got "%s"`, expected, diags)
		}
	})
}

func TestLoadConfig_moduleNotFound_callNoModules(t *testing.T) {
	withinFixtureDir(t, "module_not_found", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		config, diags := loader.LoadConfig(".", CallNoModule)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		if config.Module.SourceDir != "." {
			t.Fatalf("Root module path: want=%s, got=%s", ".", config.Module.SourceDir)
		}
		if len(config.Children) != 0 {
			t.Fatalf("Root module has children unexpectedly: %#v", config.Children)
		}
	})
}

func TestLoadConfig_moduleNotFound_callNoModules_withArgDir(t *testing.T) {
	withinFixtureDir(t, ".", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		config, diags := loader.LoadConfig("module_not_found", CallNoModule)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		if config.Module.SourceDir != "module_not_found" {
			t.Fatalf("Root module path: want=%s, got=%s", "module_not_found", config.Module.SourceDir)
		}
		if len(config.Children) != 0 {
			t.Fatalf("Root module has children unexpectedly: %#v", config.Children)
		}
	})
}

func TestLoadConfig_invalidConfiguration(t *testing.T) {
	withinFixtureDir(t, "invalid_configuration", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		_, diags := loader.LoadConfig(".", CallNoModule)
		if !diags.HasErrors() {
			t.Fatal("Expected error is not occurred")
		}

		expected := "resource.tf:3,23-29: Missing newline after argument; An argument definition must end with a newline."
		if diags.Error() != expected {
			t.Fatalf(`Expected error is "%s", but got "%s"`, expected, diags)
		}
	})
}

func TestLoadConfig_circularReferencingModules(t *testing.T) {
	withinFixtureDir(t, "circular_referencing_modules", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		_, diags := loader.LoadConfig(".", CallAllModule)
		if !diags.HasErrors() {
			t.Fatal("Expected error is not occurred")
		}

		file := filepath.Join("module2", "main.tf")
		expected := fmt.Sprintf(`%s:1,1-17: Module stack level too deep; This configuration has nested modules more than 10 levels deep. This is mainly caused by circular references. current path: module.module1.module.module2.module.module1.module.module2.module.module1.module.module2.module.module1.module.module2.module.module1.module.module2`, file)
		if diags.Error() != expected {
			t.Fatalf(`Expected error is "%s", but got "%s"`, expected, diags)
		}
	})
}

func TestLoadValuesFiles(t *testing.T) {
	withinFixtureDir(t, "values_files", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		ret, diags := loader.LoadValuesFiles(".", "cli1.tfvars", "cli2.tfvars")
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		expected := []InputValues{
			{
				"default": {
					Value: cty.StringVal("terraform.tfvars"),
				},
			},
			{
				"auto1": {
					Value: cty.StringVal("auto1.auto.tfvars"),
				},
			},
			{
				"auto2": {
					Value: cty.StringVal("auto2.auto.tfvars"),
				},
			},
			{
				"cli1": {
					Value: cty.StringVal("cli1.tfvars"),
				},
			},
			{
				"cli2": {
					Value: cty.StringVal("cli2.tfvars"),
				},
			},
		}

		if !reflect.DeepEqual(expected, ret) {
			t.Fatalf("Unexpected input values are received: expected=%#v actual=%#v", expected, ret)
		}

		want := []string{
			"auto1.auto.tfvars",
			"auto2.auto.tfvars",
			"cli1.tfvars",
			"cli2.tfvars",
			"terraform.tfvars",
		}
		loadedFiles := []string{}
		for name := range loader.Files() {
			loadedFiles = append(loadedFiles, name)
		}
		opt := cmpopts.SortSlices(func(x, y string) bool { return x > y })
		if diff := cmp.Diff(want, loadedFiles, opt); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestLoadValuesFiles_withBaseDir(t *testing.T) {
	withinFixtureDir(t, "values_files", func(dir string) {
		// The current dir is test-fixtures/values_files, but the base dir is test-fixtures
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, filepath.Dir(dir))
		if err != nil {
			t.Fatal(err)
		}
		// Files passed manually are relative to the current directory.
		ret, diags := loader.LoadValuesFiles(
			".",
			"cli1.tfvars",
			"cli2.tfvars",
		)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		expected := []InputValues{
			{
				"default": {
					Value: cty.StringVal("terraform.tfvars"),
				},
			},
			{
				"auto1": {
					Value: cty.StringVal("auto1.auto.tfvars"),
				},
			},
			{
				"auto2": {
					Value: cty.StringVal("auto2.auto.tfvars"),
				},
			},
			{
				"cli1": {
					Value: cty.StringVal("cli1.tfvars"),
				},
			},
			{
				"cli2": {
					Value: cty.StringVal("cli2.tfvars"),
				},
			},
		}

		if !reflect.DeepEqual(expected, ret) {
			t.Fatalf("Unexpected input values are received: expected=%#v actual=%#v", expected, ret)
		}

		want := []string{
			filepath.Join("values_files", "auto1.auto.tfvars"),
			filepath.Join("values_files", "auto2.auto.tfvars"),
			filepath.Join("values_files", "cli1.tfvars"),
			filepath.Join("values_files", "cli2.tfvars"),
			filepath.Join("values_files", "terraform.tfvars"),
		}
		loadedFiles := []string{}
		for name := range loader.Files() {
			loadedFiles = append(loadedFiles, name)
		}
		opt := cmpopts.SortSlices(func(x, y string) bool { return x > y })
		if diff := cmp.Diff(want, loadedFiles, opt); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestLoadValuesFiles_withArgDir(t *testing.T) {
	withinFixtureDir(t, ".", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		ret, diags := loader.LoadValuesFiles(
			"values_files",
			filepath.Join("values_files", "cli1.tfvars"),
			filepath.Join("values_files", "cli2.tfvars"),
		)
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		expected := []InputValues{
			{
				"default": {
					Value: cty.StringVal("terraform.tfvars"),
				},
			},
			{
				"auto1": {
					Value: cty.StringVal("auto1.auto.tfvars"),
				},
			},
			{
				"auto2": {
					Value: cty.StringVal("auto2.auto.tfvars"),
				},
			},
			{
				"cli1": {
					Value: cty.StringVal("cli1.tfvars"),
				},
			},
			{
				"cli2": {
					Value: cty.StringVal("cli2.tfvars"),
				},
			},
		}

		if !reflect.DeepEqual(expected, ret) {
			t.Fatalf("Unexpected input values are received: expected=%#v actual=%#v", expected, ret)
		}

		want := []string{
			filepath.Join("values_files", "auto1.auto.tfvars"),
			filepath.Join("values_files", "auto2.auto.tfvars"),
			filepath.Join("values_files", "cli1.tfvars"),
			filepath.Join("values_files", "cli2.tfvars"),
			filepath.Join("values_files", "terraform.tfvars"),
		}
		loadedFiles := []string{}
		for name := range loader.Files() {
			loadedFiles = append(loadedFiles, name)
		}
		opt := cmpopts.SortSlices(func(x, y string) bool { return x > y })
		if diff := cmp.Diff(want, loadedFiles, opt); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestLoadValuesFiles_invalidValuesFile(t *testing.T) {
	withinFixtureDir(t, "invalid_values_files", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		_, diags := loader.LoadValuesFiles(".")
		if !diags.HasErrors() {
			t.Fatal("Expected error is not occurred")
		}

		expected := "terraform.tfvars:3,1-9: Unexpected \"resource\" block; Blocks are not allowed here."
		if diags.Error() != expected {
			t.Fatalf(`Expected error is "%s", but got "%s"`, expected, diags)
		}
	})
}

func TestLoadConfigDirFiles_loader(t *testing.T) {
	withinFixtureDir(t, "v0.15.0_module", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		files, diags := loader.LoadConfigDirFiles(".")
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		want := []string{"module.tf"}
		loadedFiles := []string{}
		for name := range files {
			loadedFiles = append(loadedFiles, name)
		}
		opt := cmpopts.SortSlices(func(x, y string) bool { return x > y })
		if diff := cmp.Diff(want, loadedFiles, opt); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestLoadConfigDirFiles_loader_withBaseDir(t *testing.T) {
	withinFixtureDir(t, "v0.15.0_module", func(dir string) {
		// The current dir is test-fixtures/v0.15.0_module, but the base dir is test-fixtures
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, filepath.Dir(dir))
		if err != nil {
			t.Fatal(err)
		}
		files, diags := loader.LoadConfigDirFiles(".")
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		want := []string{filepath.Join("v0.15.0_module", "module.tf")}
		loadedFiles := []string{}
		for name := range files {
			loadedFiles = append(loadedFiles, name)
		}
		opt := cmpopts.SortSlices(func(x, y string) bool { return x > y })
		if diff := cmp.Diff(want, loadedFiles, opt); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestLoadConfigDirFiles_loader_withArgDir(t *testing.T) {
	withinFixtureDir(t, ".", func(dir string) {
		loader, err := NewLoader(afero.Afero{Fs: afero.NewOsFs()}, dir)
		if err != nil {
			t.Fatal(err)
		}
		files, diags := loader.LoadConfigDirFiles("v0.15.0_module")
		if diags.HasErrors() {
			t.Fatal(diags)
		}

		want := []string{filepath.Join("v0.15.0_module", "module.tf")}
		loadedFiles := []string{}
		for name := range files {
			loadedFiles = append(loadedFiles, name)
		}
		opt := cmpopts.SortSlices(func(x, y string) bool { return x > y })
		if diff := cmp.Diff(want, loadedFiles, opt); diff != "" {
			t.Fatal(diff)
		}
	})
}

func withinFixtureDir(t *testing.T, dir string, test func(string)) {
	t.Helper()

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	workingDir := filepath.Join(currentDir, "test-fixtures", dir)

	t.Chdir(workingDir)
	test(workingDir)
}

func testChildModule(t *testing.T, config *Config, key string, wantPath string) {
	t.Helper()

	if _, exists := config.Children[key]; !exists {
		t.Fatalf("`%s` module is not loaded, are submodules downloaded?: %#v", key, config.Children)
	}
	modulePath := config.Children[key].Module.SourceDir
	if modulePath != wantPath {
		t.Fatalf("`%s` module path: want=%s, got=%s", key, wantPath, modulePath)
	}
}
