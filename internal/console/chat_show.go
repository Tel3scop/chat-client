package console

import (
	"github.com/spf13/cobra"
)

func (c *Console) runShow(_ *cobra.Command, _ []string) {
	c.whereIAm()
	c.showMyChats()
}
