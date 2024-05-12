package console

import (
	"github.com/spf13/cobra"
)

func (c *Console) runShow(cmd *cobra.Command, args []string) {
	c.whereIAm()
	c.showMyChats()
}
