package update

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/origins"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

type Fields struct {
	OriginKey            string
	ApplicationID        int64
	Name                 string
	OriginType           string
	Addresses            []string
	OriginProtocolPolicy string
	HostHeader           string
	OriginPath           string
	HmacAuthentication   string
	HmacRegionName       string
	HmacAccessKey        string
	HmacSecretKey        string
	Path                 string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.OriginsUpdateUsage,
		Short:         msg.OriginsUpdateShortDescription,
		Long:          msg.OriginsUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion origins update --application-id 1673635839 --origin-key "58755fef-e830-4ea4-b9e0-6481f1ef496d" --name "ffcafe222sdsdffdf" --addresses "httpbin.org" --host-header "asdf.safe" --origin-type "single_origin" --origin-protocol-policy "http" --origin-path "/requests" --hmac-authentication "false"
        $ azion origins update --application-id 1673635839 --origin-key "58755fef-e830-4ea4-b9e0-6481f1ef496d" --name "drink coffe" --addresses "asdfg.asd" --host-header "host"
        $ azion origins update --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.UpdateOriginsRequest{}
			if cmd.Flags().Changed("in") {
				var (
					file *os.File
					err  error
				)
				if fields.Path == "-" {
					file = os.Stdin
				} else {
					file, err = os.Open(fields.Path)
					if err != nil {
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.Path)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("origin-key") {
					return msg.ErrorMandatoryUpdateFlags
				}
				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}
				if cmd.Flags().Changed("addresses") {
					request.SetAddresses(prepareAddresses(fields.Addresses))
				}
				if cmd.Flags().Changed("host-header") {
					request.SetHostHeader(fields.HostHeader)
				}
				if cmd.Flags().Changed("origin-type") {
					request.SetOriginType(fields.OriginType)
				}
				if cmd.Flags().Changed("origin-protocol-policy") {
					request.SetOriginProtocolPolicy(fields.OriginProtocolPolicy)
				}
				if cmd.Flags().Changed("origin-path") {
					request.SetOriginPath(fields.OriginPath)
				}
				if cmd.Flags().Changed("hmac-authentication") {
					hmacAuth, err := strconv.ParseBool(fields.HmacAuthentication)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorHmacAuthenticationFlag, fields.HmacAuthentication)
					}
					request.SetHmacAuthentication(hmacAuth)
				}
				if cmd.Flags().Changed("hmac-region-name") {
					request.SetHmacRegionName(fields.HmacRegionName)
				}
				if cmd.Flags().Changed("hmac-access-key") {
					request.SetHmacAccessKey(fields.HmacAccessKey)
				}
				if cmd.Flags().Changed("hmac-secret-key") {
					request.SetHmacSecretKey(fields.HmacSecretKey)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.UpdateOrigins(context.Background(), fields.ApplicationID, fields.OriginKey, &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateOrigin.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.OriginsUpdateOutputSuccess, response.GetOriginKey())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&fields.OriginKey, "origin-key", "o", "", msg.OriginsUpdateFlagOriginKey)
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.OriginsCreateFlagEdgeApplicationId)
	flags.StringVar(&fields.Name, "name", "", msg.OriginsCreateFlagName)
	flags.StringVar(&fields.OriginType, "origin-type", "", msg.OriginsCreateFlagOriginType)
	flags.StringSliceVar(&fields.Addresses, "addresses", []string{}, msg.OriginsCreateFlagAddresses)
	flags.StringVar(&fields.OriginProtocolPolicy, "origin-protocol-policy", "", msg.OriginsCreateFlagOriginProtocolPolicy)
	flags.StringVar(&fields.HostHeader, "host-header", "", msg.OriginsCreateFlagHostHeader)
	flags.StringVar(&fields.OriginPath, "origin-path", "", msg.OriginsCreateFlagOriginPath)
	flags.StringVar(&fields.HmacAuthentication, "hmac-authentication", "", msg.OriginsCreateFlagHmacAuthentication)
	flags.StringVar(&fields.HmacRegionName, "hmac-region-name", "", msg.OriginsCreateFlagHmacRegionName)
	flags.StringVar(&fields.HmacAccessKey, "hmac-access-key", "", msg.OriginsCreateFlagHmacAccessKey)
	flags.StringVar(&fields.HmacSecretKey, "hmac-secret-key", "", msg.OriginsCreateFlagHmacSecretKey)
	flags.StringVar(&fields.Path, "in", "", msg.OriginsCreateFlagIn)
	flags.BoolP("help", "h", false, msg.OriginsCreateHelpFlag)
	return cmd
}

func prepareAddresses(addrs []string) (addresses []sdk.CreateOriginsRequestAddresses) {
	var addr sdk.CreateOriginsRequestAddresses
	for _, v := range addrs {
		addr.Address = v
		addresses = append(addresses, addr)
	}
	return
}
