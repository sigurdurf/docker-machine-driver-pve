package driver

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var pveNetworkDeviceNamePattern = regexp.MustCompile(`^net[0-9]+$`)

func getMACFromPveNetworkDevice(device string) string {
	models := []string{
		"e1000",
		"e1000-82540em",
		"e1000-82544gc",
		"e1000-82545em",
		"e1000e",
		"i82551",
		"i82557b",
		"i82559er",
		"ne2k_isa",
		"ne2k_pci",
		"pcnet",
		"rtl8139",
		"virtio",
		"vmxnet3",
	}

	for _, param := range strings.Split(device, ",") {
		//nolint:mnd
		values := strings.SplitN(param, "=", 2)

		//nolint:mnd
		if len(values) != 2 {
			continue
		}

		if slices.Contains(models, values[0]) {
			return values[1]
		}
	}

	return ""
}

func isProxmoxCloudInitDrive(device string) bool {
	return strings.Contains(strings.ToLower(device), "cloudinit")
}

func parsePveNetworkDevice(device string) (string, string, error) {
	trimmedDevice := strings.TrimSpace(device)
	if trimmedDevice == "" {
		return "", "", fmt.Errorf("network device must not be empty")
	}

	values := strings.SplitN(trimmedDevice, "=", 2)
	if len(values) != 2 {
		return "", "", fmt.Errorf("network device must be in '<device>=<configuration>' format")
	}

	deviceName := strings.ToLower(strings.TrimSpace(values[0]))
	if !pveNetworkDeviceNamePattern.MatchString(deviceName) {
		return "", "", fmt.Errorf("network device name '%s' must match 'net<index>'", values[0])
	}

	deviceConfiguration := strings.TrimSpace(values[1])
	if deviceConfiguration == "" {
		return "", "", fmt.Errorf("network device '%s' configuration must not be empty", deviceName)
	}

	return deviceName, deviceConfiguration, nil
}
