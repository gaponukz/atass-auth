package architecture

func domainLayer() []string {
	return []string{
		"auth/src/domain/entities",
		"auth/src/domain/errors",
	}
}

func applicationLayer() []string {
	return []string{
		"auth/src/application/dto",
		"auth/src/application/usecases/passreset",
		"auth/src/application/usecases/routes",
		"auth/src/application/usecases/session",
		"auth/src/application/usecases/settings",
		"auth/src/application/usecases/show_routes",
		"auth/src/application/usecases/signin",
		"auth/src/application/usecases/signup",
	}
}

func infrastructureLayer() []string {
	return []string{
		"auth/src/infrastructure/config",
		"auth/src/infrastructure/logger",
		"auth/src/infrastructure/notifier",
		"auth/src/infrastructure/security",
		"auth/src/infrastructure/storage",
	}
}

func interfaceLayer() []string {
	return []string{
		"controller",
		"event_handler",
	}
}
