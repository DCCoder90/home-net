package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/DCCoder90/home-net/pulumi/pkg/config"
)

var validRestartPolicies = map[string]bool{
	"":               true,
	"no":             true,
	"always":         true,
	"on-failure":     true,
	"unless-stopped": true,
}

func main() {
	configDir := flag.String("config-dir", "../config", "path to the config directory")
	flag.Parse()

	var errs []string

	serviceFiles, _ := filepath.Glob(filepath.Join(*configDir, "services", "*.yaml"))
	for _, f := range serviceFiles {
		errs = append(errs, validateServiceFile(f)...)
	}

	stackFiles, _ := filepath.Glob(filepath.Join(*configDir, "stacks", "*.yaml"))
	for _, f := range stackFiles {
		errs = append(errs, validateStackFile(f)...)
	}

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}
		fmt.Fprintf(os.Stderr, "\n%d error(s) found\n", len(errs))
		os.Exit(1)
	}

	fmt.Println("All config files valid.")
}

func validateServiceFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("%s: cannot read file: %v", path, err)}
	}

	var raw map[string]config.ServiceConfig
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&raw); err != nil && err != io.EOF {
		return []string{fmt.Sprintf("%s: YAML error: %v", path, err)}
	}

	var errs []string
	for key, svc := range raw {
		errs = append(errs, validateService(fmt.Sprintf("%s[%s]", path, key), svc)...)
	}
	return errs
}

func validateStackFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("%s: cannot read file: %v", path, err)}
	}

	var raw map[string]config.StackConfig
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&raw); err != nil && err != io.EOF {
		return []string{fmt.Sprintf("%s: YAML error: %v", path, err)}
	}

	var errs []string
	for stackKey, stack := range raw {
		stackPrefix := fmt.Sprintf("%s[%s]", path, stackKey)
		if len(stack.Services) == 0 {
			errs = append(errs, fmt.Sprintf("%s: services map is required and must not be empty", stackPrefix))
		}
		for svcKey, svc := range stack.Services {
			errs = append(errs, validateService(fmt.Sprintf("%s.services[%s]", stackPrefix, svcKey), svc)...)
		}
	}
	return errs
}

func validateService(prefix string, svc config.ServiceConfig) []string {
	var errs []string

	if svc.ServiceName == "" {
		errs = append(errs, fmt.Sprintf("%s: service_name is required", prefix))
	}
	if svc.ImageName == "" {
		errs = append(errs, fmt.Sprintf("%s: image_name is required", prefix))
	}
	if !validRestartPolicies[svc.RestartPolicy] {
		errs = append(errs, fmt.Sprintf("%s: invalid restart_policy %q (valid: no, always, on-failure, unless-stopped)", prefix, svc.RestartPolicy))
	}
	if svc.NetworkMode != "" && svc.Network != nil && len(svc.Network.Networks) > 0 {
		errs = append(errs, fmt.Sprintf("%s: network_mode and network.networks are mutually exclusive", prefix))
	}
	if svc.DNS.Enabled && svc.DNS.DomainName == "" {
		errs = append(errs, fmt.Sprintf("%s: dns.domain_name is required when dns.enabled is true", prefix))
	}
	if svc.Auth != nil && svc.Auth.Proxy != nil && svc.Auth.Proxy.Enabled && !svc.DNS.Enabled {
		errs = append(errs, fmt.Sprintf("%s: dns.enabled must be true when auth.proxy.enabled is true", prefix))
	}

	return errs
}
