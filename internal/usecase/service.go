package usecase

type Service struct {
	MovieProvider MovieProvider
}

func New(provider MovieProvider) *Service {
	return &Service{
		MovieProvider: provider,
	}
}
