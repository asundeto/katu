package errorhandler

import (
	"errors"
)

var (
    ErrNoRecord           = errors.New("models: теңізілген жазба табылмады")
    ErrInvalidCredentials = errors.New("models: жарамсыз кілтсөздер")
    ErrDuplicateEntry     = errors.New("қайталау енгізу")
    ErrNoComments         = errors.New("models: постқа түсініктеме таппалмады")
    ErrZeroCode           = errors.New("Уақытша таңдалған жеткізу")
    ErrGoogleInfo         = errors.New("Google-дан алынған ақпарат")
    ErrGitHUBInfo         = errors.New("GitHub-тан алынған ақпарат")
    ErrServerError        = errors.New("Сервер қатесі")
    ErrComment            = errors.New("Дұрыс мәнді енгізіңіз")
    ErrIncorrectValue     = errors.New("Дұрыс пішім!")
    ErrPostTitle          = errors.New("Тақырып ұзындығы 3-тен 40 символға дейін болуы керек!")
    ErrPostContent        = errors.New("Мазмұн ұзындығы 5-тен 1000 символға дейін болуы керек!")
    ErrPostImageExtension = errors.New("Суреттің форматы: .png .jpg .gif")
    ErrPostImageSize      = errors.New("Сурет өлшемі 20mb-дан асқан!")
    ErrUploadImage        = errors.New("Дұрыс суретті жүктеу!")
)

var (
    ErrTooManyRequests = errors.New("Көп сұрау")
)
