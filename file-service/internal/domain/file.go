package domain

import "time"

type File struct {
    ID              string    `json:"id"`
    Name            string    `json:"name"`
    Size            int64     `json:"size"`
    CID             string    `json:"cid"`
    UploadedAt      time.Time `json:"uploadedAt"`
    DownloadKeyword string    `json:"downloadKeyword"`
    DeleteKeyword   string    `json:"deleteKeyword"`
}
