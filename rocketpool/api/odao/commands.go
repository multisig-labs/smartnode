package odao

import (
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/utils/api"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

// Register subcommands
func RegisterSubcommands(command *cli.Command, name string, aliases []string) {
	command.Subcommands = append(command.Subcommands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the Rocket Pool oracle DAO",
		Subcommands: []cli.Command{

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get oracle DAO status",
				UsageText: "rocketpool api odao status",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(getStatus(c))
					return nil

				},
			},

			{
				Name:      "members",
				Aliases:   []string{"m"},
				Usage:     "Get the oracle DAO members",
				UsageText: "rocketpool api odao members",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(getMembers(c))
					return nil

				},
			},

			{
				Name:      "proposals",
				Aliases:   []string{"p"},
				Usage:     "Get the oracle DAO proposals",
				UsageText: "rocketpool api odao proposals",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(getProposals(c))
					return nil

				},
			},

			{
				Name:      "proposal-details",
				Aliases:   []string{"d"},
				Usage:     "Get details of a proposal",
				UsageText: "rocketpool api odao proposal-details proposal-id",
				Action: func(c *cli.Context) error {

					// Validate args
					var err error
					if err = cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					id, err := cliutils.ValidateUint("proposal-id", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(getProposal(c, id))
					return nil

				},
			},

			{
				Name:      "can-propose-invite",
				Usage:     "Check whether the node can propose inviting a new member",
				UsageText: "rocketpool api odao can-propose-invite member-address member-id member-url",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 3); err != nil {
						return err
					}
					memberAddress, err := cliutils.ValidateAddress("member address", c.Args().Get(0))
					if err != nil {
						return err
					}
					memberId, err := cliutils.ValidateDAOMemberID("member ID", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeInvite(c, memberAddress, memberId, c.Args().Get(2)))
					return nil

				},
			},
			{
				Name:      "propose-invite",
				Aliases:   []string{"i"},
				Usage:     "Propose inviting a new member",
				UsageText: "rocketpool api odao propose-invite member-address member-id member-url",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 3); err != nil {
						return err
					}
					memberAddress, err := cliutils.ValidateAddress("member address", c.Args().Get(0))
					if err != nil {
						return err
					}
					memberId, err := cliutils.ValidateDAOMemberID("member ID", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeInvite(c, memberAddress, memberId, c.Args().Get(2)))
					return nil

				},
			},

			{
				Name:      "can-propose-leave",
				Usage:     "Check whether the node can propose leaving the oracle DAO",
				UsageText: "rocketpool api odao can-propose-leave",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeLeave(c))
					return nil

				},
			},
			{
				Name:      "propose-leave",
				Aliases:   []string{"l"},
				Usage:     "Propose leaving the oracle DAO",
				UsageText: "rocketpool api odao propose-leave",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeLeave(c))
					return nil

				},
			},

			{
				Name:      "can-propose-kick",
				Usage:     "Check whether the node can propose kicking a member",
				UsageText: "rocketpool api odao can-propose-kick member-address fine-amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					memberAddress, err := cliutils.ValidateAddress("member address", c.Args().Get(0))
					if err != nil {
						return err
					}
					fineAmountWei, err := cliutils.ValidatePositiveOrZeroWeiAmount("fine amount", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeKick(c, memberAddress, fineAmountWei))
					return nil

				},
			},
			{
				Name:      "propose-kick",
				Aliases:   []string{"k"},
				Usage:     "Propose kicking a member",
				UsageText: "rocketpool api odao propose-kick member-address fine-amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					memberAddress, err := cliutils.ValidateAddress("member address", c.Args().Get(0))
					if err != nil {
						return err
					}
					fineAmountWei, err := cliutils.ValidatePositiveOrZeroWeiAmount("fine amount", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeKick(c, memberAddress, fineAmountWei))
					return nil

				},
			},

			{
				Name:      "can-cancel-proposal",
				Usage:     "Check whether the node can cancel a proposal",
				UsageText: "rocketpool api odao can-cancel-proposal proposal-id",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalId, err := cliutils.ValidatePositiveUint("proposal ID", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canCancelProposal(c, proposalId))
					return nil

				},
			},
			{
				Name:      "cancel-proposal",
				Aliases:   []string{"c"},
				Usage:     "Cancel a proposal made by the node",
				UsageText: "rocketpool api odao cancel-proposal proposal-id",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalId, err := cliutils.ValidatePositiveUint("proposal ID", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(cancelProposal(c, proposalId))
					return nil

				},
			},

			{
				Name:      "can-vote-proposal",
				Usage:     "Check whether the node can vote on a proposal",
				UsageText: "rocketpool api odao can-vote-proposal proposal-id",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalId, err := cliutils.ValidatePositiveUint("proposal ID", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canVoteOnProposal(c, proposalId))
					return nil

				},
			},
			{
				Name:      "vote-proposal",
				Aliases:   []string{"v"},
				Usage:     "Vote on a proposal",
				UsageText: "rocketpool api odao vote-proposal proposal-id support",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					proposalId, err := cliutils.ValidatePositiveUint("proposal ID", c.Args().Get(0))
					if err != nil {
						return err
					}
					support, err := cliutils.ValidateBool("support", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(voteOnProposal(c, proposalId, support))
					return nil

				},
			},

			{
				Name:      "can-execute-proposal",
				Usage:     "Check whether the node can execute a proposal",
				UsageText: "rocketpool api odao can-execute-proposal proposal-id",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalId, err := cliutils.ValidatePositiveUint("proposal ID", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canExecuteProposal(c, proposalId))
					return nil

				},
			},
			{
				Name:      "execute-proposal",
				Aliases:   []string{"x"},
				Usage:     "Execute a proposal",
				UsageText: "rocketpool api odao execute-proposal proposal-id",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalId, err := cliutils.ValidatePositiveUint("proposal ID", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(executeProposal(c, proposalId))
					return nil

				},
			},

			{
				Name:      "can-join",
				Usage:     "Check whether the node can join the oracle DAO",
				UsageText: "rocketpool api odao can-join",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(canJoin(c))
					return nil

				},
			},
			{
				Name:      "join-approve-ggp",
				Aliases:   []string{"j1"},
				Usage:     "Approves the GGP bond transfer prior to join the oracle DAO",
				UsageText: "rocketpool api odao join-approve-ggp",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(approveGgp(c))
					return nil

				},
			},
			{
				Name:      "join",
				Aliases:   []string{"j2"},
				Usage:     "Join the oracle DAO (requires an executed invite proposal)",
				UsageText: "rocketpool api odao join tx-hash",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					hash, err := cliutils.ValidateTxHash("tx-hash", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(waitForApprovalAndJoin(c, hash))
					return nil

				},
			},

			{
				Name:      "can-leave",
				Usage:     "Check whether the node can leave the oracle DAO",
				UsageText: "rocketpool api odao can-leave",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(canLeave(c))
					return nil

				},
			},
			{
				Name:      "leave",
				Aliases:   []string{"e"},
				Usage:     "Leave the oracle DAO (requires an executed leave proposal)",
				UsageText: "rocketpool api odao leave bond-refund-address",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					bondRefundAddress, err := cliutils.ValidateAddress("bond refund address", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(leave(c, bondRefundAddress))
					return nil

				},
			},

			{
				Name:      "can-propose-members-quorum",
				Usage:     "Check whether the node can propose the members.quorum setting",
				UsageText: "rocketpool api odao can-propose-members-quorum value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					quorum, err := cliutils.ValidateFraction("quorum", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingMembersQuorum(c, quorum))
					return nil

				},
			},
			{
				Name:      "propose-members-quorum",
				Usage:     "Propose updating the members.quorum setting",
				UsageText: "rocketpool api odao propose-members-quorum value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					quorum, err := cliutils.ValidateFraction("quorum", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingMembersQuorum(c, quorum))
					return nil

				},
			},

			{
				Name:      "can-propose-members-ggpbond",
				Usage:     "Check whether the node can propose the members.ggpbond setting",
				UsageText: "rocketpool api odao can-propose-members-ggpbond value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					bondAmountWei, err := cliutils.ValidateWeiAmount("GGP bond amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingMembersGgpBond(c, bondAmountWei))
					return nil

				},
			},
			{
				Name:      "propose-members-ggpbond",
				Usage:     "Propose updating the members.ggpbond setting",
				UsageText: "rocketpool api odao propose-members-ggpbond value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					bondAmountWei, err := cliutils.ValidateWeiAmount("GGP bond amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingMembersGgpBond(c, bondAmountWei))
					return nil

				},
			},

			{
				Name:      "can-propose-members-minipool-unbonded-max",
				Usage:     "Check whether the node can propose the members.minipool.unbonded.max setting",
				UsageText: "rocketpool api odao can-propose-members-minipool-unbonded-max value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					unbondedMinipoolMax, err := cliutils.ValidateUint("maximum unbonded minipool count", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingMinipoolUnbondedMax(c, unbondedMinipoolMax))
					return nil

				},
			},
			{
				Name:      "propose-members-minipool-unbonded-max",
				Usage:     "Propose updating the members.minipool.unbonded.max setting",
				UsageText: "rocketpool api odao propose-members-minipool-unbonded-max value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					unbondedMinipoolMax, err := cliutils.ValidateUint("maximum unbonded minipool count", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingMinipoolUnbondedMax(c, unbondedMinipoolMax))
					return nil

				},
			},

			{
				Name:      "can-propose-proposal-cooldown",
				Usage:     "Check whether the node can propose the proposal.cooldown setting",
				UsageText: "rocketpool api odao can-propose-proposal-cooldown value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalCooldownBlocks, err := cliutils.ValidateUint("proposal cooldown period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingProposalCooldown(c, proposalCooldownBlocks))
					return nil

				},
			},
			{
				Name:      "propose-proposal-cooldown",
				Usage:     "Propose updating the proposal.cooldown setting",
				UsageText: "rocketpool api odao propose-proposal-cooldown value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalCooldownBlocks, err := cliutils.ValidateUint("proposal cooldown period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingProposalCooldown(c, proposalCooldownBlocks))
					return nil

				},
			},

			{
				Name:      "can-propose-proposal-vote-timespan",
				Usage:     "Check whether the node can propose the proposal.vote.time setting",
				UsageText: "rocketpool api odao can-propose-proposal-vote-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalVoteTimespan, err := cliutils.ValidateUint("proposal voting period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingProposalVoteTimespan(c, proposalVoteTimespan))
					return nil

				},
			},
			{
				Name:      "propose-proposal-vote-timespan",
				Usage:     "Propose updating the proposal.vote.time setting",
				UsageText: "rocketpool api odao propose-proposal-vote-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalVoteTimespan, err := cliutils.ValidateUint("proposal voting period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingProposalVoteTimespan(c, proposalVoteTimespan))
					return nil

				},
			},

			{
				Name:      "can-propose-proposal-vote-delay-timespan",
				Usage:     "Check whether the node can propose the proposal.vote.delay.time setting",
				UsageText: "rocketpool api odao can-propose-proposal-vote-delay-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalDelayTimespan, err := cliutils.ValidateUint("proposal delay period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingProposalVoteDelayTimespan(c, proposalDelayTimespan))
					return nil

				},
			},
			{
				Name:      "propose-proposal-vote-delay-timespan",
				Usage:     "Propose updating the proposal.vote.delay.time setting",
				UsageText: "rocketpool api odao propose-proposal-vote-delay-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalDelayTimespan, err := cliutils.ValidateUint("proposal delay period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingProposalVoteDelayTimespan(c, proposalDelayTimespan))
					return nil

				},
			},

			{
				Name:      "can-propose-proposal-execute-timespan",
				Usage:     "Check whether the node can propose the proposal.execute.time setting",
				UsageText: "rocketpool api odao can-propose-proposal-execute-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalExecuteTimespan, err := cliutils.ValidateUint("proposal execution period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingProposalExecuteTimespan(c, proposalExecuteTimespan))
					return nil

				},
			},
			{
				Name:      "propose-proposal-execute-timespan",
				Usage:     "Propose updating the proposal.execute.time setting",
				UsageText: "rocketpool api odao propose-proposal-execute-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalExecuteTimespan, err := cliutils.ValidateUint("proposal execution period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingProposalExecuteTimespan(c, proposalExecuteTimespan))
					return nil

				},
			},

			{
				Name:      "can-propose-proposal-action-timespan",
				Usage:     "Check whether the node can propose the proposal.action.time setting",
				UsageText: "rocketpool api odao can-propose-proposal-action-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalActionTimespan, err := cliutils.ValidateUint("proposal action period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingProposalActionTimespan(c, proposalActionTimespan))
					return nil

				},
			},
			{
				Name:      "propose-proposal-action-timespan",
				Usage:     "Propose updating the proposal.action.time setting",
				UsageText: "rocketpool api odao propose-proposal-action-timespan value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					proposalActionTimespan, err := cliutils.ValidateUint("proposal action period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingProposalActionTimespan(c, proposalActionTimespan))
					return nil

				},
			},

			{
				Name:      "can-propose-scrub-period",
				Usage:     "Check whether the node can propose the minipool.scrub.period setting",
				UsageText: "rocketpool api odao can-propose-scrub-period value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					scrubPeriod, err := cliutils.ValidateUint("scrub period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canProposeSettingScrubPeriod(c, scrubPeriod))
					return nil

				},
			},
			{
				Name:      "propose-scrub-period",
				Usage:     "Propose updating the minipool.scrub.period setting",
				UsageText: "rocketpool api odao propose-scrub-period value",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					scrubPeriod, err := cliutils.ValidateUint("scrub period", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(proposeSettingScrubPeriod(c, scrubPeriod))
					return nil

				},
			},

			{
				Name:      "get-member-settings",
				Usage:     "Get the ODAO settings related to ODAO members",
				UsageText: "rocketpool api odao get-member-settings",
				Action: func(c *cli.Context) error {

					// Run
					api.PrintResponse(getMemberSettings(c))
					return nil

				},
			},
			{
				Name:      "get-proposal-settings",
				Usage:     "Get the ODAO settings related to ODAO proposals",
				UsageText: "rocketpool api odao get-proposal-settings",
				Action: func(c *cli.Context) error {

					// Run
					api.PrintResponse(getProposalSettings(c))
					return nil

				},
			},
			{
				Name:      "get-minipool-settings",
				Usage:     "Get the ODAO settings related to minipools",
				UsageText: "rocketpool api odao get-minipool-settings",
				Action: func(c *cli.Context) error {

					// Run
					api.PrintResponse(getMinipoolSettings(c))
					return nil

				},
			},
		},
	})
}
