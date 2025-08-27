package queue

import (
    "fmt"
    "log/slog"
    "os"
    "sync"
    "time"
)

type DBQueueRequest struct {
    ID         string
    Type       string // "team_create", "user_update", bla bla bla
    Payload    interface{}
    Handler    func() error
    CreatedAt  time.Time
    RetryCount int
    MaxRetries int
}

type DBQueue struct {
    queue      []DBQueueRequest
    mutex      sync.Mutex
    processing bool
    logger     *slog.Logger
}

var globalDBQueue *DBQueue

func InitDBQueue(logger *slog.Logger) {
    globalDBQueue = &DBQueue{
        queue:  make([]DBQueueRequest, 0),
        logger: logger,
    }
    if logger != nil {
        logger.Info("DB queue initialized")
    }
}

func AddToQueue(req DBQueueRequest) {
    if globalDBQueue == nil {
        panic("DB queue not initialized. Call InitDBQueue first.")
    }
    globalDBQueue.add(req)
}

func (q *DBQueue) add(req DBQueueRequest) {
    q.mutex.Lock()
    defer q.mutex.Unlock()

    req.ID = fmt.Sprintf("%d", time.Now().UnixNano())
    req.CreatedAt = time.Now()
    if req.MaxRetries == 0 {
        req.MaxRetries = 3
    }

    q.queue = append(q.queue, req)

    if q.logger != nil {
        q.logger.Info("Added to DB queue",
            "request_id", req.ID,
            "type", req.Type,
            "queue_size", len(q.queue))
    }

    go q.process()
}

func (q *DBQueue) process() {
    q.mutex.Lock()
    if q.processing {
        q.mutex.Unlock()
        return
    }
    q.processing = true
    q.mutex.Unlock()

    if q.logger != nil {
        q.logger.Info("Starting DB queue processing")
    }

    defer func() {
        q.mutex.Lock()
        q.processing = false
        q.mutex.Unlock()
        if q.logger != nil {
            q.logger.Info("DB queue processing finished")
        }
    }()

    for {
        q.mutex.Lock()
        if len(q.queue) == 0 {
            q.mutex.Unlock()
            break
        }

        req := q.queue[0]
        q.queue = q.queue[1:]
        remaining := len(q.queue)
        q.mutex.Unlock()

        if q.logger != nil {
            q.logger.Info("Processing DB request",
                "request_id", req.ID,
                "type", req.Type,
                "retry_count", req.RetryCount,
                "remaining_in_queue", remaining)
        }

        err := req.Handler()
        if err != nil {
            q.handleError(req, err)
        } else {
            if q.logger != nil {
                q.logger.Info("DB request completed successfully",
                    "request_id", req.ID,
                    "type", req.Type)
            }
        }
    }
}

func (q *DBQueue) handleError(req DBQueueRequest, err error) {
    if q.logger != nil {
        q.logger.Error("DB request failed",
            "request_id", req.ID,
            "error", err,
            "retry_count", req.RetryCount)
    }

    if req.RetryCount < req.MaxRetries {
        req.RetryCount++
        if q.logger != nil {
            q.logger.Info("Retrying DB request",
                "request_id", req.ID,
                "retry_count", req.RetryCount)
        }
        time.Sleep(time.Duration(req.RetryCount) * time.Second)
        q.mutex.Lock()
        q.queue = append(q.queue, req)
        q.mutex.Unlock()
    } else {
        q.logFailedRequest(req, err)
    }
}

func (q *DBQueue) logFailedRequest(req DBQueueRequest, err error) {
    if q.logger != nil {
        q.logger.Error("DB request failed permanently",
            "request_id", req.ID,
            "type", req.Type,
            "error", err,
            "retries", req.RetryCount)
    }

    file, fileErr := os.OpenFile("failed_db_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if fileErr != nil {
        return
    }
    defer file.Close()

    logEntry := fmt.Sprintf("[%s] FAILED: ID=%s, Type=%s, Error=%s, Retries=%d\n",
        time.Now().Format(time.RFC3339),
        req.ID,
        req.Type,
        err.Error(),
        req.RetryCount)

    file.WriteString(logEntry)
}

func GetDBQueueStats() map[string]interface{} {
    if globalDBQueue == nil {
        return map[string]interface{}{"error": "queue not initialized"}
    }
    globalDBQueue.mutex.Lock()
    defer globalDBQueue.mutex.Unlock()
    return map[string]interface{}{
        "queue_size":  len(globalDBQueue.queue),
        "processing":  globalDBQueue.processing,
    }
}