package datatype

import "fin4-core/server/ethereum"

//ServiceContainer defines our service container type
type ServiceContainer struct {
	Config          Config
	AssetService    AssetService
	TimelineService TimelineService
	TokenService    TokenService
	UserService     UserService
	FileStorage     FileStorage
	Ethereum        *ethereum.Ethereum
}
