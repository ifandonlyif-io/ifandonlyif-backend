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

// blocklists godoc
// @Summary      Report from discord bot
// @Description  report api for discord bot
// @Tags         Discord Bot
// @Accept */*
// @produce application/json
// @Success      201  {string}  StatusCreated
// @Router       /discord/report [POST]
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

	channel, _ := server.store.GetChannelsByGuildId(context.Background(), guildId)

	if channel.GuildID == "" {
		response.Message = "unrecognized channel"
		apply, _ := server.store.GetApplianceByGuildId(e.Request().Context(), guildId)

		if apply.GuildID != "" {
			response.Data = apply
			return e.JSON(http.StatusForbidden, response)
		}

		return e.JSON(http.StatusForbidden, response)
	}

	if channel.LockedAt.Valid {
		response.Message = "this channel have been locked"
		response.Data = channel

		return e.JSON(http.StatusLocked, response)
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

// blocklists godoc
// @Summary      Apply from discord bot
// @Description  apply api for discord bot
// @Tags         Discord Bot
// @Accept */*
// @produce application/json
// @Success      201  {string}  StatusCreated
// @Router       /discord/apply [POST]
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

// blocklists godoc
// @Summary      Appliances from discord bot
// @Description  get all appliances, for admin page
// @Tags         Admin Page
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /discord/appliances [GET]
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

// blocklists godoc
// @Summary      Approve discord channel
// @Description  approve discord channel to report project
// @Tags         Admin Page
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /discord/approve/:id [PATCH]
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

// blocklists godoc
// @Summary      Get all discord channels
// @Description  get all discord channels can report projects
// @Tags         Admin Page
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /discord/channels [GET]
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

// blocklists godoc
// @Summary      Lock discord channel
// @Description  lock discord channel to prevent report project
// @Tags         Admin Page
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /discord/channel/lock/:id [PATCH]
func (server *Server) lockChannel(e echo.Context) error {
	response := reportResponse{
		Message: "",
		Data:    nil,
	}

	id, _ := uuid.Parse(e.Param("id"))

	channel, err := server.store.LockDiscordChannel(e.Request().Context(), id)
	if err != nil {
		response.Message = err.Error()
		response.Data = err
		return e.JSON(http.StatusInternalServerError, response)
	}

	response.Message = "channel locked"
	response.Data = channel

	return e.JSON(http.StatusOK, response)
}

// blocklists godoc
// @Summary      Unlock discord channel
// @Description  unlock discord channel to allow report project
// @Tags         Admin Page
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /discord/channel/unlock/:id [PATCH]
func (server *Server) UnlockChannel(e echo.Context) error {
	response := reportResponse{
		Message: "",
		Data:    nil,
	}

	id, _ := uuid.Parse(e.Param("id"))

	channel, err := server.store.UnlockDiscordChannel(e.Request().Context(), id)
	if err != nil {
		response.Message = err.Error()
		response.Data = err
		return e.JSON(http.StatusInternalServerError, response)
	}

	response.Message = "channel locked"
	response.Data = channel

	return e.JSON(http.StatusOK, response)

}
