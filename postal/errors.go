package postal

type CCDownError string

func (err CCDownError) Error() string {
    return string(err)
}

type UAADownError string

func (err UAADownError) Error() string {
    return string(err)
}

type UAAUserNotFoundError string

func (err UAAUserNotFoundError) Error() string {
    return string(err)
}

type UAAGenericError string

func (err UAAGenericError) Error() string {
    return string(err)
}

type TemplateLoadError string

func (err TemplateLoadError) Error() string {
    return string(err)
}
