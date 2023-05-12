package api

import (
	"context"
	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

func (server *Server) report(e echo.Context) error {
	response := reportResponse{
		Message: "",
		Data:    nil,
	}

	inputUrl := e.FormValue("url")
	guildId := e.FormValue("guild_id")
	guildName := e.FormValue("guild_name")
	reporterName := e.FormValue("reporter_name")
	reporterAvatar := e.FormValue("reporter_avatar")

	if !server.checkChannelExist(guildId) {
		response.Message = "unrecognized channel"
		apply, _ := server.store.GetApplianceByGuildId(e.Request().Context(), guildId)

		if apply.GuildID != "" {
			response.Data = apply
			return e.JSON(http.StatusForbidden, response)
		}

		return e.JSON(http.StatusForbidden, response)
	}

	_, err := url.ParseRequestURI(inputUrl)

	if err != nil {
		response.Message = "invalid url input"
		response.Data = err
		return e.JSON(http.StatusUnprocessableEntity, response)
	}

	blocklistByUrl, err := server.store.GetReportBlocklistByUrl(e.Request().Context(), inputUrl)

	if len(blocklistByUrl) >= 1 {
		response.Message = "this link had been reported"
		response.Data = blocklistByUrl[0]
		return e.JSON(http.StatusConflict, response)
	}

	blocklist, err := server.store.CreateReportBlocklist(e.Request().Context(), db.CreateReportBlocklistParams{
		HttpAddress:    inputUrl,
		GuildID:        guildId,
		GuildName:      guildName,
		ReporterName:   reporterName,
		ReporterAvatar: reporterAvatar,
	})

	if err != nil {
		response.Message = "some error occurred"
		response.Data = err
		return e.JSON(http.StatusInternalServerError, response)
	}

	response.Message = "report success"
	response.Data = blocklist
	return e.JSON(http.StatusCreated, response)
}

func (server *Server) checkChannelExist(guildId string) bool {
	channel, _ := server.store.GetChannelsByGuildId(context.Background(), guildId)
	return channel.GuildID != ""
}

func (server *Server) getReportNFTs(e echo.Context) error {
	blocklists, _ := server.store.ListReportBlocklists(e.Request().Context())
	return e.JSON(http.StatusOK, blocklists)
}

type reportResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (server *Server) apply(e echo.Context) error {
	response := reportResponse{
		Message: "",
		Data:    nil,
	}

	guildId := e.FormValue("guild_id")
	channelName := e.FormValue("channel_name")

	applied, err := server.store.GetApplianceByGuildId(e.Request().Context(), guildId)
	if applied.GuildID != "" {
		response.Message = "this link have applied before"
		response.Data = applied

		return e.JSON(http.StatusConflict, response)
	}

	appliance, err := server.store.CreateAppliance(e.Request().Context(), db.CreateApplianceParams{
		ChannelName: channelName,
		GuildID:     guildId,
	})

	if err != nil {
		response.Message = "Internal Server Error"
		response.Data = err

		return e.JSON(http.StatusInternalServerError, response)
	}

	response.Message = "apply success"
	response.Data = appliance

	return e.JSON(http.StatusCreated, response)
}
