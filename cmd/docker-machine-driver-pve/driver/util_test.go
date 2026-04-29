package driver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getMACFromPveNetworkDevice(t *testing.T) {
	tests := map[string]string{
		"":                                       "",
		",":                                      "",
		"bridge=vmbr1":                           "",
		"e1000=BC:24:11:45:CD:E8,bridge=vmbr0":   "BC:24:11:45:CD:E8",
		"e1000e=BC:24:11:E1:F0:71,bridge=vmbr3":  "BC:24:11:E1:F0:71",
		"rtl8139=BC:24:11:18:BB:08,bridge=vmbr1": "BC:24:11:18:BB:08",
		"virtio=BC:24:11:87:63:EC,bridge=vmbr1":  "BC:24:11:87:63:EC",
		"vmxnet3=BC:24:11:EB:05:E9,bridge=vmbr4": "BC:24:11:EB:05:E9",
	}

	for deviceConfiguration, expectedAddress := range tests {
		t.Run(
			deviceConfiguration,
			func(t *testing.T) {
				require.Equal(
					t,
					expectedAddress,
					getMACFromPveNetworkDevice(deviceConfiguration),
				)
			},
		)
	}
}

func Test_isProxmoxCloudInitDrive(t *testing.T) {
	tests := map[string]bool{
		"":                                       false,
		"scsi1: none,media=cdrom":                false,
		"scsi1: local-lvm:cloudinit,media=cdrom": true,
		"ide2: local-lvm:CloudInit,media=cdrom":  true,
		"sata0: local-lvm:cloudinit,size=4G,media=cdrom": true,
	}

	for deviceConfiguration, expected := range tests {
		t.Run(deviceConfiguration, func(t *testing.T) {
			require.Equal(t, expected, isProxmoxCloudInitDrive(deviceConfiguration))
		})
	}
}

func Test_parsePveNetworkDevice(t *testing.T) {
	tests := map[string]struct {
		expectedName          string
		expectedConfiguration string
		expectedError         string
	}{
		"": {
			expectedError: "network device must not be empty",
		},
		"net=virtio,bridge=vmbr0": {
			expectedError: "network device name 'net' must match 'net<index>'",
		},
		"foo=virtio,bridge=vmbr0": {
			expectedError: "network device name 'foo' must match 'net<index>'",
		},
		"net1": {
			expectedError: "network device must be in '<device>=<configuration>' format",
		},
		"net1=": {
			expectedError: "network device 'net1' configuration must not be empty",
		},
		"net1=virtio,bridge=vmbr0": {
			expectedName:          "net1",
			expectedConfiguration: "virtio,bridge=vmbr0",
		},
		" NET2 = virtio=BC:24:11:87:63:EC,bridge=vmbr1 ": {
			expectedName:          "net2",
			expectedConfiguration: "virtio=BC:24:11:87:63:EC,bridge=vmbr1",
		},
	}

	for input, test := range tests {
		t.Run(input, func(t *testing.T) {
			name, configuration, err := parsePveNetworkDevice(input)

			if test.expectedError != "" {
				require.EqualError(t, err, test.expectedError)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedName, name)
			require.Equal(t, test.expectedConfiguration, configuration)
		})
	}
}
