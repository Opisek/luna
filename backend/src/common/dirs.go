package common

func (env *Environmental) getBasePath() string {
	return env.DATA_PATH
}

func (env *Environmental) GetKeysPath() string {
	return env.getBasePath() + "/keys"
}
