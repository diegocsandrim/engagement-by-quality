package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

func main() {
	repos := []string{
		"AdguardTeam/AdguardFilters",
		"Alluxio/alluxio",
		"ampproject/amphtml",
		"angular/angular-cli",
		"angular/components",
		"ansible/ansible",
		"ant-design/ant-design",
		"apache/beam",
		"apache/incubator-mxnet",
		"apache/kafka",
		"apache/spark",
		"apple/swift",
		"ARMmbed/mbed-os",
		"Automattic/wp-calypso",
		"bitcoin/bitcoin",
		"brave/browser-laptop",
		"brianchandotcom/liferay-portal",
		"ceph/ceph",
		"CleverRaven/Cataclysm-DDA",
		"cms-sw/cms-sw.github.io",
		"cms-sw/cmssw",
		"cockroachdb/cockroach",
		"CocoaPods/Specs",
		"code-dot-org/code-dot-org",
		"DefinitelyTyped/DefinitelyTyped",
		"dimagi/commcare-hq",
		"dotnet/cli",
		"dotnet/coreclr",
		"dotnet/corefx",
		"dotnet/roslyn",
		"eclipse/che",
		"edx/edx-platform",
		"elastic/elasticsearch",
		"elastic/kibana",
		"electron/electron",
		"facebook/react",
		"fastlane/fastlane",
		"firehol/blocklist-ipsets",
		"flutter/flutter",
		"freeCodeCamp/freeCodeCamp",
		"gatsbyjs/gatsby",
		"gentoo/gentoo",
		"godotengine/godot",
		"golang/go",
		"grpc/grpc",
		"hacks-guide/Guide_3DS",
		"hashicorp/terraform",
		"helm/charts",
		"home-assistant/home-assistant",
		"Homebrew/homebrew-cask",
		"homebrew/homebrew-core",
		"ionic-team/ionic",
		"istio/istio",
		"jlippold/tweakCompatible",
		"jlord/patchwork",
		"joomla/joomla-cms",
		"JuliaLang/julia",
		"keybase/client",
		"kubernetes/kubernetes",
		"kubernetes/test-infra",
		"kubernetes/website",
		"laravel/framework",
		"magento/magento2",
		"ManageIQ/manageiq",
		"mapsme/omim",
		"MarlinFirmware/Marlin",
		"Microsoft/TypeScript",
		"Microsoft/vscode",
		"MicrosoftDocs/azure-docs",
		"moby/moby",
		"mui-org/material-ui",
		"nextcloud/server",
		"NixOS/nixpkgs",
		"nodejs/node",
		"odoo/odoo",
		"openshift/openshift-ansible",
		"openshift/origin",
		"oppia/oppia",
		"owncloud/core",
		"PaddlePaddle/Paddle",
		"pandas-dev/pandas",
		"phalcon/docs",
		"pingcap/tidb",
		"python/cpython",
		"pytorch/pytorch",
		"rails/rails",
		"rust-lang/crates.io-index",
		"rust-lang/rust",
		"saltstack/salt",
		"scikit-learn/scikit-learn",
		"servo/servo",
		"symfony/symfony",
		"tensorflow/tensorflow",
		"tgstation/tgstation",
		"vgstation-coders/vgstation13",
		"web-platform-tests/wpt",
		"WordPress/gutenberg",
		"XX-net/XX-Net",
		"zephyrproject-rtos/zephyr",
		"zulip/zulip",
	}

	for _, repo := range repos {
		// time.Sleep(1 * time.Second)
		url := fmt.Sprintf("https://api.github.com/repos/%s/languages", repo)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "token 3ba1abaf0bf06fb0cf2986d134fda22e95f879ae")
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		if resp.StatusCode != http.StatusOK {
			fmt.Println(bodyStr)
			panic(resp.StatusCode)
		}
		// fmt.Println(bodyStr)'
		too := math.Min(10, float64(len(bodyStr)))
		if strings.Contains(bodyStr[0:int64(too)], "\"Go\"") {
			fmt.Println(repo)
		}
	}
}
