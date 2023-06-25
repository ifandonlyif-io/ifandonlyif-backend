package api

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
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

func (server *Server) appliance(e echo.Context) error {
	response := reportResponse{
		Message: "",
		Data:    nil,
	}

	appliances, _ := server.store.GetAllAppliances(e.Request().Context())

	response.Data = appliances
	response.Message = "get appliances"

	return e.JSON(http.StatusOK, response)
}

type approveRequest struct {
	IsApprove bool `json:"is_approved"`
}

func (server *Server) approve(e echo.Context) error {
	response := reportResponse{
		Message: "",
		Data:    nil,
	}

	req := &approveRequest{}
	applianceId := e.Param("id")

	if err := e.Bind(req); err != nil {
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	id, _ := uuid.Parse(applianceId)

	appliance, _ := server.store.GetApplianceChannelById(e.Request().Context(), id)

	if !appliance.VerifiedAt.Valid {
		updateAppliance, err := server.store.UpdateApplianceChannel(e.Request().Context(), db.UpdateApplianceChannelParams{
			ID: id,
			IsApproved: sql.NullBool{
				Bool:  req.IsApprove,
				Valid: true,
			},
		})

		if err != nil {
			response.Message = err.Error()
			return e.JSON(http.StatusInternalServerError, response)
		}

		if req.IsApprove {
			channel, err := server.store.CreateChannel(e.Request().Context(), db.CreateChannelParams{
				Name:    appliance.ChannelName,
				GuildID: appliance.GuildID,
			})
			if err != nil {
				response.Message = err.Error()
				return e.JSON(http.StatusInternalServerError, response)
			}
			response.Message = "appliance approved"
			response.Data = struct {
				Channel   interface{} `json:"channel"`
				Appliance interface{} `json:"appliance"`
			}{
				Channel:   channel,
				Appliance: updateAppliance,
			}

		} else {
			response.Message = "appliance not approved"
			response.Data = struct {
				Appliance interface{} `json:"appliance"`
			}{
				Appliance: updateAppliance,
			}

		}
	} else {
		response.Message = "this appliance have verified before"
		response.Data = struct {
			Appliance interface{} `json:"appliance"`
		}{
			Appliance: appliance,
		}
	}

	return e.JSON(200, response)
}

func (server *Server) channels(e echo.Context) error {
	response := reportResponse{
		Message: "",
		Data:    nil,
	}

	channels, _ := server.store.GetAllChannels(e.Request().Context())

	response.Data = channels
	response.Message = "get channels"

	return e.JSON(200, response)
}
