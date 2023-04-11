package powerflex

import (
	"fmt"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	types "github.com/dell/goscaleio/types/v1"
)

// getFirstSystem - finds available first system and returns it.
func getFirstSystem(rc *goscaleio.Client) (*goscaleio.System, error) {
	allSystems, err := rc.GetSystems()
	if err != nil {
		return nil, fmt.Errorf("Error in goscaleio GetSystems")
	}
	system, err := rc.FindSystem(allSystems[0].ID, "", "")
	if err != nil {
		return nil, fmt.Errorf("Error in goscaleio FindSystem")
	}
	return system, nil
}

// getNewProtectionDomainEx function to get Protection Domain
func getNewProtectionDomainEx(c *goscaleio.Client, pdID string, pdName string, href string) (*goscaleio.ProtectionDomain, error) {
	system, err := getFirstSystem(c)
	if err != nil {
		return nil, err
	}
	pdr := goscaleio.NewProtectionDomainEx(c, &scaleiotypes.ProtectionDomain{})
	if pdID != "" {
		protectionDomain, err := system.FindProtectionDomain(pdID, "", "")
		pdr.ProtectionDomain = protectionDomain
		if err != nil {
			return nil, err
		}
	} else {
		protectionDomain, err := system.FindProtectionDomain("", pdName, "")
		pdr.ProtectionDomain = protectionDomain
		if err != nil {
			return nil, err
		}
	}
	return pdr, nil
}

func getSdcType(c *goscaleio.Client, sdcID string) (*goscaleio.Sdc, error) {
	system, err := getFirstSystem(c)
	if err != nil {
		return nil, err
	}

	sdcs, err1 := system.GetSdc()
	if err1 != nil {
		return nil, err1
	}

	var sdcType *types.Sdc
	for sdc := range sdcs {
		if sdcs[sdc].ID == sdcID {
			sdcType = &sdcs[sdc]
		}
	}

	c4 := goscaleio.NewSdc(c, sdcType)
	return c4, nil
}

func getVolumeType(c *goscaleio.Client, volID string) (*goscaleio.Volume, error) {
	volumes, err := c.GetVolume("", volID, "", "", false)
	if err != nil {
		// resp.Diagnostics.AddError(
		// 	"Error getting volume",
		// 	"unexpected error: "+err.Error(),
		// )
		return nil, err
	}

	volume := volumes[0]
	volType := goscaleio.NewVolume(c)
	volType.Volume = volume
	return volType, nil
}
