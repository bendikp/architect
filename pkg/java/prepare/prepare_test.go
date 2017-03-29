package prepare_test

import (
	"github.com/skatteetaten/architect/pkg/java/prepare"
	"testing"
	"fmt"
)

func TestPrepare(t *testing.T) {
	err := prepare.Prepare("", map[string]string{"rsc": "3711"}, "/home/m87950/Downloads/innsynamelding-core-feature_AOS_1280_Endre_navn_paa_properties-20170302.081857-4-Leveransepakke.zip")

	if err != nil {
		fmt.Println(err)
	}
}
