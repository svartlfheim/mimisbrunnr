package cmd

// var dummyConfig *config.AppConfig = &config.AppConfig{
// 	HTTP: config.HTTPConfig{
// 		Port:       "8080",
// 		ListenHost: "0.0.0.0",
// 	},
// 	RDB: config.RDBConfig{
// 		Driver:   "postgres",
// 		Host:     "localhost",
// 		Port:     "5432",
// 		Username: "dummy",
// 		Password: "dummy",
// 		Schema:   "dummy",
// 		Database: "dummy",
// 	},
// }

// func Test_DI_GetRootCommandRegistry(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetRootCommandRegistry()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetRootCommandRegistry())
// }

// func Test_DI_GetRDBConnManager(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetRDBConnManager()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetRDBConnManager())
// }

// func Test_DI_GetRDBConnManager_panic(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	// It panics cause the conn manager requires database creds
// 	di := NewDIContainer(l.Logger, fs, &config.AppConfig{})

// 	assert.Panics(t, func() {
// 		di.GetRDBConnManager()
// 	})
// }

// func Test_DI_GetPostgresSCMIntegrationsRepository(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetPostgresSCMIntegrationsRepository()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetPostgresSCMIntegrationsRepository())
// }

// func Test_DI_GetSCMIntegrationsManager(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetSCMIntegrationsManager()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetSCMIntegrationsManager())
// }

// func Test_DI_GetErrorHandlingJsonUnmarshaller(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetErrorHandlingJsonUnmarshaller()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetErrorHandlingJsonUnmarshaller())
// }

// func Test_DI_GetSCMIntegrationsController(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetSCMIntegrationsController()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetSCMIntegrationsController())
// }

// func Test_DI_GetProjectsController(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetProjectsController()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetProjectsController())
// }

// func Test_DI_GetServer(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetServer()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetServer())
// }

// func Test_DI_GetValidator(t *testing.T) {
// 	l := zerologmocks.NewLogger()
// 	fs := afero.NewMemMapFs()
// 	di := NewDIContainer(l.Logger, fs, dummyConfig)

// 	// As long as it doesn't panic we're good
// 	v := di.GetValidator()

// 	// ensure it's a singleton
// 	assert.Same(t, v, di.GetValidator())
// }
