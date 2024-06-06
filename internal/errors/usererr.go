package errorhandler

import "errors"

var (
    ErrUserAgrement = errors.New("Пайдалы ақпаратты қабылдаңыз!")
    ErrShortUsername = errors.New("Пайдаланушы аты өте қысқа!")
    ErrLongUsername = errors.New("Пайдаланушы аты өте ұзын!")
    ErrUsernameStart = errors.New("Пайдаланушы аты санмен басталмауы керек!")

    ErrLongUsernameSymbols = errors.New("Пайдаланушы аты 10 символдан үлкен болмауы керек")
    ErrLongEmailSymbols = errors.New("Электрондық пошта 30 символдан үлкен болмауы керек")
    ErrQuery = errors.New("сұрау қатесі")
    ErrAlreadyExistUsername = errors.New("Пайдаланушы аты қазірдің бірінде бар")
    ErrAlreadyExistEmail = errors.New("Электрондық пошта кіріс")

    ErrEnterCorrectEmail = errors.New("Дұрыс электрондық поштаны енгізіңіз!")
    ErrPasswordMismatch = errors.New("Парольдер сәйкес келмейді!")
    ErrLowPassword = errors.New("Парольдер міндетті екі үлкен сандардан 7 символдан көп болуы керек [$#%!?.*]")
    ErrAuthServer = errors.New("Сервер қатесі! Кейінірек көріңіз")

    ErrEmailExist = errors.New("Электрондық пошта жоқ!")
    ErrIncorrectPassword = errors.New("Пароль дұрыс емес!")
)
