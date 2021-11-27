package sql

import (
	"fmt"
	"strings"
	"testing"
)

func TestUrlLength(t *testing.T) {
	start := "/api/getBulkPrices?symbol="
	end := "dogeusdt_lunausdt_fttusdt_xlmusdt_vetusdt_icpusdt_thetausdt_algousdt_xmrusdt_xtzusdt_egldusdt_axsusdt_iotausdt_ftmusdt_hbarusdt_neousdt_wavesusdt_mkrusdt_nearusdt_bttusdt_chzusdt_stxusdt_dcrusdt_xemusdt_omgusdt_zecusdt_sushiusdt_enjusdt_manausdt_yfiusdt_iostusdt_qtumusdt_batusdt_zilusdt_icxusdt_grtusdt_celousdt_zenusdt_renusdt_scusdt_zrxusdt_ontusdt_nanousdt_crvusdt_bntusdt_fetusdt_umausdt_iotxusdt_lrcusdt_sandusdt_srmusdt_kavausdt_kncusdt"
	fmt.Println(" ", len([]byte(start)), " ", len([]byte(end)))
	arrs := strings.Split(end, "_")
	fmt.Println(len(arrs), " ", len([]byte("_")))
	fmt.Println(443 / 4)
	fmt.Println("", len(arrs), " ", arrs[0]+"usdt")
}
