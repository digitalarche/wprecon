package cli

import (
	"fmt"
	"os"

	"github.com/blackcrw/wprecon/cli/cmd"
	. "github.com/blackcrw/wprecon/cli/config"
	"github.com/blackcrw/wprecon/internal/pkg/banner"
	"github.com/blackcrw/wprecon/pkg/gohttp"
	"github.com/blackcrw/wprecon/pkg/printer"
	"github.com/blackcrw/wprecon/tools/wordpress/extensions"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:    "wprecon",
	Short:  "Wordpress Recon",
	Long:   `wprecon (Wordpress Recon) is a tool for wordpress exploration!`,
	PreRun: cmd.RootOptionsPreRun,
	Run:    cmd.RootOptionsRun,
	PostRun: func(cmd *cobra.Command, args []string) {
		printer.Done("Total requests:", fmt.Sprint(InfosWprecon.TotalRequests))
	},
}

var fuzzer = &cobra.Command{
	Use:     "fuzzer",
	Aliases: []string{"fuzz"},
	Short:   "Fuzzing directories or logins.",
	Long:    "This subcommand is mainly focused on fuzzing directories or logins.",
	Run:     cmd.FuzzerOptionsRun,
	PostRun: func(cmd *cobra.Command, args []string) {
		printer.Done("Total requests:", fmt.Sprint(InfosWprecon.TotalRequests))
	},
}

// Execute ::
func Execute() {
	if err := root.Execute(); err != nil {
		os.Exit(0)
	}
}

func init() {
	cobra.OnInitialize(ibanner)

	root.PersistentFlags().StringP("url", "u", "", "Target URL (Ex: http(s)://example.com/). "+printer.Required)
	root.PersistentFlags().StringP("scripts", "", "", "Auxiliary scripts.")
	root.PersistentFlags().BoolP("random-agent", "", false, "Use randomly selected HTTP(S) User-Agent header value.")
	root.PersistentFlags().BoolP("tor", "", false, "Use Tor anonymity network")
	root.PersistentFlags().BoolP("disable-tls-checks", "", false, "Disables SSL/TLS certificate verification.")
	root.PersistentFlags().BoolP("verbose", "v", false, "Verbosity mode.")
	root.PersistentFlags().BoolP("force", "f", false, "Forces wprecon to not check if the target is running WordPress and forces other executions.")

	root.Flags().BoolP("aggressive-mode", "", false, "Activates the aggressive mode of wprecon.")
	root.Flags().BoolP("detection-waf", "", false, "I will try to detect if the target is using any WAF Wordpress.")

	root.MarkPersistentFlagRequired("url")

	root.SetHelpTemplate(banner.HelpMain)

	// fuzzer
	// fuzzer.Flags().StringP("wordlist", "w", "", "Wordlist")
	fuzzer.Flags().StringP("usernames", "U", "", "Set usernames attack passwords.")
	fuzzer.Flags().StringP("passwords", "P", "", "Set wordlist attack passwords.")
	fuzzer.Flags().BoolP("backup-file", "B", false, "Performs a fuzzing to try to find the backup file if it exists.")
	// fuzzer.Flags().StringP("attack-method", "M", "xml-rpc", "Avaliable: xml-rpc and wp-login")

	fuzzer.SetHelpTemplate(banner.HelpFuzzer)
	root.AddCommand(fuzzer)
}

func ibanner() {
	InfosWprecon.Target, _ = root.Flags().GetString("url")
	InfosWprecon.Force, _ = root.Flags().GetBool("force")
	InfosWprecon.Verbose, _ = root.Flags().GetBool("verbose")
	InfosWprecon.OtherInformationsBool["http.options.tor"], _ = root.Flags().GetBool("tor")
	InfosWprecon.OtherInformationsBool["http.options.randomuseragent"], _ = root.Flags().GetBool("random-agent")
	InfosWprecon.OtherInformationsBool["http.options.tlscertificateverify"], _ = root.Flags().GetBool("tlscertificateverify")
	InfosWprecon.OtherInformationsString["scripts.name"], _ = root.Flags().GetString("scripts")

	go func() {
		response := extensions.SimpleRequest(InfosWprecon.Target, "")

		InfosWprecon.OtherInformationsString["target.http.index.raw"] = response.Raw
	}()

	if isURL, err := gohttp.IsURL(InfosWprecon.Target); isURL {
		banner.SBanner(
			InfosWprecon.Target,
			InfosWprecon.OtherInformationsBool["http.options.tor"],
			InfosWprecon.Verbose)
	} else if err != nil {
		printer.Fatal(err)
	} else {
		banner.Banner()
		printer.Fatal(err)
	}
}
