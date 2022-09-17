package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	_ "github.com/ifandonlyif-io/ifandonlyif-backend/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

type GasPrices struct {
	Average int32 `json:"average"`
}

func (server *Server) RunCronFetchGas() {
	// Create a Resty Client
	client := resty.New()

	cronjob := cron.New()

	// cronjob.AddFunc("@hourly", func() {
	cronjob.AddFunc("*/1 * * * *", func() {
		resp, err := client.R().
			EnableTrace().
			Get("https://ethgasstation.info/api/ethgasAPI.json")

		if err != nil {
			fmt.Println("No response from request")
		}

		var gp GasPrices

		if err := json.Unmarshal(resp.Body(), &gp); err != nil { // Parse []byte to the go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		createGas, dberr := server.store.CreateGasPrice(context.Background(), sql.NullInt32{Int32: int32(gp.Average), Valid: true})
		if dberr != nil {
			return
		}
		fmt.Println("cron times : ", createGas)

	})

	cronjob.Start()
	fmt.Println("Fetch Gas Cron Job Started !!!!")

}

// gasInfo godoc
// @Summary      gasInfo
// @Description  get 24 hours gas prices
// @Tags         gasInfo
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /gasInfo [GET]
func (server *Server) GasHandler(c echo.Context) (err error) {
	getGasInfo, err := server.store.GetAveragePriceByLastDay(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusAccepted, getGasInfo)
}
