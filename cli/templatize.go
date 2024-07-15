package cli

import (
	"context"

	"github.com/spf13/cobra"
	"go.keploy.io/server/v2/config"
	replaySvc "go.keploy.io/server/v2/pkg/service/replay"
	"go.keploy.io/server/v2/utils"
	"go.uber.org/zap"
)

func init() {
	Register("Templatize", Templatize)
}

func Templatize(ctx context.Context, logger *zap.Logger, conf *config.Config, serviceFactory ServiceFactory, cmdConfigurator CmdConfigurator) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "templatize",
		Short:   "templatize the keploy testcases for re-record",
		Example: `keploy templatize -c "/path/to/user/app"`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return cmdConfigurator.Validate(ctx, cmd)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Read the testcases from the path provided.
				// utils.ReadTempValues(testSet)
				// Get the replay service.
				svc, err := serviceFactory.GetService(ctx, "normalize")
				if err != nil {
					utils.LogError(logger, err, "failed to get service")
					return nil
				}
				var replay replaySvc.Service
				var ok bool
				if replay, ok = svc.(replaySvc.Service); !ok {
					utils.LogError(logger, nil, "service doesn't satisfy replay service interface")
					return nil
				}
				if err := replay.Templatize(ctx, []string{}); err != nil {
					utils.LogError(logger, err, "failed to templatize test cases")
					return nil
				}
			return nil
		},
	}

	err := cmdConfigurator.AddFlags(cmd)
	if err != nil {
		utils.LogError(logger, err, "failed to add templatize flags")
		return nil
	}

	return cmd
}
