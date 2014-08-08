package models

import (
    "errors"
)

var (
    ErrRecordNotFound  = errors.New("Record Not Found")
    ErrDuplicateRecord = errors.New("Duplicate Record")
)
