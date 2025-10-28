package stargo

func (s *App) LogWarnf(format string, v ...any) {
	s.logger.Warnf(format, v...)

}
func (s *App) LogDebugf(format string, v ...any) {
	s.logger.Debugf(format, v...)
}
func (s *App) LogErrorf(format string, v ...any) {
	s.logger.Errorf(format, v...)
}
func (s *App) LogFatalf(format string, v ...any) {
	s.logger.Fatalf(format, v...)
}
func (s *App) LogInfof(format string, v ...any) {
	s.logger.Infof(format, v...)
}
