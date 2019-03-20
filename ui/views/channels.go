package views

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
)

const (
	CHANNELS        = "channels"
	CHANNELS_HEADER = "header"
)

type Channels struct {
	*gocui.View
	items   []*models.Channel
	network *network.Network
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	headerView, err := g.SetView(CHANNELS_HEADER, x0, y0, x1, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	headerView.Frame = false
	headerView.BgColor = gocui.ColorGreen
	headerView.FgColor = gocui.ColorBlack | gocui.AttrBold
	displayChannelsHeader(headerView)

	c.View, err = g.SetView(CHANNELS, x0, y0+1, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.View.Frame = false

	err = c.update(context.Background())
	if err != nil {
		return err
	}

	c.display()
	return nil
}

func displayChannelsHeader(v *gocui.View) {
	fmt.Fprintln(v, fmt.Sprintf("%-9s  %-19s  %12s  %12s  %s",
		"status",
		"id",
		"local",
		"capacity",
		"pub_key",
	))
}

func (c *Channels) Refresh(g *gocui.Gui) error {
	var err error
	c.View, err = g.View(CHANNELS)
	if err != nil {
		return err
	}

	err = c.update(context.Background())
	if err != nil {
		return err
	}

	c.display()
	return nil
}

func (c *Channels) update(ctx context.Context) error {
	channels, err := c.network.ListChannels(ctx)
	if err != nil {
		return err
	}

	c.items = channels
	return nil
}

func (c *Channels) display() {
	for _, item := range c.items {
		line := fmt.Sprintf("%s  %s  %s  %12d  %s",
			active(item),
			chartID(item),
			color.Cyan(fmt.Sprintf("%12d", item.LocalBalance)),
			item.Capacity,
			item.RemotePubKey,
		)
		fmt.Fprintln(c.View, line)
	}
}

func NewChannels(network *network.Network) *Channels {
	return &Channels{network: network}
}

func active(c *models.Channel) string {
	if c.Active {
		return color.Green(fmt.Sprintf("%-9s", "active"))
	}
	return color.Red(fmt.Sprintf("%-9s", "inactive"))
}

func chartID(c *models.Channel) string {
	id := fmt.Sprintf("%-19d", c.ID)
	index := int(c.LocalBalance * int64(len(id)) / c.Capacity)

	var buffer bytes.Buffer
	buffer.WriteString(color.Cyan(id[:index]))
	buffer.WriteString(id[index:])

	return buffer.String()
}
