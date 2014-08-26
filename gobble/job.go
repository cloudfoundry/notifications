package gobble

import "encoding/json"

type Job struct {
    ID       int    `db:"id"`
    WorkerID string `db:"worker_id"`
    Payload  string `db:"payload"`
    Version  int64  `db:"version"`
}

func NewJob(data interface{}) Job {
    payload, err := json.Marshal(data)
    if err != nil {
        panic(err)
    }

    return Job{
        Payload: string(payload),
    }
}

func (job Job) Unmarshal(v interface{}) error {
    return json.Unmarshal([]byte(job.Payload), v)
}
