package queue

import (
    "fmt"
    "log/slog"
    "os"
    "sync"
    "time"

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

type TeamQueue struct {
    queue      []TeamQueueRequest
    mutex      sync.Mutex
    processing bool
    logger     *slog.Logger
}

var globalTeamQueue *TeamQueue

func InitTeamQueue(logger *slog.Logger) {
    globalTeamQueue = &TeamQueue{
        queue:  make([]TeamQueueRequest, 0),
        logger: logger,
    }
    
    if logger != nil {
        logger.Info("Team queue initialized")
    }
}

func AddToQueue(queueReq TeamQueueRequest) {
    if globalTeamQueue == nil {
        panic("Team queue not initialized. Call InitTeamQueue first.")
    }
    
    globalTeamQueue.add(queueReq)
}

func (tq *TeamQueue) add(queueReq TeamQueueRequest) {
    tq.mutex.Lock()
    defer tq.mutex.Unlock()
    
    queueReq.ID = uuid.New().String()
    queueReq.CreatedAt = time.Now()
    if queueReq.MaxRetries == 0 {
        queueReq.MaxRetries = 3
    }
    
    tq.queue = append(tq.queue, queueReq)
    
    if tq.logger != nil {
        tq.logger.Info("Added to team queue", 
            "request_id", queueReq.ID,
            "type", queueReq.Type,
            "user_id", queueReq.UserID,
            "team_id", queueReq.TeamID,
            "queue_size", len(tq.queue))
    }
    
    go tq.process()
}

func (tq *TeamQueue) process() {
    tq.mutex.Lock()
    if tq.processing {
        tq.mutex.Unlock()
        return
    }
    tq.processing = true
    tq.mutex.Unlock()

    if tq.logger != nil {
        tq.logger.Info("Starting team queue processing")
    }

    defer func() {
        tq.mutex.Lock()
        tq.processing = false
        tq.mutex.Unlock()
        
        if tq.logger != nil {
            tq.logger.Info("Team queue processing finished")
        }
    }()

    for {
        tq.mutex.Lock()
        if len(tq.queue) == 0 {
            tq.mutex.Unlock()
            break
        }

        req := tq.queue[0]
        tq.queue = tq.queue[1:]
        remaining := len(tq.queue)
        tq.mutex.Unlock()

        if tq.logger != nil {
            tq.logger.Info("Processing team request", 
                "request_id", req.ID,
                "type", req.Type,
                "retry_count", req.RetryCount,
                "remaining_in_queue", remaining)
        }

        err := req.Handler()
        
        if err != nil {
            tq.handleError(req, err)
        } else {
            if tq.logger != nil {
                tq.logger.Info("Team request completed successfully", 
                    "request_id", req.ID,
                    "type", req.Type)
            }
        }
    }
}

func (tq *TeamQueue) handleError(req TeamQueueRequest, err error) {
    if tq.logger != nil {
        tq.logger.Error("Team request failed", 
            "request_id", req.ID,
            "error", err,
            "retry_count", req.RetryCount)
    }
    
    if req.RetryCount < req.MaxRetries {
        req.RetryCount++
        
        if tq.logger != nil {
            tq.logger.Info("Retrying team request", 
                "request_id", req.ID,
                "retry_count", req.RetryCount)
        }
        
        time.Sleep(time.Duration(req.RetryCount) * time.Second)
        
        tq.mutex.Lock()
        tq.queue = append(tq.queue, req)
        tq.mutex.Unlock()
    } else {
        tq.logFailedRequest(req, err)
    }
}

func (tq *TeamQueue) logFailedRequest(req TeamQueueRequest, err error) {
    if tq.logger != nil {
        tq.logger.Error("Team request failed permanently", 
            "request_id", req.ID,
            "type", req.Type,
            "user_id", req.UserID,
            "team_id", req.TeamID,
            "error", err,
            "retries", req.RetryCount)
    }
    
    file, fileErr := os.OpenFile("failed_team_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Log to file
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

func GetQueueStats() map[string]interface{} {
    if globalTeamQueue == nil {
        return map[string]interface{}{"error": "queue not initialized"}
    }
    
    globalTeamQueue.mutex.Lock()
    defer globalTeamQueue.mutex.Unlock()
    
    return map[string]interface{}{
        "queue_size":  len(globalTeamQueue.queue),
        "processing":  globalTeamQueue.processing,
    }
}