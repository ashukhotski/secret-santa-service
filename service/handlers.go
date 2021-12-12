// handlers.go
package service

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handlers struct {
	logger *log.Logger
	repo   ServiceRepo
}

func NewHandlers(l *log.Logger, r ServiceRepo) *Handlers {
	return &Handlers{
		logger: l,
		repo:   r,
	}
}

func (h *Handlers) SetupRoutes(mux *mux.Router) {
	loggingMiddleware := LoggingMiddleware(h.logger)
	mux.Use(loggingMiddleware)
	mux.HandleFunc("/get", h.GetHandler).Methods(http.MethodPost)
	mux.HandleFunc("/initialize", h.InitializeHandler).Methods(http.MethodPost)
	mux.HandleFunc("/participate", h.ParticipateHandler).Methods(http.MethodPost)
	mux.HandleFunc("/randomize", h.RandomizeHandler).Methods(http.MethodPost)
}

func (h *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	req := &SlackRequest{ChannelId: String(r.PostForm.Get("channel_id")),
		ChannelName:    String(r.PostForm.Get("channel_name")),
		Command:        String(r.PostForm.Get("command")),
		EnterpriseId:   nil,
		EnterpriseName: nil,
		ResponseUrl:    r.PostForm.Get("response_url"),
		TeamDomain:     String(r.PostForm.Get("team_domain")),
		TeamId:         String(r.PostForm.Get("team_id")),
		Text:           String(r.PostForm.Get("text")),
		Token:          String(r.PostForm.Get("token")),
		TriggerId:      String(r.PostForm.Get("trigger_id")),
		UserId:         r.PostForm.Get("user_id"),
		UserName:       r.PostForm.Get("user_name"),
	}

	t := time.Now()
	y := t.Year()

	p, err := h.repo.GetParticipantById(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, req.UserId, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}
	
	if !p.IsMatched || p.YourMatchId == nil {
		channelId := ""
		if req.ChannelId != nil {
			channelId = *req.ChannelId
		}
		err = errors.New("Secret Santa " + strconv.Itoa(y) + " pairs have not been matched yet for the Slack channel <#" + channelId + ">")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	// Slack
	yourMatchId := ""
	if p.YourMatchId != nil {
		yourMatchId = *p.YourMatchId
	}
	yourMatchAddress := ""
	if p.YourMatchAddress != nil {
		yourMatchAddress = *p.YourMatchAddress
	}
	msg := "Your match is <@" + yourMatchId + ">. Prepare your gift and send it to " + yourMatchAddress + ". Thank you and happy New Year!"
	err = SendSlackMessage(req.ResponseUrl, ResponseTypeEphemeral, msg)
	if err != nil {
		h.logger.Println(err)
	}

	//w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, msg})
}

func (h *Handlers) InitializeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	req := &SlackRequest{ChannelId: String(r.PostForm.Get("channel_id")),
		ChannelName:    String(r.PostForm.Get("channel_name")),
		Command:        String(r.PostForm.Get("command")),
		EnterpriseId:   nil,
		EnterpriseName: nil,
		ResponseUrl:    r.PostForm.Get("response_url"),
		TeamDomain:     String(r.PostForm.Get("team_domain")),
		TeamId:         String(r.PostForm.Get("team_id")),
		Text:           String(r.PostForm.Get("text")),
		Token:          String(r.PostForm.Get("token")),
		TriggerId:      String(r.PostForm.Get("trigger_id")),
		UserId:         r.PostForm.Get("user_id"),
		UserName:       r.PostForm.Get("user_name"),
	}
	
	if req.Text == nil || (req.Text != nil && len(*req.Text) < 5) {
		err = errors.New("Please provide a valid postal address by typing it after the command")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	t := time.Now()
	y := t.Year()

	pCount, _ := h.repo.CountAllParticipants(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, y)
	if pCount > 0 {
		err := errors.New("Secret Santa " + strconv.Itoa(y) + " has already been initialized for this Slack channel")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	p := &Participant{req.Text,
		req.ChannelId,
		req.EnterpriseId,
		true,
		false,
		req.ResponseUrl,
		req.TeamId,
		req.UserId,
		req.UserName,
		nil,
		nil,
		nil,
	}
	err = h.repo.RegisterParticipant(r.Context(), p, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	// Slack
	channelId := ""
	if req.ChannelId != nil {
		channelId = *req.ChannelId
	}
	msg := "<@" + req.UserId + "> just initiated Secret Santa " + strconv.Itoa(y) + " for the Slack channel <#" + channelId +">"
	err = SendSlackMessage(req.ResponseUrl, ResponseTypeInChannel, msg)
	if err != nil {
		h.logger.Println(err.Error())
	}

	//w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeInChannel, msg})
}

func (h *Handlers) ParticipateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	req := &SlackRequest{ChannelId: String(r.PostForm.Get("channel_id")),
		ChannelName:    String(r.PostForm.Get("channel_name")),
		Command:        String(r.PostForm.Get("command")),
		EnterpriseId:   nil,
		EnterpriseName: nil,
		ResponseUrl:    r.PostForm.Get("response_url"),
		TeamDomain:     String(r.PostForm.Get("team_domain")),
		TeamId:         String(r.PostForm.Get("team_id")),
		Text:           String(r.PostForm.Get("text")),
		Token:          String(r.PostForm.Get("token")),
		TriggerId:      String(r.PostForm.Get("trigger_id")),
		UserId:         r.PostForm.Get("user_id"),
		UserName:       r.PostForm.Get("user_name"),
	}
	
	if req.Text == nil || (req.Text != nil && len(*req.Text) < 5) {
		err = errors.New("Please provide a valid postal address by typing it after the command")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	t := time.Now()
	y := t.Year()

	pCount, err := h.repo.CountAllParticipants(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	if pCount < 1 {
		err := errors.New("Secret Santa has not been initialized yet")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	p := &Participant{req.Text, req.ChannelId, req.EnterpriseId, false, false, req.ResponseUrl, req.TeamId, req.UserId, req.UserName, nil, nil, nil}
	err = h.repo.RegisterParticipant(r.Context(), p, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	// Slack
	channelId := ""
	if p.ChannelId != nil {
		channelId = *p.ChannelId
	}
	msg := "<@" + p.UserId + "> just enrolled in Secret Santa " + strconv.Itoa(y) + " for the Slack channel <#" + channelId + ">"
	err = SendSlackMessage(p.ResponseUrl, ResponseTypeInChannel, msg)
	if err != nil {
		h.logger.Println(err.Error())
	}

	//w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeInChannel, msg})
}

func (h *Handlers) RandomizeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	req := &SlackRequest{ChannelId: String(r.PostForm.Get("channel_id")),
		ChannelName:    String(r.PostForm.Get("channel_name")),
		Command:        String(r.PostForm.Get("command")),
		EnterpriseId:   nil,
		EnterpriseName: nil,
		ResponseUrl:    r.PostForm.Get("response_url"),
		TeamDomain:     String(r.PostForm.Get("team_domain")),
		TeamId:         String(r.PostForm.Get("team_id")),
		Text:           String(r.PostForm.Get("text")),
		Token:          String(r.PostForm.Get("token")),
		TriggerId:      String(r.PostForm.Get("trigger_id")),
		UserId:         r.PostForm.Get("user_id"),
		UserName:       r.PostForm.Get("user_name"),
	}

	t := time.Now()
	y := t.Year()

	p, err := h.repo.GetParticipantById(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, req.UserId, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	if !p.IsHost {
		err := errors.New("You are not the host of this secret santa party, hence cannot randomize pairs")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	pCount, err := h.repo.CountAllParticipants(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	if pCount < 1 {
		err := errors.New("Secret Santa " + strconv.Itoa(y) + " has not been initialized yet for this Slack channel")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	mCount, err := h.repo.CountMatchedParticipants(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	if mCount > 0 && mCount == pCount {
		err := errors.New("Secret Santa " + strconv.Itoa(y) + " pairs for this Slack channel have already been matched")
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	poolA, err := h.repo.GetUnmatchedParticipants(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	poolB := make([]Participant, len(poolA))
	copy(poolB, poolA)
	rand.Seed(time.Now().UnixNano())
	for len(poolA) > 0 {
		participantA := poolA[0]
		if !participantA.IsMatched && participantA.YourMatchId == nil {
			i := rand.Intn(len(poolB))
			participantB := poolB[i]
			if participantA == participantB && len(poolA) != 1 {
				continue
			}
			err := h.repo.UpdateParticipantMatch(r.Context(), &participantB, &participantA, y)
			if err != nil {
				h.logger.Println(err)
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
				return
			}

			poolA = append(poolA[:0], poolA[1:]...)
			poolB = append(poolB[:i], poolB[i+1:]...)
		}
	}

	matchedParticipants, err := h.repo.GetAllParticipants(r.Context(), req.ChannelId, req.EnterpriseId, req.TeamId, y)
	if err != nil {
		h.logger.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeEphemeral, err.Error()})
		return
	}

	// Slack
	for _, participant := range matchedParticipants {
		yourMatchId := ""
		if participant.YourMatchId != nil {
			yourMatchId = *participant.YourMatchId
		}
		yourMatchAddress := ""
		if participant.YourMatchAddress != nil {
			yourMatchAddress = *participant.YourMatchAddress
		}
		msg := "Your match is <@" + yourMatchId + ">. Prepare your gift and send it to " + yourMatchAddress + ". Thank you and happy New Year!"
		err = SendSlackMessage(participant.ResponseUrl, ResponseTypeEphemeral, msg)
		if err != nil {
			h.logger.Println(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&SlackMessage{ResponseTypeInChannel, "Secret Santa pairs have been randomized!"})
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		start := time.Now()
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
			logger.Println("time", start, "method", r.Method, "path", r.URL.EscapedPath(), "duration", time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}
