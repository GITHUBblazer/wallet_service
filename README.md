wallet-service/
|-- cmd/
|   |-- main.go
|-- internal/
|   |-- api/
|   |   |-- api.go
|   |   |-- wallet_api.go
|   |-- config/
|   |   |-- config.go
|   |-- database/
|   |   |-- database.go
|   |-- logger/
|   |   |-- log.go
|   |-- model/
|   |   |-- wallet.go
|   |-- repository/
|   |   |-- interface
|   |   |   |-- interface.go
|   |   |-- postgres
|   |   |   |-- postgres.go
|   |   |-- repository.go
|   |-- service/
|   |   |-- service.go
|   |   |-- wallet_service.go
|-- test/
|   |-- repository_test.go
|   |-- service_test.go
|--.env
|--.gitignore
|-- Dockerfile
|-- docker-compose.yml
|-- go.mod
|-- go.sum
|-- golangci.yaml
|-- README.md

main.go：项目的入口文件，负责初始化配置、数据库连接、日志记录等，然后启动 HTTP 服务器并注册路由。
2.2 internal目录
api目录
api.go：定义了 HTTP 路由和启动 HTTP 服务器的函数。
wallet_api.go：包含了处理各种 API 请求的处理器函数，如存款、取款、转账、查询余额和查询交易历史等。
config目录
config.go：用于读取和解析配置文件，提供配置信息给其他模块使用。
database目录
database.go：负责初始化和管理与 PostgreSQL 数据库的连接，提供数据库操作的基础方法。
logger目录
logger.go：实现了日志记录功能，提供不同级别的日志记录方法。
models目录
transaction.go：定义了交易记录的数据结构，包括交易 ID、交易类型、金额、时间等字段。
wallet.go：定义了钱包的数据结构，包括用户 ID、余额、最后更新时间等字段。
repository目录
repository.go：包含了与数据库交互的方法，如插入交易记录、更新钱包余额、查询钱包余额和交易历史等。
service目录
service.go：实现了钱包服务的业务逻辑，包括存款、取款、转账、查询余额和查询交易历史等功能，调用repository中的方法与数据库交互。
2.3 pkg目录
decimal目录
decimal.go：用于处理精确的十进制计算，确保金额计算的准确性，避免浮点数计算带来的精度问题。
2.4 test目录
e2e目录
e2e_test.go：进行端到端的 API 测试，模拟用户的实际操作，验证整个系统的功能是否正常。
unit目录
handlers_test.go：对handlers.go中的处理器函数进行单元测试，测试各个 API 端点的功能是否正确。
service_test.go：对service.go中的服务函数进行单元测试，测试业务逻辑的正确性。
2.5 其他文件
.gitignore：指定哪些文件或目录不需要被 Git 跟踪。
Dockerfile：用于构建项目的 Docker 镜像，定义了镜像的基础环境、依赖安装和项目的复制等操作。
docker-compose.yml：用于定义和运行多个容器化服务，包括 PostgreSQL 数据库和 Redis（如果需要），方便在本地进行开发和测试。
go.mod和go.sum：Go 语言的模块管理文件，记录项目的依赖关系和版本信息。
golangci.yaml：用于配置golangci-lint的检查规则，确保代码质量。
