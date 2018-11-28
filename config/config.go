package config

const (
    // ParseCityList name
    ParseCityList = "ParseCityList"

    // ParseCity name
    ParseCity = "ParseCity"

    // ProfileParser name
    ProfileParser = "ProfileParser"

    // NilParser name
    NilParser = "NilParser"

    // ElasticIndex ElasticSearch存储的index
    ElasticIndex = "dating_profile"

    // ItemServerRPC Endpoints
    ItemServerRPC = "ItemSaverService.Save"

    // CrawlServiceRPC Endpoints
    CrawlServiceRPC = "CrawlService.Process"

    // CrawlServiceRPCPort rpc port
    CrawlServiceRPCPort = 9987

    //ItemServerRPCPort port
    ItemServerRPCPort = 1234

    // QPS rate limiting
    QPS = 10

    // UpdateClientPoolSecond update client pool
    UpdateClientPoolSecond = 3
)
