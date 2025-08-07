package controllers

import (
	"encoding/json"
	"labyrinth/internal/controllers/middleware"
	"labyrinth/internal/protocol"
	"labyrinth/internal/router"
	"log/slog"
	"net/http"
	"time"
	"sync"
	"os"
	"fmt"

	_ "labyrinth/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/google/uuid"
)

type TeamQueueRequest struct {
    ID          string
    Type        string // "join", "leave", "create"
    UserID      uuid.UUID
    TeamID      string
    UserEmail   string
    TeamName    string
    Handler     func() error
    CreatedAt   time.Time
    RetryCount  int
    MaxRetries  int
}

var (
    teamQueue       []TeamQueueRequest
    teamMutex       sync.Mutex
    processingTeam  bool
    queueLogger     *slog.Logger
)

func InitTeamQueue(logger *slog.Logger) {
    queueLogger = logger
    queueLogger.Info("Team queue initialized")
}


func processTeamQueue() {
    teamMutex.Lock()
    if processingTeam {
        teamMutex.Unlock()
        return
    }
    processingTeam = true
    teamMutex.Unlock()

    if queueLogger != nil {
        queueLogger.Info("Starting team queue processing")
    }

    defer func() {
        teamMutex.Lock()
        processingTeam = false
        teamMutex.Unlock()
        
        if queueLogger != nil {
            queueLogger.Info("Team queue processing finished")
        }
    }()

    for {
        teamMutex.Lock()
        if len(teamQueue) == 0 {
            teamMutex.Unlock()
            break
        }

        req := teamQueue[0]
        teamQueue = teamQueue[1:]
        remaining := len(teamQueue)
        teamMutex.Unlock()

        if queueLogger != nil {
            queueLogger.Info("Processing team request", 
                "request_id", req.ID,
                "type", req.Type,
                "retry_count", req.RetryCount,
                "remaining_in_queue", remaining)
        }

        err := req.Handler()
        
        if err != nil {
            if queueLogger != nil {
                queueLogger.Error("Team request failed", 
                    "request_id", req.ID,
                    "error", err,
                    "retry_count", req.RetryCount)
            }
            
            // Retry logic
            if req.RetryCount < req.MaxRetries {
                req.RetryCount++
                
                if queueLogger != nil {
                    queueLogger.Info("Retrying team request", 
                        "request_id", req.ID,
                        "retry_count", req.RetryCount)
                }
                
                time.Sleep(time.Duration(req.RetryCount) * time.Second)
                
                teamMutex.Lock()
                teamQueue = append(teamQueue, req)
                teamMutex.Unlock()
            } else {
                logFailedRequest(req, err)
            }
        } else {
            if queueLogger != nil {
                queueLogger.Info("Team request completed successfully", 
                    "request_id", req.ID,
                    "type", req.Type)
            }
        }
    }
}



func addToTeamQueue(queueReq TeamQueueRequest) {
    teamMutex.Lock()
    defer teamMutex.Unlock()
    
    queueReq.ID = uuid.New().String()
    queueReq.CreatedAt = time.Now()
    if queueReq.MaxRetries == 0 {
        queueReq.MaxRetries = 3
    }
    
    teamQueue = append(teamQueue, queueReq)
    
    if queueLogger != nil {
        queueLogger.Info("Added to team queue", 
            "request_id", queueReq.ID,
            "type", queueReq.Type,
            "user_id", queueReq.UserID,
            "team_id", queueReq.TeamID,
            "queue_size", len(teamQueue))
    }
    
    go processTeamQueue()
}

func logFailedRequest(req TeamQueueRequest, err error) {
    if queueLogger != nil {
        queueLogger.Error("Team request failed permanently", 
            "request_id", req.ID,
            "type", req.Type,
            "user_id", req.UserID,
            "team_id", req.TeamID,
            "error", err,
            "retries", req.RetryCount)
    }
    
    file, fileErr := os.OpenFile("failed_team_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if fileErr != nil {
        return
    }
    defer file.Close()
    
    logEntry := fmt.Sprintf("[%s] FAILED: ID=%s, Type=%s, UserID=%s, TeamID=%s, Error=%s, Retries=%d\n",
        time.Now().Format(time.RFC3339),
        req.ID,
        req.Type, 
        req.UserID.String(),
        req.TeamID,
        err.Error(),
        req.RetryCount)
    
    file.WriteString(logEntry)
}

func HandleAll(rtr *router.Router) {
	// GET Routes here
	rtr.HandleFunc("/api", Get(DefaultHandler(rtr)))
	rtr.HandleFunc("/api/team", middleware.Authorized(rtr, Get(GetTeamHandler(rtr))))
	rtr.HandleFunc("/api/game", Get(GameConfigHandler(rtr)))

	// POST Routes
	rtr.HandleFunc("/api/user/status", middleware.Authorized(rtr, Post(TeamMemberStatusUpdateHandler(rtr))))
	rtr.HandleFunc("/api/team/create", middleware.Authorized(rtr, Post(TeamCreationHandler(rtr))))
	rtr.HandleFunc("/api/team/update", middleware.Authorized(rtr, Post(TeamUpdateHandler(rtr))))
	rtr.HandleFunc("/api/team/leave", middleware.Authorized(rtr, Post(LeaveTeamHandler(rtr))))

	rtr.HandleFunc("/api/eventlistener", Get(TeamChannelEventHandler(rtr)))

	rtr.Handle("/swagger/", http.StripPrefix("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3100/swagger/doc.json"),
	)))
}

func DefaultHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, err := w.Write([]byte("/api is up")); err != nil {

			rtr.Logger.Error("error serving /api", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		}

	})

}

func TeamChannelEventHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//userEmail := r.Context().Value("email").(string)

		query := r.URL.Query()
		teamId := query.Get("team_id")

		//user, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		//if err != nil {
		//
		//	http.Error(w, "internal error occurred", http.StatusInternalServerError)
		//	rtr.Logger.Error("internal error", "error", err)
		//	return
		//}

		//team, err := rtr.State.DB.GetTeamByUserId(context.Background(), user.ID)

		//if err != nil {

		//	http.Error(w, "internal error occurred, is the user in a valid team?", http.StatusInternalServerError)
		//	return
		//}

		teamChannel := rtr.State.ChanPool.GetChannel(teamId)
		listenerChannel := make(chan protocol.Packet)
		teamChannel.AddMember(listenerChannel)

		//w.Header().Add("Content-Type", "text/event-stream")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {

			http.Error(w, "Could not create flusher", http.StatusInternalServerError)
			return

		}

		//for eventMessage := range listenerChannel {

		//	if err := json.NewEncoder(w).Encode(eventMessage); err != nil {

		//		rtr.Logger.Debug("http stream write failed")

		//	}
		//	flusher.Flush()
		//}

		for eventMessage := range listenerChannel {
			// Convert event message to JSON string
			data, err := json.Marshal(eventMessage)
			if err != nil {
				rtr.Logger.Debug("failed to marshal event message", slog.String("error", err.Error()))
				continue
			}

			// Construct SSE message
			sseMessage := "data: " + string(data) + "\n\n"

			// Write the SSE message to the response
			_, err = w.Write([]byte(sseMessage))
			if err != nil {
				rtr.Logger.Debug("http stream write failed", slog.String("error", err.Error()))
				return
			}

			// Flush the output to ensure immediate delivery of the message
			flusher.Flush()
		}

	})

}
