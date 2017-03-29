package acceptance_test

import (
	"io/ioutil"
	"net"
	"os"

	"github.com/containernetworking/cni/pkg/ns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("errors", func() {
	BeforeEach(func() {
		cniEnv = map[string]string{
			"CNI_IFNAME":      "eth0",
			"CNI_CONTAINERID": "apricot",
			"CNI_PATH":        paths.CNIPath,
		}

		var err error
		containerNS, err = ns.NewNS()
		Expect(err).NotTo(HaveOccurred())

		cniEnv["CNI_NETNS"] = containerNS.Path()

		dataDir, err = ioutil.TempDir("", "cni-data-dir-")
		Expect(err).NotTo(HaveOccurred())

		flannelSubnetBaseIP, flannelSubnetCIDR, _ := net.ParseCIDR("10.255.30.0/24")
		_, fullNetwork, _ = net.ParseCIDR("10.255.0.0/16")
		flannelSubnet = &net.IPNet{
			IP:   flannelSubnetBaseIP,
			Mask: flannelSubnetCIDR.Mask,
		}
		subnetEnvFile = writeSubnetEnvFile(flannelSubnet.String(), fullNetwork.String())
		cniStdin = cniConfig(dataDir, subnetEnvFile)
	})

	AfterEach(func() {
		containerNS.Close() // don't bother checking errors here
		mustSucceed("iptables", "-t", "nat", "-F")
		Expect(os.RemoveAll(subnetEnvFile)).To(Succeed())
		Expect(os.RemoveAll(dataDir)).To(Succeed())
	})

	Describe("errors on ADD", func() {
		Context("when the subnet file is missing", func() {
			BeforeEach(func() {
				cniStdin = cniConfig(dataDir, "/path/does/not/exist")
			})

			It("exits with nonzero status and prints a CNI error result as JSON to stdout", func() {
				session := startCommand("ADD", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(1))

				Expect(session.Out.Contents()).To(MatchJSON(`{
				"code": 100,
				"msg": "discovering network info",
				"details": "open /path/does/not/exist: no such file or directory"
			}`))
			})
		})

		Context("when the subnet file is corrupt", func() {
			BeforeEach(func() {
				subnetEnvFile = writeSubnetEnvFile("bad-subnet", fullNetwork.String())
				cniStdin = cniConfig(dataDir, subnetEnvFile)
			})

			It("exits with nonzero status and prints a CNI error result as JSON to stdout", func() {
				session := startCommand("ADD", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(1))

				Expect(session.Out.Contents()).To(MatchJSON(`{
				"code": 100,
				"msg": "discovering network info",
				"details": "unable to parse flannel subnet file"
			}`))
			})
		})

		Context("when the ipam plugin errors on add", func() {
			BeforeEach(func() {
				subnetEnvFile = writeSubnetEnvFile("10.255.30.0/33", fullNetwork.String())
				cniStdin = cniConfig(dataDir, subnetEnvFile)
			})
			It("exits with nonzero status and prints a CNI error result as JSON to stdout", func() {
				session := startCommand("ADD", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(1))

				Expect(session.Out.Contents()).To(MatchJSON(`{
				"code": 100,
				"msg": "ipam plugin failed",
				"details": "invalid CIDR address: 10.255.30.0/33"
			}`))
			})
		})

		Context("when the veth manager fails to create a veth pair", func() {
			It("exits with nonzero status and prints a CNI error", func() {
				cniEnv["CNI_IFNAME"] = "some-bad-eth-name"
				cniStdin = cniConfig(dataDir, subnetEnvFile)
				session := startCommand("ADD", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(1))

				Expect(session.Out.Contents()).To(MatchJSON(`{
				"code": 100,
				"msg": "creating config",
				"details": "IfName cannot be longer than 15 characters"
			}`))
			})
		})
	})

	Describe("errors on DEL", func() {
		Context("when the ipam plugin errors on del", func() {
			BeforeEach(func() {
				subnetEnvFile = writeSubnetEnvFile("10.255.30.0/33", fullNetwork.String())
				cniStdin = cniConfig(dataDir, subnetEnvFile)
			})

			It("exits with zero status but logs the error", func() {
				session := startCommand("DEL", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(0))

				Expect(string(session.Err.Contents())).To(ContainSubstring(`invalid CIDR address: 10.255.30.0/33`))
			})
		})

		Context("when the network namespace doesn't exist", func() {
			BeforeEach(func() {
				cniEnv["CNI_NETNS"] = "/tmp/not/there"
			})
			It("exits with zero status but logs the error", func() {
				session := startCommand("DEL", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(0))

				Expect(session.Err).To(gbytes.Say(`opening-netns.*/tmp/not/there.*no such file or directory`))
			})
		})

		Context("when the interface isn't present inside the container", func() {
			It("exits with zero status, but logs the error", func() {
				cniEnv["CNI_IFNAME"] = "not-there"
				session := startCommand("DEL", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(0))
				Expect(string(session.Err.Contents())).To(ContainSubstring(`"deviceName":"not-there","message":"Link not found"`))
			})
		})

		Context("when the subnet file is missing", func() {
			BeforeEach(func() {
				cniStdin = cniConfig(dataDir, "/path/does/not/exist")
			})

			It("exits with nonzero status and prints a CNI error result as JSON to stdout", func() {
				session := startCommand("DEL", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(1))

				Expect(session.Out.Contents()).To(MatchJSON(`{
				"code": 100,
				"msg": "discovering network info",
				"details": "open /path/does/not/exist: no such file or directory"
			}`))
			})
		})

		Context("when the subnet file is corrupt", func() {
			BeforeEach(func() {
				subnetEnvFile = writeSubnetEnvFile("bad-subnet", fullNetwork.String())
				cniStdin = cniConfig(dataDir, subnetEnvFile)
			})

			It("exits with nonzero status and prints a CNI error result as JSON to stdout", func() {
				session := startCommand("DEL", cniStdin)
				Eventually(session, cmdTimeout).Should(gexec.Exit(1))

				Expect(session.Out.Contents()).To(MatchJSON(`{
				"code": 100,
				"msg": "discovering network info",
				"details": "unable to parse flannel subnet file"
			}`))
			})
		})
	})
})
